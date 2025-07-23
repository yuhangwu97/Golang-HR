package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gin-project/models"
	"gin-project/utils"

	"gorm.io/gorm"
)

type EnhancedSalaryServiceInterface interface {
	// Salary Component Management
	CreateSalaryComponent(component *models.SalaryComponent) (*models.SalaryComponent, error)
	UpdateSalaryComponent(id uint, component *models.SalaryComponent) (*models.SalaryComponent, error)
	DeleteSalaryComponent(id uint) error
	GetSalaryComponents(params ComponentQueryParams) (*utils.PaginationResponse, error)
	GetSalaryComponentByID(id uint) (*models.SalaryComponent, error)

	// Salary Grade Management
	CreateSalaryGrade(grade *models.SalaryGrade) (*models.SalaryGrade, error)
	UpdateSalaryGrade(id uint, grade *models.SalaryGrade) (*models.SalaryGrade, error)
	DeleteSalaryGrade(id uint) error
	GetSalaryGrades(params GradeQueryParams) (*utils.PaginationResponse, error)
	GetSalaryGradeByID(id uint) (*models.SalaryGrade, error)

	// Salary Structure Management
	CreateSalaryStructure(structure *models.SalaryStructure) (*models.SalaryStructure, error)
	UpdateSalaryStructure(id uint, structure *models.SalaryStructure) (*models.SalaryStructure, error)
	DeleteSalaryStructure(id uint) error
	GetSalaryStructures(params StructureQueryParams) (*utils.PaginationResponse, error)
	GetSalaryStructureByID(id uint) (*models.SalaryStructure, error)
	GetApplicableStructure(employeeID uint) (*models.SalaryStructure, error)

	// Payroll Period Management
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

	// Approval Workflow
	ReviewSalary(id uint, reviewerID uint, notes string, approve bool) (*models.EnhancedSalary, error)
	ApproveSalary(id uint, approverID uint, notes string) (*models.EnhancedSalary, error)
	RejectSalary(id uint, approverID uint, reason string) (*models.EnhancedSalary, error)
	BulkApproveSalaries(salaryIDs []uint, approverID uint, notes string) (*BulkApprovalResult, error)

	// Payment Processing
	CreatePaymentBatch(batch *models.PaymentBatch, salaryIDs []uint, userID uint) (*models.PaymentBatch, error)
	ProcessPaymentBatch(batchID uint, userID uint) (*models.PaymentBatch, error)
	GetPaymentBatches(params BatchQueryParams) (*utils.PaginationResponse, error)
	GetPaymentBatchByID(id uint) (*models.PaymentBatch, error)

	// Salary Adjustments
	CreateSalaryAdjustment(adjustment *models.SalaryAdjustment) (*models.SalaryAdjustment, error)
	ApproveSalaryAdjustment(id uint, approverID uint) (*models.SalaryAdjustment, error)
	GetSalaryAdjustments(params AdjustmentQueryParams) (*utils.PaginationResponse, error)

	// Analytics and Reporting
	GetSalaryAnalytics(params AnalyticsParams) (*SalaryAnalytics, error)
	GetDepartmentSalaryReport(departmentID uint, periodID uint) (*DepartmentSalaryReport, error)
	ExportSalaryReport(params ExportParams) ([]byte, string, error)

	// Personal Salary Management
	GetPersonalSalaryDetail(params PersonalSalaryParams) (*PersonalSalaryDetail, error)
	GetSalaryHistory(params SalaryHistoryParams) (*SalaryHistoryResponse, error)
	GetPersonalPayrollRecords(params PersonalPayrollParams) (*utils.PaginationResponse, error)
	GetPersonalSalaryDashboard(employeeID uint) (*PersonalSalaryDashboard, error)

	// Formula Engine
	EvaluateFormula(formula string, context FormulaContext) (float64, error)
	ValidateFormula(formula string) error
}

type EnhancedSalaryService struct {
	db *gorm.DB
}

// Query Parameters
type ComponentQueryParams struct {
	Category string
	Type     string
	Status   string
	Page     int
	PageSize int
}

type GradeQueryParams struct {
	Level    int
	Currency string
	Status   string
	Page     int
	PageSize int
}

type StructureQueryParams struct {
	DepartmentID *uint
	PositionID   *uint
	JobLevelID   *uint
	Status       string
	Page         int
	PageSize     int
}

type PeriodQueryParams struct {
	PeriodType string
	Year       int
	Status     string
	Page       int
	PageSize   int
}

type EnhancedSalaryQueryParams struct {
	EmployeeID      *uint
	PayrollPeriodID *uint
	DepartmentID    *uint
	Status          string
	Page            int
	PageSize        int
}

type BatchQueryParams struct {
	PayrollPeriodID *uint
	Status          string
	CreatedBy       *uint
	Page            int
	PageSize        int
}

type AdjustmentQueryParams struct {
	EmployeeID       *uint
	AdjustmentType   string
	Status           string
	EffectiveDateFrom *time.Time
	EffectiveDateTo   *time.Time
	Page             int
	PageSize         int
}

type AnalyticsParams struct {
	DepartmentID *uint
	PeriodID     *uint
	Year         *int
	Month        *int
}

type ExportParams struct {
	PeriodID     *uint
	DepartmentID *uint
	Format       string // excel, csv, pdf
	Template     string
}

// Personal Salary Parameters
type PersonalSalaryParams struct {
	EmployeeID uint
	PeriodID   *uint
	Year       *int
	Month      *int
}

type SalaryHistoryParams struct {
	EmployeeID uint
	Limit      int
}

type PersonalPayrollParams struct {
	EmployeeID uint
	Page       int
	PageSize   int
}

// Result Types
type EnhancedBatchResult struct {
	Total      int                        `json:"total"`
	Success    int                        `json:"success"`
	Failed     int                        `json:"failed"`
	Results    []EnhancedCalculateItem    `json:"results"`
	TotalAmount float64                   `json:"total_amount"`
}

type EnhancedCalculateItem struct {
	EmployeeID   uint                  `json:"employee_id"`
	EmployeeName string                `json:"employee_name"`
	Success      bool                  `json:"success"`
	Message      string                `json:"message"`
	Salary       *models.EnhancedSalary `json:"salary,omitempty"`
	GrossAmount  float64               `json:"gross_amount"`
	NetAmount    float64               `json:"net_amount"`
}

type BulkApprovalResult struct {
	Total     int    `json:"total"`
	Success   int    `json:"success"`
	Failed    int    `json:"failed"`
	Results   []ApprovalItem `json:"results"`
}

type ApprovalItem struct {
	SalaryID uint   `json:"salary_id"`
	Success  bool   `json:"success"`
	Message  string `json:"message"`
}

type SalaryDetailUpdate struct {
	ComponentID   uint     `json:"component_id"`
	ManualValue   *float64 `json:"manual_value"`
	Notes         string   `json:"notes"`
}

type SalaryAnalytics struct {
	TotalEmployees    int                    `json:"total_employees"`
	TotalGrossAmount  float64                `json:"total_gross_amount"`
	TotalNetAmount    float64                `json:"total_net_amount"`
	AverageGross      float64                `json:"average_gross"`
	AverageNet        float64                `json:"average_net"`
	MedianGross       float64                `json:"median_gross"`
	MedianNet         float64                `json:"median_net"`
	DepartmentBreakdown []DepartmentAnalytics `json:"department_breakdown"`
	GradeBreakdown    []GradeAnalytics       `json:"grade_breakdown"`
	ComponentBreakdown []ComponentAnalytics  `json:"component_breakdown"`
}

type DepartmentAnalytics struct {
	DepartmentID   uint    `json:"department_id"`
	DepartmentName string  `json:"department_name"`
	EmployeeCount  int     `json:"employee_count"`
	TotalGross     float64 `json:"total_gross"`
	TotalNet       float64 `json:"total_net"`
	AverageGross   float64 `json:"average_gross"`
	AverageNet     float64 `json:"average_net"`
}

type GradeAnalytics struct {
	GradeID     uint    `json:"grade_id"`
	GradeName   string  `json:"grade_name"`
	EmployeeCount int   `json:"employee_count"`
	AverageGross  float64 `json:"average_gross"`
	AverageNet    float64 `json:"average_net"`
}

type ComponentAnalytics struct {
	ComponentID   uint    `json:"component_id"`
	ComponentName string  `json:"component_name"`
	TotalAmount   float64 `json:"total_amount"`
	AverageAmount float64 `json:"average_amount"`
}

type DepartmentSalaryReport struct {
	Department      *models.Department    `json:"department"`
	PayrollPeriod   *models.PayrollPeriod `json:"payroll_period"`
	Summary         DepartmentAnalytics   `json:"summary"`
	Employees       []EmployeeSalaryInfo  `json:"employees"`
	GeneratedAt     time.Time             `json:"generated_at"`
}

type EmployeeSalaryInfo struct {
	Employee      *models.Employee       `json:"employee"`
	Salary        *models.EnhancedSalary `json:"salary"`
	Components    []ComponentBreakdown   `json:"components"`
}

type ComponentBreakdown struct {
	Component *models.SalaryComponent `json:"component"`
	Amount    float64                 `json:"amount"`
}

// Personal Salary Response Types
type PersonalSalaryDetail struct {
	Salary     *models.EnhancedSalary `json:"salary"`
	Components []ComponentBreakdown   `json:"components"`
	Period     *models.PayrollPeriod  `json:"period"`
	Employee   *models.Employee       `json:"employee"`
}

type SalaryHistoryResponse struct {
	History []SalaryHistoryItem `json:"history"`
	Summary SalaryHistorySummary `json:"summary"`
}

type SalaryHistoryItem struct {
	Period    *models.PayrollPeriod  `json:"period"`
	Salary    *models.EnhancedSalary `json:"salary"`
	NetSalary float64                `json:"net_salary"`
	Month     string                 `json:"month"`
	Year      int                    `json:"year"`
}

type SalaryHistorySummary struct {
	TotalRecords   int     `json:"total_records"`
	AverageNet     float64 `json:"average_net"`
	HighestNet     float64 `json:"highest_net"`
	LowestNet      float64 `json:"lowest_net"`
	LatestNet      float64 `json:"latest_net"`
	TrendDirection string  `json:"trend_direction"` // up, down, stable
}

type PersonalSalaryDashboard struct {
	CurrentSalary   *PersonalSalaryDetail `json:"current_salary"`
	RecentHistory   []SalaryHistoryItem   `json:"recent_history"`
	Statistics      SalaryHistorySummary  `json:"statistics"`
	PendingPayrolls []models.EnhancedPayrollRecord `json:"pending_payrolls"`
	NextPayDate     *time.Time            `json:"next_pay_date"`
	YearToDateEarnings float64           `json:"year_to_date_earnings"`
}

type FormulaContext struct {
	Employee      *models.Employee      `json:"employee"`
	BaseSalary    float64               `json:"base_salary"`
	Components    map[string]float64    `json:"components"`
	Variables     map[string]interface{} `json:"variables"`
}

func NewEnhancedSalaryService(db *gorm.DB) EnhancedSalaryServiceInterface {
	return &EnhancedSalaryService{
		db: db,
	}
}

// ========================= Salary Component Management =========================

func (s *EnhancedSalaryService) CreateSalaryComponent(component *models.SalaryComponent) (*models.SalaryComponent, error) {
	if err := s.validateSalaryComponent(component); err != nil {
		return nil, err
	}

	if err := s.db.Create(component).Error; err != nil {
		return nil, fmt.Errorf("failed to create salary component: %w", err)
	}

	return component, nil
}

func (s *EnhancedSalaryService) UpdateSalaryComponent(id uint, component *models.SalaryComponent) (*models.SalaryComponent, error) {
	if err := s.validateSalaryComponent(component); err != nil {
		return nil, err
	}

	component.ID = id
	if err := s.db.Save(component).Error; err != nil {
		return nil, fmt.Errorf("failed to update salary component: %w", err)
	}

	return component, nil
}

func (s *EnhancedSalaryService) DeleteSalaryComponent(id uint) error {
	// Check if component is being used in any structures
	var count int64
	if err := s.db.Model(&models.SalaryStructureComponent{}).Where("component_id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("salary component is being used in salary structures and cannot be deleted")
	}

	return s.db.Delete(&models.SalaryComponent{}, id).Error
}

func (s *EnhancedSalaryService) GetSalaryComponents(params ComponentQueryParams) (*utils.PaginationResponse, error) {
	var components []models.SalaryComponent
	var total int64

	query := s.db.Model(&models.SalaryComponent{})

	// Apply filters
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
	if err := query.Offset(offset).Limit(params.PageSize).Order("sort ASC, created_at DESC").Find(&components).Error; err != nil {
		return nil, err
	}

	response := utils.CreatePaginationResponse(components, params.Page, params.PageSize, total)
	return &response, nil
}

func (s *EnhancedSalaryService) GetSalaryComponentByID(id uint) (*models.SalaryComponent, error) {
	var component models.SalaryComponent
	if err := s.db.First(&component, id).Error; err != nil {
		return nil, err
	}
	return &component, nil
}

// ========================= Salary Calculation =========================

func (s *EnhancedSalaryService) CalculateEmployeeSalary(employeeID, periodID uint, userID uint) (*models.EnhancedSalary, error) {
	// Get employee and period
	var employee models.Employee
	if err := s.db.Preload("Department").Preload("Position").Preload("JobLevel").First(&employee, employeeID).Error; err != nil {
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

	// Get applicable salary structure
	structure, err := s.GetApplicableStructure(employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get salary structure: %w", err)
	}

	// Create new salary record
	salary := &models.EnhancedSalary{
		EmployeeID:      employeeID,
		PayrollPeriodID: periodID,
		StructureID:     &structure.ID,
		Status:          models.SalaryStatusDraft,
		CalculatedBy:    &userID,
		CalculatedAt:    &time.Time{},
		Version:         1,
	}

	now := time.Now()
	salary.CalculatedAt = &now

	// Calculate salary components
	totalGross := 0.0
	totalDeductions := 0.0
	var details []models.SalaryDetail

	for _, structComp := range structure.Components {
		if structComp.Component == nil {
			continue
		}

		component := structComp.Component
		calculatedValue, err := s.calculateComponentValue(component, &employee, structure, salary)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate component %s: %w", component.Name, err)
		}

		detail := models.SalaryDetail{
			ComponentID:    component.ID,
			Component:      component,
			CalculatedValue: calculatedValue,
			FinalValue:     calculatedValue,
		}

		// Apply manual override if exists
		if structComp.DefaultValue > 0 {
			detail.ManualValue = &structComp.DefaultValue
			detail.FinalValue = structComp.DefaultValue
		}

		details = append(details, detail)

		// Add to totals based on category
		switch component.Category {
		case models.ComponentCategoryDeduction, models.ComponentCategoryTax:
			totalDeductions += detail.FinalValue
		default:
			totalGross += detail.FinalValue
		}
	}

	salary.GrossSalary = totalGross
	salary.TotalDeductions = totalDeductions
	salary.NetSalary = totalGross - totalDeductions
	salary.Status = models.SalaryStatusCalculated

	// Save salary record
	if err := s.db.Create(salary).Error; err != nil {
		return nil, fmt.Errorf("failed to create salary record: %w", err)
	}

	// Save salary details
	for i := range details {
		details[i].SalaryID = salary.ID
	}
	if err := s.db.Create(&details).Error; err != nil {
		return nil, fmt.Errorf("failed to create salary details: %w", err)
	}

	// Load complete salary record with relationships
	if err := s.db.Preload("Employee").Preload("PayrollPeriod").Preload("Structure").
		Preload("Components.Component").First(salary, salary.ID).Error; err != nil {
		return nil, err
	}

	return salary, nil
}

func (s *EnhancedSalaryService) GetApplicableStructure(employeeID uint) (*models.SalaryStructure, error) {
	var employee models.Employee
	if err := s.db.Preload("Department").Preload("Position").Preload("JobLevel").First(&employee, employeeID).Error; err != nil {
		return nil, err
	}

	var structure models.SalaryStructure
	query := s.db.Preload("Components.Component").Where("status = ?", "active")

	// Try to find exact match first
	if employee.JobLevelID > 0 {
		err := query.Where("job_level_id = ?", employee.JobLevelID).First(&structure).Error
		if err == nil {
			return &structure, nil
		}
	}

	if employee.PositionID > 0 {
		err := query.Where("position_id = ?", employee.PositionID).First(&structure).Error
		if err == nil {
			return &structure, nil
		}
	}

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

func (s *EnhancedSalaryService) calculateComponentValue(component *models.SalaryComponent, employee *models.Employee, structure *models.SalaryStructure, salary *models.EnhancedSalary) (float64, error) {
	switch component.Type {
	case models.ComponentTypeFixed:
		return component.DefaultAmount, nil

	case models.ComponentTypePercentage:
		// Extract percentage from formula (e.g., "15%" -> 0.15)
		percentage, err := s.parsePercentage(component.Formula)
		if err != nil {
			return 0, err
		}
		return employee.BaseSalary * percentage, nil

	case models.ComponentTypeFormula:
		context := FormulaContext{
			Employee:   employee,
			BaseSalary: employee.BaseSalary,
			Components: make(map[string]float64),
			Variables:  make(map[string]interface{}),
		}
		
		// Add common variables
		context.Variables["base_salary"] = employee.BaseSalary
		context.Variables["employee_id"] = employee.ID
		context.Variables["department_id"] = employee.DepartmentID
		
		return s.EvaluateFormula(component.Formula, context)

	case models.ComponentTypeManual:
		return component.DefaultAmount, nil

	default:
		return component.DefaultAmount, nil
	}
}

func (s *EnhancedSalaryService) parsePercentage(formula string) (float64, error) {
	formula = strings.TrimSpace(formula)
	if strings.HasSuffix(formula, "%") {
		percentStr := strings.TrimSuffix(formula, "%")
		percent, err := strconv.ParseFloat(percentStr, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid percentage format: %s", formula)
		}
		return percent / 100.0, nil
	}
	return 0, fmt.Errorf("invalid percentage format: %s", formula)
}

// ========================= Formula Engine =========================

func (s *EnhancedSalaryService) EvaluateFormula(formula string, context FormulaContext) (float64, error) {
	// Basic formula evaluation - can be enhanced with a proper expression evaluator
	formula = strings.TrimSpace(formula)
	
	// Replace variables in formula
	for key, value := range context.Variables {
		placeholder := fmt.Sprintf("{%s}", key)
		var valueStr string
		switch v := value.(type) {
		case float64:
			valueStr = fmt.Sprintf("%.2f", v)
		case int:
			valueStr = fmt.Sprintf("%d", v)
		case uint:
			valueStr = fmt.Sprintf("%d", v)
		default:
			valueStr = fmt.Sprintf("%v", v)
		}
		formula = strings.ReplaceAll(formula, placeholder, valueStr)
	}

	// Simple expression evaluation
	// This is a basic implementation - in production, use a proper expression evaluator
	if result, err := s.evaluateSimpleExpression(formula); err == nil {
		return result, nil
	}

	return 0, fmt.Errorf("unable to evaluate formula: %s", formula)
}

func (s *EnhancedSalaryService) evaluateSimpleExpression(expr string) (float64, error) {
	// Basic implementation for simple expressions like "5000 * 0.08"
	expr = strings.ReplaceAll(expr, " ", "")
	
	// Handle simple multiplication
	if strings.Contains(expr, "*") {
		parts := strings.Split(expr, "*")
		if len(parts) == 2 {
			left, err1 := strconv.ParseFloat(parts[0], 64)
			right, err2 := strconv.ParseFloat(parts[1], 64)
			if err1 == nil && err2 == nil {
				return left * right, nil
			}
		}
	}
	
	// Handle simple addition
	if strings.Contains(expr, "+") {
		parts := strings.Split(expr, "+")
		total := 0.0
		for _, part := range parts {
			if val, err := strconv.ParseFloat(part, 64); err == nil {
				total += val
			}
		}
		return total, nil
	}
	
	// Try to parse as a number
	if val, err := strconv.ParseFloat(expr, 64); err == nil {
		return val, nil
	}
	
	return 0, fmt.Errorf("cannot evaluate expression: %s", expr)
}

func (s *EnhancedSalaryService) ValidateFormula(formula string) error {
	// Basic validation - can be enhanced
	if strings.TrimSpace(formula) == "" {
		return errors.New("formula cannot be empty")
	}
	return nil
}

// ========================= Validation =========================

func (s *EnhancedSalaryService) validateSalaryComponent(component *models.SalaryComponent) error {
	if component.Code == "" {
		return errors.New("salary component code is required")
	}
	if component.Name == "" {
		return errors.New("salary component name is required")
	}
	if component.Category == "" {
		return errors.New("salary component category is required")
	}
	if component.Type == "" {
		return errors.New("salary component type is required")
	}

	// Validate formula if type is formula or percentage
	if component.Type == models.ComponentTypeFormula || component.Type == models.ComponentTypePercentage {
		if err := s.ValidateFormula(component.Formula); err != nil {
			return fmt.Errorf("invalid formula: %w", err)
		}
	}

	// Validate amount ranges
	if component.MinAmount != nil && component.MaxAmount != nil {
		if *component.MinAmount > *component.MaxAmount {
			return errors.New("minimum amount cannot be greater than maximum amount")
		}
	}

	return nil
}

// Stub implementations for remaining methods - these would be fully implemented in production

func (s *EnhancedSalaryService) BatchCalculateSalaries(periodID uint, departmentID *uint, employeeIDs []uint, userID uint) (*EnhancedBatchResult, error) {
	// Implementation would batch calculate salaries for multiple employees
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) GetEnhancedSalaries(params EnhancedSalaryQueryParams) (*utils.PaginationResponse, error) {
	// Implementation would return paginated enhanced salary records
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) GetEnhancedSalaryByID(id uint) (*models.EnhancedSalary, error) {
	// Implementation would return a specific enhanced salary record
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) UpdateSalaryDetails(salaryID uint, details []SalaryDetailUpdate, userID uint) (*models.EnhancedSalary, error) {
	// Implementation would update salary detail components
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) ReviewSalary(id uint, reviewerID uint, notes string, approve bool) (*models.EnhancedSalary, error) {
	// Implementation would handle salary review process
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) ApproveSalary(id uint, approverID uint, notes string) (*models.EnhancedSalary, error) {
	// Implementation would handle salary approval
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) RejectSalary(id uint, approverID uint, reason string) (*models.EnhancedSalary, error) {
	// Implementation would handle salary rejection
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) BulkApproveSalaries(salaryIDs []uint, approverID uint, notes string) (*BulkApprovalResult, error) {
	// Implementation would handle bulk salary approval
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) CreatePaymentBatch(batch *models.PaymentBatch, salaryIDs []uint, userID uint) (*models.PaymentBatch, error) {
	// Implementation would create payment batches
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) ProcessPaymentBatch(batchID uint, userID uint) (*models.PaymentBatch, error) {
	// Implementation would process payment batches
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) GetPaymentBatches(params BatchQueryParams) (*utils.PaginationResponse, error) {
	// Implementation would return payment batches
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) GetPaymentBatchByID(id uint) (*models.PaymentBatch, error) {
	// Implementation would return a specific payment batch
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) CreateSalaryAdjustment(adjustment *models.SalaryAdjustment) (*models.SalaryAdjustment, error) {
	// Implementation would create salary adjustments
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) ApproveSalaryAdjustment(id uint, approverID uint) (*models.SalaryAdjustment, error) {
	// Implementation would approve salary adjustments
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) GetSalaryAdjustments(params AdjustmentQueryParams) (*utils.PaginationResponse, error) {
	// Implementation would return salary adjustments
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) GetSalaryAnalytics(params AnalyticsParams) (*SalaryAnalytics, error) {
	// Implementation would return salary analytics
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) GetDepartmentSalaryReport(departmentID uint, periodID uint) (*DepartmentSalaryReport, error) {
	// Implementation would generate department salary reports
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) ExportSalaryReport(params ExportParams) ([]byte, string, error) {
	// Implementation would export salary reports
	return nil, "", errors.New("not implemented yet")
}

// Additional stub implementations for remaining methods

func (s *EnhancedSalaryService) CreateSalaryGrade(grade *models.SalaryGrade) (*models.SalaryGrade, error) {
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) UpdateSalaryGrade(id uint, grade *models.SalaryGrade) (*models.SalaryGrade, error) {
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) DeleteSalaryGrade(id uint) error {
	return errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) GetSalaryGrades(params GradeQueryParams) (*utils.PaginationResponse, error) {
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) GetSalaryGradeByID(id uint) (*models.SalaryGrade, error) {
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) CreateSalaryStructure(structure *models.SalaryStructure) (*models.SalaryStructure, error) {
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) UpdateSalaryStructure(id uint, structure *models.SalaryStructure) (*models.SalaryStructure, error) {
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) DeleteSalaryStructure(id uint) error {
	return errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) GetSalaryStructures(params StructureQueryParams) (*utils.PaginationResponse, error) {
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) GetSalaryStructureByID(id uint) (*models.SalaryStructure, error) {
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) CreatePayrollPeriod(period *models.PayrollPeriod) (*models.PayrollPeriod, error) {
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) UpdatePayrollPeriod(id uint, period *models.PayrollPeriod) (*models.PayrollPeriod, error) {
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) GetPayrollPeriods(params PeriodQueryParams) (*utils.PaginationResponse, error) {
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
	if err := query.Offset(offset).Limit(params.PageSize).
		Order("year DESC, period_type ASC, start_date DESC").
		Find(&periods).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &utils.PaginationResponse{
		Data:       periods,
		TotalItems: total,
		TotalPages: totalPages,
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}

func (s *EnhancedSalaryService) GetPayrollPeriodByID(id uint) (*models.PayrollPeriod, error) {
	return nil, errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) LockPayrollPeriod(id uint) error {
	return errors.New("not implemented yet")
}

func (s *EnhancedSalaryService) UnlockPayrollPeriod(id uint) error {
	return errors.New("not implemented yet")
}

// ========================= Personal Salary Management =========================

func (s *EnhancedSalaryService) GetPersonalSalaryDetail(params PersonalSalaryParams) (*PersonalSalaryDetail, error) {
	var salary models.EnhancedSalary
	query := s.db.Preload("Employee").Preload("PayrollPeriod").Preload("Components.Component")
	
	// Build query based on parameters
	query = query.Where("employee_id = ?", params.EmployeeID)
	
	if params.PeriodID != nil {
		query = query.Where("payroll_period_id = ?", *params.PeriodID)
	} else if params.Year != nil && params.Month != nil {
		// Find period by year/month
		var period models.PayrollPeriod
		err := s.db.Where("year = ? AND month = ?", *params.Year, *params.Month).First(&period).Error
		if err != nil {
			return nil, err
		}
		query = query.Where("payroll_period_id = ?", period.ID)
	}
	
	if err := query.Order("created_at DESC").First(&salary).Error; err != nil {
		return nil, err
	}
	
	// Build component breakdown
	var components []ComponentBreakdown
	for _, detail := range salary.Components {
		components = append(components, ComponentBreakdown{
			Component: detail.Component,
			Amount:    detail.FinalValue,
		})
	}
	
	return &PersonalSalaryDetail{
		Salary:     &salary,
		Components: components,
		Period:     salary.PayrollPeriod,
		Employee:   salary.Employee,
	}, nil
}

func (s *EnhancedSalaryService) GetSalaryHistory(params SalaryHistoryParams) (*SalaryHistoryResponse, error) {
	var salaries []models.EnhancedSalary
	query := s.db.Preload("PayrollPeriod").Where("employee_id = ?", params.EmployeeID)
	
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	
	if err := query.Order("created_at DESC").Find(&salaries).Error; err != nil {
		return nil, err
	}
	
	// Build history items
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
	
	// Calculate summary
	summary := SalaryHistorySummary{
		TotalRecords: len(salaries),
		HighestNet:   maxNet,
		LowestNet:    minNet,
	}
	
	if len(salaries) > 0 {
		summary.AverageNet = totalNet / float64(len(salaries))
		summary.LatestNet = salaries[0].NetSalary
		
		// Determine trend (simple: compare latest with previous)
		if len(salaries) > 1 {
			latest := salaries[0].NetSalary
			previous := salaries[1].NetSalary
			if latest > previous {
				summary.TrendDirection = "up"
			} else if latest < previous {
				summary.TrendDirection = "down"
			} else {
				summary.TrendDirection = "stable"
			}
		} else {
			summary.TrendDirection = "stable"
		}
	}
	
	return &SalaryHistoryResponse{
		History: historyItems,
		Summary: summary,
	}, nil
}

func (s *EnhancedSalaryService) GetPersonalPayrollRecords(params PersonalPayrollParams) (*utils.PaginationResponse, error) {
	var records []models.EnhancedPayrollRecord
	var total int64
	
	// Count total records
	countQuery := s.db.Model(&models.EnhancedPayrollRecord{}).
		Joins("JOIN enhanced_salaries ON enhanced_payroll_records.salary_id = enhanced_salaries.id").
		Where("enhanced_salaries.employee_id = ?", params.EmployeeID)
	
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, err
	}
	
	// Get paginated records
	offset := (params.Page - 1) * params.PageSize
	query := s.db.Preload("Salary.PayrollPeriod").Preload("PaymentBatch").
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

func (s *EnhancedSalaryService) GetPersonalSalaryDashboard(employeeID uint) (*PersonalSalaryDashboard, error) {
	// Get current salary
	currentParams := PersonalSalaryParams{
		EmployeeID: employeeID,
	}
	currentSalary, err := s.GetPersonalSalaryDetail(currentParams)
	if err != nil {
		// If no current salary, set to nil instead of returning error
		currentSalary = nil
	}
	
	// Get recent history (last 6 months)
	historyParams := SalaryHistoryParams{
		EmployeeID: employeeID,
		Limit:      6,
	}
	historyResponse, err := s.GetSalaryHistory(historyParams)
	if err != nil {
		return nil, err
	}
	
	// Get pending payroll records
	var pendingPayrolls []models.EnhancedPayrollRecord
	s.db.Preload("Salary.PayrollPeriod").
		Joins("JOIN enhanced_salaries ON enhanced_payroll_records.salary_id = enhanced_salaries.id").
		Where("enhanced_salaries.employee_id = ? AND enhanced_payroll_records.status IN ?", 
			employeeID, []string{"pending", "processing"}).
		Find(&pendingPayrolls)
	
	// Calculate year-to-date earnings
	currentYear := time.Now().Year()
	var ytdEarnings float64
	s.db.Model(&models.EnhancedSalary{}).
		Joins("JOIN payroll_periods ON enhanced_salaries.payroll_period_id = payroll_periods.id").
		Where("enhanced_salaries.employee_id = ? AND payroll_periods.year = ? AND enhanced_salaries.status = ?", 
			employeeID, currentYear, "paid").
		Select("COALESCE(SUM(net_salary), 0)").Scan(&ytdEarnings)
	
	// Find next pay date (next scheduled payroll period)
	var nextPayDate *time.Time
	var nextPeriod models.PayrollPeriod
	err = s.db.Where("start_date > ? AND status IN ?", time.Now(), []string{"open", "calculated"}).
		Order("start_date ASC").First(&nextPeriod).Error
	if err == nil && nextPeriod.PayDate != nil {
		nextPayDate = nextPeriod.PayDate
	}
	
	return &PersonalSalaryDashboard{
		CurrentSalary:      currentSalary,
		RecentHistory:      historyResponse.History,
		Statistics:         historyResponse.Summary,
		PendingPayrolls:    pendingPayrolls,
		NextPayDate:        nextPayDate,
		YearToDateEarnings: ytdEarnings,
	}, nil
}