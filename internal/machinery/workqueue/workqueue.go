/*
Copyright 2020 The routerd Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package workqueue

import (
	"sync"
	"time"
)

type WorkQueue interface {
	Add(item interface{})
	AddAfter(item interface{}, after time.Duration)
	Get() (item interface{}, shutdown bool)
	Done(item interface{})
	ShutDown()
}

type workQueue struct {
	queue        []interface{}
	dirty        set
	processing   set
	cond         *sync.Cond
	shuttingDown bool
}

func NewWorkQueue() WorkQueue {
	return &workQueue{
		dirty:      set{},
		processing: set{},
		cond:       sync.NewCond(&sync.Mutex{}),
	}
}

// Add a new item to the workqueue
func (q *workQueue) Add(item interface{}) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	if q.shuttingDown {
		return
	}
	if q.dirty.has(item) {
		return
	}

	q.dirty.insert(item)
	if q.processing.has(item) {
		return
	}
	q.queue = append(q.queue, item)
	q.cond.Signal()
}

// Adds the item after the given duration
func (q *workQueue) AddAfter(item interface{}, after time.Duration) {
	go func() {
		<-time.After(after)
		q.Add(item)
	}()
}

// Get blocks until there is an item to process or the workqueue is shutting down.
func (q *workQueue) Get() (item interface{}, shutdown bool) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	for len(q.queue) == 0 && !q.shuttingDown {
		q.cond.Wait()
	}
	if len(q.queue) == 0 {
		// We must be shutting down.
		return nil, true
	}

	item, q.queue = q.queue[0], q.queue[1:]

	q.processing.insert(item)
	q.dirty.delete(item)

	return item, false
}

// Done removes an item from the workqueue.
func (q *workQueue) Done(item interface{}) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	q.processing.delete(item)
	if q.dirty.has(item) {
		q.queue = append(q.queue, item)
		q.cond.Signal()
	}
}

// stops the workqueue and shuts down all workers.
func (q *workQueue) ShutDown() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	q.shuttingDown = true
	q.cond.Broadcast()
}

// set helper
type set map[interface{}]struct{}

func (s set) has(item interface{}) bool {
	_, exists := s[item]
	return exists
}

func (s set) insert(item interface{}) {
	s[item] = struct{}{}
}

func (s set) delete(item interface{}) {
	delete(s, item)
}
