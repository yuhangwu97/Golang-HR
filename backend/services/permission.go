package services

import (
	"gin-project/models"
	"gin-project/utils"

	"gorm.io/gorm"
)

type PermissionServiceInterface interface {
	GetPermissions(params PermissionQueryParams) (*utils.PaginationResponse, error)
	GetPermissionByID(id uint) (*models.Permission, error)
	CreatePermission(permission *models.Permission) (*models.Permission, error)
	UpdatePermission(permission *models.Permission) (*models.Permission, error)
	DeletePermission(id uint) error
	GetPermissionsByResource(resource string) ([]*models.Permission, error)
	GetPermissionTree() ([]*PermissionNode, error)
	GetPermissionStatistics() (*PermissionStatistics, error)
}

type PermissionService struct {
	db *gorm.DB
}

type PermissionQueryParams struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
	Status   string `json:"status"`
	Keyword  string `json:"keyword"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type PermissionNode struct {
	Resource    string            `json:"resource"`
	Permissions []*models.Permission `json:"permissions"`
}

type PermissionStatistics struct {
	Total        int64            `json:"total"`
	ByResource   map[string]int64 `json:"by_resource"`
	ByAction     map[string]int64 `json:"by_action"`
	ByStatus     map[string]int64 `json:"by_status"`
	RoleCounts   map[string]int64 `json:"role_counts"`
}

func NewPermissionService(db *gorm.DB) PermissionServiceInterface {
	return &PermissionService{
		db: db,
	}
}

// InjectDependencies implements DependencyInjector interface
func (ps *PermissionService) InjectDependencies(deps ...interface{}) error {
	for _, dep := range deps {
		switch d := dep.(type) {
		case *gorm.DB:
			ps.db = d
		}
	}
	return nil
}

func (ps *PermissionService) GetPermissions(params PermissionQueryParams) (*utils.PaginationResponse, error) {
	var permissions []models.Permission
	var total int64

	query := ps.db.Model(&models.Permission{})

	// 过滤条件
	if params.Resource != "" {
		query = query.Where("resource = ?", params.Resource)
	}

	if params.Action != "" {
		query = query.Where("action = ?", params.Action)
	}

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
	if err := query.Offset(offset).Limit(params.PageSize).
		Order("resource ASC, action ASC").Find(&permissions).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))
	
	return &utils.PaginationResponse{
		Data:       permissions,
		TotalItems: total,
		TotalPages: totalPages,
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}

func (ps *PermissionService) GetPermissionByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	err := ps.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (ps *PermissionService) CreatePermission(permission *models.Permission) (*models.Permission, error) {
	// 检查权限编码是否已存在
	var count int64
	ps.db.Model(&models.Permission{}).Where("code = ?", permission.Code).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("权限编码已存在")
	}

	// 设置默认状态
	if permission.Status == "" {
		permission.Status = "active"
	}

	err := ps.db.Create(permission).Error
	if err != nil {
		return nil, err
	}

	return ps.GetPermissionByID(permission.ID)
}

func (ps *PermissionService) UpdatePermission(permission *models.Permission) (*models.Permission, error) {
	// 检查权限编码是否已被其他权限使用
	var count int64
	ps.db.Model(&models.Permission{}).Where("code = ? AND id != ?", permission.Code, permission.ID).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("权限编码已存在")
	}

	err := ps.db.Model(permission).Updates(permission).Error
	if err != nil {
		return nil, err
	}

	return ps.GetPermissionByID(permission.ID)
}

func (ps *PermissionService) DeletePermission(id uint) error {
	// 检查是否有角色使用此权限
	var count int64
	ps.db.Model(&models.RolePermission{}).Where("permission_id = ?", id).Count(&count)
	if count > 0 {
		return utils.NewValidationError("无法删除已被角色使用的权限")
	}

	return ps.db.Delete(&models.Permission{}, id).Error
}

func (ps *PermissionService) GetPermissionsByResource(resource string) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := ps.db.Where("resource = ? AND status = ?", resource, "active").
		Order("action ASC").Find(&permissions).Error
	return permissions, err
}

func (ps *PermissionService) GetPermissionTree() ([]*PermissionNode, error) {
	var permissions []*models.Permission
	err := ps.db.Where("status = ?", "active").
		Order("resource ASC, action ASC").Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	// 按资源分组
	resourceMap := make(map[string][]*models.Permission)
	for _, permission := range permissions {
		resourceMap[permission.Resource] = append(resourceMap[permission.Resource], permission)
	}

	var tree []*PermissionNode
	for resource, perms := range resourceMap {
		tree = append(tree, &PermissionNode{
			Resource:    resource,
			Permissions: perms,
		})
	}

	return tree, nil
}

func (ps *PermissionService) GetPermissionStatistics() (*PermissionStatistics, error) {
	var total int64
	ps.db.Model(&models.Permission{}).Count(&total)

	// 按资源统计
	var resourceStats []struct {
		Resource string `json:"resource"`
		Count    int64  `json:"count"`
	}
	ps.db.Model(&models.Permission{}).
		Select("resource, COUNT(*) as count").
		Group("resource").
		Scan(&resourceStats)

	byResource := make(map[string]int64)
	for _, stat := range resourceStats {
		byResource[stat.Resource] = stat.Count
	}

	// 按操作统计
	var actionStats []struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}
	ps.db.Model(&models.Permission{}).
		Select("action, COUNT(*) as count").
		Group("action").
		Scan(&actionStats)

	byAction := make(map[string]int64)
	for _, stat := range actionStats {
		byAction[stat.Action] = stat.Count
	}

	// 按状态统计
	var statusStats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	ps.db.Model(&models.Permission{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusStats)

	byStatus := make(map[string]int64)
	for _, stat := range statusStats {
		byStatus[stat.Status] = stat.Count
	}

	// 角色数量统计
	var roleStats []struct {
		PermissionName string `json:"permission_name"`
		Count          int64  `json:"count"`
	}
	ps.db.Model(&models.RolePermission{}).
		Select("permissions.name as permission_name, COUNT(*) as count").
		Joins("LEFT JOIN permissions ON role_permissions.permission_id = permissions.id").
		Group("permissions.name").
		Scan(&roleStats)

	roleCounts := make(map[string]int64)
	for _, stat := range roleStats {
		roleCounts[stat.PermissionName] = stat.Count
	}

	return &PermissionStatistics{
		Total:      total,
		ByResource: byResource,
		ByAction:   byAction,
		ByStatus:   byStatus,
		RoleCounts: roleCounts,
	}, nil
}