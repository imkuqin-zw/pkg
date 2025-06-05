// Package snowflake provides a snowflake id generator.
package snowflake

import (
	"github.com/imkuqin-zw/pkg/basic/snowflake"
	"github.com/imkuqin-zw/yggdrasil/pkg/config"
	"github.com/imkuqin-zw/yggdrasil/pkg/logger"
)

var global *snowflake.Snowflake

// Init 初始化snowflake
func Init() {
	cfg := &snowflake.Config{}
	if err := config.Get("snowflake").Scan(cfg); err != nil {
		logger.FatalField("fault to load snowflake config", logger.Err(err))
	}
	var err error
	global, err = snowflake.NewSnowflake(cfg)
	if err != nil {
		logger.FatalField("fault to init snowflake", logger.Err(err))
	}
}

// FetchID 获取id
func FetchID() int64 {
	return global.FetchID()
}

// WorkerID 获取当前worker id
func WorkerID() int64 {
	return global.WorkerID()
}

// Release 释放当前worker id
func Release() error {
	if global == nil {
		return nil
	}
	err := global.ReleaseWorkerID()
	if err != nil {
		logger.ErrorField("fault to release worker id", logger.Err(err))
	}
	global = nil
	return err
}
