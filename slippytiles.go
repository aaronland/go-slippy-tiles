package slippytiles

import (
       "net/http"
)

type Config struct {
	Cache struct {
		Name string
		Path string
	}

	Layers struct {
		URL     string
		Formats []string
	}
}

type Cache interface {
     Get (string) ([]byte, error)
     Set (string, []byte) error
     Unset (string) error     
}

type Provider interface {
     Handler () http.Handler
     Cache () Cache
}
