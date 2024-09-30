package sqlite

import (
	"context"
	"encoding/json"
	m "lkserver/internal/models"
	"strings"
)

type education struct {
	source *src
	repo   *sqliteRepo
}

func (s *sqliteRepo) initEducation() error {
	edu := &education{
		source: s.db,
		repo:   s,
	}
	err := edu.source.Exec(`
		CREATE TABLE IF NOT EXISTS education (
			individ BLOB,
			institut BLOB,
			year INTEGER,
			type TEXT,
			specialty BLOB
		)
	`)
	if err != nil {
		return m.HandleError(err, "sqliteRepo.initEducation")
	}
	err = edu.source.Exec(`
		CREATE INDEX IF NOT EXISTS idx_education_individ ON education(individ);
		CREATE INDEX IF NOT EXISTS idx_education_institut ON education(institut);
		CREATE INDEX IF NOT EXISTS idx_education_specialty ON education(specialty);
		CREATE INDEX IF NOT EXISTS idx_education_type ON education(type);
	`)
	if err != nil {
		return m.HandleError(err, "sqliteRepo.initRankHistory")
	}

	var mockEdu []*draftEdu
	json.Unmarshal([]byte(mockEducation), &mockEdu)
	var count int64
	edu.source.db.QueryRow(`select count(*) from education`).Scan(&count)
	if count == 0 {
		for _, row := range mockEdu {
			record, err := row.update(edu.repo)
			if err != nil {
				return m.HandleError(err, "sqliteRepo.initEducation")
			}
			err = edu.Save(context.Background(), record)
			if err != nil {
				return m.HandleError(err, "sqliteRepo.initEducation")
			}
		}
	}

	edu.GetByIin(mockEdu[1].Iin)

	s.education = edu

	return nil
}

func (e *education) Close() {}

func (e *education) GetByIin(iin string) ([]*m.Education, error) {
	individ, err := e.repo.individuals.GetByIin(iin)
	if err != nil {
		return nil, m.HandleError(err, "education.GetByIin")
	}

	edu, err := e.Get(context.Background(), individ)
	if err != nil {
		return nil, m.HandleError(err, "education.GetByIin")
	}

	return edu, nil

}

func (e *education) Get(ctx context.Context, individ *m.Individuals) ([]*m.Education, error) {

	rows, err := e.source.db.QueryContext(ctx, "SELECT individ, institut, year, type, specialty FROM education WHERE individ = ?", individ.Key)
	if err != nil {
		return nil, m.HandleError(err, "education.Get")
	}
	defer rows.Close()
	select {
	case <-ctx.Done():
		return nil, m.HandleError(ctx.Err(), "education.Get")
	default:
	}
	var records []*m.Education
	for rows.Next() {
		record := &draftEdu{}
		if err := rows.Scan(&record.Individ, &record.Institut, &record.Year, &record.Type, &record.Specialty); err != nil {
			return nil, m.HandleError(ctx.Err(), "education.Get")
		}
		edu, err := record.update(e.repo)
		if err != nil {
			return nil, m.HandleError(ctx.Err(), "education.Get")
		}
		edu.Individual = individ
		records = append(records, edu)
	}

	return records, nil
}

func (e *education) Save(ctx context.Context, edu *m.Education) error {
	return m.HandleError(e.source.ExecContextInTransaction(ctx, insertEduQuery,
		edu.Individual.Key,
		edu.EducationInstitution.Key,
		edu.YearOfCompletion,
		edu.Type,
		edu.Specialty.Key,
	))
}

type draftEdu struct {
	Iin       string     `json:"individual"`
	Individ   m.JSONByte `json:"-"`
	Institut  m.JSONByte `json:"institution"`
	Year      int        `json:"year"`
	Specialty m.JSONByte `json:"specialty"`
	Type      string     `json:"type"`
}

func (draft *draftEdu) update(r *sqliteRepo) (*m.Education, error) {
	var individ *m.Individuals
	spec, err := r.specialties.Get(draft.Specialty)
	if err != nil {
		return nil, m.HandleError(err)
	}
	institut, err := r.institutions.Get(draft.Institut)
	if err != nil {
		return nil, m.HandleError(err)
	}
	if strings.TrimSpace(draft.Iin) != "" {
		individ, err = r.individuals.GetByIin(draft.Iin)
		if err != nil {
			return nil, m.HandleError(err)
		}
	}

	return &m.Education{
		YearOfCompletion:     draft.Year,
		Type:                 draft.Type,
		Individual:           individ,
		EducationInstitution: institut,
		Specialty:            spec,
	}, nil
}

var insertEduQuery string = `
INSERT OR REPLACE INTO education (
			individ,
			institut,
			year,
			type,
			specialty)
VALUES (?, ?, ?, ?, ?)
`
var mockEducation string = `
[
    {"individual": "821019000888", "institution": "521451f0-1c6a-4647-b27d-d2204cd9e992", "year": 2007, "type": "civil", "specialty": "18fcf8f9-9705-460c-a15e-a7427c0d8f8c"},
    {"individual": "821019000888", "institution": "c86bbe0f-22bd-4020-adb3-360dc44c936d", "year": 2003, "type": "military", "specialty": "cf5ed9ea-332a-4688-a0b2-a3e3f155faf5"},
    {"individual": "821019000888", "institution": "650b8bf2-8fd5-48d9-9160-59ab0f14b979", "year": 2019, "type": "military", "specialty": "bf95bccb-0232-46ad-b81d-cd45ce5c2952"},
    {"individual": "910702000888", "institution": "2dd3483c-cee2-4e93-9ee8-761c82b55d38", "year": 2018, "type": "civil", "specialty": "18fcf8f9-9705-460c-a15e-a7427c0d8f8c"},
    {"individual": "910702000888", "institution": "c86bbe0f-22bd-4020-adb3-360dc44c936d", "year": 2003, "type": "military", "specialty": "cf5ed9ea-332a-4688-a0b2-a3e3f155faf5"},
    {"individual": "910702000888", "institution": "650b8bf2-8fd5-48d9-9160-59ab0f14b979", "year": 2019, "type": "military", "specialty": "bf95bccb-0232-46ad-b81d-cd45ce5c2952"},
    {"individual": "851204000888", "institution": "2dd3483c-cee2-4e93-9ee8-761c82b55d38", "year": 2018, "type": "civil", "specialty": "18fcf8f9-9705-460c-a15e-a7427c0d8f8c"},
    {"individual": "851204000888", "institution": "c86bbe0f-22bd-4020-adb3-360dc44c936d", "year": 2003, "type": "military", "specialty": "bf95bccb-0232-46ad-b81d-cd45ce5c2952"},
    {"individual": "851204000888", "institution": "650b8bf2-8fd5-48d9-9160-59ab0f14b979", "year": 2019, "type": "military", "specialty": "bf95bccb-0232-46ad-b81d-cd45ce5c2952"}
]`
