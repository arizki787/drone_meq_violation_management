package constants

type ViolationNoteType string

const (
	NoteObservation ViolationNoteType = "OBSERVATION"
	NoteFollowUp    ViolationNoteType = "FOLLOW_UP"
	NoteResolution  ViolationNoteType = "RESOLUTION"
	NoteOther       ViolationNoteType = "OTHER"
)
