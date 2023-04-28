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
	"fmt"
)

const (
	GeometryContains   = "s_contains"
	GeometryCrosses    = "s_crosses"
	GeometryDisjoint   = "s_disjoint"
	GeometryEquals     = "s_equals"
	GeometryIntersects = "s_intersects"
	GeometryOverlaps   = "s_overlaps"
	GeometryTouches    = "s_touches"
	GeometryWithin     = "s_within"
)

type SpatialComparison struct {
	Name  string
	Left  SpatialExpression
	Right SpatialExpression
}

var (
	_ Expression        = (*SpatialComparison)(nil)
	_ BooleanExpression = (*SpatialComparison)(nil)
	_ json.Marshaler    = (*SpatialComparison)(nil)
)

func (*SpatialComparison) expression()        {}
func (*SpatialComparison) scalarExpression()  {}
func (*SpatialComparison) booleanExpression() {}

func (e *SpatialComparison) MarshalJSON() ([]byte, error) {
	args := []Expression{e.Left, e.Right}
	return marshalOp(e.Name, args)
}

type SpatialExpression interface {
	Expression
	spatialExpression()
}

type Geometry struct {
	Value any
}

var (
	_ Expression          = (*Geometry)(nil)
	_ SpatialExpression   = (*Geometry)(nil)
	_ ArrayItemExpression = (*Geometry)(nil)
	_ json.Marshaler      = (*Geometry)(nil)
)

func (*Geometry) expression()          {}
func (*Geometry) spatialExpression()   {}
func (*Geometry) arrayItemExpression() {}

func (e *Geometry) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Value)
}

var geometryTypes = map[string]bool{
	"Point":              true,
	"LineString":         true,
	"Polygon":            true,
	"MultiPoint":         true,
	"MultiLineString":    true,
	"MultiPolygon":       true,
	"GeometryCollection": true,
}

func decodeGeometry(value map[string]any) (*Geometry, error) {
	geomType, ok := value["type"].(string)
	if !ok {
		return nil, errors.New("geometry missing type")
	}

	if !geometryTypes[geomType] {
		return nil, fmt.Errorf("unexpected geometry type: %s", geomType)
	}

	if geomType == "GeometryCollection" {
		if _, ok := value["geometries"].([]any); !ok {
			return nil, errors.New("expected geometries array in geometry collection")
		}
	} else {
		if _, ok := value["coordinates"].([]any); !ok {
			return nil, errors.New("expected coordinates in geometry")
		}
	}

	return &Geometry{Value: value}, nil
}

type BoundingBox struct {
	Extent []float64
}

var (
	_ Expression          = (*BoundingBox)(nil)
	_ SpatialExpression   = (*BoundingBox)(nil)
	_ ArrayItemExpression = (*BoundingBox)(nil)
	_ json.Marshaler      = (*BoundingBox)(nil)
)

func (*BoundingBox) expression()          {}
func (*BoundingBox) spatialExpression()   {}
func (*BoundingBox) arrayItemExpression() {}

func (e *BoundingBox) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string][]float64{"bbox": e.Extent})
}

func decodeBoundingBox(bbox []any) (*BoundingBox, error) {
	count := len(bbox)
	if count != 4 && count != 6 {
		return nil, fmt.Errorf("expected 4 or 6 bbox values, found %d", count)
	}

	extent := make([]float64, len(bbox))
	for i, v := range bbox {
		b, ok := v.(float64)
		if !ok {
			return nil, fmt.Errorf("trouble decoding bbox value %v", v)
		}
		extent[i] = b
	}

	return &BoundingBox{Extent: extent}, nil
}
