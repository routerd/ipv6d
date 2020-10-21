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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-logr/logr"
	"sigs.k8s.io/yaml"

	"routerd.net/ipv6d/internal/machinery/runtime"
)

type MetaRepository struct {
	scheme   *runtime.Scheme
	vkToRepo map[runtime.VersionKind]*Repository
}

func NewMetaRepository(scheme *runtime.Scheme) (*MetaRepository, error) {
	mr := &MetaRepository{
		scheme:   scheme,
		vkToRepo: map[runtime.VersionKind]*Repository{},
	}

	vks := scheme.KnownObjectKinds()
	for _, vk := range vks {
		obj, err := scheme.New(vk)
		if err != nil {
			return nil, err
		}

		vkList, err := scheme.ListVersionKind(obj)
		if err != nil {
			return nil, err
		}

		listObj, err := scheme.New(vkList)
		if err != nil {
			return nil, err
		}

		repo, err := NewRepository(scheme, obj.(Object), listObj)
		if err != nil {
			return nil, err
		}

		mr.vkToRepo[vk] = repo
		mr.vkToRepo[vkList] = repo
	}
	return mr, nil
}

func (r *MetaRepository) LoadFromFileSystem(
	log logr.Logger, folder string,
) error {
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// unmarshal
		fileBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		documents := bytes.Split(fileBytes, []byte("\n---"))
		for _, document := range documents {
			vk := runtime.VersionKind{}
			if err := yaml.Unmarshal(document, &vk); err != nil {
				return err
			}

			obj, err := r.scheme.New(vk)
			if err != nil {
				return fmt.Errorf("importing file %s: %w", path, err)
			}
			if err := yaml.Unmarshal(document, obj); err != nil {
				return err
			}

			// store
			log.Info(fmt.Sprintf(
				"imported %s %s", obj.GetObjectKind(), obj.(Object).GetName()))
			if err := r.Create(context.TODO(), obj.(Object)); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *MetaRepository) Run(stopCh <-chan struct{}) {
	for _, repo := range r.vkToRepo {
		go repo.Run(stopCh)
	}
	<-stopCh
}

func (r *MetaRepository) repoForObj(obj runtime.Object) (*Repository, error) {
	vk, err := r.scheme.VersionKind(obj)
	if err != nil {
		return nil, err
	}
	repo, ok := r.vkToRepo[vk]
	if !ok {
		return nil, fmt.Errorf("no repository registered for type %T", obj)
	}
	return repo, nil
}

func (r *MetaRepository) Get(ctx context.Context, key string, obj Object) error {
	repo, err := r.repoForObj(obj)
	if err != nil {
		return err
	}
	return repo.Get(ctx, key, obj)
}

func (r *MetaRepository) List(ctx context.Context, listObj ObjectList) error {
	repo, err := r.repoForObj(listObj)
	if err != nil {
		return err
	}
	return repo.List(ctx, listObj)
}

func (r *MetaRepository) Watch(ctx context.Context, obj Object) (Watcher, error) {
	repo, err := r.repoForObj(obj)
	if err != nil {
		return nil, err
	}
	return repo.Watch(ctx, obj)
}

func (r *MetaRepository) Create(ctx context.Context, obj Object) error {
	repo, err := r.repoForObj(obj)
	if err != nil {
		return err
	}
	return repo.Create(ctx, obj)
}

func (r *MetaRepository) Delete(ctx context.Context, obj Object) error {
	repo, err := r.repoForObj(obj)
	if err != nil {
		return err
	}
	return repo.Delete(ctx, obj)
}

func (r *MetaRepository) Update(ctx context.Context, obj Object) error {
	repo, err := r.repoForObj(obj)
	if err != nil {
		return err
	}
	return repo.Update(ctx, obj)
}

func (r *MetaRepository) UpdateStatus(ctx context.Context, obj Object) error {
	repo, err := r.repoForObj(obj)
	if err != nil {
		return err
	}
	return repo.UpdateStatus(ctx, obj)
}
