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

func TestProperty(t *testing.T) {
	cases := []struct {
		filter *filter.Filter
		data   string
	}{
		{
			filter: &filter.Filter{
				Expression: &filter.Comparison{
					Name:  filter.Equals,
					Left:  &filter.Property{"left"},
					Right: &filter.Property{"right"},
				},
			},
			data: `{
				"op": "=",
				"args": [{"property": "left"}, {"property": "right"}]
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
