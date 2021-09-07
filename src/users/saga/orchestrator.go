package saga

import (
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

const (
	UserChannel    string = "UserChannel"
	AuthChannel    string = "AuthChannel"
	ContentChannel string = "ContentChannel"
	ReplyChannel   string = "ReplyChannel"
	UserService    string = "User"
	AuthService    string = "Auth"
	ContentService string = "Content"
	ActionStart    string = "Start"
	ActionDone     string = "Done"
	ActionError    string = "Error"
	ActionRollback string = "Rollback"
)

type Orchestrator struct {
	c *redis.Client
	r *redis.PubSub
	l *logrus.Logger
}

func NewOrchestrator(l *logrus.Logger) *Orchestrator {
	var err error
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	if _, err = client.Ping().Result(); err != nil {
		l.Error("creating redis client %s", err)
	}

	o := &Orchestrator{
		c: client,
		r: client.Subscribe(UserChannel, AuthChannel, ContentChannel, ReplyChannel),
		l: l,
	}

	return o
}

func (o Orchestrator) Start() {
	var err error
	if _, err = o.r.Receive(); err != nil {
		o.l.Fatalf("error setting up redis %s \n", err)
	}
	ch := o.r.Channel()
	defer func() { _ = o.r.Close() }()

	o.l.Errorf("starting the redis client")
	for {
		select {
		case msg := <-ch:
			m := Message{}
			if err = json.Unmarshal([]byte(msg.Payload), &m); err != nil {
				o.l.Errorf("bad message %s", err)
				continue
			}

			switch msg.Channel {
			case ReplyChannel:
				if m.Action != ActionDone {
					o.Rollback(m)
					continue
				}

				switch m.Service {
				case UserService:
					o.Next(UserChannel, UserService, m)
				case AuthService:
					o.Next(AuthChannel, AuthService, m)
				case ContentService:
					o.Next(ContentChannel, ContentService, m)
				}
			}
		}
	}
}

func (o Orchestrator) Next(channel, service string, message Message) {
	var err error
	message.Action = ActionStart
	message.Service = service
	if err = o.c.Publish(channel, message).Err(); err != nil {
		o.l.Errorf("error publishing start-message to %s channel", channel)
	}
	o.l.Info("start message published to channel :%s", channel)
}

func (o Orchestrator) Rollback(m Message) {
	var err error
	var channel string
	switch m.Service {
	case UserService:
		channel = UserChannel
	case AuthService:
		channel = AuthChannel
	case ContentService:
		channel = ContentChannel
	}
	m.Action = ActionRollback
	if err = o.c.Publish(channel, m).Err(); err != nil {
		o.l.Errorf("error publising rollback message to %s channel", channel)
	}
}
