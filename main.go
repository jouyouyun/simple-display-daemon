package main

import (
	"os"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

const (
	dbusDest = "com.deepin.DisplayManager"
	dbusPath = "/com/deepin/DisplayManager"
	dbusIFC  = dbusDest
)

var logger = log.NewLogger("DisplayManager")

func main() {
	m, err := newManager()
	if err != nil {
		logger.Error("Failed to new manager:", err)
		return
	}

	err = dbus.InstallOnSession(m)
	if err != nil {
		logger.Error("Failed to install session bus:", err)
		m.destroy()
		return
	}
	dbus.DealWithUnhandledMessage()

	m.init()
	m.drawBackground(defaultBackgroundFile, int(m.width), int(m.height))
	m.inhibit()
	m.grabAccels()
	go m.handleEventChanged()
	launchApps()
	err = dbus.Wait()
	if err != nil {
		logger.Error("Lost dbus connection:", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
