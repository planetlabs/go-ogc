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
)

type Filter struct {
	Expression BooleanExpression
}

var (
	_ Expression        = (*Filter)(nil)
	_ BooleanExpression = (*Filter)(nil)
	_ json.Marshaler    = (*Filter)(nil)
	_ json.Unmarshaler  = (*Filter)(nil)
)

func (*Filter) expression()        {}
func (*Filter) scalarExpression()  {}
func (*Filter) booleanExpression() {}

func (f *Filter) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Expression)
}

func (f *Filter) UnmarshalJSON(data []byte) error {
	expression, err := unmarshalExpression(data)
	if err != nil {
		return err
	}

	booleanExpression, ok := expression.(BooleanExpression)
	if !ok {
		return errors.New("expected a boolean expression")
	}

	f.Expression = booleanExpression
	return nil
}
