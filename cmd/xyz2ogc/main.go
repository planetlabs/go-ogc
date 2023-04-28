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

package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/alecthomas/kong"
	"github.com/go-playground/validator/v10"
	"github.com/planetlabs/go-ogc/cmd/xyz2ogc/internal/common"
	"github.com/planetlabs/go-ogc/cmd/xyz2ogc/internal/generate"
	"github.com/planetlabs/go-ogc/cmd/xyz2ogc/internal/serve"
)

var CLI struct {
	Serve    ServeCmd    `cmd:"" help:"Serve the OGC API - Tiles metadata."`
	Generate GenerateCmd `cmd:"" help:"Generate the OGC API - Tiles metadata."`
}

type ServeCmd struct {
	Config string `help:"Configuration file." default:"config.toml" type:"existingfile"`
}

func (c *ServeCmd) Run() error {
	config, configErr := getConfig(c.Config)
	if configErr != nil {
		return configErr
	}

	port := 0
	var origin *url.URL
	if config.Serve != nil {
		port = config.Serve.Port
		if config.Serve.Origin != "" {
			o, err := url.Parse(config.Serve.Origin)
			if err != nil {
				return fmt.Errorf("failed to parse %q as a URL: %w", config.Serve.Origin, err)
			}
			origin = o
		}
	}

	options := &serve.Options{
		Port:   port,
		Origin: origin,
		Tiles:  config.Tiles,
	}
	server, serverErr := serve.New(options)
	if serverErr != nil {
		return serverErr
	}

	return server.Start()
}

type GenerateCmd struct {
	Config string `help:"Configuration file." default:"config.toml" type:"existingfile"`
}

func (c *GenerateCmd) Run() error {
	config, configErr := getConfig(c.Config)
	if configErr != nil {
		return configErr
	}

	directory := "dist"
	var origin *url.URL
	if config.Generate != nil {
		if config.Generate.Directory != "" {
			directory = config.Generate.Directory
		}
		o, err := url.Parse(config.Generate.Origin)
		if err != nil {
			return fmt.Errorf("failed to parse %q as a URL: %w", config.Generate.Origin, err)
		}
		origin = o
	}

	options := &generate.Options{
		Tiles:  config.Tiles,
		Origin: origin,
		Dir:    directory,
	}

	if err := generate.Generate(options); err != nil {
		return err
	}

	count := len(config.Tiles)
	plural := "s"
	if count == 1 {
		plural = ""
	}
	fmt.Printf("wrote metadata for %d tileset%s to %q\n", count, plural, directory)
	return nil
}

func getConfig(configPath string) (*common.Config, error) {
	configData, readErr := os.ReadFile(configPath)
	if readErr != nil {
		return nil, readErr
	}

	config := &common.Config{}
	tomlErr := toml.Unmarshal(configData, config)
	if tomlErr != nil {
		return nil, tomlErr
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("invalid %q file: %w", configPath, err)
	}

	return config, nil
}

func main() {
	ctx := kong.Parse(&CLI, kong.UsageOnError())
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
