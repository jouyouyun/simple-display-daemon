package main

import (
	"github.com/BurntSushi/xgbutil/xgraphics"
	"image"
)

const defaultBackgroundFile = "/usr/share/backgrounds/deepin/desktop.jpg"

func (m *Manager) drawBackground(bg string, w, h int) {
	img := xgraphics.New(m.xu, image.Rect(0, 0, w, h))
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			img.SetBGRA(i, j, xgraphics.BGRA{B: 113, G: 52, R: 48, A: 0})
		}
	}
	defer img.Destroy()

	err := img.XSurfaceSet(m.root)
	if err != nil {
		logger.Error("Failed to set surface:", err)
		return
	}

	err = img.XDrawChecked()
	if err != nil {
		logger.Error("Failed to draw:", err)
		return
	}

	img.XPaint(m.root)
}
