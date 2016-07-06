# go-tile-proxy

What is the simplest TMS tile proxy?

## Install

```
make build
```

_See note below about installing [dependencies](#dependencies)._

## Usage

Still too soon.

### tile-proxy

```
./bin/tile-proxy -config config.json -cors
HIT cache/osm/all/16/10486/25367.mvt
HIT cache/osm/all/16/10486/25368.mvt
HIT cache/osm/all/16/10485/25368.mvt
HIT cache/osm/all/16/10485/25367.mvt
HIT cache/osm/all/16/10492/25367.mvt
HIT cache/osm/all/16/10492/25368.mvt
HIT cache/osm/all/16/10493/25367.mvt
HIT cache/osm/all/16/10493/25368.mvt
```

### tile-proxy config files

```
{
	"cache": { "name": "Disk", "path": "cache/" },
	"layers": {
		"osm/all": { "url": "https://vector.mapzen.com/osm/all/{z}/{x}/{y}.{fmt}", "formats": ["mvt", "topojson"] }
	}
}
```

### mapzen-slippy-map

This does not happen like-magic yet...

```
var s = slippy.map.scene();
var cfg = s.config.sources['osm'];
var url = cfg.url.replace('https://vector.mapzen.com', 'http://localhost:9191');
cfg.url = url;
s.setDataSource('osm', cfg);
```

## Dependencies

### Vendoring

Vendoring has been disabled for the time being because when trying to load this package as a vendored dependency in _another_ package it all goes pear-shape with errors like this:

```
make deps
# cd /Users/local/mapzen/mapzen-slippy-map/www-server/vendor/src/github.com/whosonfirst/go-httpony; git submodule update --init --recursive
fatal: no submodule mapping found in .gitmodules for path 'vendor/src/golang.org/x/net'
package github.com/whosonfirst/go-httpony: exit status 128
make: *** [deps] Error 1
```

_Note that's the actual error. It is copy-pasted from a different package with a similar issue. The problem is the same fot this package but always seems to involve something in `github.com/jtacoma/uritemplates`._

I have no idea and would welcome suggestions...

## See also

* https://github.com/thisisaaronland/mapzen-slippy-map
