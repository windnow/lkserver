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
	"ref":             m.Desc(types.Coordinators, map[string]string{"ru": "Идентификатор"}, 0),
	"report_ref":      m.Desc(types.Report, map[string]string{"ru": "Идентификатор рапорта"}, 1),
	"coordinator_ref": m.Desc(types.Users, map[string]string{"ru": "Согласующий"}, 2),
	"who_author_ref":  m.Desc(types.Users, map[string]string{"ru": "Кто добавил"}, 3),
	"when_added":      m.Desc(types.Date, map[string]string{"ru": "Дата добавления"}, 4),
}
