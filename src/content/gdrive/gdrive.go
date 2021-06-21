package gdrive

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"saltgram/content/data"
	"saltgram/internal"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2/google"
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
	serviceCreds  = internal.GetEnvOrDefault("SALT_GDRIVE_CREDS", "../../secrets/saltgram-service-key.json")
	publicId      string
	profilesId    string
	publicBaseUrl = "https://drive.google.com/uc?export=view&id="
)

func NewGDrive(l *logrus.Logger) *GDrive {
	ctx := context.Background()
	gDrive := &GDrive{l: l, ctx: ctx}
	gDrive.getServiceClient()
	gDrive.initFolders()
	return gDrive
}

func (g *GDrive) initFolders() {
	folders, err := g.QueryFiles("name='public'")
	if err != nil {
		g.l.Fatalf("failed to query public folder: %v", err)
	}
	if len(folders) == 0 {
		public, err := g.CreateFolder("public", []string{"root"}, true)
		if err != nil {
			g.l.Fatalf("failed to create public folder: %v", err)
		}
		publicId = public.Id
	} else {
		publicId = folders[0].Id
	}
	folders, err = g.QueryFiles("name='profiles' and '" + publicId + "' in parents")
	if err != nil {
		g.l.Fatalf("failed to query profiles folder: %v", err)
	}
	if len(folders) == 0 {
		profile, err := g.CreateFolder("profiles", []string{publicId}, true)
		if err != nil {
			g.l.Fatalf("failedto create profiles folder: %v", err)
		}
		profilesId = profile.Id
	} else {
		profilesId = folders[0].Id
	}
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
		TokenURL: google.JWTTokenURL,
	}

	client := config.Client(context.Background())
	srv, err := drive.NewService(g.ctx, option.WithHTTPClient(client))
	if err != nil {
		g.l.Fatalf("failed to get gdrive service: %v\n", err)
	}
	g.s = srv
}

func (g *GDrive) CreateUserFolder(userId string) (string, string, string, error) {
	userFolder, err := g.CreateFolder(userId, []string{profilesId}, true)
	if err != nil {
		g.l.Errorf("failed to create user folder: %v", err)
		return "", "", "", err
	}
	postsFolder, err := g.CreateFolder("posts", []string{userFolder.Id}, true)
	if err != nil {
		g.l.Errorf("failed to create posts folder: %v", err)
		return "", "", "", err
	}
	storyFolder, err := g.CreateFolder("stories", []string{userFolder.Id}, true)
	if err != nil {
		g.l.Errorf("failed to create stories folder: %v", err)
		return "", "", "", err
	}
	return userFolder.Id, postsFolder.Id, storyFolder.Id, nil
}

func (g *GDrive) UploadProfilePicture(profileFolderId string, data io.Reader) (string, error) {
	profileImg, err := g.QueryFiles("name='profile' and '" + profileFolderId + "' in parents")
	if err != nil {
		g.l.Errorf("failed to query profile folder: %v", err)
		return "", err
	}
	var profileImgId string
	if len(profileImg) == 0 {
		profile, _, err := g.CreateFile("profile", []string{profileFolderId}, data, true)
		if err != nil {
			g.l.Errorf("failed to create user profile: %v", err)
			return "", err
		}
		profileImgId = profile.Id
	} else {
		f := &drive.File{
			Name:    profileImg[0].Name,
			Parents: profileImg[0].Parents,
		}
		profile, err := g.s.Files.Update(profileImg[0].Id, f).Media(data).Do()
		if err != nil {
			g.l.Errorf("failed to update profile picture: %v", err)
			return "", err
		}
		profileImgId = profile.Id
	}
	return publicBaseUrl + profileImgId, nil
}

func (g *GDrive) UploadPost(postsFolderId, filename string, data io.Reader) (string, data.EMimeType, error) {
	post, mimeType, err := g.CreateFile(filename, []string{postsFolderId}, data, true)
	if err != nil {
		g.l.Errorf("failed uploading post: %v", err)
		return "", mimeType, err
	}
	return publicBaseUrl + post.Id, mimeType, nil
}

func (g *GDrive) UploadStory(storiesFolderId, filename string, data io.Reader) (string, data.EMimeType, error) {
	story, mimeType, err := g.CreateFile(filename, []string{storiesFolderId}, data, true)
	if err != nil {
		g.l.Errorf("failed uploading story: %v", err)
		return "", mimeType, err
	}
	return publicBaseUrl + story.Id, mimeType, nil
}

func (g *GDrive) CreateFolder(name string, parentIds []string, isPublic bool) (*drive.File, error) {
	f := &drive.File{
		MimeType: "application/vnd.google-apps.folder",
		Name:     name,
		Parents:  parentIds,
	}

	createdFolder, err := g.s.Files.Create(f).Do()
	if err != nil {
		return nil, err
	}

	_, err = g.s.Permissions.Create(createdFolder.Id, &drive.Permission{
		Type:         "user",
		EmailAddress: "bezbednovic@gmail.com",
		Role:         "reader",
	}).Do()
	if err != nil {
		g.l.Errorf("failed to create bezbednovic permissions for folder: %v, error: %v\n", f.Name, err)
		return createdFolder, err
	}

	if isPublic {
		_, err := g.s.Permissions.Create(createdFolder.Id, &drive.Permission{
			Type: "anyone",
			Role: "reader",
			// AllowFileDiscovery: true, Maybe too much?
		}).Do()
		if err != nil {
			g.l.Errorf("failed to create public permissions for folder: %v, error: %v\n", f.Name, err)
			return createdFolder, err
		}
	}
	return createdFolder, nil
}

func (g *GDrive) CreateFile(name string, parentIds []string, content io.Reader, isPublic bool) (*drive.File, data.EMimeType, error) {
	mimeType := data.EMimeType(data.EMimeType_IMAGE)
	f := &drive.File{
		Name:    name,
		Parents: parentIds,
	}
	createdFile, err := g.s.Files.Create(f).Media(content).Do()
	if err != nil {
		g.l.Errorf("failed to upload file: %v, error:%v", f.Name, err)
		return nil, mimeType, err
	}

	if strings.Contains(createdFile.MimeType, "video") {
		mimeType = data.EMimeType_VIDEO
	}

	_, err = g.s.Permissions.Create(createdFile.Id, &drive.Permission{
		Type:         "user",
		EmailAddress: "bezbednovic@gmail.com",
		Role:         "reader",
	}).Fields("id").Do()

	if err != nil {
		g.l.Errorf("failed to create bezbednovic permissions for folder: %v, error: %v", createdFile.Name, err)
		return createdFile, mimeType, err
	}

	if isPublic {
		_, err = g.s.Permissions.Create(createdFile.Id, &drive.Permission{
			Type: "anyone",
			Role: "reader",
		}).Do()
		if err != nil {
			g.l.Errorf("failed to set public permissions for file: %v, error: %v", createdFile.Id, err)
			return createdFile, mimeType, err
		}
	}

	if err != nil {
		g.l.Errorf("failed to create private permissions for file: %v, error: %v", createdFile.Id, err)
		return createdFile, mimeType, err
	}

	return createdFile, mimeType, nil
}

func (g *GDrive) DeleteFile(fileId string) error {
	return g.s.Files.Delete(fileId).Do()
}

func (g *GDrive) QueryFiles(query string) ([]*drive.File, error) {
	r, err := g.s.Files.List().Q(query).Do()
	return r.Files, err
}

func (g *GDrive) GetFiles() ([]*drive.File, error) {
	r, err := g.s.Files.List().Fields("files(id, name, parents)").Do()
	return r.Files, err
}
