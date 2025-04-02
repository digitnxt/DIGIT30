package model

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

// Account represents the structure for the accounts table.
type Account struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	AccountName string         `gorm:"column:accountname;unique;not null" json:"accountname"`
	AdminEmail  string         `json:"admin_email"`
	AdminPhone  string         `json:"admin_phone"`
	Config      datatypes.JSON `gorm:"type:jsonb;not null" json:"config"`
	CreatedAt   time.Time      `json:"created_at"`
}
