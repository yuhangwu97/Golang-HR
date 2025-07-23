package controllers

import (
	"gin-project/models"
	"gin-project/services"
	"gin-project/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func NewUserControllerWithService(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// InjectDependencies implements DependencyInjector interface
func (uc *UserController) InjectDependencies(deps ...interface{}) error {
	for _, dep := range deps {
		switch d := dep.(type) {
		case services.UserService:
			uc.userService = d
		}
	}
	return nil
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User
	if !utils.BindAndValidate(c, &user) {
		return
	}

	createdUser, err := uc.userService.CreateUser(&user)
	if err != nil {
		utils.HandleServiceError(c, err, "User")
		return
	}

	utils.SuccessResponse(c, 201, "User created successfully", createdUser)
}

func (uc *UserController) GetUsers(c *gin.Context) {
	users, err := uc.userService.GetAllUsers()
	if err != nil {
		utils.HandleServiceError(c, err, "Users")
		return
	}

	utils.SuccessResponse(c, 200, "Users retrieved successfully", users)
}

func (uc *UserController) GetUser(c *gin.Context) {
	id, ok := utils.GetValidatedParamID(c)
	if !ok {
		return
	}

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	user, err := uc.userService.GetUserByID(uint(userID))
	if err != nil {
		utils.HandleServiceError(c, err, "User")
		return
	}

	utils.SuccessResponse(c, 200, "User retrieved successfully", user)
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	id, ok := utils.GetValidatedParamID(c)
	if !ok {
		return
	}

	var user models.User
	if !utils.BindAndValidate(c, &user) {
		return
	}

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	user.ID = uint(userID)
	updatedUser, err := uc.userService.UpdateUser(&user)
	if err != nil {
		utils.HandleServiceError(c, err, "User")
		return
	}

	utils.SuccessResponse(c, 200, "User updated successfully", updatedUser)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id, ok := utils.GetValidatedParamID(c)
	if !ok {
		return
	}

	userID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, 400, "Invalid user ID")
		return
	}

	err = uc.userService.DeleteUser(uint(userID))
	if err != nil {
		utils.HandleServiceError(c, err, "User")
		return
	}

	utils.SuccessResponse(c, 200, "User deleted successfully", nil)
}
