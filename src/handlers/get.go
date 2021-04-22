package handlers

import (
	"net/http"
	"saltgram/data"
)

// TODO(Jovan): REMOVE!
func (e *Emails) GetAllResets(w http.ResponseWriter, r *http.Request) {
	data.ToJSON(data.GetAllResets(), w)
}

func (e *Emails) GetAll(w http.ResponseWriter, r *http.Request) {
	err := data.ToJSON(data.GetAllActivations(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

