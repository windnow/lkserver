package sqlite

import (
	"errors"
	"lkserver/internal/repository"
)

func NewSQLiteProvider(dataFile string) (*repository.Repo, error) {
	db, err := newDB(dataFile)
	if err != nil {
		return nil, err
	}
	repo := &sqliteRepo{
		db: db,
	}

	if err := repo.initUserRepo(); err != nil {
		return nil, err
	}

	return &repository.Repo{
		User: repo.userRepo,
	}, errors.New("Testing")
}
