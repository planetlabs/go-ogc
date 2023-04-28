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

// Structs for the Common spec.
// http://docs.ogc.org/DRAFTS/19-072.html

// Link is used to connect resources in an OGC API.
//
// Implementations that conform to the Common Core Conformance class (http://www.opengis.net/spec/ogcapi-common-1/1.0/conf/core)
// must link to associated resources with the fields below.
type Link struct {
	Href      string `json:"href"`
	Rel       string `json:"rel"`
	Type      string `json:"type,omitempty"`
	HrefLang  string `json:"hreflang,omitempty"`
	Title     string `json:"title,omitempty"`
	Length    int    `json:"length,omitempty"`
	Templated bool   `json:"templated,omitempty"`
	VarBase   string `json:"varBase,omitempty"`
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
