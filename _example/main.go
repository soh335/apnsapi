package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/soh335/apnsapi"
	"golang.org/x/crypto/pkcs12"
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
)

func main() {
	flag.Parse()
	if err := _main(); err != nil {
		log.Fatal(err)
	}
}

func _main() error {
	cert, err := loadCertificate()
	if err != nil {
		return err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	tr := &http.Transport{
		TLSClientConfig: config,
	}
	c := &http.Client{
		Transport: tr,
	}

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

	header := &apnsapi.Header{ApnsTopic: *topic}

	_, err = client.Do(*token, header, b.Bytes())
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
