package types

type FinishedPayload struct {
	SubmissionId string `json:"submission_id"`
	Output       string `json:"output"`
	TimeTaken    int32  `json:"time_taken"`
	SQSKey       string
}
