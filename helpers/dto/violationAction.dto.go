package dto

type ViolationActionCreateRequest struct {
	ActionTaken string `gorm:"type:varchar(255)" json:"action_taken"`
	TakenByPkid int64  `gorm:"not null" json:"taken_by_pkid"`
	Outcome     string `gorm:"type:varchar(255)" json:"outcome"`
}

type ViolationActionCreateResponse struct {
	Pkid          int64  `gorm:"primary_key;auto_increment" json:"pkid"`
	ViolationPkid int64  `gorm:"not null" json:"violation_pkid"`
	ActionTaken   string `gorm:"type:varchar(255)" json:"action_taken"`
	TakenByPkid   int64  `gorm:"not null" json:"taken_by_pkid"`
	Outcome       string `gorm:"type:varchar(255)" json:"outcome"`
}
