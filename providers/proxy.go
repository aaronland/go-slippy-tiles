package providers

import (
       "fmt"
       "github.com/jtacoma/uritemplates"
       "github.com/thisisaaronland/go-slippy-tiles"
       "github.com/thisisaaronland/go-slippy-tiles/caches"
       "io"
       "io/ioutil"
       "net/http"
       "regexp"
       )

type ProxyProvider struct {
        slippytiles.Provider
	config	slippytiles.Config
	cache slippytiles.Cache
}

func NewProxyProvider(config slippytiles.Config) (*ProxyProvider, error){

     cache, err := caches.NewCacheFromConfig(config)

     if err != nil {
     	return nil, err
     }
     
     p := ProxyProvider{
       config: config,
       cache: cache,
     }

     return &p, nil
}

func (p ProxyProvider) Template() (*uritemplates.UriTemplate, error) {
	template, err := uritemplates.Parse(p.config.Layer.URL)
	return template, err
}

func (p ProxyProvider) Cache() slippytiles.Cache {
     return p.cache
}

func (p ProxyProvider) Handler() http.Handler {

     re, _ := regexp.Compile(`/(.*)/(\d+)/(\d+)/(\d+).(\w+)$`)
	
     fn := func(rsp http.ResponseWriter, req *http.Request){
     
		url := req.URL
		path := url.Path
		query := url.RawQuery

		if !re.MatchString(path) {
			http.Error(rsp, "404 Not found", http.StatusNotFound)
			return
		}

		/*
		if !*refresh {

			body, err := p.Cache.Get(path)

			if err == nil {

				if *cors {
				   rsp.Header().Set("Access-Control-Allow-Origin", "*")
				}

				rsp.Write(body)
				return
			}
		}
		*/
		
		m := re.FindStringSubmatch(path)
		layer := m[1]

		provider, ok := p.config.Layers[layer]

		if !ok {
			http.Error(rsp, "404 Not found", http.StatusNotFound)
			return
		}

		template, err := p.Template()

		if err != nil {
			http.Error(rsp, "500 Server Error", http.StatusInternalServerError)
			return
		}

		values := make(map[string]interface{})
		values["z"] = m[2]
		values["x"] = m[3]
		values["y"] = m[4]

		if len(provider.config.Formats) >= 1 {

			format := m[5]
			ok := false

			for _, f := range provider.Formats {
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

		fmt.Println("FETCH", source)

		client := &http.Client{}
		r, err := client.Get(source)

		if err != nil && err != io.EOF {
			http.Error(rsp, "502 Bad Gateway", http.StatusBadGateway)
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(rsp, "500 Server Error", http.StatusInternalServerError)
			return
		}

		if r.StatusCode == 200 {
			go p.Cache.Set(path, body)
		}

		rsp.Header().Set("Access-Control-Allow-Origin", "*")

		rsp.Write(body)
		return
	}

	return http.HandlerFunc(fn)
}
