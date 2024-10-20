package catalogs

import (
	"database/sql"
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
)

type Organization struct {
	Ref               m.JSONByte `json:"ref"`
	Code              string     `json:"code"`
	Description       string     `json:"description"`
	Title             string     `json:"title"`
	ShortTitle        string     `json:"short_title"`
	ConventionalTitle string     `json:"conventional_title"`
}

var OrganizationMETA m.META = m.META{
	"ref":                m.Description{Type: types.Organization, Labels: map[string]string{"ru": "Идентификатор"}},
	"code":               m.Description{Type: types.String, Labels: map[string]string{"ru": "Код"}},
	"description":        m.Description{Type: types.String, Labels: map[string]string{"ru": "Наименование"}},
	"title":              m.Description{Type: types.String, Labels: map[string]string{"ru": "Полное наименование"}},
	"short_title":        m.Description{Type: types.String, Labels: map[string]string{"ru": "Сокращенное юр. наименование"}},
	"conventional_title": m.Description{Type: types.String, Labels: map[string]string{"ru": "Условное наименование"}},
}

func (o *Organization) Scan(rows *sql.Rows) error {
	return m.HandleError(rows.Scan(
		&o.Ref,
		&o.Code,
		&o.Description,
		&o.Title,
		&o.ShortTitle,
		&o.ConventionalTitle,
	), "Organization.Scan")
}
