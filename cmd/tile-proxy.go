package main

import (
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

func main() {

	var host = flag.String("host", "localhost", "...")
	var port = flag.Int("port", 9191, "...")
	var cache = flag.String("cache", "", "...")
	var apikey = flag.String("api-key", "", "...")

	flag.Parse()

	mz_template, _ := uritemplates.Parse("https://vector.mapzen.com/osm/all/{z}/{x}/{y}.{fmt}?api_key={key}")

	re, _ := regexp.Compile(`/(\d+)/(\d+)/(\d+).(\w+)$`)

	handler := func(rsp http.ResponseWriter, req *http.Request) {

		url := req.URL
		path := url.Path

		if !re.MatchString(path) {
			http.Error(rsp, "404 Not found", http.StatusNotFound)
			return
		}

		local_path := filepath.Join(*cache, path)

		_, err := os.Stat(local_path)

		if !os.IsNotExist(err) {

			body, err := ioutil.ReadFile(local_path)

			if err == nil {

				// something something something headers?
				// (20160524/thisisaaronland)

				rsp.Write(body)
				return
			}

			fmt.Println("failed to read file", local_path, err)
		}

		m := re.FindStringSubmatch(path)

		values := make(map[string]interface{})
		values["z"] = m[1]
		values["x"] = m[2]
		values["y"] = m[3]
		values["fmt"] = m[4]
		values["key"] = *apikey

		source, err := mz_template.Expand(values)

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

			fmt.Println("caching", local_path)

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

		for k, v := range r.Header {
			for _, vv := range v {
				rsp.Header().Add(k, vv)
			}
		}

		rsp.Write(body)
		return
	}

	proxyHandler := http.HandlerFunc(handler)

	endpoint := fmt.Sprintf("%s:%d", *host, *port)
	err := http.ListenAndServe(endpoint, proxyHandler)

	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
