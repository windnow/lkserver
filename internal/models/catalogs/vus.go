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
	"ref":   m.Description{Type: types.Vus, Labels: map[string]string{"ru": "Идентификатор"}},
	"code":  m.Description{Type: types.Vus, Labels: map[string]string{"ru": "Код"}},
	"title": m.Description{Type: types.Vus, Labels: map[string]string{"ru": "Наименование"}},
}
