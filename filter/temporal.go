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
	"errors"
	"fmt"
	"time"
)

const (
	TimeAfter        = "t_after"
	TimeBefore       = "t_before"
	TimeContains     = "t_contains"
	TimeDisjoint     = "t_disjoint"
	TimeDuring       = "t_during"
	TimeEquals       = "t_equals"
	TimeFinishedBy   = "t_finishedBy"
	TimeFinishes     = "t_finishes"
	TimeIntersects   = "t_intersects"
	TimeMeets        = "t_meets"
	TimeMetBy        = "t_metBy"
	TimeOverlappedBy = "t_overlappedBy"
	TimeOverlaps     = "t_overlaps"
	TimeStartedBy    = "t_startedBy"
	TimeStarts       = "t_starts"
)

type TemporalComparison struct {
	Name  string
	Left  TemporalExpression
	Right TemporalExpression
}

var (
	_ Expression        = (*TemporalComparison)(nil)
	_ BooleanExpression = (*TemporalComparison)(nil)
	_ json.Marshaler    = (*TemporalComparison)(nil)
)

func (*TemporalComparison) expression()        {}
func (*TemporalComparison) scalarExpression()  {}
func (*TemporalComparison) booleanExpression() {}

func (e *TemporalComparison) MarshalJSON() ([]byte, error) {
	args := []Expression{e.Left, e.Right}
	return marshalOp(e.Name, args)
}

type TemporalExpression interface {
	Expression
	temporalExpression()
}

type InstantExpression interface {
	TemporalExpression
	ScalarExpression
}

type Date struct {
	Value time.Time
}

var (
	_ Expression         = (*Date)(nil)
	_ TemporalExpression = (*Date)(nil)
	_ InstantExpression  = (*Date)(nil)
	_ json.Marshaler     = (*Date)(nil)
)

func (*Date) expression()         {}
func (*Date) scalarExpression()   {}
func (*Date) temporalExpression() {}

func (e *Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"date": e.Value.Format(time.DateOnly)})
}

func decodeDate(value string) (*Date, error) {
	date, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return nil, fmt.Errorf("unable to parse date: %w", err)
	}
	return &Date{Value: date}, nil
}

type Timestamp struct {
	Value time.Time
}

var (
	_ Expression         = (*Timestamp)(nil)
	_ TemporalExpression = (*Timestamp)(nil)
	_ InstantExpression  = (*Timestamp)(nil)
	_ json.Marshaler     = (*Timestamp)(nil)
)

func (*Timestamp) expression()         {}
func (*Timestamp) scalarExpression()   {}
func (*Timestamp) temporalExpression() {}

func (e *Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{"timestamp": e.Value.Format(time.RFC3339Nano)})
}

func decodeTimestamp(value string) (*Timestamp, error) {
	timestamp, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return nil, fmt.Errorf("unable to parse timestamp: %w", err)
	}
	return &Timestamp{Value: timestamp}, nil
}

type Interval struct {
	Start InstantExpression
	End   InstantExpression
}

var (
	_ Expression         = (*Interval)(nil)
	_ TemporalExpression = (*Interval)(nil)
	_ json.Marshaler     = (*Interval)(nil)
)

func (*Interval) expression()         {}
func (*Interval) temporalExpression() {}

func (e *Interval) MarshalJSON() ([]byte, error) {
	items := make([]any, 2)
	if e.Start == nil {
		items[0] = ".."
	} else {
		items[0] = e.Start
	}
	if e.End == nil {
		items[1] = ".."
	} else {
		items[1] = e.End
	}
	return json.Marshal(map[string]any{"interval": items})
}

const nilInstant = ".."

func decodeInterval(values []any) (*Interval, error) {
	if len(values) != 2 {
		return nil, fmt.Errorf("expected 2 items for interval, found %d", len(values))
	}

	var start InstantExpression
	startValue, err := decodeExpression(values[0])
	if err != nil {
		return nil, fmt.Errorf("trouble parsing interval start: %w", err)
	}
	switch s := startValue.(type) {
	case *String:
		if s.Value != nilInstant {
			return nil, fmt.Errorf("expected date or timestamp expression, got %s", s.Value)
		}
	case *Date:
		start = s
	case *Timestamp:
		start = s
	case *Property:
		start = s
	case *Function:
		start = s
	default:
		return nil, errors.New("unsupported start expression in interval")
	}

	var end InstantExpression
	endValue, err := decodeExpression(values[1])
	if err != nil {
		return nil, fmt.Errorf("trouble parsing interval end: %w", err)
	}
	switch s := endValue.(type) {
	case *String:
		if s.Value != nilInstant {
			return nil, fmt.Errorf("expected date or timestamp expression, got %s", s.Value)
		}
	case *Date:
		end = s
	case *Timestamp:
		end = s
	case *Property:
		end = s
	case *Function:
		end = s
	default:
		return nil, errors.New("unsupported end expression in interval")
	}

	if start == nil && end == nil {
		return nil, errors.New("interval start or end must be provided")
	}

	return &Interval{Start: start, End: end}, nil
}
