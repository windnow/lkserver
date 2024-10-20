package catalogs

import (
	"database/sql"
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
)

type Devision struct {
	Ref         m.JSONByte `json:"ref"`
	OwnerRef    m.JSONByte `json:"owner_ref"`
	ParentRef   m.JSONByte `json:"parent_ref"`
	Code        string     `json:"code"`
	Description string     `json:"description"`
}

var DevisionMETA m.META = m.META{
	"ref":         m.Description{Type: types.Devision, Labels: map[string]string{"ru": "Идентификатор"}},
	"owner_ref":   m.Description{Type: types.Organization, Labels: map[string]string{"ru": "Владелец"}},
	"parent_ref":  m.Description{Type: types.Devision, Labels: map[string]string{"ru": "Родитель"}},
	"code":        m.Description{Type: types.String, Labels: map[string]string{"ru": "Код"}},
	"description": m.Description{Type: types.String, Labels: map[string]string{"ru": "Наименование"}},
}

func (d *Devision) Scan(rows *sql.Rows) error {
	return m.HandleError(rows.Scan(
		&d.Ref,
		&d.OwnerRef,
		&d.ParentRef,
		&d.Code,
		&d.Description,
	), "Devision.Scan")
}
