package sqlite

import (
	"context"
	"encoding/json"
	. "lkserver/internal/models"
)

type UserRepository struct {
	source *src
}

func (u *UserRepository) FindUser(iin, pin string) (*User, error) {

	user, err := u.GetUser(iin)
	if err != nil {
		return nil, HandleError(err, "UserRepository.FindUser")
	}
	user.Pin = pin
	if !user.CheckPassword() {
		return nil, HandleError(ErrInvalidCredentials, "UserRepository.FindUser")
	}
	return user, nil

}
func (u *UserRepository) GetUser(iin string) (*User, error) {
	user := &User{Iin: iin}
	err := u.source.db.QueryRow(`
		SELECT 
			guid,
			hash
		FROM users 
		WHERE iin=?`,
		user.Iin).Scan(&user.Key, &user.PasswordHash)
	if err != nil {
		return nil, HandleError(err, "UserRepository.GetUser")
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
		return HandleError(err, "sqliteRepo.initUserRepo")
	}
	if err := userRepo.source.Exec(`
		CREATE INDEX IF NOT EXISTS idx_users_iin ON users(iin);
		`); err != nil {
		return HandleError(err, "sqliteRepo.initUserRepo")
	}

	var count int64
	userRepo.source.db.QueryRow(`select count(*) from users`).Scan(&count)
	if count == 0 {
		var users []User
		json.Unmarshal([]byte(mockUserData), &users)
		for _, user := range users {
			if err := user.GeneratePasswordHash(); err != nil {
				return HandleError(err, "sqliteRepo.initUserRepo")
			}
			if err := userRepo.Save(context.Background(), &user); err != nil {
				return HandleError(err, "sqliteRepo.initUserRepo")
			}
		}
	}
	userRepo.FindUser("821019000888", "82")

	r.userRepo = userRepo
	return nil
}

func (u *UserRepository) Save(ctx context.Context, user *User) error {

	return u.source.ExecContextInTransaction(ctx, insertUserQuery,
		user.Key[:],
		user.Iin,
		user.PasswordHash[:])

}

var insertUserQuery string = `INSERT OR REPLACE INTO users(guid, iin, hash) VALUES (?, ?, ?)`
var mockUserData string = `[
    {
		"key": "c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e",
        "iin": "821019000888",
        "pin": "82"
    },
    {
		"key": "8c272f7c-6c2c-4dba-bba5-4062005b2400",
        "iin": "851204000888",
        "pin": "85"
    },
    {
		"key": "f31c6a0f-b07c-4632-8949-2f24fde4fc26",
        "iin": "910702000888",
        "pin": "91"
    }
]`
