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
			c.Assert(r.URL.String(), check.Equals, "/apps")
			jsonData, _ := json.Marshal([]App{
				App{Name: "myapp", Ip: "myapp.mytsuru.com"},
			})
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
		} else {
			c.Assert(r.Method, check.Equals, "POST")
			c.Assert(r.URL.String(), check.Equals, "/apps/myapp/start")
		}
	}))
	defer ts.Close()
	os.Setenv("TSURU_HOST", ts.URL)
	os.Setenv("TSURU_TOKEN", "123")

	startApp(App{Name: "myapp"})
}

func (s *Suite) TestGetAppIsProxy(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Method, check.Equals, "GET")
		c.Assert(r.URL.String(), check.Equals, "/apps")
		c.Assert(r.Header.Get("Authorization"), check.Equals, "bearer 123")

		jsonData, _ := json.Marshal([]App{
			App{Name: "tsuru-caffeine-proxy", Ip: "", Cname: []string{"proxy.mytsuru.com"}},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}))
	defer ts.Close()
	os.Setenv("TSURU_HOST", ts.URL)
	os.Setenv("TSURU_TOKEN", "123")
	os.Setenv("TSURU_APP_PROXY", "tsuru-caffeine-proxy")

	app, err := getApp("proxy.mytsuru.com")
	c.Assert(err, check.ErrorMatches, "App tsuru-caffeine-proxy is the proxy itself")
	c.Assert(app.Name, check.Equals, "tsuru-caffeine-proxy")
}

func (s *Suite) TestGetAppFoundByIp(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Method, check.Equals, "GET")
		c.Assert(r.URL.String(), check.Equals, "/apps")
		c.Assert(r.Header.Get("Authorization"), check.Equals, "bearer 123")

		jsonData, _ := json.Marshal([]App{
			App{Name: "myapp0", Ip: "myapp", Cname: []string{}},
			App{Name: "myapp-name", Ip: "myapp.mytsuru.com", Cname: []string{}},
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}))
	defer ts.Close()
	os.Setenv("TSURU_HOST", ts.URL)
	os.Setenv("TSURU_TOKEN", "123")

	app, err := getApp("myapp.mytsuru.com")
	c.Assert(err, check.IsNil)
	c.Assert(app.Name, check.Equals, "myapp-name")
}

func (s *Suite) TestGetAppFoundByCname(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Method, check.Equals, "GET")
		c.Assert(r.URL.String(), check.Equals, "/apps")
		c.Assert(r.Header.Get("Authorization"), check.Equals, "bearer 123")

		jsonData, _ := json.Marshal([]App{
			App{Name: "app1", Ip: "myapp-cname", Cname: []string{"cname1.example.com"}},
			App{Name: "real-app-name", Ip: "", Cname: []string{"cname2.example.com", "myapp-cname.mytsuru.com"}},
			App{Name: "app2", Ip: "app2", Cname: []string{"app2.mytsuru.com"}},
		})

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}))
	defer ts.Close()
	os.Setenv("TSURU_HOST", ts.URL)
	os.Setenv("TSURU_TOKEN", "123")

	app, err := getApp("myapp-cname.mytsuru.com")
	c.Assert(err, check.IsNil)
	c.Assert(app.Name, check.Equals, "real-app-name")
}

func (s *Suite) TestGetAppNotFound(c *check.C) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Assert(r.Method, check.Equals, "GET")
		c.Assert(r.URL.String(), check.Equals, "/apps")
		c.Assert(r.Header.Get("Authorization"), check.Equals, "bearer 123")

		jsonData, _ := json.Marshal([]App{
			App{Name: "app-name", Ip: "app-ip", Cname: []string{"cname1.example.com", "cname2.example.com"}},
		})

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}))
	defer ts.Close()
	os.Setenv("TSURU_HOST", ts.URL)
	os.Setenv("TSURU_TOKEN", "123")

	app, err := getApp("myapp.mytsuru.com")
	c.Assert(err, check.ErrorMatches, "App myapp.mytsuru.com not found")
	c.Assert(app, check.IsNil)
}
