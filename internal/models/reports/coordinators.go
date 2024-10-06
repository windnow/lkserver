package reports

import m "lkserver/internal/models"

type Coordinators struct {
	Ref            m.JSONByte `json:"ref"`
	ReportRef      m.JSONByte `json:"report_ref"`
	CoordinatorRef m.JSONByte `json:"coordinator_ref"`
	WhoAuthor      m.JSONByte `json:"who_author_ref"`
	WhenAdded      m.JSONTime `json:"when_added_ref"`
}
