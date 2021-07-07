package saga

import (
	"encoding/json"
	"saltgram/auth/data"

	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type redisClient struct {
	l  *logrus.Logger
	db *data.DBConn
}

func NewRedisClient(l *logrus.Logger, db *data.DBConn) *redisClient {
	return &redisClient{
		l:  l,
		db: db,
	}
}

func (rc *redisClient) Connection() {
	var err error
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	if _, err = client.Ping().Result(); err != nil {
		rc.l.Fatalf("creating redis client %s", err)
	}

	pubsub := client.Subscribe(ReplyChannel, AuthChannel)
	if _, err = pubsub.Receive(); err != nil {
		rc.l.Fatalf("subscribing %s", err)
	}

	defer func() { _ = pubsub.Close() }()
	ch := pubsub.Channel()

	rc.l.Info("starting the auth service")
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
			case AuthChannel:
				if m.Action == ActionStart {
					refreshClaims := data.RefreshClaims{
						Username: m.Username,
						StandardClaims: jwt.StandardClaims{
							// TODO(Jovan): Make programmatic?
							ExpiresAt: time.Now().UTC().AddDate(0, 6, 0).Unix(),
							Issuer:    "SaltGram",
						},
					}
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
					jws, err := token.SignedString([]byte(os.Getenv("REF_SECRET_KEY")))
					if err != nil {
						rc.l.Errorf("failure signing refresh token")
						sendToReplyChannel(client, &m, ActionError, UserService, AuthService, rc.l)
					}
					err = data.AddRefreshToken(rc.db, m.Username, jws)
					if err != nil {
						rc.l.Errorf("adding refresh token: %v\n", err)
						sendToReplyChannel(client, &m, ActionError, UserService, AuthService, rc.l)
					}
					sendToReplyChannel(client, &m, ActionDone, ContentService, AuthService, rc.l)
				}

				if m.Action == ActionRollback {
					data.DeleteRefreshToken(rc.db, m.Username)
					if m.SenderService == ContentService {
						sendToReplyChannel(client, &m, ActionError, UserService, AuthService, rc.l)
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
