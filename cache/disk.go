package cache

import (
       "github.com/thisisaaronland/go-tileproxy"
       "io/ioutil"
       "os"
       "path"
)

type DiskCache struct {
     tiles.Cache
     root string
}

func NewDiskCache(root string) (DiskCache, error) {

     _, err := os.Stat(root)

     if os.IsNotExist(err) {
     	return nil, err
     }
     
     c := DiskCache{
       root: root,
     }

     return c, nil
}

func (c DiskCache) Get(path string, body []byte) error {

     abs_path := path.Join(c.root, path)

     _, err := os.Stat(abs_path)

     if os.IsNotExist(err) {
     	return nil, err
     }

     return ioutil.ReadFile(abs_path)
}


func (c DiskCache) Set(path string, body []byte) error {

     abs_path := path.Join(c.root, path)

     root := filepath.Dir(abs_path)

     _, err = os.Stat(root)

     if os.IsNotExist(err) {
     	os.MkdirAll(root, 0755)
     }

     fh, err := os.Create(abs_path)

     if err != nil {
     	return err
     }

     defer fh.Close()
     fh.Write(body)
     fh.Sync()

     return nil
}

func (c DiskCache) Unset(path string) error {

     abs_path := path.Join(c.root, path)

     _, err := os.Stat(abs_path)

     if os.IsNotExist(err) {
     	return nil
     }

     return os.Remove(abs_path)
}     
