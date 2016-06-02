package slippytiles

import (
	"encoding/json"
	"github.com/jtacoma/uritemplates"
	"io/ioutil"
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

func NewConfigFromFile(file string) (*Config, error) {

	body, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	c := Config{}
	err = json.Unmarshal(body, &c)

	if err != nil {
		return nil, err
	}

	return &c, nil
}
