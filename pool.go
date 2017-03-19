package pool

import (
	"errors"
	"net"
)

// Pool is a pool for net.Conn
type Pool struct {
	conns  chan net.Conn
	dialer func() (net.Conn, error)
}

// NewPool returns a new net.Conn pool
func NewPool(poolSize int, dialer func() (net.Conn, error)) (*Pool, error) {
	if poolSize <= 0 {
		return nil, errors.New("invalid poolSize")
	}
	pool := Pool{make(chan net.Conn, poolSize), dialer}
	conn, err := dialer()
	if err != nil {
		return nil, err
	}
	pool.conns <- conn
	return &pool, nil
}

// Get removes a net.Conn from the Pool, and returns it to the caller.
// If there is no net.Conn in the pool, Get returns the result of calling dialer
func (p *Pool) Get() (net.Conn, error) {
	select {
	case conn := <-p.conns:
		return conn, nil
	default:
		conn, err := p.dialer()
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}

// Put adds conn to the pool.
func (p *Pool) Put(conn net.Conn) error {
	select {
	case p.conns <- conn:
		return nil
	default:
		return conn.Close()
	}
}
