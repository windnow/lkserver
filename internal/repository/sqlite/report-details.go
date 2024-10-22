package sqlite

import (
	"context"
	"database/sql"
	m "lkserver/internal/models"
	"lkserver/internal/repository"
)

func (repo *reportsRepo) getReportProcessor(procGuid m.JSONByte) (repository.ReportDetails, error) {
	reportTypeByKey, err := repo.GetTypeCode(procGuid)
	if err != nil {
		return nil, err
	}

	processor, err := repo.factory.GetReportProcessor(reportTypeByKey)
	if err != nil {
		return nil, err
	}

	return processor, nil

}

func (repo *reportsRepo) SaveDetails(tx *sql.Tx, ctx context.Context, report *m.Report, data any) error {

	processor, err := repo.getReportProcessor(report.Type)
	if err != nil {
		return m.HandleError(err, "reportsRepo.SaveDetails")
	}
	return processor.Save(tx, ctx, report.Ref, data)
}

func (repo *reportsRepo) GetDetails(ctx context.Context, report *m.Report) (any, map[string]m.META, error) {

	processor, err := repo.getReportProcessor(report.Type)
	if err != nil {
		return nil, nil, m.HandleError(err, "reportsRepo.GetDetails")
	}

	return processor.Get(ctx, report.Ref)
}
