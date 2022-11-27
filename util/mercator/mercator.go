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

// Package mercator provides transforms for converting geographic coordinates to
// and from Web Mercator coordinates.
package mercator

import "math"

const (
	Radius = 6378137
	Edge   = math.Pi * Radius
)

// Forward transforms geographic coordinates to Spherical Mercator.
func Forward(input []float64) []float64 {
	output := make([]float64, len(input))
	for i := 0; i < len(input); i += 2 {
		output[i] = (Edge * input[i]) / 180
		lat := input[i+1]
		if lat > 90 {
			lat = 90
		} else if lat < -90 {
			lat = -90
		}
		y := Radius * math.Log(math.Tan((math.Pi*(lat+90))/360))
		if y > Edge {
			y = Edge
		} else if y < -Edge {
			y = -Edge
		}
		output[i+1] = y
	}
	return output
}

// Inverse transforms Spherical Mercator coordinates to geographic.
func Inverse(input []float64) []float64 {
	output := make([]float64, len(input))
	for i := 0; i < len(input); i += 2 {
		output[i] = (180 * input[i]) / Edge
		output[i+1] = (360*math.Atan(math.Exp(input[i+1]/Radius)))/math.Pi - 90
	}
	return output
}
