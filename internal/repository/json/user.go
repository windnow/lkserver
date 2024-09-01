package json

import (
	"errors"
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
type userRepo struct {
	users []User
}

func (r *repo) initUserRepo() error {
	repo := &userRepo{}
	if err := initFile(r.dataDir+"/users.json", &repo.users); err != nil {
		return err
	}
	r.user = repo
	return nil
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

func (r *userRepo) GetUser(iin string) (*models.User, error) {

	for _, user := range r.users {
		if user.Iin == iin {
			return user.compile()
		}
	}
	return nil, errors.New("NOT FOUND")
}

func (r *userRepo) FindUser(iin, pin string) (*models.User, error) {

	for _, user := range r.users {
		if user.Iin == iin && user.Pin == pin {
			return user.compile()
		}
	}
	return nil, errors.New("NOT FOUND")

}

func (r *userRepo) Close() {}
