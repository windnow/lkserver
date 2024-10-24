package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	m "lkserver/internal/models"
	"lkserver/internal/models/catalogs"
	"lkserver/internal/models/types"
	"os"
)

type cato struct {
	source *src
}

func (repo *sqliteRepo) initCato() error {

	src := &cato{
		source: repo.db,
	}
	query := fmt.Sprintf(createCatoQuery, types.Cato)

	if err := src.source.Exec(query); err != nil {
		return m.HandleError(err, "sqliteRepo.initCato")
	}

	var count int64
	src.source.db.QueryRow(fmt.Sprintf(`select count(*) from %[1]s`, types.Cato)).Scan(&count)
	if count == 0 {
		if err := src.loadData("data/cato.json"); err != nil {
			return m.HandleError(err, "sqliteRepo.initCato")
		}
	}

	repo.catalogs.Cato = src
	return nil
}

func (c *cato) loadData(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return m.HandleError(err, "cato.loadData")
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return m.HandleError(err, "cato.loadData")
	}
	data := []*catalogs.Cato{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return m.HandleError(err, "cato.loadData")
	}

	tx, err := c.source.db.Begin()
	if err != nil {
		return m.HandleError(err, "cato.loadData")
	}
	defer tx.Rollback()

	len := len(data)
	batchSize := 50
	for i := 0; i < len; i += batchSize {
		end := i + batchSize
		if end > len {
			end = len
		}
		if err := c.saveData(tx, data[i:end]); err != nil {
			return m.HandleError(err, "cato.loadData")
		}
	}
	return m.HandleError(tx.Commit(), "cato.loadData")
}

func (c *cato) saveData(tx *sql.Tx, data []*catalogs.Cato) error {
	query := fmt.Sprintf("INSERT INTO %[1]s (ref, parentRef, code, description, title, k1, k2, k3, k4, k5) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", types.Cato)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return m.HandleError(err, "cato.saveData")
	}
	defer stmt.Close()

	for _, cc := range data {
		_, err := stmt.Exec(cc.Ref, cc.ParentRef, cc.Code, cc.Description, cc.Title, cc.K1, cc.K2, cc.K3, cc.K4, cc.K5)
		if err != nil {
			return m.HandleError(err, "cato.saveData")
		}
	}

	return nil
}

func (c *cato) Get(ctx context.Context, Ref m.JSONByte) (*catalogs.Cato, error) {

	query := fmt.Sprintf(`SELECT ref, parentRef, code, description, title, k1, k2, k3, k4, k5 FROM %[1]s WHERE ref = ?`, types.Cato)

	result, err := c.query(ctx, query, Ref)
	if err != nil {
		return nil, m.HandleError(err, "cato.Get")
	}
	if len(result) != 1 {
		m.HandleError(errors.New("WRONG RESULT LENGTH"), "cato.Get")
	}

	return result[0], nil
}

func (c *cato) List(ctx context.Context, parentRef m.JSONByte, limits ...int64) ([]*catalogs.Cato, error) {

	limit, offset := limitations(limits)

	var result []*catalogs.Cato
	var err error

	if parentRef.Blank() {
		query := fmt.Sprintf(`SELECT ref, parentRef, code, description, title, k1, k2, k3, k4, k5 FROM %[1]s WHERE parentRef IS NULL LIMIT %[2]d OFFSET %[3]d`, types.Cato, limit, offset)
		result, err = c.query(ctx, query)
	} else {
		query := fmt.Sprintf(`SELECT ref, parentRef, code, description, title, k1, k2, k3, k4, k5 FROM %[1]s WHERE parentRef = ? LIMIT %[2]d OFFSET %[3]d`, types.Cato, limit, offset)
		result, err = c.query(ctx, query, parentRef)
	}

	if err != nil {
		return nil, m.HandleError(err, "cato.List")
	}
	return result, nil

}

func (c *cato) Count(ctx context.Context) uint64 {
	var count uint64
	if err := c.source.db.QueryRow(fmt.Sprintf(`select count(ref) from %[1]s`, types.Cato)).Scan(&count); err != nil {
		return 0
	}
	return count
}

func (c *cato) Find(ctx context.Context, title string, limits ...int64) ([]*catalogs.Cato, error) {

	limit, offset := limitations(limits)
	query := fmt.Sprintf(`SELECT ref, parentRef, code, description, title, k1, k2, k3, k4, k5 FROM %[1]s WHERE title like ? LIMIT %[2]d OFFSET %[3]d`, types.Cato, limit, offset)
	result, err := c.query(ctx, query, "%"+title+"%")
	if err != nil {
		return nil, m.HandleError(err, "cato.Find")
	}

	return result, nil

}

func (c *cato) query(ctx context.Context, query string, args ...any) ([]*catalogs.Cato, error) {
	stmt, err := c.source.db.Prepare(query)
	if err != nil {
		return nil, m.HandleError(err, "cato.query")
	}
	defer stmt.Close()
	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, m.HandleError(err, "cato.query")
	}
	defer rows.Close()
	var result []*catalogs.Cato
	for rows.Next() {
		record := &catalogs.Cato{}
		if err := rows.Scan(
			&record.Ref,
			&record.ParentRef,
			&record.Code,
			&record.Description,
			&record.Title,
			&record.K1,
			&record.K2,
			&record.K3,
			&record.K4,
			&record.K5,
		); err != nil {
			return nil, m.HandleError(err, "cato.query")
		}
		result = append(result, record)
	}

	return result, nil
}

func limitations(vals []int64) (int64, int64) {
	var LIMIT int64 = 20
	var OFFSET int64 = 0

	if len(vals) > 0 && vals[0] > 0 {
		LIMIT = vals[0]
	}
	if len(vals) > 1 && vals[1] > 0 {
		OFFSET = vals[1]
	}

	return LIMIT, OFFSET
}

var createCatoQuery string = `
CREATE TABLE IF NOT EXISTS %[1]s (
	ref         BLOB PRIMARY KEY,
	parentRef   BLOB,
	code        TEXT,
	description TEXT,
	title       TEXT,
	k1          TEXT,
	k2          TEXT,
	k3          TEXT,
	k4          TEXT,
	k5          TEXT
);

CREATE INDEX IF NOT EXISTS idx_%[1]s_code ON %[1]s(code);
CREATE INDEX IF NOT EXISTS idx_%[1]s_description ON %[1]s(description);
`
