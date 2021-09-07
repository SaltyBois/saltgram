package saga

import "encoding/json"

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

type Message struct {
	Service         string
	SenderService   string
	Action          string
	UserId          uint64
	Username        string
	ProfileFolderId string
	PostsFolderId   string
	StoriesFolderId string
	Description     string
	PhoneNumber     string
	Gender          string
	DateOfBirth     int64
	WebSite         string
	PrivateProfile  bool
	Email           string
}

func (m Message) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}
