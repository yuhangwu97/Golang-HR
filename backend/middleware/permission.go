package middleware

import (
	"gin-project/config"
	"gin-project/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RequirePermission middleware checks if the user has the required permission
func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by auth middleware)
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
			c.Abort()
			return
		}

		userIDStr, ok := userIDInterface.(string)
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, "无效的用户ID")
			c.Abort()
			return
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "无效的用户ID格式")
			c.Abort()
			return
		}

		// Get user service from container
		userService := config.GetContainer().UserService()
		hasPermission, err := userService.HasPermission(uint(userID), resource, action)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "权限检查失败")
			c.Abort()
			return
		}

		if !hasPermission {
			utils.ErrorResponse(c, http.StatusForbidden, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole middleware checks if the user has one of the required roles
func RequireRole(roleCodes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
			c.Abort()
			return
		}

		userIDStr, ok := userIDInterface.(string)
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, "无效的用户ID")
			c.Abort()
			return
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "无效的用户ID格式")
			c.Abort()
			return
		}

		// Get user service from container
		userService := config.GetContainer().UserService()
		roles, err := userService.GetUserRoles(uint(userID))
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "角色检查失败")
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, role := range roles {
			if role.Code == "super_admin" {
				hasRole = true
				break
			}
			for _, requiredCode := range roleCodes {
				if role.Code == requiredCode {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			utils.ErrorResponse(c, http.StatusForbidden, "角色权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAdmin middleware checks if the user is an admin
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// RequireHR middleware checks if the user is HR or admin
func RequireHR() gin.HandlerFunc {
	return RequireRole("admin", "hr")
}

// CheckResourceAccess middleware checks if user can access a specific resource by ID
func CheckResourceAccess(resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "用户未认证")
			c.Abort()
			return
		}

		userIDStr, ok := userIDInterface.(string)
		if !ok {
			utils.ErrorResponse(c, http.StatusUnauthorized, "无效的用户ID")
			c.Abort()
			return
		}

		userID, err := strconv.ParseUint(userIDStr, 10, 32)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "无效的用户ID格式")
			c.Abort()
			return
		}

		// Get resource ID from URL parameter
		resourceIDStr := c.Param("id")
		resourceID, err := strconv.ParseUint(resourceIDStr, 10, 32)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "无效的资源ID")
			c.Abort()
			return
		}

		// For certain resources, check ownership
		switch resourceType {
		case "employee":
			// Check if user can access this employee record
			if !canAccessEmployee(uint(userID), uint(resourceID)) {
				utils.ErrorResponse(c, http.StatusForbidden, "无权访问此员工信息")
				c.Abort()
				return
			}
		case "salary":
			// Check if user can access this salary record
			if !canAccessSalary(uint(userID), uint(resourceID)) {
				utils.ErrorResponse(c, http.StatusForbidden, "无权访问此薪资信息")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// Helper functions for resource access control
func canAccessEmployee(userID, employeeID uint) bool {
	// Get user and check if they can access this employee
	userService := config.GetContainer().UserService()
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return false
	}

	// Admin can access all employees
	roles, err := userService.GetUserRoles(userID)
	if err != nil {
		return false
	}

	for _, role := range roles {
		if role.Code == "admin" || role.Code == "hr" {
			return true
		}
	}

	// User can only access their own employee record
	return user.EmployeeID != nil && *user.EmployeeID == employeeID
}

func canAccessSalary(userID, salaryID uint) bool {
	// Get user service
	userService := config.GetContainer().UserService()

	// Check if user has permission to view all salaries
	hasPermission, err := userService.HasPermission(userID, "salary", "read")
	if err != nil {
		return false
	}

	if hasPermission {
		return true
	}

	// Check if this is the user's own salary record
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return false
	}

	if user.EmployeeID == nil {
		return false
	}

	// Get salary service to check ownership
	salaryService := config.GetContainer().SalaryService()
	salary, err := salaryService.GetSalaryByID(salaryID)
	if err != nil {
		return false
	}

	return salary.EmployeeID == *user.EmployeeID
}
