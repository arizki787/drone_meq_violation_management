package model

import (
	"BE-DroneMEQ-Violation-Management/helpers/constants"
	"time"

	"gorm.io/gorm"
)

type Violation struct {
	Pkid           int64                       `gorm:"primary_key;auto_increment" json:"pkid"`
	ReportedByPkid int64                       `gorm:"not null" json:"reported_by_pkid"` // Field is now exported (starts with uppercase)
	Latitude       float64                     `gorm:"not null" json:"latitude"`         // Fixed missing closing quote for json tag
	Longitude      float64                     `gorm:"not null" json:"longitude"`        // Fixed missing closing quote for json tag
	ReportedAt     time.Time                   `gorm:"not null" json:"reported_at"`      // Field is now exported (starts with uppercase)
	ViolationType  constants.ViolationType     `gorm:"not null" json:"violation_type"`
	Status         constants.ViolationStatus   `gorm:"not null" json:"status"`
	Severity       constants.ViolationSeverity `gorm:"not null" json:"severity"`
	Description    string                      `gorm:"type:text" json:"description"`

	// Audit Columns
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
	ViolationAction         []ViolationAction         `gorm:"foreignKey:ViolationPkid;references:Pkid"`
	ViolationViolationMedia []ViolationViolationMedia `gorm:"foreignKey:ViolationPkid;references:Pkid"`
	ViolationNote           []ViolationNote           `gorm:"foreignKey:ViolationPkid;references:Pkid"`
}

func (Violation) TableName() string {
	return "violations"
}
