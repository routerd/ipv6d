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

package errors

import (
	"errors"
	"fmt"

	"routerd.net/ipv6d/internal/machinery/runtime"
)

type ErrConflict struct {
	Key string
	VK  runtime.VersionKind
}

func (e ErrConflict) Error() string {
	return fmt.Sprintf("%s: %s conflicting resource version", e.VK, e.Key)
}

func IsConflict(err error) bool {
	_, ok := errors.Unwrap(err).(ErrConflict)
	return ok
}

type ErrAlreadyExists struct {
	Key string
	VK  runtime.VersionKind
}

func (e ErrAlreadyExists) Error() string {
	return fmt.Sprintf("%s: %s already exists", e.VK, e.Key)
}

func IsAlreadyExists(err error) bool {
	_, ok := errors.Unwrap(err).(ErrAlreadyExists)
	return ok
}

type ErrNotFound struct {
	Key string
	VK  runtime.VersionKind
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("%s: %s not found", e.VK, e.Key)
}

func IsNotFound(err error) bool {
	_, ok := errors.Unwrap(err).(ErrNotFound)
	return ok
}
