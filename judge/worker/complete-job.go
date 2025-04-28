package worker

import (
	"encoding/json"
	"fmt"
	"kc-backend/judge/sqs"
	"kc-backend/judge/types"
	"time"
)

func (w *Worker) completeJob(payload types.FinishedPayload) error {

	// execute the finishing thing
	// 1. Publish to redis, and set the key in redis with a TTL,
	// 2. Push to local queue, for persist
	fmt.Println("The room name that reached completeJob() is", payload.RoomName)
	val, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err.Error())
	}
	w.RedisClient.Publish(w.context, "roomName", val)
	w.RedisClient.Set(w.context, payload.RoomName+"-output", val, 10*time.Minute)

	fmt.Println("The key is ", payload.SQSKey)
	err = sqs.DeleteMessage(w.context, &payload.SQSKey)

	// delete this message from the sqs
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
