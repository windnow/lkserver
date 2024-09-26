package sqlite

import (
	"context"
	"lkserver/internal/models"
	"lkserver/internal/repository"
)

type UserRepository struct {
	source *src
}

func (u *UserRepository) FindUser(iin, pin string) (*models.User, error) {

	user, err := u.GetUser(iin)
	if err != nil {
		return nil, err
	}
	user.Pin = pin
	if !user.CheckPassword() {
		return nil, repository.ErrInvalidCredentials
	}
	return user, nil

}
func (u *UserRepository) GetUser(iin string) (*models.User, error) {
	user := &models.User{Iin: iin}
	err := u.source.db.QueryRow(`
		SELECT 
			guid,
			hash
		FROM users 
		WHERE iin=?`,
		user.Iin).Scan(&user.Key, &user.PasswordHash)
	if err != nil {
		return nil, handleQueryError(err)
	}
	return user, nil
}

func (u *UserRepository) Close() {
	//nop
}

func (r *sqliteRepo) initUserRepo() error {
	userRepo := &UserRepository{
		source: r.db,
	}

	if err := userRepo.source.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			guid BLOB PRIMARY KEY,
			iin TEXT UNIQUE,
			hash BLOB
		)
	`); err != nil {
		return err
	}
	if err := userRepo.source.Exec(`
		CREATE INDEX IF NOT EXISTS idx_individuals_iin ON individuals(iin);
		`); err != nil {
		return err
	}

	var count int64
	userRepo.source.db.QueryRow(`select count(*) from users`).Scan(&count)
	if count == 0 {
		id1, _ := GenerateUUID()
		id2, _ := GenerateUUID()
		id3, _ := GenerateUUID()
		users := []models.User{
			{Key: id1, Iin: "821019000888", Pin: "82"},
			{Key: id2, Iin: "851204000888", Pin: "85"},
			{Key: id3, Iin: "910702000888", Pin: "91"},
		}
		for _, user := range users {
			if err := user.GeneratePasswordHash(); err != nil {
				return err
			}
			if err := userRepo.Save(context.Background(), &user); err != nil {
				return err
			}
		}
	}
	userRepo.FindUser("821019000888", "82")

	r.userRepo = userRepo
	return nil
}

func (u *UserRepository) Save(ctx context.Context, user *models.User) error {

	return u.source.ExecContextInTransaction(ctx, insertUserQuery,
		user.Key,
		user.Iin,
		user.PasswordHash[:])

}

var insertUserQuery string = `INSERT OR REPLACE INTO users(guid, iin, hash) VALUES (?, ?, ?)`
