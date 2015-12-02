package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type AppSuite struct{}

var _ = check.Suite(&AppSuite{})

func (s *AppSuite) TestStartApp(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Method, check.Equals, "POST")
		c.Assert(r.URL.String(), check.Equals, "/apps/myapp/start")
		c.Assert(r.Header.Get("Authorization"), check.Equals, "bearer 123")
	}))
	defer ts.Close()
	os.Setenv("TSURU_HOST", ts.URL)
	os.Setenv("TOKEN", "123")

	req, err := http.NewRequest("GET", "http://myapp.mytsuru.com", nil)
	c.Assert(err, check.IsNil)
	startApp(req)
}

func (s *AppSuite) TestAppNameCaffeine(c *check.C) {
	req, err := http.NewRequest("GET", "http://tsuru-caffeine-proxy.mytsuru.com", nil)
	c.Assert(err, check.IsNil)
	req.Header.Add("Host", "tsuru-caffeine-proxy.mytsuru.com")
	app, err := appName(req)
	c.Assert(err, check.ErrorMatches, "invalid app name")
	c.Assert(app, check.Equals, "")
}

func (s *AppSuite) TestAppNameApp(c *check.C) {
	req, err := http.NewRequest("GET", "http://myapp.mytsuru.com", nil)
	c.Assert(err, check.IsNil)
	req.Header.Add("Host", "myapp.mytsuru.com")
	app, err := appName(req)
	c.Assert(err, check.IsNil)
	c.Assert(app, check.Equals, "myapp")
}
