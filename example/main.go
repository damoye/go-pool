package main

import (
	"fmt"
	"log"
	"net"

	"github.com/damoye/gopool"
)

func main() {
	dialer := func() (net.Conn, error) { return net.Dial("tcp", "127.0.0.1:80") }
	p, err := pool.NewPool(30, dialer)
	if err != nil {
		log.Fatal(err)
	}
	conn, err := p.Get()
	if err != nil {
		log.Fatal(err)
	}
	// Do something with comm
	_, err = fmt.Fprint(conn, "hello world")
	if err != nil {
		conn.Close()
	} else {
		p.Put(conn)
	}
}
