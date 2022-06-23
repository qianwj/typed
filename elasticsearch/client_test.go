package elasticsearch

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"strings"
	"testing"
)

func TestAA(t *testing.T) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://69.234.217.20:9200"},
		Username:  "qianwj",
		Password:  "123456",
		CACert:    []byte("-----BEGIN CERTIFICATE-----\n    MIIDVjCCAj6gAwIBAgITLu3QwuXyKK6YBRA80SPaE3TJDDANBgkqhkiG9w0BAQsF\n    ADA0MTIwMAYDVQQDEylFbGFzdGljIENlcnRpZmljYXRlIFRvb2wgQXV0b2dlbmVy\n    YXRlZCBDQTAeFw0yMDA2MjYwNTIwMDhaFw0yMzA2MjYwNTIwMDhaMA8xDTALBgNV\n    BAMTBGVzMDEwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDDhWGSUvTC\n    P/Qj0y73i7hVEtZvWuD4XFMQutWf0j3p/Yx3Aau8OUPkmyR/X+Xzh6oXX378y9I/\n    AJIhrwKTyq5XMQ4hJWROehTw/AsbBNhe0ermM2VEG8XQFybZstH8xCxPexkLvxy4\n    hO4eI69AH77RO5nH39jEdPQ1fuJpJ4MqDRM3nUpniRUox9uN1L5b5gQDYr+sgo+o\n    6XdVnU85SYeJs1K++osai78U0oF9w8Rs7m0asSiwo91SoU4PqrAHl/d4UzFvR3mb\n    +4FYBS+vIYl1wfyhOx9TEEgWecb3haVYslVwYffEPldGn1UOFq+CKpzvc4pBAjuj\n    acas3eOX8ZNdAgMBAAGjgYUwgYIwHQYDVR0OBBYEFNI418iE2grM+Af/4GoUMS4R\n    HJ41MB8GA1UdIwQYMBaAFJHzG4Tq2xms8ePU1+avug5BF+ySMDUGA1UdEQQuMCyC\n    CWxvY2FsaG9zdIITZXMuaW50ZXJuYWwuanV6aWJvdIcEfwAAAYIEZXMwMTAJBgNV\n    HRMEAjAAMA0GCSqGSIb3DQEBCwUAA4IBAQBWm/76OROWDx7Xql1kC20rjLnHhK8g\n    ddFVfkXk1WGkpldlO8CK+MFbPLFs07kBKq9XwwozLIxDJqOVwQmaxAHoZvidaa+U\n    TdoVEYnca4uL6VCOH5CsfFLuXROXgajKOkxoWkMjLuNwKxsvqRr5aVlDBnRs57Ck\n    1kfOPrn9uDrBZzruInARrt5JmPorXQwvg6nUDyPeIHJM+iUcjdXNU+W/F12qiN07\n    C9FOTO5DbAZeX4uBExek6LDQ8JsD92qOAaL+51kpFFymvHntQhI4xshcl8CRPXSj\n    hSlobegL2HLhdMsKZURtwvZNn2aWdwMxQxBuxR4xhGkq0sWWmS9xjYNk\n    -----END CERTIFICATE-----\n"),
	})
	var (
		r map[string]interface{}
	)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		t.FailNow()
	}
	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	log.Println(strings.Repeat("~", 37))
}
