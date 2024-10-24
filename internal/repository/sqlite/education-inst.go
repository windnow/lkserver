package sqlite

import (
	"context"
	"encoding/json"
	"fmt"
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
)

type eduInstitutions struct {
	source *src
}

func (r *sqliteRepo) initEducationInstitutions() (err error) {
	eduRepo := &eduInstitutions{
		source: r.db,
	}

	if err := eduRepo.source.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %[1]s (
			ref BLOB PRIMARY KEY,
			title TEXT
		)
	`, types.Institutions)); err != nil {
		return m.HandleError(err, "sqliteRepo.initEducationInstitutions")
	}
	var count int64
	eduRepo.source.db.QueryRow(fmt.Sprintf(`select count(*) from %[1]s`, types.Institutions)).Scan(&count)
	var eduInst []m.EducationInstitution
	json.Unmarshal([]byte(mockInstitutions), &eduInst)
	if count == 0 {
		for _, ei := range eduInst {
			if err := eduRepo.Save(context.Background(), &ei); err != nil {
				return m.HandleError(err, "sqliteRepo.initEducationInstitutions")
			}
		}
	}
	eduRepo.Get(eduInst[1].Key)
	r.institutions = eduRepo

	return nil
}

func (eduRepo *eduInstitutions) Close() {}

func (i *eduInstitutions) Get(key m.JSONByte) (*m.EducationInstitution, error) {

	institut := &m.EducationInstitution{Key: key}
	err := i.source.db.QueryRow(fmt.Sprintf(`
		SELECT title
		FROM %[1]s
		WHERE ref = ? 
	`, types.Institutions), institut.Key).Scan(
		&institut.Title,
	)

	if err != nil {
		return nil, m.HandleError(err, "eduInstitution.Get")
	}

	return institut, nil

}

func (eduRepo *eduInstitutions) Save(ctx context.Context, ei *m.EducationInstitution) error {

	if ei.Key.Blank() {
		k, err := m.GenerateUUID()
		if err != nil {
			return m.HandleError(err, "EduInstitutions.Save")
		}
		ei.Key = k
	}

	return m.HandleError(eduRepo.source.ExecContextInTransaction(ctx, saveQuery, nil,
		ei.Key,
		ei.Title,
	))
}

var saveQuery string = fmt.Sprintf(`INSERT OR REPLACE INTO %[1]s (ref, title) VALUES (?, ?)`, types.Institutions)
var mockInstitutions string = `
[
    {"key": "521451f0-1c6a-4647-b27d-d2204cd9e992", "title": "РГКП «Актауский государственный университет имени Ш. Есенова»"},
    {"key": "c86bbe0f-22bd-4020-adb3-360dc44c936d", "title": "Алматинское высшее военное училище ВС РК"},
    {"key": "650b8bf2-8fd5-48d9-9160-59ab0f14b979", "title": "Национальный университет обороны"},
    {"key": "2dd3483c-cee2-4e93-9ee8-761c82b55d38", "title": "Алматинская республиканская школа Жас Улан им. Б.Момышулы"}
]
`
