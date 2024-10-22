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

type organization struct {
	source *src
}

func (s *sqliteRepo) initOrganization() error {
	src := &organization{
		source: s.db,
	}

	if err := src.source.Exec(createOrganizationsQuery); err != nil {
		return m.HandleError(err, "sqliteRepo.initOrganization")
	}

	if src.Count(context.Background()) == 0 {
		if err := src.loadData("data/orgs.json"); err != nil {
			return m.HandleError(err, "sqliteRepo.initOrganization")
		}
	}
	s.catalogs.Organization = src

	return nil
}

func (o *organization) loadData(path string) error {

	file, err := os.Open(path)
	if err != nil {
		return m.HandleError(err, "organization.loadData")
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return m.HandleError(err, "organization.loadData")
	}
	data := []*catalogs.Organization{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return m.HandleError(err, "organization.loadData")
	}

	tx, err := o.source.db.Begin()
	if err != nil {
		return m.HandleError(err, "organization.loadData")
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
			return m.HandleError(err, "organization.loadData")
		}
	}
	return m.HandleError(tx.Commit(), "organization.loadData")
}

func (o *organization) saveData(tx *sql.Tx, data []*catalogs.Organization) error {
	query := fmt.Sprintf("INSERT INTO %[1]s (ref, code, description, title, short_title, conventional_title) VALUES (?, ?, ?, ?, ?, ?)", types.Organization)

	stmt, err := tx.Prepare(query)
	if err != nil {
		return m.HandleError(err, "organization.saveData")
	}
	defer stmt.Close()

	for _, cc := range data {
		_, err := stmt.Exec(
			cc.Ref,
			strings.TrimSpace(cc.Code),
			strings.TrimSpace(cc.Description),
			strings.TrimSpace(cc.Title),
			strings.TrimSpace(cc.ShortTitle),
			strings.TrimSpace(cc.ConventionalTitle),
		)
		if err != nil {
			return m.HandleError(err, "organization.saveData")
		}
	}

	return nil
}

func newOrganization() m.Scanable {
	return &catalogs.Organization{}
}

func (o *organization) Get(ctx context.Context, Ref m.JSONByte) (*catalogs.Organization, error) {

	query := fmt.Sprintf(`SELECT ref, code, description, title, short_title, conventional_title FROM %[1]s WHERE ref = ?`, types.Organization)
	rows, err := m.Query[*catalogs.Organization](o.source.db, ctx, newOrganization, nil, query, Ref)
	if err != nil {
		return nil, m.HandleError(err, "organization.Get")
	}
	if len(rows) == 0 {
		return nil, m.HandleError(m.ErrNotFound, "organization.Get")
	}

	return rows[0], nil

}

func (o *organization) List(ctx context.Context, limits ...int64) ([]*catalogs.Organization, error) {
	limit, offset := limitations(limits)
	query := fmt.Sprintf(`SELECT ref, code, description, title, short_title, conventional_title FROM %[1]s LIMIT %[2]d OFFSET %[3]d`, types.Organization, limit, offset)
	rows, err := m.Query[*catalogs.Organization](o.source.db, ctx, newOrganization, nil, query)
	if err != nil {
		return nil, m.HandleError(err, "organization.List")
	}
	return rows, nil
}

func (o *organization) Find(ctx context.Context, pattern string, limits ...int64) ([]*catalogs.Organization, error) {
	limit, offset := limitations(limits)
	query := fmt.Sprintf(`SELECT ref, code, description, title, short_title, conventional_title FROM %[1]s
    WHERE code LIKE ? OR description LIKE ? OR title LIKE ? OR short_title LIKE ? OR conventional_title LIKE ?
	LIMIT %[2]d OFFSET %[3]d`, types.Organization, limit, offset)
	param := "%" + pattern + "%"
	rows, err := m.Query[*catalogs.Organization](o.source.db, ctx, newOrganization, nil, query, param, param, param, param, param)
	if err != nil {
		return nil, m.HandleError(err, "organization.Find")
	}
	return rows, nil
}

func (o *organization) Save(ctx context.Context, org *catalogs.Organization, tx *sql.Tx) error {
	return m.HandleError(o.saveData(tx, []*catalogs.Organization{org}))
}

func (o *organization) Count(ctx context.Context) int64 {
	var count int64
	o.source.db.QueryRow(fmt.Sprintf(`select count(ref) from %[1]s`, types.Organization)).Scan(&count)
	return count
}

var createOrganizationsQuery = fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %[1]s (
    ref                BLOB PRIMARY KEY,
    code               TEXT,
    description        TEXT,
    title              TEXT,
    short_title        TEXT,
    conventional_title TEXT
);

CREATE INDEX IF NOT EXISTS idx_%[1]s_code ON %[1]s(code);
CREATE INDEX IF NOT EXISTS idx_%[1]s_description ON %[1]s(description);
CREATE INDEX IF NOT EXISTS idx_%[1]s_title ON %[1]s(title);
CREATE INDEX IF NOT EXISTS idx_%[1]s_short_title ON %[1]s(short_title);
CREATE INDEX IF NOT EXISTS idx_%[1]s_conventional_title ON %[1]s(conventional_title);
`, types.Organization)
