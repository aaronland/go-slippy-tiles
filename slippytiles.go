package slippytiles

import (
       "net/http"
)

type Cache interface {
     Get (string) ([]byte, error)
     Set (string, []byte) error
     Unset (string) error     
}

type Provider interface {
     Handler () http.Handler
     Cache () Cache
}
