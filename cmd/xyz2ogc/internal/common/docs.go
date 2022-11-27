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
	"bytes"
	"html/template"
)

// Docs is the content for the docs page.
var docs = `<!DOCTYPE html>
<html>
  <head>
    <title>XYZ Tiles</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>
      body {
        margin: 0;
        padding: 0;
      }
    </style>
  </head>
  <body>
    <redoc spec-url="{{ .ApiPath }}" hide-download-button expand-responses="200,201" required-props-first></redoc>
    <script src="https://cdn.jsdelivr.net/npm/redoc@latest/bundles/redoc.standalone.js"></script>
  </body>
</html>
`

type docsData struct {
	ApiPath string
}

func GetDocs(apiPath string) ([]byte, error) {
	t, err := template.New("index.html").Parse(docs)
	if err != nil {
		return nil, err
	}
	buffer := &bytes.Buffer{}
	data := &docsData{
		ApiPath: apiPath,
	}
	if err := t.Execute(buffer, data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
