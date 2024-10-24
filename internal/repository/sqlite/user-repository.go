package sqlite

import (
	"context"
	"encoding/json"
	"fmt"
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
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

func (u *UserRepository) Get(guid m.JSONByte) (*m.User, error) {

	user := &m.User{
		Key: guid,
	}

	query := fmt.Sprintf("SELECT iin, name, individ from %[1]s WHERE ref = ?", types.Users)
	if err := u.source.db.QueryRow(query, user.Key).Scan(&user.Iin, &user.Name, &user.Individual); err != nil {
		return nil, m.HandleError(err, "UserRepository.Get")
	}

	return user, nil

}

func (u *UserRepository) List(ctx context.Context, limits ...int64) ([]*m.User, error) {
	limit, offset := limitations(limits)
	query := fmt.Sprintf("SELECT ref, iin, name, individ, hash FROM %[1]s LIMIT %[2]d OFFSET %[3]d", types.Users, limit, offset)

	return u.query(ctx, query)
}

func (u *UserRepository) query(ctx context.Context, query string, args ...any) ([]*m.User, error) {

	rows, err := u.source.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, m.HandleError(err, "UserRepository.query")
	}
	var result []*m.User
	for rows.Next() {
		user := &m.User{}
		err := rows.Scan(
			&user.Key, &user.Iin, &user.Name, &user.Individual, &user.PasswordHash,
		)
		if err != nil {
			return nil, m.HandleError(err, "UserRepository.query")
		}
		result = append(result, user)
	}

	return result, nil

}

func (u *UserRepository) Find(ctx context.Context, pattern string, limits ...int64) ([]*m.User, error) {
	limit, offset := limitations(limits)
	query := fmt.Sprintf("SELECT ref, iin, name, individ, hash from %[1]s WHERE name like ? OR iin like ? LIMIT %[2]d OFFSET %[3]d", types.Users, limit, offset)

	arg := "%" + pattern + "%"
	return u.query(ctx, query, arg, arg)

}

func (u *UserRepository) GetUser(iin string) (*m.User, error) {
	user := &m.User{Iin: iin}
	err := u.source.db.QueryRow(fmt.Sprintf(`
		SELECT 
			ref,
			name,
			individ,
			hash
		FROM %[1]s 
		WHERE iin=?`, types.Users),
		user.Iin).Scan(&user.Key, &user.Name, &user.Individual, &user.PasswordHash)
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
			name TEXT,
			individ BLOB,
			hash BLOB,

			FOREIGN KEY (individ) REFERENCES %[2]s(ref)
		);

		CREATE INDEX IF NOT EXISTS idx_%[1]s_iin ON %[1]s(iin);
		CREATE INDEX IF NOT EXISTS idx_%[1]s_individ ON %[1]s(individ);

	`, types.Users, types.Individuals)); err != nil {
		return m.HandleError(err, "sqliteRepo.initUserRepo")
	}

	var count int64
	userRepo.source.db.QueryRow(fmt.Sprintf(`select count(*) from %[1]s`, types.Users)).Scan(&count)
	if count == 0 {
		var users []m.User
		json.Unmarshal([]byte(mockUserData), &users)
		for _, user := range users {
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

func (u *UserRepository) Count() uint64 {
	var count uint64
	if err := u.source.db.QueryRow(fmt.Sprintf(`select count(*) from %[1]s`, types.Users)).Scan(&count); err != nil {
		return 0
	}

	return count

}

func (u *UserRepository) Save(ctx context.Context, user *m.User) error {
	if user.Key.Blank() {
		Key, err := m.GenerateUUID()
		if err != nil {
			return m.HandleError(err, "UserRepository.Save")
		}
		user.Key = Key
	}

	return m.HandleError(u.source.ExecContextInTransaction(ctx, insertUserQuery, nil,
		user.Key,
		user.Individual,
		user.Iin,
		user.Name,
		user.PasswordHash[:]))

}

var insertUserQuery string = fmt.Sprintf(`INSERT OR REPLACE INTO %[1]s(ref, individ, iin, name, hash) VALUES (?, ?, ?, ?, ?)`, types.Users)
var mockUserData string = `[
    {
		"key": "c9aba8d6-351a-4d85-a8b6-9427ea2f8c8e",
        "iin": "821019000888",
		"name": "Усенбаев Дархан Жаксылыкович",
		"individ":"27f74b66-cba7-486d-a263-81b6cb9a3e57",
        "pin": "82"
    },
    {
		"key": "8c272f7c-6c2c-4dba-bba5-4062005b2400",
        "iin": "851204000888",
		"name": "Каримов Кайрат Ганиевич",
		"individ":"52efc72d-ba0d-4f87-ae73-e902936395fe",
        "pin": "85"
    },
    {
		"key": "f31c6a0f-b07c-4632-8949-2f24fde4fc26",
        "iin": "910702000888",
		"name": "Асетов Алинур Дарханулы",
		"individ":"19db2753-68f9-4b5d-998a-727e347a958a",
        "pin": "91"
    }
]`
