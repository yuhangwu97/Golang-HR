package services

import (
	"gin-project/models"
	"gin-project/utils"

	"gorm.io/gorm"
)

type RoleServiceInterface interface {
	GetRoles(params RoleQueryParams) (*utils.PaginationResponse, error)
	GetRoleByID(id uint) (*models.Role, error)
	CreateRole(role *models.Role) (*models.Role, error)
	UpdateRole(role *models.Role) (*models.Role, error)
	DeleteRole(id uint) error
	GetRolePermissions(roleID uint) ([]*models.Permission, error)
	AssignPermissions(roleID uint, permissionIDs []uint) error
	RemovePermissions(roleID uint, permissionIDs []uint) error
	GetRoleUsers(roleID uint) ([]*models.User, error)
	GetRoleStatistics() (*RoleStatistics, error)
}

type RoleService struct {
	db *gorm.DB
}

type RoleQueryParams struct {
	Status   string `json:"status"`
	Keyword  string `json:"keyword"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type RoleStatistics struct {
	Total            int64            `json:"total"`
	ByStatus         map[string]int64 `json:"by_status"`
	UserCounts       map[string]int64 `json:"user_counts"`
	PermissionCounts map[string]int64 `json:"permission_counts"`
}

func NewRoleService(db *gorm.DB) RoleServiceInterface {
	return &RoleService{
		db: db,
	}
}

// InjectDependencies implements DependencyInjector interface
func (rs *RoleService) InjectDependencies(deps ...interface{}) error {
	for _, dep := range deps {
		switch d := dep.(type) {
		case *gorm.DB:
			rs.db = d
		}
	}
	return nil
}

func (rs *RoleService) GetRoles(params RoleQueryParams) (*utils.PaginationResponse, error) {
	var roles []models.Role
	var total int64

	query := rs.db.Model(&models.Role{})

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
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&roles).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))
	
	return &utils.PaginationResponse{
		Data:       roles,
		TotalItems: total,
		TotalPages: totalPages,
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}

func (rs *RoleService) GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	err := rs.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (rs *RoleService) CreateRole(role *models.Role) (*models.Role, error) {
	// 检查角色编码是否已存在
	var count int64
	rs.db.Model(&models.Role{}).Where("code = ?", role.Code).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("角色编码已存在")
	}

	// 检查角色名称是否已存在
	rs.db.Model(&models.Role{}).Where("name = ?", role.Name).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("角色名称已存在")
	}

	// 设置默认状态
	if role.Status == "" {
		role.Status = "active"
	}

	err := rs.db.Create(role).Error
	if err != nil {
		return nil, err
	}

	return rs.GetRoleByID(role.ID)
}

func (rs *RoleService) UpdateRole(role *models.Role) (*models.Role, error) {
	// 检查角色编码是否已被其他角色使用
	var count int64
	rs.db.Model(&models.Role{}).Where("code = ? AND id != ?", role.Code, role.ID).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("角色编码已存在")
	}

	// 检查角色名称是否已被其他角色使用
	rs.db.Model(&models.Role{}).Where("name = ? AND id != ?", role.Name, role.ID).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("角色名称已存在")
	}

	err := rs.db.Model(role).Updates(role).Error
	if err != nil {
		return nil, err
	}

	return rs.GetRoleByID(role.ID)
}

func (rs *RoleService) DeleteRole(id uint) error {
	// 检查是否有用户使用此角色
	var count int64
	rs.db.Model(&models.UserRole{}).Where("role_id = ?", id).Count(&count)
	if count > 0 {
		return utils.NewValidationError("无法删除已被用户使用的角色")
	}

	// 删除角色权限关联
	rs.db.Where("role_id = ?", id).Delete(&models.RolePermission{})

	return rs.db.Delete(&models.Role{}, id).Error
}

func (rs *RoleService) GetRolePermissions(roleID uint) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := rs.db.Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ? AND permissions.status = ?", roleID, "active").
		Find(&permissions).Error
	return permissions, err
}

func (rs *RoleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	// 检查角色是否存在
	var role models.Role
	if err := rs.db.First(&role, roleID).Error; err != nil {
		return utils.NewValidationError("角色不存在")
	}

	// 删除现有权限关联
	rs.db.Where("role_id = ?", roleID).Delete(&models.RolePermission{})

	// 添加新的权限关联
	for _, permissionID := range permissionIDs {
		rolePermission := &models.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
		}
		rs.db.Create(rolePermission)
	}

	return nil
}

func (rs *RoleService) RemovePermissions(roleID uint, permissionIDs []uint) error {
	return rs.db.Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).
		Delete(&models.RolePermission{}).Error
}

func (rs *RoleService) GetRoleUsers(roleID uint) ([]*models.User, error) {
	var users []*models.User
	err := rs.db.Table("users").
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ? AND users.status = ?", roleID, "active").
		Find(&users).Error
	return users, err
}

func (rs *RoleService) GetRoleStatistics() (*RoleStatistics, error) {
	var total int64
	rs.db.Model(&models.Role{}).Count(&total)

	// 按状态统计
	var statusStats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	rs.db.Model(&models.Role{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusStats)

	byStatus := make(map[string]int64)
	for _, stat := range statusStats {
		byStatus[stat.Status] = stat.Count
	}

	// 用户数量统计
	var userStats []struct {
		RoleName string `json:"role_name"`
		Count    int64  `json:"count"`
	}
	rs.db.Model(&models.UserRole{}).
		Select("roles.name as role_name, COUNT(*) as count").
		Joins("LEFT JOIN roles ON user_roles.role_id = roles.id").
		Group("roles.name").
		Scan(&userStats)

	userCounts := make(map[string]int64)
	for _, stat := range userStats {
		userCounts[stat.RoleName] = stat.Count
	}

	// 权限数量统计
	var permissionStats []struct {
		RoleName string `json:"role_name"`
		Count    int64  `json:"count"`
	}
	rs.db.Model(&models.RolePermission{}).
		Select("roles.name as role_name, COUNT(*) as count").
		Joins("LEFT JOIN roles ON role_permissions.role_id = roles.id").
		Group("roles.name").
		Scan(&permissionStats)

	permissionCounts := make(map[string]int64)
	for _, stat := range permissionStats {
		permissionCounts[stat.RoleName] = stat.Count
	}

	return &RoleStatistics{
		Total:            total,
		ByStatus:         byStatus,
		UserCounts:       userCounts,
		PermissionCounts: permissionCounts,
	}, nil
}