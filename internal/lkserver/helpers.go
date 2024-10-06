package lkserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const timeStampSuffix = "InitTimeStamp"

func getParam(param string, r *http.Request) (string, error) {
	vars := mux.Vars(r)
	value := vars[param]
	if value == "" {
		return "", fmt.Errorf("VALUE \"%s\" IS MISSING", param)
	}

	return value, nil
}

func (s *lkserver) sessionAddValue(w http.ResponseWriter, r *http.Request, key, value string) error {
	session, err := s.sessionStore.Get(r, s.config.SessionsKey)
	if err != nil {
		return err
	}

	session.Values[key] = value
	session.Values[key+timeStampSuffix] = time.Now().Unix()
	if err := s.sessionStore.Save(r, w, session); err != nil {
		return err
	}

	return nil
}

func (s *lkserver) sessionGetValue(w http.ResponseWriter, r *http.Request, key string) (interface{}, error) {
	session, err := s.sessionStore.Get(r, s.config.SessionsKey)
	if err != nil {
		return nil, err
	}

	value, ok := session.Values[key]
	if !ok {
		s.sessionDeleteValue(w, r, key)
		return nil, errNotFound
	}
	initTimestamp := key + timeStampSuffix
	sessionAge, ok := session.Values[initTimestamp].(int64)
	if !ok {
		return nil, errNotFound
	}
	if time.Now().Unix()-sessionAge > int64(s.config.SessionMaxAge) {
		return nil, errNotFound
	}

	return value, nil

}

func (s *lkserver) sessionDeleteValue(w http.ResponseWriter, r *http.Request, key string) {
	session, err := s.sessionStore.Get(r, s.config.SessionsKey)
	if err != nil {
		return
	}
	_, ok := session.Values[key]
	if !ok {
		return
	}
	delete(session.Values, key)

	initTimestamp := key + timeStampSuffix
	_, ok = session.Values[initTimestamp]
	if ok {
		delete(session.Values, initTimestamp)
	}
	s.sessionStore.Save(r, w, session)
}

func (s *lkserver) error(w http.ResponseWriter, code int, err error, placeholder ...error) bool {

	if err == nil {
		return false
	}
	if len(placeholder) > 0 {
		err = placeholder[0]
	}
	s.respond(w, code, map[string]string{"error": err.Error()})
	return true

}

func (s *lkserver) respond(w http.ResponseWriter, code int, data interface{}) {

	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}

}
