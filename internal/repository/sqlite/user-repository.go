package sqlite

import (
	"context"
	"database/sql"
	"errors"
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) Close() {
	//nop
}

func (r *sqliteRepo) initUserRepo() error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			guid BLOB PRIMARY KEY,
			iin TEXT UNIQUE,
			hash BLOB
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
	tx, err := u.source.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		`INSERT OR REPLACE INTO users(guid, iin, hash) VALUES (?, ?, ?)`,
		user.Key, user.Iin, user.PasswordHash[:])
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
