package orchestrator

import (
	"context"
	"kc-backend/judge/state"
	"kc-backend/judge/worker"
	"sync"
)

// Orchestrator this is the main orchestrator, that queries the SQS, and puts it in the local queue, and also scales up or down the number of workers that it has running currently.
// This is with the assumption that SQS has a hide timing for certain queue entities, when they are read by someone, until a timeout they become unavailable, so running multiple such
// orchestrators should not be a problem
type Orchestrator struct {
	ctx         context.Context
	Workers     []worker.Worker
	workerMutex sync.Mutex
	pendingJobs *state.PendingJobs_t
}

var JOBS_PER_WORKER = 10
