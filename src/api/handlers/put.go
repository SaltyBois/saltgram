package handlers

import (
	"context"
	"net/http"
	"saltgram/protos/email/premail"

	"github.com/gorilla/mux"
)

func (e *Email) Activate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	_, err := e.ec.Activate(context.Background(), &premail.ActivateRequest{Token: token})
	if err != nil {
		e.l.Printf("[ERROR] activating email: %v\n", err)
		http.Error(w, "Failed to activate email", http.StatusBadRequest)
		return
	}

	w.Write([]byte("200 - OK"))
}
