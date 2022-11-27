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
package generate

import (
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/planetlabs/go-ogc/cmd/xyz2ogc/internal/common"
)

type Options struct {
	Dir    string
	Index  string
	Origin *url.URL
	Tiles  []*common.TileSetConfig
}

type Generator struct {
	dir    string
	index  string
	origin *url.URL
	tiles  []*common.TileSetConfig
}

const defaultIndex = "index.json"

func Generate(options *Options) error {
	index := options.Index
	if index == "" {
		index = defaultIndex
	}
	if options.Origin == nil {
		return errors.New("missing Origin option")
	}
	generator := &Generator{
		dir:    options.Dir,
		index:  index,
		origin: options.Origin,
		tiles:  options.Tiles,
	}
	return generator.Generate()
}

func (g *Generator) Generate() error {
	if err := g.generateRoot(); err != nil {
		return err
	}

	if err := g.generateConformance(); err != nil {
		return err
	}

	if err := g.generateApi(); err != nil {
		return err
	}

	if err := g.generateDocs("../api/index.json"); err != nil {
		return err
	}

	if err := g.generateTileMatrixSetList(); err != nil {
		return err
	}

	if err := g.generateTileMatrixSets(); err != nil {
		return err
	}

	if err := g.generateTileSetList(); err != nil {
		return err
	}

	if err := g.generateTileSets(); err != nil {
		return err
	}

	return nil
}

func (g *Generator) generateRoot() error {
	filePath := path.Join(g.dir, g.index)
	dirPath := path.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0750); err != nil {
		return err
	}
	resource, resourceErr := common.GetRoot(g.origin)
	if resourceErr != nil {
		return resourceErr
	}
	data, jsonErr := json.Marshal(resource)
	if jsonErr != nil {
		return jsonErr
	}
	return os.WriteFile(filePath, data, 0666)
}

func (g *Generator) generateConformance() error {
	filePath := path.Join(g.dir, "conformance", g.index)
	dirPath := path.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0750); err != nil {
		return err
	}
	resource, resourceErr := common.GetConformance(g.origin)
	if resourceErr != nil {
		return resourceErr
	}
	data, jsonErr := json.Marshal(resource)
	if jsonErr != nil {
		return jsonErr
	}
	return os.WriteFile(filePath, data, 0666)
}

func (g *Generator) generateDocs(apiPath string) error {
	filePath := path.Join(g.dir, "docs", "index.html")
	dirPath := path.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0750); err != nil {
		return err
	}
	data, dataErr := common.GetDocs(apiPath)
	if dataErr != nil {
		return dataErr
	}
	return os.WriteFile(filePath, data, 0666)
}

func (g *Generator) generateApi() error {
	filePath := path.Join(g.dir, "api", g.index)
	dirPath := path.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0750); err != nil {
		return err
	}
	resource, resourceErr := common.GetSchema()
	if resourceErr != nil {
		return resourceErr
	}
	data, jsonErr := json.Marshal(resource)
	if jsonErr != nil {
		return jsonErr
	}
	return os.WriteFile(filePath, data, 0666)
}

func (g *Generator) generateTileMatrixSetList() error {
	filePath := path.Join(g.dir, "tileMatrixSets", g.index)
	dirPath := path.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0750); err != nil {
		return err
	}
	resource, resourceErr := common.GetTileMatrixSetList(g.origin, g.tiles)
	if resourceErr != nil {
		return resourceErr
	}
	data, jsonErr := json.Marshal(resource)
	if jsonErr != nil {
		return jsonErr
	}
	return os.WriteFile(filePath, data, 0666)
}

func (g *Generator) generateTileMatrixSets() error {
	tileMatrixSetList, tmsErr := common.GetTileMatrixSetList(g.origin, g.tiles)
	if tmsErr != nil {
		return tmsErr
	}
	for _, tileMatrixSetItem := range tileMatrixSetList.TileMatrixSets {
		resource, resourceErr := common.GetTileMatrixSet(tileMatrixSetItem.ID, g.origin)
		if resourceErr != nil {
			return resourceErr
		}
		filePath := path.Join(g.dir, "tileMatrixSets", resource.ID, g.index)
		dirPath := path.Dir(filePath)
		if err := os.MkdirAll(dirPath, 0750); err != nil {
			return err
		}
		data, jsonErr := json.Marshal(resource)
		if jsonErr != nil {
			return jsonErr
		}
		if err := os.WriteFile(filePath, data, 0666); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) generateTileSetList() error {
	filePath := path.Join(g.dir, "tiles", g.index)
	dirPath := path.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0750); err != nil {
		return err
	}
	resource, resourceErr := common.GetTileSetList(g.origin, g.tiles)
	if resourceErr != nil {
		return resourceErr
	}
	data, jsonErr := json.Marshal(resource)
	if jsonErr != nil {
		return jsonErr
	}
	return os.WriteFile(filePath, data, 0666)
}

func (g *Generator) generateTileSets() error {
	for i, tileSetConfig := range g.tiles {
		filePath := path.Join(g.dir, "tiles", strconv.Itoa(i), g.index)
		dirPath := path.Dir(filePath)
		if err := os.MkdirAll(dirPath, 0750); err != nil {
			return err
		}
		resource, resourceErr := common.GetTileSet(g.origin, i, tileSetConfig)
		if resourceErr != nil {
			return resourceErr
		}
		data, jsonErr := json.Marshal(resource)
		if jsonErr != nil {
			return jsonErr
		}
		if err := os.WriteFile(filePath, data, 0666); err != nil {
			return err
		}
	}
	return nil
}
