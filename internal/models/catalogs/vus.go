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
	"ref":   types.Vus,
	"code":  types.Date,
	"title": types.String,
}
