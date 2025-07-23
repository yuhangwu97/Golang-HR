package services

import (
	"gin-project/models"
	"gin-project/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserServiceInterface interface {
	GetUsers(params UserQueryParams) (*utils.PaginationResponse, error)
	GetUserByID(id uint) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id uint) error
	ChangePassword(userID uint, oldPassword, newPassword string) error
	ResetPassword(userID uint, newPassword string) error
	GetUserRoles(userID uint) ([]*models.Role, error)
	AssignRoles(userID uint, roleIDs []uint) error
	RemoveRoles(userID uint, roleIDs []uint) error
	GetUserPermissions(userID uint) ([]*models.Permission, error)
	HasPermission(userID uint, resource, action string) (bool, error)
	GetUserStatistics() (*UserStatistics, error)
	UpdateLastLogin(userID uint) error
}

type UserService struct {
	db *gorm.DB
}

type UserQueryParams struct {
	Status   string `json:"status"`
	Keyword  string `json:"keyword"`
	RoleID   uint   `json:"role_id"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type UserStatistics struct {
	Total       int64            `json:"total"`
	ByStatus    map[string]int64 `json:"by_status"`
	RoleCounts  map[string]int64 `json:"role_counts"`
	LoginCounts map[string]int64 `json:"login_counts"`
}

func NewUserService(db *gorm.DB) UserServiceInterface {
	return &UserService{
		db: db,
	}
}

// InjectDependencies implements DependencyInjector interface
func (us *UserService) InjectDependencies(deps ...interface{}) error {
	for _, dep := range deps {
		switch d := dep.(type) {
		case *gorm.DB:
			us.db = d
		}
	}
	return nil
}

func (us *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := us.db.Preload("Employee").Find(&users).Error
	return users, err
}

func (us *UserService) GetUsers(params UserQueryParams) (*utils.PaginationResponse, error) {
	var users []models.User
	var total int64

	query := us.db.Model(&models.User{}).Preload("Employee")

	// 过滤条件
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.RoleID > 0 {
		query = query.Joins("JOIN user_roles ON users.id = user_roles.user_id").
			Where("user_roles.role_id = ?", params.RoleID)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))
	
	return &utils.PaginationResponse{
		Data:       users,
		TotalItems: total,
		TotalPages: totalPages,
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}

func (us *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := us.db.Preload("Employee").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) CreateUser(user *models.User) (*models.User, error) {
	// 检查用户名是否已存在
	var count int64
	us.db.Model(&models.User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("用户名已存在")
	}

	// 检查邮箱是否已存在
	us.db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("邮箱已存在")
	}

	// 加密密码
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	// 设置默认状态
	if user.Status == "" {
		user.Status = "active"
	}

	err := us.db.Create(user).Error
	if err != nil {
		return nil, err
	}

	return us.GetUserByID(user.ID)
}

func (us *UserService) UpdateUser(user *models.User) (*models.User, error) {
	// 检查用户名是否已被其他用户使用
	var count int64
	us.db.Model(&models.User{}).Where("username = ? AND id != ?", user.Username, user.ID).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("用户名已存在")
	}

	// 检查邮箱是否已被其他用户使用
	us.db.Model(&models.User{}).Where("email = ? AND id != ?", user.Email, user.ID).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("邮箱已存在")
	}

	// 如果密码不为空，加密密码
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	err := us.db.Model(user).Updates(user).Error
	if err != nil {
		return nil, err
	}

	return us.GetUserByID(user.ID)
}

func (us *UserService) DeleteUser(id uint) error {
	// 删除用户角色关联
	us.db.Where("user_id = ?", id).Delete(&models.UserRole{})

	return us.db.Delete(&models.User{}, id).Error
}

func (us *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	var user models.User
	if err := us.db.First(&user, userID).Error; err != nil {
		return utils.NewValidationError("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return utils.NewValidationError("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return us.db.Model(&user).Update("password", string(hashedPassword)).Error
}

func (us *UserService) ResetPassword(userID uint, newPassword string) error {
	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return us.db.Model(&models.User{}).Where("id = ?", userID).Update("password", string(hashedPassword)).Error
}

func (us *UserService) GetUserRoles(userID uint) ([]*models.Role, error) {
	var roles []*models.Role
	err := us.db.Table("roles").
		Joins("JOIN user_roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ? AND roles.status = ?", userID, "active").
		Find(&roles).Error
	return roles, err
}

func (us *UserService) AssignRoles(userID uint, roleIDs []uint) error {
	// 检查用户是否存在
	var user models.User
	if err := us.db.First(&user, userID).Error; err != nil {
		return utils.NewValidationError("用户不存在")
	}

	// 删除现有角色关联
	us.db.Where("user_id = ?", userID).Delete(&models.UserRole{})

	// 添加新的角色关联
	for _, roleID := range roleIDs {
		userRole := &models.UserRole{
			UserID: userID,
			RoleID: roleID,
		}
		us.db.Create(userRole)
	}

	return nil
}

func (us *UserService) RemoveRoles(userID uint, roleIDs []uint) error {
	return us.db.Where("user_id = ? AND role_id IN ?", userID, roleIDs).
		Delete(&models.UserRole{}).Error
}

func (us *UserService) GetUserPermissions(userID uint) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := us.db.Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.status = ?", userID, "active").
		Distinct().
		Find(&permissions).Error
	return permissions, err
}

func (us *UserService) HasPermission(userID uint, resource, action string) (bool, error) {
	var count int64
	err := us.db.Table("permissions").
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.resource = ? AND permissions.action = ? AND permissions.status = ?",
			userID, resource, action, "active").
		Count(&count).Error
	
	return count > 0, err
}

func (us *UserService) GetUserStatistics() (*UserStatistics, error) {
	var total int64
	us.db.Model(&models.User{}).Count(&total)

	// 按状态统计
	var statusStats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	us.db.Model(&models.User{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusStats)

	byStatus := make(map[string]int64)
	for _, stat := range statusStats {
		byStatus[stat.Status] = stat.Count
	}

	// 角色数量统计
	var roleStats []struct {
		RoleName string `json:"role_name"`
		Count    int64  `json:"count"`
	}
	us.db.Model(&models.UserRole{}).
		Select("roles.name as role_name, COUNT(*) as count").
		Joins("LEFT JOIN roles ON user_roles.role_id = roles.id").
		Group("roles.name").
		Scan(&roleStats)

	roleCounts := make(map[string]int64)
	for _, stat := range roleStats {
		roleCounts[stat.RoleName] = stat.Count
	}

	// 登录统计（简化版，实际应该根据具体需求实现）
	loginCounts := make(map[string]int64)
	loginCounts["total_logins"] = 0 // 这里应该根据登录日志表统计

	return &UserStatistics{
		Total:       total,
		ByStatus:    byStatus,
		RoleCounts:  roleCounts,
		LoginCounts: loginCounts,
	}, nil
}

func (us *UserService) UpdateLastLogin(userID uint) error {
	now := time.Now()
	return us.db.Model(&models.User{}).Where("id = ?", userID).
		Update("last_login_at", &now).Error
}