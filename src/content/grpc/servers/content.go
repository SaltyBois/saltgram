package servers

import (
	"bytes"
	"context"
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

func (c *Content) PostProfile(ctx context.Context, r *prcontent.PostProfileRequest) (*prcontent.PostProfileResponse, error) {
	if len(r.Image) == 0 {
		c.l.Errorf("bad request, empty image")
		return &prcontent.PostProfileResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}
	f, err := c.g.CreateFile("profile", []string{"root"}, bytes.NewReader(r.Image), true)
	if err != nil {
		c.l.Errorf("failed to upload profile image: %v", err)
		return &prcontent.PostProfileResponse{}, status.Error(codes.Internal, "Internal server error")
	}
	return &prcontent.PostProfileResponse{Url: f.WebViewLink}, nil
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

func (c *Content) AddSharedMedia(ctx context.Context, r *prcontent.AddSharedMediaRequest) (*prcontent.AddSharedMediaResponse, error) {

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

	err := c.db.AddSharedMedia(&sharedMedia)
	if err != nil {
		c.l.Errorf("failure adding user: %v\n", err)
		return &prcontent.AddSharedMediaResponse{}, status.Error(codes.InvalidArgument, "Bad request")
	}

	return &prcontent.AddSharedMediaResponse{}, nil
}
