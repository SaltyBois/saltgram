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
