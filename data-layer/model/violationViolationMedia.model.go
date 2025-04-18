package model

type ViolationViolationMedia struct {
	Pkid               int64 `gorm:"primary_key;auto_increment" json:"pkid"`
	ViolationPkid      int64 `gorm:"not null" json:"violation_pkid"`
	ViolationMediaPkid int64 `gorm:"not null" json:"violation_media_pkid"`

	// relationships
	Violation      Violation      `gorm:"foreignKey:ViolationPkid;references:Pkid"`
	ViolationMedia ViolationMedia `gorm:"foreignKey:ViolationMediaPkid;references:Pkid"`
}

func (ViolationViolationMedia) TableName() string {
	return "violation_violation_media"
}
