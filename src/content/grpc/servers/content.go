package servers

import (
	"bytes"
	"context"
	"io"
	"saltgram/content/data"
	"saltgram/content/gdrive"
	"saltgram/protos/content/prcontent"

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

func (c *Content) GetStoriesIndividual(ctx context.Context, r *prcontent.GetStoriesIndividualRequest) (*prcontent.GetStoriesIndividualResponse, error) {
	stories, err := c.db.GetStoriesByUserAsMedia(r.UserId)
	if err != nil {
		c.l.Errorf("failed to get users stories: %v", err)
		return &prcontent.GetStoriesIndividualResponse{}, status.Error(codes.Internal, "Internal error")
	}

	s := []*prcontent.Media{}

	for _, st := range stories {
		s = append(s, data.DataToPRMedia(st))
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
		Name: r.Name,
		UserID: r.UserId,
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
	post := &data.Post{
		UserID: r.UserId,
		SharedMedia: data.SharedMedia{},
	}
	err := c.db.AddPost(post)
	if err != nil {
		c.l.Errorf("failed to create post: %v", err)
		return &prcontent.CreatePostResponse{}, status.Error(codes.Internal, "Internal error")
	}

	return &prcontent.CreatePostResponse{PostId: post.ID}, nil
}

func (c *Content) CreateStory(ctx context.Context, r *prcontent.CreateStoryRequest) (*prcontent.CreateStoryResponse, error) {
	story := &data.Story{
		UserID: r.UserId,
		SharedMedia: data.SharedMedia{},
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

func (c *Content) GetPostsByUser(r *prcontent.GetPostsRequest, stream prcontent.Content_GetPostsByUserServer) error {
	posts, err := c.db.GetPostByUser(r.UserId)
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
			})
		}
		post := &prcontent.Post{
			Id:     p.ID,
			UserId: p.UserID,
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
}

func (c *Content) GetPostsByUserReaction(r *prcontent.GetPostsRequest, stream prcontent.Content_GetPostsByUserReactionServer) error {
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
			})
		}
		post := &prcontent.Post{
			Id:     p.ID,
			UserId: p.UserID,
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
		tags = append(tags, data.Tag{
			Value: tag.Value,
		})
	}
	location := r.GetInfo().Media.Location

	media := &data.Media{
		Filename:      r.GetInfo().Media.Filename,
		Tags:          tags,
		Description:   r.GetInfo().Media.Description,
		Location: data.Location{
			Country: location.Country,
			State:   location.State,
			ZipCode: location.ZipCode,
			Street:  location.Street,
		},
		AddedOn: r.GetInfo().Media.AddedOn,
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
	url, err := c.g.UploadStory(r.GetInfo().StoriesFolderId, media.Filename, &imageData)
	if err != nil {
		c.l.Errorf("failed to upload story: %v", err)
		return status.Error(codes.Internal, "Internal error")
	}
	media.URL = url
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
		tags = append(tags, data.Tag{
			Value: tag.Value,
		})
	}

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
			Street:  location.Street,
		},
		AddedOn: r.GetInfo().Media.AddedOn,
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

	url, err := c.g.UploadPost(r.GetInfo().PostsFolderId, media.Filename, &imageData)
	if err != nil {
		c.l.Errorf("failure to upload post: %v", err)
		return status.Error(codes.Internal, "Internal error")
	}

	media.URL = url
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

func (c *Content) GetReactionByPost(r *prcontent.GetReactionsRequest, stream prcontent.Content_GetReactionsServer) error {
	reactions, err := c.db.GetReactionByPostId(r.PostId)
	if err != nil {
		c.l.Errorf("failure getting reactions: %v\n", err)
		return err
	}
	for _, r := range *reactions {
		err = stream.Send(&prcontent.GetReactionsResponse{
			Id:           r.ID,
			UserId:       r.UserID,
			ReactionType: r.ReactionType,
		})
		if err != nil {
			c.l.Errorf("failure getting reactions response: %v\n", err)
			return err
		}
	}
	return nil
}
