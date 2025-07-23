package services

import (
	"fmt"
	"gin-project/models"

	"gorm.io/gorm"
)

type DepartmentServiceInterface interface {
	GetDepartments(parentID *uint, includeChildren bool) ([]*models.Department, error)
	GetDepartmentByID(id uint) (*models.Department, error)
	CreateDepartment(department *models.Department) (*models.Department, error)
	UpdateDepartment(department *models.Department) (*models.Department, error)
	DeleteDepartment(id uint) error
	GetDepartmentTree() ([]*models.Department, error)
	MoveDepartment(id uint, parentID *uint) error
	GetDepartmentStatistics() (*DepartmentStatistics, error)
	GetDepartmentHierarchy(id uint) (*DepartmentHierarchy, error)
	GetAllSubDepartments(id uint) ([]*models.Department, error)
	GetDepartmentPath(id uint) ([]*models.Department, error)
	GetDepartmentLevel(id uint) (int, error)
	BulkUpdateDepartmentSort(updates []DepartmentSortUpdate) error
	GetDepartmentChart() (*DepartmentChart, error)
}

type DepartmentService struct {
	db *gorm.DB
}

type DepartmentStatistics struct {
	Total          int64            `json:"total"`
	ByLevel        map[string]int64 `json:"by_level"`
	WithManager    int64            `json:"with_manager"`
	EmployeeCounts map[uint]int64   `json:"employee_counts"`
}

type DepartmentHierarchy struct {
	Department    *models.Department   `json:"department"`
	Parent        *models.Department   `json:"parent,omitempty"`
	Children      []*models.Department `json:"children"`
	Ancestors     []*models.Department `json:"ancestors"`
	Descendants   []*models.Department `json:"descendants"`
	Level         int                  `json:"level"`
	EmployeeCount int64                `json:"employee_count"`
}

type DepartmentSortUpdate struct {
	ID   uint `json:"id"`
	Sort int  `json:"sort"`
}

type DepartmentChart struct {
	Nodes []DepartmentNode `json:"nodes"`
	Edges []DepartmentEdge `json:"edges"`
}

type DepartmentNode struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Code          string `json:"code"`
	ManagerName   string `json:"manager_name,omitempty"`
	EmployeeCount int64  `json:"employee_count"`
	Level         int    `json:"level"`
}

type DepartmentEdge struct {
	From uint `json:"from"`
	To   uint `json:"to"`
}

func NewDepartmentService(db *gorm.DB) DepartmentServiceInterface {
	return &DepartmentService{
		db: db,
	}
}

// InjectDependencies implements DependencyInjector interface
func (ds *DepartmentService) InjectDependencies(deps ...interface{}) error {
	for _, dep := range deps {
		switch d := dep.(type) {
		case *gorm.DB:
			ds.db = d
		}
	}
	return nil
}

func (ds *DepartmentService) GetDepartments(parentID *uint, includeChildren bool) ([]*models.Department, error) {
	var departments []*models.Department
	query := ds.db.Model(&models.Department{}).
		Preload("Manager").Where("is_active = ?", true)

	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	if includeChildren {
		query = query.Preload("Children")
	}

	err := query.Order("sort ASC, created_at ASC").Find(&departments).Error
	return departments, err
}

func (ds *DepartmentService) GetDepartmentByID(id uint) (*models.Department, error) {
	var department models.Department
	err := ds.db.Preload("Parent").
		Preload("Children").
		First(&department, id).Error
	if err != nil {
		return nil, err
	}
	
	// 手动加载正确的部门负责人
	if department.ManagerID != nil {
		var manager models.Employee
		if err := ds.db.Where("id = ? AND deleted_at IS NULL", *department.ManagerID).First(&manager).Error; err == nil {
			department.Manager = &manager
		}
	}
	
	return &department, nil
}

func (ds *DepartmentService) CreateDepartment(department *models.Department) (*models.Department, error) {
	if !department.IsActive {
		department.IsActive = true
	}

	// 生成部门编码
	if department.Code == "" {
		department.Code = ds.generateDepartmentCode()
	}

	err := ds.db.Create(department).Error
	if err != nil {
		return nil, err
	}

	return ds.GetDepartmentByID(department.ID)
}

func (ds *DepartmentService) UpdateDepartment(department *models.Department) (*models.Department, error) {
	err := ds.db.Model(department).Updates(department).Error
	if err != nil {
		return nil, err
	}

	return ds.GetDepartmentByID(department.ID)
}

func (ds *DepartmentService) DeleteDepartment(id uint) error {
	// 检查是否有子部门
	var childCount int64
	ds.db.Model(&models.Department{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		return gorm.ErrInvalidData
	}

	// 检查是否有员工
	var employeeCount int64
	ds.db.Model(&models.Employee{}).Where("department_id = ?", id).Count(&employeeCount)
	if employeeCount > 0 {
		return gorm.ErrInvalidData
	}

	return ds.db.Delete(&models.Department{}, id).Error
}

func (ds *DepartmentService) GetDepartmentTree() ([]*models.Department, error) {
	var departments []*models.Department
	err := ds.db.Where("is_active = ?", true).
		Order("sort ASC, created_at ASC").
		Find(&departments).Error
	if err != nil {
		return nil, err
	}

	// 获取各部门的员工数量
	var employeeCounts []struct {
		DepartmentID uint  `json:"department_id"`
		Count        int64 `json:"count"`
	}
	ds.db.Model(&models.Employee{}).
		Select("department_id, COUNT(*) as count").
		Where("department_id IS NOT NULL").
		Group("department_id").
		Scan(&employeeCounts)

	// 构建员工数量映射
	empCountMap := make(map[uint]int64)
	for _, count := range employeeCounts {
		empCountMap[count.DepartmentID] = count.Count
	}

	// 构建树形结构
	departmentMap := make(map[uint]*models.Department)
	var roots []*models.Department

	// 第一遍遍历，建立映射并设置员工数量
	for _, dept := range departments {
		departmentMap[dept.ID] = dept
		dept.Children = []*models.Department{}
		// 设置员工数量到部门对象中
		dept.EmployeeCount = int(empCountMap[dept.ID])
		
		// 手动加载正确的部门负责人
		if dept.ManagerID != nil {
			var manager models.Employee
			if err := ds.db.Where("id = ? AND deleted_at IS NULL", *dept.ManagerID).First(&manager).Error; err == nil {
				dept.Manager = &manager
			}
		}
	}

	// 第二遍遍历，建立父子关系
	for _, dept := range departments {
		if dept.ParentID != nil {
			if parent, exists := departmentMap[*dept.ParentID]; exists {
				parent.Children = append(parent.Children, dept)
			}
		} else {
			roots = append(roots, dept)
		}
	}

	// 第三遍遍历，计算子单元数量
	for _, dept := range departments {
		dept.SubunitCount = len(dept.Children)
	}

	return roots, nil
}

func (ds *DepartmentService) MoveDepartment(id uint, parentID *uint) error {
	// 检查是否会形成循环引用
	if parentID != nil && ds.wouldCreateCycle(id, *parentID) {
		return gorm.ErrInvalidData
	}

	return ds.db.Model(&models.Department{}).
		Where("id = ?", id).
		Update("parent_id", parentID).Error
}

func (ds *DepartmentService) GetDepartmentStatistics() (*DepartmentStatistics, error) {
	var total int64
	ds.db.Model(&models.Department{}).Where("is_active = ?", true).Count(&total)

	// 有管理者的部门数量
	var withManager int64
	ds.db.Model(&models.Department{}).
		Where("is_active = ? AND manager_id IS NOT NULL", true).
		Count(&withManager)

	// 各部门员工数量
	var employeeCounts map[uint]int64 = make(map[uint]int64)
	var deptEmployeeCounts []struct {
		DepartmentID uint  `json:"department_id"`
		Count        int64 `json:"count"`
	}
	ds.db.Model(&models.Employee{}).
		Select("department_id, COUNT(*) as count").
		Where("department_id IS NOT NULL").
		Group("department_id").
		Scan(&deptEmployeeCounts)

	for _, item := range deptEmployeeCounts {
		employeeCounts[item.DepartmentID] = item.Count
	}

	return &DepartmentStatistics{
		Total:          total,
		WithManager:    withManager,
		EmployeeCounts: employeeCounts,
	}, nil
}

// 检查是否会形成循环引用
func (ds *DepartmentService) wouldCreateCycle(deptID, newParentID uint) bool {
	if deptID == newParentID {
		return true
	}

	// 递归检查新父部门的所有祖先
	var parent models.Department
	if err := ds.db.First(&parent, newParentID).Error; err != nil {
		return false
	}

	if parent.ParentID != nil {
		return ds.wouldCreateCycle(deptID, *parent.ParentID)
	}

	return false
}

func (ds *DepartmentService) generateDepartmentCode() string {
	var count int64
	ds.db.Model(&models.Department{}).Count(&count)
	return fmt.Sprintf("DEPT%04d", count+1)
}

// GetDepartmentHierarchy 获取部门层级信息
func (ds *DepartmentService) GetDepartmentHierarchy(id uint) (*DepartmentHierarchy, error) {
	// 获取部门基本信息
	department, err := ds.GetDepartmentByID(id)
	if err != nil {
		return nil, err
	}

	// 获取层级
	level, err := ds.GetDepartmentLevel(id)
	if err != nil {
		return nil, err
	}

	// 获取祖先部门
	ancestors, err := ds.GetDepartmentPath(id)
	if err != nil {
		return nil, err
	}

	// 获取所有子部门
	descendants, err := ds.GetAllSubDepartments(id)
	if err != nil {
		return nil, err
	}

	// 获取员工数量
	var employeeCount int64
	ds.db.Model(&models.Employee{}).Where("department_id = ?", id).Count(&employeeCount)

	return &DepartmentHierarchy{
		Department:    department,
		Parent:        department.Parent,
		Children:      department.Children,
		Ancestors:     ancestors,
		Descendants:   descendants,
		Level:         level,
		EmployeeCount: employeeCount,
	}, nil
}

// GetAllSubDepartments 获取所有子部门（递归）
func (ds *DepartmentService) GetAllSubDepartments(id uint) ([]*models.Department, error) {
	var allSubDepts []*models.Department

	// 递归获取子部门
	err := ds.collectSubDepartments(id, &allSubDepts)
	if err != nil {
		return nil, err
	}

	return allSubDepts, nil
}

func (ds *DepartmentService) collectSubDepartments(parentID uint, result *[]*models.Department) error {
	var children []*models.Department
	err := ds.db.Where("parent_id = ? AND is_active = ?", parentID, true).
		Preload("Manager").Find(&children).Error
	if err != nil {
		return err
	}

	for _, child := range children {
		*result = append(*result, child)
		// 递归获取子部门的子部门
		err = ds.collectSubDepartments(child.ID, result)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetDepartmentPath 获取部门路径（从根到当前部门）
func (ds *DepartmentService) GetDepartmentPath(id uint) ([]*models.Department, error) {
	var path []*models.Department

	err := ds.collectAncestors(id, &path)
	if err != nil {
		return nil, err
	}

	// 反转路径，使其从根部门开始
	for i := len(path)/2 - 1; i >= 0; i-- {
		opp := len(path) - 1 - i
		path[i], path[opp] = path[opp], path[i]
	}

	return path, nil
}

func (ds *DepartmentService) collectAncestors(id uint, result *[]*models.Department) error {
	var dept models.Department
	err := ds.db.Preload("Manager").First(&dept, id).Error
	if err != nil {
		return err
	}

	*result = append(*result, &dept)

	if dept.ParentID != nil {
		return ds.collectAncestors(*dept.ParentID, result)
	}

	return nil
}

// GetDepartmentLevel 获取部门层级（根部门为1）
func (ds *DepartmentService) GetDepartmentLevel(id uint) (int, error) {
	level := 1
	currentID := id

	for {
		var dept models.Department
		err := ds.db.First(&dept, currentID).Error
		if err != nil {
			return 0, err
		}

		if dept.ParentID == nil {
			break
		}

		level++
		currentID = *dept.ParentID
	}

	return level, nil
}

// BulkUpdateDepartmentSort 批量更新部门排序
func (ds *DepartmentService) BulkUpdateDepartmentSort(updates []DepartmentSortUpdate) error {
	return ds.db.Transaction(func(tx *gorm.DB) error {
		for _, update := range updates {
			err := tx.Model(&models.Department{}).
				Where("id = ?", update.ID).
				Update("sort", update.Sort).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// GetDepartmentChart 获取组织架构图数据
func (ds *DepartmentService) GetDepartmentChart() (*DepartmentChart, error) {
	var departments []*models.Department
	err := ds.db.Where("is_active = ?", true).
		Preload("Manager").Find(&departments).Error
	if err != nil {
		return nil, err
	}

	// 获取各部门员工数量
	var employeeCounts []struct {
		DepartmentID uint  `json:"department_id"`
		Count        int64 `json:"count"`
	}
	ds.db.Model(&models.Employee{}).
		Select("department_id, COUNT(*) as count").
		Where("department_id IS NOT NULL").
		Group("department_id").
		Scan(&employeeCounts)

	empCountMap := make(map[uint]int64)
	for _, count := range employeeCounts {
		empCountMap[count.DepartmentID] = count.Count
	}

	var nodes []DepartmentNode
	var edges []DepartmentEdge

	for _, dept := range departments {
		// 获取部门层级
		level, _ := ds.GetDepartmentLevel(dept.ID)

		// 创建节点
		node := DepartmentNode{
			ID:            dept.ID,
			Name:          dept.Name,
			Code:          dept.Code,
			EmployeeCount: empCountMap[dept.ID],
			Level:         level,
		}

		if dept.Manager != nil {
			node.ManagerName = dept.Manager.Name
		}

		nodes = append(nodes, node)

		// 创建边
		if dept.ParentID != nil {
			edge := DepartmentEdge{
				From: *dept.ParentID,
				To:   dept.ID,
			}
			edges = append(edges, edge)
		}
	}

	return &DepartmentChart{
		Nodes: nodes,
		Edges: edges,
	}, nil
}
