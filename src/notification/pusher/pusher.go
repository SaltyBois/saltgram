package pusher

import (
	"github.com/pusher/pusher-http-go"
	"github.com/sirupsen/logrus"
)

type Pusehr struct {
	client pusher.Client
	l      *logrus.Logger
}

func NewPusher(l *logrus.Logger) *Pusehr {
	pusherClient := pusher.Client{
		AppID:   "1233087",
		Key:     "2c3e3d192fa386ba691c",
		Secret:  "3f47f19ea84fdb3a6688",
		Cluster: "eu",
		Secure:  true,
	}
	return &Pusehr{
		l:      l,
		client: pusherClient,
	}
}

func (p *Pusehr) PushNotification(message string) {
	data := map[string]string{"message": message}
	p.client.Trigger("my-channel", "my-event", data)
}
