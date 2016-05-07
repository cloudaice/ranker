package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

var (
	httpClient       *http.Client
	client           *github.Client
	LimitQueue       chan struct{}
	LimitConns       = 100
	ConnTimeout      = 2 * time.Second
	ReadWriteTimeout = 5 * time.Second
)

var (
	LimitErr = fmt.Errorf("Limit Conns Error")

	AccessToken string
)

type YmlConf struct {
	AccessToken string `yaml:"AccessToken"`
}

func InitConfig() {
	data, err := ioutil.ReadFile("ranker.yaml")
	if err != nil {
		panic(err)
	}
	var c = &YmlConf{}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		panic(err)
	}
	AccessToken = c.AccessToken
}

func init() {
	InitConfig()
	LimitQueue = make(chan struct{}, LimitConns)
	httpClient = &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				select {
				case LimitQueue <- struct{}{}:
				default:
					return nil, LimitErr
				}
				conn, err := net.DialTimeout(network, addr, ConnTimeout)
				if err != nil {
					<-LimitQueue
					return nil, err
				}
				return &LimitConn{conn}, nil
			},
		},
		Timeout: ReadWriteTimeout,
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: AccessToken})
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
}

type LimitConn struct {
	net.Conn
}

func (lc *LimitConn) Close() error {
	<-LimitQueue
	return lc.Conn.Close()
}
