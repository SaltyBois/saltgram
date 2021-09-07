package handlers

import "net/http"

func (a *Agent) IsAgent(w http.ResponseWriter, r *http.Request) {
	jws, err := getUserJWS(r)
	if err != nil {
		a.l.Errorf("failed getting jws: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	u, err := a.db.GetByJWS(jws)
	if err != nil {
		a.l.Errorf("failed getting user: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	if !u.Agent {
		http.Error(w, "Not agent", http.StatusBadRequest)
		return
	}
	w.Write([]byte("OK"))
}