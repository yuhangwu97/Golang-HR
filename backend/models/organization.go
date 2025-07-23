package models

import (
	"time"
	"gorm.io/gorm"
)

// OrganizationUnitType 组织单元类型
type OrganizationUnitType string

const (
	CompanyUnit       OrganizationUnitType = "company"        // 公司/法人实体
	BusinessUnit      OrganizationUnitType = "business_unit"  // 业务单元
	DepartmentUnit    OrganizationUnitType = "department"     // 部门
	TeamUnit          OrganizationUnitType = "team"           // 团队
	CostCenterUnit    OrganizationUnitType = "cost_center"    // 成本中心
	LocationUnit      OrganizationUnitType = "location"       // 地理位置
	ProjectUnit       OrganizationUnitType = "project"        // 项目组
)


// OrganizationUnit 组织单元 - 支持多种类型的组织结构
type OrganizationUnit struct {
	ID                uint                   `json:"id" gorm:"primaryKey"`
	Name              string                 `json:"name" gorm:"size:200;not null;comment:组织单元名称"`
	Code              string                 `json:"code" gorm:"uniqueIndex;size:50;comment:组织单元编码"`
	Type              OrganizationUnitType   `json:"type" gorm:"size:20;not null;comment:组织单元类型"`
	ParentID          *uint                  `json:"parent_id" gorm:"comment:上级组织单元ID"`
	Parent            *OrganizationUnit      `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children          []*OrganizationUnit    `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Level             int                    `json:"level" gorm:"default:1;comment:组织层级"`
	Sort              int                    `json:"sort" gorm:"default:0;comment:排序"`
	
	// 基本信息
	Description       string                 `json:"description" gorm:"type:text;comment:描述"`
	ShortName         string                 `json:"short_name" gorm:"size:50;comment:简称"`
	ExternalID        string                 `json:"external_id" gorm:"size:100;comment:外部系统ID"`
	
	// 地理和成本信息
	CountryCode       string                 `json:"country_code" gorm:"size:10;comment:国家代码"`
	CurrencyCode      string                 `json:"currency_code" gorm:"size:10;comment:货币代码"`
	TimeZone          string                 `json:"time_zone" gorm:"size:50;comment:时区"`
	CostCenter        string                 `json:"cost_center" gorm:"size:50;comment:成本中心编码"`
	
	// 联系信息
	Address           string                 `json:"address" gorm:"size:500;comment:地址"`
	Phone             string                 `json:"phone" gorm:"size:50;comment:电话"`
	Email             string                 `json:"email" gorm:"size:100;comment:邮箱"`
	Website           string                 `json:"website" gorm:"size:200;comment:网站"`
	
	// 管理信息
	ManagerID         *uint                  `json:"manager_id" gorm:"comment:负责人ID"`
	Manager           *Employee              `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
	FunctionalManagerID *uint                `json:"functional_manager_id" gorm:"comment:功能负责人ID"`
	FunctionalManager *Employee              `json:"functional_manager,omitempty" gorm:"foreignKey:FunctionalManagerID"`
	
	// 有效期
	EffectiveDate     *time.Time             `json:"effective_date" gorm:"comment:生效日期"`
	ExpirationDate    *time.Time             `json:"expiration_date" gorm:"comment:失效日期"`
	
	// 状态和配置
	Status            string                 `json:"status" gorm:"size:20;default:active;comment:状态"`
	IsActive          bool                   `json:"is_active" gorm:"default:true;comment:是否激活"`
	IsHeadquarters    bool                   `json:"is_headquarters" gorm:"default:false;comment:是否总部"`
	AllowSubunits     bool                   `json:"allow_subunits" gorm:"default:true;comment:允许下级单元"`
	
	// 员工分配
	EmployeeAssignments []*EmployeeAssignment `json:"employee_assignments,omitempty" gorm:"foreignKey:OrganizationUnitID"`
	
	// 计算字段（不存储在数据库中）
	EmployeeCount     int                    `json:"employee_count" gorm:"-"`
	SubunitCount      int                    `json:"subunit_count" gorm:"-"`
	
	// 系统字段
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	DeletedAt         gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

// EmployeeAssignment 员工组织分配 - 支持多种分配类型
type EmployeeAssignment struct {
	ID                  uint                   `json:"id" gorm:"primaryKey"`
	EmployeeID          uint                   `json:"employee_id" gorm:"not null;comment:员工ID"`
	Employee            *Employee              `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	OrganizationUnitID  uint                   `json:"organization_unit_id" gorm:"not null;comment:组织单元ID"`
	OrganizationUnit    *OrganizationUnit      `json:"organization_unit,omitempty" gorm:"foreignKey:OrganizationUnitID"`
	
	// 分配信息
	AssignmentType      AssignmentType         `json:"assignment_type" gorm:"size:20;default:primary;comment:分配类型"`
	ManagementType      ManagementType         `json:"management_type" gorm:"size:20;default:line;comment:管理类型"`
	PositionID          *uint                  `json:"position_id" gorm:"comment:职位ID"`
	Position            *Position              `json:"position,omitempty" gorm:"foreignKey:PositionID"`
	JobLevelID          *uint                  `json:"job_level_id" gorm:"comment:职级ID"`
	JobLevel            *JobLevel              `json:"job_level,omitempty" gorm:"foreignKey:JobLevelID"`
	
	// 管理关系
	DirectManagerID     *uint                  `json:"direct_manager_id" gorm:"comment:直接上级ID"`
	DirectManager       *Employee              `json:"direct_manager,omitempty" gorm:"foreignKey:DirectManagerID"`
	FunctionalManagerID *uint                  `json:"functional_manager_id" gorm:"comment:功能上级ID"`
	FunctionalManager   *Employee              `json:"functional_manager,omitempty" gorm:"foreignKey:FunctionalManagerID"`
	
	// 工作安排
	WorkPercentage      float64                `json:"work_percentage" gorm:"type:decimal(5,2);default:100.00;comment:工作占比"`
	IsPrimary           bool                   `json:"is_primary" gorm:"default:false;comment:是否主要分配"`
	
	// 有效期
	EffectiveDate       *time.Time             `json:"effective_date" gorm:"comment:生效日期"`
	ExpirationDate      *time.Time             `json:"expiration_date" gorm:"comment:失效日期"`
	
	// 状态
	Status              string                 `json:"status" gorm:"size:20;default:active;comment:状态"`
	Reason              string                 `json:"reason" gorm:"size:200;comment:分配原因"`
	
	// 系统字段
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	DeletedAt           gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

// ChangeType 变更类型
type ChangeType string

const (
	ChangeTypeCreate     ChangeType = "create"      // 创建
	ChangeTypeUpdate     ChangeType = "update"      // 更新
	ChangeTypeDelete     ChangeType = "delete"      // 删除
	ChangeTypeMove       ChangeType = "move"        // 移动
	ChangeTypeActivate   ChangeType = "activate"    // 激活
	ChangeTypeDeactivate ChangeType = "deactivate"  // 停用
	ChangeTypeAssign     ChangeType = "assign"      // 分配员工
	ChangeTypeUnassign   ChangeType = "unassign"    // 取消分配
	ChangeTypeReassign   ChangeType = "reassign"    // 重新分配
	ChangeTypeMerge      ChangeType = "merge"       // 合并
	ChangeTypeSplit      ChangeType = "split"       // 拆分
)

// ChangeStatus 变更状态
type ChangeStatus string

const (
	ChangeStatusPending   ChangeStatus = "pending"    // 待审批
	ChangeStatusApproved  ChangeStatus = "approved"   // 已审批
	ChangeStatusRejected  ChangeStatus = "rejected"   // 已拒绝
	ChangeStatusImplemented ChangeStatus = "implemented" // 已执行
	ChangeStatusCancelled ChangeStatus = "cancelled"  // 已取消
)

// OrganizationChange 组织变更记录
type OrganizationChange struct {
	ID                  uint                   `json:"id" gorm:"primaryKey"`
	ChangeType          ChangeType             `json:"change_type" gorm:"size:20;comment:变更类型"`
	EntityType          string                 `json:"entity_type" gorm:"size:30;comment:实体类型"`
	EntityID            uint                   `json:"entity_id" gorm:"comment:实体ID"`
	EntityName          string                 `json:"entity_name" gorm:"size:200;comment:实体名称"`
	
	// 变更内容
	FieldName           string                 `json:"field_name" gorm:"size:100;comment:字段名"`
	OldValue            string                 `json:"old_value" gorm:"type:text;comment:旧值"`
	NewValue            string                 `json:"new_value" gorm:"type:text;comment:新值"`
	ChangeReason        string                 `json:"change_reason" gorm:"type:text;comment:变更原因"`
	ChangeDescription   string                 `json:"change_description" gorm:"type:text;comment:变更描述"`
	
	// 生效信息
	EffectiveDate       *time.Time             `json:"effective_date" gorm:"comment:生效日期"`
	ProposedDate        *time.Time             `json:"proposed_date" gorm:"comment:提议日期"`
	ImplementedDate     *time.Time             `json:"implemented_date" gorm:"comment:执行日期"`
	
	// 审批信息
	Status              ChangeStatus           `json:"status" gorm:"size:20;default:pending;comment:状态"`
	InitiatorID         uint                   `json:"initiator_id" gorm:"comment:发起人ID"`
	Initiator           *Employee              `json:"initiator,omitempty" gorm:"foreignKey:InitiatorID"`
	ApproverID          *uint                  `json:"approver_id" gorm:"comment:审批人ID"`
	Approver            *Employee              `json:"approver,omitempty" gorm:"foreignKey:ApproverID"`
	ApprovalDate        *time.Time             `json:"approval_date" gorm:"comment:审批日期"`
	ApprovalNote        string                 `json:"approval_note" gorm:"type:text;comment:审批意见"`
	
	// 影响分析
	ImpactedEmployees   int                    `json:"impacted_employees" gorm:"comment:影响员工数"`
	ImpactDescription   string                 `json:"impact_description" gorm:"type:text;comment:影响描述"`
	ImpactedUnits       string                 `json:"impacted_units" gorm:"type:text;comment:影响的组织单元"`
	
	// 业务信息
	BusinessJustification string               `json:"business_justification" gorm:"type:text;comment:业务理由"`
	RiskAssessment      string                 `json:"risk_assessment" gorm:"type:text;comment:风险评估"`
	CommunicationPlan   string                 `json:"communication_plan" gorm:"type:text;comment:沟通计划"`
	
	// 关联变更
	ParentChangeID      *uint                  `json:"parent_change_id" gorm:"comment:父变更ID"`
	ParentChange        *OrganizationChange    `json:"parent_change,omitempty" gorm:"foreignKey:ParentChangeID"`
	RelatedChanges      []*OrganizationChange  `json:"related_changes,omitempty" gorm:"foreignKey:ParentChangeID"`
	
	// 附件和文档
	AttachmentURLs      string                 `json:"attachment_urls" gorm:"type:text;comment:附件链接"`
	DocumentationURLs   string                 `json:"documentation_urls" gorm:"type:text;comment:文档链接"`
	
	// 系统字段
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
}

// OrganizationHistory 组织架构历史记录
type OrganizationHistory struct {
	ID                  uint                   `json:"id" gorm:"primaryKey"`
	SnapshotDate        time.Time              `json:"snapshot_date" gorm:"comment:快照日期"`
	SnapshotReason      string                 `json:"snapshot_reason" gorm:"size:200;comment:快照原因"`
	
	// 组织单元快照
	UnitID              uint                   `json:"unit_id" gorm:"comment:组织单元ID"`
	UnitName            string                 `json:"unit_name" gorm:"size:200;comment:组织单元名称"`
	UnitCode            string                 `json:"unit_code" gorm:"size:50;comment:组织单元编码"`
	UnitType            OrganizationUnitType   `json:"unit_type" gorm:"size:20;comment:组织单元类型"`
	ParentID            *uint                  `json:"parent_id" gorm:"comment:上级组织单元ID"`
	ParentName          string                 `json:"parent_name" gorm:"size:200;comment:上级组织单元名称"`
	Level               int                    `json:"level" gorm:"comment:组织层级"`
	
	// 管理信息快照
	ManagerID           *uint                  `json:"manager_id" gorm:"comment:负责人ID"`
	ManagerName         string                 `json:"manager_name" gorm:"size:100;comment:负责人姓名"`
	FunctionalManagerID *uint                  `json:"functional_manager_id" gorm:"comment:功能负责人ID"`
	FunctionalManagerName string               `json:"functional_manager_name" gorm:"size:100;comment:功能负责人姓名"`
	
	// 统计信息快照
	EmployeeCount       int                    `json:"employee_count" gorm:"comment:员工数量"`
	DirectReports       int                    `json:"direct_reports" gorm:"comment:直接下属数"`
	SubunitCount        int                    `json:"subunit_count" gorm:"comment:下级单元数"`
	
	// 状态快照
	Status              string                 `json:"status" gorm:"size:20;comment:状态"`
	IsActive            bool                   `json:"is_active" gorm:"comment:是否激活"`
	
	// 变更信息
	ChangeID            *uint                  `json:"change_id" gorm:"comment:关联变更ID"`
	ChangeType          ChangeType             `json:"change_type" gorm:"size:20;comment:变更类型"`
	ChangedBy           uint                   `json:"changed_by" gorm:"comment:变更人ID"`
	ChangedByName       string                 `json:"changed_by_name" gorm:"size:100;comment:变更人姓名"`
	
	// 组织路径
	HierarchyPath       string                 `json:"hierarchy_path" gorm:"type:text;comment:层级路径"`
	PathIDs             string                 `json:"path_ids" gorm:"size:500;comment:路径ID序列"`
	
	// 系统字段
	CreatedAt           time.Time              `json:"created_at"`
}

// OrganizationHierarchy 组织层级视图
type OrganizationHierarchy struct {
	ID               uint                      `json:"id"`
	Name             string                    `json:"name"`
	Code             string                    `json:"code"`
	Type             OrganizationUnitType      `json:"type"`
	Level            int                       `json:"level"`
	ParentID         *uint                     `json:"parent_id"`
	ParentName       string                    `json:"parent_name,omitempty"`
	ManagerName      string                    `json:"manager_name,omitempty"`
	EmployeeCount    int                       `json:"employee_count"`
	DirectReports    int                       `json:"direct_reports"`
	TotalSubunits    int                       `json:"total_subunits"`
	IsActive         bool                      `json:"is_active"`
	Path             string                    `json:"path"` // 层级路径，如: "Company/IT/Development"
}

// EmployeeOrganizationView 员工组织视图
type EmployeeOrganizationView struct {
	EmployeeID         uint                   `json:"employee_id"`
	EmployeeName       string                 `json:"employee_name"`
	EmployeeCode       string                 `json:"employee_code"`
	
	// 主要分配
	PrimaryUnitID      uint                   `json:"primary_unit_id"`
	PrimaryUnitName    string                 `json:"primary_unit_name"`
	PrimaryUnitType    OrganizationUnitType   `json:"primary_unit_type"`
	PrimaryPosition    string                 `json:"primary_position,omitempty"`
	
	// 直接上级
	DirectManagerID    *uint                  `json:"direct_manager_id,omitempty"`
	DirectManagerName  string                 `json:"direct_manager_name,omitempty"`
	
	// 功能上级
	FunctionalManagerID *uint                 `json:"functional_manager_id,omitempty"`
	FunctionalManagerName string              `json:"functional_manager_name,omitempty"`
	
	// 额外分配
	AdditionalAssignments []EmployeeAssignment `json:"additional_assignments,omitempty"`
	
	// 组织路径
	OrganizationPath   string                 `json:"organization_path"`
	
	// 状态
	Status             string                 `json:"status"`
	EffectiveDate      *time.Time             `json:"effective_date,omitempty"`
}