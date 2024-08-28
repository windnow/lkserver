package repository

import "lkserver/internal/models"

type Repo struct {
	User     UserProvider
	Contract ContractProvider
}

type UserProvider interface {
	FindUser(iin, pin string) (*models.User, error)
	GetUser(iin string) (*models.User, error)
	Close()
}

func (r *Repo) Close() {
	r.User.Close()
	r.Contract.Close()
}

type ContractProvider interface {
	Close()
}
