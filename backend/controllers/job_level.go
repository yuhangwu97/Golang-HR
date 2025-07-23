package controllers

import (
	"net/http"
	"strconv"

	"gin-project/models"
	"gin-project/services"
	"gin-project/utils"

	"github.com/gin-gonic/gin"
)

type JobLevelController struct {
	jobLevelService services.JobLevelServiceInterface
}

func NewJobLevelController(jobLevelService services.JobLevelServiceInterface) *JobLevelController {
	return &JobLevelController{
		jobLevelService: jobLevelService,
	}
}

// GetJobLevels 获取职级列表
func (jlc *JobLevelController) GetJobLevels(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	params := services.JobLevelQueryParams{
		Status:   status,
		Keyword:  keyword,
		Page:     page,
		PageSize: pageSize,
	}

	result, err := jlc.jobLevelService.GetJobLevels(params)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取职级列表失败")
		return
	}

	utils.SuccessResponse(c, 200, "获取职级列表成功", result)
}

// GetJobLevel 获取职级详情
func (jlc *JobLevelController) GetJobLevel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的职级ID")
		return
	}

	jobLevel, err := jlc.jobLevelService.GetJobLevelByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "职级不存在")
		return
	}

	utils.SuccessResponse(c, 200, "获取职级详情成功", jobLevel)
}

// CreateJobLevel 创建职级
func (jlc *JobLevelController) CreateJobLevel(c *gin.Context) {
	var req models.JobLevel
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	jobLevel, err := jlc.jobLevelService.CreateJobLevel(&req)
	if err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "创建职级失败")
		return
	}

	utils.SuccessResponse(c, 201, "创建职级成功", jobLevel)
}

// UpdateJobLevel 更新职级
func (jlc *JobLevelController) UpdateJobLevel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的职级ID")
		return
	}

	var req models.JobLevel
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	req.ID = uint(id)
	jobLevel, err := jlc.jobLevelService.UpdateJobLevel(&req)
	if err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "更新职级失败")
		return
	}

	utils.SuccessResponse(c, 200, "更新职级成功", jobLevel)
}

// DeleteJobLevel 删除职级
func (jlc *JobLevelController) DeleteJobLevel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的职级ID")
		return
	}

	err = jlc.jobLevelService.DeleteJobLevel(uint(id))
	if err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "删除职级失败")
		return
	}

	utils.SuccessResponse(c, 200, "删除职级成功", nil)
}

// GetJobLevelsByLevel 按等级获取职级
func (jlc *JobLevelController) GetJobLevelsByLevel(c *gin.Context) {
	jobLevels, err := jlc.jobLevelService.GetJobLevelsByLevel()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取职级失败")
		return
	}

	utils.SuccessResponse(c, 200, "获取职级成功", jobLevels)
}

// GetJobLevelStatistics 获取职级统计
func (jlc *JobLevelController) GetJobLevelStatistics(c *gin.Context) {
	stats, err := jlc.jobLevelService.GetJobLevelStatistics()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "获取职级统计失败")
		return
	}

	utils.SuccessResponse(c, 200, "获取职级统计成功", stats)
}

// ValidateSalaryRange 验证薪资范围
func (jlc *JobLevelController) ValidateSalaryRange(c *gin.Context) {
	var req struct {
		JobLevelID uint    `json:"job_level_id" binding:"required"`
		Salary     float64 `json:"salary" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "请求参数错误")
		return
	}

	err := jlc.jobLevelService.ValidateSalaryRange(req.JobLevelID, req.Salary)
	if err != nil {
		if validationErr, ok := err.(*utils.ValidationError); ok {
			utils.ErrorResponse(c, http.StatusBadRequest, validationErr.Message)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "验证薪资范围失败")
		return
	}

	utils.SuccessResponse(c, 200, "薪资范围验证通过", nil)
}

// GetSalaryRangeByLevel 根据职级等级获取薪资范围
func (jlc *JobLevelController) GetSalaryRangeByLevel(c *gin.Context) {
	level, err := strconv.Atoi(c.Param("level"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "无效的职级等级")
		return
	}

	salaryRange, err := jlc.jobLevelService.GetSalaryRangeByLevel(level)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "职级不存在")
		return
	}

	utils.SuccessResponse(c, 200, "获取薪资范围成功", salaryRange)
}
