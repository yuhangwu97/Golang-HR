package services

import (
	"fmt"
	"time"

	"gin-project/models"
	"gin-project/utils"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	repo *utils.MySQLRepository
}

func NewAuthService(db *gorm.DB, redisClient *redis.Client) *AuthService {
	return &AuthService{
		repo: utils.NewMySQLRepository(db, redisClient),
	}
}

// InjectDependencies implements DependencyInjector interface
func (s *AuthService) InjectDependencies(deps ...interface{}) error {
	for _, dep := range deps {
		switch d := dep.(type) {
		case *utils.MySQLRepository:
			s.repo = d
		}
	}
	return nil
}

func (s *AuthService) Register(user *models.User) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user.Password = string(hashedPassword)

	// Set default role will be handled after user creation

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	// Clear password before returning
	user.Password = ""

	// Invalidate users cache
	s.repo.InvalidateCache("users:*")

	return user, nil
}

func (s *AuthService) Login(email, password string) (*models.LoginResponse, error) {
	var user models.User
	err := s.repo.FindOne(&user, "email = ?", email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	// Get user's primary role - for now use "user" as default
	// TODO: implement proper role lookup from UserRole table
	userRole := "admin"
	token, err := s.generateJWT(fmt.Sprintf("%d", user.ID), userRole)
	if err != nil {
		return nil, err
	}

	// Store session in Redis with 24 hour expiration
	sessionKey := fmt.Sprintf("session:%d", user.ID)
	s.repo.RedisSet(sessionKey, token, 24*time.Hour)

	user.Password = ""
	return &models.LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) generateJWT(userID, role string) (string, error) {
	return utils.DefaultJWTService.GenerateToken(userID, role)
}

// ValidateToken validates a JWT token and returns user information
func (s *AuthService) ValidateToken(tokenString string) (*models.User, error) {
	// 验证JWT token
	claims, err := utils.DefaultJWTService.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 从数据库获取用户信息
	var user models.User
	userIDStr := claims.UserID
	err = s.repo.FindByID(&user, userIDStr)
	if err != nil {
		return nil, err
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, fmt.Errorf("user account is not active")
	}

	// 检查Redis中的session是否存在
	sessionKey := fmt.Sprintf("session:%d", user.ID)
	exists, err := s.repo.RedisExists(sessionKey)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("session expired")
	}

	// 清除密码字段
	user.Password = ""
	return &user, nil
}

// RefreshToken refreshes an existing JWT token
func (s *AuthService) RefreshToken(tokenString string) (*models.LoginResponse, error) {
	// 验证当前token
	claims, err := utils.DefaultJWTService.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	var user models.User
	err = s.repo.FindByID(&user, claims.UserID)
	if err != nil {
		return nil, err
	}

	// 检查用户状态
	if user.Status != "active" {
		return nil, fmt.Errorf("user account is not active")
	}

	// 生成新的token
	newToken, err := s.generateJWT(claims.UserID, claims.Role)
	if err != nil {
		return nil, err
	}

	// 更新Redis中的session
	sessionKey := fmt.Sprintf("session:%d", user.ID)
	s.repo.RedisSet(sessionKey, newToken, 24*time.Hour)

	// 清除密码字段
	user.Password = ""
	return &models.LoginResponse{
		Token: newToken,
		User:  user,
	}, nil
}

// Logout logs out a user by removing their session
func (s *AuthService) Logout(userID uint) error {
	sessionKey := fmt.Sprintf("session:%d", userID)
	return s.repo.RedisDelete(sessionKey)
}
