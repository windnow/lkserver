package lkserver

import (
	"errors"
	"lkserver/internal/models"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *lkserver) handleGetUserInfo() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		guid := vars["guid"]
		if guid == "" {
			s.error(w, http.StatusBadRequest, errors.New("MISSING TYPE GUID"))
			return
		}

		GUID, err := models.ParseJSONByteFromString(guid)
		if err != nil {
			s.error(w, http.StatusBadRequest, errors.New("WRONG GUID FORMAT"))
			return
		}
		result, err := s.usersService.GetUserInfo(GUID)
		if s.error(w, http.StatusInternalServerError, err) {
			return
		}

		s.respond(w, http.StatusOK, result)
	}

}