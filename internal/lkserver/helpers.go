package lkserver

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func getParam(param string, r *http.Request) (string, error) {
	vars := mux.Vars(r)
	value := vars[param]
	if value == "" {
		return "", fmt.Errorf("VALUE \"%s\" IS MISSING", param)
	}

	return value, nil
}
