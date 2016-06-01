package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/thisisaaronland/go-slippy-tiles"
	"github.com/thisisaaronland/go-slippy-tiles/providers"
	"github.com/thisisaaronland/go-slippy-tiles/caches"	
	"github.com/whosonfirst/go-httpony/tls"	
	"net/http"
	"os"
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

func main() {

	var host = flag.String("host", "localhost", "...")
	var port = flag.Int("port", 9191, "...")
	var cors = flag.Bool("cors", false, "...")
	var tls = flag.Bool("tls", false, "...") // because CA warnings in browsers...
	var tls_cert = flag.String("tls-cert", "", "...")
	var tls_key = flag.String("tls-key", "", "...")
	var refresh = flag.Bool("refresh", false, "...")
	var cfg = flag.String("config", "", "...")

	flag.Parse()

	body, err := ioutil.ReadFile(*cfg)

	if err != nil {
		panic(err)
	}

	config := Config{}
	err = json.Unmarshal(body, &config)

	if err != nil {
		panic(err)
	}

	cache := config.Cache

	if cache.Name != "Disk" {
		err = errors.New("unsupported cache type")
		panic(err)
	}

	_, err = os.Stat(cache.Path)

	if os.IsNotExist(err) {
		err = errors.New("invalid cache path")
		panic(err)
	}

	cache, err := caches.NewDiskCache()

	if err != nil {
	   panic(err)
	}

	provider, err := providers.NewProxyProvider()

	if err != nil {
	   panic(err)
	}

	handler := provider.Handler()

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	if *tls {

		var cert string
		var key string

		if *tls_cert == "" && *tls_key == "" {

		   	root, err := httpony.EnsureTLSRoot()

			if err != nil {
				panic(err)
			}

			cert, key, err = httpony.GenerateTLSCert(*host, root)

			if err != nil {
				panic(err)
			}

		} else {
			cert = *tls_cert
			key = *tls_key
		}

		fmt.Printf("start and listen for requests at https://%s\n", endpoint)
		err = http.ListenAndServeTLS(endpoint, cert, key, handler)
		
	} else {
	
		fmt.Printf("start and listen for requests at http://%s\n", endpoint)
		err = http.ListenAndServe(endpoint, handler)
	}

	if err != nil {	
		panic(err)
	}

	os.Exit(0)
}
