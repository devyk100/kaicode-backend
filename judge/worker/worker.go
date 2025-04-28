package worker

import (
	"context"
	"github.com/redis/go-redis/v9"
	"kc-backend/judge/docker"
	//"kc-backend/judge/orchestrator"
	"kc-backend/judge/state"
)

// Worker This is the main worker that runs and processes all the code
type Worker struct {
	context         context.Context
	dockerContainer docker.Docker
	pendingJobs     *state.PendingJobs_t
	IsRunning       bool
	RedisClient     *redis.Client
}
