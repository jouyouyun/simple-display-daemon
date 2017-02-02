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

const drmClassPath = "/sys/class/drm/"

func canReadEDID(name string) bool {
	list := getCardByName(name)
	logger.Debug("[canReadEDID] Find card by name:", name, list)
	if len(list) == 0 {
		return false
	}
	for _, card := range list {
		var file = drmClassPath + card + "/edid"
		content, err := ioutil.ReadFile(file)
		if err != nil {
			logger.Errorf("Read edid file '%s' failed: %v", file, err)
			continue
		}
		if len(content) != 0 {
			return true
		}
	}
	return false
}

func getCardByName(name string) []string {
	cards := getAllCards()
	if len(cards) == 0 {
		return nil
	}
	logger.Debug("[getCardByName] all cards:", cards)

	if strings.Contains(name, "-") {
		array := strings.Split(name, "-")
		name = strings.Join(array, "")
	}
	var list []string
	for _, card := range cards {
		array := strings.Split(card, "-")
		tmp := strings.Join(array[1:], "")
		if tmp == name {
			list = append(list, card)
		}
	}
	return list
}

func getAllCards() []string {
	finfos, err := ioutil.ReadDir(drmClassPath)
	if err != nil {
		logger.Error("Read card list failed:", err)
		return nil
	}

	var cards []string
	for _, finfo := range finfos {
		if !strings.Contains(finfo.Name(), "card") {
			continue
		}

		if finfo.Name() == "card0" {
			continue
		}
		cards = append(cards, finfo.Name())
	}
	return cards
}
