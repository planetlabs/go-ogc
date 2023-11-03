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
	"errors"
	"fmt"
	"testing"

	"github.com/planetlabs/go-ogc/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectionMarshal(t *testing.T) {
	cases := []struct {
		name       string
		collection *api.Collection
		expected   string
	}{
		{
			name: "minimal",
			collection: &api.Collection{
				Id:          "test-1",
				Title:       "Test Collection",
				Description: "Test collection description.",
				Extent: &api.Extent{
					Spatial: &api.SpatialExtent{
						Bbox: [][]float64{{1, 2, 3, 4}},
					},
				},
				Links: []*api.Link{},
			},
			expected: `{
				"id": "test-1",
				"title": "Test Collection",
				"description": "Test collection description.",
				"extent": {
					"spatial": {
						"bbox": [[1, 2, 3, 4]]
					}
				},
				"links": []
			}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.collection)
			require.NoError(t, err)
			assert.JSONEq(t, tc.expected, string(actual))
		})
	}

}

func TestFeatureMarshal(t *testing.T) {
	cases := []struct {
		name     string
		feature  *api.Feature
		expected string
	}{
		{
			name: "basic",
			feature: &api.Feature{
				Id:       "foo",
				Geometry: nil,
				Properties: map[string]interface{}{
					"one": "foo",
					"two": "bar",
				},
				Links: []*api.Link{
					{
						Href: "http://example.com/resource.json",
						Rel:  "self",
						Type: "application/geo+json",
					},
				},
			},
			expected: `{
				"type": "Feature",
				"id": "foo",
				"geometry": null,
				"properties": {
					"one": "foo",
					"two": "bar"
				},
				"links": [
					{
						"href": "http://example.com/resource.json",
						"rel": "self",
						"type": "application/geo+json"
					}
				]
			}`,
		},
		{
			name: "minimal",
			feature: &api.Feature{
				Geometry: nil,
				Properties: map[string]interface{}{
					"one": "foo",
					"two": "bar",
				},
			},
			expected: `{
				"type": "Feature",
				"geometry": null,
				"properties": {
					"one": "foo",
					"two": "bar"
				}
			}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := json.Marshal(tc.feature)
			require.NoError(t, err)
			assert.JSONEq(t, tc.expected, string(actual))
		})
	}
}

var (
	_ api.Extension = (*FeatureExtension)(nil)
	_ api.Extension = (*CollectionExtension)(nil)
)

type FeatureExtension struct {
	RootFoo   string
	NestedBar string
}

func (e *FeatureExtension) Encode(featureMap map[string]any) error {
	featureMap["test:foo"] = e.RootFoo
	propertiesMap, ok := featureMap["properties"].(map[string]any)
	if !ok {
		return fmt.Errorf("expected properties on a feature ")
	}
	propertiesMap["test:bar"] = e.NestedBar
	return nil
}

func (e *FeatureExtension) Decode(data []byte) error {
	return errors.New("not implemented")
}

func (e *FeatureExtension) URI() string {
	return "https://example.com/test-extension"
}

type CollectionExtension struct {
	Classes []string
}

func (e *CollectionExtension) Encode(collectionMap map[string]any) error {
	collectionMap["classes"] = e.Classes
	return nil
}

func (e *CollectionExtension) Decode(data []byte) error {
	return json.Unmarshal(data, e)
}

func (e *CollectionExtension) URI() string {
	return "https://example.com/test-extension"
}

func TestFeatureMarshalExtension(t *testing.T) {
	feature := &api.Feature{
		Geometry: nil,
		Properties: map[string]interface{}{
			"one": "core-property",
		},
		Extensions: []api.Extension{
			&FeatureExtension{
				RootFoo:   "foo-value",
				NestedBar: "bar-value",
			},
		},
	}

	expected := `{
		"type": "Feature",
		"geometry": null,
		"properties": {
			"one": "core-property",
			"test:bar": "bar-value"
		},
		"test:foo": "foo-value",
		"conformsTo": [
			"https://example.com/test-extension"
		]
	}`

	actual, err := json.Marshal(feature)
	require.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
}

func TestCollectionMarshalExtension(t *testing.T) {
	collection := &api.Collection{
		Id: "test-collection",
		Extensions: []api.Extension{
			&CollectionExtension{
				Classes: []string{"foo", "bar"},
			},
		},
		Links: []*api.Link{},
	}

	expected := `{
		"id": "test-collection",
		"classes": ["foo", "bar"],
		"links": []
	}`

	actual, err := json.Marshal(collection)
	require.NoError(t, err)
	assert.JSONEq(t, expected, string(actual))
}

func TestCollectionUnmarshalExtension(t *testing.T) {
	data := `{
		"id": "test-collection",
		"classes": ["foo", "bar"],
		"links": []
	}`

	extension := &CollectionExtension{}
	collection := &api.Collection{
		Extensions: []api.Extension{extension},
	}

	err := json.Unmarshal([]byte(data), collection)
	require.NoError(t, err)
	assert.Equal(t, []string{"foo", "bar"}, extension.Classes)
}
