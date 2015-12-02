package main

import (
	"net/http"
	"net/http/httptest"
	"os"

	"gopkg.in/check.v1"
)

func (s *Suite) TestStartApp(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Method, check.Equals, "POST")
		c.Assert(r.URL.String(), check.Equals, "/apps/myapp/start")
		c.Assert(r.Header.Get("Authorization"), check.Equals, "bearer 123")
	}))
	defer ts.Close()
	os.Setenv("TSURU_HOST", ts.URL)
	os.Setenv("TOKEN", "123")

	startApp("myapp.mytsuru.com")
}

func (s *Suite) TestAppNameCaffeine(c *check.C) {
	app, err := appName("tsuru-caffeine-proxy.mytsuru.com")
	c.Assert(err, check.ErrorMatches, "invalid app name")
	c.Assert(app, check.Equals, "")
}

func (s *Suite) TestAppNameApp(c *check.C) {
	app, err := appName("myapp.mytsuru.com")
	c.Assert(err, check.IsNil)
	c.Assert(app, check.Equals, "myapp")
}
