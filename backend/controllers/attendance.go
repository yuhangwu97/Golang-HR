package controllers

import (
	"net/http"
	"strconv"
	"time"

	"gin-project/models"
	"gin-project/services"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

type AttendanceController struct {
	attendanceService services.AttendanceServiceInterface
}

func NewAttendanceController(attendanceService services.AttendanceServiceInterface) *AttendanceController {
	return &AttendanceController{
		attendanceService: attendanceService,
	}
}

// CheckIn 签到
func (ac *AttendanceController) CheckIn(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	var req struct {
		Location string `json:"location"`
		Remark   string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	attendance, err := ac.attendanceService.CheckIn(userID, req.Location, req.Remark)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "签到失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", attendance)
}

// CheckOut 签退
func (ac *AttendanceController) CheckOut(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	var req struct {
		Location string `json:"location"`
		Remark   string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	attendance, err := ac.attendanceService.CheckOut(userID, req.Location, req.Remark)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "签退失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", attendance)
}

// GetAttendanceRecords 获取考勤记录
func (ac *AttendanceController) GetAttendanceRecords(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Query("employee_id"), 10, 32)
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	params := services.AttendanceQueryParams{
		EmployeeID: uint(employeeID),
		StartDate:  startDate,
		EndDate:    endDate,
		Page:       page,
		PageSize:   pageSize,
	}

	result, err := ac.attendanceService.GetAttendanceRecords(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取考勤记录失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", result)
}

// GetTodayAttendance 获取今日考勤状态
func (ac *AttendanceController) GetTodayAttendance(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	attendance, err := ac.attendanceService.GetTodayAttendance(userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取今日考勤失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", attendance)
}

// GetAttendanceStatistics 获取考勤统计
func (ac *AttendanceController) GetAttendanceStatistics(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Query("employee_id"), 10, 32)
	month := c.Query("month")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	stats, err := ac.attendanceService.GetAttendanceStatistics(uint(employeeID), month)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取考勤统计失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", stats)
}

// CreateLeave 创建请假申请
func (ac *AttendanceController) CreateLeave(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未登录")
		return
	}

	var req models.Leave
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	req.EmployeeID = userID
	leave, err := ac.attendanceService.CreateLeave(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建请假申请失败")
		return
	}

	utils.SuccessResponse(c, 200, "操作成功", leave)
}

// GetLeaveRecords 获取请假记录
func (ac *AttendanceController) GetLeaveRecords(c *gin.Context) {
	employeeID, _ := strconv.ParseUint(c.Query("employee_id"), 10, 32)
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	params := services.LeaveQueryParams{
		EmployeeID: uint(employeeID),
		Status:     status,
		Page:       page,
		PageSize:   pageSize,
	}

	result, err := ac.attendanceService.GetLeaveRecords(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取请假记录失败")
		return
	}

	utils.SuccessResponse(c, 200, "获取请假记录成功", result)
}

// ApproveLeave 审批请假
func (ac *AttendanceController) ApproveLeave(c *gin.Context) {
	userID := c.GetUint("user_id")
	leaveID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的请假ID")
		return
	}

	var req struct {
		Status      string `json:"status" binding:"required"`
		ApproveNote string `json:"approve_note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	leave, err := ac.attendanceService.ApproveLeave(uint(leaveID), userID, req.Status, req.ApproveNote)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "审批失败")
		return
	}

	utils.SuccessResponse(c, 200, "操作成功", leave)
}

// GetPendingLeaves 获取待审批请假
func (ac *AttendanceController) GetPendingLeaves(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	result, err := ac.attendanceService.GetPendingLeaves(userID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取待审批请假失败")
		return
	}

	utils.SuccessResponse(c, 200, "获取待审批请假成功", result)
}
