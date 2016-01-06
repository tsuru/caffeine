package main

import (
	"os"
	"time"

	"github.com/rafaeljusto/redigomock"
	"gopkg.in/check.v1"
)

func (s *Suite) TestRestoreRoute(c *check.C) {
	conn := redigomock.NewConn()
	appAddress := "tsuru-caffeine-proxy.mytsuru.com"
	conn.Command("LTRIM", HIPACHE_PREFIX+appAddress, 0, 0).Expect("OK")
	restoreRoute(appAddress, conn)
}

func (s *Suite) TestWaitBeforeProxy(c *check.C) {
	os.Setenv("WAIT_BEFORE_PROXY", "1")
	start := time.Now()
	waitBeforeProxy()
	runTime := time.Since(start).Seconds()
	c.Assert(runTime >= 1, check.Equals, true)
}

func (s *Suite) TestHipacheRedisAddrWithDefaultConfiguration(c *check.C) {
	os.Setenv("HIPACHE_REDIS_HOST", "")
	os.Setenv("HIPACHE_REDIS_PORT", "")
	c.Assert(hipacheRedisAddr(), check.Equals, "localhost:6379")
}

func (s *Suite) TestHipacheRedisAddrWithHostAndPort(c *check.C) {
	os.Setenv("HIPACHE_REDIS_HOST", "192.168.50.4")
	os.Setenv("HIPACHE_REDIS_PORT", "8989")
	c.Assert(hipacheRedisAddr(), check.Equals, "192.168.50.4:8989")
}

func (s *Suite) TestHipacheRedisMaxConnDefaultConfiguration(c *check.C) {
	os.Setenv("HIPACHE_REDIS_MAX_CONN", "")
	c.Assert(hipacheRedisMaxConn(), check.Equals, 10)
}

func (s *Suite) TestHipacheRedisMaxConnConfiguredValue(c *check.C) {
	os.Setenv("HIPACHE_REDIS_MAX_CONN", "50")
	c.Assert(hipacheRedisMaxConn(), check.Equals, 50)
}
