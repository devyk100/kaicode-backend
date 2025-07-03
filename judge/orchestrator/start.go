package orchestrator

import (
	"context"
	"kc-backend/judge/sqs"
	"kc-backend/judge/state"
)

// Start starts the orchestrator, given the  context, and the state that you're going to use for maintaining the pending job reference
func (o *Orchestrator) Start(ctx context.Context, pendingJobRef *state.PendingJobs_t) {
	o.ctx = ctx
	sqs.InitSQSClient()
	o.pendingJobs = pendingJobRef
	if o.pendingJobs == nil {
		// fmt.Println("The pending Jobs reference passed was nil, make sure the memory was handled properly")
	}
}
