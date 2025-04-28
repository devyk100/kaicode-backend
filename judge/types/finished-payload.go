package types

type FinishedPayload struct {
	SubmissionId string `json:"submission_id"`
	Output       string `json:"output"`
	TimeTaken    int32  `json:"time_taken"`
	RoomName     string `json:"room_name"`
	SQSKey       string
}
