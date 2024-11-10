package model

import (
	"time"

	"gorm.io/gorm"
)

func (MachineReports) TableName() string {
	return "MachineReports"
}

type MachineReportStatus string

const (
	ReportPending    MachineReportStatus = "Pending"
	ReportInProgress MachineReportStatus = "In Progress"
	ReportFixed      MachineReportStatus = "Fixed"
	ReportCanceled   MachineReportStatus = "Cancel"
)

type MachineReports struct {
	ReportID          string              `json:"report_id" gorm:"column:report_id"`
	UserID            string              `json:"user_id" gorm:"column:user_id"`
	ReportDescription string              `json:"report_desc" gorm:"column:report_desc"`
	MacineSerial      string              `json:"machine_serial" gorm:"column:machine_serial"`
	ReportStatus      MachineReportStatus `json:"report_status" gorm:"column:report_status"`
	BranchID          string              `json:"-"  gorm:"column:branch_id"`
	CreatedAt         time.Time           `json:"created_at" gorm:"column:created_at"`
	DeletedAt         gorm.DeletedAt      `json:"deleted_at" gorm:"column:deleted_at;index" swaggertype:"string" example:"null"`
}

type AddMachineReportDTO struct {
	ReportDescription string `json:"report_desc" validate:"required"`
	MacineSerial      string `json:"machine_serial" validate:"required"`
}

type UpdateMachineReportStatusDTO struct {
	ReportID     string              `json:"report_id" validate:"required"`
	ReportStatus MachineReportStatus `json:"report_status" validate:"required"`
}

type MachineReportDetail struct {
	ReportID          string              `json:"report_id"`
	UserID            string              `json:"user_id,omitempty"`
	ReportDescription string              `json:"report_desc"`
	MacineSerial      string              `json:"machine_serial"`
	ReportStatus      MachineReportStatus `json:"report_status"`
	CreatedAt         time.Time           `json:"created_at"`
	DeletedAt         *gorm.DeletedAt     `json:"deleted_at,omitempty" gorm:"column:deleted_at;index" swaggertype:"string" example:"null"`
	BranchInfo        BranchDetail        `json:"branch"`
}

type MachineReportsRepository interface {
	CreateMachineReport(newReport *MachineReports) error
	FindMachinereportByID(machineReportID string) (*MachineReports, error)
	FindMachineReportByUserID(userID string) (*[]MachineReports, error)
	FindMachineReportByBranch(branchID string, userID string,userRole string) (*[]MachineReports, error)
	GetAll() (*[]MachineReports, error)
	UpdateMachineReportStatus(updateReport UpdateMachineReportStatusDTO) error
	DeleteMachineReport(reportID string) error
}

type MachineReportsUsecase interface {
	CreateMachineReport(newReport *AddMachineReportDTO, userID string) (*interface{}, error)
	FindMachineReportByUserID(userID string) (*[]interface{}, error)
	FindMachineReportByBranch(branchID string, userID string, userRole string) (*[]interface{}, error)
	GetAll() (*[]interface{}, error)
	UpdateMachineReportStatus(updateReport UpdateMachineReportStatusDTO, userID string, userRole string) (*interface{}, error)
	DeleteMachineReport(reportID string, userID string, userRole string) error
}
