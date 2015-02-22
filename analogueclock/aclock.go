package analogueclock

import (
	"github.com/kasworld/go-sdlgui"
	"github.com/kasworld/htmlcolors"
)

type Clock struct {
	*sdlgui.Control
	font *sdlgui.Font
	bg   htmlcolors.RGBA
	fg   htmlcolors.RGBA
}

func New() *Clock {
	c := Clock{}
	return &c
}
