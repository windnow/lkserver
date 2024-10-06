package services

import (
	"context"
	"errors"
	"fmt"
	"lkserver/internal/models/reports"
	"lkserver/internal/repository"
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

func (s *ReportService) Save(ctx context.Context, data interface{}) error {

	reportData, ok := data.(*ReportData)
	if !ok {
		return fmt.Errorf("INCORRECT DATA STRUCTURE")
	}

	tx, err := s.reports.GetTransaction(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.reports.Save(tx, ctx, reportData.Head)
	if err != nil {
		return err
	}
	err = s.reports.SaveCoordinators(tx, ctx, reportData.Coordinators)
	if err != nil {
		return nil
	}
	err = s.reports.SaveDetails(tx, ctx, reportData.Head, reportData.Details)
	if err != nil {
		return nil
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
