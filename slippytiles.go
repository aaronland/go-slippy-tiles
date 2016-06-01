package slippytiles

import (
	"github.com/jtacoma/uritemplates"
	"net/http"
)

type Config struct {
	Cache  CacheConfig
	Layers LayersConfig
}

type CacheConfig struct {
	Name string
	Path string
}

type LayersConfig map[string]LayerConfig

type LayerConfig struct {
	URL     string
	Formats []string
}

func (l LayerConfig) URITemplate() (*uritemplates.UriTemplate, error) {
	template, err := uritemplates.Parse(l.URL)
	return template, err
}

type Cache interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Unset(string) error
}

type Provider interface {
	Handler() http.Handler
	Cache() Cache
}
