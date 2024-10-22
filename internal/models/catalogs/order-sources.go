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
	"ref":         m.Desc(types.OrderSource, map[string]string{"ru": "Идентификатор"}, 0),
	"num":         m.Desc(types.Number, map[string]string{"ru": "Номер"}, 1),
	"description": m.Desc(types.String, map[string]string{"ru": "Наименование"}, 2),
}

func (o *OrderSource) Scan(rows *sql.Rows) error {
	return m.HandleError(rows.Scan(&o.Ref, &o.Number, &o.Description), "OrderSource.Scan")
}
