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

// Package snowflake_test test the snowflake.
package snowflake_test

import (
	"testing"
	"time"

	"github.com/imkuqin-zw/pkg/snowflake"
	"github.com/imkuqin-zw/pkg/snowflake/worker/gorm"
	xgorm "github.com/imkuqin-zw/yggdrasil/contrib/gorm"
	_ "github.com/imkuqin-zw/yggdrasil/contrib/gorm/driver/sqlite"
	"github.com/imkuqin-zw/yggdrasil/pkg/config"
	"github.com/imkuqin-zw/yggdrasil/pkg/logger"
	gormDB "gorm.io/gorm"
)

var db *gormDB.DB

func init() {
	_ = config.Set("snowflake.workerName", "gorm")
	_ = config.Set("gorm.default.driver", "sqlite")
	_ = config.Set("gorm.default.nameStrategy.singularTable", true)
	_ = config.Set("gorm.default.dsn", "file:shared_db?mode=memory&cache=shared")
	logger.SetLevel(logger.LvWarn)
	db = xgorm.NewDB("default")
	if err := db.AutoMigrate(&gorm.SnowflakeWorker{}); err != nil {
		logger.ErrorField("fault to migrate table")
	}
}

func truncate(db *gormDB.DB) {
	_ = db.Exec("DELETE FROM snowflake_worker").Error
}

func TestSnowflake_normal(t *testing.T) {
	defer truncate(db)
	snowflake.Init()
	ids := make(map[int64]struct{}, 10000)
	for i := 0; i < 10000; i++ {
		id := snowflake.FetchID()
		if _, ok := ids[id]; ok {
			t.Errorf("duplicate id %d", id)
			return
		}
		ids[id] = struct{}{}
	}
	_ = snowflake.Release()
	snowflake.Init()
	time.Sleep(time.Millisecond)
	defer func() {
		_ = snowflake.Release()
	}()
	for i := 0; i < 10000; i++ {
		id := snowflake.FetchID()
		if _, ok := ids[id]; ok {
			t.Errorf("duplicate id 2 %d %d", id, i)
			return
		}
		ids[id] = struct{}{}
	}
}
