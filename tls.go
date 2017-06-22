package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"time"
)

var CAcert *x509.Certificate
var CAkey *rsa.PrivateKey

func readPEM(path string) *pem.Block {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic(err)
	}
	block, _ := pem.Decode(data)
	return block
}

func init() {
	var err error
	if CAcert, err = x509.ParseCertificate(readPEM("CA.crt").Bytes); err != nil {
		log.Panic(err)
	}
	if CAkey, err = x509.ParsePKCS1PrivateKey(readPEM("CA.key").Bytes); err != nil {
		log.Panic(err)
	}
}

func getCert(server string) *x509.Certificate {
	conn, err := tls.Dial("tcp", server, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil
	}
	defer conn.Close()

	if err := conn.Handshake(); err != nil {
		return nil
	}
	for _, cert := range conn.ConnectionState().PeerCertificates {
		if len(cert.DNSNames) != 0 {
			return cert
		}
	}
	return nil
}

func generateCert(server string) *tls.Certificate {
	realcert := getCert(server)
	if realcert == nil {
		log.Panic("Cannot get certificate from the remote host")
	}
	serial := realcert.SerialNumber.String()
	realcert.SignatureAlgorithm = x509.SHA256WithRSA
	realcert.SerialNumber.SetInt64(time.Now().UnixNano())

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Panic("Serial: " + serial + err.Error())
	}

	fakecert, err := x509.CreateCertificate(rand.Reader, realcert, CAcert, &priv.PublicKey, CAkey)
	if err != nil {
		log.Panic("Serial: " + serial + err.Error())
	}

	return &tls.Certificate{
		PrivateKey:  priv,
		Certificate: [][]byte{fakecert},
	}
}
