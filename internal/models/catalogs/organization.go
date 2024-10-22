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
	"ref":                m.Desc(types.Organization, map[string]string{"ru": "Идентификатор"}, 0),
	"code":               m.Desc(types.String, map[string]string{"ru": "Код"}, 1),
	"description":        m.Desc(types.String, map[string]string{"ru": "Наименование"}, 2),
	"title":              m.Desc(types.String, map[string]string{"ru": "Полное наименование"}, 3),
	"short_title":        m.Desc(types.String, map[string]string{"ru": "Сокращенное юр. наименование"}, 4),
	"conventional_title": m.Desc(types.String, map[string]string{"ru": "Условное наименование"}, 5),
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
