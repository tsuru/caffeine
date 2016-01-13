package main

import (
	"os"
	"time"

	"gopkg.in/check.v1"
)

func (s *Suite) TestWaitBeforeProxy(c *check.C) {
	did_wait := false
	os.Setenv("WAIT_BEFORE_PROXY", "10")
	waitBeforeProxy(func(duration time.Duration) {
		expected := time.Duration(10) * time.Second
		c.Assert(duration, check.Equals, expected)
		did_wait = true
	})
	c.Assert(did_wait, check.Equals, true)
}
