package main

import "time"

func (m *Manager) waitEDPConnected() {
	if canReadModes("eDP1") {
		logger.Info("eDP1 has modes, not wait......")
		return
	}

	for {
		time.Sleep(time.Second * 5)
		if canReadModes("eDP1") {
			logger.Info("eDP1 connected, reset position......")
			runApp("xrandr --output eDP1 --off")
			runApp("xrandr --output eDP1 --pos 2047x0")
			runApp("xrandr --output eDP1 --pos 2048x0")
			return
		}
	}
}
