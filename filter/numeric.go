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

type NumericExpression interface {
	ScalarExpression
	numericExpression()
}

type Number struct {
	Value float64
}

var (
	_ Expression          = (*Number)(nil)
	_ NumericExpression   = (*Number)(nil)
	_ ArrayItemExpression = (*Number)(nil)
	_ json.Marshaler      = (*Number)(nil)
)

func (*Number) expression()          {}
func (*Number) scalarExpression()    {}
func (*Number) numericExpression()   {}
func (*Number) arrayItemExpression() {}

func (e *Number) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Value)
}

func (e *Number) String() string {
	return toString(e)
}
