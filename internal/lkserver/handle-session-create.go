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

		u, err := s.repo.FindUser(req.Iin, req.Pin)
		if err != nil {
			s.error(w, http.StatusUnauthorized, err)
			return
		}

		if err := s.updateSession(w, r, "user_iin", u.Iin); err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		u.Sanitize()
		s.respond(w, http.StatusOK, u)

	}
}

func (s *lkserver) updateSession(w http.ResponseWriter, r *http.Request, key, value string) error {
	session, err := s.sessionStore.Get(r, s.config.SessionsKey)
	if err != nil {
		return err
	}

	session.Values[key] = value
	if err := s.sessionStore.Save(r, w, session); err != nil {
		return err
	}

	return nil
}

func (s *lkserver) error(w http.ResponseWriter, code int, err error) {

	s.respond(w, code, map[string]string{"error": err.Error()})

}

func (s *lkserver) respond(w http.ResponseWriter, code int, data interface{}) {

	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}

}
