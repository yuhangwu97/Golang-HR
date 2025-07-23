package controllers

import (
	"gin-project/models"
	"gin-project/services"
	"gin-project/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	roleService       services.RoleServiceInterface
	permissionService services.PermissionServiceInterface
}

func NewRoleController(roleService services.RoleServiceInterface, permissionService services.PermissionServiceInterface) *RoleController {
	return &RoleController{
		roleService:       roleService,
		permissionService: permissionService,
	}
}

// InjectDependencies implements DependencyInjector interface
func (rc *RoleController) InjectDependencies(deps ...interface{}) error {
	for _, dep := range deps {
		switch d := dep.(type) {
		case services.RoleServiceInterface:
			rc.roleService = d
		case services.PermissionServiceInterface:
			rc.permissionService = d
		}
	}
	return nil
}

// GetRoles 获取角色列表
func (rc *RoleController) GetRoles(c *gin.Context) {
	var params services.RoleQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// 设置默认分页参数
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.PageSize <= 0 {
		params.PageSize = 20
	}

	response, err := rc.roleService.GetRoles(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取角色列表失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取角色列表成功", response)
}

// GetRole 获取角色详情
func (rc *RoleController) GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	role, err := rc.roleService.GetRoleByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "角色")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取角色详情成功", role)
}

// CreateRole 创建角色
func (rc *RoleController) CreateRole(c *gin.Context) {
	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	createdRole, err := rc.roleService.CreateRole(&role)
	if err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建角色失败")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "创建角色成功", createdRole)
}

// UpdateRole 更新角色
func (rc *RoleController) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var role models.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	role.ID = uint(id)
	updatedRole, err := rc.roleService.UpdateRole(&role)
	if err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新角色失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "更新角色成功", updatedRole)
}

// DeleteRole 删除角色
func (rc *RoleController) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	if err := rc.roleService.DeleteRole(uint(id)); err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除角色失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "删除角色成功", nil)
}

// GetRolePermissions 获取角色权限
func (rc *RoleController) GetRolePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	permissions, err := rc.roleService.GetRolePermissions(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取角色权限失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取角色权限成功", permissions)
}

// AssignPermissions 分配权限给角色
func (rc *RoleController) AssignPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := rc.roleService.AssignPermissions(uint(id), req.PermissionIDs); err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "分配权限失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "分配权限成功", nil)
}

// RemovePermissions 移除角色权限
func (rc *RoleController) RemovePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := rc.roleService.RemovePermissions(uint(id), req.PermissionIDs); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "移除权限失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "移除权限成功", nil)
}

// GetRoleUsers 获取角色用户
func (rc *RoleController) GetRoleUsers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	users, err := rc.roleService.GetRoleUsers(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取角色用户失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取角色用户成功", users)
}

// GetRoleStatistics 获取角色统计信息
func (rc *RoleController) GetRoleStatistics(c *gin.Context) {
	statistics, err := rc.roleService.GetRoleStatistics()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取角色统计信息失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取角色统计信息成功", statistics)
}