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

type ReportService struct {
	reports repository.ReportProvider
}

func NewReportService(repo *repository.Repo) *ReportService {
	return &ReportService{
		reports: repo.Reports,
	}
}

func (s *ReportService) GetTypes(types []string) ([]*reports.ReportTypes, error) {
	return s.reports.GetTypes(types)
}

func (s *ReportService) GetStructure(reportType string) (interface{}, error) {
	return s.reports.GetStructure(reportType)
}

func (s *ReportService) Save(ctx context.Context, data interface{}) error {

	reportData, ok := data.(*reports.ReportData)
	if !ok {
		return fmt.Errorf("INCORRECT DATA STRUCTURE")
	}

	tx, err := s.reports.GetTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if reportData.Head.Ref.Blank() {
		uuid, err := models.GenerateUUID()
		if err != nil {
			return err
		}
		reportData.Head.Ref = uuid
	}

	err = s.reports.Save(tx, ctx, reportData.Head)
	if err != nil {
		return err
	}

	user, ok := ctx.Value("user").(*models.User)
	if !ok {
		return errors.New("ERROR ON GET USER")
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

	err = s.reports.SaveCoordinators(tx, ctx, reportData.Coordinators)
	if err != nil {
		return err
	}
	err = s.reports.SaveDetails(tx, ctx, reportData.Head, reportData.Details)
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
