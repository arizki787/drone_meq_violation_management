package constants

type ViolationType string

const (
	ViolationIllegalFishing     ViolationType = "ILLEGAL_FISHING"
	ViolationPollution          ViolationType = "POLLUTION"
	ViolationHabitatDestruction ViolationType = "HABITAT_DESTRUCTION"
	ViolationAssetMissing       ViolationType = "ASSET_MISSING"
	ViolationOther              ViolationType = "OTHER"
)
