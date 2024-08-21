package lkserver

import (
	"context"
	"lkserver/internal/models"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type key int

const (
	CTXKEYREQUESTID key = iota
	CTXUSER
)

func (s *lkserver) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		w.Header().Set("X-Request-ID", id.String())
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CTXKEYREQUESTID, id)))
	})
}

func (s *lkserver) authUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		u, _ := s.getSessionUser(r)

		if u == nil {
			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CTXUSER, u)))
		}
	})
}

func (s *lkserver) checkUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, err := s.getSessionUser(r)
		if err != nil {
			s.error(w, http.StatusNotFound, err)
			return
		}
		next.ServeHTTP(w, r)

	})
}

func (s *lkserver) getSessionUser(r *http.Request) (*models.User, error) {

	session, err := s.sessionStore.Get(r, s.config.SessionsKey)
	if err != nil {
		return nil, err
	}

	iin, ok := session.Values["user_iin"]
	if !ok {
		return nil, errUnautorized
	}

	u, err := s.repo.GetUser(iin.(string))
	if err != nil {
		return nil, errNotFound
	}
	return u, nil

}
func (s *lkserver) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(CTXKEYREQUESTID),
		})

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK, 0}

		next.ServeHTTP(rw, r)
		name := "-"
		uany := r.Context().Value(CTXUSER)
		if uany != nil {
			u := uany.(*models.User)
			if u != nil {
				name = u.Name
			}
		}

		logger.Infof("%s %s %s %d %v %d", name, r.Method, r.RequestURI, rw.code, time.Since(start), rw.size)
	})
}
