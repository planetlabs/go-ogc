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

package serve

import (
	"net/http"
	"net/url"
	"strings"
)

// getProxyBase constructs a base URL for the proxied request.
// For example a request to http://foo.example.com/service/api with
// headers applied from a proxy would result in a URL like http://api.example.com/service.
// The path from this base URL can be joined with an internal path when serializing links.
func getProxyBase(req *http.Request) *url.URL {
	scheme := req.URL.Scheme
	if header := req.Header.Get("X-Forwarded-Proto"); header != "" {
		scheme = header
	}
	if scheme == "" {
		scheme = "http"
	}

	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	if header := req.Header.Get("X-Forwarded-Host"); header != "" {
		host = header
	}

	path := "/"
	if header := req.Header.Get("X-Original-Path"); header != "" {
		header = strings.TrimSuffix(header, "/")
		originPath := strings.TrimSuffix(req.URL.Path, "/")
		path = strings.TrimSuffix(header, originPath)
	}

	return &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}
}
