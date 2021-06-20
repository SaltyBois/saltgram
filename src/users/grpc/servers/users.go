package servers

import (
	"context"
	"os"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/content/prcontent"
	"saltgram/protos/email/premail"
	"saltgram/protos/users/prusers"
	"saltgram/users/data"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Users struct {
	prusers.UnimplementedUsersServer
	l  *logrus.Logger
	db *data.DBConn
	ac prauth.AuthClient
	ec premail.EmailClient
	cc prcontent.ContentClient
}

func NewUsers(l *logrus.Logger, db *data.DBConn, ac prauth.AuthClient, ec premail.EmailClient, cc prcontent.ContentClient) *Users {
	return &Users{
		l:  l,
		db: db,
		ac: ac,
		ec: ec,
		cc: cc,
	}
}

func (u *Users) UpdateProfilePicture(ctx context.Context, r *prusers.UpdateProfilePictureRequest) (*prusers.UpdateProfilePictureResponse, error) {
	err := u.db.UpdateProfilePicture(r.Url, r.Username)
	if err != nil {
		u.l.Errorf("failed to update profile picture: %v", err)
		return &prusers.UpdateProfilePictureResponse{}, status.Error(codes.Internal, "Internal error")
	}
	return &prusers.UpdateProfilePictureResponse{}, nil
}

func (u *Users) ResetPassword(ctx context.Context, r *prusers.UserResetRequest) (*prusers.UserResetResponse, error) {
	err := data.ResetPassword(u.db, r.Email, r.Password)
	if err != nil {
		u.l.Errorf("failure resetting password: %v\n", err)
		return &prusers.UserResetResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}
	return &prusers.UserResetResponse{}, nil
}

func (u *Users) GetByUsername(ctx context.Context, r *prusers.GetByUsernameRequest) (*prusers.GetByUsernameResponse, error) {
	user, err := u.db.GetUserByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] username: %v\n", r.Username)
		u.l.Errorf("failure getting user by username: %v\n", err)
		return &prusers.GetByUsernameResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prusers.GetByUsernameResponse{
		Id:             user.ID,
		Email:          user.Email,
		FullName:       user.FullName,
		Username:       user.Username,
		Role:           user.Role,
		HashedPassword: user.HashedPassword,
	}, nil
}

func (u *Users) GetRole(ctx context.Context, r *prusers.RoleRequest) (*prusers.RoleResponse, error) {
	role, err := data.GetRole(u.db, r.Username)
	if err != nil {
		u.l.Errorf("failure getting role: %v\n", err)
		return &prusers.RoleResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prusers.RoleResponse{Role: role}, nil
}

func (u *Users) ChangePassword(ctx context.Context, r *prusers.ChangeRequest) (*prusers.ChangeResponse, error) {
	u.l.Infof("changing password for: %v\n", r.Username)
	err := data.ChangePassword(u.db, r.Username, r.OldPlainPassword, r.NewPlainPassword)
	if err != nil {
		u.l.Errorf("failure attepmting to change password: %v\n", err)
		return &prusers.ChangeResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}
	return &prusers.ChangeResponse{}, nil
}

func (u *Users) VerifyEmail(ctx context.Context, r *prusers.VerifyEmailRequest) (*prusers.VerifyEmailResponse, error) {
	u.l.Infof("Verifying email: %v\n", r.Email)

	err := data.VerifyEmail(u.db, r.Email)
	if err != nil {
		u.l.Errorf("failure verifying email: %v\n", err)
		return &prusers.VerifyEmailResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prusers.VerifyEmailResponse{}, nil
}

func (u *Users) Register(ctx context.Context, r *prusers.RegisterRequest) (*prusers.RegisterResponse, error) {
	u.l.Infof("registering %v\n", r.Email)

	refreshClaims := data.RefreshClaims{
		Username: r.Username,
		StandardClaims: jwt.StandardClaims{
			// TODO(Jovan): Make programmatic?
			ExpiresAt: time.Now().UTC().AddDate(0, 6, 0).Unix(),
			Issuer:    "SaltGram",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	jws, err := token.SignedString([]byte(os.Getenv("REF_SECRET_KEY")))

	if err != nil {
		u.l.Errorf("failure signing refresh token")
		return &prusers.RegisterResponse{}, status.Error(codes.Internal, "Internal server error")
	}

	user := data.User{
		Email:          r.Email,
		Username:       r.Username,
		FullName:       r.FullName,
		HashedPassword: r.Password,
		Role:           "user", // TODO(Jovan): For now
	}

	err = u.db.AddUser(&user)
	if err != nil {
		u.l.Errorf("failure adding user: %v\n", err)
		return &prusers.RegisterResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	resp, err := u.cc.CreateUserFolder(context.Background(), &prcontent.CreateUserFolderRequest{UserId: strconv.FormatUint(user.ID, 10)})

	if err != nil {
		u.l.Errorf("failed to create profile folders: %v", err)
		return &prusers.RegisterResponse{}, status.Error(codes.Internal, "Internal error")
	}

	profile := data.Profile{
		Username:        r.Username,
		UserID:          user.ID,
		Taggable:        false,
		Public:          false,
		Description:     r.Description,
		PhoneNumber:     r.PhoneNumber,
		Gender:          r.Gender,
		DateOfBirth:     time.Unix(r.DateOfBirth, 0),
		WebSite:         r.WebSite,
		PrivateProfile:  r.PrivateProfile,
		ProfileFolderId: resp.ProfileFolderId,
		PostsFolderId:   resp.PostsFolderId,
		StoriesFolderId: resp.StoryFolderId,
	}

	err = u.db.AddProfile(&profile)
	if err != nil {
		u.l.Printf("[ERROR] adding profile: %v\n", err)
	}

	// NOTE(Jovan): Saving refresh token
	_, err = u.ac.AddRefresh(context.Background(), &prauth.AddRefreshRequest{
		Username: r.Username,
		Token:    jws,
	})

	if err != nil {
		u.l.Errorf("failure adding refresh token: %v\n", err)
		return &prusers.RegisterResponse{}, status.Error(codes.Internal, "Internal server error")
	}

	go func() {
		_, err := u.ec.SendActivation(context.Background(), &premail.SendActivationRequest{Email: r.Email})
		if err != nil {
			u.l.Errorf("failure sending activation request: %v\n", err)
		}
	}()

	return &prusers.RegisterResponse{}, nil
}

func (u *Users) CheckPassword(ctx context.Context, r *prusers.CheckPasswordRequest) (*prusers.CheckPasswordResponse, error) {
	hashedPass, err := data.IsPasswordValid(u.db, r.Username, r.Password)
	if err != nil {
		u.l.Errorf("failure checking password: %v\n", err)
		return &prusers.CheckPasswordResponse{}, err
	}
	return &prusers.CheckPasswordResponse{HashedPass: hashedPass}, nil
}

func (u *Users) CheckEmail(ctx context.Context, r *prusers.CheckEmailRequest) (*prusers.CheckEmailResponse, error) {
	return &prusers.CheckEmailResponse{
		Verified: data.IsEmailVerified(u.db, r.Username),
	}, nil
}

func (u *Users) GetProfileByUsername(ctx context.Context, r *prusers.ProfileRequest) (*prusers.ProfileResponse, error) {
	//var response = &prusers.ProfileResponse{}

	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.ProfileResponse{}, err
	}

	user, err := u.db.GetUserByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting user: %v\n", err)
		return &prusers.ProfileResponse{}, err
	}

	isFollowing, err := data.CheckIfFollowing(u.db, r.User, profile.UserID)
	if err != nil {
		u.l.Printf("[ERROR] geting followers")
		return &prusers.ProfileResponse{}, err
	}

	following, err := data.GetFollowingCount(u.db, r.Username)
	if err != nil {
		return &prusers.ProfileResponse{}, err
	}

	followers, err := data.GetFollowerCount(u.db, r.Username)
	if err != nil {
		return &prusers.ProfileResponse{}, err
	}

	dateStr := strconv.FormatInt(profile.DateOfBirth.Unix(), 10)
	date, err := strconv.ParseInt(dateStr, 10, 64)

	if err != nil {
		u.l.Errorf("failed t oparse date string: %v", err)
		return &prusers.ProfileResponse{}, status.Error(codes.Internal, "Internal error")
	}

	return &prusers.ProfileResponse{
		Username:        profile.Username,
		Followers:       followers,
		Following:       following,
		FullName:        user.FullName,
		Description:     profile.Description,
		IsFollowing:     isFollowing,
		IsPublic:        profile.Public,
		PhoneNumber:     profile.PhoneNumber,
		Gender:          profile.Gender,
		DateOfBirth:     date,
		WebSite:         profile.WebSite,
		ProfileFolderId: profile.ProfileFolderId,
		PostsFolderId:   profile.PostsFolderId,
		StoriesFolderId: profile.StoriesFolderId,
		UserId:          profile.UserID,
		ProfilePictureURL: profile.ProfilePictureURL,
	}, nil
}

func (u *Users) Follow(ctx context.Context, r *prusers.FollowRequest) (*prusers.FollowRespose, error) {
	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.FollowRespose{}, err
	}
	profileToFollow, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile to follow: %v\n", err)
		return &prusers.FollowRespose{}, err
	}

	isFollowing, err := data.CheckIfFollowing(u.db, profile.Username, profileToFollow.UserID)
	if err != nil {
		u.l.Printf("[ERROR] geting followers")
		return &prusers.FollowRespose{}, err
	}

	if isFollowing {
		u.l.Printf("[WARNING] Already following")
		return &prusers.FollowRespose{}, nil
	}

	if !profileToFollow.Public {
		err = data.CreateFollowRequest(u.db, profileToFollow, profile)
		if err != nil {
			u.l.Printf("[ERROR] creating following request")
			return &prusers.FollowRespose{}, err
		}
		return &prusers.FollowRespose{}, nil

	}

	data.SetFollow(u.db, profile, profileToFollow)
	return &prusers.FollowRespose{}, nil
}

func (u *Users) GetFollowers(r *prusers.FollowerRequest, stream prusers.Users_GetFollowersServer) error {
	followers, err := data.GetFollowers(u.db, r.Username)
	if err != nil {
		return err
	}
	for _, profile := range followers {
		err = stream.Send(&prusers.ProfileFollower{
			Username: profile.Username,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Users) GetFollowing(r *prusers.FollowerRequest, stream prusers.Users_GetFollowersServer) error {
	followers, err := data.GetFollowing(u.db, r.Username)
	if err != nil {
		return err
	}
	for _, profile := range followers {
		err = stream.Send(&prusers.ProfileFollower{
			Username: profile.Username,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Users) UpdateProfile(ctx context.Context, r *prusers.UpdateRequest) (*prusers.UpdateResponse, error) {
	user, err := u.db.GetUserByUsername(r.OldUsername)
	if err != nil {
		u.l.Printf("[ERROR] geting user: %v\n", err)
		return &prusers.UpdateResponse{}, err
	}

	profile, err := u.db.GetProfileByUsername(r.OldUsername)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.UpdateResponse{}, err
	}

	user.Username = r.NewUsername
	user.Email = r.Email
	user.FullName = r.FullName
	profile.Username = r.NewUsername
	profile.Public = r.Public
	profile.Taggable = r.Taggable
	profile.PhoneNumber = r.PhoneNumber
	profile.Gender = r.Gender
	profile.DateOfBirth = time.Unix(r.DateOfBirth, 0)
	profile.WebSite = r.WebSite
	profile.PrivateProfile = r.PrivateProfile
	profile.Description = r.Description

	err = u.db.UpdateUser(user)
	if err != nil {
		u.l.Printf("[ERROR] failed to update user%v\n", err)
		return &prusers.UpdateResponse{}, err
	}
	err = u.db.UpdateProfile(profile)
	if err != nil {
		u.l.Printf("[ERROR] failed to update profile%v\n", err)
		return &prusers.UpdateResponse{}, err
	}

	return &prusers.UpdateResponse{}, nil

}


func (u *Users) GetSearchedUsers(ctx context.Context, r *prusers.SearchRequest) (*prusers.SearchResponse, error) {
	users, err := u.db.GetAllUsersByUsernameSubstring(r.Query)
	if err != nil {
		u.l.Printf("[ERROR] geting user: %v\n", err)
		return &prusers.SearchResponse{}, err
	}

	searchedUsers := []*prusers.SearchedUser{}

	for i := 0; i < len(users); i++ {
		su := users[i]
		searchedUsers = append(searchedUsers, &prusers.SearchedUser{
			Username: su.Username,
			ProfilePictureAddress: "PLEASE ADD PROFILE PICTURE ADDRESS HERE!"} )
	}

	return &prusers.SearchResponse{SearchedUser: searchedUsers}, nil

}
