package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gin-project/models"
	"gin-project/utils"

	"gorm.io/gorm"
)

type SalaryServiceInterface interface {
	// Legacy Salary Management
	GetSalaryRecords(params SalaryQueryParams) (*utils.PaginationResponse, error)
	GetSalaryByID(id uint) (*models.Salary, error)
	CreateSalary(salary *models.Salary) (*models.Salary, error)
	UpdateSalary(salary *models.Salary) (*models.Salary, error)
	DeleteSalary(id uint) error
	CalculateSalary(employeeID uint, month string) (*models.Salary, error)
	BatchCalculateSalary(month string, departmentID uint, employeeIDs []uint) (*BatchCalculateResult, error)
	ApproveSalary(id uint, status, remark string) (*models.Salary, error)
	GetSalaryStatistics(month string, departmentID uint) (*SalaryStatistics, error)
	ExportSalaryReport(month string, departmentID uint, format string) ([]byte, string, error)
	GetEmployeeSalary(userID uint, month string) (*models.Salary, error)
	CreatePayrollRecord(payroll *models.PayrollRecord) (*models.PayrollRecord, error)
	ProcessPayroll(salaryID uint, paymentMethod, bankAccount string) (*models.PayrollRecord, error)
	GetPayrollRecords(salaryID uint) ([]models.PayrollRecord, error)
	UpdatePayrollStatus(id uint, status, remark string) error

	// Enhanced Salary Component Management
	CreateSalaryComponent(component *models.SalaryComponent) (*models.SalaryComponent, error)
	UpdateSalaryComponent(id uint, component *models.SalaryComponent) (*models.SalaryComponent, error)
	DeleteSalaryComponent(id uint) error
	GetSalaryComponents(params ComponentQueryParams) (*utils.PaginationResponse, error)
	GetSalaryComponentByID(id uint) (*models.SalaryComponent, error)

	// Enhanced Salary Grade Management
	CreateSalaryGrade(grade *models.SalaryGrade) (*models.SalaryGrade, error)
	UpdateSalaryGrade(id uint, grade *models.SalaryGrade) (*models.SalaryGrade, error)
	DeleteSalaryGrade(id uint) error
	GetSalaryGrades(params GradeQueryParams) (*utils.PaginationResponse, error)
	GetSalaryGradeByID(id uint) (*models.SalaryGrade, error)

	// Enhanced Salary Structure Management
	CreateSalaryStructure(structure *models.SalaryStructure) (*models.SalaryStructure, error)
	UpdateSalaryStructure(id uint, structure *models.SalaryStructure) (*models.SalaryStructure, error)
	DeleteSalaryStructure(id uint) error
	GetSalaryStructures(params StructureQueryParams) (*utils.PaginationResponse, error)
	GetSalaryStructureByID(id uint) (*models.SalaryStructure, error)
	GetApplicableStructure(employeeID uint) (*models.SalaryStructure, error)

	// Enhanced Payroll Period Management
	CreatePayrollPeriod(period *models.PayrollPeriod) (*models.PayrollPeriod, error)
	UpdatePayrollPeriod(id uint, period *models.PayrollPeriod) (*models.PayrollPeriod, error)
	GetPayrollPeriods(params PeriodQueryParams) (*utils.PaginationResponse, error)
	GetPayrollPeriodByID(id uint) (*models.PayrollPeriod, error)
	LockPayrollPeriod(id uint) error
	UnlockPayrollPeriod(id uint) error

	// Enhanced Salary Management
	CalculateEmployeeSalary(employeeID, periodID uint, userID uint) (*models.EnhancedSalary, error)
	BatchCalculateSalaries(periodID uint, departmentID *uint, employeeIDs []uint, userID uint) (*EnhancedBatchResult, error)
	GetEnhancedSalaries(params EnhancedSalaryQueryParams) (*utils.PaginationResponse, error)
	GetEnhancedSalaryByID(id uint) (*models.EnhancedSalary, error)
	UpdateSalaryDetails(salaryID uint, details []SalaryDetailUpdate, userID uint) (*models.EnhancedSalary, error)

	// Enhanced Approval Workflow
	ReviewSalary(id uint, reviewerID uint, notes string, approve bool) (*models.EnhancedSalary, error)
	ApproveEnhancedSalary(id uint, approverID uint, notes string) (*models.EnhancedSalary, error)
	RejectSalary(id uint, approverID uint, reason string) (*models.EnhancedSalary, error)
	BulkApproveSalaries(salaryIDs []uint, approverID uint, notes string) (*BulkApprovalResult, error)

	// Enhanced Payment Processing
	CreatePaymentBatch(batch *models.PaymentBatch, salaryIDs []uint, userID uint) (*models.PaymentBatch, error)
	ProcessPaymentBatch(batchID uint, userID uint) (*models.PaymentBatch, error)
	GetPaymentBatches(params BatchQueryParams) (*utils.PaginationResponse, error)
	GetPaymentBatchByID(id uint) (*models.PaymentBatch, error)

	// Enhanced Salary Adjustments
	CreateSalaryAdjustment(adjustment *models.SalaryAdjustment) (*models.SalaryAdjustment, error)
	ApproveSalaryAdjustment(id uint, approverID uint) (*models.SalaryAdjustment, error)
	GetSalaryAdjustments(params AdjustmentQueryParams) (*utils.PaginationResponse, error)

	// Enhanced Analytics and Reporting
	GetSalaryAnalytics(params AnalyticsParams) (*SalaryAnalytics, error)
	GetDepartmentSalaryReport(departmentID uint, periodID uint) (*DepartmentSalaryReport, error)
	ExportEnhancedSalaryReport(params ExportParams) ([]byte, string, error)

	// Personal Salary Management
	GetPersonalSalaryDetail(params PersonalSalaryParams) (*PersonalSalaryDetail, error)
	GetSalaryHistory(params SalaryHistoryParams) (*SalaryHistoryResponse, error)
	GetPersonalPayrollRecords(params PersonalPayrollParams) (*utils.PaginationResponse, error)
	GetPersonalSalaryDashboard(employeeID uint) (*PersonalSalaryDashboard, error)

	// Formula Engine
	EvaluateFormula(formula string, context FormulaContext) (float64, error)
	ValidateFormula(formula string) error
}

type SalaryService struct {
	db *gorm.DB
}

type SalaryQueryParams struct {
	EmployeeID uint
	Month      string
	Status     string
	Page       int
	PageSize   int
}

type BatchCalculateResult struct {
	Total     int                   `json:"total"`
	Success   int                   `json:"success"`
	Failed    int                   `json:"failed"`
	Results   []SalaryCalculateItem `json:"results"`
}

type SalaryCalculateItem struct {
	EmployeeID   uint    `json:"employee_id"`
	EmployeeName string  `json:"employee_name"`
	Success      bool    `json:"success"`
	Message      string  `json:"message"`
	Salary       *models.Salary `json:"salary,omitempty"`
}

type SalaryStatistics struct {
	TotalEmployees  int     `json:"total_employees"`
	TotalSalary     float64 `json:"total_salary"`
	AverageSalary   float64 `json:"average_salary"`
	MaxSalary       float64 `json:"max_salary"`
	MinSalary       float64 `json:"min_salary"`
	DepartmentStats []DepartmentSalaryStat `json:"department_stats"`
}

type DepartmentSalaryStat struct {
	DepartmentID   uint    `json:"department_id"`
	DepartmentName string  `json:"department_name"`
	EmployeeCount  int     `json:"employee_count"`
	TotalSalary    float64 `json:"total_salary"`
	AverageSalary  float64 `json:"average_salary"`
}

// Enhanced Query Parameters are defined at the top

func NewSalaryService(db *gorm.DB) SalaryServiceInterface {
	return &SalaryService{
		db: db,
	}
}

// InjectDependencies implements DependencyInjector interface
func (ss *SalaryService) InjectDependencies(deps ...interface{}) error {
	for _, dep := range deps {
		switch d := dep.(type) {
		case *gorm.DB:
			ss.db = d
		}
	}
	return nil
}

func (s *SalaryService) GetSalaryRecords(params SalaryQueryParams) (*utils.PaginationResponse, error) {
	var salaries []models.Salary
	var total int64

	query := s.db.Model(&models.Salary{}).Preload("Employee").Preload("Employee.Department")

	if params.EmployeeID > 0 {
		query = query.Where("employee_id = ?", params.EmployeeID)
	}
	if params.Month != "" {
		query = query.Where("month = ?", params.Month)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Find(&salaries).Error; err != nil {
		return nil, err
	}

	response := utils.CreatePaginationResponse(salaries, params.Page, params.PageSize, total)
	return &response, nil
}

func (s *SalaryService) GetSalaryByID(id uint) (*models.Salary, error) {
	var salary models.Salary
	if err := s.db.Preload("Employee").Preload("Employee.Department").First(&salary, id).Error; err != nil {
		return nil, err
	}
	return &salary, nil
}

func (s *SalaryService) CreateSalary(salary *models.Salary) (*models.Salary, error) {
	if err := s.validateSalary(salary); err != nil {
		return nil, err
	}

	if err := s.db.Create(salary).Error; err != nil {
		return nil, err
	}
	return salary, nil
}

func (s *SalaryService) UpdateSalary(salary *models.Salary) (*models.Salary, error) {
	if err := s.validateSalary(salary); err != nil {
		return nil, err
	}

	if err := s.db.Save(salary).Error; err != nil {
		return nil, err
	}
	return salary, nil
}

func (s *SalaryService) DeleteSalary(id uint) error {
	return s.db.Delete(&models.Salary{}, id).Error
}

func (s *SalaryService) CalculateSalary(employeeID uint, month string) (*models.Salary, error) {
	var employee models.Employee
	if err := s.db.Preload("Department").First(&employee, employeeID).Error; err != nil {
		return nil, errors.New("员工不存在")
	}

	existingSalary := &models.Salary{}
	if err := s.db.Where("employee_id = ? AND month = ?", employeeID, month).First(existingSalary).Error; err == nil {
		return nil, errors.New("该月份薪资已存在")
	}

	attendance, err := s.getEmployeeAttendance(employeeID, month)
	if err != nil {
		return nil, err
	}

	salary := &models.Salary{
		EmployeeID:  employeeID,
		Month:       month,
		BaseSalary:  employee.BaseSalary,
		Bonus:       s.calculateBonus(employee, attendance),
		Allowance:   s.calculateAllowance(employee),
		Deduction:   s.calculateDeduction(employee, attendance),
		Status:      "calculated",
	}

	salary.GrossSalary = salary.BaseSalary + salary.Bonus + salary.Allowance - salary.Deduction
	salary.Tax = s.calculateTax(salary.GrossSalary)
	salary.SocialSecurity = s.calculateSocialSecurity(salary.BaseSalary)
	salary.HousingFund = s.calculateHousingFund(salary.BaseSalary)
	salary.NetSalary = salary.GrossSalary - salary.Tax - salary.SocialSecurity - salary.HousingFund

	if err := s.db.Create(salary).Error; err != nil {
		return nil, err
	}

	return salary, nil
}

func (s *SalaryService) BatchCalculateSalary(month string, departmentID uint, employeeIDs []uint) (*BatchCalculateResult, error) {
	var employees []models.Employee
	query := s.db.Where("status = ?", "active")

	if departmentID > 0 {
		query = query.Where("department_id = ?", departmentID)
	}
	if len(employeeIDs) > 0 {
		query = query.Where("id IN ?", employeeIDs)
	}

	if err := query.Find(&employees).Error; err != nil {
		return nil, err
	}

	result := &BatchCalculateResult{
		Total:   len(employees),
		Results: make([]SalaryCalculateItem, 0, len(employees)),
	}

	for _, employee := range employees {
		item := SalaryCalculateItem{
			EmployeeID:   employee.ID,
			EmployeeName: employee.Name,
		}

		salary, err := s.CalculateSalary(employee.ID, month)
		if err != nil {
			item.Success = false
			item.Message = err.Error()
			result.Failed++
		} else {
			item.Success = true
			item.Message = "计算成功"
			item.Salary = salary
			result.Success++
		}

		result.Results = append(result.Results, item)
	}

	return result, nil
}

func (s *SalaryService) ApproveSalary(id uint, status, remark string) (*models.Salary, error) {
	var salary models.Salary
	if err := s.db.First(&salary, id).Error; err != nil {
		return nil, err
	}

	salary.Status = status
	salary.Remark = remark

	if err := s.db.Save(&salary).Error; err != nil {
		return nil, err
	}

	return &salary, nil
}

func (s *SalaryService) GetSalaryStatistics(month string, departmentID uint) (*SalaryStatistics, error) {
	var stats SalaryStatistics

	query := s.db.Model(&models.Salary{}).
		Joins("JOIN employees ON salaries.employee_id = employees.id")

	if month != "" {
		query = query.Where("salaries.month = ?", month)
	}
	if departmentID > 0 {
		query = query.Where("employees.department_id = ?", departmentID)
	}

	var aggregateResult struct {
		Count   int     `gorm:"column:count"`
		Total   float64 `gorm:"column:total"`
		Average float64 `gorm:"column:average"`
		Max     float64 `gorm:"column:max"`
		Min     float64 `gorm:"column:min"`
	}

	if err := query.Select("COUNT(*) as count, SUM(net_salary) as total, AVG(net_salary) as average, MAX(net_salary) as max, MIN(net_salary) as min").
		Scan(&aggregateResult).Error; err != nil {
		return nil, err
	}

	stats.TotalEmployees = aggregateResult.Count
	stats.TotalSalary = aggregateResult.Total
	stats.AverageSalary = aggregateResult.Average
	stats.MaxSalary = aggregateResult.Max
	stats.MinSalary = aggregateResult.Min

	var departmentStats []DepartmentSalaryStat
	deptQuery := s.db.Model(&models.Salary{}).
		Select("departments.id as department_id, departments.name as department_name, COUNT(*) as employee_count, SUM(salaries.net_salary) as total_salary, AVG(salaries.net_salary) as average_salary").
		Joins("JOIN employees ON salaries.employee_id = employees.id").
		Joins("JOIN departments ON employees.department_id = departments.id").
		Group("departments.id, departments.name")

	if month != "" {
		deptQuery = deptQuery.Where("salaries.month = ?", month)
	}
	if departmentID > 0 {
		deptQuery = deptQuery.Where("departments.id = ?", departmentID)
	}

	if err := deptQuery.Scan(&departmentStats).Error; err != nil {
		return nil, err
	}

	stats.DepartmentStats = departmentStats

	return &stats, nil
}

func (s *SalaryService) ExportSalaryReport(month string, departmentID uint, format string) ([]byte, string, error) {
	return nil, "", errors.New("导出功能暂未实现")
}

func (s *SalaryService) GetEmployeeSalary(userID uint, month string) (*models.Salary, error) {
	var employee models.Employee
	if err := s.db.Where("id = ?", userID).First(&employee).Error; err != nil {
		return nil, errors.New("员工不存在")
	}

	var salary models.Salary
	query := s.db.Where("employee_id = ?", employee.ID)
	if month != "" {
		query = query.Where("month = ?", month)
	}

	if err := query.Order("month DESC").First(&salary).Error; err != nil {
		return nil, errors.New("未找到薪资记录")
	}

	return &salary, nil
}

func (s *SalaryService) validateSalary(salary *models.Salary) error {
	if salary.EmployeeID == 0 {
		return errors.New("员工ID不能为空")
	}
	if salary.Month == "" {
		return errors.New("薪资月份不能为空")
	}
	return nil
}

func (s *SalaryService) getEmployeeAttendance(employeeID uint, month string) ([]models.Attendance, error) {
	var attendances []models.Attendance
	startDate := month + "-01"
	endDate := month + "-31"

	if err := s.db.Where("employee_id = ? AND date BETWEEN ? AND ?", employeeID, startDate, endDate).
		Find(&attendances).Error; err != nil {
		return nil, err
	}

	return attendances, nil
}

func (s *SalaryService) calculateBonus(employee models.Employee, attendances []models.Attendance) float64 {
	if len(attendances) == 0 {
		return 0
	}

	totalWorkDays := len(attendances)
	normalWorkDays := 22

	if totalWorkDays >= normalWorkDays {
		return employee.BaseSalary * 0.1
	}

	return 0
}

func (s *SalaryService) calculateAllowance(employee models.Employee) float64 {
	return 500
}

func (s *SalaryService) calculateDeduction(employee models.Employee, attendances []models.Attendance) float64 {
	deduction := 0.0
	normalWorkDays := 22
	actualWorkDays := len(attendances)

	if actualWorkDays < normalWorkDays {
		missedDays := normalWorkDays - actualWorkDays
		dailySalary := employee.BaseSalary / 30
		deduction = float64(missedDays) * dailySalary
	}

	return deduction
}

func (s *SalaryService) calculateTax(grossSalary float64) float64 {
	taxableIncome := grossSalary - 5000

	if taxableIncome <= 0 {
		return 0
	}

	var tax float64
	if taxableIncome <= 3000 {
		tax = taxableIncome * 0.03
	} else if taxableIncome <= 12000 {
		tax = 3000*0.03 + (taxableIncome-3000)*0.1
	} else if taxableIncome <= 25000 {
		tax = 3000*0.03 + 9000*0.1 + (taxableIncome-12000)*0.2
	} else if taxableIncome <= 35000 {
		tax = 3000*0.03 + 9000*0.1 + 13000*0.2 + (taxableIncome-25000)*0.25
	} else if taxableIncome <= 55000 {
		tax = 3000*0.03 + 9000*0.1 + 13000*0.2 + 10000*0.25 + (taxableIncome-35000)*0.3
	} else if taxableIncome <= 80000 {
		tax = 3000*0.03 + 9000*0.1 + 13000*0.2 + 10000*0.25 + 20000*0.3 + (taxableIncome-55000)*0.35
	} else {
		tax = 3000*0.03 + 9000*0.1 + 13000*0.2 + 10000*0.25 + 20000*0.3 + 25000*0.35 + (taxableIncome-80000)*0.45
	}

	return tax
}

func (s *SalaryService) calculateSocialSecurity(baseSalary float64) float64 {
	return baseSalary * 0.08
}

func (s *SalaryService) calculateHousingFund(baseSalary float64) float64 {
	return baseSalary * 0.12
}

func (s *SalaryService) CreatePayrollRecord(payroll *models.PayrollRecord) (*models.PayrollRecord, error) {
	if err := s.validatePayrollRecord(payroll); err != nil {
		return nil, err
	}

	if err := s.db.Create(payroll).Error; err != nil {
		return nil, err
	}

	return payroll, nil
}

func (s *SalaryService) ProcessPayroll(salaryID uint, paymentMethod, bankAccount string) (*models.PayrollRecord, error) {
	var salary models.Salary
	if err := s.db.Preload("Employee").First(&salary, salaryID).Error; err != nil {
		return nil, errors.New("薪资记录不存在")
	}

	if salary.Status != "approved" {
		return nil, errors.New("薪资记录尚未审批通过")
	}

	now := time.Now()
	payroll := &models.PayrollRecord{
		SalaryID:      salaryID,
		PaymentDate:   &now,
		PaymentMethod: paymentMethod,
		BankAccount:   bankAccount,
		PaymentAmount: salary.NetSalary,
		Status:        "processing",
		Remark:        "系统自动发放",
	}

	if err := s.db.Create(payroll).Error; err != nil {
		return nil, err
	}

	return payroll, nil
}

func (s *SalaryService) GetPayrollRecords(salaryID uint) ([]models.PayrollRecord, error) {
	var records []models.PayrollRecord
	if err := s.db.Where("salary_id = ?", salaryID).
		Preload("Processor").
		Order("created_at DESC").
		Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func (s *SalaryService) UpdatePayrollStatus(id uint, status, remark string) error {
	return s.db.Model(&models.PayrollRecord{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": status,
			"remark": remark,
		}).Error
}

func (s *SalaryService) validatePayrollRecord(payroll *models.PayrollRecord) error {
	if payroll.SalaryID == 0 {
		return errors.New("薪资记录ID不能为空")
	}
	if payroll.PaymentAmount <= 0 {
		return errors.New("发放金额必须大于0")
	}
	return nil
}

// ========================= Enhanced Salary Component Management =========================

func (s *SalaryService) CreateSalaryComponent(component *models.SalaryComponent) (*models.SalaryComponent, error) {
	if err := s.validateSalaryComponent(component); err != nil {
		return nil, err
	}

	if err := s.db.Create(component).Error; err != nil {
		return nil, fmt.Errorf("failed to create salary component: %w", err)
	}

	return component, nil
}

func (s *SalaryService) UpdateSalaryComponent(id uint, component *models.SalaryComponent) (*models.SalaryComponent, error) {
	if err := s.validateSalaryComponent(component); err != nil {
		return nil, err
	}

	component.ID = id
	if err := s.db.Save(component).Error; err != nil {
		return nil, fmt.Errorf("failed to update salary component: %w", err)
	}

	return component, nil
}

func (s *SalaryService) DeleteSalaryComponent(id uint) error {
	return s.db.Delete(&models.SalaryComponent{}, id).Error
}

func (s *SalaryService) GetSalaryComponents(params ComponentQueryParams) (*utils.PaginationResponse, error) {
	var components []models.SalaryComponent
	var total int64

	query := s.db.Model(&models.SalaryComponent{})

	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&components).Error; err != nil {
		return nil, err
	}

	response := utils.CreatePaginationResponse(components, params.Page, params.PageSize, total)
	return &response, nil
}

func (s *SalaryService) GetSalaryComponentByID(id uint) (*models.SalaryComponent, error) {
	var component models.SalaryComponent
	if err := s.db.First(&component, id).Error; err != nil {
		return nil, err
	}
	return &component, nil
}

// ========================= Enhanced Salary Grade Management =========================

func (s *SalaryService) CreateSalaryGrade(grade *models.SalaryGrade) (*models.SalaryGrade, error) {
	if err := s.db.Create(grade).Error; err != nil {
		return nil, err
	}
	return grade, nil
}

func (s *SalaryService) UpdateSalaryGrade(id uint, grade *models.SalaryGrade) (*models.SalaryGrade, error) {
	grade.ID = id
	if err := s.db.Save(grade).Error; err != nil {
		return nil, err
	}
	return grade, nil
}

func (s *SalaryService) DeleteSalaryGrade(id uint) error {
	return s.db.Delete(&models.SalaryGrade{}, id).Error
}

func (s *SalaryService) GetSalaryGrades(params GradeQueryParams) (*utils.PaginationResponse, error) {
	var grades []models.SalaryGrade
	var total int64

	query := s.db.Model(&models.SalaryGrade{})

	if params.Level > 0 {
		query = query.Where("level = ?", params.Level)
	}
	if params.Currency != "" {
		query = query.Where("currency = ?", params.Currency)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("level ASC").Find(&grades).Error; err != nil {
		return nil, err
	}

	response := utils.CreatePaginationResponse(grades, params.Page, params.PageSize, total)
	return &response, nil
}

func (s *SalaryService) GetSalaryGradeByID(id uint) (*models.SalaryGrade, error) {
	var grade models.SalaryGrade
	if err := s.db.First(&grade, id).Error; err != nil {
		return nil, err
	}
	return &grade, nil
}

// ========================= Enhanced Salary Structure Management =========================

func (s *SalaryService) CreateSalaryStructure(structure *models.SalaryStructure) (*models.SalaryStructure, error) {
	if err := s.db.Create(structure).Error; err != nil {
		return nil, err
	}
	return structure, nil
}

func (s *SalaryService) UpdateSalaryStructure(id uint, structure *models.SalaryStructure) (*models.SalaryStructure, error) {
	structure.ID = id
	if err := s.db.Save(structure).Error; err != nil {
		return nil, err
	}
	return structure, nil
}

func (s *SalaryService) DeleteSalaryStructure(id uint) error {
	return s.db.Delete(&models.SalaryStructure{}, id).Error
}

func (s *SalaryService) GetSalaryStructures(params StructureQueryParams) (*utils.PaginationResponse, error) {
	var structures []models.SalaryStructure
	var total int64

	query := s.db.Model(&models.SalaryStructure{}).Preload("Components")

	if params.DepartmentID != nil {
		query = query.Where("department_id = ?", *params.DepartmentID)
	}
	if params.PositionID != nil {
		query = query.Where("position_id = ?", *params.PositionID)
	}
	if params.JobLevelID != nil {
		query = query.Where("job_level_id = ?", *params.JobLevelID)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&structures).Error; err != nil {
		return nil, err
	}

	response := utils.CreatePaginationResponse(structures, params.Page, params.PageSize, total)
	return &response, nil
}

func (s *SalaryService) GetSalaryStructureByID(id uint) (*models.SalaryStructure, error) {
	var structure models.SalaryStructure
	if err := s.db.Preload("Components").First(&structure, id).Error; err != nil {
		return nil, err
	}
	return &structure, nil
}

func (s *SalaryService) GetApplicableStructure(employeeID uint) (*models.SalaryStructure, error) {
	var employee models.Employee
	if err := s.db.Preload("Department").Preload("Position").First(&employee, employeeID).Error; err != nil {
		return nil, err
	}

	var structure models.SalaryStructure
	query := s.db.Preload("Components").Where("status = ?", "active")

	// Try to find structure by department first
	if employee.DepartmentID > 0 {
		err := query.Where("department_id = ?", employee.DepartmentID).First(&structure).Error
		if err == nil {
			return &structure, nil
		}
	}

	// Fall back to default structure
	if err := query.Where("is_default = ?", true).First(&structure).Error; err != nil {
		return nil, errors.New("no applicable salary structure found")
	}

	return &structure, nil
}

// ========================= Enhanced Payroll Period Management =========================

func (s *SalaryService) CreatePayrollPeriod(period *models.PayrollPeriod) (*models.PayrollPeriod, error) {
	if err := s.db.Create(period).Error; err != nil {
		return nil, err
	}
	return period, nil
}

func (s *SalaryService) UpdatePayrollPeriod(id uint, period *models.PayrollPeriod) (*models.PayrollPeriod, error) {
	period.ID = id
	if err := s.db.Save(period).Error; err != nil {
		return nil, err
	}
	return period, nil
}

func (s *SalaryService) GetPayrollPeriods(params PeriodQueryParams) (*utils.PaginationResponse, error) {
	var periods []models.PayrollPeriod
	var total int64

	query := s.db.Model(&models.PayrollPeriod{})

	if params.PeriodType != "" {
		query = query.Where("period_type = ?", params.PeriodType)
	}
	if params.Year > 0 {
		query = query.Where("year = ?", params.Year)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("start_date DESC").Find(&periods).Error; err != nil {
		return nil, err
	}

	response := utils.CreatePaginationResponse(periods, params.Page, params.PageSize, total)
	return &response, nil
}

func (s *SalaryService) GetPayrollPeriodByID(id uint) (*models.PayrollPeriod, error) {
	var period models.PayrollPeriod
	if err := s.db.First(&period, id).Error; err != nil {
		return nil, err
	}
	return &period, nil
}

func (s *SalaryService) LockPayrollPeriod(id uint) error {
	return s.db.Model(&models.PayrollPeriod{}).Where("id = ?", id).Update("status", "locked").Error
}

func (s *SalaryService) UnlockPayrollPeriod(id uint) error {
	return s.db.Model(&models.PayrollPeriod{}).Where("id = ?", id).Update("status", "open").Error
}

func (s *SalaryService) validateSalaryComponent(component *models.SalaryComponent) error {
	if component.Code == "" {
		return errors.New("salary component code is required")
	}
	if component.Name == "" {
		return errors.New("salary component name is required")
	}
	return nil
}

// ========================= Enhanced Salary Management =========================

func (s *SalaryService) CalculateEmployeeSalary(employeeID, periodID uint, userID uint) (*models.EnhancedSalary, error) {
	// Get employee and period
	var employee models.Employee
	if err := s.db.Preload("Department").Preload("Position").First(&employee, employeeID).Error; err != nil {
		return nil, fmt.Errorf("employee not found: %w", err)
	}

	var period models.PayrollPeriod
	if err := s.db.First(&period, periodID).Error; err != nil {
		return nil, fmt.Errorf("payroll period not found: %w", err)
	}

	// Check if salary already exists
	var existingSalary models.EnhancedSalary
	if err := s.db.Where("employee_id = ? AND payroll_period_id = ?", employeeID, periodID).First(&existingSalary).Error; err == nil {
		return nil, errors.New("salary for this period already exists")
	}

	// Create new salary record with basic calculation
	salary := &models.EnhancedSalary{
		EmployeeID:      employeeID,
		PayrollPeriodID: periodID,
		GrossSalary:     employee.BaseSalary,
		NetSalary:       employee.BaseSalary * 0.8, // Simple calculation for demo
		TotalDeductions: employee.BaseSalary * 0.2,
		Status:          "calculated",
		CalculatedBy:    &userID,
		Version:         1,
	}

	now := time.Now()
	salary.CalculatedAt = &now

	// Save salary record
	if err := s.db.Create(salary).Error; err != nil {
		return nil, fmt.Errorf("failed to create salary record: %w", err)
	}

	// Load complete salary record with relationships
	if err := s.db.Preload("Employee").Preload("PayrollPeriod").First(salary, salary.ID).Error; err != nil {
		return nil, err
	}

	return salary, nil
}

func (s *SalaryService) BatchCalculateSalaries(periodID uint, departmentID *uint, employeeIDs []uint, userID uint) (*EnhancedBatchResult, error) {
	var employees []models.Employee
	query := s.db.Where("status = ?", "active")

	if departmentID != nil {
		query = query.Where("department_id = ?", *departmentID)
	}
	if len(employeeIDs) > 0 {
		query = query.Where("id IN ?", employeeIDs)
	}

	if err := query.Find(&employees).Error; err != nil {
		return nil, err
	}

	result := &EnhancedBatchResult{
		Total:   len(employees),
		Results: make([]EnhancedCalculateItem, 0, len(employees)),
	}

	for _, employee := range employees {
		item := EnhancedCalculateItem{
			EmployeeID:   employee.ID,
			EmployeeName: employee.Name,
		}

		salary, err := s.CalculateEmployeeSalary(employee.ID, periodID, userID)
		if err != nil {
			item.Success = false
			item.Message = err.Error()
			result.Failed++
		} else {
			item.Success = true
			item.Message = "计算成功"
			item.Salary = salary
			item.GrossAmount = salary.GrossSalary
			item.NetAmount = salary.NetSalary
			result.Success++
			result.TotalAmount += salary.NetSalary
		}

		result.Results = append(result.Results, item)
	}

	return result, nil
}

func (s *SalaryService) GetEnhancedSalaries(params EnhancedSalaryQueryParams) (*utils.PaginationResponse, error) {
	var salaries []models.EnhancedSalary
	var total int64

	query := s.db.Model(&models.EnhancedSalary{}).
		Preload("Employee").
		Preload("Employee.Department").
		Preload("PayrollPeriod")

	if params.EmployeeID != nil {
		query = query.Where("employee_id = ?", *params.EmployeeID)
	}
	if params.PayrollPeriodID != nil {
		query = query.Where("payroll_period_id = ?", *params.PayrollPeriodID)
	}
	if params.DepartmentID != nil {
		query = query.Joins("JOIN employees ON enhanced_salaries.employee_id = employees.id").
			Where("employees.department_id = ?", *params.DepartmentID)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&salaries).Error; err != nil {
		return nil, err
	}

	response := utils.CreatePaginationResponse(salaries, params.Page, params.PageSize, total)
	return &response, nil
}

func (s *SalaryService) GetEnhancedSalaryByID(id uint) (*models.EnhancedSalary, error) {
	var salary models.EnhancedSalary
	if err := s.db.Preload("Employee").Preload("PayrollPeriod").First(&salary, id).Error; err != nil {
		return nil, err
	}
	return &salary, nil
}

func (s *SalaryService) UpdateSalaryDetails(salaryID uint, details []SalaryDetailUpdate, userID uint) (*models.EnhancedSalary, error) {
	var salary models.EnhancedSalary
	if err := s.db.First(&salary, salaryID).Error; err != nil {
		return nil, errors.New("salary record not found")
	}

	if salary.Status != "calculated" && salary.Status != "reviewed" {
		return nil, errors.New("salary cannot be modified in current status")
	}

	// Update salary details (simplified implementation)
	salary.Status = "modified"
	if err := s.db.Save(&salary).Error; err != nil {
		return nil, err
	}

	return &salary, nil
}

// ========================= Enhanced Approval Workflow =========================

func (s *SalaryService) ReviewSalary(id uint, reviewerID uint, notes string, approve bool) (*models.EnhancedSalary, error) {
	var salary models.EnhancedSalary
	if err := s.db.First(&salary, id).Error; err != nil {
		return nil, errors.New("salary record not found")
	}

	if salary.Status != "calculated" {
		return nil, errors.New("salary is not in calculated status")
	}

	if approve {
		salary.Status = "reviewed"
	} else {
		salary.Status = "rejected"
	}

	salary.ReviewedBy = &reviewerID
	now := time.Now()
	salary.ReviewedAt = &now

	if err := s.db.Save(&salary).Error; err != nil {
		return nil, err
	}

	return &salary, nil
}

func (s *SalaryService) ApproveEnhancedSalary(id uint, approverID uint, notes string) (*models.EnhancedSalary, error) {
	var salary models.EnhancedSalary
	if err := s.db.First(&salary, id).Error; err != nil {
		return nil, errors.New("salary record not found")
	}

	if salary.Status != "reviewed" {
		return nil, errors.New("salary must be reviewed before approval")
	}

	salary.Status = "approved"
	salary.ApprovedBy = &approverID
	now := time.Now()
	salary.ApprovedAt = &now

	if err := s.db.Save(&salary).Error; err != nil {
		return nil, err
	}

	return &salary, nil
}

func (s *SalaryService) RejectSalary(id uint, approverID uint, reason string) (*models.EnhancedSalary, error) {
	var salary models.EnhancedSalary
	if err := s.db.First(&salary, id).Error; err != nil {
		return nil, errors.New("salary record not found")
	}

	salary.Status = "rejected"
	salary.ApprovedBy = &approverID
	now := time.Now()
	salary.ApprovedAt = &now

	if err := s.db.Save(&salary).Error; err != nil {
		return nil, err
	}

	return &salary, nil
}

func (s *SalaryService) BulkApproveSalaries(salaryIDs []uint, approverID uint, notes string) (*BulkApprovalResult, error) {
	result := &BulkApprovalResult{
		Total:   len(salaryIDs),
		Results: make([]ApprovalItem, 0, len(salaryIDs)),
	}

	for _, salaryID := range salaryIDs {
		item := ApprovalItem{SalaryID: salaryID}

		_, err := s.ApproveEnhancedSalary(salaryID, approverID, notes)
		if err != nil {
			item.Success = false
			item.Message = err.Error()
			result.Failed++
		} else {
			item.Success = true
			item.Message = "批准成功"
			result.Success++
		}

		result.Results = append(result.Results, item)
	}

	return result, nil
}

// ========================= Enhanced Payment Processing =========================

func (s *SalaryService) CreatePaymentBatch(batch *models.PaymentBatch, salaryIDs []uint, userID uint) (*models.PaymentBatch, error) {
	batch.CreatedBy = &userID
	batch.Status = "pending"

	if err := s.db.Create(batch).Error; err != nil {
		return nil, err
	}

	// Create payroll records for each salary
	for _, salaryID := range salaryIDs {
		var salary models.EnhancedSalary
		if err := s.db.First(&salary, salaryID).Error; err != nil {
			continue
		}

		payrollRecord := models.EnhancedPayrollRecord{
			SalaryID:      salaryID,
			PaymentAmount: salary.NetSalary,
			Status:        "pending",
		}

		s.db.Create(&payrollRecord)
	}

	return batch, nil
}

func (s *SalaryService) ProcessPaymentBatch(batchID uint, userID uint) (*models.PaymentBatch, error) {
	var batch models.PaymentBatch
	if err := s.db.First(&batch, batchID).Error; err != nil {
		return nil, errors.New("payment batch not found")
	}

	batch.Status = "processed"
	batch.ProcessedBy = &userID

	if err := s.db.Save(&batch).Error; err != nil {
		return nil, err
	}

	// Update all payroll records in this batch
	s.db.Model(&models.EnhancedPayrollRecord{}).Where("batch_id = ?", batchID).Update("status", "paid")

	// Update salary statuses
	s.db.Model(&models.EnhancedSalary{}).
		Joins("JOIN enhanced_payroll_records ON enhanced_salaries.id = enhanced_payroll_records.salary_id").
		Where("enhanced_payroll_records.batch_id = ?", batchID).
		Update("status", "paid")

	return &batch, nil
}

func (s *SalaryService) GetPaymentBatches(params BatchQueryParams) (*utils.PaginationResponse, error) {
	var batches []models.PaymentBatch
	var total int64

	query := s.db.Model(&models.PaymentBatch{}).Preload("Creator")

	if params.PayrollPeriodID != nil {
		query = query.Where("payroll_period_id = ?", *params.PayrollPeriodID)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.CreatedBy != nil {
		query = query.Where("created_by = ?", *params.CreatedBy)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&batches).Error; err != nil {
		return nil, err
	}

	response := utils.CreatePaginationResponse(batches, params.Page, params.PageSize, total)
	return &response, nil
}

func (s *SalaryService) GetPaymentBatchByID(id uint) (*models.PaymentBatch, error) {
	var batch models.PaymentBatch
	if err := s.db.Preload("Creator").First(&batch, id).Error; err != nil {
		return nil, err
	}
	return &batch, nil
}

// ========================= Enhanced Salary Adjustments =========================

func (s *SalaryService) CreateSalaryAdjustment(adjustment *models.SalaryAdjustment) (*models.SalaryAdjustment, error) {
	if err := s.db.Create(adjustment).Error; err != nil {
		return nil, err
	}
	return adjustment, nil
}

func (s *SalaryService) ApproveSalaryAdjustment(id uint, approverID uint) (*models.SalaryAdjustment, error) {
	var adjustment models.SalaryAdjustment
	if err := s.db.First(&adjustment, id).Error; err != nil {
		return nil, errors.New("salary adjustment not found")
	}

	adjustment.Status = "approved"
	adjustment.ApprovedBy = &approverID
	now := time.Now()
	adjustment.ApprovedAt = &now

	if err := s.db.Save(&adjustment).Error; err != nil {
		return nil, err
	}

	return &adjustment, nil
}

func (s *SalaryService) GetSalaryAdjustments(params AdjustmentQueryParams) (*utils.PaginationResponse, error) {
	var adjustments []models.SalaryAdjustment
	var total int64

	query := s.db.Model(&models.SalaryAdjustment{}).Preload("Employee")

	if params.EmployeeID != nil {
		query = query.Where("employee_id = ?", *params.EmployeeID)
	}
	if params.AdjustmentType != "" {
		query = query.Where("adjustment_type = ?", params.AdjustmentType)
	}
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&adjustments).Error; err != nil {
		return nil, err
	}

	response := utils.CreatePaginationResponse(adjustments, params.Page, params.PageSize, total)
	return &response, nil
}

// ========================= Enhanced Analytics and Reporting =========================

func (s *SalaryService) GetSalaryAnalytics(params AnalyticsParams) (*SalaryAnalytics, error) {
	analytics := &SalaryAnalytics{}

	query := s.db.Model(&models.EnhancedSalary{}).
		Joins("JOIN employees ON enhanced_salaries.employee_id = employees.id")

	if params.DepartmentID != nil {
		query = query.Where("employees.department_id = ?", *params.DepartmentID)
	}
	if params.PeriodID != nil {
		query = query.Where("enhanced_salaries.payroll_period_id = ?", *params.PeriodID)
	}

	var result struct {
		Count      int     `gorm:"column:count"`
		TotalGross float64 `gorm:"column:total_gross"`
		TotalNet   float64 `gorm:"column:total_net"`
		AvgGross   float64 `gorm:"column:avg_gross"`
		AvgNet     float64 `gorm:"column:avg_net"`
	}

	query.Select("COUNT(*) as count, SUM(gross_salary) as total_gross, SUM(net_salary) as total_net, AVG(gross_salary) as avg_gross, AVG(net_salary) as avg_net").
		Scan(&result)

	analytics.TotalEmployees = result.Count
	analytics.TotalGrossAmount = result.TotalGross
	analytics.TotalNetAmount = result.TotalNet
	analytics.AverageGross = result.AvgGross
	analytics.AverageNet = result.AvgNet

	return analytics, nil
}

func (s *SalaryService) GetDepartmentSalaryReport(departmentID uint, periodID uint) (*DepartmentSalaryReport, error) {
	var department models.Department
	if err := s.db.First(&department, departmentID).Error; err != nil {
		return nil, errors.New("department not found")
	}

	var period models.PayrollPeriod
	if err := s.db.First(&period, periodID).Error; err != nil {
		return nil, errors.New("payroll period not found")
	}

	// Get analytics for this department
	params := AnalyticsParams{
		DepartmentID: &departmentID,
		PeriodID:     &periodID,
	}

	analytics, err := s.GetSalaryAnalytics(params)
	if err != nil {
		return nil, err
	}

	report := &DepartmentSalaryReport{
		Department:    &department,
		PayrollPeriod: &period,
		Summary: DepartmentAnalytics{
			DepartmentID:   departmentID,
			DepartmentName: department.Name,
			EmployeeCount:  analytics.TotalEmployees,
			TotalGross:     analytics.TotalGrossAmount,
			TotalNet:       analytics.TotalNetAmount,
			AverageGross:   analytics.AverageGross,
			AverageNet:     analytics.AverageNet,
		},
		GeneratedAt: time.Now(),
	}

	return report, nil
}

func (s *SalaryService) ExportEnhancedSalaryReport(params ExportParams) ([]byte, string, error) {
	// Simplified implementation
	filename := fmt.Sprintf("salary_report_%s_%d.%s", time.Now().Format("20060102"), time.Now().Unix(), params.Format)
	return []byte("Mock export data"), filename, nil
}

// ========================= Personal Salary Management =========================

func (s *SalaryService) GetPersonalSalaryDetail(params PersonalSalaryParams) (*PersonalSalaryDetail, error) {
	var salary models.EnhancedSalary
	query := s.db.Preload("Employee").Preload("PayrollPeriod")

	query = query.Where("employee_id = ?", params.EmployeeID)

	if params.PeriodID != nil {
		query = query.Where("payroll_period_id = ?", *params.PeriodID)
	}

	if err := query.Order("created_at DESC").First(&salary).Error; err != nil {
		return nil, err
	}

	return &PersonalSalaryDetail{
		Salary:   &salary,
		Period:   salary.PayrollPeriod,
		Employee: salary.Employee,
	}, nil
}

func (s *SalaryService) GetSalaryHistory(params SalaryHistoryParams) (*SalaryHistoryResponse, error) {
	var salaries []models.EnhancedSalary
	query := s.db.Preload("PayrollPeriod").Where("employee_id = ?", params.EmployeeID)

	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}

	if err := query.Order("created_at DESC").Find(&salaries).Error; err != nil {
		return nil, err
	}

	var historyItems []SalaryHistoryItem
	var totalNet float64
	var maxNet, minNet float64

	for i, salary := range salaries {
		item := SalaryHistoryItem{
			Period:    salary.PayrollPeriod,
			Salary:    &salary,
			NetSalary: salary.NetSalary,
		}

		if salary.PayrollPeriod != nil {
			item.Year = salary.PayrollPeriod.Year
			item.Month = fmt.Sprintf("%04d-%02d", salary.PayrollPeriod.Year, *salary.PayrollPeriod.Month)
		}

		historyItems = append(historyItems, item)
		totalNet += salary.NetSalary

		if i == 0 || salary.NetSalary > maxNet {
			maxNet = salary.NetSalary
		}
		if i == 0 || salary.NetSalary < minNet {
			minNet = salary.NetSalary
		}
	}

	summary := SalaryHistorySummary{
		TotalRecords: len(salaries),
		HighestNet:   maxNet,
		LowestNet:    minNet,
	}

	if len(salaries) > 0 {
		summary.AverageNet = totalNet / float64(len(salaries))
		summary.LatestNet = salaries[0].NetSalary
		summary.TrendDirection = "stable"
	}

	return &SalaryHistoryResponse{
		History: historyItems,
		Summary: summary,
	}, nil
}

func (s *SalaryService) GetPersonalPayrollRecords(params PersonalPayrollParams) (*utils.PaginationResponse, error) {
	var records []models.EnhancedPayrollRecord
	var total int64

	countQuery := s.db.Model(&models.EnhancedPayrollRecord{}).
		Joins("JOIN enhanced_salaries ON enhanced_payroll_records.salary_id = enhanced_salaries.id").
		Where("enhanced_salaries.employee_id = ?", params.EmployeeID)

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	query := s.db.Preload("Salary.PayrollPeriod").
		Joins("JOIN enhanced_salaries ON enhanced_payroll_records.salary_id = enhanced_salaries.id").
		Where("enhanced_salaries.employee_id = ?", params.EmployeeID).
		Order("enhanced_payroll_records.created_at DESC").
		Offset(offset).Limit(params.PageSize)

	if err := query.Find(&records).Error; err != nil {
		return nil, err
	}

	return &utils.PaginationResponse{
		Data:       records,
		TotalItems: total,
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}

func (s *SalaryService) GetPersonalSalaryDashboard(employeeID uint) (*PersonalSalaryDashboard, error) {
	currentParams := PersonalSalaryParams{
		EmployeeID: employeeID,
	}
	currentSalary, err := s.GetPersonalSalaryDetail(currentParams)
	if err != nil {
		currentSalary = nil
	}

	historyParams := SalaryHistoryParams{
		EmployeeID: employeeID,
		Limit:      6,
	}
	historyResponse, err := s.GetSalaryHistory(historyParams)
	if err != nil {
		return nil, err
	}

	var pendingPayrolls []models.EnhancedPayrollRecord
	s.db.Preload("Salary.PayrollPeriod").
		Joins("JOIN enhanced_salaries ON enhanced_payroll_records.salary_id = enhanced_salaries.id").
		Where("enhanced_salaries.employee_id = ? AND enhanced_payroll_records.status IN ?",
			employeeID, []string{"pending", "processing"}).
		Find(&pendingPayrolls)

	currentYear := time.Now().Year()
	var ytdEarnings float64
	s.db.Model(&models.EnhancedSalary{}).
		Joins("JOIN payroll_periods ON enhanced_salaries.payroll_period_id = payroll_periods.id").
		Where("enhanced_salaries.employee_id = ? AND payroll_periods.year = ? AND enhanced_salaries.status = ?",
			employeeID, currentYear, "paid").
		Select("COALESCE(SUM(net_salary), 0)").Scan(&ytdEarnings)

	return &PersonalSalaryDashboard{
		CurrentSalary:      currentSalary,
		RecentHistory:      historyResponse.History,
		Statistics:         historyResponse.Summary,
		PendingPayrolls:    pendingPayrolls,
		YearToDateEarnings: ytdEarnings,
	}, nil
}

// ========================= Formula Engine =========================

func (s *SalaryService) EvaluateFormula(formula string, context FormulaContext) (float64, error) {
	// Basic formula evaluation implementation
	if formula == "" {
		return 0, errors.New("formula cannot be empty")
	}

	// Simple implementation for demo
	if formula == "base_salary * 0.1" {
		return context.BaseSalary * 0.1, nil
	}

	return 0, errors.New("formula evaluation not implemented")
}

func (s *SalaryService) ValidateFormula(formula string) error {
	if strings.TrimSpace(formula) == "" {
		return errors.New("formula cannot be empty")
	}
	return nil
}