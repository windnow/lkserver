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
	"ref":         m.Desc(types.Devision, map[string]string{"ru": "Идентификатор"}, 0),
	"owner_ref":   m.Desc(types.Organization, map[string]string{"ru": "Владелец"}, 4),
	"parent_ref":  m.Desc(types.Devision, map[string]string{"ru": "Родитель"}, 2),
	"code":        m.Desc(types.String, map[string]string{"ru": "Код"}, 3),
	"description": m.Desc(types.String, map[string]string{"ru": "Наименование"}, 1),
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
