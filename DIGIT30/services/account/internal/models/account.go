package models

import (
	"encoding/json"
	"time"
)

// Account represents the structure of an account record.
type Account struct {
	ID          int             `json:"id"`
	AccountName string          `json:"accountname"`
	AdminEmail  string          `json:"admin_email"`
	AdminPhone  string          `json:"admin_phone"`
	Config      json.RawMessage `json:"config"`
	CreatedAt   time.Time       `json:"created_at"`
}
