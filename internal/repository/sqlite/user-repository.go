package sqlite

import (
	"context"
	"encoding/json"
	"fmt"
	m "lkserver/internal/models"
)

type UserRepository struct {
	source *src
	repo   *sqliteRepo
}

func (u *UserRepository) FindUser(iin, pin string) (*m.User, error) {

	user, err := u.GetUser(iin)
	if err != nil {
		return nil, m.HandleError(err, "UserRepository.FindUser")
	}
	user.Pin = pin
	if !user.CheckPassword() {
		return nil, m.HandleError(m.ErrInvalidCredentials, "UserRepository.FindUser")
	}
	return user, nil

}
func (u *UserRepository) GetUser(iin string) (*m.User, error) {
	user := &m.User{Iin: iin}
	err := u.source.db.QueryRow(fmt.Sprintf(`
		SELECT 
			ref,
			hash
		FROM %[1]s 
		WHERE iin=?`, tabUsers),
		user.Iin).Scan(&user.Key, &user.PasswordHash)
	if err != nil {
		return nil, m.HandleError(err, "UserRepository.GetUser")
	}
	return user, nil
}

func (u *UserRepository) Close() {
	//nop
}

func (r *sqliteRepo) initUserRepo() error {
	userRepo := &UserRepository{
		source: r.db,
		repo:   r,
	}

	if err := userRepo.source.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %[1]s (
			ref BLOB PRIMARY KEY,
			iin TEXT UNIQUE,
			individ BLOB,
			hash BLOB,

			FOREIGN KEY (individ) REFERENCES %[2]s(ref)
		);

		CREATE INDEX IF NOT EXISTS idx_%[1]s_iin ON %[1]s(iin);
		CREATE INDEX IF NOT EXISTS idx_%[1]s_individ ON %[1]s(individ);

	`, tabUsers, tabIndividuals)); err != nil {
		return m.HandleError(err, "sqliteRepo.initUserRepo")
	}

	type mockUser struct {
		Key        m.JSONByte `json:"key"`
		Iin        string     `json:"iin"`
		Pin        string     `json:"pin,omitempty"`
		Individual m.JSONByte `json:"individ"`
	}

	var count int64
	userRepo.source.db.QueryRow(fmt.Sprintf(`select count(*) from %[1]s`, tabUsers)).Scan(&count)
	if count == 0 {
		var mock_users []mockUser
		json.Unmarshal([]byte(mockUserData), &mock_users)
		for _, mock_user := range mock_users {
			Individ, err := userRepo.repo.individuals.Get(mock_user.Individual)
			if err != nil {
				return m.HandleError(err, "sqliteRepo.initUserRepo")
			}
			user := m.User{
				Key:        mock_user.Key,
				Iin:        mock_user.Iin,
				Pin:        mock_user.Pin,
				Individual: Individ,
			}
			if err := user.GeneratePasswordHash(); err != nil {
				return m.HandleError(err, "sqliteRepo.initUserRepo")
			}
			if err := userRepo.Save(context.Background(), &user); err != nil {
				return m.HandleError(err, "sqliteRepo.initUserRepo")
			}
		}
		user := m.User{
			Iin: "830119399019",
			Pin: "83",
		}
		err := user.GeneratePasswordHash()
		if err != nil {
			return m.HandleError(err, "sqliteRepo.initUserRepo")
		}
		userRepo.Save(context.Background(), &user)
	}
	_, err := userRepo.FindUser("821019000888", "82")
	if err != nil {
		return m.HandleError(err, "sqliteRepo.initUserRepo")
	}

	r.userRepo = userRepo
	return nil
}

func (u *UserRepository) Save(ctx context.Context, user *m.User) error {
	if user.Key.Blank() {
		Key, err := m.GenerateUUID()
		if err != nil {
			return m.HandleError(err, "UserRepository.Save")
		}
		user.Key = Key
	}
	var key m.JSONByte
	if user.Individual != nil {
		key = user.Individual.Key
	}

	return m.HandleError(u.source.ExecContextInTransaction(ctx, insertUserQuery,
		user.Key,
		key,
		user.Iin,
		user.PasswordHash[:]))

}

var insertUserQuery string = fmt.Sprintf(`INSERT OR REPLACE INTO %[1]s(ref, individ, iin, hash) VALUES (?, ?, ?, ?)`, tabUsers)
var mockUserData string = `[
    {
		"key": "c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e",
        "iin": "821019000888",
		"individ":"27f74b66-cba7-486d-a263-81b6cb9a3e57",
        "pin": "82"
    },
    {
		"key": "8c272f7c-6c2c-4dba-bba5-4062005b2400",
        "iin": "851204000888",
		"individ":"52efc72d-ba0d-4f87-ae73-e902936395fe",
        "pin": "85"
    },
    {
		"key": "f31c6a0f-b07c-4632-8949-2f24fde4fc26",
        "iin": "910702000888",
		"individ":"19db2753-68f9-4b5d-998a-727e347a958a",
        "pin": "91"
    }
]`
