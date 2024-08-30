package lkserver

import (
	"encoding/json"
	"net/http"
)

func (s *lkserver) handleSessionCreate() http.HandlerFunc {
	type request struct {
		Iin string `json:"iin"`
		Pin string `json:"pin"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}

		u, err := s.repo.User.FindUser(req.Iin, req.Pin)
		if err != nil {
			s.error(w, http.StatusUnauthorized, err)
			return
		}

		if err := s.sessionAddValue(w, r, "user_iin", u.Iin); err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		u.Sanitize()
		s.respond(w, http.StatusOK, u)

	}
}

func (s *lkserver) handleSessionDestroy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.sessionDeleteValue(w, r, "user_iin")
		s.respond(w, http.StatusOK, struct{ status string }{status: "Ok"})
	}
}
