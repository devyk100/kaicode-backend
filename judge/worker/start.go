package worker

import (
	"context"
	"fmt"
	"kc-backend/judge/state"
)

// Start Starts the worker, basically starts the docker container that will be used
func (w *Worker) Start(ctx context.Context, pendingJobs *state.PendingJobs_t) error {
	w.context = ctx
	w.IsRunning = true
	w.pendingJobs = pendingJobs
	err := w.dockerContainer.StartContainer(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
