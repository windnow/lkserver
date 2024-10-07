package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	m "lkserver/internal/models"
	"lkserver/internal/models/reports"
)

var (
	ErrWrongType = fmt.Errorf("INVALID REPORT DATA TYPE")
)

type DepartureOnBusinessTrip struct {
	source *src
}

func (r *DepartureOnBusinessTrip) GetStructure() interface{} {
	return &reports.ReportData{Details: &reports.BussinesTripDetails{}}
}

func (r *DepartureOnBusinessTrip) Get(ref m.JSONByte, tx ...*sql.Tx) (any, error) {
	details := &reports.BussinesTripDetails{}
	query := fmt.Sprintf(`
	SELECT report, supervisor, acting_supervisor, basis, transport_type
	FROM %[1]s WHERE report = ?
	`, tabBusinessTrip)

	row := r.source.db.QueryRow
	if len(tx) > 0 && tx[0] != nil {
		row = tx[0].QueryRow
	}

	if err := row(query, ref).Scan(
		&details.ReportRef,
		&details.Supervisor,
		&details.ActingSupervisor,
		&details.Basis,
		&details.TransportType,
	); err != nil {
		return nil, m.HandleError(err)
	}
	return details, nil
}
func (r *DepartureOnBusinessTrip) Save(tx *sql.Tx, ctx context.Context, report m.JSONByte, data any) error {
	details, ok := data.(*reports.BussinesTripDetails)
	if !ok {
		m.HandleError(ErrWrongType, "DepartureOnBusinessTrip.Save")
	}
	query := fmt.Sprintf(`
	INSERT OR REPLACE INTO %[1]s (
		report, supervisor, acting_supervisor, basis, transport_type	
	) VALUES (?, ?, ?, ?, ?)
	`, tabBusinessTrip)

	_, err := tx.ExecContext(ctx, query, report,
		details.Supervisor)
	if err != nil {
		m.HandleError(ErrWrongType, "DepartureOnBusinessTrip.Save")
	}
	return errors.ErrUnsupported
}
func (r *DepartureOnBusinessTrip) Init() error {

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %[1]s (
			report				BLOB,
			supervisor			BLOB,
			acting_supervisor	BLOB,
			basis				TEXT,
			transport_type		TEXT,

			FOREIGN KEY (report) REFERENCES %[2]s(ref),
			FOREIGN KEY (supervisor) REFERENCES %[3]s(ref),
			FOREIGN KEY (acting_supervisor) REFERENCES %[3]s(ref)
		);
	
		CREATE INDEX IF NOT EXISTS idx_%[1]s_report ON    %[1]s(report);
		`, tabBusinessTrip, tabReport, tabUsers)

	err := r.source.Exec(query)
	if err != nil {
		return m.HandleError(err, "sqliteRepo.initEducation")
	}
	return nil
}
