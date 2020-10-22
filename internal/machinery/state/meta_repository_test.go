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

	"github.com/stretchr/testify/require"

	"routerd.net/ipv6d/internal/testutil"
)

func TestMetaRepository(t *testing.T) {
	t.Run("LoadFromFileSystem", func(t *testing.T) {
		log := testutil.NewLogger(t)
		r, err := NewMetaRepository(log, testScheme)
		require.NoError(t, err)

		// Run
		stopCh := make(chan struct{})
		go r.Run(stopCh)
		defer close(stopCh)

		err = r.LoadFromFileSystem("./testdata")
		require.NoError(t, err)

		ctx := context.Background()
		t1 := &testObject{}
		require.NoError(t, r.Get(ctx, "test1", t1))

		t2 := &testObject{}
		require.NoError(t, r.Get(ctx, "test2", t2))

		t3 := &testObject{}
		require.NoError(t, r.Get(ctx, "test3", t3))
	})
}
