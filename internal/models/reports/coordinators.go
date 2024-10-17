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
	"ref":             m.Description{Type: types.Coordinators, Labels: map[string]string{"ru": "Идентификатор"}},
	"report_ref":      m.Description{Type: types.Report, Labels: map[string]string{"ru": "Идентификатор рапорта"}},
	"coordinator_ref": m.Description{Type: types.Users, Labels: map[string]string{"ru": "Согласующий"}},
	"who_author_ref":  m.Description{Type: types.Users, Labels: map[string]string{"ru": "Кто добавил"}},
	"when_added":      m.Description{Type: types.Date, Labels: map[string]string{"ru": "Дата добавления"}},
}
