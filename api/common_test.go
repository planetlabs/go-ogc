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

package api_test

import (
	"encoding/json"
	"testing"

	"github.com/planetlabs/go-ogc/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkMarshal(t *testing.T) {
	cases := []struct {
		name     string
		link     *api.Link
		expected string
	}{
		{
			name: "basic",
			link: &api.Link{
				Href: "http://example.com/resource.json",
				Rel:  "self",
				Type: "application/json",
			},
			expected: `{
				"href": "http://example.com/resource.json",
				"rel": "self",
				"type": "application/json"
			}`,
		},
		{
			name: "additional fields",
			link: &api.Link{
				Href: "http://example.com/resource.json",
				Rel:  "self",
				Type: "application/json",
				AdditionalFields: map[string]interface{}{
					"one": "foo",
					"two": "bar",
				},
			},
			expected: `{
				"href": "http://example.com/resource.json",
				"rel": "self",
				"type": "application/json",
				"one": "foo",
				"two": "bar"
			}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.link)
			require.NoError(t, err)
			assert.JSONEq(t, tc.expected, string(data))
		})
	}
}
