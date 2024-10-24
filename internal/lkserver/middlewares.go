package lkserver

import (
	"context"
	"lkserver/internal/models"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ContextKey string

const (
	CTXKEYREQUESTID ContextKey = "requestId"
	CTXUSER         ContextKey = "user"
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

		u, _ := s.getSessionUser(w, r)

		if u == nil {
			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CTXUSER, u)))
		}
	})
}

func (s *lkserver) checkUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := r.Context().Value(CTXUSER)
		if user == nil {
			s.error(w, http.StatusNotFound, errUnautorized)
			return

		}
		next.ServeHTTP(w, r)

	})
}

func (s *lkserver) getSessionUser(w http.ResponseWriter, r *http.Request) (*models.User, error) {

	iin, err := s.sessionGetValue(w, r, "user_iin")
	if err != nil {
		return nil, errUnautorized
	}

	u, err := s.repo.User.GetUser(iin.(string))
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
		name := "UNAUTORIZED"
		u, ok := r.Context().Value(CTXUSER).(*models.User)
		if ok && u != nil {
			name = u.Name
		}

		logger.Infof("%d %d %v %s %s %s", rw.code, rw.size, time.Since(start), name, r.Method, r.RequestURI)
	})
}
