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
	"bytes"
	"encoding/json"
	"fmt"
)

func marshalOp(name string, args []Expression) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "")

	fmt.Fprintf(buffer, `{"op":"%s","args":`, name)
	if err := encoder.Encode(args); err != nil {
		return nil, fmt.Errorf("failed to encode %q expression arguments: %w", name, err)
	}
	buffer.WriteString(`}`)

	return buffer.Bytes(), nil
}

func toString(e json.Marshaler) string {
	v, err := e.MarshalJSON()
	if err != nil {
		return fmt.Sprintf("%#v", e)
	}
	return string(v)
}

func sliceToString[T Expression](e []T) string {
	buffer := &bytes.Buffer{}
	buffer.WriteString("[")
	for i, item := range e {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(item.String())
	}
	buffer.WriteString("]")
	return buffer.String()
}

var argCount = map[string]int{
	notOp:               1,
	likeOp:              2,
	betweenOp:           3,
	inOp:                2,
	isNullOp:            1,
	caseInsensitiveOp:   1,
	accentInsensitiveOp: 1,
	Equals:              2,
	NotEquals:           2,
	LessThan:            2,
	LessThanOrEquals:    2,
	ArrayContainedBy:    2,
	ArrayContains:       2,
	ArrayEquals:         2,
	ArrayOverlaps:       2,
	GreaterThan:         2,
	GreaterThanOrEquals: 2,
	GeometryContains:    2,
	GeometryCrosses:     2,
	GeometryDisjoint:    2,
	GeometryEquals:      2,
	GeometryIntersects:  2,
	GeometryOverlaps:    2,
	GeometryTouches:     2,
	GeometryWithin:      2,
	TimeAfter:           2,
	TimeBefore:          2,
	TimeContains:        2,
	TimeDisjoint:        2,
	TimeDuring:          2,
	TimeEquals:          2,
	TimeFinishedBy:      2,
	TimeFinishes:        2,
	TimeIntersects:      2,
	TimeMeets:           2,
	TimeMetBy:           2,
	TimeOverlappedBy:    2,
	TimeOverlaps:        2,
	TimeStartedBy:       2,
	TimeStarts:          2,
}

func decodeOp(name string, encodedArgs []any) (Expression, error) {
	if fixedArgCount, ok := argCount[name]; ok && len(encodedArgs) != fixedArgCount {
		return nil, fmt.Errorf("expected %d args for %q op, found %d", fixedArgCount, name, len(encodedArgs))
	}

	args := make([]Expression, len(encodedArgs))
	for i, arg := range encodedArgs {
		argument, err := decodeExpression(arg)
		if err != nil {
			return nil, fmt.Errorf("trouble decoding arg %d for %q op: %w", i, name, err)
		}
		args[i] = argument
	}

	switch name {
	case notOp:
		boolArg, ok := args[0].(BooleanExpression)
		if !ok {
			return nil, fmt.Errorf("expected boolean arg for %q op", name)
		}
		return &Not{Arg: boolArg}, nil

	case andOp:
		if len(args) < 2 {
			return nil, fmt.Errorf("expected at least two args for %q op", name)
		}
		boolArgs, err := toBooleanArgs(name, args)
		if err != nil {
			return nil, err
		}
		return &And{Args: boolArgs}, nil

	case orOp:
		if len(args) < 2 {
			return nil, fmt.Errorf("expected at least two args for %q op", name)
		}
		boolArgs, err := toBooleanArgs(name, args)
		if err != nil {
			return nil, err
		}
		return &Or{Args: boolArgs}, nil

	case likeOp:
		str, ok := args[0].(CharacterExpression)
		if !ok {
			return nil, fmt.Errorf("expected a character expression for arg 0 of %q op", name)
		}
		pattern, ok := args[1].(PatternExpression)
		if !ok {
			return nil, fmt.Errorf("expected a pattern expression for arg 1 of %q op", name)
		}
		return &Like{Value: str, Pattern: pattern}, nil

	case betweenOp:
		numericArgs, err := toNumericArgs(name, args)
		if err != nil {
			return nil, err
		}
		return &Between{Value: numericArgs[0], Low: numericArgs[1], High: numericArgs[2]}, nil

	case inOp:
		item, ok := args[0].(ScalarExpression)
		if !ok {
			return nil, fmt.Errorf("expected scalar expression for arg 0 of %q op", name)
		}
		array, ok := args[1].(Array)
		if !ok {
			return nil, fmt.Errorf("expected an array for arg 1 of %q", name)
		}
		list := make([]ScalarExpression, len(array))
		for i, v := range array {
			scalar, ok := v.(ScalarExpression)
			if !ok {
				return nil, fmt.Errorf("expected scalar expression for item %d of arg 1 of %q op", i, name)
			}
			list[i] = scalar
		}
		return &In{Item: item, List: list}, nil

	case isNullOp:
		return &IsNull{Value: args[0]}, nil

	case caseInsensitiveOp:
		c, ok := args[0].(CharacterExpression)
		if !ok {
			return nil, fmt.Errorf("expected character expression in casei, got %v", args[0])
		}

		return &CaseInsensitive{Value: c}, nil

	case accentInsensitiveOp:
		c, ok := args[0].(CharacterExpression)
		if !ok {
			return nil, fmt.Errorf("expected character expression in accenti, got %v", args[0])
		}

		return &AccentInsensitive{Value: c}, nil

	case Equals, NotEquals, LessThan, LessThanOrEquals, GreaterThan, GreaterThanOrEquals:
		scalarArgs, err := toScalarArgs(name, args)
		if err != nil {
			return nil, err
		}
		return &Comparison{Name: name, Left: scalarArgs[0], Right: scalarArgs[1]}, nil

	case ArrayContainedBy, ArrayContains, ArrayEquals, ArrayOverlaps:
		arrayArgs, err := toArrayArgs(name, args)
		if err != nil {
			return nil, err
		}
		return &ArrayComparison{Name: name, Left: arrayArgs[0], Right: arrayArgs[1]}, nil

	case GeometryContains, GeometryCrosses, GeometryDisjoint, GeometryEquals, GeometryIntersects, GeometryOverlaps, GeometryTouches, GeometryWithin:
		spatialArgs, err := toSpatialArgs(name, args)
		if err != nil {
			return nil, err
		}
		return &SpatialComparison{Name: name, Left: spatialArgs[0], Right: spatialArgs[1]}, nil

	case TimeAfter, TimeBefore, TimeContains, TimeDisjoint, TimeDuring, TimeEquals, TimeFinishedBy, TimeFinishes, TimeIntersects, TimeMeets, TimeMetBy, TimeOverlappedBy, TimeOverlaps, TimeStartedBy, TimeStarts:
		temporalArgs, err := toTemporalArgs(name, args)
		if err != nil {
			return nil, err
		}
		return &TemporalComparison{Name: name, Left: temporalArgs[0], Right: temporalArgs[1]}, nil
	default:
		function := &Function{Op: name}
		if len(args) > 0 {
			function.Args = args
		}
		return function, nil
	}
}

func toScalarArgs(name string, args []Expression) ([]ScalarExpression, error) {
	scalarArgs := make([]ScalarExpression, len(args))
	for i, arg := range args {
		scalarArg, ok := arg.(ScalarExpression)
		if !ok {
			return nil, fmt.Errorf("expected arg %d to be a scalar expression for %q op", i, name)
		}
		scalarArgs[i] = scalarArg
	}
	return scalarArgs, nil
}

func toBooleanArgs(name string, args []Expression) ([]BooleanExpression, error) {
	boolArgs := make([]BooleanExpression, len(args))
	for i, arg := range args {
		boolArg, ok := arg.(BooleanExpression)
		if !ok {
			return nil, fmt.Errorf("expected arg %d to be a boolean expression for %q op", i, name)
		}
		boolArgs[i] = boolArg
	}
	return boolArgs, nil
}

func toNumericArgs(name string, args []Expression) ([]NumericExpression, error) {
	numericArgs := make([]NumericExpression, len(args))
	for i, arg := range args {
		numericArg, ok := arg.(NumericExpression)
		if !ok {
			return nil, fmt.Errorf("expected arg %d to be a numeric expression for %q op", i, name)
		}
		numericArgs[i] = numericArg
	}
	return numericArgs, nil
}

func toArrayArgs(name string, args []Expression) ([]ArrayExpression, error) {
	arrayArgs := make([]ArrayExpression, len(args))
	for i, arg := range args {
		arrayArg, ok := arg.(ArrayExpression)
		if !ok {
			return nil, fmt.Errorf("expected arg %d to be an array expression for %q op", i, name)
		}
		arrayArgs[i] = arrayArg
	}
	return arrayArgs, nil
}

func toSpatialArgs(name string, args []Expression) ([]SpatialExpression, error) {
	spatialArgs := make([]SpatialExpression, len(args))
	for i, arg := range args {
		spatialArg, ok := arg.(SpatialExpression)
		if !ok {
			return nil, fmt.Errorf("expected arg %d to be a spatial expression for %q op", i, name)
		}
		spatialArgs[i] = spatialArg
	}
	return spatialArgs, nil
}

func toTemporalArgs(name string, args []Expression) ([]TemporalExpression, error) {
	temporalArgs := make([]TemporalExpression, len(args))
	for i, arg := range args {
		temporalArg, ok := arg.(TemporalExpression)
		if !ok {
			return nil, fmt.Errorf("expected arg %d to be a temporal expression for %q op", i, name)
		}
		temporalArgs[i] = temporalArg
	}
	return temporalArgs, nil
}
