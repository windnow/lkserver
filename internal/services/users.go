package services

import (
	"fmt"
	m "lkserver/internal/models"
	"lkserver/internal/repository"
)

type UserService struct {
	provider *repository.Repo
}

func NewUsersService(r *repository.Repo) *UserService {
	return &UserService{
		provider: r,
	}
}

func (s *UserService) GetUserInfo(guid m.JSONByte) (any, error) {
	type result struct {
		Key        m.JSONByte `json:"key"`
		IndividKey m.JSONByte `json:"individKey"`
		Iin        string     `json:"iin"`
		Login      string     `json:"login"`
		Name       string     `json:"name"`
	}
	user, err := s.provider.User.Get(guid)
	if err != nil {
		return nil, m.HandleError(err, "UserService.GetUserInfo")
	}

	individ, err := s.provider.Individuals.Get(user.Individual)
	if err != nil {
		return nil, m.HandleError(err, "UserService.GetUserInfo")
	}

	return &result{
		Key:        user.Key,
		Iin:        individ.IndividualNumber,
		IndividKey: individ.Key,
		Login:      user.Iin,
		Name:       fmt.Sprintf("%s %s %s", individ.LastName, individ.FirstName, individ.Patronymic),
	}, nil

}
