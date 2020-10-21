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

package controller

import (
	"time"

	"github.com/go-logr/logr"

	"routerd.net/ipv6d/internal/machinery/workqueue"
)

// Result of a reconcile operation
type Result struct {
	Requeue      bool
	RequeueAfter time.Duration
}

type Reconciler interface {
	Reconcile(key string) (Result, error)
}

type Controller struct {
	log        logr.Logger
	reconciler Reconciler
	queue      workqueue.WorkQueue
}

func NewController(
	log logr.Logger,
	reconciler Reconciler,
) *Controller {
	return &Controller{
		log:        log,
		reconciler: reconciler,
		queue:      workqueue.NewWorkQueue(),
	}
}

func (c *Controller) Run(stopCh <-chan struct{}) {
	defer c.queue.ShutDown()
	go c.worker()
	<-stopCh
}

func (c *Controller) Add(key string) {
	c.queue.Add(key)
}

func (c *Controller) AddAfter(key string, after time.Duration) {
	c.queue.AddAfter(key, after)
}

func (c *Controller) worker() {
	for c.processNext() {
	}
}

func (c *Controller) processNext() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Done(item)

	res, err := c.reconciler.Reconcile(item.(string))
	if err != nil {
		c.log.Error(err, "Reconcile error", "request", item)
		c.queue.AddAfter(item, 5*time.Second)
		return true
	}

	if res.RequeueAfter > 0 {
		c.queue.AddAfter(item, res.RequeueAfter)
		return true
	}
	if res.Requeue {
		c.queue.Add(item)
		return true
	}

	c.log.Info("Successfully Reconciled", "request", item)
	return true
}
