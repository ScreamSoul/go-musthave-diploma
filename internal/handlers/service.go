package handlers

import "net/http"

func (s *UserLoyaltyServer) PingStorage(w http.ResponseWriter, r *http.Request) {
	if !s.loyaltyRepo.Ping(r.Context()) {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
