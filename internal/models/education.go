package models

import "lkserver/internal/models/types"

type Education struct {
	Individual           *Individuals          `json:"individual"`
	EducationInstitution *EducationInstitution `json:"institution"`
	YearOfCompletion     int                   `json:"year"`
	Specialty            *Specialties          `json:"specialty"`
	Type                 string                `json:"type"` // military | civil
}

var EducationMETA = META{
	"individual":  Desc(types.Individuals, map[string]string{"ru": "Физическое лицо"}, 1),
	"institution": Desc(types.Institutions, map[string]string{"ru": "Учебное заведение"}, 2),
	"year":        Desc(types.Number, map[string]string{"ru": "Год окончания"}, 3),
	"specialty":   Desc(types.Specialties, map[string]string{"ru": "Специальность"}, 4),
	"type":        Desc(types.String, map[string]string{"ru": "Тип"}, 4),
}
