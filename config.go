package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"pkg.deepin.io/lib/strv"
	"strings"
)

var (
	outputConfigFile = os.Getenv("HOME") + "/.config/deepin/display-manager/outputs.json"
	appConfigFile    = os.Getenv("HOME") + "/.config/deepin/display-manager/apps.json"
)

type AppInfo struct {
	Name string
	X    int16
	Y    int16
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

func newAppInfosFromFile(file string) ([]AppInfo, error) {
	var infos []AppInfo
	err := jsonUnmarshalFromFile(file, &infos)
	if err != nil {
		return nil, err
	}
	return infos, nil
}

func jsonUnmarshalFromFile(file string, value interface{}) error {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, value)
}
