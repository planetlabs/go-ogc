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
	"errors"
	"time"
)

type RecordCore struct {
	Id          string
	Type        string
	Title       string
	Description string
	Time        time.Time
	Created     time.Time
	Updated     time.Time
}

var (
	_ Extension = (*RecordCore)(nil)
)

func (r *RecordCore) URI() string {
	return "http://www.opengis.net/spec/ogcapi-records-1/1.0/req/record-core"
}

func (r *RecordCore) Encode(featureMap map[string]any) error {
	propertiesMap, ok := featureMap["properties"].(map[string]any)
	if !ok {
		return errors.New("missing properties")
	}

	// required id
	id, _ := featureMap["id"].(string)
	if id == "" {
		id = r.Id
		if id == "" {
			return errors.New("missing id")
		}
	}
	featureMap["id"] = id

	// required time
	if r.Time.IsZero() {
		featureMap["time"] = nil
	} else {
		featureMap["time"] = r.Time.Format(time.RFC3339Nano)
	}

	// required type
	if r.Type == "" {
		return errors.New("missing type")
	}
	propertiesMap["type"] = r.Type

	// required title
	if r.Title == "" {
		return errors.New("missing title")
	}
	propertiesMap["title"] = r.Title

	// optional description
	if r.Description != "" {
		propertiesMap["description"] = r.Description
	}

	// optional created
	if !r.Created.IsZero() {
		propertiesMap["created"] = r.Created.Format(time.RFC3339Nano)
	}

	// optional updated
	if !r.Updated.IsZero() {
		propertiesMap["updated"] = r.Updated.Format(time.RFC3339Nano)
	}

	return nil
}
