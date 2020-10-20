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

	v1 "github.com/routerd/ipv6d/api/v1"
	"github.com/routerd/ipv6d/internal/runtime"
)

var testScheme = runtime.NewScheme()

type testObjectList struct {
	v1.TypeMeta `json:",inline"`
	Items       []testObject `json:"items"`
}

func (o *testObjectList) DeepCopy() *testObjectList {
	new := &testObjectList{}
	if err := copier.Copy(new, o); err != nil {
		panic(err)
	}
	return new
}

func (o *testObjectList) DeepCopyObject() runtime.Object {
	return o.DeepCopy()
}

type testObject struct {
	v1.TypeMeta   `json:",inline"`
	v1.ObjectMeta `json:"metadata"`
}

func (o *testObject) DeepCopy() *testObject {
	new := &testObject{}
	if err := copier.Copy(new, o); err != nil {
		panic(err)
	}
	return new
}

func (o *testObject) DeepCopyObject() runtime.Object {
	return o.DeepCopy()
}

func init() {
	testScheme.AddKnownTypes("v1", &testObject{}, &testObjectList{})
}
