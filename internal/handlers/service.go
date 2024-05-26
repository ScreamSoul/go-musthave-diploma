package handlers

import "net/http"

func (ms *UserLoyaltyServer) PingStorage(w http.ResponseWriter, r *http.Request) {
	if !ms.store.Ping(r.Context()) {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
