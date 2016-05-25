package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/jtacoma/uritemplates"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

type ProxyConfig struct {
	Cache struct {
		Name string
		Path string
	}

	Layers map[string]ProxyProvider
}

type ProxyProvider struct {
	URL     string
	Formats []string
}

func (p ProxyProvider) Template() (*uritemplates.UriTemplate, error) {
	template, err := uritemplates.Parse(p.URL)
	return template, err
}

func main() {

	var host = flag.String("host", "localhost", "...")
	var port = flag.Int("port", 9191, "...")
	var cors = flag.Bool("cors", false, "Enable CORS headers")
	var cfg = flag.String("config", "", "...")

	flag.Parse()

	body, err := ioutil.ReadFile(*cfg)

	if err != nil {
		panic(err)
	}

	config := ProxyConfig{}
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

	re, _ := regexp.Compile(`/([^/]+)/(\d+)/(\d+)/(\d+).(\w+)$`)

	handler := func(rsp http.ResponseWriter, req *http.Request) {

		url := req.URL
		path := url.Path

		if !re.MatchString(path) {
			http.Error(rsp, "404 Not found", http.StatusNotFound)
			return
		}

		local_path := filepath.Join(cache.Path, path)

		_, err := os.Stat(local_path)

		if !os.IsNotExist(err) {

			body, err := ioutil.ReadFile(local_path)

			if err == nil {

				// something something something headers?
				// (20160524/thisisaaronland)

				if *cors {
					rsp.Header().Set("Access-Control-Allow-Origin", "*")
				}

				rsp.Write(body)
				return
			}

			fmt.Println("failed to read file", local_path, err)
		}

		m := re.FindStringSubmatch(path)
		layer := m[1]

		provider, ok := config.Layers[layer]

		if !ok {
			http.Error(rsp, "404 Not found", http.StatusNotFound)
			return
		}

		template, err := provider.Template()

		if err != nil {
			http.Error(rsp, "500 Server Error", http.StatusInternalServerError)
			return
		}

		values := make(map[string]interface{})
		values["z"] = m[2]
		values["x"] = m[3]
		values["y"] = m[4]

		if len(provider.Formats) >= 1 {

			format := m[5]
			ok := false

			for _, f := range provider.Formats {
				if format == f {
					ok = true
					break
				}
			}

			if !ok {
				http.Error(rsp, "404 Not found", http.StatusNotFound)
				return
			}

			values["fmt"] = format
		}

		source, err := template.Expand(values)

		if err != nil {
			http.Error(rsp, "500 Server Error", http.StatusInternalServerError)
			return
		}

		client := &http.Client{}
		r, err := client.Get(source)

		if err != nil && err != io.EOF {
			http.Error(rsp, "502 Bad Gateway", http.StatusBadGateway)
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(rsp, "500 Server Error", http.StatusInternalServerError)
			return
		}

		if r.StatusCode == 200 {

			// fmt.Println("caching", local_path)

			go func(local_path string, body []byte) {

				root := filepath.Dir(local_path)

				_, err = os.Stat(root)

				if os.IsNotExist(err) {
					os.MkdirAll(root, 0755)
				}

				fh, err := os.Create(local_path)

				if err != nil {
					fmt.Println(err)
					return
				}

				defer fh.Close()

				fh.Write(body)
				fh.Sync()

			}(local_path, body)
		}

		// HOW DO WE cache headers?
		// (20160524/thisisaaronland)

		/*
			for k, v := range r.Header {
				for _, vv := range v {
					rsp.Header().Add(k, vv)
				}
			}
		*/

		if *cors {
			rsp.Header().Set("Access-Control-Allow-Origin", "*")
		}

		rsp.Write(body)
		return
	}

	proxyHandler := http.HandlerFunc(handler)

	endpoint := fmt.Sprintf("%s:%d", *host, *port)
	err = http.ListenAndServe(endpoint, proxyHandler)

	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
