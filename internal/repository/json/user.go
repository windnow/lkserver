package json

import (
	"errors"
	"fmt"
	"lkserver/internal/models"
	"time"
)

type User struct {
	Name      string `json:"name"`
	Iin       string `json:"iin"`
	Pin       string `json:"pin"`
	BirthDate string `json:"birth_date"`
	Image     string `json:"image"`
}
type UserRepo struct {
	dataDir string
	users   []User
}

func NewUserRepo(dataDir string) (*UserRepo, error) {
	repo := &UserRepo{dataDir: dataDir}
	if err := repo.init(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (u *User) compile() (*models.User, error) {
	BirthDate, err := time.Parse("2006-01-02", u.BirthDate)
	if err != nil {
		return nil, err
	}
	return &models.User{
		Name:      u.Name,
		Iin:       u.Iin,
		Pin:       u.Pin,
		BirthDate: models.JSONTime(BirthDate),
		Image:     u.Image,
	}, nil
}

func (r *UserRepo) init() error {

	return initFile(fmt.Sprintf("%s/users.json", r.dataDir), &r.users)

}

func (r *UserRepo) GetUser(iin string) (*models.User, error) {

	for _, user := range r.users {
		if user.Iin == iin {
			return user.compile()
		}
	}
	return nil, errors.New("NOT FOUND")
}

func (r *UserRepo) FindUser(iin, pin string) (*models.User, error) {

	for _, user := range r.users {
		if user.Iin == iin && user.Pin == pin {
			return user.compile()
		}
	}
	return nil, errors.New("NOT FOUND")

}

func (r *UserRepo) Close() {}
