package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
	"gorm.io/gorm"
)

// CustomDate handles date unmarshaling from frontend date strings
type CustomDate struct {
	time.Time
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	dateStr := strings.Trim(string(data), `"`)
	if dateStr == "" || dateStr == "null" {
		return nil
	}
	
	// Try different date formats
	layouts := []string{
		"2006-01-02T15:04:05Z07:00", // RFC3339
		"2006-01-02T15:04:05",       // ISO format without timezone
		"2006-01-02",                // Date only
		"2006-01-02 15:04:05",       // MySQL datetime format
	}
	
	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			cd.Time = t
			return nil
		}
	}
	
	return fmt.Errorf("cannot parse date: %s", dateStr)
}

// MarshalJSON implements the json.Marshaler interface
func (cd CustomDate) MarshalJSON() ([]byte, error) {
	if cd.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, cd.Time.Format("2006-01-02"))), nil
}

// Value implements the driver.Valuer interface for database storage
func (cd CustomDate) Value() (driver.Value, error) {
	if cd.Time.IsZero() {
		return nil, nil
	}
	return cd.Time, nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (cd *CustomDate) Scan(value interface{}) error {
	if value == nil {
		cd.Time = time.Time{}
		return nil
	}
	if t, ok := value.(time.Time); ok {
		cd.Time = t
		return nil
	}
	return fmt.Errorf("cannot scan %T into CustomDate", value)
}

// AssignmentType 员工分配类型
type AssignmentType string

const (
	PrimaryAssignment    AssignmentType = "primary"     // 主要分配
	AdditionalAssignment AssignmentType = "additional"  // 额外分配
	TemporaryAssignment  AssignmentType = "temporary"   // 临时分配
	ProjectAssignment    AssignmentType = "project"     // 项目分配
)

// ManagementType 管理类型
type ManagementType string

const (
	LineManagement       ManagementType = "line"        // 直线管理
	MatrixManagement     ManagementType = "matrix"      // 矩阵管理
	FunctionalManagement ManagementType = "functional"  // 功能管理
)

// Employee represents an employee in the HR system
type Employee struct {
	ID                 uint                   `json:"id" gorm:"primaryKey"`
	EmployeeID         string                 `json:"employee_id" gorm:"uniqueIndex;size:20;comment:员工工号"`
	Name               string                 `json:"name" gorm:"size:100;not null;comment:姓名"`
	Email              string                 `json:"email" gorm:"uniqueIndex;size:100;not null;comment:邮箱"`
	Phone              string                 `json:"phone" gorm:"size:20;comment:手机号"`
	Avatar             string                 `json:"avatar" gorm:"size:255;comment:头像"`
	Gender             string                 `json:"gender" gorm:"size:10;comment:性别"`
	Birthday           *CustomDate            `json:"birthday" gorm:"comment:生日"`
	IDCard             string                 `json:"id_card" gorm:"size:18;comment:身份证号"`
	Status             string                 `json:"status" gorm:"size:20;default:active;comment:状态"`
	
	// 主要职位信息（保持向后兼容）
	DepartmentID       uint                   `json:"department_id" gorm:"comment:主要部门ID"`
	Department         *Department            `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	PositionID         uint                   `json:"position_id" gorm:"comment:主要职位ID"`
	Position           *Position              `json:"position,omitempty" gorm:"foreignKey:PositionID"`
	JobLevelID         uint                   `json:"job_level_id" gorm:"comment:职级ID"`
	JobLevel           *JobLevel              `json:"job_level,omitempty" gorm:"foreignKey:JobLevelID"`
	ManagerID          *uint                  `json:"manager_id" gorm:"comment:直接上级ID"`
	Manager            *Employee              `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
	
	// 扩展多重分配支持
	FunctionalManagerID *uint                 `json:"functional_manager_id" gorm:"comment:功能上级ID"`
	FunctionalManager   *Employee             `json:"functional_manager,omitempty" gorm:"foreignKey:FunctionalManagerID"`
	SecondaryDepartmentID *uint               `json:"secondary_department_id" gorm:"comment:次要部门ID"`
	SecondaryDepartment   *Department         `json:"secondary_department,omitempty" gorm:"foreignKey:SecondaryDepartmentID"`
	WorkPercentage      float64               `json:"work_percentage" gorm:"type:decimal(5,2);default:100.00;comment:工作占比"`
	AssignmentType      AssignmentType        `json:"assignment_type" gorm:"size:20;default:primary;comment:分配类型"`
	ManagementType      ManagementType        `json:"management_type" gorm:"size:20;default:line;comment:管理类型"`
	
	// 入职信息
	HireDate           *CustomDate            `json:"hire_date" gorm:"comment:入职日期"`
	ProbationEndDate   *CustomDate            `json:"probation_end_date" gorm:"comment:试用期结束日期"`
	ContractStartDate  *CustomDate            `json:"contract_start_date" gorm:"comment:合同开始日期"`
	ContractEndDate    *CustomDate            `json:"contract_end_date" gorm:"comment:合同结束日期"`
	ContractType       string                 `json:"contract_type" gorm:"size:20;comment:合同类型"`
	
	// 薪资信息
	BaseSalary         float64                `json:"base_salary" gorm:"type:decimal(10,2);comment:基本薪资"`
	
	// 联系信息
	Address            string                 `json:"address" gorm:"size:255;comment:地址"`
	EmergencyContact   string                 `json:"emergency_contact" gorm:"size:100;comment:紧急联系人"`
	EmergencyPhone     string                 `json:"emergency_phone" gorm:"size:20;comment:紧急联系人电话"`
	
	// 教育背景
	Education          string                 `json:"education" gorm:"size:20;comment:学历"`
	School             string                 `json:"school" gorm:"size:100;comment:毕业学校"`
	Major              string                 `json:"major" gorm:"size:100;comment:专业"`
	
	// 工作经验
	WorkExperience     []WorkExperience       `json:"work_experience,omitempty" gorm:"foreignKey:EmployeeID"`
	
	// 多重分配关系（与现有organization逻辑集成）
	DepartmentAssignments []DepartmentAssignment `json:"department_assignments,omitempty" gorm:"foreignKey:EmployeeID"`
	
	// 系统字段
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	DeletedAt          gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

// DepartmentType 部门类型 - 扩展支持organization功能
type DepartmentType string

const (
	CompanyDept      DepartmentType = "company"        // 公司/法人实体
	BusinessUnitDept DepartmentType = "business_unit"  // 业务单元
	StandardDept     DepartmentType = "department"     // 标准部门
	TeamDept         DepartmentType = "team"           // 团队
	CostCenterDept   DepartmentType = "cost_center"    // 成本中心
	LocationDept     DepartmentType = "location"       // 地理位置
	ProjectDept      DepartmentType = "project"        // 项目组
)

// Department represents a department in the organization - 扩展支持organization功能
type Department struct {
	ID          uint                   `json:"id" gorm:"primaryKey"`
	Name        string                 `json:"name" gorm:"size:100;not null;comment:部门名称"`
	Code        string                 `json:"code" gorm:"uniqueIndex;size:20;comment:部门编码"`
	
	// 组织层级信息
	ParentID    *uint                  `json:"parent_id" gorm:"comment:上级部门ID"`
	Parent      *Department            `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children    []*Department          `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Level       int                    `json:"level" gorm:"default:1;comment:组织层级"`
	
	// 扩展组织类型和属性
	Type        DepartmentType         `json:"type" gorm:"size:20;default:department;comment:部门类型"`
	ShortName   string                 `json:"short_name" gorm:"size:50;comment:简称"`
	ExternalID  string                 `json:"external_id" gorm:"size:100;comment:外部系统ID"`
	
	// 地理和成本信息
	CountryCode  string                `json:"country_code" gorm:"size:10;comment:国家代码"`
	CurrencyCode string                `json:"currency_code" gorm:"size:10;comment:货币代码"`
	TimeZone     string                `json:"time_zone" gorm:"size:50;comment:时区"`
	CostCenter   string                `json:"cost_center" gorm:"size:50;comment:成本中心编码"`
	
	// 联系信息
	Address     string                 `json:"address" gorm:"size:500;comment:地址"`
	Phone       string                 `json:"phone" gorm:"size:50;comment:电话"`
	Email       string                 `json:"email" gorm:"size:100;comment:邮箱"`
	Website     string                 `json:"website" gorm:"size:200;comment:网站"`
	
	// 管理信息
	ManagerID   *uint                  `json:"manager_id" gorm:"comment:部门负责人ID"`
	Manager     *Employee              `json:"manager,omitempty" gorm:"foreignKey:ManagerID;references:ID"`
	FunctionalManagerID *uint          `json:"functional_manager_id" gorm:"comment:功能负责人ID"`
	FunctionalManager   *Employee      `json:"functional_manager,omitempty" gorm:"foreignKey:FunctionalManagerID"`
	
	// 有效期和状态
	EffectiveDate  *time.Time            `json:"effective_date" gorm:"comment:生效日期"`
	ExpirationDate *time.Time            `json:"expiration_date" gorm:"comment:失效日期"`
	IsActive       bool                  `json:"is_active" gorm:"default:true;comment:是否激活"`
	IsHeadquarters bool                  `json:"is_headquarters" gorm:"default:false;comment:是否总部"`
	AllowSubunits  bool                  `json:"allow_subunits" gorm:"default:true;comment:允许下级单元"`
	
	Description string                 `json:"description" gorm:"type:text;comment:部门描述"`
	Status      string                 `json:"status" gorm:"size:20;default:active;comment:状态"`
	Sort        int                    `json:"sort" gorm:"default:0;comment:排序"`
	
	// 计算字段（不存储在数据库中）
	EmployeeCount     int              `json:"employee_count" gorm:"-"`
	SubunitCount      int              `json:"subunit_count" gorm:"-"`
	
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

// Position represents a job position
type Position struct {
	ID           uint                   `json:"id" gorm:"primaryKey"`
	Name         string                 `json:"name" gorm:"size:100;not null;comment:职位名称"`
	Code         string                 `json:"code" gorm:"uniqueIndex;size:20;comment:职位编码"`
	DepartmentID uint                   `json:"department_id" gorm:"comment:所属部门ID"`
	Department   *Department            `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	ParentID     *uint                  `json:"parent_id" gorm:"comment:上级职位ID"`
	Parent       *Position              `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children     []*Position            `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Level        int                    `json:"level" gorm:"default:1;comment:职位层级"`
	Sort         int                    `json:"sort" gorm:"default:0;comment:排序"`
	Description  string                 `json:"description" gorm:"type:text;comment:职位描述"`
	Requirements string                 `json:"requirements" gorm:"type:text;comment:任职要求"`
	Status       string                 `json:"status" gorm:"size:20;default:active;comment:状态"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	DeletedAt    gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

// User represents a system user with authentication
type User struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Username    string         `json:"username" gorm:"uniqueIndex;size:50;not null;comment:用户名"`
	Email       string         `json:"email" gorm:"uniqueIndex;size:100;not null;comment:邮箱"`
	Password    string         `json:"-" gorm:"size:255;not null;comment:密码"`
	EmployeeID  *uint          `json:"employee_id" gorm:"comment:关联员工ID"`
	Employee    *Employee      `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	Status      string         `json:"status" gorm:"size:20;default:active;comment:状态"`
	LastLoginAt *time.Time     `json:"last_login_at" gorm:"comment:最后登录时间"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// Role represents a system role
type Role struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;size:50;not null;comment:角色名称"`
	Code        string         `json:"code" gorm:"uniqueIndex;size:50;not null;comment:角色编码"`
	Description string         `json:"description" gorm:"type:text;comment:角色描述"`
	Status      string         `json:"status" gorm:"size:20;default:active;comment:状态"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// Permission represents a system permission
type Permission struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null;comment:权限名称"`
	Code        string         `json:"code" gorm:"uniqueIndex;size:100;not null;comment:权限编码"`
	Resource    string         `json:"resource" gorm:"size:100;comment:资源"`
	Action      string         `json:"action" gorm:"size:50;comment:操作"`
	Description string         `json:"description" gorm:"type:text;comment:权限描述"`
	Status      string         `json:"status" gorm:"size:20;default:active;comment:状态"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"comment:用户ID"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	RoleID    uint      `json:"role_id" gorm:"comment:角色ID"`
	Role      *Role     `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	CreatedAt time.Time `json:"created_at"`
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	ID           uint        `json:"id" gorm:"primaryKey"`
	RoleID       uint        `json:"role_id" gorm:"comment:角色ID"`
	Role         *Role       `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	PermissionID uint        `json:"permission_id" gorm:"comment:权限ID"`
	Permission   *Permission `json:"permission,omitempty" gorm:"foreignKey:PermissionID"`
	CreatedAt    time.Time   `json:"created_at"`
}

// JobLevel represents job levels
type JobLevel struct {
	ID          uint                   `json:"id" gorm:"primaryKey"`
	Name        string                 `json:"name" gorm:"size:100;not null;comment:职级名称"`
	Code        string                 `json:"code" gorm:"uniqueIndex;size:20;comment:职级编码"`
	Level       int                    `json:"level" gorm:"comment:职级等级"`
	MinSalary   float64                `json:"min_salary" gorm:"type:decimal(10,2);comment:最低薪资"`
	MaxSalary   float64                `json:"max_salary" gorm:"type:decimal(10,2);comment:最高薪资"`
	Description string                 `json:"description" gorm:"type:text;comment:职级描述"`
	Status      string                 `json:"status" gorm:"size:20;default:active;comment:状态"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

// WorkExperience represents work experience
type WorkExperience struct {
	ID          uint                   `json:"id" gorm:"primaryKey"`
	EmployeeID  uint                   `json:"employee_id" gorm:"not null;comment:员工ID"`
	Company     string                 `json:"company" gorm:"size:100;comment:公司名称"`
	Position    string                 `json:"position" gorm:"size:100;comment:职位"`
	StartDate   *time.Time             `json:"start_date" gorm:"comment:开始日期"`
	EndDate     *time.Time             `json:"end_date" gorm:"comment:结束日期"`
	Description string                 `json:"description" gorm:"type:text;comment:工作描述"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// Attendance represents attendance records
type Attendance struct {
	ID            uint                   `json:"id" gorm:"primaryKey"`
	EmployeeID    uint                   `json:"employee_id" gorm:"not null;comment:员工ID"`
	Employee      *Employee              `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	Date          time.Time              `json:"date" gorm:"type:date;comment:考勤日期"`
	CheckInTime   *time.Time             `json:"check_in_time" gorm:"comment:签到时间"`
	CheckOutTime  *time.Time             `json:"check_out_time" gorm:"comment:签退时间"`
	WorkHours     float64                `json:"work_hours" gorm:"type:decimal(4,2);comment:工作小时数"`
	Status        string                 `json:"status" gorm:"size:20;comment:考勤状态"`
	Remark        string                 `json:"remark" gorm:"size:255;comment:备注"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// Leave represents leave requests
type Leave struct {
	ID          uint                   `json:"id" gorm:"primaryKey"`
	EmployeeID  uint                   `json:"employee_id" gorm:"not null;comment:员工ID"`
	Employee    *Employee              `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	Type        string                 `json:"type" gorm:"size:20;comment:请假类型"`
	StartDate   time.Time              `json:"start_date" gorm:"comment:开始日期"`
	EndDate     time.Time              `json:"end_date" gorm:"comment:结束日期"`
	Days        float64                `json:"days" gorm:"type:decimal(4,1);comment:请假天数"`
	Reason      string                 `json:"reason" gorm:"type:text;comment:请假原因"`
	Status      string                 `json:"status" gorm:"size:20;default:pending;comment:审批状态"`
	ApproverID  *uint                  `json:"approver_id" gorm:"comment:审批人ID"`
	Approver    *Employee              `json:"approver,omitempty" gorm:"foreignKey:ApproverID"`
	ApproveTime *time.Time             `json:"approve_time" gorm:"comment:审批时间"`
	ApproveNote string                 `json:"approve_note" gorm:"type:text;comment:审批意见"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// Salary represents salary records
type Salary struct {
	ID             uint                   `json:"id" gorm:"primaryKey"`
	EmployeeID     uint                   `json:"employee_id" gorm:"not null;comment:员工ID"`
	Employee       *Employee              `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	Month          string                 `json:"month" gorm:"size:7;comment:薪资月份"`
	BaseSalary     float64                `json:"base_salary" gorm:"type:decimal(10,2);comment:基本薪资"`
	Bonus          float64                `json:"bonus" gorm:"type:decimal(10,2);comment:奖金"`
	Allowance      float64                `json:"allowance" gorm:"type:decimal(10,2);comment:津贴"`
	Deduction      float64                `json:"deduction" gorm:"type:decimal(10,2);comment:扣款"`
	GrossSalary    float64                `json:"gross_salary" gorm:"type:decimal(10,2);comment:应发薪资"`
	Tax            float64                `json:"tax" gorm:"type:decimal(10,2);comment:个人所得税"`
	SocialSecurity float64                `json:"social_security" gorm:"type:decimal(10,2);comment:社保"`
	HousingFund    float64                `json:"housing_fund" gorm:"type:decimal(10,2);comment:公积金"`
	NetSalary      float64                `json:"net_salary" gorm:"type:decimal(10,2);comment:实发薪资"`
	Status         string                 `json:"status" gorm:"size:20;default:draft;comment:状态"`
	Remark         string                 `json:"remark" gorm:"type:text;comment:备注"`
	PayrollRecords []PayrollRecord        `json:"payroll_records,omitempty" gorm:"foreignKey:SalaryID"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// PayrollRecord represents salary payment records
type PayrollRecord struct {
	ID            uint                   `json:"id" gorm:"primaryKey"`
	SalaryID      uint                   `json:"salary_id" gorm:"not null;comment:薪资记录ID"`
	Salary        *Salary                `json:"salary,omitempty" gorm:"foreignKey:SalaryID"`
	PaymentDate   *time.Time             `json:"payment_date" gorm:"comment:发放日期"`
	PaymentMethod string                 `json:"payment_method" gorm:"size:20;comment:发放方式"`
	BankAccount   string                 `json:"bank_account" gorm:"size:50;comment:银行账户"`
	PaymentAmount float64                `json:"payment_amount" gorm:"type:decimal(10,2);comment:发放金额"`
	Status        string                 `json:"status" gorm:"size:20;default:pending;comment:发放状态"`
	ProcessorID   *uint                  `json:"processor_id" gorm:"comment:处理人ID"`
	Processor     *Employee              `json:"processor,omitempty" gorm:"foreignKey:ProcessorID"`
	Remark        string                 `json:"remark" gorm:"type:text;comment:备注"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// Recruitment represents recruitment positions
type Recruitment struct {
	ID               uint                   `json:"id" gorm:"primaryKey"`
	Title            string                 `json:"title" gorm:"size:200;not null;comment:招聘职位"`
	DepartmentID     uint                   `json:"department_id" gorm:"comment:部门ID"`
	Department       *Department            `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	PositionID       uint                   `json:"position_id" gorm:"comment:职位ID"`
	Position         *Position              `json:"position,omitempty" gorm:"foreignKey:PositionID"`
	Count            int                    `json:"count" gorm:"comment:招聘人数"`
	MinSalary        float64                `json:"min_salary" gorm:"type:decimal(10,2);comment:最低薪资"`
	MaxSalary        float64                `json:"max_salary" gorm:"type:decimal(10,2);comment:最高薪资"`
	Requirements     string                 `json:"requirements" gorm:"type:text;comment:任职要求"`
	Description      string                 `json:"description" gorm:"type:text;comment:职位描述"`
	Status           string                 `json:"status" gorm:"size:20;default:open;comment:状态"`
	PublishDate      *time.Time             `json:"publish_date" gorm:"comment:发布日期"`
	CloseDate        *time.Time             `json:"close_date" gorm:"comment:关闭日期"`
	RecruiterID      uint                   `json:"recruiter_id" gorm:"comment:招聘负责人ID"`
	Recruiter        *Employee              `json:"recruiter,omitempty" gorm:"foreignKey:RecruiterID"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// Candidate represents job candidates
type Candidate struct {
	ID             uint                   `json:"id" gorm:"primaryKey"`
	Name           string                 `json:"name" gorm:"size:100;not null;comment:姓名"`
	Email          string                 `json:"email" gorm:"size:100;comment:邮箱"`
	Phone          string                 `json:"phone" gorm:"size:20;comment:手机号"`
	Gender         string                 `json:"gender" gorm:"size:10;comment:性别"`
	Age            int                    `json:"age" gorm:"comment:年龄"`
	Education      string                 `json:"education" gorm:"size:20;comment:学历"`
	School         string                 `json:"school" gorm:"size:100;comment:毕业学校"`
	Major          string                 `json:"major" gorm:"size:100;comment:专业"`
	Experience     int                    `json:"experience" gorm:"comment:工作年限"`
	ExpectedSalary float64                `json:"expected_salary" gorm:"type:decimal(10,2);comment:期望薪资"`
	Resume         string                 `json:"resume" gorm:"size:255;comment:简历文件"`
	RecruitmentID  uint                   `json:"recruitment_id" gorm:"comment:招聘职位ID"`
	Recruitment    *Recruitment           `json:"recruitment,omitempty" gorm:"foreignKey:RecruitmentID"`
	Status         string                 `json:"status" gorm:"size:20;default:pending;comment:状态"`
	Source         string                 `json:"source" gorm:"size:50;comment:简历来源"`
	Remark         string                 `json:"remark" gorm:"type:text;comment:备注"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// Performance represents performance evaluations
type Performance struct {
	ID           uint                   `json:"id" gorm:"primaryKey"`
	EmployeeID   uint                   `json:"employee_id" gorm:"not null;comment:员工ID"`
	Employee     *Employee              `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	Period       string                 `json:"period" gorm:"size:20;comment:考核周期"`
	Year         int                    `json:"year" gorm:"comment:考核年度"`
	Quarter      int                    `json:"quarter" gorm:"comment:考核季度"`
	Score        float64                `json:"score" gorm:"type:decimal(5,2);comment:考核得分"`
	Level        string                 `json:"level" gorm:"size:20;comment:考核等级"`
	Goals        string                 `json:"goals" gorm:"type:text;comment:工作目标"`
	Achievement  string                 `json:"achievement" gorm:"type:text;comment:工作成果"`
	Improvement  string                 `json:"improvement" gorm:"type:text;comment:改进建议"`
	Status       string                 `json:"status" gorm:"size:20;default:draft;comment:状态"`
	EvaluatorID  uint                   `json:"evaluator_id" gorm:"comment:评估人ID"`
	Evaluator    *Employee              `json:"evaluator,omitempty" gorm:"foreignKey:EvaluatorID"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// DepartmentAssignment 员工部门多重分配记录 - 复用现有Employee结构
type DepartmentAssignment struct {
	ID                  uint                   `json:"id" gorm:"primaryKey"`
	EmployeeID          uint                   `json:"employee_id" gorm:"not null;comment:员工ID"`
	Employee            *Employee              `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	DepartmentID        uint                   `json:"department_id" gorm:"not null;comment:部门ID"`
	Department          *Department            `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	
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

// OrganizationChangeLog represents organization structure change history
type OrganizationChangeLog struct {
	ID          uint                   `json:"id" gorm:"primaryKey"`
	EntityType  string                 `json:"entity_type" gorm:"size:20;comment:实体类型(department/position/job_level)"`
	EntityID    uint                   `json:"entity_id" gorm:"comment:实体ID"`
	Action      string                 `json:"action" gorm:"size:20;comment:操作类型(create/update/delete/move)"`
	OldValue    string                 `json:"old_value" gorm:"type:text;comment:旧值"`
	NewValue    string                 `json:"new_value" gorm:"type:text;comment:新值"`
	Field       string                 `json:"field" gorm:"size:50;comment:变更字段"`
	OperatorID  uint                   `json:"operator_id" gorm:"comment:操作人ID"`
	Operator    *Employee              `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
	Reason      string                 `json:"reason" gorm:"type:text;comment:变更原因"`
	CreatedAt   time.Time              `json:"created_at"`
}