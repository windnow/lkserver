package services

import (
	"lkserver/internal/models"
	"lkserver/internal/models/reports"
)

type ReportData struct {
	Head         *models.Report          `json:"head"`
	Coordinators []*reports.Coordinators `json:"coordinators"`
	Details      any                     `json:"details"`
}
