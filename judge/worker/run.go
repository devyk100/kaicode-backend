package worker

import (
	"fmt"
	"kc-backend/judge/types"
)

func (w *Worker) Run() {
	for w.IsRunning {
		w.mainLoop()
	}
}

func (w *Worker) mainLoop() {
	////return nil
	if len(w.pendingJobs.Jobs) == 0 {
		return
	}
	var finishedPayload types.FinishedPayload
	if w.pendingJobs == nil {
		fmt.Println("WARNING: THE POINTER TO THE STATE PENDING JOB IS NIL")
		return
	}
	w.pendingJobs.Mutex.Lock()
	for k, j := range w.pendingJobs.Jobs {
		fmt.Printf("Running job %s\n", k)
		// exec this code
		delete(w.pendingJobs.Jobs, k)
		finishedPayload = w.ExecCode(k, j)
		break
	}
	w.pendingJobs.Mutex.Unlock()

	fmt.Println("Executed the code, this is it", finishedPayload)

	err := w.completeJob(finishedPayload)
	if err != nil {
		fmt.Println("ERROR: There was an error executing the code", err.Error())
		return
	}
}
