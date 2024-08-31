package repository

import (
	"errors"
	"io"
	"lkserver/internal/models"
)

var (
	ErrNotFound     = errors.New("NOT FOUND")
	ErrRefIntegrity = errors.New("REFERENCE INTEGRITY IS VIOLATED")
)

type Repo struct {
	User         UserProvider
	Individuals  IndividualsProvider
	Contract     ContractProvider
	Ranks        RankProvider
	RanksHistory RankHistoryProvider
}

func (r *Repo) Close() {
	r.User.Close()
	r.Contract.Close()
	r.Ranks.Close()
	r.RanksHistory.Close()
}

type UserProvider interface {
	FindUser(iin, pin string) (*models.User, error)
	GetUser(iin string) (*models.User, error)
	Close()
}

type IndividualsProvider interface {
	Get(iin string) (*models.Individuals, error)
}

type ContractProvider interface {
	Close()
}

type FileProvider interface {
	// Return reader of file (io.Reader)
	GetFile(fileId string) (io.Reader, string, error)
}

type RankProvider interface {
	Get(id int) (*models.Rank, error)
	Close()
}

type RankHistoryProvider interface {
	GetLast(people *models.Individuals) (*models.RankHistory, error)
	GetHistory(people *models.Individuals) ([]*models.RankHistory, error)
	Close()
}
