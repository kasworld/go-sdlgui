package sdlgui

import (
	"github.com/kasworld/quadtree"
)

type ControlIList []ControlI

func (s ControlIList) Len() int {
	return len(s)
}
func (s ControlIList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ControlIList) Less(i, j int) bool {
	return s[i].GetZ() < s[j].GetZ()
}

type ControlI interface {
	quadtree.QuadTreeObjI

	MouseOver(x, y int, btnstate uint32)
	MouseButton(x, y int, btnnum uint8, btnstate uint8)
	MouseWheel(x, y int, dx int32, dy int32, btnstate uint32)

	UpdateContents()
	DrawSurface()
	UpdateToWindow()

	SetWindow(w *Window)
	GetZ() int
	Show(visible bool)
	IsVisible() bool
	IsTransparent() bool
}
