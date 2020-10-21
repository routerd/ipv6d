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

package v1

import "routerd.net/ipv6d/internal/machinery/runtime"

type TypeMeta struct {
	Kind    string `json:"kind"`
	Version string `json:"version"`
}

func (m *TypeMeta) GetObjectKind() runtime.ObjectKind { return m }
func (m *TypeMeta) GetKind() string                   { return m.Kind }
func (m *TypeMeta) SetKind(kind string)               { m.Kind = kind }
func (m *TypeMeta) GetVersion() string                { return m.Version }
func (m *TypeMeta) SetVersion(version string)         { m.Version = version }

func (m *TypeMeta) GetVersionKind() runtime.VersionKind {
	return runtime.VersionKind{Version: m.Version, Kind: m.Kind}
}

func (m *TypeMeta) SetVersionKind(vk runtime.VersionKind) {
	m.Version = vk.Version
	m.Kind = vk.Kind
}

type ObjectMeta struct {
	Name            string `json:"name"`
	Generation      int64  `json:"generation"`
	ResourceVersion string `json:"resourceVersion"`
}

func (m *ObjectMeta) GetName() string {
	return m.Name
}

func (m *ObjectMeta) SetName(name string) {
	m.Name = name
}

func (m *ObjectMeta) GetGeneration() int64 {
	return m.Generation
}

func (m *ObjectMeta) SetGeneration(generation int64) {
	m.Generation = generation
}

func (m *ObjectMeta) GetResourceVersion() string {
	return m.ResourceVersion
}

func (m *ObjectMeta) SetResourceVersion(resourceVersion string) {
	m.ResourceVersion = resourceVersion
}

type Object interface {
	GetName() string
	SetName(name string)
	GetGeneration() int64
	SetGeneration(generation int64)
	GetResourceVersion() string
	SetResourceVersion(resourceVersion string)
}
