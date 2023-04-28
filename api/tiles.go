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

package api

// Structs for use with the Tiles spec.
// https://docs.ogc.org/is/20-057/20-057.html

// Data types for tilesets.
const (
	TileDataTypeVector   = "vector"
	TileDataTypeMap      = "map"
	TileDataTypeCoverage = "coverage"
)

// TileSet includes metadata for a tileset.
//
// Implementations that conform to the Tileset Conformance class (http://www.opengis.net/spec/ogcapi-tiles-1/1.0/conf/tileset)
// must provide these metadata fields.
type TileSet struct {
	Title               string                `json:"title,omitempty"`
	Description         string                `json:"description,omitempty"`
	Keywords            []string              `json:"keywords,omitempty"`
	Version             string                `json:"version,omitempty"`
	PointOfContact      string                `json:"pointOfContact,omitempty"`
	Attribution         string                `json:"attribution,omitempty"`
	License             string                `json:"license,omitempty"`
	AccessConstraints   []string              `json:"accessConstraints,omitempty"`
	MediaTypes          []string              `json:"mediaTypes,omitempty"`
	DataType            string                `json:"dataType"`
	TileMatrixSetURI    string                `json:"tileMatrixSetURI,omitempty"`
	TileMatrixSetLimits []*TileMatrixSetLimit `json:"tileMatrixSetLimits,omitempty"`
	CRS                 string                `json:"crs"`
	Epoch               float64               `json:"epoch,omitempty"`
	Links               []*Link               `json:"links"`
	BoundingBox         *BoundingBox          `json:"boundingBox,omitempty"`
	CenterPoint         *TilePoint            `json:"centerPoint,omitempty"`
}

type TilePoint struct {
	Coordinates      []float64 `json:"coordinates"`
	CRS              string    `json:"crs,omitempty"`
	TileMatrix       string    `json:"tileMatrix"`
	ScaleDenominator float64   `json:"scaleDenominator,omitempty"`
	CellSize         float64   `json:"cellSize,omitempty"`
}

// TileMatrixSetLimit adds metadata to a tileset describing limits on the referenced tiling scheme.
//
// Implementations that conform to the TileMatrixSetLimits Conformance class (http://www.opengis.net/spec/tms/2.0/conf/tilematrixsetlimits)
// may provide these metadata fields.
type TileMatrixSetLimit struct {
	TileMatrix string `json:"tileMatrix"`
	MinTileRow int    `json:"minTileRow"`
	MaxTileRow int    `json:"maxTileRow"`
	MinTileCol int    `json:"minTileCol"`
	MaxTileCol int    `json:"maxTileCol"`
}

// TileSetItem includes metadata for a a list of tilesets.
//
// Implementations that conform to the Tileset List Conformance class (http://www.opengis.net/spec/ogcapi-tiles-1/1.0/conf/tilesets-list)
// must provide these metadata fields when rendering a list of tilesets.
type TileSetItem struct {
	Title            string  `json:"title,omitempty"`
	DataType         string  `json:"dataType"`
	CRS              string  `json:"crs"`
	TileMatrixSetURI string  `json:"tileMatrixSetURI,omitempty"`
	Links            []*Link `json:"links"`
}

// TileSetList represents a list of tilesets.
//
// Implementations that conform to the Tileset List Conformance class (http://www.opengis.net/spec/ogcapi-tiles-1/1.0/conf/tilesets-list)
// must provide a list of tilesets.
type TileSetList struct {
	TileSets []*TileSetItem `json:"tilesets"`
	Links    []*Link        `json:"links,omitempty"`
}
