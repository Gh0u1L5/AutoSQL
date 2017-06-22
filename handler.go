package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func hijack(w http.ResponseWriter) (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("Failed to cast writer to hijacker.")
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		return nil, nil, err
	}
	if len := bufrw.Reader.Buffered(); len > 0 {
		buf, _ := bufrw.Peek(len)
		log.Printf("Non empty bufrw: \"%s\"", string(buf))
	}
	return conn, bufrw, nil
}

func handleRegularRequest(w http.ResponseWriter, req *http.Request) {
	if req.URL.Scheme == "" {
		req.URL.Scheme = "https"
		req.URL.Host = req.Host
		go scanURL("https://" + req.Host + req.RequestURI)
	} else {
		go scanURL(req.RequestURI)
	}
	req.RequestURI = ""

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Print(err)
		return
	}

	for key, values := range resp.Header {
		w.Header().Set(key, values[0])
		for _, value := range values[1:] {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func handleConnectRequest(w http.ResponseWriter, req *http.Request) {
	ln := &listener{0, make(chan net.Conn), make(chan bool)}
	srv := &http.Server{
		Handler: http.HandlerFunc(handleRequest),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*generateCert(req.Host)},
		},
	}
	go srv.Serve(tls.NewListener(ln, srv.TLSConfig))

	conn, _, err := hijack(w)
	if err != nil {
		log.Print(err)
		return
	}
	conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	ln.port <- conn
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodConnect {
		handleConnectRequest(w, req)
	} else {
		handleRegularRequest(w, req)
	}
}
