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
	"ref":    types.ReportType,
	"parent": types.ReportType,
	"code":   types.String,
	"title":  types.String,
}

type ReportData struct {
	Head         *m.Report       `json:"head"`
	Coordinators []*Coordinators `json:"coordinators"`
	Details      any             `json:"details"`
}
