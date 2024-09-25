package sqlite

import "lkserver/internal/models"

type UserRepository struct {
	source *src
}

func (u *UserRepository) FindUser(iin, pin string) (*models.User, error) {
	return nil, nil
}
func (u *UserRepository) GetUser(iin string) (*models.User, error) {
	return nil, nil
}

func (u *UserRepository) Close() {
	//nop
}

func (r *sqliteRepo) initUserRepo() error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			guid BLOB PRIMARY KEY,
			iin TEXT,
			hash TEXT
		)
	`

	userRepo := &UserRepository{
		source: r.db,
	}
	_, err := userRepo.source.db.Exec(query)
	if err != nil {
		return err
	}

	var count int64
	userRepo.source.db.QueryRow(`select count(*) from users`).Scan(&count)
	if count == 0 {

	}

	r.userRepo = userRepo
	return nil
}
