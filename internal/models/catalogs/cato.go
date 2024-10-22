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
	"ref":         m.Desc(types.Cato, map[string]string{"ru": "Идентификатор"}, 0),
	"parentRef":   m.Desc(types.Cato, map[string]string{"ru": "Родитель"}, 8),
	"code":        m.Desc(types.String, map[string]string{"ru": "Код"}, 9),
	"k1":          m.Desc(types.String, map[string]string{"ru": "К1"}, 3),
	"k2":          m.Desc(types.String, map[string]string{"ru": "К2"}, 4),
	"k3":          m.Desc(types.String, map[string]string{"ru": "К3"}, 5),
	"k4":          m.Desc(types.String, map[string]string{"ru": "К4"}, 6),
	"k5":          m.Desc(types.String, map[string]string{"ru": "К5"}, 7),
	"description": m.Desc(types.String, map[string]string{"ru": "Краткое наименование"}, 1),
	"title":       m.Desc(types.String, map[string]string{"ru": "Полное наименование"}, 2),
}
