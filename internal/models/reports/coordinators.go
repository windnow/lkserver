package reports

import (
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
)

type Coordinators struct {
	Ref            m.JSONByte `json:"ref"`
	ReportRef      m.JSONByte `json:"report_ref"`
	CoordinatorRef m.JSONByte `json:"coordinator_ref"`
	WhoAuthor      m.JSONByte `json:"who_author_ref"`
	WhenAdded      m.JSONTime `json:"when_added"`
}

var CoordinatorsMETA = m.META{
	"ref":             types.Coordinators,
	"report_ref":      types.Report,
	"coordinator_ref": types.Users,
	"who_author_ref":  types.Users,
	"when_added":      types.Date,
}
