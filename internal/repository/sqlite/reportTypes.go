package sqlite

import (
	"context"
	"encoding/json"
	"fmt"
	m "lkserver/internal/models"
	r "lkserver/internal/models/reports"
	"strings"
)

func InitReportTypes(repo *reportsRepo) error {

	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %[1]s (
			ref BLOB PRIMARY KEY,
			parent BLOB,
			code TEXT UNIQUE,
			title TEXT	
		);
		CREATE INDEX IF NOT EXISTS idx_%[1]s_code ON %[1]s(code);
	`, tabReportType)

	if _, err := repo.source.db.Exec(query); err != nil {
		return err
	}

	var mockTypes []*r.ReportTypes
	if err := json.Unmarshal([]byte(mockReportTypesData), &mockTypes); err != nil {
		return m.HandleError(err, "sqliteRepo.InitReportTypes")
	}
	ctx := context.Background()
	dbTypes, err := repo.GetTypes(ctx, []string{"УбытиеВСлужебнКомандировку"})
	if err != nil {
		return m.HandleError(err, "sqliteRepo.InitReportTypes")
	}
	codes := map[string]struct{ code string }{}
	for _, rType := range dbTypes {
		codes[rType.Code] = struct{ code string }{code: rType.Code}
	}

	for _, mType := range mockTypes {
		_, ok := codes[mType.Code]
		if ok {
			continue
		}
		if err := repo.SaveType(ctx, mType); err != nil {
			return m.HandleError(err, "sqliteRepo.InitReportTypes")
		}
	}

	return nil

}

func (repo *reportsRepo) GetTypes(ctx context.Context, types []string) ([]*r.ReportTypes, error) {

	var conditions = ""
	var args []any

	if len(types) > 0 {
		placeholders := make([]string, len(types))
		args = make([]any, len(types))
		for i, t := range types {
			placeholders[i] = "?"
			args[i] = t
		}
		conditions = fmt.Sprintf(` WHERE code in(%s)`, strings.Join(placeholders, ", "))
	}

	rows, err := repo.source.db.QueryContext(ctx, fmt.Sprintf(`SELECT ref, parent, code, title from %[1]s%[2]s`, tabReportType, conditions), args...)
	if err != nil {
		return nil, m.HandleError(err, "reportsRepo.GetTypes")
	}
	defer rows.Close()

	result := make([]*r.ReportTypes, 0, 10)
	for rows.Next() {
		value := &r.ReportTypes{}
		if err := rows.Scan(&value.Ref, &value.ParentRef, &value.Code, &value.Title); err != nil {
			return nil, m.HandleError(err, "reportsRepo.GetTypes")
		}
		result = append(result, value)
	}
	if err := rows.Err(); err != nil {
		return nil, m.HandleError(err, "reportsRepo.GetTypes")
	}

	return result, nil

}

func (repo *reportsRepo) SaveType(ctx context.Context, rt *r.ReportTypes) error {

	if rt.Ref.Blank() {
		ref, err := m.GenerateUUID()
		if err != nil {
			return m.HandleError(err, "reportsRepo.SaveTypes")
		}
		rt.Ref = ref
	}

	return m.HandleError(repo.source.ExecContextInTransaction(ctx,
		fmt.Sprintf(`INSERT OR REPLACE INTO %[1]s (ref, parent, code, title) VALUES (?, ?, ?, ?)`, tabReportType),

		rt.Ref,
		rt.ParentRef,
		rt.Code,
		rt.Title,
	))

}

var mockReportTypesData = `
[
{"ref":"fcf8e381-ea56-43ea-a83f-c2059a3aa329", "code": "0001", "title":"Об убытии в служебные командировки"}
]
`