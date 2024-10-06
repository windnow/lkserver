package reports

import (
	"database/sql"
	m "lkserver/internal/models"
)

type ReportDetails interface {
	Get(ref m.JSONByte, reportData any) error
	Save(tx *sql.Tx, report *m.Report, reporData any) error
	Init() error
}
