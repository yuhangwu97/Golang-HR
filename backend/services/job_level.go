package services

import (
	"gin-project/models"
	"gin-project/utils"

	"gorm.io/gorm"
)

type JobLevelServiceInterface interface {
	GetJobLevels(params JobLevelQueryParams) (*utils.PaginationResponse, error)
	GetJobLevelByID(id uint) (*models.JobLevel, error)
	CreateJobLevel(jobLevel *models.JobLevel) (*models.JobLevel, error)
	UpdateJobLevel(jobLevel *models.JobLevel) (*models.JobLevel, error)
	DeleteJobLevel(id uint) error
	GetJobLevelsByLevel() ([]*models.JobLevel, error)
	GetJobLevelStatistics() (*JobLevelStatistics, error)
	ValidateSalaryRange(jobLevelID uint, salary float64) error
	GetSalaryRangeByLevel(level int) (*SalaryRange, error)
}

type JobLevelService struct {
	db *gorm.DB
}

type JobLevelQueryParams struct {
	Status   string `json:"status"`
	Keyword  string `json:"keyword"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type JobLevelStatistics struct {
	Total         int64                `json:"total"`
	ByStatus      map[string]int64     `json:"by_status"`
	SalaryRanges  []SalaryRangeByLevel `json:"salary_ranges"`
	EmployeeCount map[string]int64     `json:"employee_count"`
}

type SalaryRange struct {
	MinSalary float64 `json:"min_salary"`
	MaxSalary float64 `json:"max_salary"`
}

type SalaryRangeByLevel struct {
	Level     int     `json:"level"`
	Name      string  `json:"name"`
	MinSalary float64 `json:"min_salary"`
	MaxSalary float64 `json:"max_salary"`
}

func NewJobLevelService(db *gorm.DB) JobLevelServiceInterface {
	return &JobLevelService{
		db: db,
	}
}

func (jls *JobLevelService) GetJobLevels(params JobLevelQueryParams) (*utils.PaginationResponse, error) {
	var jobLevels []models.JobLevel
	var total int64

	query := jls.db.Model(&models.JobLevel{})

	// 过滤条件
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("level ASC").Find(&jobLevels).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &utils.PaginationResponse{
		Data:       jobLevels,
		TotalItems: total,
		TotalPages: totalPages,
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}

func (jls *JobLevelService) GetJobLevelByID(id uint) (*models.JobLevel, error) {
	var jobLevel models.JobLevel
	err := jls.db.First(&jobLevel, id).Error
	if err != nil {
		return nil, err
	}
	return &jobLevel, nil
}

func (jls *JobLevelService) CreateJobLevel(jobLevel *models.JobLevel) (*models.JobLevel, error) {
	// 检查职级编码是否已存在
	var count int64
	jls.db.Model(&models.JobLevel{}).Where("code = ?", jobLevel.Code).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("职级编码已存在")
	}

	// 检查职级等级是否已存在
	jls.db.Model(&models.JobLevel{}).Where("level = ?", jobLevel.Level).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("职级等级已存在")
	}

	// 验证薪资范围
	if jobLevel.MinSalary >= jobLevel.MaxSalary {
		return nil, utils.NewValidationError("最低薪资必须小于最高薪资")
	}

	// 设置默认状态
	if jobLevel.Status == "" {
		jobLevel.Status = "active"
	}

	err := jls.db.Create(jobLevel).Error
	if err != nil {
		return nil, err
	}

	return jls.GetJobLevelByID(jobLevel.ID)
}

func (jls *JobLevelService) UpdateJobLevel(jobLevel *models.JobLevel) (*models.JobLevel, error) {
	// 检查职级编码是否已被其他职级使用
	var count int64
	jls.db.Model(&models.JobLevel{}).Where("code = ? AND id != ?", jobLevel.Code, jobLevel.ID).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("职级编码已存在")
	}

	// 检查职级等级是否已被其他职级使用
	jls.db.Model(&models.JobLevel{}).Where("level = ? AND id != ?", jobLevel.Level, jobLevel.ID).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("职级等级已存在")
	}

	// 验证薪资范围
	if jobLevel.MinSalary >= jobLevel.MaxSalary {
		return nil, utils.NewValidationError("最低薪资必须小于最高薪资")
	}

	err := jls.db.Model(jobLevel).Updates(jobLevel).Error
	if err != nil {
		return nil, err
	}

	return jls.GetJobLevelByID(jobLevel.ID)
}

func (jls *JobLevelService) DeleteJobLevel(id uint) error {
	// 检查是否有员工使用此职级
	var count int64
	jls.db.Model(&models.Employee{}).Where("job_level_id = ?", id).Count(&count)
	if count > 0 {
		return utils.NewValidationError("无法删除已被员工使用的职级")
	}

	return jls.db.Delete(&models.JobLevel{}, id).Error
}

func (jls *JobLevelService) GetJobLevelsByLevel() ([]*models.JobLevel, error) {
	var jobLevels []*models.JobLevel
	err := jls.db.Where("status = ?", "active").
		Order("level ASC").Find(&jobLevels).Error
	return jobLevels, err
}

func (jls *JobLevelService) GetJobLevelStatistics() (*JobLevelStatistics, error) {
	var total int64
	jls.db.Model(&models.JobLevel{}).Count(&total)

	// 按状态统计
	var statusStats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	jls.db.Model(&models.JobLevel{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusStats)

	byStatus := make(map[string]int64)
	for _, stat := range statusStats {
		byStatus[stat.Status] = stat.Count
	}

	// 薪资范围统计
	var salaryRanges []SalaryRangeByLevel
	jls.db.Model(&models.JobLevel{}).
		Select("level, name, min_salary, max_salary").
		Where("status = ?", "active").
		Order("level ASC").
		Scan(&salaryRanges)

	// 员工数量统计
	var employeeStats []struct {
		JobLevelName string `json:"job_level_name"`
		Count        int64  `json:"count"`
	}
	jls.db.Model(&models.Employee{}).
		Select("job_levels.name as job_level_name, COUNT(*) as count").
		Joins("LEFT JOIN job_levels ON employees.job_level_id = job_levels.id").
		Group("job_levels.name").
		Scan(&employeeStats)

	employeeCount := make(map[string]int64)
	for _, stat := range employeeStats {
		employeeCount[stat.JobLevelName] = stat.Count
	}

	return &JobLevelStatistics{
		Total:         total,
		ByStatus:      byStatus,
		SalaryRanges:  salaryRanges,
		EmployeeCount: employeeCount,
	}, nil
}

func (jls *JobLevelService) ValidateSalaryRange(jobLevelID uint, salary float64) error {
	var jobLevel models.JobLevel
	err := jls.db.First(&jobLevel, jobLevelID).Error
	if err != nil {
		return err
	}

	if salary < jobLevel.MinSalary || salary > jobLevel.MaxSalary {
		return utils.NewValidationError("薪资超出职级范围")
	}

	return nil
}

func (jls *JobLevelService) GetSalaryRangeByLevel(level int) (*SalaryRange, error) {
	var jobLevel models.JobLevel
	err := jls.db.Where("level = ? AND status = ?", level, "active").First(&jobLevel).Error
	if err != nil {
		return nil, err
	}

	return &SalaryRange{
		MinSalary: jobLevel.MinSalary,
		MaxSalary: jobLevel.MaxSalary,
	}, nil
}
