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
	"ref":         types.Cato,
	"parentRef":   types.Cato,
	"code":        types.String,
	"k1":          types.String,
	"k2":          types.String,
	"k3":          types.String,
	"k4":          types.String,
	"k5":          types.String,
	"description": types.String,
	"title":       types.String,
}
