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

// Package worker define the worker interface
package worker

// Info the worker base information struct
type Info struct {
	WorkerID     int64
	OverLastTime int64
	BackLastTime int64
}

// Worker the interface of the worker
type Worker interface {
	// GetWorkerInfo get worker info
	GetWorkerInfo() (*Info, error)
	// WorkerIDBitLength get worker id bit length
	WorkerIDBitLength() byte
	// ReleaseWorkerID release worker id
	ReleaseWorkerID() error
	// UpdateOverLastTime update the over last time
	UpdateOverLastTime(overLastTime int64) error
	// UpdateBackLastTime update back last time
	UpdateBackLastTime(backLastTime int64) error
}

var workerIDAllocatorBuilders = map[string]func() Worker{}

// RegisterWorkerBuilder register worker builder
func RegisterWorkerBuilder(name string, builder func() Worker) {
	workerIDAllocatorBuilders[name] = builder
}

// NewWorker new a worker by name
func NewWorker(name string) Worker {
	return workerIDAllocatorBuilders[name]()
}
