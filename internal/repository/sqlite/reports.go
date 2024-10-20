package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
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

func (r *reportsRepo) GetTransaction(ctx context.Context) (*sql.Tx, error) {
	return r.source.db.BeginTx(ctx, nil)
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

func (repo *reportsRepo) Save(tx *sql.Tx, ctx context.Context, report *m.Report) error {

	if report.Ref.Blank() {
		guid, err := m.GenerateUUID()
		if err != nil {
			return m.HandleError(err, "reportsRepo.Save")
		}
		report.Ref = guid
	}

	return m.HandleError(repo.source.ExecContextInTransaction(ctx, saveReportQuery, tx,
		report.Ref,
		report.Type,
		report.Date.Unix(),
		report.Number,
		report.RegNumber,
		report.Author,
	), "reportRepo.Save")
}

func (repo *reportsRepo) GetStructure(reportType string) (interface{}, error) {

	processor, err := repo.factory.GetReportProcessor(reportType)
	if err != nil {
		return nil, err
	}

	return processor.GetStructure(), nil
}

func (repo *reportsRepo) List(ctx context.Context, userKey m.JSONByte) ([]*m.Report, error) {

	query := fmt.Sprintf("SELECT ref, type, date, number, reg_number, author FROM %[1]s WHERE author = ?", types.Report)

	rows, err := repo.source.db.QueryContext(ctx, query, userKey)
	if err != nil {
		return nil, err
	}
	var result []*m.Report
	for rows.Next() {
		report := &m.Report{}
		if err = rows.Scan(&report.Ref, &report.Type, &report.Date, &report.Number, &report.RegNumber, &report.Author); err != nil {
			return nil, err
		}
		result = append(result, report)
	}

	return result, nil
}

func (repo *reportsRepo) Get(guid m.JSONByte) (*m.Report, error) {
	query := fmt.Sprintf(`
		SELECT type, date, number, reg_number, author
		FROM %[1]s
		WHERE ref = ?
	`, types.Report)
	report := &m.Report{Ref: guid}

	if err := repo.source.db.QueryRow(query, guid).Scan(
		&report.Type,
		&report.Date,
		&report.Number,
		&report.RegNumber,
		&report.Author,
	); err != nil {
		return nil, err
	}

	return report, nil
}

func InitReports(repo *reportsRepo) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %[1]s (
			ref 		BLOB	PRIMARY KEY,
			type 		BLOB	NOT NULL,
			date 		INTEGER NOT NULL,
			number 		TEXT	NOT NULL,
			reg_number	TEXT,
			author		BLOB	NOT NULL,

			FOREIGN KEY (type) REFERENCES %[2]s(ref)
			FOREIGN KEY (author) REFERENCES %[3]s(ref)
		)
	`, types.Report, types.ReportType, types.Users)

	if err := repo.source.Exec(query); err != nil {
		return m.HandleError(err, "reportsRepo.InitReports")
	}
	return nil
}

var saveReportQuery = fmt.Sprintf(`
    INSERT OR REPLACE INTO %[1]s (
        ref, type, date, number, reg_number, author	 
	) VALUES (?, ?, ?, ?, ?, ?)
`, types.Report)
