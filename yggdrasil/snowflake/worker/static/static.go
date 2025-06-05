// Package static is a snowflake worker builder based on static config
package static

import (
	basicwoker "github.com/imkuqin-zw/pkg/basic/snowflake/worker"
	"github.com/imkuqin-zw/pkg/basic/snowflake/worker/static"
	"github.com/imkuqin-zw/pkg/yggdrasil/snowflake/worker"
	"github.com/imkuqin-zw/yggdrasil/pkg/config"
	"github.com/imkuqin-zw/yggdrasil/pkg/logger"
)

func init() {
	worker.RegisterWorkerBuilder("static", NewWorker)
}

// NewWorker returns a new static worker
func NewWorker() (basicwoker.Worker, error) {
	cfg := static.Config{}
	if err := config.Get("snowflake.worker.static").Scan(&cfg); err != nil {
		logger.ErrorField("fault to load snowflake worker config", logger.Err(err))
		return nil, err
	}
	return static.NewWorker(&cfg)
}
