# go-tile-proxy

Too soon.

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

## See also

* https://github.com/thisisaaronland/mapzen-slippy-map
