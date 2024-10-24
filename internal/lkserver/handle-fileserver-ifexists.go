package lkserver

import (
	"net/http"
	"os"
)

func (s *lkserver) handleFileServerIfExists() {
	_, err := os.Stat(s.config.StaticFilesPath)
	if err != nil {
		return
	}
	var staticFilesPrefix = "/"
	fileServer := http.FileServer(http.Dir(s.config.StaticFilesPath))
	// sp := http.StripPrefix(staticFilesPrefix, fileServer)
	s.PathPrefix(staticFilesPrefix).Handler(fileServer) //sp)
}
