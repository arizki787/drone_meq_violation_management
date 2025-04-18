package constants

type ViolationStatus string

const (
	StatusReported    ViolationStatus = "REPORTED"
	StatusUnderReview ViolationStatus = "UNDER_REVIEW"
	StatusResolved    ViolationStatus = "RESOLVED"
	StatusDismissed   ViolationStatus = "DISMISSED"
)
