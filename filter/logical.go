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

const (
	orOp  = "or"
	andOp = "and"
	notOp = "not"
)

type Not struct {
	Arg BooleanExpression
}

var (
	_ Expression        = (*Not)(nil)
	_ BooleanExpression = (*Not)(nil)
	_ json.Marshaler    = (*Not)(nil)
)

func (*Not) expression()        {}
func (*Not) scalarExpression()  {}
func (*Not) booleanExpression() {}

func (e *Not) MarshalJSON() ([]byte, error) {
	args := []Expression{e.Arg}
	return marshalOp(notOp, args)
}

type And struct {
	Args []BooleanExpression
}

var (
	_ Expression        = (*And)(nil)
	_ BooleanExpression = (*And)(nil)
	_ json.Marshaler    = (*And)(nil)
)

func (*And) expression()        {}
func (*And) scalarExpression()  {}
func (*And) booleanExpression() {}

func (e *And) MarshalJSON() ([]byte, error) {
	args := make([]Expression, len(e.Args))
	for i, arg := range e.Args {
		args[i] = arg
	}
	return marshalOp(andOp, args)
}

type Or struct {
	Args []BooleanExpression
}

var (
	_ Expression        = (*Or)(nil)
	_ BooleanExpression = (*Or)(nil)
	_ json.Marshaler    = (*Or)(nil)
)

func (*Or) expression()        {}
func (*Or) scalarExpression()  {}
func (*Or) booleanExpression() {}

func (e *Or) MarshalJSON() ([]byte, error) {
	args := make([]Expression, len(e.Args))
	for i, arg := range e.Args {
		args[i] = arg
	}
	return marshalOp(orOp, args)
}
