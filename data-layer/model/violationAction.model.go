package model

import (
	"BE-DroneMEQ-Violation-Management/helpers/constants"
	"time"

	"gorm.io/gorm"
)

type ViolationAction struct {
	Pkid          int64                      `gorm:"primary_key;auto_increment" json:"pkid"`
	ViolationPkid int64                      `gorm:"not null" json:"violation_pkid"`        // Field is now exported (starts with uppercase)
	ActionTaken   string                     `gorm:"type:varchar(255)" json:"action_taken"` // Field is now exported (starts with uppercase)
	TakenByPkid   int64                      `gorm:"not null" json:"taken_by_pkid"`
	ActionTakenAt time.Time                  `json:"action_taken_at"`
	Outcome       constants.ViolationOutcome `json:"outcome"`
	// Audit Column
	OfficePkid  *int64         `gorm:"not null" json:"office_pkid"`
	CreatedBy   string         `json:"created_by"`
	CreatedHost string         `json:"created_host"`
	CreatedAt   time.Time      `gorm:"type:timestamp with time zone" json:"created_at"`
	UpdatedBy   *string        `json:"updated_by"`
	UpdatedHost *string        `json:"updated_host"`
	UpdatedAt   time.Time      `gorm:"type:timestamp with time zone" json:"updated_at"`
	DeletedBy   *string        `json:"deleted_by,omitempty"`
	DeletedHost *string        `json:"deleted_host,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	IsDeleted   bool           `gorm:"default:false" json:"is_deleted"`

	// Relationships
	Violation Violation `gorm:"foreignKey:ViolationPkid;references:Pkid"`
}

func (ViolationAction) TableName() string {
	return "violation_action"
}
