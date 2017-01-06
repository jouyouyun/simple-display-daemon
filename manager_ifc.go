package main

func (m *Manager) ListConnectedOutput() []OutputInfo {
	m.outputLocker.Lock()
	defer m.outputLocker.Unlock()
	connected := m.outputInfos.ListValidOutputs().ListConnectionOutputs()
	var infos []OutputInfo
	for _, output := range connected {
		infos = append(infos, OutputInfo{
			Name:   output.Name,
			X:      output.Crtc.X,
			Y:      output.Crtc.Y,
			Width:  output.Crtc.Width,
			Height: output.Crtc.Height,
		})
	}
	return infos
}
