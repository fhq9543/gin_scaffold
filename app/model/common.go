package model

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Page struct {
	Limit  uint `json:"limit" uri:"limit" form:"limit"`
	Offset uint `json:"offset" uri:"offset" form:"offset"`
}
