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

type OrganizationController struct {
	organizationService services.OrganizationServiceInterface
}

func NewOrganizationController(organizationService services.OrganizationServiceInterface) *OrganizationController {
	return &OrganizationController{
		organizationService: organizationService,
	}
}

func (oc *OrganizationController) parseQueryParams(c *gin.Context, paramNames []string) map[string]interface{} {
	params := make(map[string]interface{})
	
	for _, name := range paramNames {
		if value := c.Query(name); value != "" {
			if name == "is_active" || name == "is_primary" {
				params[name] = value == "true"
			} else {
				params[name] = value
			}
		}
	}
	
	return params
}

func (oc *OrganizationController) getUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

func (oc *OrganizationController) parseUintParam(c *gin.Context, paramName string) (uint, error) {
	id, err := strconv.ParseUint(c.Param(paramName), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// GetOrganizationUnits 获取组织单元列表
func (oc *OrganizationController) GetOrganizationUnits(c *gin.Context) {
	params := oc.parseQueryParams(c, []string{"parent_id", "type", "status", "is_active", "keyword", "page", "size"})
	
	result, err := oc.organizationService.GetOrganizationUnits(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "获取组织单元列表成功", result)
}

// GetOrganizationUnit 获取单个组织单元
func (oc *OrganizationController) GetOrganizationUnit(c *gin.Context) {
	id, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的组织单元ID")
		return
	}
	
	unit, err := oc.organizationService.GetOrganizationUnit(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "获取组织单元成功", unit)
}

// CreateOrganizationUnit 创建组织单元
func (oc *OrganizationController) CreateOrganizationUnit(c *gin.Context) {
	var unit models.OrganizationUnit
	if err := c.ShouldBindJSON(&unit); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}
	
	userID, exists := oc.getUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	
	createdUnit, err := oc.organizationService.CreateOrganizationUnit(&unit, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusCreated, "创建组织单元成功", createdUnit)
}

// UpdateOrganizationUnit 更新组织单元
func (oc *OrganizationController) UpdateOrganizationUnit(c *gin.Context) {
	id, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的组织单元ID")
		return
	}
	
	var updates models.OrganizationUnit
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}
	
	userID, exists := oc.getUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	
	updatedUnit, err := oc.organizationService.UpdateOrganizationUnit(id, &updates, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "更新组织单元成功", updatedUnit)
}

// DeleteOrganizationUnit 删除组织单元
func (oc *OrganizationController) DeleteOrganizationUnit(c *gin.Context) {
	id, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的组织单元ID")
		return
	}
	
	userID, exists := oc.getUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	
	err = oc.organizationService.DeleteOrganizationUnit(id, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "删除组织单元成功", nil)
}

// GetOrganizationTree 获取组织架构树
func (oc *OrganizationController) GetOrganizationTree(c *gin.Context) {
	tree, err := oc.organizationService.GetOrganizationTree()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "获取组织架构树成功", tree)
}

// GetSubunits 获取子单元
func (oc *OrganizationController) GetSubunits(c *gin.Context) {
	parentID, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的父单元ID")
		return
	}
	
	subunits, err := oc.organizationService.GetSubunits(parentID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "获取子单元成功", subunits)
}

// MoveUnit 移动组织单元
func (oc *OrganizationController) MoveUnit(c *gin.Context) {
	unitID, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的组织单元ID")
		return
	}
	
	var req struct {
		ParentID *uint `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}
	
	userID, exists := oc.getUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	
	err = oc.organizationService.MoveUnit(unitID, req.ParentID, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "移动组织单元成功", nil)
}

// GetHierarchyPath 获取层级路径
func (oc *OrganizationController) GetHierarchyPath(c *gin.Context) {
	unitID, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的组织单元ID")
		return
	}
	
	path, err := oc.organizationService.GetHierarchyPath(unitID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "获取层级路径成功", gin.H{"path": path})
}

// GetUnitAssignments 获取组织单元的员工分配
func (oc *OrganizationController) GetUnitAssignments(c *gin.Context) {
	unitID, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的组织单元ID")
		return
	}
	
	params := oc.parseQueryParams(c, []string{"assignment_type", "status", "is_primary", "keyword", "page", "size"})
	
	result, err := oc.organizationService.GetUnitAssignments(unitID, params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "获取员工分配成功", result)
}

// AssignEmployee 分配员工到组织单元
func (oc *OrganizationController) AssignEmployee(c *gin.Context) {
	var assignment models.EmployeeAssignment
	if err := c.ShouldBindJSON(&assignment); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}
	
	userID, exists := oc.getUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	
	createdAssignment, err := oc.organizationService.AssignEmployee(&assignment, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusCreated, "员工分配成功", createdAssignment)
}

// UpdateAssignment 更新员工分配
func (oc *OrganizationController) UpdateAssignment(c *gin.Context) {
	id, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的分配ID")
		return
	}
	
	var updates models.EmployeeAssignment
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}
	
	userID, exists := oc.getUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	
	updatedAssignment, err := oc.organizationService.UpdateAssignment(id, &updates, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "更新员工分配成功", updatedAssignment)
}

// RemoveAssignment 移除员工分配
func (oc *OrganizationController) RemoveAssignment(c *gin.Context) {
	id, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的分配ID")
		return
	}
	
	userID, exists := oc.getUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	
	err = oc.organizationService.RemoveAssignment(id, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "移除员工分配成功", nil)
}

// GetOrganizationChanges 获取组织变更记录
func (oc *OrganizationController) GetOrganizationChanges(c *gin.Context) {
	params := oc.parseQueryParams(c, []string{"unit_id", "change_type", "status", "entity_type", "initiator_id", "start_date", "end_date", "page", "size"})
	
	result, err := oc.organizationService.GetOrganizationChanges(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "获取组织变更记录成功", result)
}

// CreateOrganizationChange 创建组织变更
func (oc *OrganizationController) CreateOrganizationChange(c *gin.Context) {
	var change models.OrganizationChange
	if err := c.ShouldBindJSON(&change); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}
	
	userID, exists := oc.getUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	
	change.InitiatorID = userID
	
	createdChange, err := oc.organizationService.CreateOrganizationChange(&change)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusCreated, "创建组织变更成功", createdChange)
}

// ApproveChange 审批变更
func (oc *OrganizationController) ApproveChange(c *gin.Context) {
	id, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的变更ID")
		return
	}
	
	var req struct {
		ApprovalNote  string     `json:"approval_note"`
		EffectiveDate *time.Time `json:"effective_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}
	
	userID, exists := oc.getUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	
	err = oc.organizationService.ApproveChange(id, userID, req.ApprovalNote, req.EffectiveDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "审批变更成功", nil)
}

// RejectChange 拒绝变更
func (oc *OrganizationController) RejectChange(c *gin.Context) {
	id, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的变更ID")
		return
	}
	
	var req struct {
		ApprovalNote string `json:"approval_note" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}
	
	userID, exists := oc.getUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	
	err = oc.organizationService.RejectChange(id, userID, req.ApprovalNote)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "拒绝变更成功", nil)
}

// GetOrganizationSnapshots 获取组织架构历史快照
func (oc *OrganizationController) GetOrganizationSnapshots(c *gin.Context) {
	params := oc.parseQueryParams(c, []string{"unit_id", "change_type", "changed_by", "start_date", "end_date", "page", "size"})
	
	result, err := oc.organizationService.GetOrganizationSnapshots(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "获取历史快照成功", result)
}

// CreateOrganizationSnapshot 创建组织架构快照
func (oc *OrganizationController) CreateOrganizationSnapshot(c *gin.Context) {
	var req struct {
		UnitID *uint  `json:"unit_id"`
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}
	
	userID, exists := oc.getUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	
	err := oc.organizationService.CreateOrganizationSnapshot(req.UnitID, req.Reason, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusCreated, "创建快照成功", nil)
}

// GetUnitHistory 获取组织单元历史记录
func (oc *OrganizationController) GetUnitHistory(c *gin.Context) {
	unitID, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的组织单元ID")
		return
	}
	
	params := oc.parseQueryParams(c, []string{"change_type", "start_date", "end_date", "page", "size"})
	
	result, err := oc.organizationService.GetUnitHistory(unitID, params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "获取单元历史记录成功", result)
}

// CompareHistoryVersions 比较历史版本
func (oc *OrganizationController) CompareHistoryVersions(c *gin.Context) {
	unitID, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的组织单元ID")
		return
	}
	
	fromDateStr := c.Query("from_date")
	toDateStr := c.Query("to_date")
	
	if fromDateStr == "" || toDateStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "请提供起始和结束日期")
		return
	}
	
	fromDate, err := time.Parse("2006-01-02 15:04:05", fromDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "起始日期格式无效")
		return
	}
	
	toDate, err := time.Parse("2006-01-02 15:04:05", toDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "结束日期格式无效")
		return
	}
	
	comparison, err := oc.organizationService.CompareHistoryVersions(unitID, fromDate, toDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "版本比较成功", comparison)
}

// RollbackToHistory 回滚到历史版本
func (oc *OrganizationController) RollbackToHistory(c *gin.Context) {
	unitID, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的组织单元ID")
		return
	}
	
	var req struct {
		HistoryID uint `json:"history_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}
	
	userID, exists := oc.getUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	
	err = oc.organizationService.RollbackToHistory(unitID, req.HistoryID, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "回滚成功", nil)
}

// GetUnitStatistics 获取组织单元统计信息
func (oc *OrganizationController) GetUnitStatistics(c *gin.Context) {
	unitID, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的组织单元ID")
		return
	}
	
	stats, err := oc.organizationService.GetUnitStatistics(unitID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "获取统计信息成功", stats)
}

// SearchUnits 搜索组织单元
func (oc *OrganizationController) SearchUnits(c *gin.Context) {
	params := oc.parseQueryParams(c, []string{"keyword", "type", "status", "is_active", "manager_id", "page", "size"})
	
	result, err := oc.organizationService.SearchUnits(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "搜索组织单元成功", result)
}

// GetOrganizationTimeline 获取组织架构时间线
func (oc *OrganizationController) GetOrganizationTimeline(c *gin.Context) {
	params := oc.parseQueryParams(c, []string{"unit_id", "start_date", "end_date", "page", "size"})
	
	result, err := oc.organizationService.GetOrganizationTimeline(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "获取时间线成功", result)
}

// GetChangeStatistics 获取变更统计
func (oc *OrganizationController) GetChangeStatistics(c *gin.Context) {
	params := oc.parseQueryParams(c, []string{"start_date", "end_date"})
	
	stats, err := oc.organizationService.GetChangeStatistics(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "获取变更统计成功", stats)
}

// GetOrganizationEvolution 获取组织架构演变图
func (oc *OrganizationController) GetOrganizationEvolution(c *gin.Context) {
	unitID, err := oc.parseUintParam(c, "id")
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的组织单元ID")
		return
	}
	
	params := oc.parseQueryParams(c, []string{"start_date", "end_date"})
	
	evolution, err := oc.organizationService.GetOrganizationEvolution(unitID, params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "获取组织架构演变图成功", evolution)
}

// ValidateUnitCode 验证组织单元编码唯一性
func (oc *OrganizationController) ValidateUnitCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "请提供组织单元编码")
		return
	}
	
	var excludeID *uint
	if excludeIDStr := c.Query("exclude_id"); excludeIDStr != "" {
		if id, err := strconv.ParseUint(excludeIDStr, 10, 32); err == nil {
			idUint := uint(id)
			excludeID = &idUint
		}
	}
	
	isUnique, err := oc.organizationService.ValidateUnitCode(code, excludeID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "验证成功", gin.H{
		"is_unique": isUnique,
		"code":      code,
	})
}

// GetUnitTypes 获取组织单元类型
func (oc *OrganizationController) GetUnitTypes(c *gin.Context) {
	types := oc.organizationService.GetUnitTypes()
	utils.SuccessResponse(c, http.StatusOK, "获取组织单元类型成功", types)
}

// BatchUpdateUnits 批量更新组织单元
func (oc *OrganizationController) BatchUpdateUnits(c *gin.Context) {
	var req struct {
		Updates []map[string]interface{} `json:"updates" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数无效: "+err.Error())
		return
	}
	
	userID, exists := oc.getUserID(c)
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
		return
	}
	
	err := oc.organizationService.BatchUpdateUnits(req.Updates, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	utils.SuccessResponse(c, http.StatusOK, "批量更新成功", nil)
}

// ExportOrganization 导出组织架构
func (oc *OrganizationController) ExportOrganization(c *gin.Context) {
	format := c.Query("format")
	if format == "" {
		format = "json"
	}
	
	data, err := oc.organizationService.ExportOrganization(format)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	
	// 设置响应头
	switch format {
	case "excel":
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename=organization.xlsx")
	default:
		c.Header("Content-Type", "application/json")
		c.Header("Content-Disposition", "attachment; filename=organization.json")
	}
	
	c.Data(http.StatusOK, c.GetHeader("Content-Type"), data)
}