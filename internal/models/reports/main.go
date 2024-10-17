package reports

import (
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
)

type ReportTypes struct {
	Ref       m.JSONByte `json:"ref"`
	ParentRef m.JSONByte `json:"parent"`
	Code      string     `json:"code"`
	Title     string     `json:"title"`
}

var ReportTypesMETA = m.META{
	"ref":    m.Description{Type: types.ReportType, Labels: map[string]string{"ru": "Идентификатор"}},
	"parent": m.Description{Type: types.ReportType, Labels: map[string]string{"ru": "Родитель"}},
	"code":   m.Description{Type: types.String, Labels: map[string]string{"ru": "Код"}},
	"title":  m.Description{Type: types.String, Labels: map[string]string{"ru": "Наименование"}},
}

type ReportData struct {
	Head         *m.Report       `json:"head"`
	Coordinators []*Coordinators `json:"coordinators"`
	Details      any             `json:"details"`
}
