package constants

type ViolationOutcome string

const (
	OutcomeSuccessful ViolationOutcome = "SUCCESSFUL"
	OutcomeFailed     ViolationOutcome = "FAILED"
	OutcomePartial    ViolationOutcome = "PARTIAL"
)
