package saga

import (
	"context"
	"encoding/json"
	"saltgram/users/data"
	"time"

	"saltgram/protos/email/premail"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type RedisClient struct {
	l  *logrus.Logger
	db *data.DBConn
	ec premail.EmailClient
	o  *Orchestrator
}

func NerRedisClient(l *logrus.Logger, db *data.DBConn, ec premail.EmailClient) *RedisClient {
	return &RedisClient{
		o:  NewOrchestrator(l),
		l:  l,
		db: db,
		ec: ec,
	}
}

func (rc *RedisClient) Start() {
	go rc.o.Start()
}

func (rc *RedisClient) Next(channel string, service string, message Message) {
	rc.o.Next(channel, service, message)
}

func (rc *RedisClient) Connection() {
	var err error
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	if _, err = client.Ping().Result(); err != nil {
		rc.l.Fatalf("creating redis client %s", err)
	}

	pubsub := client.Subscribe(ReplyChannel, UserChannel)
	if _, err = pubsub.Receive(); err != nil {
		rc.l.Fatalf("subscribing %s", err)
	}

	defer func() { _ = pubsub.Close() }()
	ch := pubsub.Channel()

	rc.l.Info("starting the user service")
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
			case UserChannel:
				if m.Action == ActionStart {
					ok := true
					profile := data.Profile{
						Username:        m.Username,
						UserID:          m.UserId,
						Taggable:        true,
						Public:          false,
						Description:     m.Description,
						PhoneNumber:     m.PhoneNumber,
						Gender:          m.Gender,
						DateOfBirth:     time.Unix(m.DateOfBirth, 0),
						WebSite:         m.WebSite,
						PrivateProfile:  m.PrivateProfile,
						ProfileFolderId: m.ProfileFolderId,
						PostsFolderId:   m.PostsFolderId,
						StoriesFolderId: m.StoriesFolderId,
						Messagable:      true,
						Verified:        false,
						AccountType:     "",
					}

					err = rc.db.AddProfile(&profile)
					if err != nil {
						rc.l.Errorf("adding profile: %v\n", err)
						ok = false
					}
					if ok {
						go func() {
							_, err := rc.ec.SendActivation(context.Background(), &premail.SendActivationRequest{Email: m.Email})
							if err != nil {
								rc.l.Errorf("failure sending activation request: %v\n", err)
							}
						}()
					} else {
						sendToReplyChannel(client, &m, ActionError, ContentService, UserService, rc.l)
					}
				}
				if m.Action == ActionRollback {
					rc.db.DeleteUser(m.Username)
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
