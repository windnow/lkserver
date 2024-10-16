package sqlite

import (
	"context"
	"database/sql"

	m "lkserver/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type src struct {
	db *sql.DB
}

func newDB(file string) (*src, error) {

	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}
	return &src{
		db: db,
	}, nil
}

func (s *src) ExecContextInTransaction(ctx context.Context, query string, tx *sql.Tx, args ...any) error {
	var err error
	localTx := false
	if tx == nil {
		tx, err = s.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		localTx = true
	}
	defer func() {
		if localTx {
			tx.Rollback()
		}
	}()
	_, err = tx.ExecContext(ctx,
		query,
		args...)
	if err != nil {
		if localTx {
			tx.Rollback()
		}
		return m.HandleError(err, "src.ExecContextInTransaction")
	}
	if localTx {
		select {
		case <-ctx.Done():
			tx.Rollback()
			return m.HandleError(ctx.Err(), "src.ExecContextInTransaction")
		default:
			return m.HandleError(tx.Commit(), "src.ExecContextInTransaction")
		}
	}
	return nil
}

func (s *src) Exec(query string, args ...any) error {
	_, err := s.db.Exec(query, args...)
	return err
}
