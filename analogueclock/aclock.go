package analogueclock

import (
	"math"
	"time"

	"github.com/kasworld/go-sdlgui"
	"github.com/kasworld/htmlcolors"
)

type Clock struct {
	*sdlgui.Control
	bg    htmlcolors.RGBA
	fg    htmlcolors.RGBA
	hhand htmlcolors.RGBA
	mhand htmlcolors.RGBA
	shand htmlcolors.RGBA
	t     time.Time
}

func New(x, y, z int, wx, wy int) *Clock {
	c := Clock{
		sdlgui.NewControl(x, y, z, wx, wy),
		htmlcolors.Black.ToRGBA(),
		htmlcolors.White.ToRGBA(),
		htmlcolors.Green.ToRGBA(),
		htmlcolors.Blue.ToRGBA(),
		htmlcolors.Red.ToRGBA(),
		time.Now(),
	}
	c.BorderSize = 2
	c.BorderType = 1
	return &c
}

func (tc *Clock) SetBG(bg htmlcolors.RGBA) {
	if bg == tc.bg {
		return
	}
	tc.ContentsChanged = true
	tc.bg = bg
}

func (tc *Clock) SetTime(t time.Time) {
	if t == tc.t {
		return
	}
	tc.ContentsChanged = true
	tc.t = t
}

func (tc *Clock) drawFace() {
	tc.Rend.SetDrawColor(tc.fg[0], tc.fg[1], tc.fg[2], tc.fg[3])
	c := tc.GetRect().Center()
	rx := tc.GetRect().SizeVector()[0]/2 - 3
	ry := tc.GetRect().SizeVector()[1]/2 - 3
	rate := float64(ry) / float64(rx)
	for x := -rx; x < rx; x++ {
		theta := math.Acos(float64(x) / float64(rx))
		y := int(float64(x) * math.Tan(theta) * rate)
		if x == 0 {
			y = ry
		}
		tc.Rend.DrawLine(c[0]+x, c[1]-y, c[0]+x, c[1]+y)
	}
}

func (tc *Clock) drawHand(co htmlcolors.RGBA, angle float64, l float64) {
	angle -= 90
	x1, y1 := tc.W/2, tc.H/2
	x2 := x1 + int(math.Cos(angle/180*math.Pi)*float64(x1)*l)
	y2 := y1 + int(math.Sin(angle/180*math.Pi)*float64(y1)*l)
	tc.Rend.SetDrawColor(co[0], co[1], co[2], co[3])
	tc.Rend.DrawLine(x1, y1, x2, y2)
}

func (tc *Clock) time2Angle() (float64, float64, float64) {
	h := float64(tc.t.Hour()%12) / 12 * 360
	m := float64(tc.t.Minute()) / 60 * 360
	s := float64(tc.t.Second()) / 60 * 360
	return h, m, s
}

func (tc *Clock) DrawSurface() {
	if !tc.ContentsChanged {
		return
	}
	tc.ContentsChanged = false

	tc.Rend.SetDrawColor(tc.bg[0], tc.bg[1], tc.bg[2], tc.bg[3])
	tc.Rend.Clear()

	h, m, s := tc.time2Angle()
	tc.drawFace()
	tc.drawHand(tc.hhand, h, 0.60)
	tc.drawHand(tc.mhand, m, 0.75)
	tc.drawHand(tc.shand, s, 0.90)

	tc.Rend.Present()
	tc.Win.AddUpdateControl(tc)
}
