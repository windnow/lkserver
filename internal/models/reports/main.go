package reports

import (
	"lkserver/internal/models"
	m "lkserver/internal/models"
)

type ReportTypes struct {
	Ref       m.JSONByte `json:"ref"`
	ParentRef m.JSONByte `json:"parent"`
	Code      string     `json:"code"`
	Title     string     `json:"title"`
}

type ReportData struct {
	Head         *models.Report  `json:"head"`
	Coordinators []*Coordinators `json:"coordinators"`
	Details      any             `json:"details"`
}
