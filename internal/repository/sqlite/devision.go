package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	m "lkserver/internal/models"
	"lkserver/internal/models/catalogs"
	"lkserver/internal/models/types"
	"os"
	"strings"
)

type devision struct {
	source *src
}

func (s *sqliteRepo) initDevision() error {
	src := &devision{
		source: s.db,
	}

	if err := src.source.Exec(createDevisionsQuery); err != nil {
		return m.HandleError(err, "sqliteRepo.initDevision")
	}

	var count int64
	src.source.db.QueryRow(fmt.Sprintf(`select count(*) from %[1]s`, types.Devision)).Scan(&count)
	if count == 0 {
		if err := src.loadData("data/org_struct.json"); err != nil {
			return m.HandleError(err, "sqliteRepo.initDevision")
		}
	}
	s.catalogs.Devision = src

	return nil
}

func (o *devision) loadData(path string) error {

	file, err := os.Open(path)
	if err != nil {
		return m.HandleError(err, "devision.loadData")
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return m.HandleError(err, "devision.loadData")
	}
	data := []*catalogs.Devision{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return m.HandleError(err, "devision.loadData")
	}

	tx, err := o.source.db.Begin()
	if err != nil {
		return m.HandleError(err, "devision.loadData")
	}
	defer tx.Rollback()

	len := len(data)
	batchSize := 50
	for i := 0; i < len; i += batchSize {
		end := i + batchSize
		if end > len {
			end = len
		}
		if err := o.saveData(tx, data[i:end]); err != nil {
			return m.HandleError(err, "devision.loadData")
		}
	}
	return m.HandleError(tx.Commit(), "devision.loadData")
}

func (o *devision) saveData(tx *sql.Tx, data []*catalogs.Devision) error {
	query := fmt.Sprintf("INSERT INTO %[1]s (ref, owner_ref, parent_ref, code, description) VALUES (?, ?, ?, ?, ?)", types.Devision)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return m.HandleError(err, "devision.saveData")
	}
	defer stmt.Close()

	for _, cc := range data {
		_, err := stmt.Exec(
			cc.Ref,
			cc.OwnerRef,
			cc.ParentRef,
			strings.TrimSpace(cc.Code),
			strings.TrimSpace(cc.Description),
		)
		if err != nil {
			return m.HandleError(err, "devision.saveData")
		}
	}

	return nil
}

func newDevision() m.Scanable {
	return &catalogs.Devision{}
}

func (o *devision) Get(ctx context.Context, Ref m.JSONByte) (*catalogs.Devision, error) {

	query := fmt.Sprintf(`SELECT ref, owner_ref, parent_ref, code, description FROM %[1]s WHERE ref = ?`, types.Devision)
	rows, err := m.Query[*catalogs.Devision](o.source.db, ctx, newDevision, nil, query, Ref)
	if err != nil {
		return nil, m.HandleError(err, "devision.Get")
	}
	if len(rows) == 0 {
		return nil, m.HandleError(m.ErrNotFound, "devision.Get")
	}

	return rows[0], nil
}

func (o *devision) List(ctx context.Context, limits ...int64) ([]*catalogs.Devision, error) {
	limit, offset := limitations(limits)
	query := fmt.Sprintf(`SELECT ref, owner_ref, parent_ref, code, description FROM %[1]s LIMIT %[2]d OFFSET %[3]d`, types.Devision, limit, offset)
	rows, err := m.Query[*catalogs.Devision](o.source.db, ctx, newDevision, nil, query)
	if err != nil {
		return nil, m.HandleError(err, "devision.List")
	}
	return rows, nil
}

func (o *devision) Find(ctx context.Context, pattern string, limits ...int64) ([]*catalogs.Devision, error) {
	limit, offset := limitations(limits)
	query := fmt.Sprintf(`SELECT ref, owner_ref, parent_ref, code, description FROM %[1]s
    WHERE code LIKE ? OR description LIKE ?
	LIMIT %[2]d OFFSET %[3]d`, types.Devision, limit, offset)
	param := "%" + pattern + "%"
	rows, err := m.Query[*catalogs.Devision](o.source.db, ctx, newDevision, nil, query, param, param)
	if err != nil {
		return nil, m.HandleError(err, "devision.Find")
	}
	return rows, nil
}

func (o *devision) Save(ctx context.Context, dev *catalogs.Devision, tx *sql.Tx) error {
	return m.HandleError(o.saveData(tx, []*catalogs.Devision{dev}))
}

func (o *devision) Count(ctx context.Context) int64 {
	var count int64
	o.source.db.QueryRow(fmt.Sprintf(`select count(ref) from %[1]s`, types.Devision)).Scan(&count)
	return count
}

var createDevisionsQuery = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %[1]s (
    ref                BLOB PRIMARY KEY,
    owner_ref          BLOB,
    parent_ref         BLOB,
    code               TEXT,
    description        TEXT
);

CREATE INDEX IF NOT EXISTS idx_%[1]s_owner_ref ON %[1]s(owner_ref);
CREATE INDEX IF NOT EXISTS idx_%[1]s_parent_ref ON %[1]s(parent_ref);
CREATE INDEX IF NOT EXISTS idx_%[1]s_code ON %[1]s(code);
CREATE INDEX IF NOT EXISTS idx_%[1]s_description ON %[1]s(description);
`, types.Devision)
