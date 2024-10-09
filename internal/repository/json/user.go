package json

import (
	"context"
	"errors"
	"lkserver/internal/models"
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
	return &models.User{
		Iin: u.Iin,
		Pin: u.Pin,
	}, nil
}

func (u *userRepo) Get(guid models.JSONByte) (*models.User, error) {

	return nil, errors.ErrUnsupported

}
func (r *userRepo) GetUser(iin string) (*models.User, error) {

	for _, user := range r.users {
		if user.Iin == iin {
			return user.compile()
		}
	}
	return nil, errors.New("NOT FOUND")
}

func (s *userRepo) Save(ctx context.Context, user *models.User) error {
	return errors.ErrUnsupported
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
