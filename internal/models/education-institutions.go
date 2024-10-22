package models

import "lkserver/internal/models/types"

type EducationInstitution struct {
	Key   JSONByte `json:"key"`
	Title string   `json:"title"`
}

var EducationInstitutionMETA = META{
	"ref":   Desc(types.Coordinators, map[string]string{"ru": "Идентификатор"}, 0),
	"title": Desc(types.String, map[string]string{"ru": "Название"}, 1),
}
