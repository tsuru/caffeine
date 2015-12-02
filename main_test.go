package main

import (
	"net/http"

	"github.com/rafaeljusto/redigomock"
	"gopkg.in/check.v1"
)

func (s *Suite) TestRestoreRoute(c *check.C) {
	conn := redigomock.NewConn()
	req, _ := http.NewRequest("GET", "http://tsuru-caffeine-proxy.mytsuru.com", nil)
	conn.Command("LTRIM", HIPACHE_PREFIX+req.Host, 0, 0).Expect("OK")
	restoreRoute(req, conn)

}
