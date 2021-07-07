package saga

import (
	"encoding/json"
	"saltgram/content/data"
	"saltgram/content/gdrive"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type redisClient struct {
	l  *logrus.Logger
	db *data.DBConn
	g  *gdrive.GDrive
}

func NerRedisClient(l *logrus.Logger, db *data.DBConn, g *gdrive.GDrive) *redisClient {
	return &redisClient{
		l:  l,
		db: db,
		g:  g,
	}
}

func (rc *redisClient) Connection() {
	var err error
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	if _, err = client.Ping().Result(); err != nil {
		rc.l.Fatalf("creating redis client %s", err)
	}

	pubsub := client.Subscribe(ReplyChannel, ContentChannel)
	if _, err = pubsub.Receive(); err != nil {
		rc.l.Fatalf("subscribing %s", err)
	}

	defer func() { _ = pubsub.Close() }()
	ch := pubsub.Channel()

	rc.l.Info("starting the content service")
	for {
		select {
		case msg := <-ch:
			m := Message{}
			err := json.Unmarshal([]byte(msg.Payload), &m)
			if err != nil {
				rc.l.Errorf("bad message")
				continue
			}

			switch msg.Channel {
			case ContentChannel:
				if m.Action == ActionStart {
					userId := strconv.FormatUint(m.UserId, 10)
					profile, posts, stories, err := rc.g.CreateUserFolder(userId)
					if err != nil {
						rc.l.Errorf("failed to create user folders: %v", err)
						sendToReplyChannel(client, &m, ActionError, AuthService, ContentService, rc.l)
					}
					m.ProfileFolderId = profile
					m.PostsFolderId = posts
					m.StoriesFolderId = stories
					sendToReplyChannel(client, &m, ActionDone, UserService, ContentService, rc.l)
				}
				if m.Action == ActionRollback {
					rc.g.DeleteFile(m.ProfileFolderId)
					rc.g.DeleteFile(m.PostsFolderId)
					rc.g.DeleteFile(m.StoriesFolderId)
					if m.SenderService == UserService {
						sendToReplyChannel(client, &m, ActionError, AuthService, ContentService, rc.l)
					}
				}
			}
		}
	}
}

func sendToReplyChannel(client *redis.Client, m *Message, action string, service string, senderService string, l *logrus.Logger) {
	var err error
	m.Action = action
	m.Service = service
	m.SenderService = senderService
	if err = client.Publish(ReplyChannel, m).Err(); err != nil {
		l.Errorf("publishing done-message to %s channel", ReplyChannel)
	}
}
