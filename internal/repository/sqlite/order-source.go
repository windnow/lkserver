package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"lkserver/internal/models"
	"lkserver/internal/models/catalogs"
	"lkserver/internal/models/types"
)

type orderSource struct {
	source *src
}

func (s *sqliteRepo) initOrderSource() error {
	src := &orderSource{
		source: s.db,
	}

	if _, err := src.source.db.Exec(createOrderSourceQuery); err != nil {
		return models.HandleError(err, "sqliteRepo.initOrderSource")
	}

	s.catalogs.OrderSource = src
	if src.Count(context.Background()) == 0 {
		if err := src.loadData(); err != nil {
			return models.HandleError(err, "sqliteRepo.initOrderSource")
		}
	}

	s.catalogs.OrderSource = src

	return nil
}

func (o *orderSource) loadData() error {
	var mockOrderSourceData []*catalogs.OrderSource
	if err := json.Unmarshal([]byte(mockOrderSources), &mockOrderSourceData); err != nil {
		return models.HandleError(err, "orderSource.loadData")
	}
	tx, err := o.source.db.Begin()
	if err != nil {
		return models.HandleError(err, "orderSource.loadData")
	}
	for _, row := range mockOrderSourceData {
		o.Save(context.Background(), row, tx)
	}
	return nil
}

func newOrderSource() models.Scanable {
	return &catalogs.OrderSource{}
}

func (o *orderSource) Get(ctx context.Context, Ref models.JSONByte) (*catalogs.OrderSource, error) {
	query := fmt.Sprintf(`SELECT ref, num, description FROM %[1]s WHERE ref = ?`, types.OrderSource)
	rows, err := models.Query[*catalogs.OrderSource](o.source.db, ctx, newOrderSource, query, Ref)
	if err != nil {
		return nil, models.HandleError(err, "orderSource.Get")
	}
	if len(rows) != 1 {
		return nil, models.HandleError(models.ErrNotFound, "orderSource.Get")
	}
	return rows[0], nil
}

func (o *orderSource) List(ctx context.Context, limits ...int64) ([]*catalogs.OrderSource, error) {
	limit, offset := limitations(limits)
	query := fmt.Sprintf(`SELECT ref, num, description FROM %[1]s LIMIT %[2]d OFFSET %[3]d`, types.OrderSource, limit, offset)
	rows, err := models.Query[*catalogs.OrderSource](o.source.db, ctx, newOrderSource, query)
	if err != nil {
		return nil, models.HandleError(err, "orderSource.List")
	}
	return rows, nil
}

func (o *orderSource) Find(ctx context.Context, pattern string, limits ...int64) ([]*catalogs.OrderSource, error) {

	limit, offset := limitations(limits)
	query := fmt.Sprintf(`SELECT ref, num, description FROM %[1]s WHERE description LIKE ? LIMIT %[2]d OFFSET %[3]d`, types.OrderSource, limit, offset)
	rows, err := models.Query[*catalogs.OrderSource](o.source.db, ctx, newOrderSource, query, "%"+pattern+"%")
	if err != nil {
		return nil, models.HandleError(err, "orderSource.Find")
	}
	return rows, nil
}

func (o *orderSource) Save(ctx context.Context, orderSource *catalogs.OrderSource, tx *sql.Tx) error {
	query := fmt.Sprintf(`INSERT OR REPLACE INTO %[1]s (ref, num, description) VALUES (?, ?, ?)`, types.OrderSource)
	return models.HandleError(o.source.ExecContextInTransaction(ctx, query, nil, orderSource.Ref, orderSource.Number, orderSource.Description))
}

func (o *orderSource) Count(ctx context.Context) int64 {
	var count int64
	o.source.db.QueryRow(fmt.Sprintf(`select count(ref) from %[1]s`, types.OrderSource)).Scan(&count)
	return count
}

var createOrderSourceQuery = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %[1]s (
    ref         BLOB PRIMARY KEY,
    num         INTEGER UNIQUE,     
    description TEXT
)`, types.OrderSource)

var mockOrderSources = `[
{"ref": "cba50c1c-4b96-4ada-a536-5b0aaad01cf9","num": 1, "Description": "Қазақстан Республикасы Қорғаныс министрінің"},
{"ref": "1f960ff5-3813-4ef3-9e18-6cebe0b0f8fd","num": 2, "Description": "Қазақстан Республикасы Қорғаныс министрінің бірінші орынбасары"},
{"ref": "83270ad3-f67b-4225-9156-849ecc5f6558","num": 3, "Description": "Қазақстан Республикасы Қарулы Күштері Бас штабы бастығының"},
{"ref": "3d23400e-0bac-4cb6-bcb9-92d1987db2b5","num": 4, "Description": "Қазақстан Республикасы Қорғаныс министрлігі армиясының орталық спорт клубының бастығының"},
{"ref": "8d4c43a0-f699-4642-b42f-05c3bd070ac0","num": 5, "Description": "Қазақстан Республикасы Ұлттық Қорғаныс университеті бастығының"},
{"ref": "86299282-a7bc-402d-9d7d-a20a42db3d18","num": 6, "Description": "Басқа / Другое"}
]`
