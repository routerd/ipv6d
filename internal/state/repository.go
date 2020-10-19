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

	v1 "github.com/routerd/ipv6d/api/v1"
	"github.com/routerd/ipv6d/internal/runtime"
	"github.com/routerd/ipv6d/internal/runtime/schema"
)

type Repository struct {
	scheme        *runtime.Scheme
	objVK, listVK schema.VersionKind
	mux           sync.Mutex
	data          map[string][]byte
}

type Object interface {
	v1.Object
	runtime.Object
}

type ObjectList interface {
	runtime.Object
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
		scheme: scheme,
		objVK:  objVK,
		listVK: listVK,
		data:   map[string][]byte{},
	}, nil
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
		return ErrNotFound{Key: key, VK: r.objVK}
	}
	if err := json.Unmarshal(data, obj); err != nil {
		return err
	}
	return nil
}

func (r *Repository) store(obj Object) error {
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

// Writer

func (r *Repository) Create(ctx context.Context, obj Object) error {
	if err := r.checkObjType(obj); err != nil {
		return err
	}

	r.mux.Lock()
	defer r.mux.Unlock()

	key := obj.GetName()
	if _, ok := r.data[key]; ok {
		return ErrAlreadyExists{Key: key, VK: r.objVK}
	}

	obj.SetGeneration(1)
	obj.SetResourceVersion("1")

	return r.store(obj)
}

func (r *Repository) Delete(ctx context.Context, obj Object) error {
	if err := r.checkObjType(obj); err != nil {
		return err
	}

	r.mux.Lock()
	defer r.mux.Unlock()

	key := obj.GetName()
	if _, ok := r.data[key]; ok {
		return ErrNotFound{Key: key, VK: r.objVK}
	}
	delete(r.data, key)
	return nil
}

func (r *Repository) Update(ctx context.Context, obj Object) error {
	if err := r.checkObjType(obj); err != nil {
		return err
	}

	r.mux.Lock()
	defer r.mux.Unlock()

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
		return ErrConflict{Key: key, VK: r.objVK}
	}

	// Update
	obj.SetGeneration(existingObj.GetGeneration() + 1)
	i, _ := strconv.Atoi(existingObj.GetResourceVersion())
	obj.SetResourceVersion(strconv.Itoa(i + 1))
	return r.store(obj)
}

type ErrConflict struct {
	Key string
	VK  schema.VersionKind
}

func (e ErrConflict) Error() string {
	return fmt.Sprintf("%s: %s conflicting resource version", e.VK, e.Key)
}

type ErrAlreadyExists struct {
	Key string
	VK  schema.VersionKind
}

func (e ErrAlreadyExists) Error() string {
	return fmt.Sprintf("%s: %s already exists", e.VK, e.Key)
}

type ErrNotFound struct {
	Key string
	VK  schema.VersionKind
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%s: %s not found", e.VK, e.Key)
}
