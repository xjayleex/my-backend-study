package main


import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"golang.org/x/net/http2"
	"io/ioutil"
	"log"
	"net/http"
)

const url = "https://localhost:8000"
const certFile = "/Users/ijaehyeon/Documents/tls/server.crt"

var (
	httpVersion = flag.Int("version", 2,"HTTP version")
)
func main() {
	flag.Parse()
	client := &http.Client{}
	caCert, err := ioutil.ReadFile(certFile)
	if err != nil {
		log.Fatalf("Reading server certificate : %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config {
		RootCAs: caCertPool,
	}

	switch *httpVersion {
	case 1:
		client.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	case 2:
		client.Transport = &http2.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Failed to Get : %s" , err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed reading response body : %s", err)
	}
	fmt.Printf("Got response %d : %s %s\n",
		resp.StatusCode, resp.Proto, string(body))
}