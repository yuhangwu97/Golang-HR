package controllers

import (
	"strconv"
	"strings"

	"gin-project/models"
	"gin-project/services"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func NewAuthControllerWithService(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// InjectDependencies implements DependencyInjector interface
func (ac *AuthController) InjectDependencies(deps ...interface{}) error {
	for _, dep := range deps {
		switch d := dep.(type) {
		case *services.AuthService:
			ac.authService = d
		}
	}
	return nil
}

func (ac *AuthController) Register(c *gin.Context) {
	var user models.User
	if !utils.BindAndValidate(c, &user) {
		return
	}

	createdUser, err := ac.authService.Register(&user)
	if err != nil {
		utils.HandleServiceError(c, err, "User")
		return
	}

	utils.SuccessResponse(c, 201, "User registered successfully", createdUser)
}

func (ac *AuthController) Login(c *gin.Context) {
	var loginReq models.LoginRequest
	if !utils.BindAndValidate(c, &loginReq) {
		return
	}

	loginResp, err := ac.authService.Login(loginReq.Email, loginReq.Password)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			utils.UnauthorizedResponse(c, "Invalid credentials")
			return
		}
		utils.HandleServiceError(c, err, "Login")
		return
	}

	utils.SuccessResponse(c, 200, "Login successful", loginResp)
}

// ValidateToken validates the JWT token and returns user information
func (ac *AuthController) ValidateToken(c *gin.Context) {
	// 从Header中获取token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.UnauthorizedResponse(c, "Authorization header required")
		return
	}

	// 检查Bearer前缀
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		utils.UnauthorizedResponse(c, "Invalid authorization header format")
		return
	}

	// 提取token
	token := authHeader[len(bearerPrefix):]
	if token == "" {
		utils.UnauthorizedResponse(c, "Token required")
		return
	}

	// 验证token
	user, err := ac.authService.ValidateToken(token)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid or expired token")
		return
	}

	utils.SuccessResponse(c, 200, "Token is valid", gin.H{
		"user":  user,
		"valid": true,
	})
}

// RefreshToken refreshes an existing JWT token
func (ac *AuthController) RefreshToken(c *gin.Context) {
	// 从Header中获取token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.UnauthorizedResponse(c, "Authorization header required")
		return
	}

	// 检查Bearer前缀
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		utils.UnauthorizedResponse(c, "Invalid authorization header format")
		return
	}

	// 提取token
	token := authHeader[len(bearerPrefix):]
	if token == "" {
		utils.UnauthorizedResponse(c, "Token required")
		return
	}

	// 刷新token
	loginResp, err := ac.authService.RefreshToken(token)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid or expired token")
		return
	}

	utils.SuccessResponse(c, 200, "Token refreshed successfully", loginResp)
}

// Logout logs out the current user
func (ac *AuthController) Logout(c *gin.Context) {
	// 从Header中获取token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		utils.UnauthorizedResponse(c, "Authorization header required")
		return
	}

	// 检查Bearer前缀
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		utils.UnauthorizedResponse(c, "Invalid authorization header format")
		return
	}

	// 提取token
	token := authHeader[len(bearerPrefix):]
	if token == "" {
		utils.UnauthorizedResponse(c, "Token required")
		return
	}

	// 验证token并获取用户ID
	claims, err := utils.DefaultJWTService.ValidateToken(token)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid token")
		return
	}

	// 解析用户ID
	userID, err := strconv.ParseUint(claims.UserID, 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid user ID in token")
		return
	}

	// 执行登出
	err = ac.authService.Logout(uint(userID))
	if err != nil {
		utils.HandleServiceError(c, err, "Logout")
		return
	}

	utils.SuccessResponse(c, 200, "Logout successful", nil)
}
