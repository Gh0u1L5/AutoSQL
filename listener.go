package main

import (
	"io"
	"net"
)

type listener struct {
	count int
	port  chan net.Conn
	done  chan bool
}

func (ln *listener) Accept() (net.Conn, error) {
	ln.count += 1
	defer func() { ln.count -= 1 }()

	select {
	case conn := <-ln.port:
		return conn, nil
	case <-ln.done:
		return nil, io.EOF
	}
}

func (ln *listener) Close() error {
	for i := 0; i < ln.count; i++ {
		ln.done <- true
	}
	return nil
}

func (ln *listener) Addr() net.Addr {
	return nil
}
