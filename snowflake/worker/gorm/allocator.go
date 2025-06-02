// Copyright 2022 The imkuqin-zw Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package gorm implements the worker  allocator
package gorm

import (
	"sync"

	"github.com/google/uuid"
	"github.com/imkuqin-zw/pkg/snowflake/worker"
	"github.com/imkuqin-zw/yggdrasil/pkg/config"
	"github.com/imkuqin-zw/yggdrasil/pkg/logger"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func init() {
	worker.RegisterWorkerBuilder("gorm", NewWorker)
}

// Config the worker allocator config
type Config struct {
	WorkerIDBitLength int8   `default:"6"`
	Business          string `default:"default"`
	DBName            string `default:"default"`
}

// WorkerIDAllocator gorm worker id allocator
type WorkerIDAllocator struct {
	sync.Mutex
	info              *worker.Info
	business          string
	workerIDBitLength byte
	flag              string
	data              *snowflakeWorkerData
}

// NewWorkerWithDB new worker with db
func NewWorkerWithDB(db *gorm.DB) worker.Worker {
	cfg := Config{}
	if err := config.Get("snowflake.worker.gorm").Scan(&cfg); err != nil {
		logger.FatalField("fault to load snowflake worker config", logger.Err(err))
	}
	w := &WorkerIDAllocator{
		workerIDBitLength: byte(cfg.WorkerIDBitLength),
		flag:              uuid.Must(uuid.NewUUID()).String(),
		business:          cfg.Business,
		data: &snowflakeWorkerData{
			db:          db,
			staticDB:    true,
			maxWorkerID: (1 << cfg.WorkerIDBitLength) - 1,
		},
	}
	return w
}

// NewWorker new worker
func NewWorker() worker.Worker {
	cfg := Config{}
	if err := config.Get("snowflake.worker.gorm").Scan(&cfg); err != nil {
		logger.FatalField("fault to load snowflake worker config", logger.Err(err))
	}
	w := &WorkerIDAllocator{
		workerIDBitLength: byte(cfg.WorkerIDBitLength),
		flag:              uuid.Must(uuid.NewUUID()).String(),
		business:          cfg.Business,
		data: &snowflakeWorkerData{
			dbName:      cfg.DBName,
			maxWorkerID: (1 << cfg.WorkerIDBitLength) - 1,
		},
	}
	return w
}

// GetWorkerInfo get worker info
func (w *WorkerIDAllocator) GetWorkerInfo() (*worker.Info, error) {
	w.Lock()
	defer w.Unlock()
	if w.info != nil {
		return w.info, nil
	}
	w.data.openDB()
	defer w.data.closeDB()
	workerInfo, err := w.data.getReleasedWorkerInfo(w.business, w.flag)
	if err != nil {
		return nil, err
	}
	if workerInfo != nil {
		w.info = workerInfo
		return workerInfo, nil
	}
	var i int
	for i < 3 {
		workerInfo, err = w.data.getNewWorker(w.business, w.flag)
		if err == nil {
			w.info = workerInfo
			return workerInfo, nil
		}

		if !errors.Is(err, errWorkerIDExist) {
			return nil, err
		}
		i++
	}
	return nil, errors.WithStack(err)
}

// WorkerIDBitLength get worker id bit length
func (w *WorkerIDAllocator) WorkerIDBitLength() byte {
	return w.workerIDBitLength
}

// ReleaseWorkerID release worker id
func (w *WorkerIDAllocator) ReleaseWorkerID() error {
	w.Lock()
	defer w.Unlock()
	if w.info == nil {
		return nil
	}
	w.data.openDB()
	defer w.data.closeDB()
	if err := w.data.releaseWorkerID(w.info.WorkerID, w.business, w.flag); err != nil {
		return err
	}
	w.info = nil
	return nil
}

// UpdateOverLastTime update over last time
func (w *WorkerIDAllocator) UpdateOverLastTime(overLastTime int64) error {
	w.Lock()
	defer w.Unlock()
	if w.info == nil {
		return errors.WithStack(errWorkerIDNotExist)
	}
	w.data.openDB()
	defer w.data.closeDB()
	if err := w.data.updateOverLastTime(w.info.WorkerID, w.business, w.flag, overLastTime); err != nil {
		return err
	}
	w.info.OverLastTime = overLastTime
	return nil
}

// UpdateBackLastTime update back last time
func (w *WorkerIDAllocator) UpdateBackLastTime(backLastTime int64) error {
	w.Lock()
	defer w.Unlock()
	if w.info == nil {
		return errors.WithStack(errWorkerIDNotExist)
	}
	w.data.openDB()
	defer w.data.closeDB()
	if err := w.data.updateBackLastTime(w.info.WorkerID, w.business, w.flag, backLastTime); err != nil {
		return err
	}
	w.info.BackLastTime = backLastTime
	return nil
}
