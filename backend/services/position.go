package services

import (
	"gin-project/models"
	"gin-project/utils"

	"gorm.io/gorm"
)

type PositionServiceInterface interface {
	GetPositions(params PositionQueryParams) (*utils.PaginationResponse, error)
	GetPositionByID(id uint) (*models.Position, error)
	CreatePosition(position *models.Position) (*models.Position, error)
	UpdatePosition(position *models.Position) (*models.Position, error)
	DeletePosition(id uint) error
	GetPositionsByDepartment(departmentID uint) ([]*models.Position, error)
	GetPositionStatistics() (*PositionStatistics, error)
	SearchPositions(query string) ([]*models.Position, error)
	BulkCreatePositions(positions []*models.Position) (*BulkResult, error)
	GetPositionTree() ([]*PositionTreeNode, error)
	GetAllPositions() ([]*models.Position, error)
}

type PositionService struct {
	db *gorm.DB
}

type PositionQueryParams struct {
	DepartmentID uint   `json:"department_id"`
	Status       string `json:"status"`
	Keyword      string `json:"keyword"`
	Page         int    `json:"page"`
	PageSize     int    `json:"page_size"`
}

type PositionStatistics struct {
	Total             int64            `json:"total"`
	ByDepartment      map[string]int64 `json:"by_department"`
	ByStatus          map[string]int64 `json:"by_status"`
	ActivePositions   int64            `json:"active_positions"`
	InactivePositions int64            `json:"inactive_positions"`
}

type PositionTreeNode struct {
	*models.Position
	Children []*PositionTreeNode `json:"children"`
}

func NewPositionService(db *gorm.DB) PositionServiceInterface {
	return &PositionService{
		db: db,
	}
}

// InjectDependencies implements DependencyInjector interface
func (ps *PositionService) InjectDependencies(deps ...interface{}) error {
	for _, dep := range deps {
		switch d := dep.(type) {
		case *gorm.DB:
			ps.db = d
		}
	}
	return nil
}

func (ps *PositionService) GetPositions(params PositionQueryParams) (*utils.PaginationResponse, error) {
	var positions []models.Position
	var total int64

	query := ps.db.Model(&models.Position{}).Preload("Department")

	// 过滤条件
	if params.DepartmentID > 0 {
		query = query.Where("department_id = ?", params.DepartmentID)
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
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&positions).Error; err != nil {
		return nil, err
	}

	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &utils.PaginationResponse{
		Data:       positions,
		TotalItems: total,
		TotalPages: totalPages,
		Page:       params.Page,
		PageSize:   params.PageSize,
	}, nil
}

func (ps *PositionService) GetPositionByID(id uint) (*models.Position, error) {
	var position models.Position
	err := ps.db.Preload("Department").First(&position, id).Error
	if err != nil {
		return nil, err
	}
	return &position, nil
}

func (ps *PositionService) CreatePosition(position *models.Position) (*models.Position, error) {
	// 检查职位编码是否已存在
	var count int64
	ps.db.Model(&models.Position{}).Where("code = ?", position.Code).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("职位编码已存在")
	}

	// 设置默认状态
	if position.Status == "" {
		position.Status = "active"
	}

	err := ps.db.Create(position).Error
	if err != nil {
		return nil, err
	}

	return ps.GetPositionByID(position.ID)
}

func (ps *PositionService) UpdatePosition(position *models.Position) (*models.Position, error) {
	// 检查职位编码是否已被其他职位使用
	var count int64
	ps.db.Model(&models.Position{}).Where("code = ? AND id != ?", position.Code, position.ID).Count(&count)
	if count > 0 {
		return nil, utils.NewValidationError("职位编码已存在")
	}

	err := ps.db.Model(position).Updates(position).Error
	if err != nil {
		return nil, err
	}

	return ps.GetPositionByID(position.ID)
}

func (ps *PositionService) DeletePosition(id uint) error {
	// 检查是否有员工使用此职位
	var count int64
	ps.db.Model(&models.Employee{}).Where("position_id = ?", id).Count(&count)
	if count > 0 {
		return utils.NewValidationError("无法删除已被员工使用的职位")
	}

	return ps.db.Delete(&models.Position{}, id).Error
}

func (ps *PositionService) GetPositionsByDepartment(departmentID uint) ([]*models.Position, error) {
	var positions []*models.Position
	err := ps.db.Where("department_id = ?", departmentID).
		Order("name ASC").Find(&positions).Error
	return positions, err
}

func (ps *PositionService) GetPositionStatistics() (*PositionStatistics, error) {
	var total int64
	ps.db.Model(&models.Position{}).Count(&total)

	// 按部门统计
	var deptStats []struct {
		DepartmentName string `json:"department_name"`
		Count          int64  `json:"count"`
	}
	ps.db.Model(&models.Position{}).
		Select("departments.name as department_name, COUNT(*) as count").
		Joins("LEFT JOIN departments ON positions.department_id = departments.id").
		Group("departments.name").
		Scan(&deptStats)

	byDepartment := make(map[string]int64)
	for _, stat := range deptStats {
		byDepartment[stat.DepartmentName] = stat.Count
	}

	// 按状态统计
	var statusStats []struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	ps.db.Model(&models.Position{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusStats)

	byStatus := make(map[string]int64)
	var activePositions, inactivePositions int64
	for _, stat := range statusStats {
		byStatus[stat.Status] = stat.Count
		if stat.Status == "active" {
			activePositions = stat.Count
		} else {
			inactivePositions += stat.Count
		}
	}

	return &PositionStatistics{
		Total:             total,
		ByDepartment:      byDepartment,
		ByStatus:          byStatus,
		ActivePositions:   activePositions,
		InactivePositions: inactivePositions,
	}, nil
}

func (ps *PositionService) SearchPositions(query string) ([]*models.Position, error) {
	var positions []*models.Position
	err := ps.db.Preload("Department").
		Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			"%"+query+"%", "%"+query+"%", "%"+query+"%").
		Where("status = ?", "active").
		Limit(20).
		Find(&positions).Error
	return positions, err
}

func (ps *PositionService) BulkCreatePositions(positions []*models.Position) (*BulkResult, error) {
	result := &BulkResult{
		Success: 0,
		Failed:  0,
		Errors:  []string{},
	}

	for _, position := range positions {
		if position.Status == "" {
			position.Status = "active"
		}

		_, err := ps.CreatePosition(position)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, err.Error())
		} else {
			result.Success++
		}
	}

	return result, nil
}

// GetAllPositions 获取所有职位（用于父级选择）
func (ps *PositionService) GetAllPositions() ([]*models.Position, error) {
	var positions []*models.Position
	err := ps.db.Preload("Department").
		Preload("Parent").
		Where("status = ?", "active").
		Order("level ASC, sort ASC, name ASC").
		Find(&positions).Error
	return positions, err
}

// GetPositionTree 获取职位树形结构
func (ps *PositionService) GetPositionTree() ([]*PositionTreeNode, error) {
	var positions []*models.Position
	
	// 先查询没有status字段条件的，如果失败再加上条件
	query := ps.db.Preload("Department").Preload("Parent")
	
	// 检查是否有status字段
	var count int64
	if err := ps.db.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_name = 'positions' AND column_name = 'status'").Count(&count).Error; err == nil && count > 0 {
		// 如果有status字段，则添加过滤条件
		query = query.Where("status = ?", "active")
	}
	
	err := query.Order("level ASC, sort ASC, name ASC").Find(&positions).Error
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	positionMap := make(map[uint]*PositionTreeNode)
	var rootNodes []*PositionTreeNode

	// 第一遍：创建所有节点
	for i := range positions {
		node := &PositionTreeNode{
			Position: positions[i], 
			Children: []*PositionTreeNode{},
		}
		positionMap[positions[i].ID] = node
	}

	// 第二遍：建立父子关系
	for i := range positions {
		node := positionMap[positions[i].ID]
		if positions[i].ParentID != nil && *positions[i].ParentID > 0 {
			// 有父级，添加到父级的children中
			if parentNode, exists := positionMap[*positions[i].ParentID]; exists {
				parentNode.Children = append(parentNode.Children, node)
			}
		} else {
			// 没有父级，是根节点
			rootNodes = append(rootNodes, node)
		}
	}

	return rootNodes, nil
}
