package data

import "saltgram/data"

type NotificationType string

const (
	LIKE          NotificationType = "LIKE"
	COMMENT       NotificationType = "COMMENT"
	FOLLOW        NotificationType = "FOLLOW"
	POST          NotificationType = "POST"
	STORY         NotificationType = "STORY"
	FOLLOWREQUEST NotificationType = "FOLLOWREQUEST"
)

type Notification struct {
	data.Identifiable
	UserID         uint64           `json:"userId" gorm:"type:numeric"`
	ReferredUserId uint64           `json:"referdUserId" gorm:"type:numeric"`
	Type           NotificationType `json:"type"`
	Seen           bool             `json:"seen"`
}

func (db *DBConn) CreateNotification(notification *Notification) error {
	return db.DB.Create(notification).Error
}

func (db *DBConn) UpdateNotification(notification *Notification) error {
	n := Notification{}
	err := db.DB.First(&n).Error
	if err != nil {
		return err
	}
	return db.DB.Save(notification).Error
}

func (db *DBConn) GetNotification(id uint64) ([]*Notification, error) {
	res := []*Notification{}
	err := db.DB.Where("user_id = ?", id).Find(&res).Error
	return res, err
}

func (db *DBConn) GetUnseenNotification(id uint64) ([]*Notification, error) {
	res := []*Notification{}
	err := db.DB.Where("user_id = ? AND seen = ?", id, false).Find(&res).Error
	return res, err
}

func (db *DBConn) NotificationSeen(id uint64) error {
	return db.DB.Model(&Notification{}).Where("user_id = ?", id).Update("seen", true).Error
}

func (db *DBConn) GetUnseenNotificationsCount(id uint64) (int64, error) {
	var count int64
	err := db.DB.Model(&Notification{}).Where("user_id = ? AND seen = ?", id, false).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
