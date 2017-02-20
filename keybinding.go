package main

import (
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
	"os"
	"os/exec"
	"pkg.deepin.io/lib/strv"
	"strings"
)

func (m *Manager) grabAccels() {
	keybind.Initialize(m.xu)

	var accels = []string{
		"mod4-r",      // startdde
		"mod4-t",      // terminal
		"mod4-delete", // logout
		"mod4-1",      // xterm
		"mod4-2",
		"mod4-3",
	}

	for _, accel := range accels {
		err := doGrab(m.xu, accel)
		if err != nil {
			logger.Errorf("Failed to grab '%s' : %v", accel, err)
		}
	}
}

func (m *Manager) handleKeyPressEvent(ev xproto.KeyPressEvent) {
	modStr := filterInvalidMod(keybind.ModifierString(ev.State))
	key := keybind.LookupString(m.xu, ev.State, ev.Detail)
	logger.Debugf("Key press event mod: %s, key: %s", modStr, key)
	accel := modStr + "-" + key
	switch {
	case isAccelEqual(m.xu, accel, "mod4-r"):
		logger.Debug("Will launch startdde")
		if isAppLaunched("startdde") {
			runApp("killall dde-session-daemon dde-launcher dde-session-initializer")
			err := runApp("systemctl --user stop startdde.scope")
			if err != nil {
				logger.Error("Stop startdde failed:", err)
			} else {
				m.inhibit()
				m.init()
				m.drawBackground(defaultBackgroundFile, int(m.width), int(m.height))
			}
			return
		}
		go func() {
			err := runApp("systemd-run  --scope --user --unit startdde /usr/bin/startdde")
			if err != nil {
				logger.Error("Failed to launch startdde:", err)
				return
			}
		}()
	case isAccelEqual(m.xu, accel, "mod4-t"):
		logger.Debug("Will launch terminal")
		go runApp("x-terminal-emulator")
	case isAccelEqual(m.xu, accel, "mod4-delete"):
		//exit
		os.Exit(0)
	case isAccelEqual(m.xu, accel, "mod4-1"):
		go launchXTerm(1)
	case isAccelEqual(m.xu, accel, "mod4-2"):
		go launchXTerm(2)
	case isAccelEqual(m.xu, accel, "mod4-3"):
		go launchXTerm(3)
	}
}

func doGrab(xu *xgbutil.XUtil, accel string) error {
	mod, codes, err := keybind.ParseString(xu, accel)
	if err != nil {
		return err
	}

	for _, code := range codes {
		err := keybind.GrabChecked(xu, xu.RootWin(), mod, code)
		if err != nil {
			return err
		}
	}
	return nil
}

func launchXTerm(screen int) {
	if isAppLaunched("xterm") {
		runApp("killall xterm")
		return
	}
	runApp("openbox &")
	runApp(fmt.Sprintf("xterm -geometry 100x40+%d+100", 100+(screen-1)*1024))
}

func isAccelEqual(xu *xgbutil.XUtil, accel1, accel2 string) bool {
	if accel1 == accel2 {
		return true
	}

	mod1, codes1, err := keybind.ParseString(xu, accel1)
	if err != nil {
		return false
	}

	mod2, codes2, err := keybind.ParseString(xu, accel2)
	if err != nil {
		return false
	}

	if mod1 != mod2 {
		return false
	}

	l1, l2 := len(codes1), len(codes2)
	if l1 != l2 {
		return false
	}

	for i := 0; i < l1; i++ {
		if codes1[i] != codes2[i] {
			return false
		}
	}
	return true
}

var invalidMod = strv.Strv{
	"mod2",
	"lock",
	"num_lock",
	"caps_lock",
}

func filterInvalidMod(mod string) string {
	if len(mod) == 0 {
		return ""
	}

	list := strings.Split(mod, "-")
	var ret []string
	for _, v := range list {
		if invalidMod.Contains(v) {
			continue
		}
		ret = append(ret, v)
	}
	return strings.Join(ret, "-")
}

func isAppLaunched(app string) bool {
	out, err := exec.Command("/bin/sh", "-c",
		"ps aux|grep -w "+app).CombinedOutput()
	if err != nil {
		logger.Warningf("Failed to check %s state: %v, %v", app, string(out), err)
		return false
	}
	lines := strv.Strv(strings.Split(string(out), "\n"))
	lines = lines.FilterEmpty()
	return len(lines) > 2
}
