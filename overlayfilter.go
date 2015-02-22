package sdlgui

import (
	"github.com/kasworld/htmlcolors"
)

type OverlayFilter []htmlcolors.RGBA

func MakeOverlayFilter(n int, b htmlcolors.RGBA) OverlayFilter {
	rtn := make(OverlayFilter, n)
	for i := 0; i < n; i++ {
		for j, v := range b {
			if int(v) < i*256/n {
				rtn[i][j] = 0
			} else {
				rtn[i][j] = v - uint8(i*256/n)
			}
		}
	}
	return rtn
}
