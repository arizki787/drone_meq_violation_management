package external

import "time"

// MetadataDTO represents user management metadata for tracking creation, update, and deletion details
type MetadataDTO struct {
	OfficePKID  int64      `json:"office_pkid"`  // Office identifier from token
	CreatedBy   string     `json:"created_by"`   // Username from token
	CreatedHost string     `json:"created_host"` // Host/IP address from token or request
	CreatedAt   time.Time  `json:"created_at"`   // Timestamp of creation
	UpdatedBy   *string    `json:"updated_by"`   // Optional: Username of the person who updated (nullable)
	UpdatedHost *string    `json:"updated_host"` // Optional: Host/IP of the person who updated (nullable)
	UpdatedAt   *time.Time `json:"updated_at"`   // Optional: Timestamp of the last update (nullable)
	DeletedBy   *string    `json:"deleted_by"`   // Optional: Username of the person who deleted (nullable)
	DeletedHost *string    `json:"deleted_host"` // Optional: Host/IP of the person who deleted (nullable)
	DeletedAt   *time.Time `json:"deleted_at"`   // Optional: Timestamp of deletion (nullable)
	IsDeleted   bool       `json:"is_deleted"`
}