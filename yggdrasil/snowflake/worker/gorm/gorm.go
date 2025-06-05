// Package gorm is a snowflake worker builder based on gorm
package gorm

import (
	basicwoker "github.com/imkuqin-zw/pkg/basic/snowflake/worker"
	"github.com/imkuqin-zw/pkg/basic/snowflake/worker/gorm"
	"github.com/imkuqin-zw/pkg/yggdrasil/snowflake/worker"
	"github.com/imkuqin-zw/yggdrasil/pkg/config"
	"github.com/imkuqin-zw/yggdrasil/pkg/logger"
)

func init() {
	worker.RegisterWorkerBuilder("gorm", NewWorker)
}

// NewWorker returns a new gorm worker
func NewWorker() (basicwoker.Worker, error) {
	cfg := gorm.Config{}
	if err := config.Get("snowflake.worker.gorm").Scan(&cfg); err != nil {
		logger.ErrorField("fault to load snowflake worker config", logger.Err(err))
		return nil, err
	}
	return gorm.NewWorker(&cfg)
}
