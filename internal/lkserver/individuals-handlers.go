package lkserver

import (
	"lkserver/internal/models"
	"net/http"
)

type individuals struct {
	*models.Individuals
	LastRank    *models.RankHistory   `json:"last_rank"`
	RankHistory []*models.RankHistory `json:"rank_history"`
}

func (s *lkserver) handleIndividualsByIIN() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		iin, err := getParam("iin", r)
		if err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}
		individ, err := s.repo.Individuals.GetByIin(iin)
		if err != nil {
			s.error(w, http.StatusNotFound, err)
			return
		}
		LastRank, _ := s.repo.RanksHistory.GetLastByIin(iin)
		RankHistory, _ := s.repo.RanksHistory.GetHistoryByIin(iin)

		full := &individuals{
			Individuals: individ,
			LastRank:    LastRank,
			RankHistory: RankHistory,
		}

		s.respond(w, http.StatusOK, full)
	}
}

func (s *lkserver) handleEducationByIIN() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		iin, err := getParam("iin", r)
		if err != nil {
			s.error(w, http.StatusBadRequest, err)
			return
		}
		edu, err := s.repo.Education.Get(iin)
		if err != nil {
			s.error(w, http.StatusNotFound, err)
			return
		}
		if len(edu) == 0 {
			s.error(w, http.StatusNotFound, err)
			return
		}

		s.respond(w, http.StatusOK, edu)
	}
}
