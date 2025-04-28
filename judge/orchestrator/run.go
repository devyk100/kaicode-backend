package orchestrator

import (
	"encoding/json"
	"fmt"
	"kc-backend/judge/sqs"
	"kc-backend/judge/types"
	"kc-backend/judge/worker"
)

func (o *Orchestrator) Run() {
	// exit gracefully
	defer func() {
		for _, workerTemp := range o.Workers {
			workerTemp.Exit()
		}
	}()

	for {
		o.mainLoop()
	}
}

func (o *Orchestrator) mainLoop() {
	message, err := sqs.ReceiveMessage()
	if err != nil {
		fmt.Println("mainLoop(): ", err.Error())
		return
	}
	if len(message) == 0 {
		fmt.Println("mainLoop(): Message is empty.")
		return
	}

	for _, message := range message {

		fmt.Println("mainLoop(): Message received", *message.Body, "the receipt handle of this message is", *message.ReceiptHandle)

		var job types.Job

		err = json.Unmarshal([]byte(*message.Body), &job)
		if err != nil {
			fmt.Println("mainLoop(): Error unmarshalling message", err.Error())
			return
		}

		o.pendingJobs.Mutex.Lock()
		o.pendingJobs.Jobs[*message.ReceiptHandle] = job
		o.pendingJobs.Mutex.Unlock()
	}

	for {
		// scale up the number of workers
		if len(o.pendingJobs.Jobs) > JOBS_PER_WORKER*len(o.Workers) {
			// scale up
			workerTemp := worker.Worker{}
			o.workerMutex.Lock()
			o.Workers = append(o.Workers, workerTemp)
			err := workerTemp.Start(o.ctx, o.pendingJobs)
			if err != nil {
				fmt.Println("mainLoop(): Error starting worker", err.Error())
				return
			}
			go workerTemp.Run()
			o.workerMutex.Unlock()
		} else {
			break
		}
	}

	for {
		// scale down the number of workers
		if len(o.pendingJobs.Jobs) < JOBS_PER_WORKER*(len(o.Workers)-1) {
			// scale down
			o.workerMutex.Lock()
			workerTemp := o.Workers[len(o.Workers)-1]
			o.Workers = o.Workers[:len(o.Workers)-1] // delete last
			workerTemp.Exit()
			o.workerMutex.Unlock()
		} else {
			break
		}
	}
}
