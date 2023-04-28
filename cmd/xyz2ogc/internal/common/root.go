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

package common

import (
	"net/url"

	"github.com/planetlabs/go-ogc/api"
)

func GetRoot(base *url.URL) (*api.Root, error) {

	root := &api.Root{
		Links: []*api.Link{
			{
				Title: "Service Root",
				Rel:   "self",
				Type:  "application/json",
				Href:  getURL(base, "/").String(),
			},
			{
				Title: "OpenAPI Service Description",
				Rel:   "service-desc",
				Type:  "application/vnd.oai.openapi+json;version=3.0",
				Href:  getURL(base, "/api").String(),
			},
			{
				Title: "OpenAPI Service Docs",
				Rel:   "service-doc",
				Type:  "text/html",
				Href:  getURL(base, "/docs/index.html").String(),
			},
			{
				Title: "OGC API Conformance",
				Rel:   "http://www.opengis.net/def/rel/ogc/1.0/conformance",
				Type:  "application/json",
				Href:  getURL(base, "/conformance").String(),
			},
			{
				Title: "Tileset List",
				Rel:   "collection",
				Type:  "application/json",
				Href:  getURL(base, "/tiles").String(),
			},
			{
				Title: "Tiling Schemes",
				Rel:   "http://www.opengis.net/def/rel/ogc/1.0/tiling-schemes",
				Type:  "application/json",
				Href:  getURL(base, "/tileMatrixSets").String(),
			},
		},
	}

	return root, nil
}
