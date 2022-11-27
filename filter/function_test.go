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

func TestFunction(t *testing.T) {
	cases := []struct {
		filter *filter.Filter
		data   string
	}{
		{
			filter: &filter.Filter{
				Expression: &filter.Comparison{
					Name: filter.Equals,
					Left: &filter.Function{
						Name: "testing",
						Args: []filter.Expression{
							&filter.Number{1},
							&filter.Number{2},
							&filter.Number{3},
						},
					},
					Right: &filter.Boolean{true},
				},
			},
			data: `{
				"op": "=",
				"args": [{"function": {"name": "testing", "args": [1, 2, 3]}}, true]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.Comparison{
					Name:  filter.Equals,
					Left:  &filter.Function{Name: "agreeable"},
					Right: &filter.Boolean{false},
				},
			},
			data: `{
				"op": "=",
				"args": [{"function": {"name": "agreeable"}}, false]
			}`,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			data, err := json.Marshal(c.filter)
			require.Nil(t, err)
			assert.JSONEq(t, c.data, string(data))

			filter := &filter.Filter{}
			require.Nil(t, json.Unmarshal([]byte(c.data), filter))
			assert.Equal(t, c.filter, filter)
		})
	}
}
