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
	"fmt"
	"testing"

	"github.com/planetlabs/go-ogc/filter"
)

func TestCharacter(t *testing.T) {
	cases := []*FilterCase{
		{
			filter: &filter.Filter{
				Expression: &filter.Comparison{
					Name:  filter.Equals,
					Left:  &filter.String{"foo"},
					Right: &filter.String{"bar"},
				},
			},
			data: `{
				"op": "=",
				"args": ["foo", "bar"]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.Comparison{
					Name:  filter.Equals,
					Left:  &filter.Property{"soup"},
					Right: &filter.CaseInsensitive{&filter.String{"chicken"}},
				},
			},
			data: `{
				"op": "=",
				"args": [{"property": "soup"}, {"op": "casei", "args": ["chicken"]}]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.Comparison{
					Name:  filter.Equals,
					Left:  &filter.Property{"soup"},
					Right: &filter.AccentInsensitive{&filter.String{"Chícken"}},
				},
			},
			data: `{
				"op": "=",
				"args": [{"property": "soup"}, {"op": "accenti", "args": ["Chícken"]}]
			}`,
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assertFilterIO(t, c)
		})
	}
}
