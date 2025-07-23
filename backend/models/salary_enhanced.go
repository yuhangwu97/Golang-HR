package models

import (
	"time"
	"gorm.io/gorm"
)

// SalaryComponent represents individual salary components (基本薪资构成)
type SalaryComponent struct {
	ID            uint                   `json:"id" gorm:"primaryKey"`
	Code          string                 `json:"code" gorm:"uniqueIndex;size:50;not null;comment:薪资组件编码"`
	Name          string                 `json:"name" gorm:"size:100;not null;comment:薪资组件名称"`
	Category      SalaryComponentCategory `json:"category" gorm:"size:20;not null;comment:组件分类"`
	Type          SalaryComponentType    `json:"type" gorm:"size:20;not null;comment:计算类型"`
	Formula       string                 `json:"formula" gorm:"type:text;comment:计算公式"`
	IsFixed       bool                   `json:"is_fixed" gorm:"default:false;comment:是否固定金额"`
	IsTaxable     bool                   `json:"is_taxable" gorm:"default:true;comment:是否计税"`
	IsRequired    bool                   `json:"is_required" gorm:"default:false;comment:是否必选项"`
	DefaultAmount float64                `json:"default_amount" gorm:"type:decimal(15,2);default:0;comment:默认金额"`
	MinAmount     *float64               `json:"min_amount" gorm:"type:decimal(15,2);comment:最小金额"`
	MaxAmount     *float64               `json:"max_amount" gorm:"type:decimal(15,2);comment:最大金额"`
	Sort          int                    `json:"sort" gorm:"default:0;comment:排序"`
	Status        string                 `json:"status" gorm:"size:20;default:active;comment:状态"`
	Description   string                 `json:"description" gorm:"type:text;comment:描述"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	DeletedAt     gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

// SalaryComponentCategory 薪资组件分类
type SalaryComponentCategory string

const (
	ComponentCategoryBase       SalaryComponentCategory = "base"        // 基本薪资
	ComponentCategoryAllowance  SalaryComponentCategory = "allowance"   // 津贴补助
	ComponentCategoryBonus      SalaryComponentCategory = "bonus"       // 奖金绩效
	ComponentCategoryDeduction  SalaryComponentCategory = "deduction"   // 扣款项目
	ComponentCategoryTax        SalaryComponentCategory = "tax"         // 税费
	ComponentCategoryBenefit    SalaryComponentCategory = "benefit"     // 福利
	ComponentCategoryInsurance  SalaryComponentCategory = "insurance"   // 保险
)

// SalaryComponentType 薪资组件计算类型
type SalaryComponentType string

const (
	ComponentTypeFixed      SalaryComponentType = "fixed"       // 固定金额
	ComponentTypePercentage SalaryComponentType = "percentage"  // 百分比
	ComponentTypeFormula    SalaryComponentType = "formula"     // 公式计算
	ComponentTypeManual     SalaryComponentType = "manual"      // 手动输入
)

// SalaryGrade 薪资等级 (Salary Grade/Band)
type SalaryGrade struct {
	ID          uint                   `json:"id" gorm:"primaryKey"`
	Code        string                 `json:"code" gorm:"uniqueIndex;size:20;not null;comment:等级编码"`
	Name        string                 `json:"name" gorm:"size:100;not null;comment:等级名称"`
	Level       int                    `json:"level" gorm:"not null;comment:等级层次"`
	MinSalary   float64                `json:"min_salary" gorm:"type:decimal(15,2);not null;comment:最低薪资"`
	MaxSalary   float64                `json:"max_salary" gorm:"type:decimal(15,2);not null;comment:最高薪资"`
	MidSalary   float64                `json:"mid_salary" gorm:"type:decimal(15,2);comment:中位薪资"`
	Currency    string                 `json:"currency" gorm:"size:3;default:CNY;comment:货币代码"`
	EffectiveDate *time.Time           `json:"effective_date" gorm:"comment:生效日期"`
	ExpiryDate    *time.Time           `json:"expiry_date" gorm:"comment:失效日期"`
	Status      string                 `json:"status" gorm:"size:20;default:active;comment:状态"`
	Description string                 `json:"description" gorm:"type:text;comment:描述"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

// SalaryStructure 薪资结构模板 (Salary Structure Template)
type SalaryStructure struct {
	ID           uint                   `json:"id" gorm:"primaryKey"`
	Code         string                 `json:"code" gorm:"uniqueIndex;size:50;not null;comment:结构编码"`
	Name         string                 `json:"name" gorm:"size:100;not null;comment:结构名称"`
	Description  string                 `json:"description" gorm:"type:text;comment:描述"`
	DepartmentID *uint                  `json:"department_id" gorm:"comment:适用部门ID"`
	Department   *Department            `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	PositionID   *uint                  `json:"position_id" gorm:"comment:适用职位ID"`
	Position     *Position              `json:"position,omitempty" gorm:"foreignKey:PositionID"`
	JobLevelID   *uint                  `json:"job_level_id" gorm:"comment:适用职级ID"`
	JobLevel     *JobLevel              `json:"job_level,omitempty" gorm:"foreignKey:JobLevelID"`
	SalaryGradeID *uint                 `json:"salary_grade_id" gorm:"comment:薪资等级ID"`
	SalaryGrade   *SalaryGrade          `json:"salary_grade,omitempty" gorm:"foreignKey:SalaryGradeID"`
	Components    []SalaryStructureComponent `json:"components,omitempty" gorm:"foreignKey:StructureID"`
	IsDefault     bool                  `json:"is_default" gorm:"default:false;comment:是否默认结构"`
	Status        string                `json:"status" gorm:"size:20;default:active;comment:状态"`
	EffectiveDate *time.Time            `json:"effective_date" gorm:"comment:生效日期"`
	ExpiryDate    *time.Time            `json:"expiry_date" gorm:"comment:失效日期"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
	DeletedAt     gorm.DeletedAt        `json:"deleted_at,omitempty" gorm:"index"`
}

// SalaryStructureComponent 薪资结构组件关联
type SalaryStructureComponent struct {
	ID                uint                   `json:"id" gorm:"primaryKey"`
	StructureID       uint                   `json:"structure_id" gorm:"not null;comment:薪资结构ID"`
	Structure         *SalaryStructure       `json:"structure,omitempty" gorm:"foreignKey:StructureID"`
	ComponentID       uint                   `json:"component_id" gorm:"not null;comment:薪资组件ID"`
	Component         *SalaryComponent       `json:"component,omitempty" gorm:"foreignKey:ComponentID"`
	DefaultValue      float64                `json:"default_value" gorm:"type:decimal(15,2);default:0;comment:默认值"`
	IsRequired        bool                   `json:"is_required" gorm:"default:false;comment:是否必填"`
	CanEdit           bool                   `json:"can_edit" gorm:"default:true;comment:是否可编辑"`
	Sort              int                    `json:"sort" gorm:"default:0;comment:排序"`
	CreatedAt         time.Time              `json:"created_at"`
}

// PayrollPeriod 薪资周期 (Pay Period)
type PayrollPeriod struct {
	ID          uint                   `json:"id" gorm:"primaryKey"`
	Name        string                 `json:"name" gorm:"size:100;not null;comment:周期名称"`
	PeriodType  PayrollPeriodType      `json:"period_type" gorm:"size:20;not null;comment:周期类型"`
	Year        int                    `json:"year" gorm:"not null;comment:年份"`
	Month       *int                   `json:"month" gorm:"comment:月份"`
	Quarter     *int                   `json:"quarter" gorm:"comment:季度"`
	StartDate   time.Time              `json:"start_date" gorm:"not null;comment:开始日期"`
	EndDate     time.Time              `json:"end_date" gorm:"not null;comment:结束日期"`
	PayDate     *time.Time             `json:"pay_date" gorm:"comment:发薪日期"`
	Status      PayrollPeriodStatus    `json:"status" gorm:"size:20;default:draft;comment:状态"`
	IsLocked    bool                   `json:"is_locked" gorm:"default:false;comment:是否锁定"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	DeletedAt   gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

// PayrollPeriodType 薪资周期类型
type PayrollPeriodType string

const (
	PeriodTypeMonthly   PayrollPeriodType = "monthly"    // 月薪
	PeriodTypeQuarterly PayrollPeriodType = "quarterly"  // 季薪
	PeriodTypeYearly    PayrollPeriodType = "yearly"     // 年薪
	PeriodTypeBonus     PayrollPeriodType = "bonus"      // 奖金
)

// PayrollPeriodStatus 薪资周期状态
type PayrollPeriodStatus string

const (
	PeriodStatusDraft     PayrollPeriodStatus = "draft"      // 草稿
	PeriodStatusOpen      PayrollPeriodStatus = "open"       // 开放
	PeriodStatusCalculated PayrollPeriodStatus = "calculated" // 已计算
	PeriodStatusReviewed  PayrollPeriodStatus = "reviewed"   // 已审核
	PeriodStatusApproved  PayrollPeriodStatus = "approved"   // 已批准
	PeriodStatusPaid      PayrollPeriodStatus = "paid"       // 已发放
	PeriodStatusClosed    PayrollPeriodStatus = "closed"     // 已关闭
)

// EnhancedSalary 增强版薪资记录
type EnhancedSalary struct {
	ID              uint                   `json:"id" gorm:"primaryKey"`
	EmployeeID      uint                   `json:"employee_id" gorm:"not null;comment:员工ID"`
	Employee        *Employee              `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	PayrollPeriodID uint                   `json:"payroll_period_id" gorm:"not null;comment:薪资周期ID"`
	PayrollPeriod   *PayrollPeriod         `json:"payroll_period,omitempty" gorm:"foreignKey:PayrollPeriodID"`
	StructureID     *uint                  `json:"structure_id" gorm:"comment:薪资结构ID"`
	Structure       *SalaryStructure       `json:"structure,omitempty" gorm:"foreignKey:StructureID"`
	
	// 薪资计算结果
	GrossSalary     float64                `json:"gross_salary" gorm:"type:decimal(15,2);default:0;comment:应发薪资"`
	TotalDeductions float64                `json:"total_deductions" gorm:"type:decimal(15,2);default:0;comment:总扣除"`
	NetSalary       float64                `json:"net_salary" gorm:"type:decimal(15,2);default:0;comment:实发薪资"`
	
	// 详细组件记录
	Components      []SalaryDetail         `json:"components,omitempty" gorm:"foreignKey:SalaryID"`
	
	// 计算和审批信息
	CalculatedBy    *uint                  `json:"calculated_by" gorm:"comment:计算人ID"`
	Calculator      *Employee              `json:"calculator,omitempty" gorm:"foreignKey:CalculatedBy"`
	CalculatedAt    *time.Time             `json:"calculated_at" gorm:"comment:计算时间"`
	ReviewedBy      *uint                  `json:"reviewed_by" gorm:"comment:审核人ID"`
	Reviewer        *Employee              `json:"reviewer,omitempty" gorm:"foreignKey:ReviewedBy"`
	ReviewedAt      *time.Time             `json:"reviewed_at" gorm:"comment:审核时间"`
	ApprovedBy      *uint                  `json:"approved_by" gorm:"comment:批准人ID"`
	Approver        *Employee              `json:"approver,omitempty" gorm:"foreignKey:ApprovedBy"`
	ApprovedAt      *time.Time             `json:"approved_at" gorm:"comment:批准时间"`
	
	// 状态和备注
	Status          SalaryStatus           `json:"status" gorm:"size:20;default:draft;comment:状态"`
	ReviewNotes     string                 `json:"review_notes" gorm:"type:text;comment:审核备注"`
	ApprovalNotes   string                 `json:"approval_notes" gorm:"type:text;comment:批准备注"`
	
	// 发放记录
	PayrollRecords  []EnhancedPayrollRecord `json:"payroll_records,omitempty" gorm:"foreignKey:SalaryID"`
	
	// 审计字段
	Version         int                    `json:"version" gorm:"default:1;comment:版本号"`
	PreviousVersion *uint                  `json:"previous_version" gorm:"comment:上一版本ID"`
	ChangeReason    string                 `json:"change_reason" gorm:"size:255;comment:变更原因"`
	
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	DeletedAt       gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

// SalaryStatus 薪资状态
type SalaryStatus string

const (
	SalaryStatusDraft     SalaryStatus = "draft"      // 草稿
	SalaryStatusCalculated SalaryStatus = "calculated" // 已计算
	SalaryStatusReviewed  SalaryStatus = "reviewed"   // 已审核
	SalaryStatusApproved  SalaryStatus = "approved"   // 已批准
	SalaryStatusPaid      SalaryStatus = "paid"       // 已发放
	SalaryStatusCancelled SalaryStatus = "cancelled"  // 已取消
	SalaryStatusRejected  SalaryStatus = "rejected"   // 已拒绝
)

// SalaryDetail 薪资详细组件记录
type SalaryDetail struct {
	ID                uint                   `json:"id" gorm:"primaryKey"`
	SalaryID          uint                   `json:"salary_id" gorm:"not null;comment:薪资记录ID"`
	Salary            *EnhancedSalary        `json:"salary,omitempty" gorm:"foreignKey:SalaryID"`
	ComponentID       uint                   `json:"component_id" gorm:"not null;comment:薪资组件ID"`
	Component         *SalaryComponent       `json:"component,omitempty" gorm:"foreignKey:ComponentID"`
	CalculatedValue   float64                `json:"calculated_value" gorm:"type:decimal(15,2);default:0;comment:计算值"`
	ManualValue       *float64               `json:"manual_value" gorm:"type:decimal(15,2);comment:手动调整值"`
	FinalValue        float64                `json:"final_value" gorm:"type:decimal(15,2);default:0;comment:最终值"`
	CalculationFormula string                `json:"calculation_formula" gorm:"type:text;comment:计算公式"`
	Notes             string                 `json:"notes" gorm:"type:text;comment:备注"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

// EnhancedPayrollRecord 增强版发放记录
type EnhancedPayrollRecord struct {
	ID              uint                   `json:"id" gorm:"primaryKey"`
	SalaryID        uint                   `json:"salary_id" gorm:"not null;comment:薪资记录ID"`
	Salary          *EnhancedSalary        `json:"salary,omitempty" gorm:"foreignKey:SalaryID"`
	PaymentBatchID  *uint                  `json:"payment_batch_id" gorm:"comment:批次ID"`
	PaymentBatch    *PaymentBatch          `json:"payment_batch,omitempty" gorm:"foreignKey:PaymentBatchID"`
	
	// 支付信息
	PaymentAmount   float64                `json:"payment_amount" gorm:"type:decimal(15,2);not null;comment:支付金额"`
	PaymentMethod   PaymentMethod          `json:"payment_method" gorm:"size:20;not null;comment:支付方式"`
	BankAccount     string                 `json:"bank_account" gorm:"size:50;comment:银行账号"`
	BankName        string                 `json:"bank_name" gorm:"size:100;comment:银行名称"`
	TransactionRef  string                 `json:"transaction_ref" gorm:"size:100;comment:交易参考号"`
	
	// 时间信息
	ScheduledDate   *time.Time             `json:"scheduled_date" gorm:"comment:计划支付日期"`
	ProcessedDate   *time.Time             `json:"processed_date" gorm:"comment:实际处理日期"`
	CompletedDate   *time.Time             `json:"completed_date" gorm:"comment:完成日期"`
	
	// 状态和处理
	Status          PayrollStatus          `json:"status" gorm:"size:20;default:pending;comment:发放状态"`
	FailureReason   string                 `json:"failure_reason" gorm:"type:text;comment:失败原因"`
	ProcessedBy     *uint                  `json:"processed_by" gorm:"comment:处理人ID"`
	Processor       *Employee              `json:"processor,omitempty" gorm:"foreignKey:ProcessedBy"`
	
	// 审计
	RetryCount      int                    `json:"retry_count" gorm:"default:0;comment:重试次数"`
	Notes           string                 `json:"notes" gorm:"type:text;comment:备注"`
	
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	DeletedAt       gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

// PaymentMethod 支付方式
type PaymentMethod string

const (
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer" // 银行转账
	PaymentMethodCash         PaymentMethod = "cash"          // 现金
	PaymentMethodCheck        PaymentMethod = "check"         // 支票
	PaymentMethodMobilePay    PaymentMethod = "mobile_pay"    // 移动支付
)

// PayrollStatus 发放状态
type PayrollStatus string

const (
	PayrollStatusPending   PayrollStatus = "pending"    // 待处理
	PayrollStatusProcessing PayrollStatus = "processing" // 处理中
	PayrollStatusCompleted PayrollStatus = "completed"  // 已完成
	PayrollStatusFailed    PayrollStatus = "failed"     // 失败
	PayrollStatusCancelled PayrollStatus = "cancelled"  // 已取消
)

// PaymentBatch 支付批次
type PaymentBatch struct {
	ID              uint                   `json:"id" gorm:"primaryKey"`
	BatchNumber     string                 `json:"batch_number" gorm:"uniqueIndex;size:50;not null;comment:批次号"`
	Name            string                 `json:"name" gorm:"size:100;not null;comment:批次名称"`
	PayrollPeriodID uint                   `json:"payroll_period_id" gorm:"not null;comment:薪资周期ID"`
	PayrollPeriod   *PayrollPeriod         `json:"payroll_period,omitempty" gorm:"foreignKey:PayrollPeriodID"`
	TotalAmount     float64                `json:"total_amount" gorm:"type:decimal(15,2);default:0;comment:总金额"`
	TotalRecords    int                    `json:"total_records" gorm:"default:0;comment:总记录数"`
	SuccessRecords  int                    `json:"success_records" gorm:"default:0;comment:成功记录数"`
	FailedRecords   int                    `json:"failed_records" gorm:"default:0;comment:失败记录数"`
	Status          PaymentBatchStatus     `json:"status" gorm:"size:20;default:draft;comment:批次状态"`
	ScheduledDate   *time.Time             `json:"scheduled_date" gorm:"comment:计划执行日期"`
	ProcessedDate   *time.Time             `json:"processed_date" gorm:"comment:实际处理日期"`
	CompletedDate   *time.Time             `json:"completed_date" gorm:"comment:完成日期"`
	CreatedBy       *uint                  `json:"created_by" gorm:"comment:创建人ID"`
	Creator         *Employee              `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	ProcessedBy     *uint                  `json:"processed_by" gorm:"comment:处理人ID"`
	Processor       *Employee              `json:"processor,omitempty" gorm:"foreignKey:ProcessedBy"`
	Notes           string                 `json:"notes" gorm:"type:text;comment:备注"`
	Records         []EnhancedPayrollRecord `json:"records,omitempty" gorm:"foreignKey:PaymentBatchID"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	DeletedAt       gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

// PaymentBatchStatus 支付批次状态
type PaymentBatchStatus string

const (
	BatchStatusDraft      PaymentBatchStatus = "draft"       // 草稿
	BatchStatusReady      PaymentBatchStatus = "ready"       // 就绪
	BatchStatusProcessing PaymentBatchStatus = "processing"  // 处理中
	BatchStatusCompleted  PaymentBatchStatus = "completed"   // 已完成
	BatchStatusFailed     PaymentBatchStatus = "failed"      // 失败
	BatchStatusCancelled  PaymentBatchStatus = "cancelled"   // 已取消
)

// SalaryAdjustment 薪资调整记录
type SalaryAdjustment struct {
	ID              uint                   `json:"id" gorm:"primaryKey"`
	EmployeeID      uint                   `json:"employee_id" gorm:"not null;comment:员工ID"`
	Employee        *Employee              `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	AdjustmentType  SalaryAdjustmentType   `json:"adjustment_type" gorm:"size:20;not null;comment:调整类型"`
	Reason          string                 `json:"reason" gorm:"size:255;not null;comment:调整原因"`
	EffectiveDate   time.Time              `json:"effective_date" gorm:"not null;comment:生效日期"`
	OldBaseSalary   float64                `json:"old_base_salary" gorm:"type:decimal(15,2);comment:原基本薪资"`
	NewBaseSalary   float64                `json:"new_base_salary" gorm:"type:decimal(15,2);comment:新基本薪资"`
	AdjustmentAmount float64               `json:"adjustment_amount" gorm:"type:decimal(15,2);comment:调整金额"`
	AdjustmentPercent float64              `json:"adjustment_percent" gorm:"type:decimal(5,2);comment:调整百分比"`
	ApprovedBy      *uint                  `json:"approved_by" gorm:"comment:批准人ID"`
	Approver        *Employee              `json:"approver,omitempty" gorm:"foreignKey:ApprovedBy"`
	ApprovedAt      *time.Time             `json:"approved_at" gorm:"comment:批准时间"`
	Status          string                 `json:"status" gorm:"size:20;default:pending;comment:状态"`
	Notes           string                 `json:"notes" gorm:"type:text;comment:备注"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	DeletedAt       gorm.DeletedAt         `json:"deleted_at,omitempty" gorm:"index"`
}

// SalaryAdjustmentType 薪资调整类型
type SalaryAdjustmentType string

const (
	AdjustmentTypePromotion  SalaryAdjustmentType = "promotion"   // 晋升
	AdjustmentTypeMarket     SalaryAdjustmentType = "market"      // 市场调整
	AdjustmentTypePerformance SalaryAdjustmentType = "performance" // 绩效调整
	AdjustmentTypeAnnual     SalaryAdjustmentType = "annual"      // 年度调整
	AdjustmentTypeCorrection SalaryAdjustmentType = "correction"  // 纠正调整
	AdjustmentTypeOther      SalaryAdjustmentType = "other"       // 其他
)

// Table names
func (SalaryComponent) TableName() string            { return "salary_components" }
func (SalaryGrade) TableName() string               { return "salary_grades" }
func (SalaryStructure) TableName() string          { return "salary_structures" }
func (SalaryStructureComponent) TableName() string { return "salary_structure_components" }
func (PayrollPeriod) TableName() string            { return "payroll_periods" }
func (EnhancedSalary) TableName() string           { return "enhanced_salaries" }
func (SalaryDetail) TableName() string             { return "salary_details" }
func (EnhancedPayrollRecord) TableName() string    { return "enhanced_payroll_records" }
func (PaymentBatch) TableName() string             { return "payment_batches" }
func (SalaryAdjustment) TableName() string         { return "salary_adjustments" }