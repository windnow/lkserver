package lkserver

import (
	"lkserver/internal/models"
	"net/http"
)

func (s *lkserver) handleGetReportTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		queryParams := r.URL.Query()["types"]

		result, err := s.reportsService.GetTypes(queryParams)
		if s.error(w, http.StatusInternalServerError, err) {
			return
		}

		if len(result) == 0 {
			s.error(w, http.StatusNotFound, models.ErrNotFound)
			return
		}

		s.respond(w, http.StatusOK, result)
	}
}
