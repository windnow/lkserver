package reports

import (
	"database/sql"
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
)

type ReportTypes struct {
	Ref       m.JSONByte `json:"ref"`
	ParentRef m.JSONByte `json:"parent"`
	Code      string     `json:"code"`
	Title     string     `json:"title"`
}

func NewReportType() m.Scanable {
	return &ReportTypes{}
}

func (r *ReportTypes) Scan(rows *sql.Rows) error {
	return rows.Scan(
		&r.Ref,
		&r.ParentRef,
		&r.Code,
		&r.Title,
	)
}

var ReportTypesMETA = m.META{
	"ref":    m.Desc(types.ReportType, map[string]string{"ru": "Идентификатор"}, 0),
	"parent": m.Desc(types.ReportType, map[string]string{"ru": "Родитель"}, 1),
	"code":   m.Desc(types.String, map[string]string{"ru": "Код"}, 2),
	"title":  m.Desc(types.String, map[string]string{"ru": "Наименование"}, 3),
}

type ReportData struct {
	Head         *m.Report       `json:"head"`
	Coordinators []*Coordinators `json:"coordinators"`
	Details      any             `json:"details"`
}
