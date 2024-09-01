package models

type Education struct {
	Individual           *Individuals          `json:"individual"`
	EducationInstitution *EducationInstitution `json:"institution"`
	YearOfCompletion     int                   `json:"year"`
	Specialty            *Specialties          `json:"specialty"`
	Type                 string                `json:"type"`
}
