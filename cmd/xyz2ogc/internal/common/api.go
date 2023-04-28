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
	"embed"

	"github.com/getkin/kin-openapi/openapi3"
)

//go:embed schema
var embedDir embed.FS

func GetSchema() (*openapi3.T, error) {
	schemaData, readErr := embedDir.ReadFile("schema/api.json")
	if readErr != nil {
		return nil, readErr
	}
	loader := openapi3.NewLoader()
	return loader.LoadFromData(schemaData)
}
