package controllers

import (
	"net/http"
	"strconv"

	"gin-project/models"
	"gin-project/services"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

type EmployeeController struct {
	employeeService services.EmployeeServiceInterface
}

func NewEmployeeController(employeeService services.EmployeeServiceInterface) *EmployeeController {
	return &EmployeeController{
		employeeService: employeeService,
	}
}

// GetEmployees 获取员工列表
func (ec *EmployeeController) GetEmployees(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	department := c.Query("department")
	status := c.Query("status")
	keyword := c.Query("keyword")
	sortBy := c.DefaultQuery("sortBy", "created_at")
	sortOrder := c.DefaultQuery("sortOrder", "desc")

	params := services.EmployeeListParams{
		Page:       page,
		PageSize:   pageSize,
		Department: department,
		Status:     status,
		Keyword:    keyword,
		SortBy:     sortBy,
		SortOrder:  sortOrder,
	}

	result, err := ec.employeeService.GetEmployeeList(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取员工列表失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", result)
}

// GetEmployee 获取员工详情
func (ec *EmployeeController) GetEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的员工ID")
		return
	}

	employee, err := ec.employeeService.GetEmployeeByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "员工不存在")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", employee)
}

// CreateEmployee 创建员工
func (ec *EmployeeController) CreateEmployee(c *gin.Context) {
	var req models.Employee
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	employee, err := ec.employeeService.CreateEmployee(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建员工失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", employee)
}

// UpdateEmployee 更新员工信息
func (ec *EmployeeController) UpdateEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的员工ID")
		return
	}

	var req models.Employee
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	req.ID = uint(id)
	employee, err := ec.employeeService.UpdateEmployee(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新员工信息失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", employee)
}

// DeleteEmployee 删除员工
func (ec *EmployeeController) DeleteEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的员工ID")
		return
	}

	err = ec.employeeService.DeleteEmployee(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除员工失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", gin.H{"message": "删除成功"})
}

// GetEmployeeStatistics 获取员工统计信息
func (ec *EmployeeController) GetEmployeeStatistics(c *gin.Context) {
	stats, err := ec.employeeService.GetEmployeeStatistics()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取统计信息失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", stats)
}

// SearchEmployees 搜索员工
func (ec *EmployeeController) SearchEmployees(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "搜索关键词不能为空")
		return
	}

	employees, err := ec.employeeService.SearchEmployees(query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "搜索失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", gin.H{
		"data": employees,
	})
}

// ExportEmployees 导出员工数据
func (ec *EmployeeController) ExportEmployees(c *gin.Context) {
	format := c.DefaultQuery("format", "csv")

	data, filename, err := ec.employeeService.ExportEmployees(format)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "导出失败")
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

// ImportEmployees 导入员工数据
func (ec *EmployeeController) ImportEmployees(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请选择文件")
		return
	}

	result, err := ec.employeeService.ImportEmployees(file)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "导入失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", result)
}

// BulkUpdateEmployees 批量更新员工
func (ec *EmployeeController) BulkUpdateEmployees(c *gin.Context) {
	var req struct {
		IDs     []uint                 `json:"ids" binding:"required"`
		Updates map[string]interface{} `json:"updates" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	result, err := ec.employeeService.BulkUpdateEmployees(req.IDs, req.Updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "批量更新失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", result)
}

// BulkDeleteEmployees 批量删除员工
func (ec *EmployeeController) BulkDeleteEmployees(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	result, err := ec.employeeService.BulkDeleteEmployees(req.IDs)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "批量删除失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", result)
}

// GetEmployeesByDepartment 获取部门员工
func (ec *EmployeeController) GetEmployeesByDepartment(c *gin.Context) {
	departmentID, err := strconv.ParseUint(c.Param("departmentId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的部门ID")
		return
	}

	employees, err := ec.employeeService.GetEmployeesByDepartment(uint(departmentID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取部门员工失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取部门员工成功", gin.H{
		"data": employees,
	})
}
