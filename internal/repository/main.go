package repository

import (
	"errors"
	"io"
	"lkserver/internal/models"
)

var (
	ErrNotFound           = errors.New("NOT FOUND")
	ErrRefIntegrity       = errors.New("REFERENCE INTEGRITY IS VIOLATED")
	ErrInvalidCredentials = errors.New("INVALID CREDENTIALS")
)

type Repo struct {
	User                 UserProvider
	Individuals          IndividualsProvider
	Contract             ContractProvider
	Ranks                RankProvider
	RanksHistory         RankHistoryProvider
	EducationInstitution EducationInstitutionProvider
	Specialties          SpecialtiesProvider
	Education            EducationProvider
}

func (r *Repo) Close() {
	r.User.Close()
	r.Contract.Close()
	r.Ranks.Close()
	r.RanksHistory.Close()
	r.EducationInstitution.Close()
	r.Specialties.Close()
	r.Education.Close()
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
	GetLast(individIin string) (*models.RankHistory, error)
	GetHistory(indivIin string) ([]*models.RankHistory, error)
	Close()
}

type EducationInstitutionProvider interface {
	Get(id int) (*models.EducationInstitution, error)
	Close()
}

type SpecialtiesProvider interface {
	Get(id int) (*models.Specialties, error)
	Close()
}

type EducationProvider interface {
	Get(individIin string) ([]*models.Education, error)
	Close()
}
