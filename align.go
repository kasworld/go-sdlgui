package sdlgui

import (
	"github.com/kasworld/go-sdl2/sdl"
	"github.com/kasworld/rect"
)

const (
	Align_Center = iota
	Align_Up
	Align_UpRight
	Align_Right
	Align_DownRight
	Align_Down
	Align_DownLeft
	Align_Left
	Align_UpLeft
)

var AlignVt = [...][2]int{
	Align_Center:    [2]int{0, 0},
	Align_Up:        [2]int{0, -1},
	Align_UpRight:   [2]int{1, -1},
	Align_Right:     [2]int{1, 0},
	Align_DownRight: [2]int{1, 1},
	Align_Down:      [2]int{0, 1},
	Align_DownLeft:  [2]int{-1, 1},
	Align_Left:      [2]int{-1, 0},
	Align_UpLeft:    [2]int{-1, -1},
}

// rect inner align
func SdlInRectAlign(dst rect.Rect, src sdl.Rect, aligndir int) sdl.Rect {
	rtn := src
	vt := AlignVt[aligndir]
	switch vt[0] {
	case -1:
		rtn.X = int32(dst.X)
	case 0:
		rtn.X = int32(dst.X + dst.W/2 - int(src.W)/2)
	case 1:
		rtn.X = int32(dst.X + dst.W - int(src.W))
	}
	switch vt[1] {
	case -1:
		rtn.Y = int32(dst.Y)
	case 0:
		rtn.Y = int32(dst.Y + dst.H/2 - int(src.H/2))
	case 1:
		rtn.Y = int32(dst.Y + dst.H - int(src.H))
	}
	return rtn
}

// point base align
func SdlDestRect(src sdl.Rect, x, y int32, aligndir int) sdl.Rect {
	rtn := sdl.Rect{x, y, src.W, src.H}
	vt := AlignVt[aligndir]
	rtn.X = x + src.W/2*(int32(vt[0])-1)
	rtn.Y = y + src.H/2*(int32(vt[1])-1)
	return rtn
}
