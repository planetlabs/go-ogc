# Utilities for working with OGC APIs

The go-ogc module provides Go packages for working with OGC APIs.

## The api package

The `api` package provides structs for implementing the following standards:

 * [OGC API – Common](https://ogcapi.ogc.org/common/)
 * [OGC API – Tiles](https://ogcapi.ogc.org/tiles/)
 * [OGC API – Features](https://ogcapi.ogc.org/features/)

## The filter package

The `filter` package provides structs for encoding and decoding CQL2 filters as JSON.

## xyz2ogc

Generate [OGC API – Tiles](https://ogcapi.ogc.org/tiles/) metadata from exiting XYZ tilesets.

### Use

The `xyz2ogc` program reads a configuration file (`config.toml` by default) for information about a list of XYZ tilesets.  The program will then generate OGC API – Tiles compliant metadata for the configured tilesets.

The `config.toml` file is a [TOML](https://toml.io/) configuration file has [an array](https://toml.io/en/v1.0.0#array-of-tables) of `[[Tiles]]` tables describing the XYZ tilesets.  At a minimum, each `[[Tiles]]` table must have a `URL` key where the value is the XYZ tile URL template.

Here is an example `config.toml` with two tilesets:

```toml
[[Tiles]]
Title = "OpenStreetMap"
URL = "https://tile.openstreetmap.org/{z}/{x}/{y}.png"

[[Tiles]]
Title = "Stamen Terrain"
URL = "https://stamen-tiles-a.a.ssl.fastly.net/terrain-background/{z}/{x}/{y}.jpg"
```

In addition to the `URL` key, a tileset can be given a `Title`.  See below for more detail on the sections in the configuration options.

With a `config.toml` file in the current directory, you can run `xyz2ogc`:

```shell
xyz2ogc serve
```

This will start a server providing OGC API - Tiles compliant metadata.  See the command output for the server URL.

In addition to serving the OGC Tiles metadata locally, you can use the `xyz2ogc generate` command to generate a set of metadata files that can be served as a static website.  See details on the `generate` command below.

### Serve Configuration

The `xyz2ogc serve` command requires a configuration file with at least one `[[Tiles]]` section.  See below for details on the tiles configuration fields.

In addition, the configuration file can have a `[Serve]` section with configuration fields specific to the `xyz2ogc serve` command.

#### Port

By default, the `serve` command will find a random open port to listen for requests.  You can provide a `Port` number in the `config.toml` file.  For example:

```toml
[Serve]
Port = 8000
```

#### Origin

By default, the `serve` command generate metadata with links that use the origin of the request.  If you want links to be written with a different origin, you can provide an `Origin` string in the `config.toml` file.  For example:

```toml
[Serve]
Origin = "https://example.com"
```

### Generate Configuration

The `xyz2ogc generate` command requires a configuration file with at least one `[[Tiles]]` section.  See below for details on the tiles configuration fields.

In addition, the configuration file can have a `[Generate]` section with configuration fields specific to the `xyz2ogc generate` command.

#### Origin

By default, the `generate` command generate metadata with links without an origin in the URL (e.g. `"/api"`).  If you want links to be written with a different origin, you can provide an `Origin` string in the `config.toml` file.  For example:

```toml
[Generate]
Origin = "https://example.com"
```

#### Directory

By default, the `generate` command write tileset metadata to a `dist` directory.  If you want to write to a different location, you can provide a `Directory` string in the `config.toml` file.  For example:

```toml
[Generate]
Directory = "public"
```

### Tiles Configuration

Both the `xyz2ogc serve` and `xyz2ogc generate` commands require one or more `[[Tiles]]` sections in the `config.toml`.  Each `[[Tiles]]` section must have a `URL` string.  See below for the other tiles configuration fields.

#### URL

The `URL` field provides the URL template of an XYZ tileset.  The URL is expected to have `{z}`, `{x}`, and `{y}` variables in it.  For example:

```toml
[[Tiles]]
URL = "https://tile.openstreetmap.org/{z}/{x}/{y}.png"
```

#### Title

The `Title` field gives a tileset a title.  For example:

```toml
[[Tiles]]
Title = "OpenStreetMap"
URL = "https://tile.openstreetmap.org/{z}/{x}/{y}.png"
```

#### MinZoom

The `MinZoom` field allows limiting the minimum zoom level.  For example:

```toml
[[Tiles]]
URL = "https://tile.openstreetmap.org/{z}/{x}/{y}.png"
MinZoom = 2
```

#### MaxZoom

The `MaxZoom` field allows limiting the maximum zoom level.  For example:

```toml
[[Tiles]]
URL = "https://tile.openstreetmap.org/{z}/{x}/{y}.png"
MaxZoom = 10
```

#### Extent

The `Extent` field allows limiting the extent of the tileset.  The extent is described by `[minLon, minLat, maxLon, maxLat]` values.  For example:

```toml
[[Tiles]]
URL = "https://example.com/limited-extent/{z}/{x}/{y}.png"
Extent = [-120, 40, -110, 50]
```
