package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Problem struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	Title            string         `gorm:"size:200;not null" json:"title"`
	Description      string         `gorm:"type:text;not null" json:"description"`
	Difficulty       string         `gorm:"size:20;not null" json:"difficulty"` // Easy, Medium, Hard
	Tags             string         `gorm:"size:255" json:"tags"`               // JSON array
	TimeLimit        int            `gorm:"not null;default:5000" json:"time_limit"`   // 毫秒
	MemoryLimit      int            `gorm:"not null;default:256" json:"memory_limit"`  // MB
	FunctionTemplate string         `gorm:"type:text" json:"function_template"` // 展示给用户的函数签名模板
	DriverCode       string         `gorm:"type:text" json:"-"`                 // 拼接在用户代码后的驱动代码（含main函数）
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	TestCases        []TestCase     `gorm:"foreignKey:ProblemID" json:"-"`
}

type TestCase struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ProblemID uint           `gorm:"not null;index" json:"problem_id"`
	Input     string         `gorm:"type:text;not null" json:"input"`
	Output    string         `gorm:"type:text;not null" json:"output"`
	IsPublic  bool           `gorm:"default:false" json:"is_public"` // 是否公开
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Submission struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UserID     uint           `gorm:"not null;index" json:"user_id"`
	ProblemID  uint           `gorm:"not null;index" json:"problem_id"`
	Code       string         `gorm:"type:text;not null" json:"code"`
	Language   string         `gorm:"size:20;not null" json:"language"` // Go, Python, etc
	Status     string         `gorm:"size:20;not null" json:"status"`   // Pending, Running, Accepted, Wrong Answer, etc
	Result     string         `gorm:"type:text" json:"result"`          // 详细结果信息
	TimeUsed   int            `json:"time_used"`                       // 毫秒
	MemoryUsed int            `json:"memory_used"`                     // KB
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
	User       *User          `gorm:"foreignKey:UserID" json:"-"`
	Problem    *Problem       `gorm:"foreignKey:ProblemID" json:"-"`
}
