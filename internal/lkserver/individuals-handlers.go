package lkserver

import (
	"net/http"
)

func (s *lkserver) handleIndividualsByIIN() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		iin, err := getParam("iin", r)
		if err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}
		individual, err := s.repo.Individuals.Get(iin)
		if err != nil {
			s.error(w, http.StatusNotFound, err)
			return
		}
		s.respond(w, http.StatusOK, individual)
	}
}
