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

import (
	"errors"
	"fmt"
	"mime"
	"net/url"
	"path"
	"strings"

	"github.com/planetlabs/go-ogc/api"
	"github.com/planetlabs/go-ogc/util/mercator"
)

func GetTileSetList(base *url.URL, tiles []*TileSetConfig) (*api.TileSetList, error) {
	tileSets := make([]*api.TileSetItem, len(tiles))
	for i, t := range tiles {
		tileSet, err := tileSetItemFromConfig(base, i, t)
		if err != nil {
			return nil, err
		}
		tileSets[i] = tileSet
	}

	tileSetList := &api.TileSetList{
		TileSets: tileSets,
		Links: []*api.Link{
			{
				Title: "Service Root",
				Rel:   "service",
				Type:  "application/json",
				Href:  getURL(base, "/").String(),
			},
		},
	}
	return tileSetList, nil
}

func GetTileSet(base *url.URL, id int, tileSetConfig *TileSetConfig) (*api.TileSet, error) {
	return tileSetFromConfig(base, id, tileSetConfig)
}

func rewriteTileTemplate(template string) string {
	template = strings.ReplaceAll(template, "{z}", "{tileMatrix}")
	template = strings.ReplaceAll(template, "{x}", "{tileCol}")
	template = strings.ReplaceAll(template, "{y}", "{tileRow}")
	return template
}

func tileSetItemFromConfig(base *url.URL, id int, t *TileSetConfig) (*api.TileSetItem, error) {
	tileSet, err := tileSetFromConfig(base, id, t)
	if err != nil {
		return nil, err
	}

	tileSetItem := &api.TileSetItem{
		Title:            t.Title,
		DataType:         tileSet.DataType,
		TileMatrixSetURI: tileSet.TileMatrixSetURI,
		CRS:              tileSet.CRS,
		Links:            tileSet.Links,
	}

	return tileSetItem, nil
}

func tileSetFromConfig(base *url.URL, id int, t *TileSetConfig) (*api.TileSet, error) {
	dataType := t.Type
	if dataType == "" {
		dataType = api.TileDataTypeMap
	}

	tileType := t.TileType
	if tileType == "" {
		tileType = mime.TypeByExtension(path.Ext(t.URL))
	}

	tileMatrixSetID := t.TileMatrixSetID
	if tileMatrixSetID == "" {
		tileMatrixSetID = defaultTileMatrixSetID
	}

	tileMatrixSet, tmsErr := api.GetTileMatrixSet(tileMatrixSetID)
	if tmsErr != nil {
		return nil, tmsErr
	}

	var boundingBox *api.BoundingBox
	if t.Extent != nil {
		if len(t.Extent) != 4 {
			return nil, errors.New("expected the Extent to have 4 items")
		}
		if tileMatrixSetID != api.TMSWebMercatorQuad {
			return nil, errors.New("can only transform extent for EPSG:3857")
		}
		mercatorExtent := mercator.Forward(t.Extent)
		boundingBox = &api.BoundingBox{
			LowerLeft:  mercatorExtent[:2],
			UpperRight: mercatorExtent[2:],
		}
	}

	var tileMatrixSetLimits []*api.TileMatrixSetLimit
	if boundingBox != nil || t.MinZoom != nil || t.MaxZoom != nil {
		minZoom := 0
		if t.MinZoom != nil {
			minZoom = *t.MinZoom
			if minZoom < 0 {
				return nil, fmt.Errorf("the MinZoom for tileset %d is less than zero: %d", id, minZoom)
			}
		}

		maxAvailableZoom := len(tileMatrixSet.TileMatrices) - 1
		maxZoom := maxAvailableZoom
		if t.MaxZoom != nil {
			maxZoom = *t.MaxZoom
			if maxZoom > maxAvailableZoom {
				return nil, fmt.Errorf("the MaxZoom for tileset %d exceeds highest zoom level: %d", id, maxAvailableZoom)
			}
		}

		limits, err := tileMatrixSet.Limits(minZoom, maxZoom, boundingBox)
		if err != nil {
			return nil, err
		}
		tileMatrixSetLimits = limits
	}

	tileSet := &api.TileSet{
		Title:               t.Title,
		DataType:            dataType,
		TileMatrixSetURI:    tileMatrixSet.URI,
		TileMatrixSetLimits: tileMatrixSetLimits,
		CRS:                 tileMatrixSet.CRS,
		BoundingBox:         boundingBox,
		Links: []*api.Link{
			{
				Title: "Tileset Metadata",
				Rel:   "self",
				Type:  "application/json",
				Href:  getURL(base, fmt.Sprintf("/tiles/%d", id)).String(),
			},
			{
				Title:     "Tile URL Template",
				Rel:       "item",
				Type:      tileType,
				Href:      rewriteTileTemplate(t.URL),
				Templated: true,
			},
			{
				Title: "Tiling Scheme",
				Rel:   "http://www.opengis.net/def/rel/ogc/1.0/tiling-scheme",
				Type:  "application/json",
				Href:  getURL(base, fmt.Sprintf("/tileMatrixSets/%s", tileMatrixSetID)).String(),
			},
			{
				Title: "Tileset List",
				Rel:   "collection",
				Type:  "application/json",
				Href:  getURL(base, "/tiles").String(),
			},
		},
	}
	return tileSet, nil
}
