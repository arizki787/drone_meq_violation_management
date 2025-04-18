package model

import (
	"time"

	"gorm.io/gorm"
)

type ViolationMedia struct {
	Pkid           int64     `gorm:"primary_key;auto_increment" json:"pkid"`
	ViolationPkid  int64     `gorm:"not null" json:"violation_pkid"` // Field is now exported (starts with uppercase)
	MediaURL       string    `json:"media_url"`                      // Field is now exported (starts with uppercase)
	MediaType      string    `json:"media_type"`
	UploadedByPkid int64     `gorm:"not null" json:"uploaded_by_pkid"`
	UploadedAt     time.Time `json:"uploaded_at"`
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

	//relationships
	Violation Violation `gorm:"foreignKey:ViolationPkid;references:Pkid"`
}

func (ViolationMedia) TableName() string {
	return "violation_media"
}
