package model

import (
	"time"

	"gorm.io/datatypes"
)

// Account represents the structure for the accounts table.
type Account struct {
	ID          int            `gorm:"primaryKey" json:"id"`
	AccountName string         `gorm:"column:accountname;unique;not null" json:"accountname"`
	AdminEmail  string         `json:"admin_email"`
	AdminPhone  string         `json:"admin_phone"`
	Config      datatypes.JSON `gorm:"type:jsonb;not null" json:"config"`
	CreatedAt   time.Time      `json:"created_at"`
}
