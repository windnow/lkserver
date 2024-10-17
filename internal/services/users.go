package services

import (
	"context"
	m "lkserver/internal/models"
	"lkserver/internal/models/types"
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

func (s *UserService) GetUserInfo(guid m.JSONByte) (*Result, error) {
	user, err := s.provider.User.Get(guid)
	if err != nil {
		return nil, m.HandleError(err, "UserService.GetUserInfo")
	}

	return &Result{
		Data: user,
		Len:  1,
		Rows: -1,
		Meta: map[string]m.META{types.Users: m.UserMETA},
	}, nil

}

func (s *UserService) UsersList(ctx context.Context, search string, limits ...int64) (*Result, error) {

	var result []*m.User
	var err error

	if search != "" {
		result, err = s.provider.User.Find(ctx, search, limits...)
	} else {
		result, err = s.provider.User.List(ctx, limits...)
	}
	if err != nil {
		return nil, m.HandleError(err, "UserService.UserList")
	}

	return &Result{
		Data: result,
		Len:  len(result),
		Rows: -1,
		Meta: map[string]m.META{types.Users: m.UserMETA},
	}, nil

}
