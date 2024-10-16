package sqlite

import (
	"context"
	"encoding/json"
	"fmt"
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
)

type specialties struct {
	source *src
}

func (r *sqliteRepo) initSpecialties() (err error) {
	s := &specialties{
		source: r.db,
	}
	if err := s.source.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %[1]s (
			ref BLOB PRIMARY KEY,
			title TEXT
		)
	`, types.Specialties)); err != nil {
		return m.HandleError(err, "sqliteRepo.initEducationInstitutions")
	}

	var specs []m.Specialties
	json.Unmarshal([]byte(mockSpecialties), &specs)

	var count int64
	s.source.db.QueryRow(fmt.Sprintf(`select count(*) from %[1]s`, types.Specialties)).Scan(&count)

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
	err := s.source.db.QueryRow(fmt.Sprintf(`
		SELECT title
		FROM %[1]s
		WHERE ref = ? 
	`, types.Specialties), spec.Key).Scan(
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

	return m.HandleError(s.source.ExecContextInTransaction(ctx, saveSpecQuery, nil,
		spec.Key,
		spec.Title,
	))

}

var saveSpecQuery = fmt.Sprintf(`INSERT OR REPLACE INTO %[1]s (ref, title) VALUES (?, ?)`, types.Specialties)

var mockSpecialties string = `
[
    {"key": "18fcf8f9-9705-460c-a15e-a7427c0d8f8c", "title": "юриспруденция"},
    {"key": "cf5ed9ea-332a-4688-a0b2-a3e3f155faf5", "title": "командная тактическая мотострелковых войск"},
    {"key": "bf95bccb-0232-46ad-b81d-cd45ce5c2952", "title": "менеджмент в военном деле"}
]
`
