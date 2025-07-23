package controllers

import (
	"net/http"
	"strconv"

	"gin-project/models"
	"gin-project/services"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

type PositionController struct {
	positionService services.PositionServiceInterface
}

func NewPositionController(positionService services.PositionServiceInterface) *PositionController {
	return &PositionController{
		positionService: positionService,
	}
}

// GetPositions 获取职位列表
func (pc *PositionController) GetPositions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	departmentID, _ := strconv.ParseUint(c.Query("department_id"), 10, 32)
	status := c.Query("status")
	keyword := c.Query("keyword")

	params := services.PositionQueryParams{
		DepartmentID: uint(departmentID),
		Status:       status,
		Keyword:      keyword,
		Page:         page,
		PageSize:     pageSize,
	}

	result, err := pc.positionService.GetPositions(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取职位列表失败")
		return
	}

	utils.SuccessResponse(c, 200, "获取职位列表成功", result)
}

// GetPosition 获取职位详情
func (pc *PositionController) GetPosition(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的职位ID")
		return
	}

	position, err := pc.positionService.GetPositionByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "职位不存在")
		return
	}

	utils.SuccessResponse(c, 200, "获取职位详情成功", position)
}

// CreatePosition 创建职位
func (pc *PositionController) CreatePosition(c *gin.Context) {
	var req models.Position
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	position, err := pc.positionService.CreatePosition(&req)
	if err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建职位失败")
		return
	}

	utils.SuccessResponse(c, 201, "创建职位成功", position)
}

// UpdatePosition 更新职位
func (pc *PositionController) UpdatePosition(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的职位ID")
		return
	}

	var req models.Position
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	req.ID = uint(id)
	position, err := pc.positionService.UpdatePosition(&req)
	if err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新职位失败")
		return
	}

	utils.SuccessResponse(c, 200, "更新职位成功", position)
}

// DeletePosition 删除职位
func (pc *PositionController) DeletePosition(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的职位ID")
		return
	}

	err = pc.positionService.DeletePosition(uint(id))
	if err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除职位失败")
		return
	}

	utils.SuccessResponse(c, 200, "删除职位成功", nil)
}

// GetPositionsByDepartment 获取部门下的职位
func (pc *PositionController) GetPositionsByDepartment(c *gin.Context) {
	departmentID, err := strconv.ParseUint(c.Param("departmentId"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的部门ID")
		return
	}

	positions, err := pc.positionService.GetPositionsByDepartment(uint(departmentID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取部门职位失败")
		return
	}

	utils.SuccessResponse(c, 200, "获取部门职位成功", positions)
}

// GetPositionStatistics 获取职位统计
func (pc *PositionController) GetPositionStatistics(c *gin.Context) {
	stats, err := pc.positionService.GetPositionStatistics()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取职位统计失败")
		return
	}

	utils.SuccessResponse(c, 200, "获取职位统计成功", stats)
}

// SearchPositions 搜索职位
func (pc *PositionController) SearchPositions(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "搜索关键词不能为空")
		return
	}

	positions, err := pc.positionService.SearchPositions(query)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "搜索职位失败")
		return
	}

	utils.SuccessResponse(c, 200, "搜索职位成功", positions)
}

// BulkCreatePositions 批量创建职位
func (pc *PositionController) BulkCreatePositions(c *gin.Context) {
	var req struct {
		Positions []*models.Position `json:"positions" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	result, err := pc.positionService.BulkCreatePositions(req.Positions)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "批量创建职位失败")
		return
	}

	utils.SuccessResponse(c, 200, "批量创建职位完成", result)
}

// GetPositionTree 获取职位树形结构
func (pc *PositionController) GetPositionTree(c *gin.Context) {
	tree, err := pc.positionService.GetPositionTree()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取职位树失败")
		return
	}

	utils.SuccessResponse(c, 200, "获取职位树成功", tree)
}

// GetAllPositions 获取所有职位（用于父级选择）
func (pc *PositionController) GetAllPositions(c *gin.Context) {
	positions, err := pc.positionService.GetAllPositions()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取所有职位失败")
		return
	}

	utils.SuccessResponse(c, 200, "获取所有职位成功", positions)
}
