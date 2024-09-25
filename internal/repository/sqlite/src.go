package sqlite

import (
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
