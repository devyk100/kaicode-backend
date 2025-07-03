package worker

import (
	"context"
	"crypto/tls"
	"kc-backend/judge/state"
	"os"

	"github.com/redis/go-redis/v9"
)

// Start Starts the worker, basically starts the docker container that will be used
func (w *Worker) Start(ctx context.Context, pendingJobs *state.PendingJobs_t) error {
	w.context = ctx
	w.IsRunning = true
	w.pendingJobs = pendingJobs
	rdb := redis.NewClient(&redis.Options{
		Addr:      os.Getenv("REDIS_URL"),
		Password:  os.Getenv("REDIS_PASSWORD"),
		TLSConfig: &tls.Config{},
		DB:        0,
	})
	w.RedisClient = rdb
	err := w.dockerContainer.StartContainer(ctx)
	if err != nil {
		// fmt.Println(err.Error())
		return err
	}
	return nil
}
