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

// Package gorm_test implements the gorm worker allocator test
package gorm_test

import (
	"testing"

	"github.com/imkuqin-zw/pkg/snowflake/worker"
	"github.com/imkuqin-zw/pkg/snowflake/worker/gorm"
	xgorm "github.com/imkuqin-zw/yggdrasil/contrib/gorm"
	_ "github.com/imkuqin-zw/yggdrasil/contrib/gorm/driver/sqlite"
	"github.com/imkuqin-zw/yggdrasil/pkg/config"
	"github.com/imkuqin-zw/yggdrasil/pkg/logger"
	"github.com/stretchr/testify/require"

	gormDB "gorm.io/gorm"
)

var db *gormDB.DB

func init() {
	_ = config.Set("gorm.default.driver", "sqlite")
	_ = config.Set("gorm.default.nameStrategy.singularTable", true)
	_ = config.Set("gorm.default.dsn", "file:shared_db?mode=memory&cache=shared")

	logger.SetLevel(logger.LvInfo)
	db = xgorm.NewDB("default")
	if err := db.AutoMigrate(&gorm.SnowflakeWorker{}); err != nil {
		logger.ErrorField("fault to migrate table")
	}
}

func truncate(db *gormDB.DB) {
	_ = db.Exec("DELETE FROM snowflake_worker").Error
}

func testSingle(t *testing.T, allocator worker.Worker) {
	defer truncate(db)
	workerID1, err := allocator.GetWorkerInfo()
	require.NoError(t, err)
	workerID2, err := allocator.GetWorkerInfo()
	require.NoError(t, err)
	require.Equal(t, workerID1.WorkerID, workerID2.WorkerID)
	require.NoError(t, allocator.ReleaseWorkerID())
	workerID3, err := allocator.GetWorkerInfo()
	require.NoError(t, err)
	require.Equal(t, workerID1.WorkerID, workerID3.WorkerID)
}

func testMulti(t *testing.T, worker1, worker2, worker3 worker.Worker) {
	workerID1, err := worker1.GetWorkerInfo()
	require.NoError(t, err)
	workerID2, err := worker2.GetWorkerInfo()
	require.NoError(t, err)
	workerID3, err := worker3.GetWorkerInfo()
	require.NoError(t, err)

	require.NotEqual(t, workerID1.WorkerID, workerID2.WorkerID, "worker id should be different")
	require.NotEqual(t, workerID1.WorkerID, workerID3.WorkerID, "worker id should be different")
	require.NotEqual(t, workerID2.WorkerID, workerID3.WorkerID, "worker id should be different")

	require.NoError(t, worker1.ReleaseWorkerID())
	require.NoError(t, worker2.ReleaseWorkerID())
	workerID4, err := worker1.GetWorkerInfo()
	require.NoError(t, err)
	workerID5, err := worker2.GetWorkerInfo()
	require.NoError(t, err)
	require.NotEqual(t, workerID4.WorkerID, workerID5.WorkerID, "worker id should be different")
	require.Equal(t, workerID1.WorkerID, workerID4.WorkerID, "worker id should be same")
	require.Equal(t, workerID2.WorkerID, workerID5.WorkerID, "worker id should be same")
}

func TestWorkerIDAllocator_StaticSingle(t *testing.T) {
	allocator := gorm.NewWorkerWithDB(db)
	testSingle(t, allocator)
}

func TestWorkerIDAllocator_StaticMulti(t *testing.T) {
	defer truncate(db)
	worker1 := gorm.NewWorkerWithDB(db)
	worker2 := gorm.NewWorkerWithDB(db)
	worker3 := gorm.NewWorkerWithDB(db)
	testMulti(t, worker1, worker2, worker3)
}

func TestWorkerIDAllocator_Single(t *testing.T) {
	allocator := gorm.NewWorker()
	testSingle(t, allocator)
}

func TestWorkerIDAllocator_Multi(t *testing.T) {
	defer truncate(db)
	worker1 := gorm.NewWorker()
	worker2 := gorm.NewWorker()
	worker3 := gorm.NewWorker()
	testMulti(t, worker1, worker2, worker3)
}
