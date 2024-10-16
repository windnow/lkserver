package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	m "lkserver/internal/models"
	"lkserver/internal/models/catalogs"
	"lkserver/internal/models/types"
)

type vus struct {
	source *src
}

func (repo *sqliteRepo) initVus() error {
	src := &vus{
		source: repo.db,
	}

	query := fmt.Sprintf(createVusQuery, types.Vus)
	if err := src.source.Exec(query); err != nil {
		return m.HandleError(err, "sqliteRepo.initVus")
	}

	var count int64
	src.source.db.QueryRow(fmt.Sprintf(`select count(*) from %[1]s`, types.Vus)).Scan(&count)
	if count == 0 {
		ctx := context.Background()
		var vusSet []*catalogs.Vus
		if err := json.Unmarshal([]byte(mockVusData), &vusSet); err != nil {
			return m.HandleError(err, "sqliteRepo.initVus")
		}
		for _, vus := range vusSet {
			if err := src.Save(ctx, vus, nil); err != nil {
				return m.HandleError(err, "sqliteRepo.initVus")
			}
		}
	}

	repo.catalogs.Vus = src
	return nil
}

func (s *vus) Get(ctx context.Context, Ref m.JSONByte) (*catalogs.Vus, error) {
	query := fmt.Sprintf("SELECT ref, code, title from %[1]s WHERE ref = ?", types.Vus)
	result, err := s.query(ctx, query, Ref)
	if err != nil {
		return nil, m.HandleError(err, "vus.Get")
	}

	if len(result) != 1 {
		m.HandleError(errors.New("WRONG RESULT LENGTH"), "vus.Get")
	}

	return result[0], nil
}

func (c *vus) Count(ctx context.Context) int64 {
	var count int64
	c.source.db.QueryRow(fmt.Sprintf(`select count(ref) from %[1]s`, types.Vus)).Scan(&count)
	return count
}

func (s *vus) query(ctx context.Context, query string, args ...any) ([]*catalogs.Vus, error) {
	stmt, err := s.source.db.Prepare(query)
	if err != nil {
		return nil, m.HandleError(err, "vus.query")
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, m.HandleError(err, "vus.query")
	}
	defer rows.Close()
	var result []*catalogs.Vus
	for rows.Next() {
		record := &catalogs.Vus{}
		if err := rows.Scan(
			&record.Ref,
			&record.Code,
			&record.Title,
		); err != nil {
			return nil, m.HandleError(err, "cato.Scan")
		}
		result = append(result, record)
	}

	return result, nil
}

func (s *vus) List(ctx context.Context, limits ...int64) ([]*catalogs.Vus, error) {
	limit, offset := limitations(limits)

	query := fmt.Sprintf("SELECT ref, code, title from %[1]s LIMIT %[2]d OFFSET %[3]d", types.Vus, limit, offset)
	result, err := s.query(ctx, query)
	if err != nil {
		return nil, m.HandleError(err, "vus.List")
	}

	return result, nil
}

func (s *vus) Find(ctx context.Context, pattern string, limits ...int64) ([]*catalogs.Vus, error) {
	limit, offset := limitations(limits)
	query := fmt.Sprintf("SELECT ref, code, title from %[1]s WHERE code like ? OR title like ? LIMIT %[2]d OFFSET %[3]d", types.Vus, limit, offset)
	pattern = "%" + pattern + "%"
	result, err := s.query(ctx, query, pattern, pattern)
	if err != nil {
		return nil, m.HandleError(err, "vus.List")
	}

	return result, nil
}

func (s *vus) Save(ctx context.Context, vus *catalogs.Vus, tx *sql.Tx) error {
	return s.source.ExecContextInTransaction(ctx, fmt.Sprintf("INSERT INTO %[1]s (ref, code, title) VALUES (?, ?, ?)", types.Vus), tx,
		vus.Ref,
		vus.Code,
		vus.Title)
}

var createVusQuery string = `
CREATE TABLE IF NOT EXISTS %[1]s (
	ref   BLOB PRIMARY KEY,
	code  TEXT,
	title TEXT
);


CREATE INDEX IF NOT EXISTS idx_%[1]s_code ON %[1]s(code);
`
var mockVusData string = `
[
	{"ref": "5e982366-826f-4b16-804d-178dac0b4ff9", "code": "7654321", "title":"Применение подразделений автоматизированных средств управления зенитными ракетными комплексами и зенитной артиллерией войсковой противовоздушной обороны"}
]
`
