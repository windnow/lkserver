package lkserver

import (
	"errors"
	"io"
	"net/http"
)

func (s *lkserver) handleFile() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fileID, err := getParam("id", r)

		if err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}

		fileReader, mimeType, err := s.fileStore.GetFile(fileID)
		if err != nil {
			s.error(w, http.StatusNotFound, errors.New("FILE NOT FOUND")) // err may contain a path to a file on the server
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
