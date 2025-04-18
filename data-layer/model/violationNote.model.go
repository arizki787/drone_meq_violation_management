package model

import (
	"BE-DroneMEQ-Violation-Management/helpers/constants"
	"time"

	"gorm.io/gorm"
)

type ViolationNote struct {
	Pkid          int64                       `gorm:"primary_key;auto_increment" json:"pkid"`
	ViolationPkid int64                       `gorm:"not null" json:"violation_pkid"` // Field is now exported (starts with uppercase)
	Note          string                      `gorm:"type:text" json:"note"`          // Field is now exported (starts with uppercase)
	NoteType      constants.ViolationNoteType `json:"note_type"`
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
}

func (ViolationNote) TableName() string {
	return "violation_note"
}
