package lkserver

import (
	"errors"
	"io"
	"net/http"
)

func (s *lkserver) handleFile() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fileID, err := getParam("id", r)

		if s.error(w, http.StatusBadRequest, err) {
			return
		}

		fileReader, mimeType, err := s.fileStore.GetFile(fileID)
		if s.error(w, http.StatusNotFound, err, errors.New("FILE NOT FOUND")) { // err may contain a path to a file on the server
			return
		}
		w.Header().Set("Content-Type", mimeType)

		_, err = io.Copy(w, fileReader)
		if s.error(w, http.StatusInternalServerError, err) {
			return
		}
	}
}
