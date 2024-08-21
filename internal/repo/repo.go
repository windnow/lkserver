package repo

import "lkserver/internal/models"

type DataProvider interface {
	FindUser(iin, pin string) (*models.User, error)
	GetUser(iin string) (*models.User, error)
	Close()
}
