package cache

import (
	"errors"
	"github.com/thisisaaronland/go-slippy-tiles"
)

func NewCacheFromConfig(config slippytiles.Config) (slippytiles.Cache, error) {

	if config.Cache.Name != "Disk" {
		err := errors.New("unsupported cache type")
		return nil, err
	}

	cache, err := NewDiskCache(config)
	return cache, err
}
