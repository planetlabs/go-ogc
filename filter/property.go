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

import "encoding/json"

type Property struct {
	Name string
}

var (
	_ Expression          = (*Property)(nil)
	_ CharacterExpression = (*Property)(nil)
	_ NumericExpression   = (*Property)(nil)
	_ ArrayExpression     = (*Property)(nil)
	_ ArrayItemExpression = (*Property)(nil)
	_ SpatialExpression   = (*Property)(nil)
	_ TemporalExpression  = (*Property)(nil)
	_ json.Marshaler      = (*Property)(nil)
)

func (*Property) expression()          {}
func (*Property) scalarExpression()    {}
func (*Property) characterExpression() {}
func (*Property) numericExpression()   {}
func (*Property) arrayExpression()     {}
func (*Property) arrayItemExpression() {}
func (*Property) spatialExpression()   {}
func (*Property) temporalExpression()  {}

func (e *Property) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"property": e.Name})
}
