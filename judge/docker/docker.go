package docker

import (
	"context"
	"os"
	"time"

	"github.com/docker/docker/client"
)

var (
	MAX_CONTAINERS       = 3
	MAX_TIMEOUT          = time.Second * 40
	IMAGE_NAME           = "code-exec-engine"
	MAX_PROCESSES  int64 = 130
	SigChan        chan os.Signal
	Running        bool = true
)

type Docker struct {
	containerId string
	ctx         context.Context
	cli         *client.Client
}
