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

package api_test

import (
	"fmt"
	"testing"

	"github.com/planetlabs/go-ogc/api"
	"github.com/planetlabs/go-ogc/util/mercator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTileMatrixSet(t *testing.T) {
	tileMatrixSet, err := api.GetTileMatrixSet(api.TMSWebMercatorQuad)
	require.Nil(t, err)

	assert.Len(t, tileMatrixSet.TileMatrices, 30)
}

func TestGetTileMatrixSetMinZoom(t *testing.T) {
	tileMatrixSet, err := api.GetTileMatrixSet(api.TMSWebMercatorQuad, api.TileMatrixSetMinZoom(10))
	require.Nil(t, err)

	assert.Len(t, tileMatrixSet.TileMatrices, 20)
}

func TestGetTileMatrixSetMaxZoom(t *testing.T) {
	tileMatrixSet, err := api.GetTileMatrixSet(api.TMSWebMercatorQuad, api.TileMatrixSetMaxZoom(10))
	require.Nil(t, err)

	assert.Len(t, tileMatrixSet.TileMatrices, 11)
}

func TestGetTileMatrixSetMinMaxZoom(t *testing.T) {
	tileMatrixSet, err := api.GetTileMatrixSet(
		api.TMSWebMercatorQuad,
		api.TileMatrixSetMinZoom(6),
		api.TileMatrixSetMaxZoom(10),
	)
	require.Nil(t, err)

	assert.Len(t, tileMatrixSet.TileMatrices, 5)
}

func TestTileMatrixLimit(t *testing.T) {
	tileMatrixSet, err := api.GetTileMatrixSet(api.TMSWebMercatorQuad)
	require.Nil(t, err)

	cases := []struct {
		level  int
		bounds []float64
		minCol int
		maxCol int
		minRow int
		maxRow int
	}{
		{
			level:  0,
			minCol: 0,
			maxCol: 0,
			minRow: 0,
			maxRow: 0,
		},
		{
			level:  0,
			bounds: mercator.Forward([]float64{-180, 0, 0, 90}),
			minCol: 0,
			maxCol: 0,
			minRow: 0,
			maxRow: 0,
		},
		{
			level:  1,
			bounds: mercator.Forward([]float64{-180, 0, 0, 90}),
			minCol: 0,
			maxCol: 0,
			minRow: 0,
			maxRow: 0,
		},
		{
			level:  1,
			bounds: mercator.Forward([]float64{0, -90, 180, 0}),
			minCol: 1,
			maxCol: 1,
			minRow: 1,
			maxRow: 1,
		},
		{
			level:  2,
			bounds: mercator.Forward([]float64{0, -90, 180, 0}),
			minCol: 2,
			maxCol: 3,
			minRow: 2,
			maxRow: 3,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			matrix := tileMatrixSet.TileMatrices[c.level]
			var boundingBox *api.BoundingBox
			if c.bounds != nil {
				boundingBox = &api.BoundingBox{
					LowerLeft:  c.bounds[:2],
					UpperRight: c.bounds[2:],
				}
			}

			limit, err := matrix.Limit(boundingBox)
			require.Nil(t, err)
			assert.Equal(t, matrix.ID, limit.TileMatrix)
			assert.Equal(t, c.minCol, limit.MinTileCol)
			assert.Equal(t, c.maxCol, limit.MaxTileCol)
			assert.Equal(t, c.minRow, limit.MinTileRow)
			assert.Equal(t, c.maxRow, limit.MaxTileRow)
		})
	}
}
