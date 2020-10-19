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

package daemon

// type ConfigLoader struct {
// 	// File system path to config folder.
// 	path string

// 	create map[versionKind]objectCreator
// }

// type objectCreator func() Object

// type versionKind struct {
// 	Version, Kind string
// }

// func (c *ConfigLoader) Load() error {
// 	if err := filepath.Walk(c.path, c.walk); err != nil {
// 		return fmt.Errorf("walking config folder %s: %w", c.path, err)
// 	}

// 	return nil
// }

// func (c *ConfigLoader) walk(path string, info os.FileInfo, err error) error {
// 	fileByte, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		return fmt.Errorf("reading file %s: %w", path, err)
// 	}

// 	obj := &object{}
// 	if err := yaml.Unmarshal(fileByte, obj); err != nil {
// 		return fmt.Errorf("unmarshal yaml from file %s: %w", path, err)
// 	}

// 	return nil
// }

// type object struct {
// 	v1.TypeMeta
// 	v1.ObjectMeta
// }
