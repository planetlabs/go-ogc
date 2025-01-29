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

const (
	Equals              = "="
	NotEquals           = "<>"
	LessThan            = "<"
	LessThanOrEquals    = "<="
	GreaterThan         = ">"
	GreaterThanOrEquals = ">="
	likeOp              = "like"
	betweenOp           = "between"
	inOp                = "in"
	isNullOp            = "isNull"
)

type Comparison struct {
	Name  string
	Left  ScalarExpression
	Right ScalarExpression
}

var (
	_ Expression        = (*Comparison)(nil)
	_ BooleanExpression = (*Comparison)(nil)
	_ json.Marshaler    = (*Comparison)(nil)
)

func (*Comparison) expression()        {}
func (*Comparison) scalarExpression()  {}
func (*Comparison) booleanExpression() {}

func (e *Comparison) MarshalJSON() ([]byte, error) {
	args := []Expression{e.Left, e.Right}
	return marshalOp(e.Name, args)
}

func (e *Comparison) String() string {
	return toString(e)
}

type Like struct {
	Value   CharacterExpression
	Pattern PatternExpression
}

var (
	_ Expression        = (*Like)(nil)
	_ BooleanExpression = (*Like)(nil)
	_ json.Marshaler    = (*Like)(nil)
)

func (*Like) expression()        {}
func (*Like) scalarExpression()  {}
func (*Like) booleanExpression() {}

func (e *Like) MarshalJSON() ([]byte, error) {
	args := []Expression{e.Value, e.Pattern}
	return marshalOp(likeOp, args)
}

func (e *Like) String() string {
	return toString(e)
}

type Between struct {
	Value NumericExpression
	Low   NumericExpression
	High  NumericExpression
}

var (
	_ Expression        = (*Between)(nil)
	_ BooleanExpression = (*Between)(nil)
	_ json.Marshaler    = (*Between)(nil)
)

func (*Between) expression()        {}
func (*Between) scalarExpression()  {}
func (*Between) booleanExpression() {}

func (e *Between) MarshalJSON() ([]byte, error) {
	args := []Expression{e.Value, e.Low, e.High}
	return marshalOp(betweenOp, args)
}

func (e *Between) String() string {
	return toString(e)
}

type ScalarList []ScalarExpression

var (
	_ Expression = (ScalarList)(nil)
)

func (ScalarList) expression()       {}
func (ScalarList) scalarExpression() {}

func (e ScalarList) String() string {
	return sliceToString(e)
}

type In struct {
	Item ScalarExpression
	List ScalarList
}

var (
	_ Expression        = (*In)(nil)
	_ BooleanExpression = (*In)(nil)
	_ json.Marshaler    = (*In)(nil)
)

func (*In) expression()        {}
func (*In) scalarExpression()  {}
func (*In) booleanExpression() {}

func (e *In) MarshalJSON() ([]byte, error) {
	args := []Expression{e.Item, e.List}
	return marshalOp(inOp, args)
}

func (e *In) String() string {
	return toString(e)
}

type IsNull struct {
	Value Expression
}

var (
	_ Expression        = (*IsNull)(nil)
	_ BooleanExpression = (*IsNull)(nil)
	_ json.Marshaler    = (*IsNull)(nil)
)

func (*IsNull) expression()        {}
func (*IsNull) scalarExpression()  {}
func (*IsNull) booleanExpression() {}

func (e *IsNull) MarshalJSON() ([]byte, error) {
	args := []Expression{e.Value}
	return marshalOp(isNullOp, args)
}

func (e *IsNull) String() string {
	return toString(e)
}
