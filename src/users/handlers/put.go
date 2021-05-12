package handlers

import (
	"net/http"
	"saltgram/users/data"

	"github.com/go-playground/validator"
)

type PasswordValidDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (dto *PasswordValidDTO) Validate() error {
	valid := validator.New()
	return valid.Struct(dto)
}

func (u *Users) IsPasswordValid(db *data.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u.l.Println("Checking if password valid")

		dto := PasswordValidDTO{}
		err := data.FromJSON(&dto, r.Body)
		if err != nil {
			u.l.Printf("[ERROR] deserializing PasswordValidDTO: %v\n", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err = dto.Validate()
		if err != nil {
			u.l.Printf("[ERROR] validating PasswordValidDTO: %v\n", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		hashedPass, err := data.IsPasswordValid(db, dto.Username, dto.Password)
		if err != nil {
			u.l.Printf("[ERROR] invalid password")
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		w.Write([]byte(hashedPass))
	}
}

// func (u *Users) Update(db *data.DBConn) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		id, err := getUserID(r)

// 		if err != nil {
// 			http.Error(w, "Invalid id", http.StatusBadRequest)
// 			return
// 		}

// 		u.l.Println("Handling PUT Users", id)
// 		// NOTE(Jovan): Safe to cast because middleware makes sure nothing's wrong
// 		user := r.Context().Value(KeyUser{}).(data.User)

// 		err = db.UpdateUser(&user)
// 		if err != nil {
// 			u.l.Printf("[ERROR] updating user: %v\n", err)
// 			http.Error(w, "Failed to update user", http.StatusNotFound)
// 			return
// 		}
// 	}
// }
