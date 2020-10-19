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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v1 "github.com/routerd/ipv6d/api/v1"
	"github.com/routerd/ipv6d/internal/runtime"
)

var scheme = runtime.NewScheme()

func init() {
	v1.AddToScheme(scheme)
}

func TestRepository(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		r, err := NewRepository(scheme, &v1.NetworkMap{}, &v1.NetworkMapList{})
		require.NoError(t, err)
		ctx := context.Background()

		// Inject some data
		r.data["test123"] = []byte(`{"metadata": {"name":"test123"}}`)

		// List
		list := &v1.NetworkMapList{}
		err = r.List(ctx, list)
		require.NoError(t, err)

		if assert.Len(t, list.Items, 1) {
			assert.Equal(t, "test123", list.Items[0].Name)
		}
	})
}
