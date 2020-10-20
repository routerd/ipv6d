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

package runtime

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/routerd/ipv6d/internal/runtime/schema"
)

type Object interface {
	GetObjectKind() schema.ObjectKind
	DeepCopyObject() Object
}

type Scheme struct {
	vkToType map[schema.VersionKind]reflect.Type
	typeToVK map[reflect.Type]schema.VersionKind
}

func NewScheme() *Scheme {
	return &Scheme{
		vkToType: map[schema.VersionKind]reflect.Type{},
		typeToVK: map[reflect.Type]schema.VersionKind{},
	}
}

func ensureStructPointer(obj Object) reflect.Type {
	rt := reflect.TypeOf(obj)
	if rt.Kind() != reflect.Ptr {
		panic("All types must be pointers to structs.")
	}

	rt = rt.Elem()
	if rt.Kind() != reflect.Struct {
		panic("All types must be pointers to structs.")
	}
	return rt
}

func (s *Scheme) AddKnownTypes(version string, types ...Object) {
	for _, obj := range types {
		rt := ensureStructPointer(obj)
		s.AddKnownTypeWithKind(schema.VersionKind{
			Version: version,
			Kind:    rt.Name(),
		}, obj)
	}
}

func (s *Scheme) AddKnownTypeWithKind(vk schema.VersionKind, obj Object) {
	if len(vk.Version) == 0 {
		panic("Version is required on all types.")
	}
	rt := ensureStructPointer(obj)

	s.vkToType[vk] = rt
	s.typeToVK[rt] = vk
}

func (s *Scheme) New(vk schema.VersionKind) (Object, error) {
	if rt, exists := s.vkToType[vk]; exists {
		new := reflect.New(rt).Interface().(Object)
		new.GetObjectKind().SetVersionKind(vk)
		return new, nil
	}
	return nil, fmt.Errorf("kind %s is not registered", vk)
}

func (s *Scheme) VersionKind(obj Object) (schema.VersionKind, error) {
	rt := ensureStructPointer(obj)

	if vk, ok := s.typeToVK[rt]; ok {
		return vk, nil
	}
	return schema.VersionKind{}, fmt.Errorf("object %T is not registered", obj)
}

func (s *Scheme) ListVersionKind(obj Object) (schema.VersionKind, error) {
	vk, err := s.VersionKind(obj)
	if err != nil {
		return schema.VersionKind{}, err
	}

	listVK := schema.VersionKind{
		Version: vk.Version,
		Kind:    vk.Kind + "List",
	}
	if _, ok := s.vkToType[listVK]; !ok {
		return schema.VersionKind{},
			fmt.Errorf("no list type for %s is not registered", obj)
	}
	return listVK, nil
}

func (s *Scheme) KnownObjectKinds() []schema.VersionKind {
	var vks []schema.VersionKind
	for vk := range s.vkToType {
		if strings.HasSuffix(vk.Kind, "List") {
			continue
		}

		vks = append(vks, vk)
	}
	return vks
}
