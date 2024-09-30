package sqlite

import (
	"context"
	"encoding/json"
	m "lkserver/internal/models"
)

type specialties struct {
	source *src
}

func (r *sqliteRepo) initSpecialties() (err error) {
	s := &specialties{
		source: r.db,
	}
	if err := s.source.Exec(`
		CREATE TABLE IF NOT EXISTS specialties (
			guid BLOB PRIMARY KEY,
			title TEXT
		)
	`); err != nil {
		return m.HandleError(err, "sqliteRepo.initEducationInstitutions")
	}

	var specs []m.Specialties
	json.Unmarshal([]byte(mockSpecialties), &specs)

	var count int64
	s.source.db.QueryRow(`select count(*) from specialties`).Scan(&count)

	if count == 0 {
		for _, spec := range specs {
			if err := s.Save(context.Background(), &spec); err != nil {
				return m.HandleError(err, "specialties.Get")
			}
		}
	}

	r.specialties = s
	return nil
}

func (s *specialties) Get(key m.JSONByte) (*m.Specialties, error) {

	spec := &m.Specialties{Key: key}
	err := s.source.db.QueryRow(`
		SELECT title
		FROM specialties
		WHERE guid = ? 
	`, spec.Key).Scan(
		&spec.Title,
	)

	if err != nil {
		return nil, m.HandleError(err, "specialties.Get")
	}

	return spec, nil

}

func (s *specialties) Close() {}

func (s *specialties) Save(ctx context.Context, spec *m.Specialties) error {
	if spec.Key.Blank() {
		k, err := m.GenerateUUID()
		if err != nil {
			return m.HandleError(err, "specialties.Save")
		}
		spec.Key = k
	}

	return m.HandleError(s.source.ExecContextInTransaction(ctx, saveSpecQuery,
		spec.Key,
		spec.Title,
	))

}

var saveSpecQuery = `INSERT OR REPLACE INTO specialties (guid, title) VALUES (?, ?)`

var mockSpecialties string = `
[
    {"key": "18fcf8f9-9705-460c-a15e-a7427c0d8f8c", "title": "юриспруденция"},
    {"key": "cf5ed9ea-332a-4688-a0b2-a3e3f155faf5", "title": "командная тактическая мотострелковых войск"},
    {"key": "bf95bccb-0232-46ad-b81d-cd45ce5c2952", "title": "менеджмент в военном деле"}
]
`
