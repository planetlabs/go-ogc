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
)

type Function struct {
	Op   string       `json:"op"`
	Args []Expression `json:"args"`
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
		"op":   e.Op,
		"args": e.Args,
	}
	if e.Args == nil {
		f["args"] = []Expression{}
	}
	return json.Marshal(f)
}
