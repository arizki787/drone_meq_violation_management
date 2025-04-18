package dto

import (
	"BE-DroneMEQ-Violation-Management/data-layer/model"
	"time"
)

// request
type ViolationNoteCreateRequest struct {
	Note     string `gorm:"type:text" json:"note"`
	NoteType string `gorm:"type:varchar(255)" json:"note_type"`
}

type ViolationNoteUpdateRequest struct {
	Note     string `gorm:"type:text" json:"note"`
	NoteType string `gorm:"type:varchar(255)" json:"note_type"`
	UpdatedBy   *string   `json:"updated_by" form:"updated_by"`
	UpdatedHost *string   `json:"updated_host"`
	UpdatedAt   time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
}

// response
type ViolationNoteCreateResponse struct {
	Pkid          int64  `gorm:"primary_key;auto_increment" json:"pkid"`
	ViolationPkid int64  `gorm:"not null" json:"violation_pkid"`
	Note          string `gorm:"type:text" json:"note"`
	NoteType      string `gorm:"type:varchar(255)" json:"note_type"`
	OfficePkid    int64  `gorm:"not null" json:"office_pkid"`
}

type ViolationNoteResponse struct {
	Pkid          int64  `gorm:"primary_key;auto_increment" json:"pkid"`
	ViolationPkid int64  `gorm:"not null" json:"violation_pkid"`
	Note          string `gorm:"type:text" json:"note"`
	NoteType      string `gorm:"type:varchar(255)" json:"note_type"`
	OfficePkid    int64  `gorm:"not null" json:"office_pkid"`
}

type ViolationNoteUpdateResponse struct {
	Pkid          int64  `gorm:"primary_key;auto_increment" json:"pkid"`
	ViolationPkid int64  `gorm:"not null" json:"violation_pkid"`
	Note          string `gorm:"type:text" json:"note"`
	NoteType      string `gorm:"type:varchar(255)" json:"note_type"`
	OfficePkid    int64  `gorm:"not null" json:"office_pkid"`
}

type GetAllViolationNoteRepositoryResponse struct {
	ViolationNotes []model.ViolationNote
	PaginationResponse
}

type ViolationNotePaginationResponse struct {
	Data []ViolationNoteResponse `json:"data"`
	PaginationResponse
}

type ViolationNoteDeleteResponse struct {
	Pkid          int64  `json:"pkid"`
	ViolationPkid int64  `json:"violation_pkid"`
	Note          string `json:"note"`
	NoteType      string `json:"note_type"`
	IsDeleted     bool   `json:"is_deleted"`
	OfficePkid    int64  `json:"office_pkid"`
}
