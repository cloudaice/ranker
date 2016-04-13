package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

var (
	client           *http.Client
	LimitQueue       chan struct{}
	LimitConns       = 100
	ConnTimeout      = 2 * time.Second
	ReadWriteTimeout = 5 * time.Second
)

var (
	LimitErr = fmt.Errorf("Limit Conns Error")
)

func init() {
	LimitQueue = make(chan struct{}, LimitConns)
	client = &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				select {
				case LimitQueue <- struct{}{}:
				default:
					return nil, LimitErr
				}
				conn, err := net.DialTimeout(network, addr, time.Second)
				if err != nil {
					<-LimitQueue
					return nil, err
				}
				return &LimitConn{conn}, nil
			},
		},
		Timeout: ReadWriteTimeout,
	}
}

type LimitConn struct {
	net.Conn
}

func (lc *LimitConn) Close() error {
	<-LimitQueue
	return lc.Conn.Close()
}
