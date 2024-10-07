package services

import (
	"context"
	"errors"
	"fmt"
	"lkserver/internal/models"
	"lkserver/internal/models/reports"
	"lkserver/internal/repository"
	"time"
)

type ReportProcessor interface {
	Check(any) error
}

type processorsFactory struct {
	processorsMap map[string]ReportProcessor
}

type ReportService struct {
	reports    repository.ReportProvider
	processors processorsFactory
}

func (f *processorsFactory) GetProcessor(reportType string) (ReportProcessor, error) {

	proc, ok := f.processorsMap[reportType]
	if !ok {
		return nil, fmt.Errorf("НЕ ОБНАРУЖЕН ОБРАБОТЧИК (%s)", reportType)
	}

	return proc, nil

}

func NewProcessors() processorsFactory {
	return processorsFactory{
		processorsMap: map[string]ReportProcessor{
			"0001": &porocessorBusinesTrip{},
		},
	}
}

func NewReportService(repo *repository.Repo) *ReportService {
	return &ReportService{
		reports:    repo.Reports,
		processors: NewProcessors(),
	}
}

func (s *ReportService) GetTypes(types []string) ([]*reports.ReportTypes, error) {
	return s.reports.GetTypes(types)
}

func (s *ReportService) GetStructure(reportType string) (interface{}, error) {
	return s.reports.GetStructure(reportType)
}

func getContextUser(ctx context.Context) (*models.User, error) {

	user, ok := ctx.Value(models.CtxKey("user")).(*models.User)
	if !ok {
		return nil, errors.New("ERROR ON GET USER")
	}

	return user, nil

}

func (s *ReportService) Save(ctx context.Context, reportType string, data interface{}) error {

	reportData, ok := data.(*reports.ReportData)
	if !ok {
		return fmt.Errorf("INCORRECT DATA STRUCTURE")
	}

	proc, err := s.processors.GetProcessor(reportType)
	if err != nil {
		return err
	}
	proc.Check(data)

	user, err := getContextUser(ctx)
	if err != nil {
		return err
	}

	tx, err := s.reports.GetTransaction(ctx)
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()

	if reportData.Head.Ref.Blank() {
		uuid, err := models.GenerateUUID()
		if err != nil {
			return err
		}
		reportData.Head.Ref = uuid
	}

	if reportData.Head.Author.Blank() {
		reportData.Head.Author = user.Key
	}

	for _, coordinator := range reportData.Coordinators {
		if coordinator.CoordinatorRef.Blank() {
			return errors.New("COORDINATOR NOT SET")
		}
		if coordinator.Ref.Blank() {
			uuid, err := models.GenerateUUID()
			if err != nil {
				return err
			}
			coordinator.Ref = uuid
		}
		coordinator.ReportRef = reportData.Head.Ref
		if coordinator.WhoAuthor.Blank() {
			coordinator.WhoAuthor = user.Key
			coordinator.WhenAdded = models.JSONTime(time.Now())
		}
	}

	err = s.reports.Save(tx, ctx, reportData.Head)
	if err != nil {
		return err
	}

	err = s.reports.SaveCoordinators(tx, ctx, reportData.Coordinators)
	if err != nil {
		return err
	}
	err = s.reports.SaveDetails(tx, ctx, reportType, reportData.Head, reportData.Details)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		tx.Rollback()
		return ctx.Err()
	default:
		err := tx.Commit()
		if err != nil {
			return err
		}
	}

	return errors.ErrUnsupported
}

func (s *ReportService) List(ctx context.Context) ([]*models.Report, error) {
	user, err := getContextUser(ctx)
	if err != nil {
		return nil, err
	}

	reports, err := s.reports.List(ctx, user.Key)
	if err != nil {
		return nil, err
	}

	return reports, nil
}
