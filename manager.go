package main

import (
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/dpms"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"pkg.deepin.io/dde/api/drandr"
	"pkg.deepin.io/lib/dbus"
	"sync"
	"time"
)

// TODO:
// 1. block output changed signal [X]
// 2. list connected output details [X]
// 3. no black screen and no suspend [X]
// 4. loop to query output state
// 5. move window to special position [X]
// 6. draw background [X]
// 7. keybinding [X]

type Manager struct {
	conn        *xgb.Conn
	xu          *xgbutil.XUtil
	root        xproto.Window
	outputInfos drandr.OutputInfos
	modeInfos   drandr.ModeInfos

	width  uint16
	height uint16

	outputLocker sync.Mutex
	eventLocker  sync.Mutex

	Changed func()
}

type OutputInfo struct {
	Name   string
	X      int16
	Y      int16
	Width  uint16
	Height uint16
}

func newManager() (*Manager, error) {
	xu, err := xgbutil.NewConn()
	if err != nil {
		return nil, err
	}
	err = randr.Init(xu.Conn())
	if err != nil {
		logger.Error("Failed to init randr:", err)
	}

	screenInfo, err := drandr.GetScreenInfo(xu.Conn())
	if err != nil {
		return nil, err
	}

	err = dpms.Init(xu.Conn())
	if err != nil {
		logger.Error("Failed to init dpms:", err)
	}

	var m = &Manager{
		conn:        xu.Conn(),
		xu:          xu,
		root:        xproto.Setup(xu.Conn()).DefaultScreen(xu.Conn()).Root,
		outputInfos: screenInfo.Outputs,
		modeInfos:   screenInfo.Modes,
	}
	m.width, m.height = screenInfo.GetScreenSize()
	return m, nil
}

func (m *Manager) init() {
	m.checkScreenStatus()
	m.joinExtendMode()
	m.updateOutputInfo()
}

func (m *Manager) destroy() {
	if m.conn == nil {
		return
	}
	m.conn.Close()
	m.conn = nil
}

func (m *Manager) checkScreenStatus() {
	// if all output was invalid, wait until output validity
	for {
		if len(m.outputInfos) != 0 && len(m.modeInfos) != 0 {
			break
		}

		err := doAction("xrandr --auto")
		if err != nil {
			logger.Warningf("Try open output failed %v, try again", err)
		}
		time.Sleep(time.Second * 2)
		m.updateOutputInfo()
	}
}

func (m *Manager) joinExtendMode() {
	connected := m.outputInfos.ListValidOutputs().ListConnectionOutputs()
	names := connected.ListNames()
	infos, _ := newOutputInfosFromFile(outputConfigFile)
	if len(infos) == 0 {
		infos, _ = newOutputInfosFromFile(defaultOutputFile)
		if len(infos) == 0 {
			goto output
		}
	}

	if !checkOutputConfigValidity(names, infos) {
		logger.Warning("Output config invalid:", infos, names)
		goto output
	} else {
		m.joinExtendModeFromConfigInfos(infos)
		return
	}
output:
	m.joinExtendModeFromOutputs(connected)
}

func (m *Manager) joinExtendModeFromOutputs(outputs drandr.OutputInfos) {
	var cmd = "xrandr "
	startx := uint16(0)
	for _, output := range outputs {
		//if !canReadEDID(output.Name) {
		//logger.Warning("Failed to read edid for:", output.Name)
		//continue
		//}
		cmd += " --output " + output.Name
		modes := m.getOutputModes(output.Name)
		var mode drandr.ModeInfo = modes.Best()
		if v := modes.QueryBySize(1024, 768); v.Width != 0 && v.Height != 0 {
			mode = v
		}
		cmd += fmt.Sprintf(" --mode %dx%d --pos %dx0 ", mode.Width, mode.Height, startx)
		if startx == 0 {
			cmd += " --primary "
		}
		startx += mode.Width
	}

	logger.Debug("[joinExtendModeFromOutputs] command:", cmd)
	err := doAction(cmd)
	if err != nil {
		logger.Error("[joinExtendModeFromOutputs] failed:", err)
	}
}

func (m *Manager) joinExtendModeFromConfigInfos(infos []OutputInfo) {
	var cmd = "xrandr "
	primary := false
	for _, info := range infos {
		//if !canReadEDID(info.Name) {
		//logger.Warning("Failed to read edid for:", info.Name)
		//continue
		//}
		cmd += " --output " + info.Name
		cmd += fmt.Sprintf(" --mode %dx%d --pos %dx%d ", info.Width, info.Height, info.X, info.Y)
		if !primary && info.X == 0 {
			primary = true
			cmd += " --primary "
		}
	}

	logger.Debug("[joinExtendModeFromConfigInfos] command:", cmd)
	err := doAction(cmd)
	if err != nil {
		logger.Error("[joinExtendModeFromConfigInfos] failed:", err)
	}
}

func (m *Manager) handleEventChanged() {
	err := randr.SelectInputChecked(m.conn, m.root,
		randr.NotifyMaskOutputChange|randr.NotifyMaskOutputProperty|
			randr.NotifyMaskCrtcChange|randr.NotifyMaskScreenChange).Check()
	if err != nil {
		logger.Error("Failed to select input event:", err)
		return
	}
	for {
		e, err := m.conn.WaitForEvent()
		if err != nil {
			continue
		}
		m.eventLocker.Lock()
		logger.Debug("[Debug] output event:", e.String())
		switch ee := e.(type) {
		case randr.NotifyEvent:
			switch ee.SubCode {
			case randr.NotifyCrtcChange:
			case randr.NotifyOutputChange:
				m.updateOutputInfo()
			case randr.NotifyOutputProperty:
			}
		case randr.ScreenChangeNotifyEvent:
			m.updateOutputInfo()
		case xproto.KeyPressEvent:
			m.handleKeyPressEvent(ee)
		}
		m.eventLocker.Unlock()
	}
}

func (m *Manager) updateOutputInfo() {
	screenInfo, err := drandr.GetScreenInfo(m.conn)
	if err != nil {
		logger.Error("Failed to get screen info:", err)
		return
	}
	m.outputLocker.Lock()
	m.outputInfos, m.modeInfos = screenInfo.Outputs, screenInfo.Modes
	m.width, m.height = screenInfo.GetScreenSize()
	m.outputLocker.Unlock()

	oldLen := len(prevConnected)
	now := m.ListConnectedOutput()
	if oldLen != len(now) {
		dbus.Emit(m, "Changed")
	}
}

func (m *Manager) getOutputModes(name string) drandr.ModeInfos {
	info := m.outputInfos.QueryByName(name)
	var modes drandr.ModeInfos
	for _, id := range info.Modes {
		modes = append(modes, m.modeInfos.Query(id))
	}
	return modes
}

// inhibit no blank and no suspend
func (m *Manager) inhibit() {
	xproto.SetScreenSaver(m.conn, 0, 0, 0, 0)
	// dpms.SetTimeouts(m.conn, 0, 0, 0)
	err := dpms.DisableChecked(m.conn).Check()
	if err != nil {
		logger.Warning("Failed to disable dpms:", err)
	}
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}
