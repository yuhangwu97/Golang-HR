package controllers

import (
	"net/http"
	"strconv"

	"gin-project/models"
	"gin-project/services"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

type DepartmentController struct {
	departmentService services.DepartmentServiceInterface
}

func NewDepartmentController(departmentService services.DepartmentServiceInterface) *DepartmentController {
	return &DepartmentController{
		departmentService: departmentService,
	}
}

// GetDepartments 获取部门列表
func (dc *DepartmentController) GetDepartments(c *gin.Context) {
	includeChildren := c.Query("includeChildren") == "true"
	parentID := c.Query("parentId")

	var parentIDPtr *uint
	if parentID != "" {
		id, err := strconv.ParseUint(parentID, 10, 32)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "无效的父部门ID")
			return
		}
		idUint := uint(id)
		parentIDPtr = &idUint
	}

	departments, err := dc.departmentService.GetDepartments(parentIDPtr, includeChildren)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取部门列表失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", gin.H{
		"data": departments,
	})
}

// GetDepartment 获取部门详情
func (dc *DepartmentController) GetDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的部门ID")
		return
	}

	department, err := dc.departmentService.GetDepartmentByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "部门不存在")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", department)
}

// CreateDepartment 创建部门
func (dc *DepartmentController) CreateDepartment(c *gin.Context) {
	var req models.Department
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	department, err := dc.departmentService.CreateDepartment(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建部门失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", department)
}

// UpdateDepartment 更新部门
func (dc *DepartmentController) UpdateDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的部门ID")
		return
	}

	var req models.Department
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	req.ID = uint(id)
	department, err := dc.departmentService.UpdateDepartment(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新部门失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", department)
}

// DeleteDepartment 删除部门
func (dc *DepartmentController) DeleteDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的部门ID")
		return
	}

	err = dc.departmentService.DeleteDepartment(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除部门失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", gin.H{"message": "删除成功"})
}

// GetDepartmentTree 获取部门树形结构
func (dc *DepartmentController) GetDepartmentTree(c *gin.Context) {
	tree, err := dc.departmentService.GetDepartmentTree()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取部门树失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", gin.H{
		"data": tree,
	})
}

// MoveDepartment 移动部门
func (dc *DepartmentController) MoveDepartment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的部门ID")
		return
	}

	var req struct {
		ParentID *uint `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	err = dc.departmentService.MoveDepartment(uint(id), req.ParentID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "移动部门失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", gin.H{"message": "移动成功"})
}

// GetDepartmentStatistics 获取部门统计信息
func (dc *DepartmentController) GetDepartmentStatistics(c *gin.Context) {
	stats, err := dc.departmentService.GetDepartmentStatistics()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取统计信息失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "操作成功", stats)
}

// GetDepartmentHierarchy 获取部门层级信息
func (dc *DepartmentController) GetDepartmentHierarchy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的部门ID")
		return
	}

	hierarchy, err := dc.departmentService.GetDepartmentHierarchy(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取部门层级信息失败")
		return
	}

	utils.SuccessResponse(c, 200, "获取部门层级信息成功", hierarchy)
}

// GetAllSubDepartments 获取所有子部门
func (dc *DepartmentController) GetAllSubDepartments(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的部门ID")
		return
	}

	subDepartments, err := dc.departmentService.GetAllSubDepartments(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取子部门失败")
		return
	}

	utils.SuccessResponse(c, 200, "获取子部门成功", subDepartments)
}

// GetDepartmentPath 获取部门路径
func (dc *DepartmentController) GetDepartmentPath(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的部门ID")
		return
	}

	path, err := dc.departmentService.GetDepartmentPath(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取部门路径失败")
		return
	}

	utils.SuccessResponse(c, 200, "获取部门路径成功", path)
}

// BulkUpdateDepartmentSort 批量更新部门排序
func (dc *DepartmentController) BulkUpdateDepartmentSort(c *gin.Context) {
	var req struct {
		Updates []services.DepartmentSortUpdate `json:"updates" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	err := dc.departmentService.BulkUpdateDepartmentSort(req.Updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新部门排序失败")
		return
	}

	utils.SuccessResponse(c, 200, "更新部门排序成功", nil)
}

// GetDepartmentChart 获取组织架构图
func (dc *DepartmentController) GetDepartmentChart(c *gin.Context) {
	chart, err := dc.departmentService.GetDepartmentChart()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取组织架构图失败")
		return
	}

	utils.SuccessResponse(c, 200, "获取组织架构图成功", chart)
}
