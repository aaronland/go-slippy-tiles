package tiles

import (
       "net/http"
)

type Cache interface {
     Get(path string) (body []byte, error)
     Set(path string, body []byte) error
}

type Provider interface {
     Cache() Cache
     Handler() http.Handler
}
