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

package api

import (
	"encoding/json"

	"github.com/go-viper/mapstructure/v2"
)

// Structs for the Common spec.
// http://docs.ogc.org/DRAFTS/19-072.html

// Link is used to connect resources in an OGC API.
//
// Implementations that conform to the Common Core Conformance class (http://www.opengis.net/spec/ogcapi-common-1/1.0/conf/core)
// must link to associated resources with the fields below.
//
// The AdditionalFields map will be populated with any additional fields when unmarshalling JSON.
// You can use this map to add additional fields when marshalling JSON.
type Link struct {
	Href             string         `mapstructure:"href"`
	Rel              string         `mapstructure:"rel"`
	Type             string         `mapstructure:"type,omitempty"`
	HrefLang         string         `mapstructure:"hreflang,omitempty"`
	Title            string         `mapstructure:"title,omitempty"`
	Length           int            `mapstructure:"length,omitempty"`
	Templated        bool           `mapstructure:"templated,omitempty"`
	VarBase          string         `mapstructure:"varBase,omitempty"`
	AdditionalFields map[string]any `mapstructure:",remain"`
}

var (
	_ json.Marshaler   = (*Link)(nil)
	_ json.Unmarshaler = (*Link)(nil)
)

func (link *Link) MarshalJSON() ([]byte, error) {
	m := map[string]any{}
	if err := mapstructure.Decode(link, &m); err != nil {
		return nil, err
	}
	return json.Marshal(m)
}

func (link *Link) UnmarshalJSON(data []byte) error {
	m := map[string]any{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	return mapstructure.Decode(m, link)
}

type Conformance struct {
	Links      []*Link  `json:"links"`
	ConformsTo []string `json:"conformsTo"`
}

type Root struct {
	Links       []*Link `json:"links"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Attribution string  `json:"attribution,omitempty"`
}

type Extension interface {
	URI() string
	Encode(map[string]any) error
	Decode([]byte) error
}
