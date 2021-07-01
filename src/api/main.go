package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"saltgram/api/handlers"
	"saltgram/internal"
	"saltgram/log"
	"saltgram/pki"
	"saltgram/protos/admin/pradmin"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/content/prcontent"
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"
	"time"

	"github.com/rs/cors"
)

func main() {
	l := log.NewLogger("saltgram-api")
	l.L.Printf("Starting API Gateway on port: %s\n", os.Getenv("SALT_API_PORT"))
	pkiHandler := pki.Init()
	cert, err := pkiHandler.RegisterSaltgramService("api-gateway")
	if err != nil {
		l.L.Fatalf("registering service for pki: %v\n", err)
	}
	s := internal.NewService(l.L)

	err = s.Init("saltgram-api-gateway", cert.CertPEM, cert.PrivateKeyPEM, pkiHandler.RootCA.CertPEM)
	if err != nil {
		l.L.Fatalf("failed to init api service: %v\n", err)
	}

	authConnection, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_AUTH_ADDR", "localhost"), os.Getenv("SALT_AUTH_PORT")))
	if err != nil {
		l.L.Fatalf("dialing auth connection: %v\n", err)
	}
	defer authConnection.Close()
	authClient := prauth.NewAuthClient(authConnection)
	authHandler := handlers.NewAuth(l.L, authClient)
	authRouter := s.S.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/refresh", authHandler.Refresh) //.Methods(http.MethodGet)
	authRouter.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
	authRouter.HandleFunc("/jwt", authHandler.GetJWT).Methods(http.MethodPost)
	authRouter.HandleFunc("/update", authHandler.UpdateJWTUsername).Methods(http.MethodPut)
	authRouter.HandleFunc("", authHandler.CheckPermissions).Methods(http.MethodPut)
	authRouter.HandleFunc("/2fa", authHandler.Get2FAQR).Methods(http.MethodPost)
	authRouter.HandleFunc("/2fa/{token}", authHandler.Authenticate2FA).Methods(http.MethodGet)

	usersConnection, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_USERS_ADDR", "localhost"), os.Getenv("SALT_USERS_PORT")))
	if err != nil {
		l.L.Fatalf("dialing users connection: %v\n", err)
	}
	defer usersConnection.Close()
	usersClient := prusers.NewUsersClient(usersConnection)
	usersHandler := handlers.NewUsers(l.L, usersClient)
	usersRouter := s.S.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/register", usersHandler.Register).Methods(http.MethodPost)
	usersRouter.HandleFunc("", usersHandler.GetByJWS).Methods(http.MethodGet)
	usersRouter.HandleFunc("/{username}", usersHandler.GetByUsernameRoute).Methods(http.MethodGet)
	usersRouter.HandleFunc("/resetpass", usersHandler.ResetPassword).Methods(http.MethodPost)
	usersRouter.HandleFunc("/changepass", usersHandler.ChangePassword).Methods(http.MethodPost)
	usersRouter.HandleFunc("/profile/{username}", usersHandler.GetProfile).Methods(http.MethodGet)
	usersRouter.HandleFunc("/profile/{username}", usersHandler.UpdateProfile).Methods(http.MethodPut)
	usersRouter.HandleFunc("/create/follow", usersHandler.Follow).Methods(http.MethodPost)
	usersRouter.HandleFunc("/unfollow", usersHandler.Unfollow).Methods(http.MethodPost)
	usersRouter.HandleFunc("/get/followers/{username}", usersHandler.GetFollowers).Methods(http.MethodGet)
	usersRouter.HandleFunc("/get/following/{username}", usersHandler.GetFollowing).Methods(http.MethodGet)
	usersRouter.HandleFunc("/search/{username}", usersHandler.SearchUsers).Methods(http.MethodGet)
	usersRouter.HandleFunc("/follow/requests/", usersHandler.GetFollowingRequest).Methods(http.MethodGet)
	usersRouter.HandleFunc("/follow/request/", usersHandler.FollowRespond).Methods(http.MethodPost)
	usersRouter.HandleFunc("/followers/detailed/{username}", usersHandler.GetFollowersDetailed).Methods(http.MethodGet)
	usersRouter.HandleFunc("/following/detailed/{username}", usersHandler.GetFollowingDetailed).Methods(http.MethodGet)
	usersRouter.HandleFunc("/check/follow/{username}", usersHandler.CheckFollowing).Methods(http.MethodGet)
	usersRouter.HandleFunc("/check/followrequest/{username}", usersHandler.CheckFollowRequest).Methods(http.MethodGet)
	usersRouter.HandleFunc("/mute/profile", usersHandler.MuteProfile).Methods(http.MethodPost)
	usersRouter.HandleFunc("/unmute/profile", usersHandler.UnmuteProfile).Methods(http.MethodPost)
	usersRouter.HandleFunc("/get/muted", usersHandler.GetMutedProfiles).Methods(http.MethodGet)
	usersRouter.HandleFunc("/check/muted/{username}", usersHandler.CheckIfMuted).Methods(http.MethodGet)
	usersRouter.HandleFunc("/block/profile", usersHandler.BlockProfile).Methods(http.MethodPost)
	usersRouter.HandleFunc("/unblock/profile", usersHandler.UnblockProfile).Methods(http.MethodPost)
	usersRouter.HandleFunc("/get/blocked", usersHandler.GetBlockedProfiles).Methods(http.MethodGet)
	usersRouter.HandleFunc("/check/blocked/{username}", usersHandler.CheckIfBlocked).Methods(http.MethodGet)
	usersRouter.HandleFunc("/add/closefrined", usersHandler.AddCloseFriend).Methods(http.MethodPost)
	usersRouter.HandleFunc("/remove/closefrined", usersHandler.RemoveCloseFriend).Methods(http.MethodPost)
	usersRouter.HandleFunc("/get/closefriend", usersHandler.GetCloseFriends).Methods(http.MethodGet)
	usersRouter.HandleFunc("/get/closefriend/following", usersHandler.GetProfilesForCloseFriends).Methods(http.MethodGet)
	usersRouter.HandleFunc("/taggableprofiles/get", usersHandler.GetTaggableProfiles).Methods(http.MethodGet)

	emailConnection, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_EMAIL_ADDR", "localhost"), os.Getenv("SALT_EMAIL_PORT")))
	if err != nil {
		l.L.Fatalf("dialing email connection")
	}

	emailClient := premail.NewEmailClient(emailConnection)
	emailHandler := handlers.NewEmail(l.L, emailClient, usersClient)
	emailRouter := s.S.PathPrefix("/email").Subrouter()
	emailRouter.HandleFunc("/activate/{token}", emailHandler.Activate).Methods(http.MethodPut)
	emailRouter.HandleFunc("/forgot", emailHandler.ForgotPassword).Methods(http.MethodPost)
	emailRouter.HandleFunc("/reset/{token}", emailHandler.ConfirmReset).Methods(http.MethodPut)

	contentConnection, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_CONTENT_ADDR", "localhost"), os.Getenv("SALT_CONTENT_PORT")))
	if err != nil {
		l.L.Fatalf("dialing content connection: %v\n", err)
	}
	defer contentConnection.Close()
	contentClient := prcontent.NewContentClient(contentConnection)
	contentHandler := handlers.NewContent(l.L, contentClient, usersClient)
	contentRouter := s.S.PathPrefix("/content").Subrouter()
	contentRouter.HandleFunc("/user", contentHandler.GetSharedMedia).Methods(http.MethodGet)
	// contentRouter.HandleFunc("/sharedmedia", contentHandler.AddSharedMedia).Methods(http.MethodPost) What?
	contentRouter.HandleFunc("/user/{id}", contentHandler.GetSharedMediaByUser).Methods(http.MethodGet)
	contentRouter.HandleFunc("/profilepicture/{id}", contentHandler.GetProfilePictureByUser).Methods(http.MethodGet)
	contentRouter.HandleFunc("/profilepicture", contentHandler.AddProfilePicture).Methods(http.MethodPost) // Mora se impl
	contentRouter.HandleFunc("/post/{id}", contentHandler.GetPostsByUser).Methods(http.MethodGet)
	contentRouter.HandleFunc("/post", contentHandler.AddPost).Methods(http.MethodPost)
	contentRouter.HandleFunc("/story", contentHandler.AddStory).Methods(http.MethodPost)
	contentRouter.HandleFunc("/story/{id}", contentHandler.GetStoriesByUser).Methods(http.MethodGet)
	contentRouter.HandleFunc("/highlight", contentHandler.AddHighlight).Methods(http.MethodPost)
	contentRouter.HandleFunc("/highlight/{id}", contentHandler.GetHighlights).Methods(http.MethodGet)
	contentRouter.HandleFunc("/comment", contentHandler.AddComment).Methods(http.MethodPost)
	contentRouter.HandleFunc("/reaction", contentHandler.AddReaction).Methods(http.MethodPost)
	contentRouter.HandleFunc("/reaction/user", contentHandler.GetPostsByUserReaction).Methods(http.MethodGet)
	contentRouter.HandleFunc("/comment/{id}", contentHandler.GetCommentsByPost).Methods(http.MethodGet)
	contentRouter.HandleFunc("/reaction/{id}", contentHandler.GetReactionsByPost).Methods(http.MethodGet)
	contentRouter.HandleFunc("/story/{id}", contentHandler.GetStoriesByUser).Methods(http.MethodGet)
	contentRouter.HandleFunc("/reaction", contentHandler.PutReaction).Methods(http.MethodPut)
	contentRouter.HandleFunc("/tag/{value}", contentHandler.GetPostsByTag).Methods(http.MethodGet)
	contentRouter.HandleFunc("/tag/search/{value}", contentHandler.GetTagsByName).Methods(http.MethodGet)
	contentRouter.HandleFunc("/location/search/{name}", contentHandler.GetLocationNames).Methods(http.MethodGet)
	contentRouter.HandleFunc("/location/{name}", contentHandler.GetContentsByLocation).Methods(http.MethodGet)

	adminConnection, err := s.GetConnection(fmt.Sprintf("%s:%s", internal.GetEnvOrDefault("SALT_ADMIN_ADDR", "localhost"), os.Getenv("SALT_ADMIN_PORT")))
	if err != nil {
		l.L.Fatalf("dialing admin connection: %v\n", err)
	}
	defer adminConnection.Close()
	adminClient := pradmin.NewAdminClient(adminConnection)
	adminHandler := handlers.NewAdmin(l.L, adminClient, usersClient, contentClient)
	adminRouter := s.S.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/verificationrequest", adminHandler.GetPendingVerifications).Methods(http.MethodGet)
	// Better suited for user router?
	adminRouter.HandleFunc("/verificationrequest", adminHandler.AddVerificationRequest).Methods(http.MethodPost)
	adminRouter.HandleFunc("/verificationrequest", adminHandler.ReviewVerificationRequest).Methods(http.MethodPut)
	// Better suited for user router?
	adminRouter.HandleFunc("/inappropriatecontent", adminHandler.SendInappropriateContentReport).Methods(http.MethodPost)
	adminRouter.HandleFunc("/inappropriatecontent", adminHandler.GetPendingReports).Methods(http.MethodGet)

	// TODO REPAIR THIS AFTER FINISHING FRONTEND
	c := cors.New(cors.Options{
		AllowedOrigins: []string{fmt.Sprintf("https://localhost:%s", os.Getenv("SALT_WEB_PORT"))},
		//AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowCredentials: true,
		Debug:            true,
	})

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("SALT_API_PORT")),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      c.Handler(s.S),
		TLSConfig:    s.TLS.TC,
	}

	go func() {
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			l.L.Fatalf("while serving: %v\n", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)

	sig := <-signalChan
	l.L.Printf("Recieved terminate, graceful shutdown with sigtype: %v\n", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(tc)
}
