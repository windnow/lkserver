package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
	"lkserver/internal/repository"
)

type reportsRepo struct {
	source  *src
	repo    *sqliteRepo
	factory *reportFactory
}

type reportFactory struct {
	detailsMap map[string]repository.ReportDetails
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

func (r *reportsRepo) GetTransaction(ctx context.Context) (*sql.Tx, error) {
	return r.source.db.BeginTx(ctx, nil)
}

func (repo *reportsRepo) META(reportType string) map[string]m.META {
	processor, err := repo.factory.GetReportProcessor(reportType)
	if err != nil {
		return nil
	}

	return processor.META()
}

func (repo *reportsRepo) GetStructure(reportType string) (interface{}, error) {

	processor, err := repo.factory.GetReportProcessor(reportType)
	if err != nil {
		return nil, err
	}

	return processor.GetStructure(), nil
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

func (repo *reportsRepo) TypesCount() uint64 {

	var count uint64
	if err := repo.source.db.QueryRow(fmt.Sprintf(`select count(ref) from %[1]s`, types.ReportType)).Scan(&count); err != nil {
		return 0
	}
	return count
}
func (repo *reportsRepo) Count(userRef, reportType m.JSONByte) uint64 {

	var row *sql.Row

	if reportType.Blank() {
		query := fmt.Sprintf(`select count(ref) from %[1]s WHERE author = ?`, types.Report)
		row = repo.source.db.QueryRow(query, userRef)
	} else {
		query := fmt.Sprintf(`select count(ref) from %[1]s WHERE author = ? AND type = ?`, types.Report)
		row = repo.source.db.QueryRow(query, userRef, reportType)
	}
	var count uint64
	if err := row.Scan(&count); err != nil {
		return 0
	}
	return count
}

func (repo *reportsRepo) List(ctx context.Context, userKey, typeRef m.JSONByte, limits ...int64) ([]*m.Report, error) {

	var result []*m.Report
	var err error

	limit, offset := limitations(limits)

	if typeRef.Blank() {
		query := fmt.Sprintf("SELECT ref, type, date, number, reg_number, author FROM %[1]s WHERE author = ? LIMIT %[2]d OFFSET %[3]d", types.Report, limit, offset)
		result, err = m.Query[*m.Report](repo.source.db, ctx, m.NewReport, nil, query, userKey)
	} else {
		query := fmt.Sprintf("SELECT ref, type, date, number, reg_number, author FROM %[1]s WHERE author = ? AND type = ? LIMIT %[2]d OFFSET %[3]d", types.Report, limit, offset)
		result, err = m.Query[*m.Report](repo.source.db, ctx, m.NewReport, nil, query, userKey, typeRef)
	}
	if err != nil {
		return nil, m.HandleError(err, "reportsRepo.List")
	}

	return result, nil
}

func (repo *reportsRepo) Get(guid m.JSONByte) (*m.Report, error) {
	query := fmt.Sprintf(`
		SELECT ref, type, date, number, reg_number, author
		FROM %[1]s
		WHERE ref = ?
	`, types.Report)
	report, err := m.Query[*m.Report](repo.source.db, context.Background(), m.NewReport, nil, query, guid)
	if err != nil {
		return nil, err
	}
	if len(report) != 1 {
		return nil, m.ErrNotFound
	}

	return report[0], nil
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
