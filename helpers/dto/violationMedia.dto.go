package dto

type ViolationMediaCreateRequest struct {
	MediaURL  string `json:"media_url" form:"media_url"`
	MediaType string `json:"media_type" form:"media_type"`
}

type ViolationMediaCreateResponse struct {
	Pkid          int64  `json:"pkid"`
	ViolationPkid int64  `json:"violation_pkid"`
	MediaURL      string `json:"media_url"`
	MediaType     string `json:"media_type"`
}
