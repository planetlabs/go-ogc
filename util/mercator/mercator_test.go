// Copyright 2023 Planet Labs PBC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mercator_test

import (
	"fmt"
	"testing"

	"github.com/planetlabs/go-ogc/util/mercator"
	"github.com/stretchr/testify/assert"
)

const delta = 1e-8

func TestForward(t *testing.T) {
	cases := []struct {
		gg []float64
		sm []float64
	}{
		{
			gg: []float64{0, 0},
			sm: []float64{0, 0},
		},
		{
			gg: []float64{-180, 0},
			sm: []float64{-mercator.Edge, 0},
		},
		{
			gg: []float64{180, 0},
			sm: []float64{mercator.Edge, 0},
		},
		{
			gg: []float64{0, -90},
			sm: []float64{0, -mercator.Edge},
		},
		{
			gg: []float64{0, 90},
			sm: []float64{0, mercator.Edge},
		},
		{
			gg: []float64{0, 91},
			sm: []float64{0, mercator.Edge},
		},
		{
			gg: []float64{0, -91},
			sm: []float64{0, -mercator.Edge},
		},
		{
			gg: []float64{-180, -90, 180, 90},
			sm: []float64{-mercator.Edge, -mercator.Edge, mercator.Edge, mercator.Edge},
		},
		{
			gg: []float64{-5.625, 52.4827802220782},
			sm: []float64{-626172.13571216376, 6887893.4928337997},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			sm := mercator.Forward(c.gg)
			assert.InDeltaSlice(t, c.sm, sm, delta)
		})
	}
}

func TestInverse(t *testing.T) {
	cases := []struct {
		gg []float64
		sm []float64
	}{
		{
			gg: []float64{0, 0},
			sm: []float64{0, 0},
		},
		{
			gg: []float64{-180, 0},
			sm: []float64{-mercator.Edge, 0},
		},
		{
			gg: []float64{180, 0},
			sm: []float64{mercator.Edge, 0},
		},
		{
			gg: []float64{-5.625, 52.4827802220782},
			sm: []float64{-626172.13571216376, 6887893.4928337997},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			gg := mercator.Inverse(c.sm)
			assert.InDeltaSlice(t, c.gg, gg, delta)
		})
	}
}
