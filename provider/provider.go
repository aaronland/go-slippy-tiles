package provider

import (
	"github.com/thisisaaronland/go-slippy-tiles"
)

func NewProviderFromConfig(config *slippytiles.Config) (slippytiles.Provider, error) {

	/*
		if config.Cache.Name != "Disk" {
			err := errors.New("unsupported cache type")
			return nil, err
		}
	*/

	provider, err := NewProxyProvider(config)
	return provider, err
}
