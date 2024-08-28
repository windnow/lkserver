package repository

import (
	"io"
	"lkserver/internal/models"
)

type Repo struct {
	User     UserProvider
	Contract ContractProvider
	Files    FileProvider
}

func (r *Repo) Close() {
	r.User.Close()
	r.Contract.Close()
}

type UserProvider interface {
	FindUser(iin, pin string) (*models.User, error)
	GetUser(iin string) (*models.User, error)
	Close()
}

type ContractProvider interface {
	Close()
}

type FileProvider interface {
	// Return reader of file (io.Reader)
	GetFile(fileId string) (io.Reader, string, error)
}
