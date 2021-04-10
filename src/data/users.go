package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	CreatedOn string `json:"-"`
	UpdatedOn string `json:"-"`
	DeletedOn string `json:"-"`
}

// NOTE(Jovan): Collection of users
type Users []*User

// NOTE(Jovan): Serializing to JSON
// NewEncoder provides better perf than json.Unmarshal
// https://golang.org/pkg/encoding/json/#NewEncoder
func (u *Users) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

// NOTE(Jovan): Deserializing from JSON
func (u *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(u)
}

// TODO(Jovan): For testing currently, use a DB of sorts
func GetUsers() Users {
	return userList
}

func AddUser(u *User) {
	u.ID = getNextID()
	userList = append(userList, u)
} 

func UpdateUser(id uint64, u *User) error {
	_, pos, err := findUser(id)
	if err != nil {
		return err
	}

	u.ID = id
	userList[pos] = u

	return nil
}

var ErrUserNotFound = fmt.Errorf("User not found")

func findUser(id uint64) (*User, int, error) {
	for i, u := range userList {
		if u.ID == id {
			return u, i, nil
		}
	}

	return nil, -1, ErrUserNotFound
}

func getNextID() uint64 {
	u := userList[len(userList) - 1]
	return u.ID + 1
}

var userList = []*User{
	&User{
		ID: 1,
		Username: "AgentSmith",
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	},
	&User{
		ID: 2,
		Username: "Neo",
		CreatedOn: time.Now().UTC().String(),
		UpdatedOn: time.Now().UTC().String(),
	},
}