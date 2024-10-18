package lkserver

import (
	"errors"
	"lkserver/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func parseInt64(strVal string) int64 {
	value, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		return 0
	}
	return value
}

func getLimits(r *http.Request) (int64, int64) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	return parseInt64(limitStr), parseInt64(offsetStr)
}

func (s *lkserver) handleGetCato() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guid := mux.Vars(r)["guid"]
		if guid == "" {
			s.error(w, http.StatusBadRequest, errors.New("CATO GUID IS MISSING"))
			return
		}

		GUID, err := models.ParseJSONByteFromString(guid)
		if s.error(w, http.StatusBadRequest, err) {
			return
		}

		result, err := s.catalogsService.GetCato(r.Context(), GUID)
		if s.error(w, http.StatusInternalServerError, err) {
			return
		}

		s.respond(w, http.StatusOK, result)

	}
}
func (s *lkserver) handleCatoList() http.HandlerFunc {
	if s.catalogsService == nil {
		panic("catalogService is nil")
	}
	return func(w http.ResponseWriter, r *http.Request) {

		limit, offset := getLimits(r)
		parent := r.URL.Query().Get("parent")
		search := r.URL.Query().Get("search")
		guid, err := models.ParseJSONByteFromString(parent)
		if s.error(w, http.StatusBadRequest, err) {
			return
		}
		result, err := s.catalogsService.CatoList(r.Context(), guid, search, limit, offset)
		if s.error(w, http.StatusInternalServerError, err) {
			return
		}
		s.respond(w, http.StatusOK, result)

	}
}

func (s *lkserver) handleGetVus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guid := mux.Vars(r)["guid"]
		if guid == "" {
			s.error(w, http.StatusBadRequest, errors.New("VUS GUID IS MISSING"))
			return
		}

		GUID, err := models.ParseJSONByteFromString(guid)
		if s.error(w, http.StatusBadRequest, err) {
			return
		}

		result, err := s.catalogsService.GetVus(r.Context(), GUID)
		if s.error(w, http.StatusInternalServerError, err) {
			return
		}

		s.respond(w, http.StatusOK, result)

	}
}

func (s *lkserver) handleVusList() http.HandlerFunc {
	if s.catalogsService == nil {
		panic("catalogService is nil")
	}
	return func(w http.ResponseWriter, r *http.Request) {
		limit, offset := getLimits(r)
		search := r.URL.Query().Get("search")
		result, err := s.catalogsService.VusList(r.Context(), search, limit, offset)
		if s.error(w, http.StatusInternalServerError, err) {
			return
		}
		s.respond(w, http.StatusOK, result)
	}
}

func (s *lkserver) handleOrgsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		limit, offset := getLimits(r)
		search := r.URL.Query().Get("search")

		result, err := s.catalogsService.OrganizationList(r.Context(), search, limit, offset)
		if s.error(w, http.StatusInternalServerError, err) {
			return
		}
		s.respond(w, http.StatusOK, result)

	}
}

func (s *lkserver) handleDevisionsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		limit, offset := getLimits(r)
		search := r.URL.Query().Get("search")

		result, err := s.catalogsService.DevisionList(r.Context(), search, limit, offset)
		if s.error(w, http.StatusInternalServerError, err) {
			return
		}
		s.respond(w, http.StatusOK, result)

	}
}

func (s *lkserver) handleGetOrganization() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guid := mux.Vars(r)["guid"]
		if guid == "" {
			s.error(w, http.StatusBadRequest, errors.New("ORGANIZATION GUID IS MISSING"))
			return
		}
		GUID, err := models.ParseJSONByteFromString(guid)
		if s.error(w, http.StatusBadRequest, err) {
			return
		}

		result, err := s.catalogsService.GetOrganization(r.Context(), GUID)
		if s.error(w, http.StatusInternalServerError, err) {
			return
		}
		s.respond(w, http.StatusOK, result)
	}
}

func (s *lkserver) handleGetDevision() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guid := mux.Vars(r)["guid"]
		if guid == "" {
			s.error(w, http.StatusBadRequest, errors.New("DEVISION GUID IS MISSING"))
			return
		}
		GUID, err := models.ParseJSONByteFromString(guid)
		if s.error(w, http.StatusBadRequest, err) {
			return
		}

		result, err := s.catalogsService.GetDevision(r.Context(), GUID)
		if s.error(w, http.StatusInternalServerError, err) {
			return
		}
		s.respond(w, http.StatusOK, result)
	}
}
