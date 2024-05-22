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
	caseInsensitiveOp   = "casei"
	accentInsensitiveOp = "accenti"
)

type CharacterExpression interface {
	ScalarExpression
	characterExpression()
}

type PatternExpression interface {
	Expression
	patternExpression()
}

type CaseInsensitive struct {
	Value CharacterExpression
}

var (
	_ Expression          = (*CaseInsensitive)(nil)
	_ CharacterExpression = (*CaseInsensitive)(nil)
	_ PatternExpression   = (*CaseInsensitive)(nil)
	_ json.Marshaler      = (*CaseInsensitive)(nil)
)

func (*CaseInsensitive) expression()          {}
func (*CaseInsensitive) scalarExpression()    {}
func (*CaseInsensitive) characterExpression() {}
func (*CaseInsensitive) patternExpression()   {}

func (e *CaseInsensitive) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"op":   "casei",
		"args": []CharacterExpression{e.Value},
	}

	return json.Marshal(m)
}

type AccentInsensitive struct {
	Value CharacterExpression
}

var (
	_ Expression          = (*AccentInsensitive)(nil)
	_ CharacterExpression = (*AccentInsensitive)(nil)
	_ PatternExpression   = (*AccentInsensitive)(nil)
	_ json.Marshaler      = (*AccentInsensitive)(nil)
)

func (*AccentInsensitive) expression()          {}
func (*AccentInsensitive) scalarExpression()    {}
func (*AccentInsensitive) characterExpression() {}
func (*AccentInsensitive) patternExpression()   {}

func (e *AccentInsensitive) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"op":   "accenti",
		"args": []CharacterExpression{e.Value},
	}

	return json.Marshal(m)
}

type String struct {
	Value string
}

var (
	_ Expression          = (*String)(nil)
	_ CharacterExpression = (*String)(nil)
	_ ArrayItemExpression = (*String)(nil)
	_ PatternExpression   = (*String)(nil)
	_ json.Marshaler      = (*String)(nil)
)

func (*String) expression()          {}
func (*String) scalarExpression()    {}
func (*String) characterExpression() {}
func (*String) patternExpression()   {}
func (*String) arrayItemExpression() {}

func (e *String) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Value)
}
