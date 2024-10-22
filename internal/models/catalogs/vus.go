package catalogs

import (
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
)

type Vus struct {
	Ref   m.JSONByte `json:"ref"`
	Code  string     `json:"code"`
	Title string     `json:"title"`
}

var VusMETA = m.META{
	"ref":   m.Desc(types.Vus, map[string]string{"ru": "Идентификатор"}, 0),
	"code":  m.Desc(types.String, map[string]string{"ru": "Код"}, 1),
	"title": m.Desc(types.String, map[string]string{"ru": "Наименование"}, 2),
}
