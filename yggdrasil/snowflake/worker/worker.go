// Package worker implements worker builder registry for snowflake worker.
package worker

import (
	"github.com/imkuqin-zw/pkg/basic/snowflake/worker"
	"github.com/pkg/errors"
)

// Builder worker构造函数
type Builder func() (worker.Worker, error)

var workerBuilders = map[string]Builder{}

// RegisterWorkerBuilder 注册worker构造函数
func RegisterWorkerBuilder(name string, builder Builder) {
	workerBuilders[name] = builder
}

// NewWorker 创建worker
func NewWorker(name string) (worker.Worker, error) {
	builder, ok := workerBuilders[name]
	if !ok {
		return nil, errors.Errorf("worker builder %s not found", name)
	}
	return builder()
}
