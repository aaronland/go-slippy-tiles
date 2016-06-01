package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/thisisaaronland/go-slippy-tiles"
	"github.com/thisisaaronland/go-slippy-tiles/provider"
	"github.com/whosonfirst/go-httpony/cors"
	"github.com/whosonfirst/go-httpony/tls"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {

	var host = flag.String("host", "localhost", "...")
	var port = flag.Int("port", 9191, "...")
	var cors_enable = flag.Bool("cors", false, "...")
	var cors_allow = flag.String("cors-allow", "*", "...")
	var tls_enable = flag.Bool("tls", false, "...") // because CA warnings in browsers...
	var tls_cert = flag.String("tls-cert", "", "...")
	var tls_key = flag.String("tls-key", "", "...")
	var cfg = flag.String("config", "", "...")

	flag.Parse()

	body, err := ioutil.ReadFile(*cfg)

	if err != nil {
		panic(err)
	}

	config := slippytiles.Config{}
	err = json.Unmarshal(body, &config)

	if err != nil {
		panic(err)
	}

	provider, err := provider.NewProviderFromConfig(config)

	if err != nil {
		panic(err)
	}

	endpoint := fmt.Sprintf("%s:%d", *host, *port)

	handler := cors.EnsureCORSHandler(provider.Handler(), *cors_enable, *cors_allow)

	if *tls_enable {

		var cert string
		var key string

		if *tls_cert == "" && *tls_key == "" {

			root, err := tls.EnsureTLSRoot()

			if err != nil {
				panic(err)
			}

			cert, key, err = tls.GenerateTLSCert(*host, root)

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
