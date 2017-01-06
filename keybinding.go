package main

import (
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
		"control-mod1-t",
		"control-mod1-delete",
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
	logger.Infof("Key press event mod: %s, key: %s", modStr, key)
	accel := modStr + "-" + key
	switch {
	case isAccelEqual(m.xu, accel, "control-mod1-t"):
		logger.Info("Will launch deepin-terminal")
		go func() {
			err := exec.Command("/bin/sh", "-c", "exec deepin-terminal").Run()
			if err != nil {
				logger.Error("Failed to open deepin-terminal:", err)
			}
		}()
	case isAccelEqual(m.xu, accel, "control-mod1-delete"):
		//exit
		os.Exit(0)
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