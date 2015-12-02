package main

import (
	"github.com/rafaeljusto/redigomock"
	"gopkg.in/check.v1"
)

func (s *Suite) TestRestoreRoute(c *check.C) {
	conn := redigomock.NewConn()
	host := "tsuru-caffeine-proxy.mytsuru.com"
	conn.Command("LTRIM", HIPACHE_PREFIX+host, 0, 0).Expect("OK")
	restoreRoute(host, conn)
}
