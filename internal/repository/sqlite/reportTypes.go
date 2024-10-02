package sqlite

import (
	"encoding/json"
	"errors"
	"fmt"
	m "lkserver/internal/models"
	"lkserver/internal/models/reports"
	r "lkserver/internal/models/reports"
)

var tab_name = "report_types"

func InitReportTypes(repo *reportsRepo) error {

	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %[1]s (
			ref BLOB PRIMARY KEY,
			parent BLOB,
			code TEXT UNIQUE,
			title TEXT	
		);
		CREATE INDEX IF NOT EXISTS idx_%[1]s_code ON %[1]s(code);
	`, tab_name)

	if _, err := repo.source.db.Exec(query); err != nil {
		return err
	}

	var mockTypes []*reports.ReportTypes
	if err := json.Unmarshal([]byte(data), &mockTypes); err != nil {
		return m.HandleError(err, "sqliteRepo.InitReportTypes")
	}
	dbTypes, err := repo.GetTypes([]string{})
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
		repo.SaveType(mType)
	}

	return nil

}

func (r *reportsRepo) GetTypes(types []string) ([]*r.ReportTypes, error) {

	return nil, errors.ErrUnsupported

}

func (r *reportsRepo) SaveType(reportType *r.ReportTypes) error {

	return errors.ErrUnsupported

}

var mockReportTypesData = `
[
{"ref":"fcf8e381-ea56-43ea-a83f-c2059a3aa329", "code": "000001", "title":"Об убытии в служебные командировки"}
]
`
