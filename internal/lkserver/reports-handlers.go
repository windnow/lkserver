package lkserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"lkserver/internal/models"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *lkserver) handleGetReportType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		guid := vars["guid"]
		if guid == "" {
			s.error(w, http.StatusBadRequest, errors.New("MISSING TYPE GUID"))
			return
		}

		var GUID models.JSONByte
		guid = fmt.Sprintf("\"%s\"", guid)
		err := json.Unmarshal([]byte(guid), &GUID)
		if err != nil {
			s.error(w, http.StatusBadRequest, errors.New("WRONG GUID FORMAT"))
			return
		}
		result, err := s.reportsService.GetTypes([]string{})
		if s.error(w, http.StatusInternalServerError, err) {
			return
		}
		for _, reportType := range result {
			if reportType.Ref == GUID {
				s.respond(w, http.StatusOK, reportType)
				return
			}
		}

		s.error(w, http.StatusNotFound, errNotFound)

	}
}
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

		ctx, err := s.setUserContext(w, r)
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
		ctx, err := s.setUserContext(w, r)
		if s.error(w, http.StatusBadRequest, err) {
			return
		}

		list, err := s.reportsService.List(ctx)
		if s.error(w, http.StatusInternalServerError, err) {
			return
		}

		s.respond(w, http.StatusOK, list)

	}
}

func (s *lkserver) setUserContext(w http.ResponseWriter, r *http.Request) (context.Context, error) {
	user, err := s.getSessionUser(w, r)
	if err != nil {
		return nil, err
	}
	ctx := context.WithValue(r.Context(), models.CtxKey("user"), user)

	return ctx, nil
}
