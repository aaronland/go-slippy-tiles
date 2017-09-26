package cache

import (
	"database/sql"
	"errors"
	"github.com/thisisaaronland/go-slippy-tiles"
)

// https://github.com/mapbox/mbtiles-spec/blob/master/1.1/spec.md

type MBTilesDB struct {
	conn *sql.DB
}

func (db *MBTilesDB) Close() error {
	return db.conn.Close()
}

type MBTilesCache struct {
	slippytiles.Cache
	db *MBTilesDB
}

func NewMBTilesCache(config *slippytiles.Config) (*MBTilesCache, error) {
	return nil, errors.New("Please write me")
}

func (c *MBTilesCache) Get(rel_path string) ([]byte, error) {
	return nil, errors.New("Please write me")
}

func (c *MBTilesCache) Set(rel_path string, body []byte) error {
	return errors.New("Please write me")
}

func (c *MBTilesCache) Unset(rel_path string) error {
	return errors.New("Please write me")
}
