package probit

import (
	"net"
	"net/http"
	"time"
)

type config struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
}

func timeoutDialer(config *config) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (c net.Conn, err error) {
		conn, err := net.DialTimeout(netw, addr, config.ConnectTimeout)
		if err != nil {
			return nil, err
		}

		err = conn.SetDeadline(time.Now().Add(config.ReadWriteTimeout))
		if err != nil {
			return nil, err
		}

		return conn, nil
	}
}

func newTimeoutClient(connectTimeout, readWriteTimeout time.Duration) *http.Client {
	config := &config{
		ConnectTimeout:   connectTimeout,
		ReadWriteTimeout: readWriteTimeout,
	}

	return &http.Client{
		Transport: &http.Transport{
			Dial: timeoutDialer(config),
		},
	}
}
