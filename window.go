package main

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
)

type WindowInfo struct {
	Id   xproto.Window
	Name string
	// Exec string
	Pid uint32
}
type WindowInfos []*WindowInfo

func (m *Manager) MoveWindow(name string, x, y int16) error {
	logger.Debug("MoveWindow:", name, x, y)
	info, err := m.getWindowByName(name)
	if err != nil {
		logger.Error("Failed to get window:", err)
		return err
	}

	logger.Debugf("Window %s id: %v", name, info.Id)
	err = ewmh.MoveWindow(m.xu, info.Id, int(x), int(y))
	if err != nil {
		logger.Errorf("Failed to move '%s' to (%d, %d)", name, x, y)
	}
	return err
}

func (m *Manager) ResizeWindow(name string, w, h int16) error {
	logger.Debug("ResizeWindow:", name, w, h)
	info, err := m.getWindowByName(name)
	if err != nil {
		logger.Error("Failed to get window:", err)
		return err
	}

	logger.Debugf("Window %s id: %v", name, info.Id)
	err = ewmh.ResizeWindow(m.xu, info.Id, int(w), int(h))
	if err != nil {
		logger.Errorf("Failed to resize '%s' to (%d, %d)", name, w, h)
	}
	return err
}

func (m *Manager) getWindowByName(name string) (*WindowInfo, error) {
	infos, err := m.getWindowInfos()
	if err != nil {
		return nil, err
	}

	info := infos.GetByName(name)
	if info == nil {
		return nil, fmt.Errorf("Not found window by name '%s'", name)
	}
	return info, nil
}

func (m *Manager) ListWindow() (string, error) {
	infos, err := m.getWindowInfos()
	if err != nil {
		return "", err
	}

	data, err := json.Marshal(infos)
	if err != nil {
		logger.Error("Failed to marshal window infos:", err)
		return "", err
	}
	return string(data), nil
}

func (m *Manager) getWindowInfos() (WindowInfos, error) {
	list, err := m.getWindowList()
	if err != nil {
		return nil, err
	}

	var infos WindowInfos
	for _, id := range list {
		info, err := m.getWindowInfo(id)
		if info == nil || err != nil {
			continue
		}
		infos = append(infos, info)
	}
	return infos, nil
}

func (m *Manager) getWindowList() ([]xproto.Window, error) {
	list, err := ewmh.ClientListGet(m.xu)
	if err == nil {
		return list, nil
	}
	logger.Debug("Failed to get client wm list:", err)
	reply, err := xproto.QueryTree(m.conn, m.root).Reply()
	if err != nil {
		logger.Error("Failed to qeury tree:", err)
		return nil, err
	}
	return reply.Children, nil

}

func (m *Manager) getWindowInfo(wid xproto.Window) (*WindowInfo, error) {
	attrs, err := xproto.GetWindowAttributes(m.conn, wid).Reply()
	if err != nil {
		return nil, err
	}

	if attrs.MapState != xproto.MapStateViewable {
		return nil, err
	}

	name, err := ewmh.WmNameGet(m.xu, wid)
	if err != nil {
		return nil, err
	}

	pid, err := ewmh.WmPidGet(m.xu, wid)
	if err != nil {
		return nil, err
	}
	return &WindowInfo{
		Id:   wid,
		Name: name,
		Pid:  uint32(pid),
	}, nil
}

func (infos WindowInfos) GetByName(name string) *WindowInfo {
	for _, info := range infos {
		if info.Name == name {
			return info
		}
	}
	return nil
}
