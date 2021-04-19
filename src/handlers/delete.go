package handlers

import (
	"net/http"
	"time"
)

// TODO(Jovan): Add delete methods
func (a *Auth) Logout (w http.ResponseWriter, r *http.Request) {
	a.l.Println("Handling Logout")
	cookie, err := r.Cookie("refresh")
	if err != nil {
		a.l.Printf("[ERROR] getting cookie: %v", err)
		http.Error(w, "Logged out", http.StatusOK)
		return
	}

	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}