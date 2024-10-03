package sqlite

import (
	"errors"
	"fmt"
	m "lkserver/internal/models"
)

var reports = "reports"

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
	err = InitReports(repo)
	if err != nil {
		return m.HandleError(err, "sqliteRepo.initReports")
	}

	s.reports = repo

	return nil

}

func InitReports(repo *reportsRepo) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %[1]s (
			ref 		BLOB	PRIMARY KEY,
			date 		INTEGER NOT NULL,
			number 		TEXT	NOT NULL,
			reg_number	TEXT	NOT NULL,
			author		BLOB	NOT NULL,

			FOREIGN KEY (author) REFERENCES user(ref)
		)
	`, reports)

	if err := repo.source.Exec(query); err != nil {
		return m.HandleError(err, "reportsRepo.InitReports")
	}
	return errors.ErrUnsupported
}
