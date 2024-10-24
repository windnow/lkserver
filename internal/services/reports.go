package services

import (
	"context"
	"errors"
	"fmt"
	"lkserver/internal/models"
	"lkserver/internal/models/reports"
	"lkserver/internal/models/types"
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
	provider   *repository.Repo
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
		provider:   repo,
		processors: NewProcessors(),
	}
}

func (s *ReportService) GetReportType(guid models.JSONByte) (*Result, error) {

	set, err := s.provider.Reports.GetTypes([]string{})
	if err != nil {
		return nil, models.HandleError(err, "ReportService.Get")
	}
	var result *reports.ReportTypes = nil
	for _, row := range set {
		if row.Ref.Equal(guid) {
			result = row
		}
	}
	if result == nil {
		return nil, models.ErrNotFound
	}

	return &Result{
		Data: result,
		Rows: s.provider.Reports.TypesCount(),
		Len:  1,
		Meta: map[string]models.META{types.ReportType: reports.ReportTypesMETA},
	}, nil

}

func (s *ReportService) GetTypes(t []string) (*Result, error) {
	result, err := s.provider.Reports.GetTypes(t)
	if err != nil {
		return nil, models.HandleError(err, "ReportService.GetTypes")
	}

	return &Result{
		Data: result,
		Len:  len(result),
		Rows: s.provider.Reports.TypesCount(),
		Meta: map[string]models.META{
			types.ReportType: reports.ReportTypesMETA,
		},
	}, nil
}

func (s *ReportService) GetStructure(reportType string) (interface{}, error) {
	return s.provider.Reports.GetStructure(reportType)
}

func (s *ReportService) NewReport(reportType string) (*Result, error) {
	meta := s.provider.Reports.META(reportType)
	mockData, err := s.provider.Reports.GetStructure(reportType)
	if err != nil {
		return nil, models.HandleError(err, "ReportService.NewReport")
	}

	result := &Result{
		Data: mockData,
		Meta: map[string]models.META{
			"head":         models.ReportMETA,
			"coordinators": reports.CoordinatorsMETA,
		},
	}

	for key, value := range meta {
		result.Meta[key] = value
	}

	return result, nil

}

func (s *ReportService) GetReportData(ctx context.Context, guid models.JSONByte) (*Result, error) {

	reportHead, err := s.provider.Reports.Get(guid)
	if err != nil {
		return nil, err
	}
	reportCoordinators, err := s.provider.Reports.GetCoordinators(ctx, reportHead.Ref)
	if err != nil {
		return nil, models.HandleError(err, "ReportService.GetReportData")
	}
	reportDetails, meta, err := s.provider.Reports.GetDetails(ctx, reportHead)
	if err != nil {
		return nil, models.HandleError(err, "ReportService.GetReportData")
	}
	reportType, err := s.provider.Reports.GetType(ctx, reportHead.Type)
	if err != nil {
		return nil, models.HandleError(err, "ReportService.GetReportData")
	}

	user, err := getContextUser(ctx)
	if err != nil {
		return nil, models.HandleError(err, "ReportService.GetReportData")
	}

	result := &Result{Data: &reports.ReportData{
		Head:         reportHead,
		Coordinators: reportCoordinators,
		Details:      reportDetails,
	},
		Len:  1,
		Rows: s.provider.Reports.Count(user.Key, reportType.Ref),
		Meta: map[string]models.META{
			"head":         models.ReportMETA,
			"coordinators": reports.CoordinatorsMETA,
		},
	}

	for key, value := range meta {
		result.Meta[key] = value
	}

	return result, nil
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

	if err := proc.Check(data); err != nil {
		return err
	}

	user, err := getContextUser(ctx)
	if err != nil {
		return err
	}

	tx, err := s.provider.Reports.GetTransaction(ctx)
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

	err = s.provider.Reports.Save(tx, ctx, reportData.Head)
	if err != nil {
		return err
	}

	err = s.provider.Reports.SaveCoordinators(tx, ctx, reportData.Coordinators)
	if err != nil {
		return err
	}
	err = s.provider.Reports.SaveDetails(tx, ctx, reportData.Head, reportData.Details)
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

	return nil
}

func (s *ReportService) List(ctx context.Context, typeCode string, limits ...int64) (*Result, error) {
	user, err := getContextUser(ctx)
	if err != nil {
		return nil, models.HandleError(err, "ReportService.List")
	}

	var reportType *reports.ReportTypes
	if len(typeCode) > 0 {
		reportType, err = s.provider.Reports.GetTypeByCode(typeCode)
		if err != nil {
			return nil, models.HandleError(err, "ReportService.List")
		}

	} else {
		reportType = &reports.ReportTypes{}
	}

	reportsList, err := s.provider.Reports.List(ctx, user.Key, reportType.Ref, limits...)
	if err != nil {
		return nil, models.HandleError(err, "ReportService.List")
	}

	return &Result{
		Data: reportsList,
		Len:  len(reportsList),
		Rows: s.provider.Reports.Count(user.Key, reportType.Ref),
		Meta: map[string]models.META{types.Report: models.ReportMETA},
	}, nil

}
