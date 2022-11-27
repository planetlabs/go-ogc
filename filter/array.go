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
package filter

import (
	"encoding/json"
	"fmt"
)

const (
	ArrayContainedBy = "a_containedBy"
	ArrayContains    = "a_contains"
	ArrayEquals      = "a_equals"
	ArrayOverlaps    = "a_overlaps"
)

type ArrayComparison struct {
	Name  string
	Left  ArrayExpression
	Right ArrayExpression
}

var (
	_ Expression        = (*ArrayComparison)(nil)
	_ BooleanExpression = (*ArrayComparison)(nil)
	_ json.Marshaler    = (*ArrayComparison)(nil)
)

func (*ArrayComparison) expression()        {}
func (*ArrayComparison) scalarExpression()  {}
func (*ArrayComparison) booleanExpression() {}

func (e *ArrayComparison) MarshalJSON() ([]byte, error) {
	args := []Expression{e.Left, e.Right}
	return marshalOp(e.Name, args)
}

type ArrayItemExpression interface {
	Expression
	arrayItemExpression()
}

type ArrayExpression interface {
	Expression
	arrayExpression()
}

type Array []ArrayItemExpression

var (
	_ Expression          = (Array)(nil)
	_ ArrayExpression     = (Array)(nil)
	_ ArrayItemExpression = (Array)(nil)
)

func (Array) expression()          {}
func (Array) arrayExpression()     {}
func (Array) arrayItemExpression() {}

func decodeArray(values []any) (Array, error) {
	items := make([]ArrayItemExpression, len(values))
	for i, value := range values {
		expression, err := decodeExpression(value)
		if err != nil {
			return nil, fmt.Errorf("trouble decoding item %d from array: %w", i, err)
		}
		item, ok := expression.(ArrayItemExpression)
		if !ok {
			return nil, fmt.Errorf("unsupported type for item %d of array", i)
		}
		items[i] = item
	}
	return items, nil
}
