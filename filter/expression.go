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

package filter

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Expression interface {
	expression()
	String() string
}

type ScalarExpression interface {
	Expression
	scalarExpression()
}

func unmarshalExpression(data []byte) (Expression, error) {
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, err
	}
	return decodeExpression(value)
}

func decodeExpression(value any) (Expression, error) {
	switch v := value.(type) {
	case bool:
		return &Boolean{Value: v}, nil
	case string:
		return &String{Value: v}, nil
	case float64:
		return &Number{Value: v}, nil
	case []any:
		return decodeArray(v)
	case map[string]any:
		if dateString, ok := v["date"].(string); ok {
			return decodeDate(dateString)
		}

		if timestampString, ok := v["timestamp"].(string); ok {
			return decodeTimestamp(timestampString)
		}

		if intervalValues, ok := v["interval"].([]any); ok {
			return decodeInterval(intervalValues)
		}

		if bbox, ok := v["bbox"].([]any); ok {
			return decodeBoundingBox(bbox)
		}

		if t, ok := v["type"].(string); ok {
			if geometryTypes[t] {
				return decodeGeometry(v)
			}
			return nil, fmt.Errorf("unexpected expression type: %s", t)
		}

		if propertyName, ok := v["property"].(string); ok {
			return &Property{Name: propertyName}, nil
		}

		if opName, ok := v["op"].(string); ok {
			args, ok := v["args"].([]any)
			if !ok {
				return nil, fmt.Errorf("expected args in %q op", opName)
			}
			return decodeOp(opName, args)
		}
	}

	return nil, errors.New("unsupported expression")
}
