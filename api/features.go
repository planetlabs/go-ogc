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

package api

import (
	"encoding/json"
	"fmt"

	"github.com/go-viper/mapstructure/v2"
)

type Collection struct {
	Id          string      `json:"id"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Links       []*Link     `json:"links"`
	Extent      *Extent     `json:"extent,omitempty"`
	ItemType    string      `json:"itemType,omitempty"`
	Crs         []string    `json:"crs,omitempty"`
	Extensions  []Extension `json:"-"`
}

var (
	_ json.Marshaler   = (*Collection)(nil)
	_ json.Unmarshaler = (*Collection)(nil)
)

func (collection Collection) MarshalJSON() ([]byte, error) {
	collectionMap := map[string]any{}
	decoder, decoderErr := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  &collectionMap,
	})
	if decoderErr != nil {
		return nil, decoderErr
	}

	decodeErr := decoder.Decode(collection)
	if decodeErr != nil {
		return nil, decodeErr
	}

	for _, extension := range collection.Extensions {
		if err := extension.Encode(collectionMap); err != nil {
			return nil, fmt.Errorf("trouble encoding feature JSON with the %q extension: %w", extension.URI(), err)
		}
	}

	return json.Marshal(collectionMap)
}

type decodedCollection Collection

func (collection *Collection) UnmarshalJSON(data []byte) error {
	d := &decodedCollection{Extensions: collection.Extensions}
	if err := json.Unmarshal(data, d); err != nil {
		return err
	}
	*collection = Collection(*d)

	for _, e := range collection.Extensions {
		if err := e.Decode(data); err != nil {
			return err
		}
	}
	return nil
}

type Extent struct {
	Spatial  *SpatialExtent  `json:"spatial,omitempty"`
	Temporal *TemporalExtent `json:"temporal,omitempty"`
}

type SpatialExtent struct {
	Bbox [][]float64 `json:"bbox"`
	Crs  string      `json:"crs,omitempty"`
}

type TemporalExtent struct {
	Interval [][]any `json:"interval"`
	Trs      string  `json:"trs,omitempty"`
}

type CollectionsList struct {
	Collections []*Collection `json:"collections"`
	Links       []*Link       `json:"links"`
}

type Feature struct {
	Id         string         `json:"id,omitempty"`
	Geometry   any            `json:"geometry"`
	Properties map[string]any `json:"properties"`
	Links      []*Link        `json:"links,omitempty"`
	Extensions []Extension    `json:"-"`
}

var _ json.Marshaler = (*Feature)(nil)

func (feature Feature) MarshalJSON() ([]byte, error) {
	featureMap := map[string]any{"type": "Feature"}
	decoder, decoderErr := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  &featureMap,
	})
	if decoderErr != nil {
		return nil, decoderErr
	}

	decodeErr := decoder.Decode(feature)
	if decodeErr != nil {
		return nil, decodeErr
	}

	extensionUris := []string{}
	lookup := map[string]bool{}

	for _, extension := range feature.Extensions {
		if err := extension.Encode(featureMap); err != nil {
			return nil, fmt.Errorf("trouble encoding feature JSON with the %q extension: %w", extension.URI(), err)
		}
		uri := extension.URI()
		if !lookup[uri] {
			extensionUris = append(extensionUris, uri)
			lookup[uri] = true
		}
	}

	if len(extensionUris) > 0 {
		featureMap["conformsTo"] = extensionUris
	}

	return json.Marshal(featureMap)
}

type FeatureCollection struct {
	Type           string     `json:"type"`
	Features       []*Feature `json:"features"`
	Links          []*Link    `json:"links,omitempty"`
	TimeStamp      string     `json:"timeStamp,omitempty"`
	NumberMatched  int        `json:"numberMatched,omitempty"`
	NumberReturned int        `json:"numberReturned,omitempty"`
}
