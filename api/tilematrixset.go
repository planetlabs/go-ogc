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
package api

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

const (
	TopLeft    = "topLeft"
	BottomLeft = "bottomLeft"
)

// Structs for the Two Dimensional Tile Matrix Set and Tile Set Metadata spec.
// https://docs.ogc.org/is/17-083r4/17-083r4.html

// TileMatrix includes for a single "level" in a tile matrix set.
//
// Implementations that conform to the Tileset Conformance class (http://www.opengis.net/spec/ogcapi-tiles-1/1.0/conf/tileset)
// must link to a tiling scheme where each level includes this metadata.
type TileMatrix struct {
	Title                string                 `json:"title,omitempty"`
	Description          string                 `json:"description,omitempty"`
	Keywords             []string               `json:"keywords,omitempty"`
	ID                   string                 `json:"id"`
	ScaleDenominator     float64                `json:"scaleDenominator"`
	CellSize             float64                `json:"cellSize"`
	CornerOfOrigin       string                 `json:"cornerOfOrigin,omitempty"`
	PointOfOrigin        []float64              `json:"pointOfOrigin"`
	TileWidth            int                    `json:"tileWidth"`
	TileHeight           int                    `json:"tileHeight"`
	MatrixWidth          int                    `json:"matrixWidth"`
	MatrixHeight         int                    `json:"matrixHeight"`
	VariableMatrixWidths []*VariableMatrixWidth `json:"variableMatrixWidths,omitempty"`
}

func quantize(value float64, delta float64) float64 {
	return delta * math.Round(value/delta)
}

// Limit returns the tile range for a single matrix level.  Pass a nil bounds to get the full tile range for a level.
func (matrix *TileMatrix) Limit(bounds *BoundingBox) (*TileMatrixSetLimit, error) {
	matrixBounds := matrix.BoundingBox()
	if bounds == nil {
		bounds = matrix.BoundingBox()
	}

	tileMapWidth := matrix.CellSize * float64(matrix.TileWidth)
	tileMapHeight := matrix.CellSize * float64(matrix.TileHeight)

	offsetLeft := math.Floor((quantize(bounds.LowerLeft[0]-matrixBounds.LowerLeft[0], matrix.CellSize)) / tileMapWidth)
	offsetRight := math.Floor(quantize(matrixBounds.UpperRight[0]-bounds.UpperRight[0], matrix.CellSize) / tileMapWidth)

	offsetBottom := math.Floor(quantize(bounds.LowerLeft[1]-matrixBounds.LowerLeft[1], matrix.CellSize) / tileMapHeight)
	offsetTop := math.Floor(quantize(matrixBounds.UpperRight[1]-bounds.UpperRight[1], matrix.CellSize) / tileMapHeight)

	limit := &TileMatrixSetLimit{
		TileMatrix: matrix.ID,
		MinTileCol: int(offsetLeft),
		MaxTileCol: matrix.MatrixWidth - 1 - int(offsetRight),
		MinTileRow: int(offsetTop),
		MaxTileRow: matrix.MatrixHeight - 1 - int(offsetBottom),
	}
	return limit, nil
}

// BoundingBox returns the map extent of a tile matrix.
func (matrix *TileMatrix) BoundingBox() *BoundingBox {
	tileMapWidth := matrix.CellSize * float64(matrix.TileWidth)
	tileMapHeight := matrix.CellSize * float64(matrix.TileHeight)

	origin := matrix.PointOfOrigin
	left := origin[0]
	right := left + tileMapWidth*float64(matrix.MatrixWidth)

	var bottom float64
	var top float64
	if matrix.CornerOfOrigin == BottomLeft {
		bottom = origin[1]
		top = bottom + tileMapHeight*float64(matrix.MatrixHeight)
	} else {
		top = origin[1]
		bottom = top - tileMapHeight*float64(matrix.MatrixHeight)
	}

	return &BoundingBox{
		LowerLeft:  []float64{left, bottom},
		UpperRight: []float64{right, top},
	}
}

type TileMatrixSetList struct {
	TileMatrixSets []*TileMatrixSetItem `json:"tileMatrixSets"`
	Links          []*Link              `json:"links,omitempty"`
}

// TileMatrixSet includes metadata tiling schema.
//
// Implementations that conform to the Tileset Conformance class (http://www.opengis.net/spec/ogcapi-tiles-1/1.0/conf/tileset)
// must link to a tiling scheme for each tileset.
type TileMatrixSet struct {
	ID                string        `json:"id,omitempty"`
	Title             string        `json:"title,omitempty"`
	Description       string        `json:"description,omitempty"`
	Keywords          []string      `json:"keywords,omitempty"`
	URI               string        `json:"uri,omitempty"`
	OrderedAxes       []string      `json:"orderedAxes,omitempty"`
	CRS               string        `json:"crs"`
	WellKnownScaleSet string        `json:"wellKnownScaleSet,omitempty"`
	BoundingBox       *BoundingBox  `json:"boundingBox,omitempty"`
	TileMatrices      []*TileMatrix `json:"tileMatrices"`
	Links             []*Link       `json:"links,omitempty"`
}

// Limits generates a slice of tile ranges for the given levels and bounds.  Pass a nil bounds to get the full tile range for a level.
func (tms *TileMatrixSet) Limits(minLevel int, maxLevel int, bounds *BoundingBox) ([]*TileMatrixSetLimit, error) {
	if bounds == nil {
		bounds = tms.BoundingBox
	}
	if minLevel > maxLevel {
		return nil, errors.New("minLevel must be less than or equal to maxLevel")
	}
	if minLevel < 0 || maxLevel >= len(tms.TileMatrices) {
		return nil, errors.New("level out of range")
	}

	limits := []*TileMatrixSetLimit{}
	for level := minLevel; level <= maxLevel; level += 1 {
		matrix := tms.TileMatrices[level]
		limit, err := matrix.Limit(bounds)
		if err != nil {
			return nil, err
		}
		limits = append(limits, limit)
	}
	return limits, nil
}

type TileMatrixSetItem struct {
	ID    string  `json:"id,omitempty"`
	Title string  `json:"title,omitempty"`
	URI   string  `json:"uri,omitempty"`
	CRS   string  `json:"crs,omitempty"`
	Links []*Link `json:"links"`
}

type VariableMatrixWidth struct {
	Coalesce   int `json:"coalesce"`
	MinTileRow int `json:"minTileRow"`
	MaxTileRow int `json:"maxTileRow"`
}

type BoundingBox struct {
	CRS         string    `json:"crs,omitempty"`
	OrderedAxes []string  `json:"orderedAxes,omitempty"`
	LowerLeft   []float64 `json:"lowerLeft"`
	UpperRight  []float64 `json:"upperRight"`
}

const (
	TMSWebMercatorQuad = "WebMercatorQuad"
)

var forwardsAxisOrder = []string{"E", "N"}

type tileMatrixSetConfig struct {
	id                  string
	minZ                int
	maxZ                int
	maxCellSize         float64
	maxScaleDenominator float64
	origin              []float64
	crs                 string
	cornerOfOrigin      string
	tileWidth           int
	tileHeight          int
	orderedAxes         []string
	wellKnownScaleSet   string
}

var tmsLookup = map[string]*tileMatrixSetConfig{
	TMSWebMercatorQuad: {
		id:                  TMSWebMercatorQuad,
		minZ:                0,
		maxZ:                29,
		maxCellSize:         156543.03392804097,
		maxScaleDenominator: 559082264.0287178,
		origin:              []float64{-20037508.342789244, 20037508.342789244},
		cornerOfOrigin:      TopLeft,
		tileWidth:           256,
		tileHeight:          256,
		crs:                 crsEPSG3857,
		orderedAxes:         forwardsAxisOrder,
		wellKnownScaleSet:   wkssGoogle,
	},
}

type TileMatrixSetOption interface {
	apply(tileMatrixSetConfig) *tileMatrixSetConfig
}

type TileMatrixSetMinZoom int

func (minZ TileMatrixSetMinZoom) apply(config tileMatrixSetConfig) *tileMatrixSetConfig {
	config.minZ = int(minZ)
	return &config
}

var _ TileMatrixSetOption = TileMatrixSetMinZoom(0)

type TileMatrixSetMaxZoom int

func (maxZ TileMatrixSetMaxZoom) apply(config tileMatrixSetConfig) *tileMatrixSetConfig {
	config.maxZ = int(maxZ)
	return &config
}

var _ TileMatrixSetOption = TileMatrixSetMaxZoom(0)

func GetTileMatrixSet(id string, options ...TileMatrixSetOption) (*TileMatrixSet, error) {
	config, ok := tmsLookup[id]
	if !ok {
		return nil, fmt.Errorf("no tile matrix set with id %q", id)
	}
	for _, option := range options {
		config = option.apply(*config)
	}
	tileMatrixSet := &TileMatrixSet{
		ID:           id,
		URI:          getTileMatrixSetURI(id),
		CRS:          config.crs,
		OrderedAxes:  config.orderedAxes,
		TileMatrices: getTileMatrices(config),
	}
	return tileMatrixSet, nil
}

func getTileMatrices(config *tileMatrixSetConfig) []*TileMatrix {
	tileMatrices := []*TileMatrix{}
	for z := config.minZ; z <= config.maxZ; z += 1 {
		factor := math.Pow(2, float64(z))
		tileMatrix := &TileMatrix{
			ID:               strconv.Itoa(z),
			CellSize:         config.maxCellSize / factor,
			ScaleDenominator: config.maxScaleDenominator / factor,
			PointOfOrigin:    config.origin,
			CornerOfOrigin:   config.cornerOfOrigin,
			TileWidth:        config.tileWidth,
			TileHeight:       config.tileHeight,
			MatrixWidth:      int(factor),
			MatrixHeight:     int(factor),
		}
		tileMatrices = append(tileMatrices, tileMatrix)
	}
	return tileMatrices
}
