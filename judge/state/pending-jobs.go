package state

import (
	"kc-backend/judge/types"
	"sync"
)

type PendingJobs_t struct {
	Jobs  map[string]types.Job
	Mutex sync.Mutex
}

var PendingJobs *PendingJobs_t
