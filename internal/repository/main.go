package repository

import (
	"context"
	"database/sql"
	"io"
	"lkserver/internal/models"
	"lkserver/internal/models/catalogs"
	"lkserver/internal/models/reports"
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
	Reports              ReportProvider
	Catalogs             *Catalogs
}

type Catalogs struct {
	Cato         CatoProvider
	Vus          VusProvider
	Organization OrganizationProvider
	Devision     DevisionProvider
	OrderSource  OrderSourceProvider
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
	Find(ctx context.Context, pattern string, limits ...int64) ([]*models.User, error)
	List(ctx context.Context, limits ...int64) ([]*models.User, error)
	FindUser(iin, pin string) (*models.User, error)
	GetUser(iin string) (*models.User, error)
	Get(guid models.JSONByte) (*models.User, error)
	Save(ctx context.Context, user *models.User) error
	Close()
}

type IndividualsProvider interface {
	Get(key models.JSONByte) (*models.Individuals, error)
	GetByIin(iin string) (*models.Individuals, error)
}

type ContractProvider interface {
	Close()
}

type FileProvider interface {
	// Return reader of file (io.Reader)
	GetFile(fileId string) (io.Reader, string, error)
}

type RankProvider interface {
	Get(key models.JSONByte) (*models.Rank, error)
	Save(ctx context.Context, rank *models.Rank) error
	Close()
}

type RankHistoryProvider interface {
	GetLastByIin(individIin string) (*models.RankHistory, error)
	GetHistoryByIin(indivIin string) ([]*models.RankHistory, error)
	Close()
}

type EducationInstitutionProvider interface {
	Get(Key models.JSONByte) (*models.EducationInstitution, error)
	Save(ctx context.Context, ei *models.EducationInstitution) error
	Close()
}

type SpecialtiesProvider interface {
	Get(key models.JSONByte) (*models.Specialties, error)
	Save(ctx context.Context, ei *models.Specialties) error
	Close()
}

type EducationProvider interface {
	GetByIin(individIin string) ([]*models.Education, error)
	Save(ctx context.Context, ei *models.Education) error
	Close()
}

type ReportProvider interface {
	GetTransaction(ctx context.Context) (*sql.Tx, error)
	GetTypes(codes []string) ([]*reports.ReportTypes, error)
	SaveType(context.Context, *reports.ReportTypes) error
	Save(tx *sql.Tx, ctx context.Context, report *models.Report) error
	SaveCoordinators(tx *sql.Tx, ctx context.Context, coordinators []*reports.Coordinators) error
	SaveDetails(tx *sql.Tx, ctx context.Context, report *models.Report, data any) error
	GetStructure(reportType string) (any, error)
	Get(guid models.JSONByte) (*models.Report, error)
	GetCoordinators(ctx context.Context, report *models.Report) ([]*reports.Coordinators, error)
	GetDetails(ctx context.Context, report *models.Report) (any, models.META, error)
	List(context.Context, models.JSONByte) ([]*models.Report, error)
}
type ReportDetails interface {
	Get(ctx context.Context, ref models.JSONByte, tx ...*sql.Tx) (any, models.META, error)
	Save(tx *sql.Tx, ctx context.Context, report models.JSONByte, data any) error
	Init() error
	GetStructure() interface{}
}

type CatoProvider interface {
	Get(ctx context.Context, Ref models.JSONByte) (*catalogs.Cato, error)
	List(ctx context.Context, parentRef models.JSONByte, limits ...int64) ([]*catalogs.Cato, error)
	Find(ctx context.Context, description string, limits ...int64) ([]*catalogs.Cato, error)
	Count(ctx context.Context) int64
}

type VusProvider interface {
	Get(ctx context.Context, Ref models.JSONByte) (*catalogs.Vus, error)
	List(ctx context.Context, limits ...int64) ([]*catalogs.Vus, error)
	Find(ctx context.Context, pattern string, limits ...int64) ([]*catalogs.Vus, error)
	Save(ctx context.Context, vus *catalogs.Vus, tx *sql.Tx) error
	Count(ctx context.Context) int64
}

type OrganizationProvider interface {
	Get(ctx context.Context, Ref models.JSONByte) (*catalogs.Organization, error)
	List(ctx context.Context, limits ...int64) ([]*catalogs.Organization, error)
	Find(ctx context.Context, pattern string, limits ...int64) ([]*catalogs.Organization, error)
	Save(ctx context.Context, org *catalogs.Organization, tx *sql.Tx) error
	Count(ctx context.Context) int64
}

type DevisionProvider interface {
	Get(ctx context.Context, Ref models.JSONByte) (*catalogs.Devision, error)
	List(ctx context.Context, limits ...int64) ([]*catalogs.Devision, error)
	Find(ctx context.Context, pattern string, limits ...int64) ([]*catalogs.Devision, error)
	Save(ctx context.Context, dev *catalogs.Devision, tx *sql.Tx) error
	Count(ctx context.Context) int64
}

type OrderSourceProvider interface {
	Get(ctx context.Context, Ref models.JSONByte) (*catalogs.OrderSource, error)
	List(ctx context.Context, limits ...int64) ([]*catalogs.OrderSource, error)
	Find(ctx context.Context, pattern string, limits ...int64) ([]*catalogs.OrderSource, error)
	Save(ctx context.Context, dev *catalogs.OrderSource, tx *sql.Tx) error
	Count(ctx context.Context) int64
}
