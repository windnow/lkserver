package services

import (
	"lkserver/internal/models"
	"lkserver/internal/models/reports"
)

type ReportData struct {
	Head         *models.Report
	Coordinators []*reports.Coordinators
	Details      any
}
