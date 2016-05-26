package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"github.com/jtacoma/uritemplates"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"
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

// https://github.com/mattrobenolt/https

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func generateCert(host string) (string, string) {
	var err error

	dir := "/tmp/https-certs/" + host + "/"
	certPath := dir + "cert.pem"
	keyPath := dir + "key.pem"

	if exists(certPath) && exists(keyPath) {
		return certPath, keyPath
	}

	log.Println("Generating new certificates...")

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("failed to generate private key: %s", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
	}

	template := x509.Certificate{
		IsCA: true,

		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, host)
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}

	err = os.MkdirAll(dir, 0700)
	if err != nil {
		log.Fatalf("Failed to write certificates: %s", err)
	}

	certOut, err := os.Create(certPath)
	if err != nil {
		log.Fatalf("failed to open cert.pem for writing: %s", err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("failed to open key.pem for writing: %s", err)
	}

	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	return certPath, keyPath
}

func main() {

	var host = flag.String("host", "localhost", "...")
	var port = flag.Int("port", 9191, "...")
	var cors = flag.Bool("cors", false, "...")
	var refresh = flag.Bool("refresh", false, "...")
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

	re, _ := regexp.Compile(`/(.*)/(\d+)/(\d+)/(\d+).(\w+)$`)

	handler := func(rsp http.ResponseWriter, req *http.Request) {

		url := req.URL
		path := url.Path
		query := url.RawQuery

		if !re.MatchString(path) {
			http.Error(rsp, "404 Not found", http.StatusNotFound)
			return
		}

		local_path := filepath.Join(cache.Path, path)

		if !*refresh {
			_, err := os.Stat(local_path)

			if !os.IsNotExist(err) {

				body, err := ioutil.ReadFile(local_path)

				if err == nil {

					fmt.Println("HIT", local_path)

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

		if query != "" {
			source = source + "?" + query
		}

		fmt.Println("FETCH", source)

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

	// cert, key := generateCert(*host)
	// http.ListenAndServeTLS(*endpoint, cert, key, proxyHandler)

	if err != nil {
		panic(err)
	}

	os.Exit(0)
}
