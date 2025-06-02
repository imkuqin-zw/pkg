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

// Package snowflake is a snowflake id generator.
package snowflake

import (
	"log/slog"
	"sync"
	"time"

	"github.com/imkuqin-zw/pkg/snowflake/worker"
	"github.com/pkg/errors"
)

// Config the snowflake config
type Config struct {
	BaseTime         int64 `default:"1582136402000"`
	SeqBitLength     byte  `default:"12"`
	MaxSeqNumber     int64
	MinSeqNumber     int64  `default:"5"`
	TopOverCostCount int    `default:"2000"`
	WorkerName       string `default:"static"`
	worker           worker.Worker
}

// WithWorker set worker
func (c *Config) WithWorker(w worker.Worker) {
	c.worker = w
}

func (c *Config) check() error {
	if c.MaxSeqNumber == 0 {
		c.MaxSeqNumber = (1 << c.SeqBitLength) - 1
	}
	if c.MinSeqNumber < 5 { // nolint: mnd
		return errors.Errorf("min seq number must  greater than 5")
	}
	if c.MinSeqNumber > c.MaxSeqNumber {
		return errors.Errorf("min seq number must  less than max seq number ")
	}
	if c.worker == nil {
		return errors.New("worker not set")
	}
	return nil
}

// Snowflake snowflake
type Snowflake struct {
	sync.Mutex
	baseTime          int64
	workerID          int64
	workerIDBitLength byte
	seqBitLength      byte
	maxSeqNumber      int64
	minSeqNumber      int64
	topOverCostCount  int

	lastTimeTick     int64
	currentSeqNumber int64
	timestampShift   byte

	isOverCost             bool
	overCostCountInOneTerm int

	turnBackTimeTick int64
	minBackTimeTick  int64
	turnBackIndex    int64

	worker worker.Worker
}

// NewSnowflake new snowflake by config
func NewSnowflake(cfg *Config) (*Snowflake, error) {
	if err := cfg.check(); err != nil {
		return nil, err
	}
	w := cfg.worker
	if w.WorkerIDBitLength()+cfg.SeqBitLength > 22 { // nolint: mnd
		return nil, errors.Errorf("worker id bit length + seq bit length must less than 22")
	}
	workerInfo, err := w.GetWorkerInfo()
	if err != nil {
		return nil, err
	}
	slog.Info("snowflake generate",
		slog.Int64("worker_id", workerInfo.WorkerID))
	snowflake := &Snowflake{
		baseTime:          cfg.BaseTime,
		workerID:          workerInfo.WorkerID,
		workerIDBitLength: w.WorkerIDBitLength(),
		seqBitLength:      cfg.SeqBitLength,
		maxSeqNumber:      cfg.MaxSeqNumber,
		minSeqNumber:      cfg.MinSeqNumber,
		topOverCostCount:  cfg.TopOverCostCount,
		timestampShift:    w.WorkerIDBitLength() + cfg.SeqBitLength,
		currentSeqNumber:  cfg.MinSeqNumber,
		worker:            w,
	}
	if workerInfo.OverLastTime >= snowflake.getCurrentTimeTick() {
		snowflake.lastTimeTick = workerInfo.OverLastTime
		snowflake.getNextTimeTick()
	}
	return snowflake, nil
}

// FetchID fetch next ID
func (w *Snowflake) FetchID() int64 {
	w.Lock()
	defer w.Unlock()
	if w.isOverCost {
		return w.nextOverCostID()
	}
	return w.nextNormalID()
}

// ReleaseWorkerID release worker id
func (w *Snowflake) ReleaseWorkerID() error {
	return w.worker.ReleaseWorkerID()
}

// WorkerID get snowflake worker id
func (w *Snowflake) WorkerID() int64 {
	return w.workerID
}

func (w *Snowflake) nextNormalID() int64 {
	currentTimeTick := w.getCurrentTimeTick()
	if currentTimeTick < w.lastTimeTick {
		if w.turnBackTimeTick < 1 {
			if !w.beginTurnBackAction() {
				return w.costID(w.lastTimeTick)
			}
		}
		return w.calcTurnBackID(w.turnBackTimeTick)
	}

	// 时间追平时
	if w.turnBackTimeTick > 0 {
		w.endTurnBackAction()
	}

	if currentTimeTick > w.lastTimeTick {
		w.lastTimeTick = currentTimeTick
		w.currentSeqNumber = w.minSeqNumber
		return w.costID(w.lastTimeTick)
	}

	if w.currentSeqNumber > w.maxSeqNumber {
		w.beginOverCostAction()
		return w.costID(w.lastTimeTick)
	}

	return w.costID(w.lastTimeTick)
}

func (w *Snowflake) nextOverCostID() int64 {
	currentTimeTick := w.getCurrentTimeTick()
	if currentTimeTick > w.lastTimeTick {
		w.endOverCostAction(currentTimeTick)
		return w.costID(w.lastTimeTick)
	}
	if w.overCostCountInOneTerm >= w.topOverCostCount {
		tick := w.getNextTimeTick()
		w.endOverCostAction(tick)
		return w.costID(w.lastTimeTick)
	}
	if w.currentSeqNumber > w.maxSeqNumber {
		w.beginOverCostAction()
		return w.costID(w.lastTimeTick)
	}
	return w.costID(w.lastTimeTick)
}

func (w *Snowflake) beginOverCostAction() {
	if err := w.worker.UpdateOverLastTime(w.lastTimeTick + 1); err != nil {
		slog.Error("fault to update over last time", "error", err)
		w.endOverCostAction(w.getNextTimeTick())
		return
	}
	w.lastTimeTick++
	w.currentSeqNumber = w.minSeqNumber
	w.isOverCost = true
	w.overCostCountInOneTerm++
}

func (w *Snowflake) endOverCostAction(currentTimeTick int64) {
	w.lastTimeTick = currentTimeTick
	w.currentSeqNumber = w.minSeqNumber
	w.isOverCost = false
	w.overCostCountInOneTerm = 0
}

func (w *Snowflake) beginTurnBackAction() bool {
	w.turnBackIndex++
	w.turnBackTimeTick = w.lastTimeTick - 1
	if w.minBackTimeTick >= w.turnBackTimeTick && w.turnBackIndex >= w.minSeqNumber {
		w.lastTimeTick = w.getNextTimeTick()
		w.endTurnBackAction()
		return false
	}
	if w.turnBackIndex == 1 {
		if err := w.worker.UpdateBackLastTime(w.lastTimeTick); err != nil {
			w.lastTimeTick = w.getNextTimeTick()
			w.endTurnBackAction()
			return false
		}
	}
	return true
}

func (w *Snowflake) endTurnBackAction() {
	w.turnBackTimeTick = 0
	w.turnBackIndex = 0
	info, _ := w.worker.GetWorkerInfo()
	w.minBackTimeTick = info.BackLastTime
}

func (w *Snowflake) getNextTimeTick() int64 {
	tempTimeTicker := w.getCurrentTimeTick()
	for tempTimeTicker <= w.lastTimeTick {
		time.Sleep(time.Millisecond)
		tempTimeTicker = w.getCurrentTimeTick()
	}
	return tempTimeTicker
}

func (w *Snowflake) getCurrentTimeTick() int64 {
	return time.Now().UnixMilli() - w.baseTime
}

func (w *Snowflake) costID(useTimeTick int64) int64 {
	result := (useTimeTick << w.timestampShift) + (w.workerID << w.seqBitLength) + w.currentSeqNumber
	w.currentSeqNumber++
	return result
}

func (w *Snowflake) calcTurnBackID(useTimeTick int64) int64 {
	result := (useTimeTick << w.timestampShift) + (w.workerID << w.seqBitLength) + w.turnBackIndex
	w.turnBackTimeTick--
	return result
}
