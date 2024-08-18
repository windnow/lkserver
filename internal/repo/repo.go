package repo

import "lkserver/internal/models"

type DataProvider interface {
	GetUser(iin, pin string) (*models.User, error)
	Close()
}
