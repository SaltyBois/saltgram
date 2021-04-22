package handlers

import (
	"net/http"
	"os"
	"saltgram/users/data"

	"github.com/dgrijalva/jwt-go"
)

func (u *Users) GetByJWS(db *data.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		jws, err := getUserJWS(r)
		if err != nil {
			u.l.Println("[ERROR] JWS not found")
			http.Error(w, "JWS not found", http.StatusBadRequest)
			return
		}

		token, err := jwt.ParseWithClaims(
			jws,
			&AccessClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			},
		)

		if err != nil {
			u.l.Printf("[ERROR] parsing claims: %v", err)
			http.Error(w, "Error parsing claims", http.StatusBadRequest)
			return
		}

		claims, ok := token.Claims.(*AccessClaims)

		if !ok {
			u.l.Println("[ERROR] unable to parse claims")
			http.Error(w, "Error parsing claims: ", http.StatusInternalServerError)
			return
		}

		user, err := db.GetUserByUsername(claims.Username)
		if err != nil {
			u.l.Println("[ERROR] fetching user", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if user.HashedPassword != claims.Password {
			u.l.Println("[ERROR] passwords do not match")
			http.Error(w, "JWT password doesn't match user's password", http.StatusUnauthorized)
			return
		}

		err = data.ToJSON(user, w)
		if err != nil {
			u.l.Println("[ERROR] serializing user ", err)
			http.Error(w, "Error serializing user", http.StatusInternalServerError)
			return
		}
	}
}

func (u *Users) GetAll(db *data.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u.l.Println("Handling GET Users")
		ul := db.GetUsers()
		err := data.ToJSON(ul, w)
		if err != nil {
			http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
		}
	}
}
