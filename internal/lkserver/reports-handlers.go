package lkserver

import (
	"context"
	"encoding/json"
	"errors"
	"lkserver/internal/models"
	"net/http"

	"github.com/gorilla/mux"
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
		vars := mux.Vars(r)
		reportType := vars["type"]
		if reportType == "" {
			s.error(w, http.StatusBadRequest, errors.New("MISSING REPORT TYPE"))
			return
		}

		ctx, err := s.getUserContext(w, r)
		if s.error(w, http.StatusBadRequest, err) {
			return
		}

		incomingData, err := s.reportsService.GetStructure(reportType)
		if s.error(w, http.StatusBadRequest, err) {
			return
		}
		if s.error(w, http.StatusInternalServerError, json.NewDecoder(r.Body).Decode(incomingData)) {
			return
		}

		if s.error(w, http.StatusBadRequest, s.reportsService.Save(ctx, reportType, incomingData)) {
			return
		}

		s.respond(w, http.StatusAccepted, incomingData)

	}
}

func (s *lkserver) handleReportsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *lkserver) getUserContext(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	user, err := s.getSessionUser(w, r)
	if err != nil {
		return nil, err
	}
	ctx := context.WithValue(r.Context(), models.CtxKey("user"), user)

	return ctx, nil
}
