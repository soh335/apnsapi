package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/soh335/apnsapi"
	"golang.org/x/crypto/pkcs12"
	"golang.org/x/net/http2"
)

var (
	cerFile     = flag.String("cer", "", "cer")
	keyFile     = flag.String("key", "", "key")
	p12File     = flag.String("p12", "", "p12")
	p12Password = flag.String("password", "", "password of p12")
	production  = flag.Bool("production", false, "connect to production server")
	msg         = flag.String("msg", "hi", "msg")
	token       = flag.String("token", "", "token")
	topic       = flag.String("topic", "", "topic")
	kid         = flag.String("kid", "", "key identifier")
	teamId      = flag.String("teamId", "", "team id")
)

func main() {
	flag.Parse()
	var client *http.Client
	var header *apnsapi.Header
	if *kid != "" && *teamId != "" {
		jwtToken, err := createJwtToken()
		if err != nil {
			log.Fatal(err)
		}
		tr := &http.Transport{}

		if err := http2.ConfigureTransport(tr); err != nil {
			log.Fatal(err)
		}

		client = &http.Client{
			Transport: tr,
		}

		header = &apnsapi.Header{ApnsTopic: *topic, Authorization: "bearer " + jwtToken}
	} else {
		cert, err := loadCertificate()
		if err != nil {
			log.Fatal(err)
		}
		config := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		tr := &http.Transport{
			TLSClientConfig: config,
		}
		if err := http2.ConfigureTransport(tr); err != nil {
			log.Fatal(err)
		}

		client = &http.Client{
			Transport: tr,
		}

		header = &apnsapi.Header{ApnsTopic: *topic}
	}
	if err := _main(client, header); err != nil {
		log.Fatal(err)
	}
}

func _main(c *http.Client, header *apnsapi.Header) error {
	var host string
	if *production {
		host = apnsapi.ProductionServer
	} else {
		host = apnsapi.DevelopmentServer
	}

	client := apnsapi.NewClient(host, c)

	var b bytes.Buffer

	p := map[string]interface{}{
		"aps": map[string]interface{}{
			"alert": *msg,
		},
	}

	if err := json.NewEncoder(&b).Encode(p); err != nil {
		return err
	}

	_, err := client.Do(*token, header, b.Bytes())
	if err != nil {
		switch err.(type) {
		case *apnsapi.ErrorResponse:
			log.Println("error response:", err.(*apnsapi.ErrorResponse).Reason)
		default:
			return err
		}
	}

	return nil
}

func loadCertificate() (tls.Certificate, error) {
	if *p12File != "" {
		cert := tls.Certificate{}
		b, err := ioutil.ReadFile(*p12File)
		if err != nil {
			return cert, err
		}
		key, cer, err := pkcs12.Decode(b, *p12Password)
		if err != nil {
			return cert, err
		}
		cert.PrivateKey = key
		cert.Certificate = [][]byte{cer.Raw}
		return cert, nil
	} else {
		return tls.LoadX509KeyPair(*cerFile, *keyFile)
	}
}

func createJwtToken() (string, error) {
	data, err := ioutil.ReadFile(*keyFile)
	if err != nil {
		return "", err
	}
	block, _ := pem.Decode(data)
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	return apnsapi.CreateToken(key.(*ecdsa.PrivateKey), *kid, *teamId)
}
