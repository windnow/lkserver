package lkserver

import (
	"context"
	"errors"
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

var errUnautorized = errors.New("NOT AUTORIZED")

func (s *lkserver) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		w.Header().Set("X-Request-ID", id.String())
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CTXKEYREQUESTID, id)))
	})
}

func (s *lkserver) authUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, s.config.SessionsKey)
		if err != nil {
			s.error(w, http.StatusInternalServerError, err)
			return
		}

		iin, ok := session.Values["user_iin"]
		if !ok {
			// s.error(w, http.StatusUnauthorized, errUnautorized) // Просто не будем устанавливать пользователя
			return
		}

		u, err := s.repo.GetUser(iin.(string))
		if err != nil {
			s.error(w, http.StatusUnauthorized, errUnautorized)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CTXUSER, u)))
	})
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
