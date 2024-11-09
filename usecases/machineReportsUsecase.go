package usecases

import (
	"errors"
	"time"
	"zuck-my-clothe/zuck-my-clothe-backend/model"
	"zuck-my-clothe/zuck-my-clothe-backend/repository"

	"github.com/google/uuid"
)

type machineReportUsecase struct {
	machineReportRepository model.MachineReportsRepository
	machineRepository       repository.MachineRepository
	branchRepository        repository.BranchReopository
	employeeContract        repository.EmployeeContractRepository
}

func CreateNewMachineReportUsecase(machineReportRepository model.MachineReportsRepository, machineRepository repository.MachineRepository, branchRepository repository.BranchReopository, employeeContract repository.EmployeeContractRepository) model.MachineReportsUsecase {
	return &machineReportUsecase{machineReportRepository: machineReportRepository,
		machineRepository: machineRepository,
		branchRepository:  branchRepository,
		employeeContract:  employeeContract,
	}
}

func toMachineReportDetail(machineReport *model.MachineReports) interface{} {
	reportData := model.MachineReportDetail{
		ReportID:          machineReport.ReportID,
		UserID:            machineReport.UserID,
		ReportDescription: machineReport.ReportDescription,
		MacineSerial:      machineReport.MacineSerial,
		ReportStatus:      machineReport.ReportStatus,
	}
	return reportData
}

func (u *machineReportUsecase) checkSith(reportID string, userID string, userRole string) error {
	if userRole == "SuperAdmin" {
		return nil
	} else if userRole == "Client" {
		return errors.New("un authorized")
	}
	report, err := u.machineReportRepository.FindMachinereportByID(reportID)
	var isEmployee bool = false
	if err != nil {
		return err
	}
	var machineSerial string = report.MacineSerial
	machine, err := u.machineRepository.GetByMachineSerial(machineSerial)
	if err != nil {
		return err
	}
	var branchID string = machine.BranchID

	branch, err := u.branchRepository.GetByBranchID(branchID)
	if err != nil {
		return err
	}

	if branch.OwnerUserID == userID {
		return nil
	}

	contractList, err := u.employeeContract.GetByBranchID(branchID)
	if err != nil {
		return err
	}
	for _, contract := range *contractList {
		if contract.UserID == userID {
			isEmployee = true
			break
		}
	}
	if isEmployee {
		return nil
	}
	return errors.New("un authorized")
}

func (u *machineReportUsecase) CreateMachineReport(newReport *model.AddMachineReportDTO, userID string) (*interface{}, error) {
	data := model.MachineReports{
		ReportID:          uuid.New().String(),
		UserID:            userID,
		ReportDescription: newReport.ReportDescription,
		MacineSerial:      newReport.MacineSerial,
		ReportStatus:      model.ReportPending,
		CreatedAt:         time.Now(),
	}
	err := u.machineReportRepository.CreateMachineReport(&data)
	if err != nil {
		return nil, err
	}
	report, err := u.machineReportRepository.FindMachinereportByID(data.ReportID)
	if err != nil {
		return nil, err
	}
	detail := toMachineReportDetail(report)
	return &detail, nil
}

func (u *machineReportUsecase) FindMachineReportByUserID(userID string) (*[]interface{}, error) {
	machineReportList, err := u.machineReportRepository.FindMachineReportByUserID(userID)
	if err != nil {
		return nil, err
	}
	var result []interface{}
	for _, machineReport := range *machineReportList {
		result = append(result, toMachineReportDetail(&machineReport))
	}
	return &result, nil
}

func (u *machineReportUsecase) FindMachineReportByBranch(branchID string, userID string) (*[]interface{}, error) {
	reportList, err := u.machineReportRepository.FindMachineReportByBranch(branchID, userID)
	if err != nil {
		return nil, err
	}

	var result []interface{}
	for _, report := range *reportList {
		result = append(result, toMachineReportDetail(&report))
	}
	return &result, nil
}

func (u *machineReportUsecase) GetAll() (*[]model.MachineReports, error) {
	reportList, err := u.machineReportRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return reportList, nil
}

func (u *machineReportUsecase) UpdateMachineReportStatus(updateReport model.UpdateMachineReportStatusDTO, userID string, userRole string) (*interface{}, error) {
	err := u.checkSith(updateReport.ReportID, userID, userRole)
	if err != nil {
		return nil, err
	}
	if err := u.machineReportRepository.UpdateMachineReportStatus(updateReport); err != nil {
		return nil, err
	}
	updatedReport, err := u.machineReportRepository.FindMachinereportByID(updateReport.ReportID)
	if err != nil {
		return nil, err
	}
	result := toMachineReportDetail(updatedReport)
	return &result, nil
}

func (u *machineReportUsecase) DeleteMachineReport(reportID string, userID string, userRole string) error {
	errchk := u.checkSith(reportID, userID, userRole)
	if errchk != nil {
		return errchk
	}
	err := u.machineReportRepository.DeleteMachineReport(reportID)
	return err
}
