package servers

import (
	"bytes"
	"context"
	"io"
	"saltgram/content/data"
	"saltgram/content/gdrive"
	"saltgram/protos/content/prcontent"
	"strconv"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Content struct {
	prcontent.UnimplementedContentServer
	l  *logrus.Logger
	db *data.DBConn
	g  *gdrive.GDrive
}

func NewContent(l *logrus.Logger, db *data.DBConn, g *gdrive.GDrive) *Content {
	return &Content{
		l:  l,
		db: db,
		g:  g,
	}
}

func (c *Content) GetPostPreviewURL(ctx context.Context, r *prcontent.GetPostPreviewURLRequest) (*prcontent.GetPostPreviewURLResponse, error) {
	post, err := c.db.GetPost(r.PostId)
	if err != nil {
		c.l.Errorf("failed to get shared media: %v", err)
		return &prcontent.GetPostPreviewURLResponse{}, status.Error(codes.Internal, "Internal error")
	}

	return &prcontent.GetPostPreviewURLResponse{Url: post.SharedMedia.Media[0].URL}, nil
}

func (c *Content) AddVerificationImage(stream prcontent.Content_AddVerificationImageServer) error {
	r, err := stream.Recv()
	if err != nil {
		c.l.Errorf("failed to receive verification meta: %v", err)
		return err
	}
	v := data.Verification{
		UserID: r.GetInfo().UserId,
	}
	imageData := bytes.Buffer{}
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			c.l.Info("verification image received")
			break
		}
		if err != nil {
			c.l.Errorf("error while receiving image data: %v", err)
			return status.Error(codes.InvalidArgument, "Bad request")
		}
		// TODO(Jovan): Prevent large uploads?
		_, err = imageData.Write(chunk.GetImage())
		if err != nil {
			c.l.Errorf("error appending image chunk data: %v", err)
			return status.Error(codes.Internal, "Internal server error")
		}
	}
	url, err := c.g.UploadVerificationImage(strconv.FormatUint(v.UserID, 10), r.GetInfo().Filename, &imageData)
	if err != nil {
		c.l.Errorf("failed to upload verification image: %v", err)
		return err
	}
	v.URL = url
	err = c.db.AddVerification(&v)
	if err != nil {
		c.l.Errorf("failed to save verification: %v", err)
		return err
	}

	err = stream.SendAndClose(&prcontent.AddVerificationImageResponse{Url: url})
	return err
}

func (c *Content) GetHighlights(ctx context.Context, r *prcontent.GetHighlightsRequest) (*prcontent.GetHighlightsResponse, error) {
	highlights, err := c.db.GetHighlights(r.UserId)
	if err != nil {
		c.l.Errorf("failed to get highlights: %v", err)
		return &prcontent.GetHighlightsResponse{}, status.Error(codes.Internal, "Internal error")
	}

	highlightsPR := []*prcontent.Highlight{}

	for _, h := range highlights {
		highlightsPR = append(highlightsPR, data.DataToPRHighlight(h))
	}

	return &prcontent.GetHighlightsResponse{
		Highlights: highlightsPR,
	}, nil
}

func (c *Content) GetStoriesIndividual(ctx context.Context, r *prcontent.GetStoriesIndividualRequest) (*prcontent.GetStoriesIndividualResponse, error) {
	stories, err := c.db.GetStoryByUser(r.UserId)
	if err != nil {
		c.l.Errorf("failed to get users stories: %v", err)
		return &prcontent.GetStoriesIndividualResponse{}, status.Error(codes.Internal, "Internal error")
	}

	s := []*prcontent.Story{}

	for _, st := range stories {
		s = append(s, data.DataToPRStory(st))
	}

	return &prcontent.GetStoriesIndividualResponse{
		Stories: s,
	}, nil
}

func (c *Content) AddHighlight(ctx context.Context, r *prcontent.AddHighlightRequest) (*prcontent.AddHighlightResponse, error) {
	storyIds := r.Stories

	stories, err := c.db.GetMediaByIds(storyIds...)
	if err != nil {
		c.l.Errorf("failed to get stories for highlight: %v", err)
		return &prcontent.AddHighlightResponse{}, status.Error(codes.Internal, "Internal error")
	}

	newHighlight := data.Highlight{
		Name:    r.Name,
		UserID:  r.UserId,
		Stories: stories,
	}
	err = c.db.CreateHighlight(&newHighlight)
	if err != nil {
		c.l.Errorf("failed to create highlight: %v", err)
		return &prcontent.AddHighlightResponse{}, status.Error(codes.Internal, "Internal error")
	}

	return &prcontent.AddHighlightResponse{}, nil
}

func (c *Content) CreatePost(ctx context.Context, r *prcontent.CreatePostRequest) (*prcontent.CreatePostResponse, error) {
	ageGroup := data.EAgeGroup_PRE20
	switch r.AgeGroup {
	case "Pre 20s":
		ageGroup = data.EAgeGroup_PRE20
	case "20s":
		ageGroup = data.EAgeGroup_20
	case "30s":
		ageGroup = data.EAgeGroup_30
	default:
		ageGroup = data.EAgeGroup_PRE20
	}
	sm := data.SharedMedia{
		IsCampaign:       r.Campaign,
		CampaignAgeGroup: data.EAgeGroup(ageGroup),
		CampaignOneTime:  r.CampaignOneTime,
		CampaignWebsite:  r.CampaignWebsite,
		CampaignStart:    r.CampaignStart,
		CampaignEnd:      r.CampaignEnd,
	}
	post := &data.Post{
		UserID:      r.UserId,
		SharedMedia: sm,
	}
	err := c.db.AddPost(post)
	if err != nil {
		c.l.Errorf("failed to create post: %v", err)
		return &prcontent.CreatePostResponse{}, status.Error(codes.Internal, "Internal error")
	}

	return &prcontent.CreatePostResponse{PostId: post.ID}, nil
}

func (c *Content) CreateStory(ctx context.Context, r *prcontent.CreateStoryRequest) (*prcontent.CreateStoryResponse, error) {
	ageGroup := data.EAgeGroup_PRE20
	switch r.AgeGroup {
	case "Pre 20s":
		ageGroup = data.EAgeGroup_PRE20
	case "20s":
		ageGroup = data.EAgeGroup_20
	case "30s":
		ageGroup = data.EAgeGroup_30
	default:
		ageGroup = data.EAgeGroup_PRE20
	}
	sm := data.SharedMedia{
		IsCampaign:       r.Campaign,
		CampaignAgeGroup: data.EAgeGroup(ageGroup),
		CampaignOneTime:  r.CampaignOneTime,
		CampaignWebsite:  r.CampaignWebsite,
		CampaignStart:    r.CampaignStart,
		CampaignEnd:      r.CampaignEnd,
	}
	story := &data.Story{
		UserID:      r.UserId,
		SharedMedia: sm,
	}
	err := c.db.AddStory(story)
	if err != nil {
		c.l.Errorf("failed to create shared media: %v", err)
		return &prcontent.CreateStoryResponse{}, status.Error(codes.Internal, "Internal error")
	}
	return &prcontent.CreateStoryResponse{StoryId: story.ID}, nil
}

func (c *Content) CreateUserFolder(ctx context.Context, r *prcontent.CreateUserFolderRequest) (*prcontent.CreateUserFolderResponse, error) {
	profile, posts, stories, err := c.g.CreateUserFolder(r.UserId)
	if err != nil {
		c.l.Errorf("failed to create user folders: %v", err)
		return &prcontent.CreateUserFolderResponse{}, status.Error(codes.Internal, "Internal error")
	}
	return &prcontent.CreateUserFolderResponse{
		ProfileFolderId: profile,
		PostsFolderId:   posts,
		StoryFolderId:   stories,
	}, nil
}

func (c *Content) GetSharedMedia(r *prcontent.SharedMediaRequest, stream prcontent.Content_GetSharedMediaServer) error {
	sharedMedias, err := c.db.GetSharedMediaByUser(r.UserId)
	if err != nil {
		c.l.Errorf("failure getting user shared media: %v\n", err)
		return err
	}
	for _, sm := range *sharedMedias {
		media := []*prcontent.Media{}
		for _, m := range sm.Media {
			media = append(media, &prcontent.Media{
				Filename:    m.Filename,
				Description: m.Description,
				AddedOn:     m.AddedOn,
				Location: &prcontent.Location{
					Country: m.Location.Country,
					State:   m.Location.State,
					ZipCode: m.Location.ZipCode,
					Street:  m.Location.Street,
				},
			})
		}
		err = stream.Send(&prcontent.SharedMediaResponse{
			Media: media,
		})
		if err != nil {
			c.l.Errorf("failure sending shared media response: %v\n", err)
			return err
		}
	}
	return nil
}

func (c *Content) GetProfilePicture(ctx context.Context, r *prcontent.GetProfilePictureRequest) (*prcontent.GetProfilePictureResponse, error) {
	profilePicture, err := c.db.GetProfilePictureByUser(r.UserId)
	if err != nil {
		return &prcontent.GetProfilePictureResponse{Url: ""}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prcontent.GetProfilePictureResponse{Url: profilePicture.URL}, nil
}

/*func (c *Content) GetPostsByUser(ctx context.Context, r *prcontent.GetPostsRequest) (*prcontent.GetPostsResponse, error) {
	posts, err := c.db.GetPostByUser(r.UserId)
	if err != nil {
		c.l.Errorf("failure getting user posts: %v\n", err)
		return &prcontent.GetPostsResponse{}, err
	}
	response := []*prcontent.Post{}

	for _, p := range *posts {
		media := []*prcontent.Media{}
		for _, m := range p.SharedMedia.Media {
			media = append(media, &prcontent.Media{
				Filename:    m.Filename,
				Description: m.Description,
				AddedOn:     m.AddedOn,
				Location: &prcontent.Location{
					Country: m.Location.Country,
					State:   m.Location.State,
					ZipCode: m.Location.ZipCode,
					Street:  m.Location.Street,
				},
				Url: m.URL,
			})
		}
		response = append(response, &prcontent.Post{
			Id:     p.ID,
			UserId: p.UserID,
			SharedMedia: &prcontent.SharedMedia{
				Media: media,
			},
		})

	}
	return &prcontent.GetPostsResponse{Post: response}, nil
}*/

func (c *Content) GetPostsByUser(r *prcontent.GetPostsRequest, stream prcontent.Content_GetPostsByUserServer) error {
	posts, err := c.db.GetPostByUser(r.UserId)
	if err != nil {
		c.l.Errorf("failure getting user posts: %v\n", err)
		return err
	}
	for _, p := range *posts {
		media := []*prcontent.Media{}
		for _, m := range p.SharedMedia.Media {
			tags := []*prcontent.Tag{}
			for _, t := range m.Tags {
				tags = append(tags, &prcontent.Tag{Id: t.ID, Value: t.Value})
			}
			userTags := []*prcontent.UserTag{}
			for _, t := range m.TaggedUsers {
				userTags = append(userTags, &prcontent.UserTag{Id: t.UserID})
			}
			media = append(media, &prcontent.Media{
				Filename:    m.Filename,
				Description: m.Description,
				AddedOn:     m.AddedOn,
				Tags:        tags,
				Location: &prcontent.Location{
					Country: m.Location.Country,
					State:   m.Location.State,
					City:    m.Location.City,
					ZipCode: m.Location.ZipCode,
					Street:  m.Location.Street,
					Name:    m.Location.Name,
				},
				Url:      m.URL,
				MimeType: prcontent.EMimeType(m.MimeType),
				UserTags: userTags,
			})
		}
		post := &prcontent.Post{
			Id:     strconv.FormatUint(p.ID, 10),
			UserId: strconv.FormatUint(p.UserID, 10),
			SharedMedia: &prcontent.SharedMedia{
				Media: media,
			},
			IsCampaign: p.SharedMedia.IsCampaign,
			CampaignWebsite: p.SharedMedia.CampaignWebsite,
		}
		err = stream.Send(&prcontent.GetPostsResponse{
			Post: post,
		})
		if err != nil {
			c.l.Errorf("failed sending post response: %v\n", err)
			return err
		}
	}
	return nil
}

/*func (c *Content) GetPostsByUserReaction(r *prcontent.GetPostsRequest, stream prcontent.Content_GetPostsByUserReactionServer) error {
	posts, err := c.db.GetPostsByReaction(r.UserId)
	if err != nil {
		c.l.Errorf("failure getting user posts: %v\n", err)
		return err
	}
	for _, p := range *posts {
		media := []*prcontent.Media{}
		for _, m := range p.SharedMedia.Media {
			media = append(media, &prcontent.Media{
				Filename:    m.Filename,
				Description: m.Description,
				AddedOn:     m.AddedOn,
				Location: &prcontent.Location{
					Country: m.Location.Country,
					State:   m.Location.State,
					ZipCode: m.Location.ZipCode,
					Street:  m.Location.Street,
				},
				Url:      m.URL,
				MimeType: prcontent.EMimeType(m.MimeType),
			})
		}
		post := &prcontent.Post{
			Id:     strconv.FormatUint(p.ID, 10),
			UserId: strconv.FormatUint(p.UserID, 10),
			SharedMedia: &prcontent.SharedMedia{
				Media: media,
			},
		}
		err = stream.Send(&prcontent.GetPostsResponse{
			Post: post,
		})
		if err != nil {
			c.l.Errorf("failed sending post response: %v\n", err)
			return err
		}
	}
	return nil
}*/

func (c *Content) GetPostsByUserReaction(ctx context.Context, r *prcontent.GetPostsByUserReactionRequest) (*prcontent.GetPostsByUserReactionResponse, error) {
	posts, err := c.db.GetPostsByReaction(r.Id)
	if err != nil {
		c.l.Errorf("failure getting posts: %v\n", err)
		return &prcontent.GetPostsByUserReactionResponse{}, err
	}

	retVal := []*prcontent.Post{}

	for _, p := range *posts {
		media := []*prcontent.Media{}
		for _, m := range p.SharedMedia.Media {
			tags := []*prcontent.Tag{}
			for _, t := range m.Tags {
				tags = append(tags, &prcontent.Tag{Id: t.ID, Value: t.Value})
			}
			userTags := []*prcontent.UserTag{}
			for _, t := range m.TaggedUsers {
				userTags = append(userTags, &prcontent.UserTag{Id: t.UserID})
			}
			media = append(media, &prcontent.Media{
				Filename:    m.Filename,
				Description: m.Description,
				AddedOn:     m.AddedOn,
				Tags:        tags,
				Location: &prcontent.Location{
					Country: m.Location.Country,
					State:   m.Location.State,
					City:    m.Location.City,
					ZipCode: m.Location.ZipCode,
					Street:  m.Location.Street,
					Name:    m.Location.Name,
				},
				Url:      m.URL,
				MimeType: prcontent.EMimeType(m.MimeType),
				UserTags: userTags,
			})
		}
		post := &prcontent.Post{
			Id:     strconv.FormatUint(p.ID, 10),
			UserId: strconv.FormatUint(p.UserID, 10),
			SharedMedia: &prcontent.SharedMedia{
				Media: media,
			},
		}

		retVal = append(retVal, post)

		if err != nil {
			c.l.Errorf("failed sending post response: %v\n", err)
			return &prcontent.GetPostsByUserReactionResponse{}, err
		}
	}
	return &prcontent.GetPostsByUserReactionResponse{Post: retVal}, nil
}

func (c *Content) AddProfilePicture(stream prcontent.Content_AddProfilePictureServer) error {
	r, err := stream.Recv()
	if err != nil {
		c.l.Errorf("failed to recieve profile metadata: %v", err)
	}

	profilePicture := data.ProfilePicture{
		UserID: r.GetInfo().UserId,
	}
	imageData := bytes.Buffer{}
	c.l.Info("receiving profile image for user: %v...", profilePicture.UserID)
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			c.l.Info("profile image received")
			break
		}
		if err != nil {
			c.l.Errorf("error while receiving image data: %v", err)
			return status.Error(codes.InvalidArgument, "Bad request")
		}
		// TODO(Jovan): Prevent large uploads?
		_, err = imageData.Write(chunk.GetImage())
		if err != nil {
			c.l.Errorf("error appending image chunk data: %v", err)
			return status.Error(codes.Internal, "Internal server error")
		}
	}

	url, err := c.g.UploadProfilePicture(r.GetInfo().ProfileFolderId, &imageData)
	if err != nil {
		c.l.Errorf("failed to upload profile picture: %v", err)
		return status.Error(codes.Internal, "Internal server error")
	}

	profilePicture.URL = url

	err = c.db.AddProfilePicture(&profilePicture)
	if err != nil {
		c.l.Errorf("failed adding profile picture: %v", err)
		return status.Error(codes.InvalidArgument, "Bad request")
	}
	err = stream.SendAndClose(&prcontent.AddProfilePictureResponse{
		Url: url,
	})
	if err != nil {
		c.l.Errorf("failed to send and close: %v", err)
	}
	return err
}

// What?
// func (c *Content) AddSharedMedia(ctx context.Context, r *prcontent.AddSharedMediaRequest) (*prcontent.AddSharedMediaResponse, error) {

// 	media := []*data.Media{}
// 	for _, m := range r.Media {
// 		media = append(media, &data.Media{
// 			Filename:    m.Filename,
// 			Tags:        []data.Tag{}, // TODO
// 			Description: m.Description,
// 			Location: data.Location{
// 				Country: m.Location.Country,
// 				State:   m.Location.State,
// 				ZipCode: m.Location.ZipCode,
// 				Street:  m.Location.Street,
// 			},
// 			AddedOn: m.AddedOn,
// 		})
// 	}
// 	sharedMedia := data.SharedMedia{
// 		Media: media,
// 	}

// 	err := c.db.AddSharedMedia(&sharedMedia)
// 	if err != nil {
// 		c.l.Errorf("failure adding user: %v\n", err)
// 		return &prcontent.AddSharedMediaResponse{}, status.Error(codes.InvalidArgument, "Bad request")
// 	}

// 	return &prcontent.AddSharedMediaResponse{}, nil
// }

func (c *Content) AddStory(stream prcontent.Content_AddStoryServer) error {
	r, err := stream.Recv()
	if err != nil {
		c.l.Errorf("failed to receive story metadata: %v", err)
		return status.Error(codes.InvalidArgument, "Invalid argument error")
	}

	tags := []data.Tag{}
	for _, tag := range r.GetInfo().Media.Tags {

		t := data.Tag{
			Value: tag.Value,
		}
		ta, err := c.db.GetIfExists(t.Value)
		if err != nil {
			if err == data.ErrTagNotExists {
				ta, err = c.db.AddTag(&t)
				if err != nil {
					c.l.Errorf("failed to add tag: %v", err)
					return status.Error(codes.Internal, "Internal error")
				}
			} else {
				c.l.Errorf("failed to get tag: %v", err)
				return status.Error(codes.Internal, "Internal error")
			}
		}

		tags = append(tags, *ta)
	}
	/////
	for _, t := range tags {
		c.l.Info("Getting tag id: ", t.ID)
	}
	/////

	userTags := []data.UserTag{}
	for _, userTag := range r.GetInfo().Media.UserTags {

		t := data.UserTag{
			UserID: userTag.Id,
		}
		ta, err := c.db.GetUserTagById(t.UserID)
		if err != nil {
			if err == data.ErrUserTagNotFound {
				ta, err = c.db.AddUserTag(&t)
				if err != nil {
					c.l.Errorf("failed to add user tag: %v", err)
					return status.Error(codes.Internal, "Internal error")
				}
			} else {
				c.l.Errorf("failed to get user tag: %v", err)
				return status.Error(codes.Internal, "Internal error")
			}
		}

		userTags = append(userTags, *ta)
	}
	/////
	for _, t := range userTags {
		c.l.Info("Getting tag id: ", t.UserID)
	}
	/////

	location := r.GetInfo().Media.Location

	media := &data.Media{
		Filename:    r.GetInfo().Media.Filename,
		Tags:        tags,
		Description: r.GetInfo().Media.Description,
		Location: data.Location{
			Country: location.Country,
			State:   location.State,
			ZipCode: location.ZipCode,
			City:    location.City,
			Street:  location.Street,
			Name:    location.Name,
		},
		AddedOn:     r.GetInfo().Media.AddedOn,
		TaggedUsers: userTags,
	}

	imageData := bytes.Buffer{}
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			c.l.Info("profile image received")
			break
		}
		if err != nil {
			c.l.Errorf("error while receiving image data: %v", err)
			return status.Error(codes.InvalidArgument, "Bad request")
		}
		// TODO(Jovan): Prevent large uploads?
		_, err = imageData.Write(chunk.GetImage())
		if err != nil {
			c.l.Errorf("error appending image chunk data: %v", err)
			return status.Error(codes.Internal, "Internal server error")
		}
	}
	url, mimeType, err := c.g.UploadStory(r.GetInfo().StoriesFolderId, media.Filename, &imageData)
	if err != nil {
		c.l.Errorf("failed to upload story: %v", err)
		return status.Error(codes.Internal, "Internal error")
	}
	media.URL = url
	media.MimeType = mimeType
	err = c.db.AddMediaToStory(r.GetInfo().StoryId, media)
	if err != nil {
		c.l.Errorf("faield to add media to story: %v", err)
		return status.Error(codes.Internal, "Internal error")
	}
	err = stream.SendAndClose(&prcontent.AddStoryResponse{})
	if err != nil {
		c.l.Errorf("failed to send and close: %v", err)
		return status.Error(codes.Internal, "Internal error")
	}
	return nil
}

func (c *Content) AddPost(stream prcontent.Content_AddPostServer) error {
	r, err := stream.Recv()
	if err != nil {
		c.l.Errorf("failed to recieve profile metadata: %v", err)
		return status.Error(codes.InvalidArgument, "Invalid argument error")
	}

	tags := []data.Tag{}
	for _, tag := range r.GetInfo().Media.Tags {

		t := data.Tag{
			Value: tag.Value,
		}
		ta, err := c.db.GetIfExists(t.Value)
		if err != nil {
			if err == data.ErrTagNotExists {
				ta, err = c.db.AddTag(&t)
				if err != nil {
					c.l.Errorf("failed to add tag: %v", err)
					return status.Error(codes.Internal, "Internal error")
				}
			} else {
				c.l.Errorf("failed to get tag: %v", err)
				return status.Error(codes.Internal, "Internal error")
			}
		}

		tags = append(tags, *ta)
	}
	/////
	for _, t := range tags {
		c.l.Info("Getting tag id: ", t.ID)
	}
	/////

	userTags := []data.UserTag{}
	for _, userTag := range r.GetInfo().Media.UserTags {

		t := data.UserTag{
			UserID: userTag.Id,
		}
		ta, err := c.db.GetUserTagById(t.UserID)
		if err != nil {
			if err == data.ErrUserTagNotFound {
				ta, err = c.db.AddUserTag(&t)
				if err != nil {
					c.l.Errorf("failed to add user tag: %v", err)
					return status.Error(codes.Internal, "Internal error")
				}
			} else {
				c.l.Errorf("failed to get user tag: %v", err)
				return status.Error(codes.Internal, "Internal error")
			}
		}

		userTags = append(userTags, *ta)
	}
	/////
	for _, t := range userTags {
		c.l.Info("Getting tag id: ", t.UserID)
	}
	/////

	location := r.GetInfo().Media.Location

	media := &data.Media{
		SharedMediaID: r.GetInfo().Media.SharedMediaId,
		Filename:      r.GetInfo().Media.Filename,
		Tags:          tags,
		Description:   r.GetInfo().Media.Description,
		Location: data.Location{
			Country: location.Country,
			State:   location.State,
			ZipCode: location.ZipCode,
			City:    location.City,
			Street:  location.Street,
			Name:    location.Name,
		},
		AddedOn:     r.GetInfo().Media.AddedOn,
		TaggedUsers: userTags,
	}

	/////
	for _, t := range media.Tags {
		c.l.Info("Getting tag from media      id: ", t.ID)
	}
	/////

	imageData := bytes.Buffer{}
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			c.l.Info("profile image received")
			break
		}
		if err != nil {
			c.l.Errorf("error while receiving image data: %v", err)
			return status.Error(codes.InvalidArgument, "Bad request")
		}
		// TODO(Jovan): Prevent large uploads?
		_, err = imageData.Write(chunk.GetImage())
		if err != nil {
			c.l.Errorf("error appending image chunk data: %v", err)
			return status.Error(codes.Internal, "Internal server error")
		}
	}

	url, mimeType, err := c.g.UploadPost(r.GetInfo().PostsFolderId, media.Filename, &imageData)
	if err != nil {
		c.l.Errorf("failure to upload post: %v", err)
		return status.Error(codes.Internal, "Internal error")
	}

	media.URL = url
	media.MimeType = mimeType

	err = c.db.AddMediaToPost(r.GetInfo().PostId, media)
	if err != nil {
		c.l.Errorf("failed to add media to shared media: %v", err)
		return status.Error(codes.Internal, "Internal error")
	}

	err = stream.SendAndClose(&prcontent.AddPostResponse{})
	if err != nil {
		c.l.Errorf("failure to send and close: %v", err)
		return status.Error(codes.Internal, "Internal error")
	}
	return nil
}

func (c *Content) AddComment(ctx context.Context, r *prcontent.AddCommentRequest) (*prcontent.AddCommentResponse, error) {

	comment := data.Comment{
		Content: r.Content,
		UserID:  r.UserId,
		PostID:  r.PostId,
	}

	err := c.db.AddComment(&comment)
	if err != nil {
		c.l.Errorf("Failed adding comment: %v\n", err)
		return &prcontent.AddCommentResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prcontent.AddCommentResponse{}, nil
}

func (c *Content) AddReaction(ctx context.Context, r *prcontent.AddReactionRequest) (*prcontent.AddReactionResponse, error) {

	reaction := data.Reaction{
		ReactionType: r.ReactionType,
		UserID:       r.UserId,
		PostID:       r.PostId,
	}

	err := c.db.AddReaction(&reaction)
	if err != nil {
		c.l.Errorf("Failed adding reaction: %v\n", err)
		return &prcontent.AddReactionResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prcontent.AddReactionResponse{}, nil
}

func (c *Content) GetReactions(r *prcontent.GetReactionsRequest, stream prcontent.Content_GetReactionsServer) error {
	reactions, err := c.db.GetReactionByPostId(r.PostId)
	if err != nil {
		c.l.Errorf("failure getting reactions: %v\n", err)
		return err
	}
	for _, r := range *reactions {
		err = stream.Send(&prcontent.GetReactionsResponse{
			Id:           strconv.FormatUint(r.ID, 10),
			UserId:       strconv.FormatUint(r.UserID, 10),
			ReactionType: r.ReactionType,
		})
		if err != nil {
			c.l.Errorf("failure getting reactions response: %v\n", err)
			return err
		}
	}
	return nil
}

/*func (c *Content) GetComments(r *prcontent.GetCommentsRequest, stream prcontent.Content_GetCommentsServer) error {
	comments, err := c.db.GetCommentByPostId(r.PostId)
	if err != nil {
		c.l.Errorf("failure getting comments: %v\n", err)
		return err
	}
	for _, cm := range *comments {
		comment := &prcontent.Comment{
			Content: cm.Content,
			UserId:  cm.UserID,
			PostId:  cm.PostID,
		}
		err = stream.Send(&prcontent.GetCommentsResponse{
			Comment: comment,
		})
		if err != nil {
			c.l.Errorf("failure getting comments response: %v\n", err)
			return err
		}
	}
	return nil
}*/

func (c *Content) GetComments(ctx context.Context, r *prcontent.GetCommentsRequest) (*prcontent.GetCommentsResponse, error) {
	comments, err := c.db.GetCommentByPostId(r.PostId)
	if err != nil {
		c.l.Errorf("failure getting comments: %v\n", err)
		return &prcontent.GetCommentsResponse{}, err
	}

	coms := []*prcontent.Comment{}
	for _, cm := range *comments {
		coms = append(coms, &prcontent.Comment{
			Content: cm.Content,
			UserId:  cm.UserID,
			PostId:  cm.PostID,
		})
	}
	return &prcontent.GetCommentsResponse{Comment: coms}, nil
}

// func (c *Content) GetStories(r *prcontent.GetStoryRequest, stream prcontent.Content_GetStoriesServer) error {
// 	stories, err := c.db.GetStoryByUser(r.UserId)
// 	if err != nil {
// 		c.l.Errorf("failure getting user stories: %v\n", err)
// 		return err
// 	}
// 	for _, p := range *stories {
// 		media := []*prcontent.Media{}
// 		for _, m := range p.SharedMedia.Media {
// 			media = append(media, &prcontent.Media{
// 				Filename:    m.Filename,
// 				Description: m.Description,
// 				AddedOn:     m.AddedOn,
// 				Location: &prcontent.Location{
// 					Country: m.Location.Country,
// 					State:   m.Location.State,
// 					ZipCode: m.Location.ZipCode,
// 					Street:  m.Location.Street,
// 				},
// 				Url: m.URL,
// 			})
// 		}
// 		story := &prcontent.Story{
// 			Id:     p.ID,
// 			UserId: p.UserID,
// 			SharedMedia: &prcontent.SharedMedia{
// 				Media: media,
// 			},
// 			CloseFriends: p.CloseFriends,
// 		}
// 		err = stream.Send(&prcontent.GetStoriesResponse{
// 			Story: story,
// 		})
// 		if err != nil {
// 			c.l.Errorf("failed getting story response: %v\n", err)
// 			return err
// 		}
// 	}
// 	return nil
// }

func (c *Content) PutReaction(ctx context.Context, r *prcontent.PutReactionRequest) (*prcontent.PutReactionResponse, error) {

	i, err := strconv.ParseUint(r.Id, 10, 64)
	if err != nil {
		c.l.Errorf("failure converting reaction id: %v\n", err)
		return &prcontent.PutReactionResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	err = data.PutReaction(c.db, r.ReactionType, i)
	if err != nil {
		c.l.Errorf("failure updating reaction: %v\n", err)
		return &prcontent.PutReactionResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prcontent.PutReactionResponse{}, nil
}

func (c *Content) SearchContent(ctx context.Context, r *prcontent.SearchContentRequest) (*prcontent.SearchContentResponse, error) {
	posts, err := c.db.GetPostsByTag(r.Value)
	if err != nil {
		c.l.Errorf("failure getting posts: %v\n", err)
		return &prcontent.SearchContentResponse{}, err
	}

	retVal := []*prcontent.Post{}

	for _, p := range *posts {
		media := []*prcontent.Media{}
		for _, m := range p.SharedMedia.Media {
			tags := []*prcontent.Tag{}
			for _, t := range m.Tags {
				tags = append(tags, &prcontent.Tag{Id: t.ID, Value: t.Value})
			}
			userTags := []*prcontent.UserTag{}
			for _, t := range m.TaggedUsers {
				userTags = append(userTags, &prcontent.UserTag{Id: t.UserID})
			}
			media = append(media, &prcontent.Media{
				Filename:    m.Filename,
				Description: m.Description,
				AddedOn:     m.AddedOn,
				Tags:        tags,
				Location: &prcontent.Location{
					Country: m.Location.Country,
					State:   m.Location.State,
					ZipCode: m.Location.ZipCode,
					City:    m.Location.City,
					Street:  m.Location.Street,
					Name:    m.Location.Name,
				},
				Url:      m.URL,
				MimeType: prcontent.EMimeType(m.MimeType),
				UserTags: userTags,
			})
		}
		post := &prcontent.Post{
			Id:     strconv.FormatUint(p.ID, 10),
			UserId: strconv.FormatUint(p.UserID, 10),
			SharedMedia: &prcontent.SharedMedia{
				Media: media,
			},
		}

		retVal = append(retVal, post)

		if err != nil {
			c.l.Errorf("failed sending post response: %v\n", err)
			return &prcontent.SearchContentResponse{}, err
		}
	}
	return &prcontent.SearchContentResponse{Post: retVal}, nil
}

func (c *Content) GetTagsByName(ctx context.Context, r *prcontent.GetTagsByNameRequest) (*prcontent.GetTagsByNameResponse, error) {
	tags, err := c.db.GetAllTagsByNameSubstring(r.Query)
	if err != nil {
		c.l.Printf("[ERROR] geting tag: %v\n", err)
		return &prcontent.GetTagsByNameResponse{}, err
	}

	var tagNames []string
	for _, tagName := range tags {
		tagNames = append(tagNames, tagName.Value)
	}

	return &prcontent.GetTagsByNameResponse{Name: tagNames}, nil
}

func (c *Content) GetLocationNames(ctx context.Context, r *prcontent.GetLocationNamesRequest) (*prcontent.GetLocationNamesResponse, error) {
	names, err := c.db.GetAllLocationNames(r.Query)
	if err != nil {
		c.l.Printf("[ERROR] geting location names: %v\n", err)
		return &prcontent.GetLocationNamesResponse{}, err
	}

	return &prcontent.GetLocationNamesResponse{Name: names}, nil
}

func (c *Content) SearchContentLocation(ctx context.Context, r *prcontent.SearchContentLocationRequest) (*prcontent.SearchContentLocationResponse, error) {
	posts, err := c.db.GetContentsByLocation(r.Name)
	if err != nil {
		c.l.Errorf("failure getting posts: %v\n", err)
		return &prcontent.SearchContentLocationResponse{}, err
	}

	retVal := []*prcontent.Post{}

	for _, p := range *posts {
		media := []*prcontent.Media{}
		for _, m := range p.SharedMedia.Media {
			tags := []*prcontent.Tag{}
			for _, t := range m.Tags {
				tags = append(tags, &prcontent.Tag{Id: t.ID, Value: t.Value})
			}
			userTags := []*prcontent.UserTag{}
			for _, t := range m.TaggedUsers {
				userTags = append(userTags, &prcontent.UserTag{Id: t.UserID})
			}
			media = append(media, &prcontent.Media{
				Filename:    m.Filename,
				Description: m.Description,
				AddedOn:     m.AddedOn,
				Tags:        tags,
				Location: &prcontent.Location{
					Country: m.Location.Country,
					State:   m.Location.State,
					ZipCode: m.Location.ZipCode,
					City:    m.Location.City,
					Street:  m.Location.Street,
					Name:    m.Location.Name,
				},
				Url:      m.URL,
				MimeType: prcontent.EMimeType(m.MimeType),
				UserTags: userTags,
			})
		}
		post := &prcontent.Post{
			Id:     strconv.FormatUint(p.ID, 10),
			UserId: strconv.FormatUint(p.UserID, 10),
			SharedMedia: &prcontent.SharedMedia{
				Media: media,
			},
		}

		retVal = append(retVal, post)

		if err != nil {
			c.l.Errorf("failed sending post response: %v\n", err)
			return &prcontent.SearchContentLocationResponse{}, err
		}
	}
	return &prcontent.SearchContentLocationResponse{Post: retVal}, nil
}

func (c *Content) SavePost(ctx context.Context, r *prcontent.SavePostRequest) (*prcontent.SavePostResponse, error) {
	savedPost := &data.SavedPost{
		UserID: r.UserId,
		PostID: r.PostId,
	}
	err := c.db.AddSavedPost(savedPost)
	if err != nil {
		c.l.Errorf("failed to save post: %v", err)
		return &prcontent.SavePostResponse{}, status.Error(codes.Internal, "Internal error")
	}

	return &prcontent.SavePostResponse{}, nil
}

func (c *Content) GetSavedPosts(ctx context.Context, r *prcontent.GetSavedPostsRequest) (*prcontent.GetSavedPostsResponse, error) {
	posts, err := c.db.GetSavedPosts(r.UserId)
	if err != nil {
		c.l.Errorf("failure getting posts: %v\n", err)
		return &prcontent.GetSavedPostsResponse{}, err
	}

	retVal := []*prcontent.Post{}

	for _, p := range *posts {
		media := []*prcontent.Media{}
		for _, m := range p.SharedMedia.Media {
			tags := []*prcontent.Tag{}
			for _, t := range m.Tags {
				tags = append(tags, &prcontent.Tag{Id: t.ID, Value: t.Value})
			}
			userTags := []*prcontent.UserTag{}
			for _, t := range m.TaggedUsers {
				userTags = append(userTags, &prcontent.UserTag{Id: t.UserID})
			}
			media = append(media, &prcontent.Media{
				Filename:    m.Filename,
				Description: m.Description,
				AddedOn:     m.AddedOn,
				Tags:        tags,
				Location: &prcontent.Location{
					Country: m.Location.Country,
					State:   m.Location.State,
					ZipCode: m.Location.ZipCode,
					City:    m.Location.City,
					Street:  m.Location.Street,
					Name:    m.Location.Name,
				},
				Url:      m.URL,
				MimeType: prcontent.EMimeType(m.MimeType),
				UserTags: userTags,
			})
		}
		post := &prcontent.Post{
			Id:     strconv.FormatUint(p.ID, 10),
			UserId: strconv.FormatUint(p.UserID, 10),
			SharedMedia: &prcontent.SharedMedia{
				Media: media,
			},
		}

		retVal = append(retVal, post)

		if err != nil {
			c.l.Errorf("failed sending post response: %v\n", err)
			return &prcontent.GetSavedPostsResponse{}, err
		}
	}
	return &prcontent.GetSavedPostsResponse{Post: retVal}, nil
}

func (c *Content) GetTaggedPosts(ctx context.Context, r *prcontent.GetTaggedPostsRequest) (*prcontent.GetTaggedPostsResponse, error) {
	posts, err := c.db.GetTaggedPostsByUser(r.UserId)
	if err != nil {
		c.l.Errorf("failure getting posts: %v\n", err)
		return &prcontent.GetTaggedPostsResponse{}, err
	}

	retVal := []*prcontent.Post{}

	for _, p := range *posts {
		media := []*prcontent.Media{}
		for _, m := range p.SharedMedia.Media {
			tags := []*prcontent.Tag{}
			for _, t := range m.Tags {
				tags = append(tags, &prcontent.Tag{Id: t.ID, Value: t.Value})
			}
			userTags := []*prcontent.UserTag{}
			for _, t := range m.TaggedUsers {
				userTags = append(userTags, &prcontent.UserTag{Id: t.UserID})
			}
			media = append(media, &prcontent.Media{
				Filename:    m.Filename,
				Description: m.Description,
				AddedOn:     m.AddedOn,
				Tags:        tags,
				Location: &prcontent.Location{
					Country: m.Location.Country,
					State:   m.Location.State,
					ZipCode: m.Location.ZipCode,
					City:    m.Location.City,
					Street:  m.Location.Street,
					Name:    m.Location.Name,
				},
				Url:      m.URL,
				MimeType: prcontent.EMimeType(m.MimeType),
				UserTags: userTags,
			})
		}
		post := &prcontent.Post{
			Id:     strconv.FormatUint(p.ID, 10),
			UserId: strconv.FormatUint(p.UserID, 10),
			SharedMedia: &prcontent.SharedMedia{
				Media: media,
			},
		}

		retVal = append(retVal, post)

		if err != nil {
			c.l.Errorf("failed sending post response: %v\n", err)
			return &prcontent.GetTaggedPostsResponse{}, err
		}
	}
	return &prcontent.GetTaggedPostsResponse{Post: retVal}, nil
}
