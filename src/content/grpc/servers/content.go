package servers

import (
	"context"
	"saltgram/content/data"
	"saltgram/protos/content/prcontent"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Content struct {
	prcontent.UnimplementedContentServer
	l  *logrus.Logger
	db *data.DBConn
}

func NewContent(l *logrus.Logger, db *data.DBConn) *Content {
	return &Content{
		l:  l,
		db: db,
	}
}

func (s *Content) GetSharedMedia(r *prcontent.SharedMediaRequest, stream prcontent.Content_GetSharedMediaServer) error {
	sharedMedias, err := s.db.GetSharedMediaByUser(r.UserId)
	if err != nil {
		s.l.Errorf("failure getting user shared media: %v\n", err)
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
			s.l.Errorf("failure sending shared media response: %v\n", err)
			return err
		}
	}
	return nil
}

func (c *Content) getProflePicture(ctx context.Context, r *prcontent.GetProfilePictureRequest) (*prcontent.GetProfilePictureResponse, error) {
	profilePicture, err := c.db.GetProfilePictureByUser(r.UserId)
	if err != nil {
		return &prcontent.GetProfilePictureResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prcontent.GetProfilePictureResponse{

		Media: &prcontent.Media{
			Filename:    profilePicture.Media.Filename,
			Description: profilePicture.Media.Description,
			AddedOn:     profilePicture.Media.AddedOn,
			Location: &prcontent.Location{
				Country: profilePicture.Media.Location.Country,
				State:   profilePicture.Media.Location.State,
				ZipCode: profilePicture.Media.Location.ZipCode,
				Street:  profilePicture.Media.Location.Street,
			},
		},
		UserId: profilePicture.UserID,
		Id:     profilePicture.ID,
	}, nil
}

func (c *Content) GetPostsByUser(r *prcontent.GetPostsRequest, stream prcontent.Content_GetPostsServer) error {
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
		post := &prcontent.Post {
			Id: p.ID,
			UserId: p.UserID,
			SharedMedia: &prcontent.SharedMedia {
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

func (c *Content) AddProflePicture(ctx context.Context, r *prcontent.AddProfilePictureRequest) (*prcontent.AddProfilePictureResponse, error) {

	profilePicture := data.ProfilePicture{
		UserID: r.UserId,
		Media: data.Media{
			Filename:    r.Media.Filename,
			Description: r.Media.Description,
			AddedOn:     r.Media.AddedOn,
			Location: data.Location{
				Country: r.Media.Location.Country,
				State:   r.Media.Location.State,
				ZipCode: r.Media.Location.ZipCode,
				Street:  r.Media.Location.Street,
			},
		},
	}

	err := c.db.AddProfilePicture(&profilePicture)
	if err != nil {
		c.l.Println("[ERROR] adding profile picture ", err)
	}

	return &prcontent.AddProfilePictureResponse{}, nil
}

func (u *Content) AddSharedMedia(ctx context.Context, r *prcontent.AddSharedMediaRequest) (*prcontent.AddSharedMediaResponse, error) {

	media := []*data.Media{}
	for _, m := range r.Media {
		media = append(media, &data.Media{
			Filename:    m.Filename,
			Tags:        []data.Tag{}, // TODO
			Description: m.Description,
			Location: data.Location{
				Country: m.Location.Country,
				State:   m.Location.State,
				ZipCode: m.Location.ZipCode,
				Street:  m.Location.Street,
			},
			AddedOn: m.AddedOn,
		})
	}
	sharedMedia := data.SharedMedia{
		Media: media,
	}

	err := u.db.AddSharedMedia(&sharedMedia)
	if err != nil {
		u.l.Errorf("failure adding user: %v\n", err)
		return &prcontent.AddSharedMediaResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prcontent.AddSharedMediaResponse{}, nil
}

func (u *Content) AddPost(ctx context.Context, r *prcontent.AddPostRequest) (*prcontent.AddPostResponse, error) {

	media := []*data.Media{}
	for _, m := range r.SharedMedia.Media {
		media = append(media, &data.Media{
			Filename:    m.Filename,
			Tags:        []data.Tag{},
			Description: m.Description,
			Location: data.Location{
				Country: m.Location.Country,
				State:   m.Location.State,
				ZipCode: m.Location.ZipCode,
				Street:  m.Location.Street,
			},
			AddedOn: m.AddedOn,
		})
	}
	post := data.Post{
		SharedMedia: data.SharedMedia{
			Media: media,
		},
		UserID: r.UserId,
	}

	err := u.db.AddPost(&post)
	if err != nil {
		u.l.Errorf("failure adding post: %v\n", err)
		return &prcontent.AddPostResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prcontent.AddPostResponse{}, nil
}

func (c *Content) AddComment(ctx context.Context, r *prcontent.AddCommentRequest) (*prcontent.AddCommentResponse, error) {

	comment := data.Comment{
		Content: r.Content,
		UserID: r.UserId,
		PostID: r.PostId,
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
		//ReactionType: r.ReactionType,
		UserID: r.UserId,
		PostID: r.PostId,
	}

	err := c.db.AddReaction(&reaction)
	if err != nil {
		c.l.Errorf("Failed adding reaction: %v\n", err)
		return &prcontent.AddReactionResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prcontent.AddReactionResponse{}, nil
}

