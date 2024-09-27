package sqlite

import (
	"context"
	"database/sql"

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
	return &src{
		db: db,
	}, nil

}

func (s *src) ExecContextInTransaction(ctx context.Context, query string, args ...any) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		query,
		args...)
	if err != nil {
		tx.Rollback()
		return err
	}
	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:
		return tx.Commit()
	}
}

func (s *src) Exec(query string, args ...any) error {
	_, err := s.db.Exec(query, args...)
	return err
}
