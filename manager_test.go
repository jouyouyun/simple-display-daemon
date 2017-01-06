package main

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestModFilter(t *testing.T) {
	Convey("Test mod filter", t, func() {
		So(filterInvalidMod("control-mod1-mod2"), ShouldEqual, "control-mod1")
	})
}

func TestOutputConfig(t *testing.T) {
	Convey("Test output config parse", t, func() {
		infos, err := newOutputInfosFromFile("testdata/outputs.json")
		So(err, ShouldBeNil)
		So(len(infos), ShouldEqual, 3)
		var target = []OutputInfo{
			{
				Name:   "eDP-1",
				X:      0,
				Y:      0,
				Width:  1024,
				Height: 768,
			},
			{
				Name:   "VGA-1",
				X:      1024,
				Y:      0,
				Width:  1024,
				Height: 768,
			},
			{
				Name:   "HDMI-1",
				X:      2048,
				Y:      0,
				Width:  1024,
				Height: 768,
			},
		}
		data1, _ := json.Marshal(infos)
		data2, _ := json.Marshal(target)
		So(string(data1), ShouldEqual, string(data2))
	})
}

func TestAppConfig(t *testing.T) {
	Convey("Test app config parse", t, func() {
		infos, err := newOutputInfosFromFile("testdata/apps.json")
		So(err, ShouldBeNil)
		So(len(infos), ShouldEqual, 3)
		var target = []OutputInfo{
			{
				Name: "app1",
				X:    0,
				Y:    0,
			},
			{
				Name: "app2",
				X:    1024,
				Y:    0,
			},
			{
				Name: "app3",
				X:    2048,
				Y:    0,
			},
		}
		data1, _ := json.Marshal(infos)
		data2, _ := json.Marshal(target)
		So(string(data1), ShouldEqual, string(data2))
	})
}
