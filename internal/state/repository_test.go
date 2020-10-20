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

package state

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/routerd/ipv6d/api/v1"
)

func TestRepository(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		r, err := NewRepository(testScheme, &testObject{}, &testObjectList{})
		require.NoError(t, err)

		// Run
		stopCh := make(chan struct{})
		go r.Run(stopCh)
		defer close(stopCh)

		// Inject some data
		r.data["test123"] = []byte(`{"metadata": {"name":"test123"}}`)

		// List
		ctx := context.Background()
		obj := &testObject{}
		err = r.Get(ctx, "test123", obj)
		require.NoError(t, err)

		assert.Equal(t, "test123", obj.Name)
	})

	t.Run("List", func(t *testing.T) {
		r, err := NewRepository(testScheme, &testObject{}, &testObjectList{})
		require.NoError(t, err)

		// Run
		stopCh := make(chan struct{})
		go r.Run(stopCh)
		defer close(stopCh)

		// Inject some data
		r.data["test123"] = []byte(`{"metadata": {"name":"test123"}}`)

		// List
		ctx := context.Background()
		list := &testObjectList{}
		err = r.List(ctx, list)
		require.NoError(t, err)

		if assert.Len(t, list.Items, 1) {
			assert.Equal(t, "test123", list.Items[0].Name)
		}
	})

	t.Run("Watch", func(t *testing.T) {
		r, err := NewRepository(testScheme, &testObject{}, &testObjectList{})
		require.NoError(t, err)

		// Run
		stopCh := make(chan struct{})
		go r.Run(stopCh)
		defer close(stopCh)

		// Watch
		ctx := context.Background()
		watcher, err := r.Watch(ctx, nil)
		require.NoError(t, err)

		var events []Event
		var wg sync.WaitGroup
		wg.Add(1) // wait for 1 event
		go func() {
			for event := range watcher.ResultChan() {
				events = append(events, event)
				wg.Done()
			}
		}()

		// generate a "Added" event
		obj := &testObject{ObjectMeta: v1.ObjectMeta{Name: "test3000"}}
		require.NoError(t, r.Create(ctx, obj))

		// Assertions
		wg.Wait()
		if assert.Len(t, events, 1) {
			assert.Equal(t, Event{New: obj}, events[0])
			assert.Equal(t, Added, events[0].Type())
		}
	})

	t.Run("Create", func(t *testing.T) {
		r, err := NewRepository(testScheme, &testObject{}, &testObjectList{})
		require.NoError(t, err)

		// Run
		stopCh := make(chan struct{})
		go r.Run(stopCh)
		defer close(stopCh)

		// Create
		ctx := context.Background()
		obj := &testObject{ObjectMeta: v1.ObjectMeta{Name: "test3000"}}
		require.NoError(t, r.Create(ctx, obj))

		assert.Equal(t,
			`{"kind":"testObject","version":"v1","metadata":{"name":"test3000","generation":1,"resourceVersion":"1"}}`,
			string(r.data["test3000"]))
	})

	t.Run("Update", func(t *testing.T) {
		r, err := NewRepository(testScheme, &testObject{}, &testObjectList{})
		require.NoError(t, err)

		// Run
		stopCh := make(chan struct{})
		go r.Run(stopCh)
		defer close(stopCh)

		// Inject some data
		r.data["test123"] = []byte(`{"metadata": {"name":"test123","generation":3,"resourceVersion":"53"}}`)

		// Update
		ctx := context.Background()
		obj := &testObject{ObjectMeta: v1.ObjectMeta{
			Name:            "test123",
			Generation:      3,
			ResourceVersion: "53",
		}}
		require.NoError(t, r.Update(ctx, obj))

		assert.Equal(t,
			`{"kind":"testObject","version":"v1","metadata":{"name":"test123","generation":4,"resourceVersion":"54"}}`,
			string(r.data["test123"]))
	})

	t.Run("Update", func(t *testing.T) {
		r, err := NewRepository(testScheme, &testObject{}, &testObjectList{})
		require.NoError(t, err)

		// Run
		stopCh := make(chan struct{})
		go r.Run(stopCh)
		defer close(stopCh)

		// Inject some data
		r.data["test123"] = []byte(`{"metadata": {"name":"test123","generation":3,"resourceVersion":"53"}}`)

		// Update
		ctx := context.Background()
		obj := &testObject{ObjectMeta: v1.ObjectMeta{
			Name:            "test123",
			Generation:      3,
			ResourceVersion: "53",
		}}
		require.NoError(t, r.Update(ctx, obj))

		assert.Equal(t,
			`{"kind":"testObject","version":"v1","metadata":{"name":"test123","generation":4,"resourceVersion":"54"}}`,
			string(r.data["test123"]))
	})

	t.Run("UpdateStatus", func(t *testing.T) {
		r, err := NewRepository(testScheme, &testObject{}, &testObjectList{})
		require.NoError(t, err)

		// Run
		stopCh := make(chan struct{})
		go r.Run(stopCh)
		defer close(stopCh)

		// Inject some data
		r.data["test123"] = []byte(`{"metadata": {"name":"test123","generation":3,"resourceVersion":"53"}}`)

		// UpdateStatus
		ctx := context.Background()
		obj := &testObject{ObjectMeta: v1.ObjectMeta{
			Name:            "test123",
			Generation:      3,
			ResourceVersion: "53",
		}}
		require.NoError(t, r.UpdateStatus(ctx, obj))

		assert.Equal(t,
			`{"kind":"testObject","version":"v1","metadata":{"name":"test123","generation":3,"resourceVersion":"54"}}`,
			string(r.data["test123"]))
	})

	t.Run("Delete", func(t *testing.T) {
		r, err := NewRepository(testScheme, &testObject{}, &testObjectList{})
		require.NoError(t, err)

		// Run
		stopCh := make(chan struct{})
		go r.Run(stopCh)
		defer close(stopCh)

		// Inject some data
		r.data["test123"] = []byte(`{"metadata": {"name":"test123","generation":3,"resourceVersion":"53"}}`)

		// Delete
		ctx := context.Background()
		obj := &testObject{ObjectMeta: v1.ObjectMeta{
			Name:            "test123",
			Generation:      3,
			ResourceVersion: "53",
		}}
		require.NoError(t, r.Delete(ctx, obj))

		assert.Zero(t, r.data["test123"])
	})
}
