package main

import (
	"io/ioutil"
	"strings"
)

var prevConnected []OutputInfo

func (m *Manager) ListConnectedOutput() []OutputInfo {
	m.outputLocker.Lock()
	defer m.outputLocker.Unlock()
	connected := m.outputInfos.ListValidOutputs().ListConnectionOutputs()
	var infos []OutputInfo
	for _, output := range connected {
		if !canReadEDID(output.Name) || output.Crtc.Width == 0 || output.Crtc.Height == 0 {
			continue
		}
		infos = append(infos, OutputInfo{
			Name:   output.Name,
			X:      output.Crtc.X,
			Y:      output.Crtc.Y,
			Width:  output.Crtc.Width,
			Height: output.Crtc.Height,
		})
	}
	prevConnected = infos
	return infos
}

const drmCard0Path = "/sys/class/drm/"

func canReadEDID(name string) bool {
	card := getCardByName(name)
	if len(card) == 0 {
		return false
	}
	var file = drmCard0Path + card + "/edid"
	content, err := ioutil.ReadFile(file)
	if err != nil {
		logger.Errorf("Read edid file '%s' failed: %v", file, err)
		return false
	}
	if len(content) == 0 {
		return false
	}
	return true
}

func getCardByName(name string) string {
	cards := getAllCards()
	if len(cards) == 0 {
		return ""
	}

	if strings.Contains(name, "-") {
		array := strings.Split(name, "-")
		name = strings.Join(array, "")
	}
	for _, card := range cards {
		array := strings.Split(card, "-")
		tmp := strings.Join(array[1:], "")
		if tmp == name {
			return card
		}
	}
	return ""
}

func getAllCards() []string {
	finfos, err := ioutil.ReadDir(drmCard0Path)
	if err != nil {
		logger.Error("Read card0 list failed:", err)
		return nil
	}

	var cards []string
	for _, finfo := range finfos {
		if !strings.Contains(finfo.Name(), "card0") {
			continue
		}

		if finfo.Name() == "card0" {
			continue
		}
		cards = append(cards, finfo.Name())
	}
	return cards
}
