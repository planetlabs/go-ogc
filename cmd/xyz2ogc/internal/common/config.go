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
package common

type Config struct {
	Serve    *ServeConfig     `validate:"omitempty"`
	Generate *GenerateConfig  `validate:"omitempty"`
	Tiles    []*TileSetConfig `validate:"required,dive"`
}

type ServeConfig struct {
	Port   int `validate:"gte=0"`
	Origin string
}

type GenerateConfig struct {
	Origin    string
	Directory string
}

type TileSetConfig struct {
	URL             string `validate:"required,url,contains={z},contains={x},contains={y}"`
	Title           string
	Type            string
	TileType        string
	TileMatrixSetID string
	MinZoom         *int
	MaxZoom         *int
	Extent          []float64
}
