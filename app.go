package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type App struct {
	Name  string
	Ip    string
	Cname []string
}

func startApp(app App) {
	host, _ := getConfig("TSURU_HOST")
	startAppURL := fmt.Sprintf("%s/apps/%s/start", host, app.Name)

	req, err := http.NewRequest("POST", startAppURL, nil)
	if err != nil {
		log.Printf("Error trying to start app %s: %s\n", app.Name, err.Error())
	}

	authToken, err := authToken()
	if err != nil {
		log.Printf("Error trying to start app %s: %s\n", app.Name, err.Error())
		return
	}

	req.Header.Add("Authorization", authToken)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		log.Printf("Error trying to start app %s\n", app.Name)
		return
	}

	log.Printf("App %s started\n", app.Name)
}

func getApp(hostname string) (*App, error) {
	if len(hostname) == 0 {
		return nil, fmt.Errorf("Empty hostname")
	}
	apps, err := listApps()
	if err != nil {
		return nil, err
	}

	return filterAppByHostname(hostname, apps)
}

func filterAppByHostname(hostname string, apps []App) (*App, error) {
	for _, app := range apps {
		if hostname == app.Ip {
			return &app, nil
		}
		for _, cname := range app.Cname {
			if hostname == cname {
				return &app, nil
			}
		}
	}

	return nil, fmt.Errorf("App %s not found", hostname)
}

func authToken() (string, error) {
	token, err := getConfig("TSURU_TOKEN")
	if err != nil {
		return token, err
	}

	return fmt.Sprintf("bearer %s", token), nil
}

func listApps() ([]App, error) {
	host, _ := getConfig("TSURU_HOST")
	listAppsURL := fmt.Sprintf("%s/apps?status=asleep", host)

	req, _ := http.NewRequest("GET", listAppsURL, nil)

	authToken, err := authToken()
	if err != nil {
		return nil, fmt.Errorf("Error trying to get apps info: %s", err.Error())
	}

	req.Header.Add("Authorization", authToken)
	resp, _ := http.DefaultClient.Do(req)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error trying to get apps info: HTTP %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error trying to get apps info: %s", err.Error())
	}

	var apps []App
	err = json.Unmarshal(body, &apps)
	if err != nil {
		return nil, fmt.Errorf("Error trying to get apps info: %s", err.Error())
	}

	return apps, nil
}
