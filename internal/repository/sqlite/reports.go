package sqlite

import (
	m "lkserver/internal/models"
)

type reportsRepo struct {
	source *src
	repo   *sqliteRepo
}

func (s *sqliteRepo) initReports() error {

	repo := &reportsRepo{
		source: s.db,
		repo:   s,
	}

	err := InitReportTypes(repo)
	if err != nil {
		return m.HandleError(err, "sqliteRepo.initReports")
	}

	s.reports = repo

	return nil

}
