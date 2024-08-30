package lkserver

import "net/http"

func (s *lkserver) handleWhoAmI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := s.getSessionUser(w, r)
		if err != nil {
			s.respond(w, http.StatusUnauthorized, err)
			return
		}

		u.Sanitize()
		s.respond(w, http.StatusOK, u)
	}
}
