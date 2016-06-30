package provider

import (
	"github.com/thisisaaronland/go-slippy-tiles"
	"github.com/thisisaaronland/go-slippy-tiles/cache"
	"io"
	"io/ioutil"
	_ "log"
	"net/http"
	"regexp"
)

type ProxyProvider struct {
	slippytiles.Provider
	config *slippytiles.Config
	cache  slippytiles.Cache
}

func NewProxyProvider(config *slippytiles.Config) (*ProxyProvider, error) {

	cache, err := cache.NewCacheFromConfig(config)

	if err != nil {
		return nil, err
	}

	p := ProxyProvider{
		config: config,
		cache:  cache,
	}

	return &p, nil
}

func (p ProxyProvider) Cache() slippytiles.Cache {
	return p.cache
}

func (p ProxyProvider) Handler(next http.Handler) http.Handler {

	re, _ := regexp.Compile(`/(.*)/(\d+)/(\d+)/(\d+).(\w+)$`)

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		url := req.URL
		path := url.Path
		query := url.RawQuery

		if re.MatchString(path) {

			cache := p.Cache()
			body, err := cache.Get(path)

			if err == nil {
				rsp.Write(body)
				return
			}

			m := re.FindStringSubmatch(path)
			layer_name := m[1]

			layer, ok := p.config.Layers[layer_name]

			if !ok {
				http.Error(rsp, "404 Not found", http.StatusNotFound)
				return
			}

			template, err := layer.URITemplate()

			if err != nil {
				http.Error(rsp, "500 Server Error", http.StatusInternalServerError)
				return
			}

			values := make(map[string]interface{})
			values["z"] = m[2]
			values["x"] = m[3]
			values["y"] = m[4]

			if len(layer.Formats) >= 1 {

				format := m[5]
				ok := false

				for _, f := range layer.Formats {
					if format == f {
						ok = true
						break
					}
				}

				if !ok {
					http.Error(rsp, "404 Not found", http.StatusNotFound)
					return
				}

				values["fmt"] = format
			}

			source, err := template.Expand(values)

			if err != nil {
				http.Error(rsp, "500 Server Error", http.StatusInternalServerError)
				return
			}

			if query != "" {
				source = source + "?" + query
			}

			client := &http.Client{}
			r, err := client.Get(source)

			if err != nil && err != io.EOF {
				http.Error(rsp, "502 Bad Gateway", http.StatusBadGateway)
				return
			}

			body, err = ioutil.ReadAll(r.Body)

			if err != nil {
				http.Error(rsp, "500 Server Error", http.StatusInternalServerError)
				return
			}

			if r.StatusCode == 200 {
				cache := p.Cache()
				go cache.Set(path, body)
			}

			// log.Println("return from source", path)

			rsp.Write(body)
			return

		}

		next.ServeHTTP(rsp, req)
	}

	return http.HandlerFunc(fn)
}
