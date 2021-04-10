package data

import (
	"encoding/json"
	"io"
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
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

// TODO(Jovan): For testing currently, use a DB of sorts
func GetUsers() Users {
	return []*User {
		&User{
			ID: 1,
			Username: "AgentSmith",
		},
		&User{
			ID: 2,
			Username: "Neo",
		},
	}
}