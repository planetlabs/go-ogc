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

package filter_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/planetlabs/go-ogc/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComparison(t *testing.T) {
	cases := []struct {
		filter *filter.Filter
		data   string
	}{
		{
			filter: &filter.Filter{
				Expression: &filter.Comparison{
					Name:  filter.Equals,
					Left:  &filter.Property{"city"},
					Right: &filter.String{"Pleasantville"},
				},
			},
			data: `{
				"op": "=",
				"args": [{"property": "city"}, "Pleasantville"]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.Comparison{
					Name:  filter.LessThan,
					Left:  &filter.Property{"population"},
					Right: &filter.Number{123456},
				},
			},
			data: `{
				"op": "<",
				"args": [{"property": "population"}, 123456]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.Comparison{
					Name:  filter.Equals,
					Left:  &filter.Property{"sunny"},
					Right: &filter.Boolean{true},
				},
			},
			data: `{
				"op": "=",
				"args": [{"property": "sunny"}, true]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.Comparison{
					Name:  filter.GreaterThan,
					Left:  &filter.Property{"income"},
					Right: &filter.Number{1e5},
				},
			},
			data: `{
				"op": ">",
				"args": [{"property": "income"}, 1e5]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.Comparison{
					Name:  filter.NotEquals,
					Left:  &filter.Property{"news"},
					Right: &filter.Boolean{false},
				},
			},
			data: `{
				"op": "<>",
				"args": [{"property": "news"}, false]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.Comparison{
					Name:  filter.GreaterThanOrEquals,
					Left:  &filter.Property{"recreation"},
					Right: &filter.Property{"work"},
				},
			},
			data: `{
				"op": ">=",
				"args": [{"property": "recreation"}, {"property": "work"}]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.Comparison{
					Name:  filter.LessThanOrEquals,
					Left:  &filter.Number{8},
					Right: &filter.Property{"rest"},
				},
			},
			data: `{
				"op": "<=",
				"args": [8, {"property": "rest"}]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.Like{
					String:  &filter.Property{"name"},
					Pattern: &filter.CaseInsensitive{&filter.String{"park"}},
				},
			},
			data: `{
				"op": "like",
				"args": [{"property": "name"}, {"op": "casei", "args": ["park"]}]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.Like{
					String:  &filter.Property{"name"},
					Pattern: &filter.AccentInsensitive{&filter.String{"Noël"}},
				},
			},
			data: `{
				"op": "like",
				"args": [{"property": "name"}, {"op": "accenti", "args": ["Noël"]}]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.Between{
					Value: &filter.Property{"depth"},
					Low:   &filter.Number{50},
					High:  &filter.Number{100},
				},
			},
			data: `{
				"op": "between",
				"args": [{"property": "depth"}, 50, 100]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.In{
					Item: &filter.Property{"item"},
					List: filter.ScalarList{
						&filter.String{"one"},
						&filter.String{"two"},
					},
				},
			},
			data: `{
				"op": "in",
				"args": [{"property": "item"}, ["one", "two"]]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.IsNull{&filter.Property{"geometry"}},
			},
			data: `{
				"op": "isNull",
				"args": [{"property": "geometry"}]
			}`,
		},
	}

	schema := getSchema(t)
	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			data, err := json.Marshal(c.filter)
			require.NoError(t, err)
			assert.JSONEq(t, c.data, string(data))

			v := map[string]any{}
			require.NoError(t, json.Unmarshal(data, &v))
			if err := schema.Validate(v); err != nil {
				t.Errorf("failed to validate\n%#v", err)
			}

			filter := &filter.Filter{}
			require.NoError(t, json.Unmarshal([]byte(c.data), filter))
			assert.Equal(t, c.filter, filter)
		})
	}
}
