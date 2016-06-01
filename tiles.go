package tiles

import (
       "net/http"
)

/*
type ProxyConfig struct {
	Cache struct {
		Name string
		Path string
	}

	Layers map[string]ProxyProvider
}
*/

type Cache interface {
     Get(path string) (body []byte, error)
     Set(path string, body []byte) error
     Unset(path string) error     
}

type Provider interface {
     Handler() http.Handler
     Cache() tiles.Cache
}
