package acme

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"
)

type Domain struct {
	Main string `json:"main"`
}

type Certificate struct {
	Domain      Domain `json:"domain"`
	Certificate string `json:"certificate"`
	Key         string `json:"key"`
}

type Account struct {
	Email      string `json:"email"`
	PrivateKey string `json:"PrivateKey"`
}

type letsencrypt struct {
	Account      Account       `json:"Account"`
	Certificates []Certificate `json:"Certificates"`
}

type ACME struct {
	Letsencrypt letsencrypt `json:"letsencrypt"`
}

func Load(path string) (tls.Certificate, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return tls.Certificate{}, err
	}
	var acme ACME
	err = json.Unmarshal(data, &acme)
	if err != nil {
		return tls.Certificate{}, err
	}

	for _, cert := range acme.Letsencrypt.Certificates {
		if strings.HasPrefix(cert.Domain.Main, "mqtt.") {
			c, err := base64.StdEncoding.DecodeString(cert.Certificate)
			if err != nil {
				return tls.Certificate{}, err
			}
			k, err := base64.StdEncoding.DecodeString(cert.Key)
			if err != nil {
				return tls.Certificate{}, err
			}
			return tls.X509KeyPair(c, k)
		}
	}

	return tls.Certificate{}, errors.New("mqtt certificate not found")
}
