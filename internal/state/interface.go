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

	v1 "github.com/routerd/ipv6d/api/v1"
	"github.com/routerd/ipv6d/internal/runtime"
)

type Client interface {
	Reader
	Writer
}

type Watcher interface {
	Close() error
	ResultChan() <-chan Event
}

type Reader interface {
	Get(ctx context.Context, key string, obj Object) error
	List(ctx context.Context, listObj ObjectList) error
	Watch(ctx context.Context, obj Object) (Watcher, error)
}

type Writer interface {
	Create(ctx context.Context, obj Object) error
	Delete(ctx context.Context, obj Object) error
	Update(ctx context.Context, obj Object) error
}

type Object interface {
	v1.Object
	runtime.Object
}

type ObjectList interface {
	runtime.Object
}
