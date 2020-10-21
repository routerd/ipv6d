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
	"github.com/jinzhu/copier"

	"routerd.net/ipv6d/internal/machinery/runtime"
)

var testScheme = runtime.NewScheme()

type testObjectList struct {
	runtime.VersionKind `json:",inline"`
	Items               []testObject `json:"items"`
}

func (obj *testObjectList) GetVersionKind() runtime.VersionKind {
	return obj.VersionKind
}

func (obj *testObjectList) SetVersionKind(vk runtime.VersionKind) {
	obj.VersionKind = vk
}

func (obj *testObjectList) GetObjectKind() runtime.ObjectKind {
	return obj
}

func (obj *testObjectList) DeepCopy() *testObjectList {
	new := &testObjectList{}
	if err := copier.Copy(new, obj); err != nil {
		panic(err)
	}
	return new
}

func (obj *testObjectList) DeepCopyObject() runtime.Object {
	return obj.DeepCopy()
}

type testObject struct {
	runtime.VersionKind `json:",inline"`
	ObjectMeta          `json:"metadata"`
}

func (obj *testObject) GetVersionKind() runtime.VersionKind {
	return obj.VersionKind
}

func (obj *testObject) SetVersionKind(vk runtime.VersionKind) {
	obj.VersionKind = vk
}

func (obj *testObject) GetObjectKind() runtime.ObjectKind {
	return obj
}

func (obj *testObject) DeepCopy() *testObject {
	new := &testObject{}
	if err := copier.Copy(new, obj); err != nil {
		panic(err)
	}
	return new
}

func (obj *testObject) DeepCopyObject() runtime.Object {
	return obj.DeepCopy()
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

func init() {
	testScheme.AddKnownTypes("v1", &testObject{}, &testObjectList{})
}
