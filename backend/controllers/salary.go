package controllers

import (
	"net/http"
	"strconv"

	"gin-project/models"
	"gin-project/services"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

type SalaryController struct {
	salaryService services.SalaryServiceInterface
}

func NewSalaryController(salaryService services.SalaryServiceInterface) *SalaryController {
	return &SalaryController{
		salaryService: salaryService,
	}
}

// ========================= Legacy Salary Management =========================

// GetSalaryRecords 获取薪资记录
func (sc *SalaryController) GetSalaryRecords(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Query("employee_id"), 10, 32)
	month := c.Query("month")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	params := services.SalaryQueryParams{
		EmployeeID: uint(employeeID),
		Month:      month,
		Status:     status,
		Page:       page,
		PageSize:   pageSize,
	}

	result, err := sc.salaryService.GetSalaryRecords(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取薪资记录失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", result)
}

// GetSalaryDetail 获取薪资详情
func (sc *SalaryController) GetSalaryDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的薪资ID")
		return
	}

	salary, err := sc.salaryService.GetSalaryByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "薪资记录不存在")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", salary)
}

// CreateSalary 创建薪资记录
func (sc *SalaryController) CreateSalary(c *gin.Context) {
	var req models.Salary
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	salary, err := sc.salaryService.CreateSalary(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建薪资记录失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", salary)
}

// UpdateSalary 更新薪资记录
func (sc *SalaryController) UpdateSalary(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的薪资ID")
		return
	}

	var req models.Salary
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	req.ID = uint(id)
	salary, err := sc.salaryService.UpdateSalary(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新薪资记录失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", salary)
}

// DeleteSalary 删除薪资记录
func (sc *SalaryController) DeleteSalary(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的薪资ID")
		return
	}

	err = sc.salaryService.DeleteSalary(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除薪资记录失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", gin.H{"message": "删除成功"})
}

// GetMySalary 获取我的薪资
func (sc *SalaryController) GetMySalary(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	month := c.Query("month")
	if month == "" {
		month = "2024-01"
	}

	salary, err := sc.salaryService.GetEmployeeSalary(userID, month)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "未找到薪资记录")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", salary)
}

// GetMyEnhancedSalary 获取我的增强版薪资详情
func (sc *SalaryController) GetMyEnhancedSalary(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	var periodID *uint
	if pID, err := strconv.ParseUint(c.Query("period_id"), 10, 32); err == nil {
		id := uint(pID)
		periodID = &id
	}

	var year *int
	if y, err := strconv.Atoi(c.Query("year")); err == nil {
		year = &y
	}

	var month *int
	if m, err := strconv.Atoi(c.Query("month")); err == nil {
		month = &m
	}

	params := services.PersonalSalaryParams{
		EmployeeID: userID,
		PeriodID:   periodID,
		Year:       year,
		Month:      month,
	}

	salaryDetail, err := sc.salaryService.GetPersonalSalaryDetail(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取薪资详情失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取成功", salaryDetail)
}

// GetMySalaryHistory 获取我的薪资历史
func (sc *SalaryController) GetMySalaryHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	
	params := services.SalaryHistoryParams{
		EmployeeID: userID,
		Limit:      limit,
	}

	history, err := sc.salaryService.GetSalaryHistory(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取薪资历史失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取成功", history)
}

// GetMyPayrollRecords 获取我的发放记录
func (sc *SalaryController) GetMyPayrollRecords(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	params := services.PersonalPayrollParams{
		EmployeeID: userID,
		Page:       page,
		PageSize:   pageSize,
	}

	records, err := sc.salaryService.GetPersonalPayrollRecords(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取发放记录失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取成功", records)
}

// GetMySalaryDashboard 获取我的薪资仪表板
func (sc *SalaryController) GetMySalaryDashboard(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	dashboard, err := sc.salaryService.GetPersonalSalaryDashboard(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取薪资仪表板失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取成功", dashboard)
}

// ========================= Enhanced Salary Component Management =========================

// CreateSalaryComponent 创建薪资组件
func (sc *SalaryController) CreateSalaryComponent(c *gin.Context) {
	var component models.SalaryComponent
	if err := c.ShouldBindJSON(&component); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	result, err := sc.salaryService.CreateSalaryComponent(&component)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建薪资组件失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "创建成功", result)
}

// UpdateSalaryComponent 更新薪资组件
func (sc *SalaryController) UpdateSalaryComponent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的组件ID")
		return
	}

	var component models.SalaryComponent
	if err := c.ShouldBindJSON(&component); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	result, err := sc.salaryService.UpdateSalaryComponent(uint(id), &component)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新薪资组件失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "更新成功", result)
}

// DeleteSalaryComponent 删除薪资组件
func (sc *SalaryController) DeleteSalaryComponent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的组件ID")
		return
	}

	err = sc.salaryService.DeleteSalaryComponent(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除薪资组件失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "删除成功", nil)
}

// GetSalaryComponents 获取薪资组件列表
func (sc *SalaryController) GetSalaryComponents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	params := services.ComponentQueryParams{
		Category: c.Query("category"),
		Type:     c.Query("type"),
		Status:   c.Query("status"),
		Page:     page,
		PageSize: pageSize,
	}

	result, err := sc.salaryService.GetSalaryComponents(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取薪资组件失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取成功", result)
}

// GetSalaryComponent 获取薪资组件详情
func (sc *SalaryController) GetSalaryComponent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的组件ID")
		return
	}

	component, err := sc.salaryService.GetSalaryComponentByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "薪资组件不存在")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取成功", component)
}

// ========================= Salary Grade Management =========================

// CreateSalaryGrade 创建薪资等级
func (sc *SalaryController) CreateSalaryGrade(c *gin.Context) {
	var grade models.SalaryGrade
	if err := c.ShouldBindJSON(&grade); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	result, err := sc.salaryService.CreateSalaryGrade(&grade)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建薪资等级失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "创建成功", result)
}

// GetSalaryGrades 获取薪资等级列表
func (sc *SalaryController) GetSalaryGrades(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	level, _ := strconv.Atoi(c.Query("level"))

	params := services.GradeQueryParams{
		Level:    level,
		Currency: c.Query("currency"),
		Status:   c.Query("status"),
		Page:     page,
		PageSize: pageSize,
	}

	result, err := sc.salaryService.GetSalaryGrades(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取薪资等级失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取成功", result)
}

// ========================= Salary Structure Management =========================

// CreateSalaryStructure 创建薪资结构
func (sc *SalaryController) CreateSalaryStructure(c *gin.Context) {
	var structure models.SalaryStructure
	if err := c.ShouldBindJSON(&structure); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	result, err := sc.salaryService.CreateSalaryStructure(&structure)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建薪资结构失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "创建成功", result)
}

// GetSalaryStructures 获取薪资结构列表
func (sc *SalaryController) GetSalaryStructures(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var departmentID *uint
	if deptID, err := strconv.ParseUint(c.Query("department_id"), 10, 32); err == nil {
		id := uint(deptID)
		departmentID = &id
	}

	var positionID *uint
	if posID, err := strconv.ParseUint(c.Query("position_id"), 10, 32); err == nil {
		id := uint(posID)
		positionID = &id
	}

	var jobLevelID *uint
	if levelID, err := strconv.ParseUint(c.Query("job_level_id"), 10, 32); err == nil {
		id := uint(levelID)
		jobLevelID = &id
	}

	params := services.StructureQueryParams{
		DepartmentID: departmentID,
		PositionID:   positionID,
		JobLevelID:   jobLevelID,
		Status:       c.Query("status"),
		Page:         page,
		PageSize:     pageSize,
	}

	result, err := sc.salaryService.GetSalaryStructures(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取薪资结构失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取成功", result)
}

// ========================= Payroll Period Management =========================

// CreatePayrollPeriod 创建薪资周期
func (sc *SalaryController) CreatePayrollPeriod(c *gin.Context) {
	var period models.PayrollPeriod
	if err := c.ShouldBindJSON(&period); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	result, err := sc.salaryService.CreatePayrollPeriod(&period)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建薪资周期失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "创建成功", result)
}

// GetPayrollPeriods 获取薪资周期列表
func (sc *SalaryController) GetPayrollPeriods(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	year, _ := strconv.Atoi(c.Query("year"))

	params := services.PeriodQueryParams{
		PeriodType: c.Query("period_type"),
		Year:       year,
		Status:     c.Query("status"),
		Page:       page,
		PageSize:   pageSize,
	}

	result, err := sc.salaryService.GetPayrollPeriods(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取薪资周期失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取成功", result)
}

// LockPayrollPeriod 锁定薪资周期
func (sc *SalaryController) LockPayrollPeriod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的周期ID")
		return
	}

	err = sc.salaryService.LockPayrollPeriod(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "锁定薪资周期失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "锁定成功", nil)
}

// UnlockPayrollPeriod 解锁薪资周期
func (sc *SalaryController) UnlockPayrollPeriod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的周期ID")
		return
	}

	err = sc.salaryService.UnlockPayrollPeriod(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "解锁薪资周期失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "解锁成功", nil)
}

// ========================= Enhanced Salary Calculation =========================

// CalculateSalary 计算薪资 (Legacy)
func (sc *SalaryController) CalculateSalary(c *gin.Context) {
	var req struct {
		EmployeeID uint   `json:"employee_id" binding:"required"`
		Month      string `json:"month" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	salary, err := sc.salaryService.CalculateSalary(req.EmployeeID, req.Month)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "计算薪资失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", salary)
}

// CalculateEmployeeSalary 计算员工薪资 (Enhanced)
func (sc *SalaryController) CalculateEmployeeSalary(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	var req struct {
		EmployeeID uint `json:"employee_id" binding:"required"`
		PeriodID   uint `json:"period_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	salary, err := sc.salaryService.CalculateEmployeeSalary(req.EmployeeID, req.PeriodID, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "计算薪资失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "计算成功", salary)
}

// BatchCalculateSalary 批量计算薪资 (Legacy)
func (sc *SalaryController) BatchCalculateSalary(c *gin.Context) {
	var req struct {
		Month        string `json:"month" binding:"required"`
		DepartmentID uint   `json:"department_id"`
		EmployeeIDs  []uint `json:"employee_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	result, err := sc.salaryService.BatchCalculateSalary(req.Month, req.DepartmentID, req.EmployeeIDs)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "批量计算薪资失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", result)
}

// BatchCalculateSalaries 批量计算薪资 (Enhanced)
func (sc *SalaryController) BatchCalculateSalaries(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	var req struct {
		PeriodID     uint   `json:"period_id" binding:"required"`
		DepartmentID *uint  `json:"department_id"`
		EmployeeIDs  []uint `json:"employee_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	result, err := sc.salaryService.BatchCalculateSalaries(req.PeriodID, req.DepartmentID, req.EmployeeIDs, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "批量计算薪资失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "批量计算成功", result)
}

// GetEnhancedSalaries 获取增强薪资记录列表
func (sc *SalaryController) GetEnhancedSalaries(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var employeeID *uint
	if empID, err := strconv.ParseUint(c.Query("employee_id"), 10, 32); err == nil {
		id := uint(empID)
		employeeID = &id
	}

	var periodID *uint
	if pID, err := strconv.ParseUint(c.Query("period_id"), 10, 32); err == nil {
		id := uint(pID)
		periodID = &id
	}

	var departmentID *uint
	if deptID, err := strconv.ParseUint(c.Query("department_id"), 10, 32); err == nil {
		id := uint(deptID)
		departmentID = &id
	}

	params := services.EnhancedSalaryQueryParams{
		EmployeeID:      employeeID,
		PayrollPeriodID: periodID,
		DepartmentID:    departmentID,
		Status:          c.Query("status"),
		Page:            page,
		PageSize:        pageSize,
	}

	result, err := sc.salaryService.GetEnhancedSalaries(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取薪资记录失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取成功", result)
}

// GetEnhancedSalary 获取增强薪资记录详情
func (sc *SalaryController) GetEnhancedSalary(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的薪资ID")
		return
	}

	salary, err := sc.salaryService.GetEnhancedSalaryByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "薪资记录不存在")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取成功", salary)
}

// UpdateSalaryDetails 更新薪资详情
func (sc *SalaryController) UpdateSalaryDetails(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的薪资ID")
		return
	}

	var req struct {
		Details []services.SalaryDetailUpdate `json:"details" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	salary, err := sc.salaryService.UpdateSalaryDetails(uint(id), req.Details, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新薪资详情失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "更新成功", salary)
}

// ========================= Approval Workflow =========================

// ApproveSalary 审批薪资 (Legacy)
func (sc *SalaryController) ApproveSalary(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的薪资ID")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
		Remark string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	salary, err := sc.salaryService.ApproveSalary(uint(id), req.Status, req.Remark)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "审批薪资失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", salary)
}

// ReviewSalary 审核薪资 (Enhanced)
func (sc *SalaryController) ReviewSalary(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的薪资ID")
		return
	}

	var req struct {
		Notes   string `json:"notes"`
		Approve bool   `json:"approve"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	salary, err := sc.salaryService.ReviewSalary(uint(id), userID, req.Notes, req.Approve)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "审核薪资失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "审核成功", salary)
}

// ApproveEnhancedSalary 批准薪资 (Enhanced)
func (sc *SalaryController) ApproveEnhancedSalary(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的薪资ID")
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	salary, err := sc.salaryService.ApproveEnhancedSalary(uint(id), userID, req.Notes)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "批准薪资失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "批准成功", salary)
}

// BulkApproveSalaries 批量批准薪资
func (sc *SalaryController) BulkApproveSalaries(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	var req struct {
		SalaryIDs []uint `json:"salary_ids" binding:"required"`
		Notes     string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	result, err := sc.salaryService.BulkApproveSalaries(req.SalaryIDs, userID, req.Notes)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "批量批准失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "批量批准成功", result)
}

// ========================= Payment Processing =========================

// ProcessPayroll 处理薪资发放 (Legacy)
func (sc *SalaryController) ProcessPayroll(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的薪资ID")
		return
	}

	var req struct {
		PaymentMethod string `json:"payment_method" binding:"required"`
		BankAccount   string `json:"bank_account" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	payroll, err := sc.salaryService.ProcessPayroll(uint(id), req.PaymentMethod, req.BankAccount)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "处理薪资发放失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "薪资发放处理成功", payroll)
}

// GetPayrollRecords 获取薪资发放记录 (Legacy)
func (sc *SalaryController) GetPayrollRecords(c *gin.Context) {
	salaryID, err := strconv.ParseUint(c.Param("salary_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的薪资ID")
		return
	}

	records, err := sc.salaryService.GetPayrollRecords(uint(salaryID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取发放记录失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取发放记录成功", records)
}

// UpdatePayrollStatus 更新发放状态 (Legacy)
func (sc *SalaryController) UpdatePayrollStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的发放记录ID")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
		Remark string `json:"remark"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	err = sc.salaryService.UpdatePayrollStatus(uint(id), req.Status, req.Remark)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新发放状态失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "更新成功", nil)
}

// CreatePaymentBatch 创建支付批次 (Enhanced)
func (sc *SalaryController) CreatePaymentBatch(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	var req struct {
		Batch     models.PaymentBatch `json:"batch" binding:"required"`
		SalaryIDs []uint              `json:"salary_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	batch, err := sc.salaryService.CreatePaymentBatch(&req.Batch, req.SalaryIDs, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建支付批次失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "创建成功", batch)
}

// ProcessPaymentBatch 处理支付批次 (Enhanced)
func (sc *SalaryController) ProcessPaymentBatch(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的批次ID")
		return
	}

	batch, err := sc.salaryService.ProcessPaymentBatch(uint(id), userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "处理支付批次失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "处理成功", batch)
}

// GetPaymentBatches 获取支付批次列表 (Enhanced)
func (sc *SalaryController) GetPaymentBatches(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var periodID *uint
	if pID, err := strconv.ParseUint(c.Query("period_id"), 10, 32); err == nil {
		id := uint(pID)
		periodID = &id
	}

	var createdBy *uint
	if cID, err := strconv.ParseUint(c.Query("created_by"), 10, 32); err == nil {
		id := uint(cID)
		createdBy = &id
	}

	params := services.BatchQueryParams{
		PayrollPeriodID: periodID,
		Status:          c.Query("status"),
		CreatedBy:       createdBy,
		Page:            page,
		PageSize:        pageSize,
	}

	result, err := sc.salaryService.GetPaymentBatches(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取支付批次失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取成功", result)
}

// ========================= Analytics and Reporting =========================

// GetSalaryStatistics 获取薪资统计 (Legacy)
func (sc *SalaryController) GetSalaryStatistics(c *gin.Context) {
	month := c.Query("month")
	departmentID, _ := strconv.ParseUint(c.Query("department_id"), 10, 32)

	stats, err := sc.salaryService.GetSalaryStatistics(month, uint(departmentID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取薪资统计失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", stats)
}

// GetSalaryAnalytics 获取薪资分析 (Enhanced)
func (sc *SalaryController) GetSalaryAnalytics(c *gin.Context) {
	var departmentID *uint
	if deptID, err := strconv.ParseUint(c.Query("department_id"), 10, 32); err == nil {
		id := uint(deptID)
		departmentID = &id
	}

	var periodID *uint
	if pID, err := strconv.ParseUint(c.Query("period_id"), 10, 32); err == nil {
		id := uint(pID)
		periodID = &id
	}

	var year *int
	if y, err := strconv.Atoi(c.Query("year")); err == nil {
		year = &y
	}

	var month *int
	if m, err := strconv.Atoi(c.Query("month")); err == nil {
		month = &m
	}

	params := services.AnalyticsParams{
		DepartmentID: departmentID,
		PeriodID:     periodID,
		Year:         year,
		Month:        month,
	}

	analytics, err := sc.salaryService.GetSalaryAnalytics(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取薪资分析失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取成功", analytics)
}

// GetDepartmentSalaryReport 获取部门薪资报告
func (sc *SalaryController) GetDepartmentSalaryReport(c *gin.Context) {
	departmentID, err := strconv.ParseUint(c.Query("department_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "部门ID参数错误")
		return
	}

	periodID, err := strconv.ParseUint(c.Query("period_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "周期ID参数错误")
		return
	}

	report, err := sc.salaryService.GetDepartmentSalaryReport(uint(departmentID), uint(periodID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取部门薪资报告失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取成功", report)
}

// ExportSalaryReport 导出薪资报表 (Legacy)
func (sc *SalaryController) ExportSalaryReport(c *gin.Context) {
	month := c.Query("month")
	departmentID, _ := strconv.ParseUint(c.Query("department_id"), 10, 32)
	format := c.DefaultQuery("format", "excel")

	data, filename, err := sc.salaryService.ExportSalaryReport(month, uint(departmentID), format)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "导出薪资报表失败")
		return
	}

	var contentType string
	switch format {
	case "excel":
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	default:
		contentType = "text/csv"
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", contentType)
	c.Data(http.StatusOK, contentType, data)
}

// ExportEnhancedSalaryReport 导出薪资报告 (Enhanced)
func (sc *SalaryController) ExportEnhancedSalaryReport(c *gin.Context) {
	var periodID *uint
	if pID, err := strconv.ParseUint(c.Query("period_id"), 10, 32); err == nil {
		id := uint(pID)
		periodID = &id
	}

	var departmentID *uint
	if deptID, err := strconv.ParseUint(c.Query("department_id"), 10, 32); err == nil {
		id := uint(deptID)
		departmentID = &id
	}

	params := services.ExportParams{
		PeriodID:     periodID,
		DepartmentID: departmentID,
		Format:       c.DefaultQuery("format", "excel"),
		Template:     c.Query("template"),
	}

	data, filename, err := sc.salaryService.ExportEnhancedSalaryReport(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "导出薪资报告失败: "+err.Error())
		return
	}

	var contentType string
	switch params.Format {
	case "excel":
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "pdf":
		contentType = "application/pdf"
	default:
		contentType = "text/csv"
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", contentType)
	c.Data(http.StatusOK, contentType, data)
}

// ValidateFormula 验证薪资公式
func (sc *SalaryController) ValidateFormula(c *gin.Context) {
	var req struct {
		Formula string `json:"formula" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误: "+err.Error())
		return
	}

	err := sc.salaryService.ValidateFormula(req.Formula)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "公式验证失败: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "公式验证成功", gin.H{"valid": true})
}