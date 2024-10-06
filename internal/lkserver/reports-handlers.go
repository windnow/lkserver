package lkserver

import (
	"context"
	"encoding/json"
	"errors"
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

func (s *lkserver) handleSaveReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reportType := r.Header.Get("X-Report-Type")
		user, err := s.getSessionUser(w, r)
		if s.error(w, http.StatusBadRequest, err) {
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		if reportType == "" {
			s.error(w, http.StatusBadRequest, errors.New("MISSING REPORT TYPE (X-Report-Type)"))
			return
		}

		incomingData, err := s.reportsService.GetStructure(reportType)
		if s.error(w, http.StatusBadRequest, err) {
			return
		}
		if s.error(w, http.StatusInternalServerError, json.NewDecoder(r.Body).Decode(incomingData)) {
			return
		}

		if s.error(w, http.StatusBadRequest, s.reportsService.Save(ctx, incomingData)) {
			return
		}

		s.respond(w, http.StatusAccepted, incomingData)

	}
}
