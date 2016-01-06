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

func startApp(hostname string) {
	appName, err := appName(hostname)
	if err != nil {
		log.Println(err)
		return
	}

	startAppURL := fmt.Sprintf("%s/apps/%s/start", getConfig("TSURU_HOST"), appName)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", startAppURL, nil)
	req.Header.Add("Authorization", authToken())
	resp, _ := client.Do(req)
	if resp.StatusCode != 200 {
		log.Printf("Error trying to start app %s\n", appName)
		return
	}

	log.Printf("App %s started\n", appName)
}

func appName(hostname string) (string, error) {
	app, err := getApp(hostname)
	if err != nil {
		return "", err
	}
	if app.Name == getConfig("TSURU_APP_PROXY") {
		return "", fmt.Errorf("App %s can't be started by itself", app.Name)
	}

	return app.Name, nil
}

func getApp(hostname string) (*App, error) {
	apps, err := listApps()
	if err != nil {
		return nil, err
	}

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

func authToken() string {
	return fmt.Sprintf("bearer %s", getConfig("TOKEN"))
}

func listApps() ([]App, error) {
	listAppsURL := fmt.Sprintf("%s/apps/", getConfig("TSURU_HOST"))

	client := &http.Client{}
	req, _ := http.NewRequest("GET", listAppsURL, nil)
	req.Header.Add("Authorization", authToken())
	resp, _ := client.Do(req)
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
