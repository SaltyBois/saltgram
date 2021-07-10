package servers

import (
	"context"
	"saltgram/notification/data"
	"saltgram/notification/pusher"
	"saltgram/protos/notifications/prnotifications"
	"saltgram/protos/users/prusers"

	"github.com/sirupsen/logrus"
)

type Notification struct {
	prnotifications.UnimplementedNotificationsServer
	l  *logrus.Logger
	db *data.DBConn
	uc prusers.UsersClient
	p  *pusher.Pusehr
}

func NewNotification(l *logrus.Logger, db *data.DBConn, uc prusers.UsersClient, p *pusher.Pusehr) *Notification {
	return &Notification{
		l:  l,
		db: db,
		uc: uc,
		p:  p,
	}
}

func (n *Notification) CreateLikeNotification(ctx context.Context, r *prnotifications.Request) (*prnotifications.Respond, error) {
	notification := data.Notification{
		UserID:         r.UserId,
		ReferredUserId: r.ReferredId,
		Type:           data.LIKE,
		Seen:           false,
	}
	err := n.db.CreateNotification(&notification)
	if err != nil {
		return &prnotifications.Respond{}, err
	}

	user, err := n.uc.GetByUserId(context.Background(), &prusers.GetByIdRequest{Id: r.ReferredId})
	if err != nil {
		return &prnotifications.Respond{}, err
	}

	message := "@" + user.Username + " liked your post"
	n.p.PushNotification(message)

	return &prnotifications.Respond{}, nil
}

func (n *Notification) CreateCommentNotification(ctx context.Context, r *prnotifications.Request) (*prnotifications.Respond, error) {
	notification := data.Notification{
		UserID:         r.UserId,
		ReferredUserId: r.ReferredId,
		Type:           data.COMMENT,
		Seen:           false,
	}
	err := n.db.CreateNotification(&notification)
	if err != nil {
		return &prnotifications.Respond{}, err
	}

	user, err := n.uc.GetByUserId(context.Background(), &prusers.GetByIdRequest{Id: r.ReferredId})
	if err != nil {
		return &prnotifications.Respond{}, err
	}

	message := "@" + user.Username + " has commented your post"
	n.p.PushNotification(message)

	return &prnotifications.Respond{}, nil
}

func (n *Notification) CreateFollowNotification(ctx context.Context, r *prnotifications.RequestUsername) (*prnotifications.Respond, error) {
	notification := data.Notification{
		UserID:         r.UserId,
		ReferredUserId: r.ReferredId,
		Type:           data.FOLLOW,
		Seen:           false,
	}
	err := n.db.CreateNotification(&notification)
	if err != nil {
		return &prnotifications.Respond{}, err
	}

	message := "@" + r.ReferredUsername + " started following you"
	n.p.PushNotification(message)

	return &prnotifications.Respond{}, nil
}

func (n *Notification) CreateFollowRequestNotification(ctx context.Context, r *prnotifications.RequestUsername) (*prnotifications.Respond, error) {
	notification := data.Notification{
		UserID:         r.UserId,
		ReferredUserId: r.ReferredId,
		Type:           data.FOLLOWREQUEST,
		Seen:           false,
	}
	err := n.db.CreateNotification(&notification)
	if err != nil {
		return &prnotifications.Respond{}, err
	}

	message := "@" + r.ReferredUsername + " send you a following request you"
	n.p.PushNotification(message)

	return &prnotifications.Respond{}, nil
}

func (n *Notification) GetUnseenNotificationsCount(ctx context.Context, r *prnotifications.NProfile) (*prnotifications.NotificationCount, error) {
	user, err := n.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: r.Username})
	if err != nil {
		n.l.Errorf("failure getting user: %v\n", err)
		return &prnotifications.NotificationCount{}, err
	}
	count, err := n.db.GetUnseenNotificationsCount(user.Id)
	if err != nil {
		n.l.Errorf("failure counting notifications: %v\n", err)
		return &prnotifications.NotificationCount{}, err
	}

	return &prnotifications.NotificationCount{Count: count}, nil
}

func (n *Notification) NotificationSeen(ctx context.Context, r *prnotifications.NProfile) (*prnotifications.Respond, error) {
	user, err := n.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: r.Username})
	if err != nil {
		n.l.Errorf("failure getting user: %v\n", err)
		return &prnotifications.Respond{}, err
	}
	err = n.db.NotificationSeen(user.Id)
	if err != nil {
		n.l.Errorf("failure updating notifications: %v\n", err)
		return &prnotifications.Respond{}, err
	}

	return &prnotifications.Respond{}, nil
}

func (n *Notification) GetNotifications(r *prnotifications.NProfile, stream prnotifications.Notifications_GetNotificationsServer) error {
	user, err := n.uc.GetByUsername(context.Background(), &prusers.GetByUsernameRequest{Username: r.Username})
	if err != nil {
		n.l.Errorf("failure getting user: %v\n", err)
		return err
	}
	notifications, err := n.db.GetNotification(user.Id)
	if err != nil {
		n.l.Errorf("failure fetching notifications: %v\n", err)
		return err
	}
	for _, notification := range notifications {
		referred, err := n.uc.GetProfileByUserId(context.Background(), &prusers.GetByIdRequest{Id: notification.ReferredUserId})
		if err != nil {
			n.l.Errorf("failure getting user: %v\n", err)
			continue
		}
		n := &prnotifications.Notification{
			Username:                      user.Username,
			ReferredUsername:              referred.Username,
			ReferredUserProfilePictureURL: referred.ProfilePictureURL,
			Type:                          string(notification.Type),
			Seen:                          notification.Seen,
		}
		err = stream.Send(n)
		if err != nil {
			return err
		}
	}
	return nil
}
