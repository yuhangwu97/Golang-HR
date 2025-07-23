package models

// type User struct {
// 	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
// 	Name        string         `gorm:"not null;size:100" json:"name" binding:"required"`
// 	Email       string         `gorm:"unique;not null;size:100" json:"email" binding:"required,email"`
// 	Password    string         `gorm:"not null;size:255" json:"password,omitempty" binding:"required,min=6"`
// 	Role        string         `gorm:"not null;size:50;default:employee" json:"role"`
// 	Status      string         `gorm:"size:20;default:active" json:"status"`
// 	Avatar      string         `gorm:"size:255" json:"avatar"`
// 	Phone       string         `gorm:"size:20" json:"phone"`
// 	EmployeeID  *uint          `gorm:"index" json:"employee_id"`
// 	Employee    *Employee      `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
// 	LastLoginAt *time.Time     `json:"last_login_at"`
// 	LoginIP     string         `gorm:"size:50" json:"login_ip"`
// 	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
// 	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
// 	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
// }

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
