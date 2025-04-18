package dto

import (
	"BE-DroneMEQ-Violation-Management/data-layer/model"
	"time"
)

// Request

type ViolationCreateRequest struct {
	ReportedByPkid int64     `json:"reported_by_pkid" form:"reported_by_pkid"`
	Latitude       float64   `json:"latitude" form:"latitude"`
	Longitude      float64   `json:"longitude" form:"longitude"`
	ReportedAt     time.Time `json:"reported_at" form:"reported_at"`
	ViolationType  string    `json:"violation_type" form:"violation_type"`
	Status         string    `json:"status" form:"status"`
	Severity       string    `json:"severity" form:"severity"`
	Description    string    `json:"description" form:"description"`
	OfficePkid     int64     `gorm:"not null" json:"office_pkid"`
}

type ViolationUpdateRequest struct {
	ReportedByPkid int64     `json:"reported_by_pkid" form:"reported_by_pkid"`
	Latitude       float64   `json:"latitude" form:"latitude"`
	Longitude      float64   `json:"longitude" form:"longitude"`
	ReportedAt     time.Time `json:"reported_at" form:"reported_at"`
	ViolationType  string    `json:"violation_type" form:"violation_type"`
	Status         string    `json:"status" form:"status"`
	Severity       string    `json:"severity" form:"severity"`
	Description    string    `json:"description" form:"description"`
	UpdatedBy      *string   `json:"updated_by" form:"updated_by"`
	UpdatedHost    *string   `json:"updated_host"`
	UpdatedAt      time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
}

// Response
type GetAllViolationRepositoryResponse struct {
	Violations []model.Violation
	PaginationResponse
}

type ViolationResponse struct {
	Pkid           int64     `json:"pkid"`
	ReportedByPkid int64     `json:"reported_by_pkid"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	ViolationType  string    `json:"violation_type"`
	ReportedAt     time.Time `json:"reported_at"`
	Status         string    `json:"status"`
	Severity       string    `json:"severity"`
	Description    string    `json:"description"`
	OfficePkid     int64     `gorm:"not null" json:"office_pkid"`
	MediaURL       string    `json:"media_url" form:"media_url"`
	MediaType      string    `json:"media_type" form:"media_type"`
}

type ViolationMediaResponse struct {
	MediaURL  string `json:"media_url"`
	MediaType string `json:"media_type"`
}

type ViolationResponseGetId struct {
	Violation      model.Violation          `json:"violation"`
	ViolationMedia []ViolationMediaResponse `json:"violation_media"`
}

type GetViolationByIdResponse struct {
	Pkid           int64     `json:"pkid"`
	ReportedByPkid int64     `json:"reported_by_pkid"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	ViolationType  string    `json:"violation_type"`
	ReportedAt     time.Time `json:"reported_at"`
	Status         string    `json:"status"`
	Severity       string    `json:"severity"`
	Description    string    `json:"description"`
	OfficePkid     int64     `gorm:"not null" json:"office_pkid"`
	MediaURL       string    `json:"media_url" form:"media_url"`
	MediaType      string    `json:"media_type" form:"media_type"`
}

type ViolationUpdateResponse struct {
	Pkid           int64     `json:"pkid"`
	Status         string    `json:"status"`
	ReportedByPkid int64     `json:"reported_by_pkid"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	ViolationType  string    `json:"violation_type"`
	ReportedAt     time.Time `json:"reported_at"`
	Severity       string    `json:"severity"`
	Description    string    `json:"description"`
}

type ViolationPaginationResponse struct {
	Data []ViolationResponse `json:"data"`
	PaginationResponse
}

type ViolationDeleteResponse struct {
	Pkid      int64 `json:"pkid"`
	IsDeleted bool  `json:"is_deleted"`
}
