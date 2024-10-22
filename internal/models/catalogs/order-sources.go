package catalogs

import (
	"database/sql"
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
)

type OrderSource struct {
	Ref         m.JSONByte `json:"ref"`
	Number      int        `json:"num"`
	Description string     `json:"description"`
}

var OrderSourceMETA = m.META{
	"ref":         m.Description{Type: types.OrderSource, Labels: map[string]string{"ru": "Идентификатор"}},
	"num":         m.Description{Type: types.Number, Labels: map[string]string{"ru": "Номер"}},
	"description": m.Description{Type: types.String, Labels: map[string]string{"ru": "Наименование"}},
}

func (o *OrderSource) Scan(rows *sql.Rows) error {
	return m.HandleError(rows.Scan(&o.Ref, &o.Number, &o.Description), "OrderSource.Scan")
}
