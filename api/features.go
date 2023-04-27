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

	"github.com/mitchellh/mapstructure"
)

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
