package catalogs

import (
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
)

type Cato struct {
	Ref         m.JSONByte `json:"ref"`
	ParentRef   m.JSONByte `json:"parentRef"`
	Code        string     `json:"code"`
	K1          string     `json:"k1"`
	K2          string     `json:"k2"`
	K3          string     `json:"k3"`
	K4          string     `json:"k4"`
	K5          string     `json:"k5"`
	Description string     `json:"description"`
	Title       string     `json:"title"`
}

var CatoMETA = m.META{
	"ref":         m.Description{Type: types.Cato, Labels: map[string]string{"ru": "Идентификатор"}},
	"parentRef":   m.Description{Type: types.Cato, Labels: map[string]string{"ru": "Родитель"}},
	"code":        m.Description{Type: types.String, Labels: map[string]string{"ru": "Код"}},
	"k1":          m.Description{Type: types.String, Labels: map[string]string{"ru": "К1"}},
	"k2":          m.Description{Type: types.String, Labels: map[string]string{"ru": "К2"}},
	"k3":          m.Description{Type: types.String, Labels: map[string]string{"ru": "К3"}},
	"k4":          m.Description{Type: types.String, Labels: map[string]string{"ru": "К4"}},
	"k5":          m.Description{Type: types.String, Labels: map[string]string{"ru": "К5"}},
	"description": m.Description{Type: types.String, Labels: map[string]string{"ru": "Краткое наименование"}},
	"title":       m.Description{Type: types.String, Labels: map[string]string{"ru": "Полное наименование"}},
}
