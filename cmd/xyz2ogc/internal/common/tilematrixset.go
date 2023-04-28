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

package common

import (
	"fmt"
	"net/url"

	"github.com/planetlabs/go-ogc/api"
)

const defaultTileMatrixSetID = api.TMSWebMercatorQuad

func GetTileMatrixSet(id string, base *url.URL) (*api.TileMatrixSet, error) {
	tileMatrixSet, err := api.GetTileMatrixSet(id)
	if err != nil {
		return nil, fmt.Errorf("no tilematrixset with id %q", id)
	}

	tileMatrixSet.Links = []*api.Link{
		{
			Title: tileMatrixSet.ID,
			Rel:   "self",
			Type:  "application/json",
			Href:  getURL(base, fmt.Sprintf("/tileMatrixSets/%s", tileMatrixSet.ID)).String(),
		},
		{
			Title: "Tiling Schemes",
			Rel:   "http://www.opengis.net/def/rel/ogc/1.0/tiling-schemes",
			Type:  "application/json",
			Href:  getURL(base, "/tileMatrixSets").String(),
		},
	}
	return tileMatrixSet, nil
}

func GetTileMatrixSetList(base *url.URL, tiles []*TileSetConfig) (*api.TileMatrixSetList, error) {
	included := map[string]bool{}
	items := []*api.TileMatrixSetItem{}
	for _, t := range tiles {
		tileMatrixSetID := t.TileMatrixSetID
		if tileMatrixSetID == "" {
			tileMatrixSetID = defaultTileMatrixSetID
		}
		if !included[tileMatrixSetID] {
			tileMatrixSet, err := api.GetTileMatrixSet(tileMatrixSetID)
			if err != nil {
				return nil, err
			}
			items = append(items, &api.TileMatrixSetItem{
				ID: tileMatrixSet.ID,
				Links: []*api.Link{
					{
						Title: tileMatrixSet.ID,
						Rel:   "self",
						Type:  "application/json",
						Href:  getURL(base, fmt.Sprintf("/tileMatrixSets/%s", tileMatrixSet.ID)).String(),
					},
				},
			})
			included[tileMatrixSetID] = true
		}
	}

	tileMatrixSetList := &api.TileMatrixSetList{
		TileMatrixSets: items,
		Links: []*api.Link{
			{
				Title: "Service Root",
				Rel:   "service",
				Type:  "application/json",
				Href:  getURL(base, "/").String(),
			},
			{
				Title: "Tiling Schemes",
				Rel:   "self",
				Type:  "application/json",
				Href:  getURL(base, "/tileMatrixSets").String(),
			},
		},
	}

	return tileMatrixSetList, nil
}
