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
	"time"

	"github.com/planetlabs/go-ogc/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordCore(t *testing.T) {
	origin := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	created := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	updated := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name     string
		record   *api.Feature
		expected string
	}{
		{
			name: "basic",
			record: &api.Feature{
				Id:       "test-id",
				Geometry: nil,
				Properties: map[string]interface{}{
					"one": "foo",
				},
				Extensions: []api.Extension{
					&api.RecordCore{
						Type:  "test-record",
						Title: "Test Record",
					},
				},
			},
			expected: `{
				"type": "Feature",
				"id": "test-id",
				"geometry": null,
				"time": null,
				"properties": {
					"one": "foo",
					"title": "Test Record",
					"type": "test-record"
				},
				"conformsTo": [
					"http://www.opengis.net/spec/ogcapi-records-1/1.0/req/record-core"
				]
			}`,
		},
		{
			name: "less basic",
			record: &api.Feature{
				Geometry: nil,
				Properties: map[string]interface{}{
					"one": "foo",
				},
				Extensions: []api.Extension{
					&api.RecordCore{
						Id:          "record-id",
						Type:        "test-record",
						Title:       "Test Record",
						Description: "Description of the record",
						Time:        origin,
						Created:     created,
						Updated:     updated,
					},
				},
			},
			expected: `{
				"type": "Feature",
				"id": "record-id",
				"geometry": null,
				"time": "2000-01-01T00:00:00Z",
				"properties": {
					"type": "test-record",
					"title": "Test Record",
					"description": "Description of the record",
					"one": "foo",
					"created": "2020-01-01T00:00:00Z",
					"updated": "2020-01-02T00:00:00Z"
				},
				"conformsTo": [
					"http://www.opengis.net/spec/ogcapi-records-1/1.0/req/record-core"
				]
			}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.record)
			require.NoError(t, err)
			assert.JSONEq(t, tc.expected, string(actual))
		})
	}
}
