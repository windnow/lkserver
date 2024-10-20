package sqlite

import (
	"context"
	"encoding/json"
	"fmt"
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
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
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %[1]s (
			individ BLOB,
			institut BLOB,
			year INTEGER,
			type TEXT,
			specialty BLOB,

			FOREIGN KEY (individ) REFERENCES %[2]s(ref),
			FOREIGN KEY (institut) REFERENCES %[3]s(ref),
			FOREIGN KEY (specialty) REFERENCES %[4]s(ref)

		);

		CREATE INDEX IF NOT EXISTS idx_%[1]s_individ ON    %[1]s(individ);
		CREATE INDEX IF NOT EXISTS idx_%[1]s_institut ON   %[1]s(institut);
		CREATE INDEX IF NOT EXISTS idx_%[1]s_specialty ON  %[1]s(specialty);
		CREATE INDEX IF NOT EXISTS idx_%[1]s_type ON       %[1]s(type);
	`, types.Education, types.Individuals, types.Institutions, types.Specialties)
	err := edu.source.Exec(query)
	if err != nil {
		return m.HandleError(err, "sqliteRepo.initEducation")
	}

	var mockEdu []*draftEdu
	if err := json.Unmarshal([]byte(mockEducation), &mockEdu); err != nil {
		return m.HandleError(err, "sqliteRepo.initEducation")
	}
	var count int64
	edu.source.db.QueryRow(fmt.Sprintf(`select count(*) from %[1]s`, types.Education)).Scan(&count)
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

	rows, err := e.source.db.QueryContext(ctx, fmt.Sprintf("SELECT individ, institut, year, type, specialty FROM %[1]s WHERE individ = ?", types.Education), individ.Key)
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
	return m.HandleError(e.source.ExecContextInTransaction(ctx, insertEduQuery, nil,
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

var insertEduQuery string = fmt.Sprintf(`
INSERT OR REPLACE INTO %[1]s (
			individ,
			institut,
			year,
			type,
			specialty)
VALUES (?, ?, ?, ?, ?)
`, types.Education)
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
