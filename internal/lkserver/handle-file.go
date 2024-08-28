package lkserver

import (
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *lkserver) handleFile() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fileID := vars["id"]

		if fileID == "" {
			s.error(w, http.StatusBadRequest, errors.New("FILE ID IS MISSING"))
			return
		}

		fileReader, mimeType, err := s.repo.Files.GetFile(fileID)
		if err != nil {
			s.error(w, http.StatusNotFound, errors.New("FILE NOT FOUND"))
			return
		}
		w.Header().Set("Content-Type", mimeType)

		_, err = io.Copy(w, fileReader)
		if err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}
	}
}
