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

package schema

type VersionKind struct {
	Version, Kind string
}

func (vk VersionKind) Empty() bool {
	return len(vk.Version) == 0 && len(vk.Kind) == 0
}

func (vk VersionKind) String() string {
	return vk.Version + "/" + vk.Kind
}

type ObjectKind interface {
	SetVersionKind(vk VersionKind)
	GetVersionKind() VersionKind
}
