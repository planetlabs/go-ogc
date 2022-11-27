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
package common

import (
	"net/url"

	"github.com/planetlabs/go-ogc/api"
)

func GetConformance(base *url.URL) (*api.Conformance, error) {
	conformance := &api.Conformance{
		Links: []*api.Link{
			{
				Rel:  "self",
				Type: "application/json",
				Href: getURL(base, "/conformance").String(),
			},
			{
				Title: "Service Root",
				Rel:   "service",
				Type:  "application/json",
				Href:  getURL(base, "/").String(),
			},
		},
		ConformsTo: []string{
			"http://www.opengis.net/spec/ogcapi-common-1/1.0/conf/core",
			"http://www.opengis.net/spec/ogcapi-common-1/1.0/conf/json",
			"http://www.opengis.net/spec/ogcapi-common-1/1.0/req/oas30",
			"http://www.opengis.net/spec/ogcapi-tiles-1/1.0/conf/core",
			"http://www.opengis.net/spec/ogcapi-tiles-1/1.0/conf/tileset",
			"http://www.opengis.net/spec/ogcapi-tiles-1/1.0/conf/tilesets-list",
			"http://www.opengis.net/spec/tms/2.0/conf/tilematrixset",
			"http://www.opengis.net/spec/tms/2.0/conf/tilesetmetadata",
			"http://www.opengis.net/spec/tms/2.0/conf/json-tilematrixset",
			"http://www.opengis.net/spec/tms/2.0/conf/json-tilesetmetadata",
			"http://www.opengis.net/spec/tms/2.0/conf/tilematrixsetlimits",
			"http://www.opengis.net/spec/tms/2.0/conf/json-tilematrixsetlimits",
		},
	}

	return conformance, nil
}
