package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"

	"gopkg.in/check.v1"
)

func (s *Suite) TestStartApp(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Header.Get("Authorization"), check.Equals, "bearer 123")

		if r.Method == "GET" {
			c.Assert(r.URL.String(), check.Equals, "/apps/?name=myapp")
			json, _ := json.Marshal([]map[string]string{map[string]string{"name": "myapp"}})
			w.Header().Set("Content-Type", "application/json")
			w.Write(json)
		} else {
			c.Assert(r.Method, check.Equals, "POST")
			c.Assert(r.URL.String(), check.Equals, "/apps/myapp/start")
		}
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

func (s *Suite) TestAppNameFound(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Method, check.Equals, "GET")
		c.Assert(r.URL.String(), check.Equals, "/apps/?name=myapp")
		c.Assert(r.Header.Get("Authorization"), check.Equals, "bearer 123")

		json, _ := json.Marshal([]map[string]string{map[string]string{"name": "myapp"}})
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}))
	defer ts.Close()
	os.Setenv("TSURU_HOST", ts.URL)
	os.Setenv("TOKEN", "123")

	app, err := appName("myapp.mytsuru.com")
	c.Assert(err, check.IsNil)
	c.Assert(app, check.Equals, "myapp")
}

func (s *Suite) TestAppNameFoundByCname(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Method, check.Equals, "GET")
		c.Assert(r.Header.Get("Authorization"), check.Equals, "bearer 123")

		var jsonData []byte
		if r.URL.String() == "/apps/?name=myapp-cname" {
			jsonData, _ = json.Marshal([]map[string]string{})
		} else {
			c.Assert(r.URL.String(), check.Equals, "/apps/")
			jsonData, _ = json.Marshal([]map[string]string{map[string]string{"name": "real-app-name", "cname": "myapp-cname"}})
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}))
	defer ts.Close()
	os.Setenv("TSURU_HOST", ts.URL)
	os.Setenv("TOKEN", "123")

	app, err := appName("myapp-cname.mytsuru.com")
	c.Assert(err, check.IsNil)
	c.Assert(app, check.Equals, "real-app-name")
}

func (s *Suite) TestAppNameNotFound(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Method, check.Equals, "GET")
		c.Assert(r.Header.Get("Authorization"), check.Equals, "bearer 123")

		var jsonData []byte
		if r.URL.String() == "/apps/?name=myapp" {
			jsonData, _ = json.Marshal([]map[string]string{})
		} else {
			c.Assert(r.URL.String(), check.Equals, "/apps/")
			jsonData, _ = json.Marshal([]map[string]string{})
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}))
	defer ts.Close()
	os.Setenv("TSURU_HOST", ts.URL)
	os.Setenv("TOKEN", "123")

	app, err := appName("myapp.mytsuru.com")
	c.Assert(err, check.ErrorMatches, "App not found")
	c.Assert(app, check.Equals, "")
}
