package controllers

import (
	"gin-project/models"
	"gin-project/services"
	"gin-project/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	permissionService services.PermissionServiceInterface
}

func NewPermissionController(permissionService services.PermissionServiceInterface) *PermissionController {
	return &PermissionController{
		permissionService: permissionService,
	}
}

// InjectDependencies implements DependencyInjector interface
func (pc *PermissionController) InjectDependencies(deps ...interface{}) error {
	for _, dep := range deps {
		switch d := dep.(type) {
		case services.PermissionServiceInterface:
			pc.permissionService = d
		}
	}
	return nil
}

// GetPermissions 获取权限列表
func (pc *PermissionController) GetPermissions(c *gin.Context) {
	var params services.PermissionQueryParams
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

	response, err := pc.permissionService.GetPermissions(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取权限列表失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取权限列表成功", response)
}

// GetPermission 获取权限详情
func (pc *PermissionController) GetPermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的权限ID")
		return
	}

	permission, err := pc.permissionService.GetPermissionByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "权限")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取权限详情成功", permission)
}

// CreatePermission 创建权限
func (pc *PermissionController) CreatePermission(c *gin.Context) {
	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	createdPermission, err := pc.permissionService.CreatePermission(&permission)
	if err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建权限失败")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "创建权限成功", createdPermission)
}

// UpdatePermission 更新权限
func (pc *PermissionController) UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的权限ID")
		return
	}

	var permission models.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	permission.ID = uint(id)
	updatedPermission, err := pc.permissionService.UpdatePermission(&permission)
	if err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新权限失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "更新权限成功", updatedPermission)
}

// DeletePermission 删除权限
func (pc *PermissionController) DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的权限ID")
		return
	}

	if err := pc.permissionService.DeletePermission(uint(id)); err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除权限失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "删除权限成功", nil)
}

// GetPermissionsByResource 按资源获取权限
func (pc *PermissionController) GetPermissionsByResource(c *gin.Context) {
	resource := c.Query("resource")
	if resource == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "资源参数不能为空")
		return
	}

	permissions, err := pc.permissionService.GetPermissionsByResource(resource)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取权限失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取权限成功", permissions)
}

// GetPermissionTree 获取权限树
func (pc *PermissionController) GetPermissionTree(c *gin.Context) {
	tree, err := pc.permissionService.GetPermissionTree()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取权限树失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取权限树成功", tree)
}

// GetPermissionStatistics 获取权限统计信息
func (pc *PermissionController) GetPermissionStatistics(c *gin.Context) {
	statistics, err := pc.permissionService.GetPermissionStatistics()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取权限统计信息失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取权限统计信息成功", statistics)
}