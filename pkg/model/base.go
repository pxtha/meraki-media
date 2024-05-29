package model

import (
	uuid "github.com/google/uuid"
	"gitlab.com/merakilab9/meracore/ginext"
	"gorm.io/gorm"
	"time"
)

type Pagination struct {
	Page     int
	PageSize int
}

type UriParse struct {
	ID []string `json:"id" uri:"id"`
}

type BaseModel struct {
	ID        uuid.UUID       `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatorID *uuid.UUID      `json:"creator_id,omitempty" gorm:"type:uuid"`
	UpdaterID *uuid.UUID      `json:"updater_id,omitempty" gorm:"type:uuid"`
	CreatedAt time.Time       `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" swaggertype:"string"`
}

type Base struct {
	ID        uuid.UUID       `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time       `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type BaseModelStripped struct {
	ID uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()" json:"id"`
}

type APIResponseList struct {
	Data []interface{} `json:"data"`
	Meta *ginext.Pager `json:"meta"`
}

type APIResponseOne struct {
	Data interface{} `json:"data"`
}
