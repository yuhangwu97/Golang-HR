package controllers

import (
	"gin-project/models"
	"gin-project/services"
	"gin-project/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SystemUserController struct {
	userService services.UserServiceInterface
}

func NewSystemUserController(userService services.UserServiceInterface) *SystemUserController {
	return &SystemUserController{
		userService: userService,
	}
}

// InjectDependencies implements DependencyInjector interface
func (suc *SystemUserController) InjectDependencies(deps ...interface{}) error {
	for _, dep := range deps {
		switch d := dep.(type) {
		case services.UserServiceInterface:
			suc.userService = d
		}
	}
	return nil
}

// GetUsers 获取用户列表
func (suc *SystemUserController) GetUsers(c *gin.Context) {
	var params services.UserQueryParams
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

	response, err := suc.userService.GetUsers(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取用户列表失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取用户列表成功", response)
}

// GetUser 获取用户详情
func (suc *SystemUserController) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	user, err := suc.userService.GetUserByID(uint(id))
	if err != nil {
		utils.NotFoundResponse(c, "用户")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取用户详情成功", user)
}

// CreateUser 创建用户
func (suc *SystemUserController) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	createdUser, err := suc.userService.CreateUser(&user)
	if err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建用户失败")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "创建用户成功", createdUser)
}

// UpdateUser 更新用户
func (suc *SystemUserController) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	user.ID = uint(id)
	updatedUser, err := suc.userService.UpdateUser(&user)
	if err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新用户失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "更新用户成功", updatedUser)
}

// DeleteUser 删除用户
func (suc *SystemUserController) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	if err := suc.userService.DeleteUser(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除用户失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "删除用户成功", nil)
}

// ChangePassword 修改密码
func (suc *SystemUserController) ChangePassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := suc.userService.ChangePassword(uint(id), req.OldPassword, req.NewPassword); err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "修改密码失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "修改密码成功", nil)
}

// ResetPassword 重置密码
func (suc *SystemUserController) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := suc.userService.ResetPassword(uint(id), req.NewPassword); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "重置密码失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "重置密码成功", nil)
}

// GetUserRoles 获取用户角色
func (suc *SystemUserController) GetUserRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	roles, err := suc.userService.GetUserRoles(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取用户角色失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取用户角色成功", roles)
}

// AssignRoles 分配角色给用户
func (suc *SystemUserController) AssignRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req struct {
		RoleIDs []uint `json:"role_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := suc.userService.AssignRoles(uint(id), req.RoleIDs); err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "分配角色失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "分配角色成功", nil)
}

// RemoveRoles 移除用户角色
func (suc *SystemUserController) RemoveRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req struct {
		RoleIDs []uint `json:"role_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	if err := suc.userService.RemoveRoles(uint(id), req.RoleIDs); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "移除角色失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "移除角色成功", nil)
}

// GetUserPermissions 获取用户权限
func (suc *SystemUserController) GetUserPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的用户ID")
		return
	}

	permissions, err := suc.userService.GetUserPermissions(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取用户权限失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取用户权限成功", permissions)
}

// GetUserStatistics 获取用户统计信息
func (suc *SystemUserController) GetUserStatistics(c *gin.Context) {
	statistics, err := suc.userService.GetUserStatistics()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取用户统计信息失败")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "获取用户统计信息成功", statistics)
}