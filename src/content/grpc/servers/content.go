package servers

import (
	"context"
	"log"
	"saltgram/content/data"
	"saltgram/protos/content/prcontent"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Content struct {
	prcontent.UnimplementedContentServer
	l  *log.Logger
	db *data.DBConn
}

func NewContent(l *log.Logger, db *data.DBConn) *Content {
	return &Content{
		l:  l,
		db: db,
	}
}

func (s *Content) getSharedMedia(r *prcontent.SharedMediaRequest, stream prcontent.Content_GetSharedMediaServer) error {
	sharedMedias, err := s.db.GetSharedMediaByUser(r.UserId)
	if err != nil {
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
		u.l.Printf("[ERROR] adding user: %v\n", err)
		return &prcontent.AddSharedMediaResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prcontent.AddSharedMediaResponse{}, nil
}
