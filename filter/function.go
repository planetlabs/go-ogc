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

type Function struct {
	Name string
	Args []Expression
}

var (
	_ Expression          = (*Function)(nil)
	_ BooleanExpression   = (*Function)(nil) // see https://github.com/opengeospatial/ogcapi-features/issues/796
	_ CharacterExpression = (*Function)(nil)
	_ NumericExpression   = (*Function)(nil)
	_ ArrayExpression     = (*Function)(nil)
	_ SpatialExpression   = (*Function)(nil)
	_ TemporalExpression  = (*Function)(nil)
	_ json.Marshaler      = (*Function)(nil)
)

func (*Function) expression()          {}
func (*Function) scalarExpression()    {}
func (*Function) booleanExpression()   {}
func (*Function) characterExpression() {}
func (*Function) numericExpression()   {}
func (*Function) arrayExpression()     {}
func (*Function) spatialExpression()   {}
func (*Function) temporalExpression()  {}

func (e *Function) MarshalJSON() ([]byte, error) {
	f := map[string]any{
		"name": e.Name,
	}
	if len(e.Args) > 0 {
		f["args"] = e.Args
	}
	return json.Marshal(map[string]any{"function": f})
}

func decodeFunction(function map[string]any) (*Function, error) {
	name, ok := function["name"].(string)
	if !ok {
		return nil, errors.New("missing function name")
	}

	argsValue, ok := function["args"]
	if !ok {
		return &Function{Name: name}, nil
	}

	argsSlice, ok := argsValue.([]any)
	if !ok {
		return nil, errors.New("expected function args to be an array")
	}

	args := make([]Expression, len(argsSlice))
	for i, arg := range argsSlice {
		argument, err := decodeExpression(arg)
		if err != nil {
			return nil, fmt.Errorf("trouble parsing function argument %d: %w", i, err)
		}
		args[i] = argument
	}
	return &Function{Name: name, Args: args}, nil
}
