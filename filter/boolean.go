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

type BooleanExpression interface {
	ScalarExpression
	booleanExpression()
}

type Boolean struct {
	Value bool
}

var (
	_ Expression          = (*Boolean)(nil)
	_ BooleanExpression   = (*Boolean)(nil)
	_ ArrayItemExpression = (*Boolean)(nil)
	_ json.Marshaler      = (*Boolean)(nil)
)

func (*Boolean) expression()          {}
func (*Boolean) scalarExpression()    {}
func (*Boolean) booleanExpression()   {}
func (*Boolean) arrayItemExpression() {}

func (e *Boolean) MarshalJSON() ([]byte, error) {
	if e.Value {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}
