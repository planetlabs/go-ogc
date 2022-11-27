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

func TestArrayComparison(t *testing.T) {
	cases := []struct {
		filter *filter.Filter
		data   string
	}{
		{
			filter: &filter.Filter{
				Expression: &filter.ArrayComparison{
					Name:  filter.ArrayContainedBy,
					Left:  &filter.Property{"array1"},
					Right: &filter.Property{"array2"},
				},
			},
			data: `{
				"op": "a_containedBy",
				"args": [{"property": "array1"}, {"property": "array2"}]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.ArrayComparison{
					Name: filter.ArrayContains,
					Left: &filter.Property{"array"},
					Right: filter.Array{
						&filter.String{"foo"},
						&filter.Number{42},
						&filter.Boolean{false},
					},
				},
			},
			data: `{
				"op": "a_contains",
				"args": [{"property": "array"}, ["foo", 42, false]]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.ArrayComparison{
					Name: filter.ArrayEquals,
					Left: &filter.Property{"array"},
					Right: filter.Array{
						&filter.BoundingBox{[]float64{4, 3, 2, 1}},
					},
				},
			},
			data: `{
				"op": "a_equals",
				"args": [{"property": "array"}, [{"bbox": [4, 3, 2, 1]}]]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.ArrayComparison{
					Name: filter.ArrayOverlaps,
					Left: &filter.Property{"food"},
					Right: filter.Array{
						&filter.String{"raddish"},
						&filter.String{"turnip"},
					},
				},
			},
			data: `{
				"op": "a_overlaps",
				"args": [{"property": "food"}, ["raddish", "turnip"]]
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
