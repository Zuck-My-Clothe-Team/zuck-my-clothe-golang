package repository

import (
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/platform"

	"gorm.io/gorm"
)

type machineReportsRepository struct {
	db *platform.Postgres
}

func CreateNewMachineReportRepository(db *platform.Postgres) model.MachineReportsRepository {
	return &machineReportsRepository{db: db}
}

func (u *machineReportsRepository) CreateMachineReport(newReport *model.MachineReports) error {
	dbTx := u.db.Create(newReport)
	return dbTx.Error
}

func (u *machineReportsRepository) FindMachinereportByID(machineReportID string) (*model.MachineReports, error) {
	machineReport := new(model.MachineReports)
	dbTx := u.db.First(machineReport, "report_id = ?", machineReportID)
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	if dbTx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return machineReport, nil
}

func (u *machineReportsRepository) FindMachineReportByUserID(userID string) (*[]model.MachineReports, error) {
	machineReportLists := new([]model.MachineReports)
	dbTx := u.db.Find(machineReportLists, "user_id = ?", userID)
	if dbTx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return machineReportLists, dbTx.Error
}

func (u *machineReportsRepository) FindMachineReportByBranch(branchID string, userID string) (*[]model.MachineReports, error) {
	machineReportLists := new([]model.MachineReports)

	dbTx := u.db.Raw(
		`SELECT "MachineReports".report_id,"MachineReports".user_id,"MachineReports".report_desc,"MachineReports".machine_serial,"MachineReports".report_status,"MachineReports".created_at,"MachineReports".deleted_at
		FROM "MachineReports"
		INNER JOIN "Machines" ON "Machines".machine_serial = "MachineReports".machine_serial
		WHERE "MachineReports".deleted_at IS NULL AND "Machines".branch_id IN (
			SELECT DISTINCT "EmployeeContracts".branch_id
			FROM "EmployeeContracts"
			WHERE "EmployeeContracts".user_id = ? AND "EmployeeContracts".branch_id = ?
			)
			OR  "Machines".branch_id IN (
			SELECT DISTINCT "Branches".branch_id
			FROM "Branches"
			WHERE "Branches".owner_user_id = ? AND "Branches".branch_id = ?
		)`,
		userID, branchID, userID, branchID).Scan(machineReportLists)

	if dbTx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return machineReportLists, nil
}

func (u *machineReportsRepository) GetAll() (*[]model.MachineReports, error) {
	machineReportLists := new([]model.MachineReports)
	dbTx := u.db.Raw(`
	SELECT *
	FROM "MachineReports"
	`).Scan(machineReportLists)
	if dbTx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if dbTx.Error != nil {
		return nil, dbTx.Error
	}
	return machineReportLists, nil
}

func (u *machineReportsRepository) UpdateMachineReportStatus(updateReport model.UpdateMachineReportStatusDTO) error {
	dbTx := u.db.Table("MachineReports").Where("report_id = ?", updateReport.ReportID).Update("report_status", updateReport.ReportStatus)
	if dbTx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return dbTx.Error
}

func (u *machineReportsRepository) DeleteMachineReport(reportID string) error {
	deletedReport := new(model.MachineReports)
	dbTx := u.db.Table("MachineReports").Where("report_id = ?", reportID).Delete(deletedReport)
	if dbTx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return dbTx.Error
}
