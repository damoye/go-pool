package pool

import (
	"io"
	"log"
	"net"
	"testing"
)

var (
	dialer = func() (net.Conn, error) { return net.Dial("tcp", "127.0.0.1:8080") }
)

func init() {
	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer l.Close()
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go func() {
				_, err := io.Copy(conn, conn)
				if err != nil {
					log.Fatal(err)
				}
				conn.Close()
			}()
		}
	}()
}

func TestNew(t *testing.T) {
	_, err := New(30, dialer)
	if err != nil {
		t.Error(err)
	}
}

func TestNewBadPoolSize(t *testing.T) {
	_, err := New(0, dialer)
	if err == nil {
		t.Error("get: nil, want: error")
	}
}

func TestGetBadDailer(t *testing.T) {
	p, err := New(30, func() (net.Conn, error) { return net.Dial("tcp", "127.0.0.1:8888") })
	if err != nil {
		t.Error(err)
	}
	_, err = p.Get()
	if err == nil {
		t.Error("get: nil, want: error")
	}
}

func TestGetPut(t *testing.T) {
	p, err := New(30, dialer)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 2; i++ {
		c, err := p.Get()
		if err != nil {
			t.Error(err)
		}
		if err := p.Put(c); err != nil {
			t.Error(err)
		}
	}
}

func TestPutFullPool(t *testing.T) {
	p, err := New(1, dialer)
	if err != nil {
		t.Error(err)
	}
	for i := 0; i < 2; i++ {
		c, err := dialer()
		if err != nil {
			t.Error(err)
		}
		if err := p.Put(c); err != nil {
			t.Error(err)
		}
	}
}
