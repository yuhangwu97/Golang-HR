package services

import (
	"encoding/csv"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	"gin-project/models"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type EmployeeServiceInterface interface {
	GetEmployeeList(params EmployeeListParams) (*EmployeeListResponse, error)
	GetEmployeeByID(id uint) (*models.Employee, error)
	CreateEmployee(employee *models.Employee) (*models.Employee, error)
	UpdateEmployee(employee *models.Employee) (*models.Employee, error)
	DeleteEmployee(id uint) error
	GetEmployeeStatistics() (*EmployeeStatistics, error)
	SearchEmployees(query string) ([]*models.Employee, error)
	ExportEmployees(format string) ([]byte, string, error)
	ImportEmployees(file *multipart.FileHeader) (*ImportResult, error)
	BulkUpdateEmployees(ids []uint, updates map[string]interface{}) (*BulkResult, error)
	BulkDeleteEmployees(ids []uint) (*BulkResult, error)
	GetEmployeesByDepartment(departmentID uint) ([]*models.Employee, error)
}

type EmployeeService struct {
	db *gorm.DB
}

func NewEmployeeService(db *gorm.DB) EmployeeServiceInterface {
	return &EmployeeService{
		db: db,
	}
}

// InjectDependencies implements DependencyInjector interface
func (es *EmployeeService) InjectDependencies(deps ...interface{}) error {
	for _, dep := range deps {
		switch d := dep.(type) {
		case *gorm.DB:
			es.db = d
		}
	}
	return nil
}

type EmployeeListParams struct {
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
	Department string `json:"department"`
	Status     string `json:"status"`
	Keyword    string `json:"keyword"`
	SortBy     string `json:"sort_by"`
	SortOrder  string `json:"sort_order"`
}

type EmployeeListResponse struct {
	Data       []*models.Employee `json:"data"`
	Pagination *Pagination        `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type EmployeeStatistics struct {
	Total            int64                    `json:"total"`
	ByDepartment     map[string]int64         `json:"by_department"`
	ByStatus         map[string]int64         `json:"by_status"`
	RecentHires      []*models.Employee       `json:"recent_hires"`
	AgeDistribution  map[string]int64         `json:"age_distribution"`
	MonthlyTrend     []MonthlyHireStats       `json:"monthly_trend"`
	DepartmentStats  []DepartmentStatItem     `json:"department_stats"`
}

type MonthlyHireStats struct {
	Month  string `json:"month"`
	Hires  int64  `json:"hires"`
	Leaves int64  `json:"leaves"`
}

type DepartmentStatItem struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	EmployeeCount int64  `json:"employee_count"`
}

type ImportResult struct {
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors"`
}

type BulkResult struct {
	Success int      `json:"success"`
	Failed  int      `json:"failed"`
	Errors  []string `json:"errors"`
}

func (es *EmployeeService) GetEmployeeList(params EmployeeListParams) (*EmployeeListResponse, error) {
	var employees []*models.Employee
	var total int64

	query := es.db.Model(&models.Employee{}).
		Preload("Department.Parent.Parent.Parent.Parent.Parent"). // 最多5层部门层级
		Preload("Position").
		Preload("JobLevel").
		Preload("Manager")

	// 过滤条件
	if params.Department != "" {
		query = query.Joins("JOIN departments ON employees.department_id = departments.id").
			Where("departments.name LIKE ?", "%"+params.Department+"%")
	}

	if params.Status != "" {
		query = query.Where("employees.status = ?", params.Status)
	}

	if params.Keyword != "" {
		query = query.Where("employees.name LIKE ? OR employees.email LIKE ? OR employees.employee_id LIKE ?",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 排序
	orderBy := params.SortBy
	if params.SortOrder == "desc" {
		orderBy += " DESC"
	}
	query = query.Order(orderBy)

	// 分页
	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Find(&employees).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &EmployeeListResponse{
		Data: employees,
		Pagination: &Pagination{
			Page:       params.Page,
			PageSize:   params.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (es *EmployeeService) GetEmployeeByID(id uint) (*models.Employee, error) {
	var employee models.Employee
	err := es.db.Preload("Department.Parent.Parent.Parent.Parent.Parent").
		Preload("Position").
		Preload("JobLevel").
		Preload("Manager").
		Preload("WorkExperience").
		First(&employee, id).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (es *EmployeeService) CreateEmployee(employee *models.Employee) (*models.Employee, error) {
	// 生成员工工号
	if employee.EmployeeID == "" {
		employee.EmployeeID = es.generateEmployeeID()
	}

	// 设置默认状态
	if employee.Status == "" {
		employee.Status = "active"
	}

	err := es.db.Create(employee).Error
	if err != nil {
		return nil, err
	}

	return es.GetEmployeeByID(employee.ID)
}

func (es *EmployeeService) UpdateEmployee(employee *models.Employee) (*models.Employee, error) {
	err := es.db.Model(employee).Updates(employee).Error
	if err != nil {
		return nil, err
	}

	return es.GetEmployeeByID(employee.ID)
}

func (es *EmployeeService) DeleteEmployee(id uint) error {
	return es.db.Delete(&models.Employee{}, id).Error
}

func (es *EmployeeService) GetEmployeeStatistics() (*EmployeeStatistics, error) {
	var total int64
	if err := es.db.Model(&models.Employee{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// 按部门统计
	var deptStats []struct {
		DepartmentName string `json:"department_name"`
		Count          int64  `json:"count"`
	}
	if err := es.db.Model(&models.Employee{}).
		Select("departments.name as department_name, COUNT(*) as count").
		Joins("LEFT JOIN departments ON employees.department_id = departments.id").
		Group("departments.name").
		Scan(&deptStats).Error; err != nil {
		return nil, err
	}

	byDepartment := make(map[string]int64)
	for _, stat := range deptStats {
		name := stat.DepartmentName
		if name == "" {
			name = "未分配部门"
		}
		byDepartment[name] = stat.Count
	}

	// 按状态统计
	var statusStats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	if err := es.db.Model(&models.Employee{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusStats).Error; err != nil {
		return nil, err
	}

	byStatus := make(map[string]int64)
	for _, stat := range statusStats {
		byStatus[stat.Status] = stat.Count
	}

	// 最近入职员工
	var recentHires []*models.Employee
	if err := es.db.Where("hire_date IS NOT NULL").
		Order("hire_date DESC").
		Limit(10).
		Preload("Department.Parent.Parent.Parent.Parent.Parent").
		Preload("Position").
		Find(&recentHires).Error; err != nil {
		return nil, err
	}

	// 年龄分布统计
	ageDistribution := es.calculateAgeDistribution()

	// 月度入职趋势（最近12个月）
	monthlyTrend := es.calculateMonthlyTrend()

	// 部门详细统计
	departmentStats := es.calculateDepartmentStats()

	return &EmployeeStatistics{
		Total:            total,
		ByDepartment:     byDepartment,
		ByStatus:         byStatus,
		RecentHires:      recentHires,
		AgeDistribution:  ageDistribution,
		MonthlyTrend:     monthlyTrend,
		DepartmentStats:  departmentStats,
	}, nil
}

func (es *EmployeeService) SearchEmployees(query string) ([]*models.Employee, error) {
	var employees []*models.Employee
	err := es.db.Where("name LIKE ? OR email LIKE ? OR employee_id LIKE ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%").
		Preload("Department.Parent.Parent.Parent.Parent.Parent").
		Preload("Position").
		Limit(20).
		Find(&employees).Error
	return employees, err
}

func (es *EmployeeService) ExportEmployees(format string) ([]byte, string, error) {
	var employees []*models.Employee
	err := es.db.Preload("Department.Parent.Parent.Parent.Parent.Parent").
		Preload("Position").
		Preload("JobLevel").
		Find(&employees).Error
	if err != nil {
		return nil, "", err
	}

	switch format {
	case "excel":
		return es.exportToExcel(employees)
	default:
		return es.exportToCSV(employees)
	}
}

func (es *EmployeeService) exportToCSV(employees []*models.Employee) ([]byte, string, error) {
	var output strings.Builder
	writer := csv.NewWriter(&output)

	// 写入标题行
	headers := []string{"工号", "姓名", "邮箱", "电话", "部门", "职位", "状态", "入职日期"}
	writer.Write(headers)

	// 写入数据行
	for _, emp := range employees {
		record := []string{
			emp.EmployeeID,
			emp.Name,
			emp.Email,
			emp.Phone,
			"",
			"",
			emp.Status,
			"",
		}

		if emp.Department != nil {
			record[4] = emp.Department.Name
		}
		if emp.Position != nil {
			record[5] = emp.Position.Name
		}
		if emp.HireDate != nil {
			record[7] = emp.HireDate.Format("2006-01-02")
		}

		writer.Write(record)
	}

	writer.Flush()
	filename := fmt.Sprintf("employees_%s.csv", time.Now().Format("20060102"))
	return []byte(output.String()), filename, nil
}

func (es *EmployeeService) exportToExcel(employees []*models.Employee) ([]byte, string, error) {
	f := excelize.NewFile()
	sheetName := "员工信息"
	f.NewSheet(sheetName)

	// 设置标题行
	headers := []string{"工号", "姓名", "邮箱", "电话", "部门", "职位", "状态", "入职日期"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入数据
	for i, emp := range employees {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), emp.EmployeeID)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), emp.Name)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), emp.Email)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), emp.Phone)
		
		if emp.Department != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), emp.Department.Name)
		}
		if emp.Position != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), emp.Position.Name)
		}
		
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), emp.Status)
		
		if emp.HireDate != nil {
			f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), emp.HireDate.Format("2006-01-02"))
		}
	}

	data, err := f.WriteToBuffer()
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("employees_%s.xlsx", time.Now().Format("20060102"))
	return data.Bytes(), filename, nil
}

func (es *EmployeeService) ImportEmployees(file *multipart.FileHeader) (*ImportResult, error) {
	// 这里简化处理，实际应该根据文件类型解析CSV或Excel
	return &ImportResult{
		Success: 0,
		Failed:  0,
		Errors:  []string{"功能暂未实现"},
	}, nil
}

func (es *EmployeeService) BulkUpdateEmployees(ids []uint, updates map[string]interface{}) (*BulkResult, error) {
	result := &BulkResult{
		Success: 0,
		Failed:  0,
		Errors:  []string{},
	}

	for _, id := range ids {
		err := es.db.Model(&models.Employee{}).Where("id = ?", id).Updates(updates).Error
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("更新员工ID %d 失败: %v", id, err))
		} else {
			result.Success++
		}
	}

	return result, nil
}

func (es *EmployeeService) BulkDeleteEmployees(ids []uint) (*BulkResult, error) {
	result := &BulkResult{
		Success: 0,
		Failed:  0,
		Errors:  []string{},
	}

	for _, id := range ids {
		err := es.db.Delete(&models.Employee{}, id).Error
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("删除员工ID %d 失败: %v", id, err))
		} else {
			result.Success++
		}
	}

	return result, nil
}

func (es *EmployeeService) GetEmployeesByDepartment(departmentID uint) ([]*models.Employee, error) {
	var employees []*models.Employee
	err := es.db.Where("department_id = ?", departmentID).
		Preload("Position").
		Preload("JobLevel").
		Find(&employees).Error
	return employees, err
}

func (es *EmployeeService) generateEmployeeID() string {
	var count int64
	es.db.Model(&models.Employee{}).Count(&count)
	// 返回纯数字格式的员工ID，匹配数据库中的mediumint类型
	return fmt.Sprintf("%d", count+1)
}

// calculateAgeDistribution 计算年龄分布
func (es *EmployeeService) calculateAgeDistribution() map[string]int64 {
	var employees []models.Employee
	ageDistribution := make(map[string]int64)
	
	// 获取所有有生日信息的员工
	if err := es.db.Where("birthday IS NOT NULL").Find(&employees).Error; err != nil {
		return ageDistribution
	}
	
	now := time.Now()
	for _, emp := range employees {
		if emp.Birthday != nil {
			age := now.Year() - emp.Birthday.Year()
			if now.YearDay() < emp.Birthday.YearDay() {
				age--
			}
			
			// 按年龄段分组
			var ageGroup string
			switch {
			case age < 25:
				ageGroup = "20-25"
			case age < 30:
				ageGroup = "26-30"
			case age < 35:
				ageGroup = "31-35"
			case age < 40:
				ageGroup = "36-40"
			case age < 45:
				ageGroup = "41-45"
			case age < 50:
				ageGroup = "46-50"
			default:
				ageGroup = "50+"
			}
			
			ageDistribution[ageGroup]++
		}
	}
	
	return ageDistribution
}

// calculateMonthlyTrend 计算月度入职趋势
func (es *EmployeeService) calculateMonthlyTrend() []MonthlyHireStats {
	var monthlyTrend []MonthlyHireStats
	now := time.Now()
	
	for i := 11; i >= 0; i-- {
		month := now.AddDate(0, -i, 0)
		monthStr := month.Format("2006-01")
		monthName := month.Format("1月")
		
		// 统计入职人数
		var hires int64
		es.db.Model(&models.Employee{}).
			Where("DATE_FORMAT(hire_date, '%Y-%m') = ?", monthStr).
			Count(&hires)
		
		// 统计离职人数（这里假设用updated_at作为离职时间，实际应该用离职日期字段）
		var leaves int64
		es.db.Model(&models.Employee{}).
			Where("status = ? AND DATE_FORMAT(updated_at, '%Y-%m') = ?", "inactive", monthStr).
			Count(&leaves)
		
		monthlyTrend = append(monthlyTrend, MonthlyHireStats{
			Month:  monthName,
			Hires:  hires,
			Leaves: leaves,
		})
	}
	
	return monthlyTrend
}

// calculateDepartmentStats 计算部门详细统计
func (es *EmployeeService) calculateDepartmentStats() []DepartmentStatItem {
	var departmentStats []DepartmentStatItem
	
	type DeptStat struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}
	
	var deptStats []DeptStat
	if err := es.db.Model(&models.Employee{}).
		Select("departments.id, departments.name, COUNT(*) as count").
		Joins("LEFT JOIN departments ON employees.department_id = departments.id").
		Group("departments.id, departments.name").
		Scan(&deptStats).Error; err != nil {
		return departmentStats
	}
	
	for _, stat := range deptStats {
		departmentStats = append(departmentStats, DepartmentStatItem{
			ID:            stat.ID,
			Name:          stat.Name,
			EmployeeCount: stat.Count,
		})
	}
	
	return departmentStats
}