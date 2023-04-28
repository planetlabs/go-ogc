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

package serve

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phayes/freeport"
	"github.com/planetlabs/go-ogc/cmd/xyz2ogc/internal/common"
	"github.com/sirupsen/logrus"
)

type Options struct {
	Port   int
	Origin *url.URL
	Tiles  []*common.TileSetConfig
}

type Server struct {
	tiles  []*common.TileSetConfig
	port   int
	origin *url.URL
	e      *echo.Echo
}

func New(options *Options) (*Server, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	port := options.Port
	if port == 0 {
		p, portErr := freeport.GetFreePort()
		if portErr != nil {
			return nil, portErr
		}
		port = p
	}
	server := &Server{
		tiles:  options.Tiles,
		port:   port,
		origin: options.Origin,
		e:      e,
	}

	e.Use(Logger(logrus.StandardLogger()))
	e.Use(ErrorCommitter())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", server.getRoot)
	e.GET("/api", server.getApi)
	e.GET("/docs/index.html", server.getDocs)
	e.GET("/conformance", server.getConformance)
	e.GET("/tiles", server.getTileSetList)
	e.GET("/tiles/:id", server.getTileSet)
	e.GET("/tileMatrixSets", server.getTileMatrixSetList)
	e.GET("/tileMatrixSets/:id", server.getTileMatrixSet)

	return server, nil
}

func (s *Server) getRoot(c echo.Context) error {
	base := s.getProxyBase(c)

	root, err := common.GetRoot(base)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, root)
}

func (s *Server) getApi(c echo.Context) error {
	schema, err := common.GetSchema()
	if err != nil {
		return err
	}
	c.Response().Header().Set(echo.HeaderContentType, "application/vnd.oai.openapi+json;version=3.0")
	return c.JSON(http.StatusOK, schema)
}

func (s *Server) getDocs(ctx echo.Context) error {
	data, dataErr := common.GetDocs("../api")
	if dataErr != nil {
		return dataErr
	}
	return ctx.HTMLBlob(http.StatusOK, data)
}

func (s *Server) getConformance(c echo.Context) error {
	base := s.getProxyBase(c)
	conformance, err := common.GetConformance(base)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, conformance)
}

func (s *Server) getTileSetList(c echo.Context) error {
	base := s.getProxyBase(c)
	tileSetList, err := common.GetTileSetList(base, s.tiles)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tileSetList)
}

func (s *Server) getTileSet(c echo.Context) error {
	id, idErr := strconv.Atoi(c.Param("id"))
	if idErr != nil {
		return notFound(fmt.Sprintf("no tileset with id %q", c.Param("id")))
	}
	base := s.getProxyBase(c)

	if id < 0 || id >= len(s.tiles) {
		return notFound(fmt.Sprintf("no tileset with id %q", c.Param("id")))
	}

	tileSet, tileSetErr := common.GetTileSet(base, id, s.tiles[id])
	if tileSetErr != nil {
		return tileSetErr
	}

	return c.JSON(http.StatusOK, tileSet)
}

func (s *Server) getTileMatrixSet(c echo.Context) error {
	id := c.Param("id")
	base := s.getProxyBase(c)

	tileMatrixSet, err := common.GetTileMatrixSet(id, base)
	if err != nil {
		return notFound(fmt.Sprintf("no tilematrixset with id %q", id))
	}
	return c.JSON(http.StatusOK, tileMatrixSet)
}

func (s *Server) getTileMatrixSetList(c echo.Context) error {
	base := s.getProxyBase(c)

	tileMatrixSetList, err := common.GetTileMatrixSetList(base, s.tiles)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tileMatrixSetList)
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("listening on http://localhost%s/\n", addr)
	return s.e.Start(addr)
}

func (s *Server) getProxyBase(c echo.Context) *url.URL {
	if s.origin != nil {
		return s.origin
	}
	return getProxyBase(c.Request())
}

func notFound(message string) error {
	return echo.NewHTTPError(http.StatusNotFound, message)
}
