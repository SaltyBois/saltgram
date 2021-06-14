package gdrive

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"saltgram/internal"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GDrive struct {
	s   *drive.Service
	l   *logrus.Logger
	ctx context.Context
}

var (
	serviceCreds = internal.GetEnvOrDefault("SALT_GDRIVE_CREDS", "../../secrets/saltgram-service-key.json")
)

func NewGDrive(l *logrus.Logger) *GDrive {
	ctx := context.Background()
	gDrive := &GDrive{l: l, ctx: ctx}
	gDrive.getServiceClient()
	return gDrive
}

func (g *GDrive) getServiceClient() {
	b, err := ioutil.ReadFile(serviceCreds)
	if err != nil {
		g.l.Fatalf("failed to load gdrive service credentials: %v\n", err)
	}
	var s = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}
	json.Unmarshal(b, &s)
	config := &jwt.Config{
		Email:      s.Email,
		PrivateKey: []byte(s.PrivateKey),
		Scopes: []string{
			drive.DriveScope,
		},
	}
	client := config.Client(context.Background())
	srv, err := drive.NewService(g.ctx, option.WithHTTPClient(client))
	if err != nil {
		g.l.Fatalf("failed to get gdrive service: %v\n", err)
	}
	g.s = srv
}

func (g *GDrive) CreateFolder(name string, parentIds []string, isPublic bool) (*drive.File, error) {
	f := &drive.File{
		MimeType: "application/vnd.google-apps.folder",
		Name:     name,
		Parents:  parentIds,
	}

	createdFile, err := g.s.Files.Create(f).Do()
	if err != nil {
		return nil, err
	}
	if isPublic {
		_, err := g.s.Permissions.Create(createdFile.Id, &drive.Permission{
			Type: "anyone",
			Role: "reader",
			// AllowFileDiscovery: true, Maybe too much?
		}).Do()
		if err != nil {
			g.l.Errorf("failed to create public permissions for file: %v, error: %v\n", f.Name, err)
			return createdFile, err
		}
	} else {
		_, err := g.s.Permissions.Create(createdFile.Id, &drive.Permission{
			Type:         "user",
			EmailAddress: "bezbednovic@gmail.com",
			Role:         "reader",
		}).Do()
		if err != nil {
			g.l.Errorf("failed to create private permissions for file: %v, error: %v\n", f.Name, err)
			return createdFile, err
		}
	}
	return createdFile, nil
}

func (g *GDrive) CreateFile(name string, parentIds []string, data io.Reader, isPublic bool) (*drive.File, error) {
	f := &drive.File{
		MimeType: "application/vnd.google-apps.photo",
		Name:     name,
		Parents:  parentIds,
	}
	createdFile, err := g.s.Files.Create(f).Media(data).Do()
	if err != nil {
		g.l.Errorf("failed to upload file: %v, error:%v", f.Name, err)
		return nil, err
	}
	if isPublic {
		_, err = g.s.Permissions.Create(createdFile.Id, &drive.Permission{
			Type: "anyone",
			Role: "reader",
		}).Do()
		if err != nil {
			g.l.Errorf("failed to set public permissions for file: %v, error: %v", createdFile.Id, err)
			return createdFile, err
		}
	} else {
		_, err = g.s.Permissions.Create(createdFile.Id, &drive.Permission{
			Type:         "user",
			EmailAddress: "bezbednovic@gmail.com",
			Role:         "reader",
		}).Do()
		if err != nil {
			g.l.Errorf("failed to create private permissions for file: %v, error: %v", createdFile.Id, err)
			return createdFile, err
		}
	}
	return createdFile, nil
}

func (g *GDrive) DeleteFile(fileId string) error {
	return g.s.Files.Delete(fileId).Do()
}

func (g *GDrive) GetFiles() ([]*drive.File, error) {
	r, err := g.s.Files.List().Fields("files(id, name)").Do()
	return r.Files, err
}