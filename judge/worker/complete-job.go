package worker

import (
	"fmt"
	"kc-backend/judge/sqs"
	"kc-backend/judge/types"
)

func (w *Worker) completeJob(payload types.FinishedPayload) error {

	// execute the finishing thing
	// 1. Publish to redis, and set the key in redis with a TTL,
	// 2. Push to local queue, for persist

	fmt.Println("The key is ", payload.SQSKey)
	err := sqs.DeleteMessage(w.context, &payload.SQSKey)

	// delete this message from the sqs
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
