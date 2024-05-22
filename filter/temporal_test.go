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

package filter_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/planetlabs/go-ogc/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemporal(t *testing.T) {
	cases := []struct {
		filter *filter.Filter
		data   string
	}{
		{
			filter: &filter.Filter{
				Expression: &filter.Comparison{
					Name:  filter.Equals,
					Left:  &filter.Property{"date"},
					Right: &filter.Date{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
			data: `{
				"op": "=",
				"args": [
					{"property": "date"},
					{"date": "2020-01-01"}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.Comparison{
					Name:  filter.Equals,
					Left:  &filter.Property{"datetime"},
					Right: &filter.Timestamp{time.Date(2023, 2, 26, 23, 53, 29, 882000000, time.UTC)},
				},
			},
			data: `{
				"op": "=",
				"args": [
					{"property": "datetime"},
					{"timestamp": "2023-02-26T23:53:29.882Z"}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name:  filter.TimeAfter,
					Left:  &filter.Property{"datetime"},
					Right: &filter.Timestamp{time.Date(2023, 2, 26, 23, 53, 29, 882000000, time.UTC)},
				},
			},
			data: `{
				"op": "t_after",
				"args": [
					{"property": "datetime"},
					{"timestamp": "2023-02-26T23:53:29.882Z"}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name:  filter.TimeBefore,
					Left:  &filter.Property{"datetime"},
					Right: &filter.Date{time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
			data: `{
				"op": "t_before",
				"args": [
					{"property": "datetime"},
					{"date": "2000-01-01"}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name: filter.TimeContains,
					Left: &filter.Interval{
						Start: &filter.Date{time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
						End:   &filter.Date{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
					},
					Right: &filter.Property{"datetime"},
				},
			},
			data: `{
				"op": "t_contains",
				"args": [
					{"interval": ["2000-01-01", "2020-01-01"]},
					{"property": "datetime"}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name: filter.TimeDisjoint,
					Left: &filter.Interval{
						Start: &filter.Property{"start0"},
						End:   &filter.Property{"end0"},
					},
					Right: &filter.Interval{
						Start: &filter.Property{"start1"},
						End:   &filter.Property{"end1"},
					},
				},
			},
			data: `{
				"op": "t_disjoint",
				"args": [
					{"interval": [{"property": "start0"}, {"property": "end0"}]},
					{"interval": [{"property": "start1"}, {"property": "end1"}]}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name: filter.TimeDuring,
					Left: &filter.Interval{
						Start: &filter.Property{"start0"},
						End:   &filter.Property{"end0"},
					},
					Right: &filter.Interval{
						Start: &filter.Property{"start1"},
						End:   &filter.Property{"end1"},
					},
				},
			},
			data: `{
				"op": "t_during",
				"args": [
					{"interval": [{"property": "start0"}, {"property": "end0"}]},
					{"interval": [{"property": "start1"}, {"property": "end1"}]}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name: filter.TimeEquals,
					Left: &filter.Function{
						Op: "later",
						Args: []filter.Expression{
							&filter.Property{"datetime"},
							&filter.String{"1 month"},
						},
					},
					Right: &filter.Date{time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
				},
			},
			data: `{
				"op": "t_equals",
				"args": [
					{
						"op": "later",
						"args": [{"property": "datetime"}, "1 month"]
					},
					{"date": "2020-01-01"}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name: filter.TimeFinishedBy,
					Left: &filter.Interval{
						Start: &filter.Property{"start0"},
						End:   &filter.Property{"end0"},
					},
					Right: &filter.Interval{
						Start: &filter.Property{"start1"},
						End:   &filter.Property{"end1"},
					},
				},
			},
			data: `{
				"op": "t_finishedBy",
				"args": [
					{"interval": [{"property": "start0"}, {"property": "end0"}]},
					{"interval": [{"property": "start1"}, {"property": "end1"}]}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name: filter.TimeFinishes,
					Left: &filter.Interval{
						Start: &filter.Property{"start0"},
						End:   &filter.Property{"end0"},
					},
					Right: &filter.Interval{
						Start: &filter.Property{"start1"},
						End:   &filter.Property{"end1"},
					},
				},
			},
			data: `{
				"op": "t_finishes",
				"args": [
					{"interval": [{"property": "start0"}, {"property": "end0"}]},
					{"interval": [{"property": "start1"}, {"property": "end1"}]}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name: filter.TimeIntersects,
					Left: &filter.Interval{
						Start: &filter.Property{"start0"},
						End:   &filter.Property{"end0"},
					},
					Right: &filter.Interval{
						Start: &filter.Property{"start1"},
						End:   &filter.Property{"end1"},
					},
				},
			},
			data: `{
				"op": "t_intersects",
				"args": [
					{"interval": [{"property": "start0"}, {"property": "end0"}]},
					{"interval": [{"property": "start1"}, {"property": "end1"}]}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name: filter.TimeIntersects,
					Left: &filter.Property{"datetime"},
					Right: &filter.Interval{
						Start: &filter.Timestamp{Value: time.Date(2020, time.November, 11, 0, 0, 0, 0, time.UTC)},
						End:   &filter.Timestamp{Value: time.Date(2020, time.November, 12, 0, 0, 0, 0, time.UTC)},
					},
				},
			},
			data: `{
				"op": "t_intersects",
				"args": [
					{"property": "datetime"},
					{"interval": ["2020-11-11T00:00:00Z", "2020-11-12T00:00:00Z"] }
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name: filter.TimeMeets,
					Left: &filter.Interval{
						Start: &filter.Property{"start0"},
						End:   &filter.Property{"end0"},
					},
					Right: &filter.Interval{
						Start: &filter.Property{"start1"},
						End:   &filter.Property{"end1"},
					},
				},
			},
			data: `{
				"op": "t_meets",
				"args": [
					{"interval": [{"property": "start0"}, {"property": "end0"}]},
					{"interval": [{"property": "start1"}, {"property": "end1"}]}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name: filter.TimeMetBy,
					Left: &filter.Interval{
						Start: &filter.Property{"start0"},
						End:   &filter.Property{"end0"},
					},
					Right: &filter.Interval{
						Start: &filter.Property{"start1"},
						End:   &filter.Property{"end1"},
					},
				},
			},
			data: `{
				"op": "t_metBy",
				"args": [
					{"interval": [{"property": "start0"}, {"property": "end0"}]},
					{"interval": [{"property": "start1"}, {"property": "end1"}]}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name: filter.TimeOverlappedBy,
					Left: &filter.Interval{
						Start: &filter.Property{"start0"},
						End:   &filter.Property{"end0"},
					},
					Right: &filter.Interval{
						Start: &filter.Property{"start1"},
						End:   &filter.Property{"end1"},
					},
				},
			},
			data: `{
				"op": "t_overlappedBy",
				"args": [
					{"interval": [{"property": "start0"}, {"property": "end0"}]},
					{"interval": [{"property": "start1"}, {"property": "end1"}]}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name: filter.TimeOverlaps,
					Left: &filter.Interval{
						Start: &filter.Property{"start0"},
						End:   &filter.Property{"end0"},
					},
					Right: &filter.Interval{
						Start: &filter.Property{"start1"},
						End:   &filter.Property{"end1"},
					},
				},
			},
			data: `{
				"op": "t_overlaps",
				"args": [
					{"interval": [{"property": "start0"}, {"property": "end0"}]},
					{"interval": [{"property": "start1"}, {"property": "end1"}]}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name: filter.TimeStartedBy,
					Left: &filter.Interval{
						Start: &filter.Property{"start0"},
						End:   &filter.Property{"end0"},
					},
					Right: &filter.Interval{
						Start: &filter.Property{"start1"},
						End:   &filter.Property{"end1"},
					},
				},
			},
			data: `{
				"op": "t_startedBy",
				"args": [
					{"interval": [{"property": "start0"}, {"property": "end0"}]},
					{"interval": [{"property": "start1"}, {"property": "end1"}]}
				]
			}`,
		},
		{
			filter: &filter.Filter{
				Expression: &filter.TemporalComparison{
					Name: filter.TimeStarts,
					Left: &filter.Interval{
						Start: &filter.Property{"start0"},
						End:   &filter.Property{"end0"},
					},
					Right: &filter.Interval{
						Start: &filter.Property{"start1"},
						End:   &filter.Property{"end1"},
					},
				},
			},
			data: `{
				"op": "t_starts",
				"args": [
					{"interval": [{"property": "start0"}, {"property": "end0"}]},
					{"interval": [{"property": "start1"}, {"property": "end1"}]}
				]
			}`,
		},
	}

	schema := getSchema(t)
	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			data, err := json.Marshal(c.filter)
			require.NoError(t, err)
			assert.JSONEq(t, c.data, string(data))

			v := map[string]any{}
			require.NoError(t, json.Unmarshal(data, &v))
			if err := schema.Validate(v); err != nil {
				t.Errorf("failed to validate\n%#v", err)
			}

			filter := &filter.Filter{}
			require.NoError(t, json.Unmarshal([]byte(c.data), filter))
			assert.Equal(t, c.filter, filter)
		})
	}
}
