/**
 * Copyright 2023 Planet Labs PBC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package filter_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/planetlabs/go-ogc/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpatial(t *testing.T) {
	cases := []struct {
		filter *filter.Filter
		data   string
	}{
		{
			filter: &filter.Filter{
				Expression: &filter.SpatialComparison{
					Name:  filter.GeometryContains,
					Left:  &filter.BoundingBox{[]float64{-180, -90, 180, 90}},
					Right: &filter.Property{"geometry"},
				},
			},
			data: `{
				"op": "s_contains",
				"args": [
					{"bbox": [-180, -90, 180, 90]},
					{"property": "geometry"}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.SpatialComparison{
					Name:  filter.GeometryCrosses,
					Left:  &filter.Property{"geom1"},
					Right: &filter.Property{"geom2"},
				},
			},
			data: `{
				"op": "s_crosses",
				"args": [
					{"property": "geom1"},
					{"property": "geom2"}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.SpatialComparison{
					Name:  filter.GeometryDisjoint,
					Left:  &filter.Property{"geom1"},
					Right: &filter.Property{"geom2"},
				},
			},
			data: `{
				"op": "s_disjoint",
				"args": [
					{"property": "geom1"},
					{"property": "geom2"}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.SpatialComparison{
					Name:  filter.GeometryEquals,
					Left:  &filter.Property{"geom1"},
					Right: &filter.Property{"geom2"},
				},
			},
			data: `{
				"op": "s_equals",
				"args": [
					{"property": "geom1"},
					{"property": "geom2"}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.SpatialComparison{
					Name:  filter.GeometryIntersects,
					Left:  &filter.Property{"geom1"},
					Right: &filter.Property{"geom2"},
				},
			},
			data: `{
				"op": "s_intersects",
				"args": [
					{"property": "geom1"},
					{"property": "geom2"}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.SpatialComparison{
					Name:  filter.GeometryOverlaps,
					Left:  &filter.Property{"geom1"},
					Right: &filter.Property{"geom2"},
				},
			},
			data: `{
				"op": "s_overlaps",
				"args": [
					{"property": "geom1"},
					{"property": "geom2"}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.SpatialComparison{
					Name:  filter.GeometryWithin,
					Left:  &filter.Property{"geom1"},
					Right: &filter.Property{"geom2"},
				},
			},
			data: `{
				"op": "s_within",
				"args": [
					{"property": "geom1"},
					{"property": "geom2"}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.SpatialComparison{
					Name:  filter.GeometryTouches,
					Left:  &filter.Property{"geom1"},
					Right: &filter.Property{"geom2"},
				},
			},
			data: `{
				"op": "s_touches",
				"args": [
					{"property": "geom1"},
					{"property": "geom2"}
				]
			}`,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			data, err := json.Marshal(c.filter)
			require.NoError(t, err)
			assert.JSONEq(t, c.data, string(data))

			filter := &filter.Filter{}
			require.NoError(t, json.Unmarshal([]byte(c.data), filter))
			assert.Equal(t, c.filter, filter)
		})
	}
}
