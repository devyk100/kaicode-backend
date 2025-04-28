package types

type Job struct {
	SubmissionId string `json:"submission_id"`
	Code         string `json:"code"`
	Language     string `json:"language"`
	RoomName     string `json:"room_name"`
	Input        string `json:"input"`
	// some aws sqs handle, to delete this from the original SQS later
}
