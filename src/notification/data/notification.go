package data

type NotificationType string

const (
	MESSAGE NotificationType = "MESSAGE"
	FOLLOWING
	POST
	STORY
	COMMENT
)

type Notification struct {
	ID           uint64               `json:"id"`
	NotificationType NotificationType `json:"notificationType" validate:"required"`
	User         User                 `json:"user"`
	UserID       uint64               `json:"userId"`
	Content      string               `json:"content" validate:"required"`
}

func (db *DBConn) AddNotification(notification *Notification) error {
	return db.DB.Create(notification).Error
}