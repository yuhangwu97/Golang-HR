package services

import (
	"fmt"
	"time"

	"gin-project/models"

	"gorm.io/gorm"
)

type AttendanceServiceInterface interface {
	CheckIn(employeeID uint, location, remark string) (*models.Attendance, error)
	CheckOut(employeeID uint, location, remark string) (*models.Attendance, error)
	GetAttendanceRecords(params AttendanceQueryParams) (*AttendanceListResponse, error)
	GetTodayAttendance(employeeID uint) (*models.Attendance, error)
	GetAttendanceStatistics(employeeID uint, month string) (*AttendanceStatistics, error)
	CreateLeave(leave *models.Leave) (*models.Leave, error)
	GetLeaveRecords(params LeaveQueryParams) (*LeaveListResponse, error)
	ApproveLeave(leaveID, approverID uint, status, note string) (*models.Leave, error)
	GetPendingLeaves(approverID uint, page, pageSize int) (*LeaveListResponse, error)
}

type AttendanceService struct {
	db *gorm.DB
}

func NewAttendanceService(db *gorm.DB) AttendanceServiceInterface {
	return &AttendanceService{db: db}
}

// InjectDependencies implements DependencyInjector interface
func (as *AttendanceService) InjectDependencies(deps ...interface{}) error {
	for _, dep := range deps {
		switch d := dep.(type) {
		case *gorm.DB:
			as.db = d
		}
	}
	return nil
}

type AttendanceQueryParams struct {
	EmployeeID uint   `json:"employee_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

type AttendanceListResponse struct {
	Data       []*models.Attendance `json:"data"`
	Pagination *Pagination          `json:"pagination"`
}

type AttendanceStatistics struct {
	TotalDays     int     `json:"total_days"`
	WorkDays      int     `json:"work_days"`
	LeaveDays     int     `json:"leave_days"`
	AbsentDays    int     `json:"absent_days"`
	LateCount     int     `json:"late_count"`
	EarlyCount    int     `json:"early_count"`
	WorkHours     float64 `json:"work_hours"`
	OvertimeHours float64 `json:"overtime_hours"`
}

type LeaveQueryParams struct {
	EmployeeID uint   `json:"employee_id"`
	Status     string `json:"status"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

type LeaveListResponse struct {
	Data       []*models.Leave `json:"data"`
	Pagination *Pagination     `json:"pagination"`
}

// CheckIn 员工签到
func (as *AttendanceService) CheckIn(employeeID uint, location, remark string) (*models.Attendance, error) {
	now := time.Now()
	today := now.Format("2006-01-02")

	// 检查今日是否已签到
	var existingAttendance models.Attendance
	err := as.db.Where("employee_id = ? AND DATE(created_at) = ?", employeeID, today).First(&existingAttendance).Error
	if err == nil {
		return nil, fmt.Errorf("今日已签到")
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 创建签到记录
	attendance := &models.Attendance{
		EmployeeID:  employeeID,
		Date:        now,
		CheckInTime: &now,
		Status:      "normal",
		Remark:      remark,
	}

	// 判断是否迟到 (假设上班时间为9:00)
	workStartTime := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
	if now.After(workStartTime) {
		attendance.Status = "late"
	}

	if err := as.db.Create(attendance).Error; err != nil {
		return nil, err
	}

	return as.GetAttendanceByID(attendance.ID)
}

// CheckOut 员工签退
func (as *AttendanceService) CheckOut(employeeID uint, location, remark string) (*models.Attendance, error) {
	now := time.Now()
	today := now.Format("2006-01-02")

	// 查找今日签到记录
	var attendance models.Attendance
	err := as.db.Where("employee_id = ? AND DATE(created_at) = ?", employeeID, today).First(&attendance).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("请先签到")
		}
		return nil, err
	}

	if attendance.CheckOutTime != nil {
		return nil, fmt.Errorf("今日已签退")
	}

	// 更新签退时间
	attendance.CheckOutTime = &now
	if attendance.Remark != "" {
		attendance.Remark += "; " + remark
	} else {
		attendance.Remark = remark
	}

	// 计算工作时长
	if attendance.CheckInTime != nil {
		duration := now.Sub(*attendance.CheckInTime)
		attendance.WorkHours = duration.Hours()
	}

	// 判断是否早退 (假设下班时间为18:00)
	workEndTime := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())
	if now.Before(workEndTime) && attendance.Status != "late" {
		attendance.Status = "early"
	}

	if err := as.db.Save(&attendance).Error; err != nil {
		return nil, err
	}

	return as.GetAttendanceByID(attendance.ID)
}

// GetAttendanceRecords 获取考勤记录
func (as *AttendanceService) GetAttendanceRecords(params AttendanceQueryParams) (*AttendanceListResponse, error) {
	var attendances []*models.Attendance
	var total int64

	query := as.db.Model(&models.Attendance{}).Preload("Employee")

	if params.EmployeeID > 0 {
		query = query.Where("employee_id = ?", params.EmployeeID)
	}

	if params.StartDate != "" {
		query = query.Where("date >= ?", params.StartDate)
	}

	if params.EndDate != "" {
		query = query.Where("date <= ?", params.EndDate)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("date DESC").Offset(offset).Limit(params.PageSize).Find(&attendances).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &AttendanceListResponse{
		Data: attendances,
		Pagination: &Pagination{
			Page:       params.Page,
			PageSize:   params.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// GetTodayAttendance 获取今日考勤
func (as *AttendanceService) GetTodayAttendance(employeeID uint) (*models.Attendance, error) {
	today := time.Now().Format("2006-01-02")

	var attendance models.Attendance
	err := as.db.Where("employee_id = ? AND DATE(created_at) = ?", employeeID, today).First(&attendance).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 今日未签到
		}
		return nil, err
	}

	return &attendance, nil
}

// GetAttendanceStatistics 获取考勤统计
func (as *AttendanceService) GetAttendanceStatistics(employeeID uint, month string) (*AttendanceStatistics, error) {
	var stats AttendanceStatistics

	// 查询该月的考勤记录
	var attendances []models.Attendance
	err := as.db.Where("employee_id = ? AND DATE_FORMAT(date, '%Y-%m') = ?", employeeID, month).Find(&attendances).Error
	if err != nil {
		return nil, err
	}

	// 统计各项数据
	stats.TotalDays = len(attendances)
	for _, att := range attendances {
		switch att.Status {
		case "normal":
			stats.WorkDays++
		case "late":
			stats.WorkDays++
			stats.LateCount++
		case "early":
			stats.WorkDays++
			stats.EarlyCount++
		case "leave":
			stats.LeaveDays++
		case "absent":
			stats.AbsentDays++
		}

		stats.WorkHours += att.WorkHours
		if att.WorkHours > 8 { // 假设标准工时为8小时
			stats.OvertimeHours += att.WorkHours - 8
		}
	}

	// 查询该月的请假记录
	var leaves []models.Leave
	as.db.Where("employee_id = ? AND DATE_FORMAT(start_date, '%Y-%m') = ? AND status = 'approved'", employeeID, month).Find(&leaves)
	for _, leave := range leaves {
		stats.LeaveDays += int(leave.Days)
	}

	return &stats, nil
}

// CreateLeave 创建请假申请
func (as *AttendanceService) CreateLeave(leave *models.Leave) (*models.Leave, error) {
	// 计算请假天数
	duration := leave.EndDate.Sub(leave.StartDate)
	leave.Days = duration.Hours() / 24

	if leave.Status == "" {
		leave.Status = "pending"
	}

	if err := as.db.Create(leave).Error; err != nil {
		return nil, err
	}

	return as.GetLeaveByID(leave.ID)
}

// GetLeaveRecords 获取请假记录
func (as *AttendanceService) GetLeaveRecords(params LeaveQueryParams) (*LeaveListResponse, error) {
	var leaves []*models.Leave
	var total int64

	query := as.db.Model(&models.Leave{}).Preload("Employee").Preload("Approver")

	if params.EmployeeID > 0 {
		query = query.Where("employee_id = ?", params.EmployeeID)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(params.PageSize).Find(&leaves).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &LeaveListResponse{
		Data: leaves,
		Pagination: &Pagination{
			Page:       params.Page,
			PageSize:   params.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// ApproveLeave 审批请假
func (as *AttendanceService) ApproveLeave(leaveID, approverID uint, status, note string) (*models.Leave, error) {
	var leave models.Leave
	if err := as.db.First(&leave, leaveID).Error; err != nil {
		return nil, err
	}

	if leave.Status != "pending" {
		return nil, fmt.Errorf("该请假已被处理")
	}

	now := time.Now()
	leave.Status = status
	leave.ApproverID = &approverID
	leave.ApproveTime = &now
	leave.ApproveNote = note

	if err := as.db.Save(&leave).Error; err != nil {
		return nil, err
	}

	return as.GetLeaveByID(leave.ID)
}

// GetPendingLeaves 获取待审批请假
func (as *AttendanceService) GetPendingLeaves(approverID uint, page, pageSize int) (*LeaveListResponse, error) {
	var leaves []*models.Leave
	var total int64

	// 这里简化处理，实际应该根据组织架构确定审批人权限
	query := as.db.Model(&models.Leave{}).Preload("Employee").Where("status = ?", "pending")

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&leaves).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &LeaveListResponse{
		Data: leaves,
		Pagination: &Pagination{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// 辅助方法
func (as *AttendanceService) GetAttendanceByID(id uint) (*models.Attendance, error) {
	var attendance models.Attendance
	err := as.db.Preload("Employee").First(&attendance, id).Error
	return &attendance, err
}

func (as *AttendanceService) GetLeaveByID(id uint) (*models.Leave, error) {
	var leave models.Leave
	err := as.db.Preload("Employee").Preload("Approver").First(&leave, id).Error
	return &leave, err
}
