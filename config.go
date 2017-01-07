package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"pkg.deepin.io/lib/strv"
	"strings"
)

const (
	defaultOutputFile = "/usr/share/multi-display-session/outputs.json"
	defaultAppsFile   = "/usr/share/multi-display-session/autostart"
)

var (
	outputConfigFile = os.Getenv("HOME") + "/.config/multi-display-session/outputs.json"
	appConfigFile    = os.Getenv("HOME") + "/.config/multi-display-session/autostart"
)

func launchApps() {
	apps, err := newAppsFromFile(appConfigFile)
	if err != nil {
		apps, err = newAppsFromFile(defaultAppsFile)
		if err != nil {
			logger.Error("No apps config found:", err)
			return
		}
	}

	for _, app := range apps {
		go runApp(app)
	}
}

func checkOutputConfigValidity(names []string, infos []OutputInfo) bool {
	if len(names) != len(infos) {
		return false
	}

	for _, info := range infos {
		if strings.Contains(strings.ToLower(info.Name), "hdmi") {
			return false
		}
		if !strv.Strv(names).Contains(info.Name) {
			return false
		}
	}
	return true
}

func newOutputInfosFromFile(file string) ([]OutputInfo, error) {
	var infos []OutputInfo
	err := jsonUnmarshalFromFile(file, &infos)
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func newAppsFromFile(file string) ([]string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := strv.Strv(strings.Split(string(content), "\n"))
	lines = lines.FilterEmpty()
	return []string(lines), nil
}

func jsonUnmarshalFromFile(file string, value interface{}) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, value)
}
