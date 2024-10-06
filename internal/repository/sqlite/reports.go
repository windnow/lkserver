package sqlite

import (
	"context"
	"fmt"
	m "lkserver/internal/models"
	"lkserver/internal/repository"
)

type reportFactory struct {
	detailsMap map[string]repository.ReportDetails
}

func NewReportFactory(source *src) *reportFactory {
	return &reportFactory{
		detailsMap: map[string]repository.ReportDetails{
			"0001": &DepartureOnBusinessTrip{source: source},
		},
	}
}

func (f *reportFactory) GetAllProcessors() []repository.ReportDetails {
	var result []repository.ReportDetails
	for _, val := range f.detailsMap {
		result = append(result, val)
	}

	return result
}

func (f *reportFactory) GetReportProcessor(reportType string) (repository.ReportDetails, error) {

	details, ok := f.detailsMap[reportType]
	if !ok {
		return nil, m.HandleError(fmt.Errorf("unsupported report type: %s", reportType))
	}

	return details, nil

}

type reportsRepo struct {
	source  *src
	repo    *sqliteRepo
	factory *reportFactory
}

func (s *sqliteRepo) initReports() error {

	repo := &reportsRepo{
		source:  s.db,
		repo:    s,
		factory: NewReportFactory(s.db),
	}

	if err := InitReportTypes(repo); err != nil {
		return m.HandleError(err, "sqliteRepo.initReports")
	}
	if err := InitReports(repo); err != nil {
		return m.HandleError(err, "sqliteRepo.initReports")
	}
	if err := InitCoordinators(repo); err != nil {
		return m.HandleError(err, "sqliteRepo.initReports")
	}
	allTypes := repo.factory.GetAllProcessors()
	for _, rt := range allTypes {
		if err := rt.Init(); err != nil {
			return m.HandleError(err, "sqliteRepo.initReports")
		}
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
	`, tabReport)

	if err := repo.source.Exec(query); err != nil {
		return m.HandleError(err, "reportsRepo.InitReports")
	}
	return nil
}

func (repo *reportsRepo) Save(ctx context.Context, report *m.Report) error {

	if report.Ref.Blank() {
		guid, err := m.GenerateUUID()
		if err != nil {
			return m.HandleError(err, "reportsRepo.Save")
		}
		report.Ref = guid
	}

	return m.HandleError(repo.source.ExecContextInTransaction(ctx,
		saveRepoQuery,
		report.Ref,
		report.Date,
		report.Number,
		report.RegNumber,
		report.Author,
	), "reportRepo.Save")

}

var saveRepoQuery = `
    INSERT OR REPLACE INTO %[1]s (
        ref, date, number, reg_number, author	 
	) VALUES (?, ?, ?, ?, ?)
`
