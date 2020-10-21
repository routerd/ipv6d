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
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"routerd.net/ipv6d/internal/machinery/errors"
	"routerd.net/ipv6d/internal/machinery/runtime"
)

var (
	_ Client = &Repository{}
)

type Repository struct {
	scheme        *runtime.Scheme
	objVK, listVK runtime.VersionKind
	mux           sync.Mutex
	eventHub      *eventHub
	data          map[string][]byte
}

func NewRepository(scheme *runtime.Scheme, obj Object, listObj ObjectList) (*Repository, error) {
	objVK, err := scheme.VersionKind(obj)
	if err != nil {
		return nil, err
	}

	listVK, err := scheme.VersionKind(listObj)
	if err != nil {
		return nil, err
	}

	return &Repository{
		scheme:   scheme,
		objVK:    objVK,
		listVK:   listVK,
		eventHub: newEventHub(),
		data:     map[string][]byte{},
	}, nil
}

func (r *Repository) Run(stopCh <-chan struct{}) {
	r.eventHub.Run(stopCh)
}

func (r *Repository) checkObjType(obj Object) error {
	vk, err := r.scheme.VersionKind(obj)
	if err != nil {
		return err
	}

	if vk == r.objVK {
		return nil
	}
	return fmt.Errorf("wrong type given to repository want %s, got %s", r.objVK, vk)
}

func (r *Repository) checkObjListType(obj ObjectList) error {
	vk, err := r.scheme.VersionKind(obj)
	if err != nil {
		return err
	}

	if vk == r.listVK {
		return nil
	}
	return fmt.Errorf("wrong type given to repository want %s, got %s", r.listVK, vk)
}

func (r *Repository) load(key string, obj Object) error {
	data, ok := r.data[key]
	if !ok {
		return errors.ErrNotFound{Key: key, VK: r.objVK}
	}
	if err := json.Unmarshal(data, obj); err != nil {
		return err
	}
	return nil
}

func (r *Repository) store(obj Object) error {
	obj.GetObjectKind().SetVersionKind(r.objVK)

	key := obj.GetName()
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	r.data[key] = data
	return nil
}

// Reader

func (r *Repository) Get(ctx context.Context, key string, obj Object) error {
	if err := r.checkObjType(obj); err != nil {
		return err
	}

	r.mux.Lock()
	defer r.mux.Unlock()

	return r.load(key, obj)
}

func (r *Repository) List(ctx context.Context, listObj ObjectList) error {
	if err := r.checkObjListType(listObj); err != nil {
		return err
	}

	r.mux.Lock()
	defer r.mux.Unlock()

	rv := reflect.ValueOf(listObj).Elem()
	for _, entryData := range r.data {
		obj, err := r.scheme.New(r.objVK)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(entryData, obj); err != nil {
			return err
		}

		rv.FieldByName("Items").Set(
			reflect.Append(rv.FieldByName("Items"), reflect.ValueOf(obj).Elem()),
		)
	}
	return nil
}

func (r *Repository) Watch(ctx context.Context, obj Object) (Watcher, error) {
	return r.eventHub.Register(), nil
}

// Writer

func (r *Repository) Create(ctx context.Context, obj Object) error {
	if err := r.checkObjType(obj); err != nil {
		return err
	}

	r.mux.Lock()
	defer r.mux.Unlock()

	key := obj.GetName()
	if _, ok := r.data[key]; ok {
		return errors.ErrAlreadyExists{Key: key, VK: r.objVK}
	}

	obj.SetGeneration(1)
	obj.SetResourceVersion("1")

	if err := r.store(obj); err != nil {
		return err
	}
	r.eventHub.Broadcast(nil, obj)
	return nil
}

func (r *Repository) Delete(ctx context.Context, obj Object) error {
	if err := r.checkObjType(obj); err != nil {
		return err
	}

	r.mux.Lock()
	defer r.mux.Unlock()

	key := obj.GetName()
	if err := r.load(key, obj); err != nil {
		return err
	}
	delete(r.data, key)
	r.eventHub.Broadcast(obj, nil)
	return nil
}

func (r *Repository) Update(ctx context.Context, obj Object) error {
	if err := r.checkObjType(obj); err != nil {
		return err
	}

	r.mux.Lock()
	defer r.mux.Unlock()

	// Get existing
	existing, err := r.scheme.New(r.objVK)
	if err != nil {
		return err
	}
	existingObj := existing.(Object)

	key := obj.GetName()
	if err := r.load(key, existingObj); err != nil {
		return err
	}

	// Check ResourceVersion
	if existingObj.GetResourceVersion() != obj.GetResourceVersion() {
		return errors.ErrConflict{Key: key, VK: r.objVK}
	}

	// Update
	obj.SetGeneration(existingObj.GetGeneration() + 1)
	i, _ := strconv.Atoi(existingObj.GetResourceVersion())
	obj.SetResourceVersion(strconv.Itoa(i + 1))

	// Ensure Status is not updated, if the field exists
	statusField := reflect.ValueOf(obj).Elem().FieldByName("Status")
	if statusField.IsValid() {
		statusField.Set(
			reflect.ValueOf(existingObj).Elem().FieldByName("Status"),
		)
	}

	if err := r.store(obj); err != nil {
		return err
	}
	r.eventHub.Broadcast(existingObj, obj)
	return nil
}

func (r *Repository) UpdateStatus(ctx context.Context, obj Object) error {
	if err := r.checkObjType(obj); err != nil {
		return err
	}

	r.mux.Lock()
	defer r.mux.Unlock()

	// Get existing
	existing, err := r.scheme.New(r.objVK)
	if err != nil {
		return err
	}
	existingObj := existing.(Object)

	key := obj.GetName()
	if err := r.load(key, existingObj); err != nil {
		return err
	}

	// Check ResourceVersion
	if existingObj.GetResourceVersion() != obj.GetResourceVersion() {
		return errors.ErrConflict{Key: key, VK: r.objVK}
	}

	// Ensure ObjectMeta and Spec is not updated
	reflect.ValueOf(obj).Elem().FieldByName("ObjectMeta").Set(
		reflect.ValueOf(existingObj).Elem().FieldByName("ObjectMeta"),
	)
	specField := reflect.ValueOf(obj).Elem().FieldByName("Spec")
	if specField.IsValid() {
		specField.Set(
			reflect.ValueOf(existingObj).Elem().FieldByName("Spec"),
		)
	}

	// Update (no generation update)
	i, _ := strconv.Atoi(existingObj.GetResourceVersion())
	obj.SetResourceVersion(strconv.Itoa(i + 1))

	if err := r.store(obj); err != nil {
		return err
	}
	r.eventHub.Broadcast(existingObj, obj)
	return nil
}
