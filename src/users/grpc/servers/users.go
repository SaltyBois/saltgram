package servers

import (
	"context"
	"fmt"
	"saltgram/protos/auth/prauth"
	"saltgram/protos/content/prcontent"
	"saltgram/protos/email/premail"
	"saltgram/protos/notifications/prnotifications"
	"saltgram/protos/users/prusers"
	"saltgram/users/data"
	"saltgram/users/saga"
	"strconv"
	"time"

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
	rc *saga.RedisClient
	nc prnotifications.NotificationsClient
}

func NewUsers(l *logrus.Logger, db *data.DBConn, ac prauth.AuthClient, ec premail.EmailClient, cc prcontent.ContentClient, nc prnotifications.NotificationsClient, rc *saga.RedisClient) *Users {
	return &Users{
		l:  l,
		db: db,
		ac: ac,
		ec: ec,
		cc: cc,
		rc: rc,
		nc: nc,
	}
}

func (u *Users) AcceptInfluencer(ctx context.Context, r *prusers.AcceptInfluencerRequest) (*prusers.AcceptInfluencerResponse, error) {
	err := u.db.RemoveInfluencerRequest(r.InfluencerId, r.CampaignId)
	if err != nil {
		u.l.Errorf("failed to remove influencer request: %v", err)
		return &prusers.AcceptInfluencerResponse{}, status.Error(codes.Internal, "Internal error")
	}
	return &prusers.AcceptInfluencerResponse{}, nil
}

func (u *Users) GetInfluencerRequests(ctx context.Context, r *prusers.GetInfluencerRequestsRequest) (*prusers.GetInfluencerRequestsResponse, error) {
	reqs, err := u.db.GetInfluencerRequests(r.InfluencerId)
	if err != nil {
		u.l.Errorf("failed to get influencer requests: %v", err)
		return &prusers.GetInfluencerRequestsResponse{}, status.Error(codes.Internal, "Internal error")
	}

	prreqs := []*prusers.Request{}
	for _, req := range *reqs {
		prreqs = append(prreqs, &prusers.Request{
			InfluencerId: req.InfluencerID,
			CampaignId:   req.CampaignID,
			Website:      req.Website,
		})
	}
	return &prusers.GetInfluencerRequestsResponse{Requests: prreqs}, nil
}

func (u *Users) InfluencerRequest(ctx context.Context, r *prusers.InfluencerRequestRequest) (*prusers.InfluencerRequestResponse, error) {
	err := u.db.AddInfluencerRequest(&data.InfluencerRequest{
		InfluencerID: r.InfluencerId,
		CampaignID:   r.CampaignId,
		Website:      r.Website,
	})
	if err != nil {
		u.l.Errorf("failed to add influencer request", err)
		return &prusers.InfluencerRequestResponse{}, status.Error(codes.Internal, "Internal error")
	}
	return &prusers.InfluencerRequestResponse{}, nil
}

func (u *Users) VerifyProfile(ctx context.Context, r *prusers.VerifyProfileRequest) (*prusers.VerifyProfileResponse, error) {
	err := u.db.VerifyProfile(r.UserId, r.AccountType)
	if err != nil {
		u.l.Errorf("failed to verify profile: %v", err)
		return &prusers.VerifyProfileResponse{}, status.Error(codes.Internal, "Internal error")
	}
	return &prusers.VerifyProfileResponse{}, nil
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

/*
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

	role := "user"
	if r.Agent {
		role = "agent"
	}

	user := data.User{
		Email:          r.Email,
		Username:       r.Username,
		FullName:       r.FullName,
		HashedPassword: r.Password,
		Role:           role,
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
		Taggable:        true,
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
		Messagable:      true,
		Verified:        false,
		AccountType:     "",
		Active:          true,
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

	if role != "agent" {
		go func() {
			_, err := u.ec.SendActivation(context.Background(), &premail.SendActivationRequest{Email: r.Email})
			if err != nil {
				u.l.Errorf("failure sending activation request: %v\n", err)
			}
		}()
	}

	return &prusers.RegisterResponse{}, nil
}
*/
func (u *Users) Register(ctx context.Context, r *prusers.RegisterRequest) (*prusers.RegisterResponse, error) {
	u.l.Infof("registering %v\n", r.Email)

	user := data.User{
		Email:          r.Email,
		Username:       r.Username,
		FullName:       r.FullName,
		HashedPassword: r.Password,
		Role:           "user", // TODO(Jovan): For now
	}

	err := u.db.AddUser(&user)
	if err != nil {
		u.l.Errorf("failure adding user: %v\n", err)
		return &prusers.RegisterResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	m := saga.Message{
		Service:        saga.AuthService,
		SenderService:  saga.UserService,
		Action:         saga.ActionStart,
		UserId:         user.ID,
		Username:       r.Username,
		Description:    r.Description,
		PhoneNumber:    r.PhoneNumber,
		Gender:         r.Gender,
		DateOfBirth:    r.DateOfBirth,
		WebSite:        r.WebSite,
		PrivateProfile: r.PrivateProfile,
		Email:          r.Email,
	}

	u.rc.Next(saga.AuthChannel, saga.AuthService, m)

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
	profile, err := u.db.GetProfileByUsername(r.User)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.ProfileResponse{}, err
	}

	user_profile, err := u.db.GetUserByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting user: %v\n", err)
		return &prusers.ProfileResponse{}, err
	}

	profile_show, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile_show: %v\n", err)
		return &prusers.ProfileResponse{}, err
	}

	isFollowing, err := data.CheckIfFollowing(u.db, profile, profile_show)
	if err != nil {
		u.l.Printf("[ERROR] geting followers")
		return &prusers.ProfileResponse{}, err
	}

	following, err := data.GetFollowingCount(u.db, profile_show)
	if err != nil {
		return &prusers.ProfileResponse{}, err
	}

	followers, err := data.GetFollowerCount(u.db, profile_show)
	if err != nil {
		return &prusers.ProfileResponse{}, err
	}

	dateStr := strconv.FormatInt(profile_show.DateOfBirth.Unix(), 10)
	date, err := strconv.ParseInt(dateStr, 10, 64)

	if err != nil {
		u.l.Errorf("failed t oparse date string: %v", err)
		return &prusers.ProfileResponse{}, status.Error(codes.Internal, "Internal error")
	}

	return &prusers.ProfileResponse{
		Username:          profile.Username,
		Followers:         followers,
		Following:         following,
		FullName:          user_profile.FullName,
		Description:       profile.Description,
		IsFollowing:       isFollowing,
		IsPublic:          !profile.PrivateProfile,
		PhoneNumber:       profile.PhoneNumber,
		Gender:            profile.Gender,
		DateOfBirth:       date,
		WebSite:           profile.WebSite,
		ProfileFolderId:   profile.ProfileFolderId,
		PostsFolderId:     profile.PostsFolderId,
		StoriesFolderId:   profile.StoriesFolderId,
		UserId:            profile.UserID,
		ProfilePictureURL: profile.ProfilePictureURL,
		Taggable:          profile.Taggable,
		Messageable:       profile.Messagable,
		Verified:          profile.Verified,
		AccountType:       profile.AccountType,
	}, nil
}

func (u *Users) Follow(ctx context.Context, r *prusers.FollowRequest) (*prusers.FollowRespose, error) {
	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.FollowRespose{}, err
	}
	profileToFollow, err := u.db.GetProfileByUsername(r.ToFollow)
	if err != nil {
		u.l.Printf("[ERROR] geting profile to follow: %v\n", err)
		return &prusers.FollowRespose{}, err
	}

	isFollowing, err := data.CheckIfFollowing(u.db, profile, profileToFollow)
	if err != nil {
		u.l.Printf("[ERROR] geting followers")
		return &prusers.FollowRespose{}, err
	}

	if isFollowing {
		u.l.Printf("[WARNING] Already following")
		return &prusers.FollowRespose{Message: "Already following"}, nil
	}

	blocked, _ := u.db.CheckIfBlocked(profile, profileToFollow)
	if blocked {
		u.db.UnblockProfile(profile, profileToFollow)
	}

	if profileToFollow.PrivateProfile {
		followRequest, _ := data.CheckForFollowingRequest(u.db, profileToFollow, profile)
		if followRequest {
			u.l.Printf("[WARNING] Follow request allready sent")
			return &prusers.FollowRespose{}, nil
		}
		err = data.CreateFollowRequest(u.db, profileToFollow, profile)
		if err != nil {
			u.l.Printf("[ERROR] creating following request")
			return &prusers.FollowRespose{}, err
		}
		_, err = u.nc.CreateFollowRequestNotification(context.Background(), &prnotifications.RequestUsername{UserId: profileToFollow.UserID, ReferredId: profile.UserID, ReferredUsername: profile.Username})
		if err != nil {
			u.l.Errorf("creating notification %v\n", err)
		}
		return &prusers.FollowRespose{Message: "PENDING"}, nil
	}

	_, err = u.nc.CreateFollowNotification(context.Background(), &prnotifications.RequestUsername{UserId: profileToFollow.UserID, ReferredId: profile.UserID, ReferredUsername: profile.Username})
	if err != nil {
		u.l.Errorf("creating notification %v\n", err)
	}

	data.SetFollow(u.db, profile, profileToFollow)
	return &prusers.FollowRespose{Message: "Following"}, nil
}

var ErrorNotFollowing = fmt.Errorf("not following")

func (u *Users) UnFollow(ctx context.Context, r *prusers.FollowRequest) (*prusers.FollowRespose, error) {
	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.FollowRespose{}, err
	}
	profileToUnfollow, err := u.db.GetProfileByUsername(r.ToFollow)
	if err != nil {
		u.l.Printf("[ERROR] geting profile to follow: %v\n", err)
		return &prusers.FollowRespose{}, err
	}

	isFollowing, err := data.CheckIfFollowing(u.db, profile, profileToUnfollow)
	if err != nil {
		u.l.Printf("[ERROR] Checking if is following")
		return &prusers.FollowRespose{}, err
	}

	if !isFollowing {
		u.l.Printf("[ERROR] Is not following")
		return &prusers.FollowRespose{Message: "Not following"}, ErrorNotFollowing
	}

	muted, _ := u.db.CheckIfMuted(profile, profileToUnfollow)
	if muted {
		u.db.UnmuteProfile(profile, profileToUnfollow)
	}

	closeFriend, _ := u.db.CheckIfCloseFriend(profile, profileToUnfollow)
	if closeFriend {
		u.db.RemoveCloseFriend(profile, profileToUnfollow)
	}

	data.Unfollow(u.db, profile, profileToUnfollow)
	return &prusers.FollowRespose{Message: "Unfollowed"}, nil
}

func (u *Users) GetFollowers(r *prusers.FollowerRequest, stream prusers.Users_GetFollowersServer) error {
	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return err
	}
	followers, err := data.GetFollowers(u.db, profile)
	if err != nil {
		u.l.Printf("[ERROR] query not working", err)
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

func (u *Users) GerFollowing(r *prusers.FollowerRequest, stream prusers.Users_GerFollowingServer) error {
	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return err
	}
	following, err := data.GetFollowing(u.db, profile)
	if err != nil {
		u.l.Printf("[ERROR] query not working", err)
		u.l.Printf("[ERROR] following", len(following))
		return err
	}
	for _, profile := range following {
		err = stream.Send(&prusers.ProfileFollower{
			Username:       profile.Username,
			Taggable:       profile.Taggable,
			ProfilePicture: profile.ProfilePictureURL,
			UserId:         strconv.FormatUint(profile.UserID, 10),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Users) GetFollowRequests(r *prusers.Profile, stream prusers.Users_GetFollowRequestsServer) error {
	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return err
	}
	profiles, err := data.GetFollowRequests(u.db, profile)
	if err != nil {
		return err
	}
	for _, profile_request := range profiles {
		err = stream.Send(&prusers.FollowingRequest{
			Username:       profile_request.Username,
			UserId:         profile_request.UserID,
			ProfilePicture: profile_request.ProfilePictureURL,
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
	profile.Messagable = r.Messageable

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

func (u *Users) GetByUserId(ctx context.Context, r *prusers.GetByIdRequest) (*prusers.GetByIdResponse, error) {
	user, err := u.db.GetUserById(r.Id)
	if err != nil {
		u.l.Printf("[ERROR] username: %v\n", r.Id)
		u.l.Errorf("failure getting user by id: %v\n", err)
		return &prusers.GetByIdResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prusers.GetByIdResponse{
		Id:             user.ID,
		Email:          user.Email,
		FullName:       user.FullName,
		Username:       user.Username,
		Role:           user.Role,
		HashedPassword: user.HashedPassword,
	}, nil

}

func (u *Users) GetSearchedUsers(ctx context.Context, r *prusers.SearchRequest) (*prusers.SearchResponse, error) {
	//users, err := u.db.GetAllUsersByUsernameSubstring(r.Query)
	profiles, err := u.db.GetAllUsersByUsernameSubstring(r.Query)
	if err != nil {
		u.l.Printf("[ERROR] geting user: %v\n", err)
		return &prusers.SearchResponse{}, err
	}

	searchedUsers := []*prusers.SearchedUser{}

	//for i := 0; i < len(users); i++ {
	for i := 0; i < len(profiles); i++ {
		//su := users[i]
		su := profiles[i]
		searchedUsers = append(searchedUsers, &prusers.SearchedUser{
			Username:              su.Username,
			ProfilePictureAddress: su.ProfilePictureURL})
	}

	return &prusers.SearchResponse{SearchedUser: searchedUsers}, nil

}

func (u *Users) SetFollowRequestRespond(ctx context.Context, r *prusers.FollowRequestRespond) (*prusers.FollowRequestSet, error) {
	UserProfile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.FollowRequestSet{}, err
	}

	profile_request, err := u.db.GetProfileByUsername(r.RequestUsername)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.FollowRequestSet{}, err
	}

	err = data.FollowRequestRespond(u.db, UserProfile, profile_request, r.Accepted)
	if err != nil {
		u.l.Printf("[ERROR] Setting respond")
		return &prusers.FollowRequestSet{}, err
	}

	u.l.Printf("[INFO] trying ot follow")
	u.l.Printf("[INFO] profile request: %v", profile_request)

	u.l.Printf("[INFO] profile: %v", UserProfile)
	if !r.Accepted {
		return &prusers.FollowRequestSet{}, nil
	}
	u.l.Printf("[INFO] trying to set follow for user %v to follow user: %v\n", profile_request, UserProfile)
	err = data.SetFollow(u.db, profile_request, UserProfile)
	if err != nil {
		u.l.Printf("[ERROR] following: %v\n", err)
		return &prusers.FollowRequestSet{}, err
	}
	_, err = u.nc.CreateFollowNotification(context.Background(), &prnotifications.RequestUsername{UserId: UserProfile.UserID, ReferredId: profile_request.UserID, ReferredUsername: profile_request.Username})
	if err != nil {
		u.l.Errorf("creating notification %v\n", err)
	}

	return &prusers.FollowRequestSet{}, nil
}

func (u *Users) GetFollowersDetailed(r *prusers.ProflieFollowRequest, stream prusers.Users_GetFollowersDetailedServer) error {
	userProfile, err := u.db.GetProfileByUsername(r.Logeduser)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return err
	}

	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return err
	}

	profiles, err := data.GetFollowers(u.db, profile)
	if err != nil {
		u.l.Printf("[ERROR]  fetching followers %v\n", err)
		return err
	}
	for _, profile := range profiles {
		following, err := data.CheckIfFollowing(u.db, userProfile, &profile)
		if err != nil {
			return err
		}
		if following {
			err = stream.Send(&prusers.ProfileFollowDetaild{
				Username:       profile.Username,
				Following:      following,
				Pending:        false,
				ProfliePicture: profile.ProfilePictureURL,
			})
			if err != nil {
				return err
			}
			continue
		}
		pending, err := data.CheckForFollowingRequest(u.db, &profile, userProfile)
		if err != nil {
			return err
		}
		err = stream.Send(&prusers.ProfileFollowDetaild{
			Username:       profile.Username,
			Following:      following,
			Pending:        pending,
			ProfliePicture: profile.ProfilePictureURL,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Users) GetFollowingDetailed(r *prusers.ProflieFollowRequest, stream prusers.Users_GetFollowingDetailedServer) error {
	userProfile, err := u.db.GetProfileByUsername(r.Logeduser)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return err
	}

	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return err
	}

	profiles, err := data.GetFollowing(u.db, profile)
	if err != nil {
		u.l.Printf("[ERROR]  fetching followers %v\n", err)
		return err
	}
	for _, profile := range profiles {
		following, err := data.CheckIfFollowing(u.db, userProfile, &profile)
		if err != nil {
			return err
		}
		if following {
			err = stream.Send(&prusers.ProfileFollowDetaild{
				Username:       profile.Username,
				Following:      following,
				Pending:        false,
				ProfliePicture: profile.ProfilePictureURL,
			})
			if err != nil {
				return err
			}
			continue
		}
		pending, err := data.CheckForFollowingRequest(u.db, &profile, userProfile)
		if err != nil {
			return err
		}
		err = stream.Send(&prusers.ProfileFollowDetaild{
			Username:       profile.Username,
			Following:      following,
			Pending:        pending,
			ProfliePicture: profile.ProfilePictureURL,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Users) CheckIfFollowing(ctx context.Context, r *prusers.ProflieFollowRequest) (*prusers.BoolResponse, error) {
	userProfile, err := u.db.GetProfileByUsername(r.Logeduser)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.BoolResponse{}, err
	}

	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.BoolResponse{}, err
	}

	following, err := data.CheckIfFollowing(u.db, userProfile, profile)
	if err != nil {
		return &prusers.BoolResponse{}, err
	}
	if following {
		return &prusers.BoolResponse{Response: true}, nil
	}
	return &prusers.BoolResponse{Response: false}, nil

}

func (u *Users) CheckForFollowingRequest(ctx context.Context, r *prusers.ProflieFollowRequest) (*prusers.BoolResponse, error) {
	userProfile, err := u.db.GetProfileByUsername(r.Logeduser)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.BoolResponse{}, err
	}

	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.BoolResponse{}, err
	}

	pending, err := data.CheckForFollowingRequest(u.db, profile, userProfile)
	if err != nil {
		return &prusers.BoolResponse{}, err
	}
	if pending {
		return &prusers.BoolResponse{Response: true}, nil
	}
	return &prusers.BoolResponse{Response: false}, nil

}

func (u *Users) MuteProfile(ctx context.Context, r *prusers.MuteProfileRequest) (*prusers.MuteProfileResponse, error) {
	userProfile, err := u.db.GetProfileByUsername(r.Logged)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.MuteProfileResponse{}, err
	}

	profile, err := u.db.GetProfileByUsername(r.Profile)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.MuteProfileResponse{}, err
	}

	err = u.db.MuteProfile(userProfile, profile)
	if err != nil {
		u.l.Printf("[ERROR] muting profile: %v\n", err)
		return &prusers.MuteProfileResponse{}, err
	}
	return &prusers.MuteProfileResponse{}, nil

}

func (u *Users) UnmuteProfile(ctx context.Context, r *prusers.UnmuteProfileRequest) (*prusers.UnmuteProfileResponse, error) {
	userProfile, err := u.db.GetProfileByUsername(r.Logged)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.UnmuteProfileResponse{}, err
	}

	profile, err := u.db.GetProfileByUsername(r.Profile)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.UnmuteProfileResponse{}, err
	}

	err = u.db.UnmuteProfile(userProfile, profile)
	if err != nil {
		u.l.Printf("[ERROR] unmuting profile: %v\n", err)
		return &prusers.UnmuteProfileResponse{}, err
	}
	return &prusers.UnmuteProfileResponse{}, nil
}

func (u *Users) GetMutedProfiles(r *prusers.Profile, stream prusers.Users_GetMutedProfilesServer) error {
	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return err
	}

	profiles, err := u.db.GetMutedProfiles(profile)
	if err != nil {
		u.l.Printf("[ERROR]  geting muted profiles %v\n", err)
		return err
	}
	for _, profile := range profiles {

		err = stream.Send(&prusers.ProfileMBCF{
			Username:          profile.Username,
			ProfilePictureURL: profile.ProfilePictureURL,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Users) CheckIfMuted(ctx context.Context, r *prusers.MuteProfileRequest) (*prusers.BoolResponse, error) {
	userProfile, err := u.db.GetProfileByUsername(r.Logged)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.BoolResponse{}, err
	}

	profile, err := u.db.GetProfileByUsername(r.Profile)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.BoolResponse{}, err
	}

	muted, err := u.db.CheckIfMuted(userProfile, profile)
	if err != nil {
		u.l.Printf("[ERROR] checking if muted: %v\n", err)
		return &prusers.BoolResponse{}, err
	}
	return &prusers.BoolResponse{Response: muted}, nil

}

func (u *Users) BlockProfile(ctx context.Context, r *prusers.BlockProfileRequest) (*prusers.BlockProfileResposne, error) {
	userProfile, err := u.db.GetProfileByUsername(r.Logged)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.BlockProfileResposne{}, err
	}

	profile, err := u.db.GetProfileByUsername(r.Profile)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.BlockProfileResposne{}, err
	}

	err = u.db.BlockProfile(userProfile, profile)
	if err != nil {
		u.l.Printf("[ERROR] blocking profile: %v\n", err)
		return &prusers.BlockProfileResposne{}, err
	}

	following, _ := data.CheckIfFollowing(u.db, profile, userProfile)
	if following {
		data.Unfollow(u.db, profile, userProfile)
	}
	follower, _ := data.CheckIfFollowing(u.db, userProfile, profile)
	if follower {
		data.Unfollow(u.db, userProfile, profile)
	}
	followRequest, _ := data.CheckForFollowingRequest(u.db, userProfile, profile)
	if followRequest {
		data.FollowRequestRespond(u.db, userProfile, profile, false)
	}
	muted, _ := u.db.CheckIfMuted(userProfile, profile)
	if muted {
		u.db.UnmuteProfile(userProfile, profile)
	}
	closeFriend, _ := u.db.CheckIfCloseFriend(userProfile, profile)
	if closeFriend {
		u.db.RemoveCloseFriend(userProfile, profile)
	}
	return &prusers.BlockProfileResposne{}, nil
}

func (u *Users) UnblockProfile(ctx context.Context, r *prusers.UnblockProfileRequest) (*prusers.UnblockProfileResposne, error) {
	userProfile, err := u.db.GetProfileByUsername(r.Logged)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.UnblockProfileResposne{}, err
	}

	profile, err := u.db.GetProfileByUsername(r.Profile)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.UnblockProfileResposne{}, err
	}

	err = u.db.UnblockProfile(userProfile, profile)
	if err != nil {
		u.l.Printf("[ERROR] unblocking profile: %v\n", err)
		return &prusers.UnblockProfileResposne{}, err
	}
	return &prusers.UnblockProfileResposne{}, nil
}

func (u *Users) GetBlockedProfiles(r *prusers.Profile, stream prusers.Users_GetBlockedProfilesServer) error {
	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return err
	}

	profiles, err := u.db.GetBlockedProfiles(profile)
	if err != nil {
		u.l.Printf("[ERROR]  geting blocked profiles %v\n", err)
		return err
	}
	for _, profile := range profiles {

		err = stream.Send(&prusers.ProfileMBCF{
			Username:          profile.Username,
			ProfilePictureURL: profile.ProfilePictureURL,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Users) CheckIfBlocked(ctx context.Context, r *prusers.BlockProfileRequest) (*prusers.BoolResponse, error) {
	userProfile, err := u.db.GetProfileByUsername(r.Logged)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.BoolResponse{}, err
	}

	profile, err := u.db.GetProfileByUsername(r.Profile)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.BoolResponse{}, err
	}

	blocked, err := u.db.CheckIfBlocked(userProfile, profile)
	if err != nil {
		u.l.Printf("[ERROR] checking if blocked: %v\n", err)
		return &prusers.BoolResponse{}, err
	}
	return &prusers.BoolResponse{Response: blocked}, nil

}

func (u *Users) AddCloseFriend(ctx context.Context, r *prusers.CloseFriendRequest) (*prusers.CloseFriendResposne, error) {
	userProfile, err := u.db.GetProfileByUsername(r.Logged)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.CloseFriendResposne{}, err
	}

	profile, err := u.db.GetProfileByUsername(r.Profile)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.CloseFriendResposne{}, err
	}
	err = u.db.AddCloseFriend(userProfile, profile)
	if err != nil {
		u.l.Printf("[ERROR] adding close friend: %v\n", err)
		return &prusers.CloseFriendResposne{}, err
	}
	return &prusers.CloseFriendResposne{}, nil
}

func (u *Users) RemoveCloseFriend(ctx context.Context, r *prusers.CloseFriendRequest) (*prusers.CloseFriendResposne, error) {
	userProfile, err := u.db.GetProfileByUsername(r.Logged)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.CloseFriendResposne{}, err
	}

	profile, err := u.db.GetProfileByUsername(r.Profile)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.CloseFriendResposne{}, err
	}
	err = u.db.RemoveCloseFriend(userProfile, profile)
	if err != nil {
		u.l.Printf("[ERROR] removing close friend %v\n", err)
		return &prusers.CloseFriendResposne{}, err
	}
	return &prusers.CloseFriendResposne{}, nil
}

func (u *Users) GetCloseFriends(r *prusers.Profile, stream prusers.Users_GetCloseFriendsServer) error {
	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return err
	}

	profiles, err := u.db.GetCloseFriends(profile)
	if err != nil {
		u.l.Printf("[ERROR]  geting close friends %v\n", err)
		return err
	}
	for _, profile := range profiles {

		err = stream.Send(&prusers.ProfileMBCF{
			Username:          profile.Username,
			ProfilePictureURL: profile.ProfilePictureURL,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Users) GetProfilesForCloseFriends(r *prusers.Profile, stream prusers.Users_GetProfilesForCloseFriendsServer) error {
	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return err
	}

	profiles, err := u.db.GetProfilesForCloseFriends(profile)
	if err != nil {
		u.l.Printf("[ERROR]  geting profiles for close friends %v\n", err)
		return err
	}
	for _, profile := range profiles {

		err = stream.Send(&prusers.ProfileMBCF{
			Username:          profile.Username,
			ProfilePictureURL: profile.ProfilePictureURL,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Users) DeleteProfile(ctx context.Context, r *prusers.Profile) (*prusers.DeleteProfileResponse, error) {
	profile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.DeleteProfileResponse{}, err
	}
	profile.Active = false
	err = u.db.UpdateProfile(profile)
	if err != nil {
		u.l.Printf("[ERROR] deleting profile: %v\n", err)
		return &prusers.DeleteProfileResponse{}, err
	}

	err = u.db.DeleteFromFollwing(profile)
	if err != nil {
		u.l.Printf("[ERROR] deleting form following %v", err)
	}

	err = u.db.DeleteFromMuted(profile)
	if err != nil {
		u.l.Printf("[ERROR] deleting form muted %v", err)
	}

	err = u.db.DeleteFromBlocked(profile)
	if err != nil {
		u.l.Printf("[ERROR] deleting form blocked %v", err)
	}

	err = u.db.DeleteFromCloseFriends(profile)
	if err != nil {
		u.l.Printf("[ERROR] deleting form close friends %v", err)
	}

	return &prusers.DeleteProfileResponse{}, nil

}

func (u *Users) CheckActive(ctx context.Context, r *prusers.Profile) (*prusers.BoolResponse, error) {
	active, err := u.db.CheckActive(r.Username)
	if err != nil {
		return &prusers.BoolResponse{}, err
	}
	return &prusers.BoolResponse{Response: active}, nil
}

func (u *Users) GetFollowingMain(r *prusers.Profile, stream prusers.Users_GetFollowingMainServer) error {
	userProfile, err := u.db.GetProfileByUsername(r.Username)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return err
	}

	profiles, err := data.GetFollowing(u.db, userProfile)
	if err != nil {
		u.l.Printf("[ERROR]  fetching followers %v\n", err)
		return err
	}
	for _, profile := range profiles {
		muted, err := u.db.CheckIfMuted(userProfile, &profile)
		if err != nil {
			return err
		}
		if muted {
			continue
		}
		err = stream.Send(&prusers.ProfileMBCF{
			Username:          profile.Username,
			ProfilePictureURL: profile.ProfilePictureURL,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *Users) GetProfileByUserId(ctx context.Context, r *prusers.GetByIdRequest) (*prusers.ProfileMBCF, error) {
	profile, err := u.db.GetProfileByUserId(r.Id)
	if err != nil {
		u.l.Printf("[ERROR] geting profile: %v\n", err)
		return &prusers.ProfileMBCF{}, err
	}

	return &prusers.ProfileMBCF{Username: profile.Username, ProfilePictureURL: profile.ProfilePictureURL}, nil
}
