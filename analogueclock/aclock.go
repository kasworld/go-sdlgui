// Copyright 2015 SeukWon Kang (kasworld@gmail.com)
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package analogueclock

import (
	"math"
	"time"

	"github.com/kasworld/go-sdlgui"
	"github.com/kasworld/htmlcolors"
)

const (
	CO_bg = uint8(iota)
	CO_face
	CO_daial
	CO_hhand
	CO_mhand
	CO_shand
	CO_End
)

type Clock struct {
	*sdlgui.Control
	colors [CO_End]htmlcolors.RGBA
	t      time.Time
}

func New(x, y, z int, wx, wy int) *Clock {
	c := Clock{
		sdlgui.NewControl(x, y, z, wx, wy),
		[CO_End]htmlcolors.RGBA{
			htmlcolors.Black.ToRGBA(),
			htmlcolors.Gray.ToRGBA(),
			htmlcolors.White.ToRGBA(),
			htmlcolors.Green.ToRGBA(),
			htmlcolors.Blue.ToRGBA(),
			htmlcolors.Red.ToRGBA(),
		},
		time.Now(),
	}
	c.BorderType = 0
	return &c
}

func (tc *Clock) SetColor(n int, co htmlcolors.RGBA) {
	if co == tc.colors[n] {
		return
	}
	tc.ContentsChanged = true
	tc.colors[n] = co
}

func (tc *Clock) SetTime(t time.Time) {
	if t == tc.t {
		return
	}
	tc.ContentsChanged = true
	tc.t = t
}

func (tc *Clock) drawFace() {
	tc.Rend.SetDrawColor(tc.colors[CO_face][0], tc.colors[CO_face][1], tc.colors[CO_face][2], tc.colors[CO_face][3])
	c := tc.GetRect().Center()
	rx := tc.GetRect().SizeVector()[0] / 2
	ry := tc.GetRect().SizeVector()[1] / 2
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

func (tc *Clock) getFacePos(angle float64, l float64) (int, int) {
	angle -= 90
	xc, yc := tc.W/2, tc.H/2
	x := xc + int(math.Cos(angle/180*math.Pi)*float64(xc)*l)
	y := yc + int(math.Sin(angle/180*math.Pi)*float64(yc)*l)
	return x, y
}

func (tc *Clock) drawHand(co htmlcolors.RGBA, angle float64, l float64) {
	// x1, y1 := tc.W/2, tc.H/2
	x1, y1 := tc.getFacePos(angle, -l/5)
	x2, y2 := tc.getFacePos(angle, l)
	tc.Rend.SetDrawColor(co[0], co[1], co[2], co[3])
	tc.Rend.DrawLine(x1, y1, x2, y2)
}

func (tc *Clock) drawDialLine(co htmlcolors.RGBA, angle float64, s, e float64) {
	x1, y1 := tc.getFacePos(angle, s)
	x2, y2 := tc.getFacePos(angle, e)
	tc.Rend.SetDrawColor(co[0], co[1], co[2], co[3])
	tc.Rend.DrawLine(x1, y1, x2, y2)
}

func (tc *Clock) drawDials() {
	for a := 0.0; a < 360; a += 6 {
		tc.drawDialLine(tc.colors[CO_daial], a, 0.95, 0.99)
	}
	for a := 0.0; a < 360; a += 6 * 5 {
		tc.drawDialLine(tc.colors[CO_daial], a, 0.90, 0.99)
	}
}

func (tc *Clock) time2Angle() (float64, float64, float64) {
	ms := float64(tc.t.Nanosecond()/1000000) / 1000
	s := (float64(tc.t.Second()) + ms)
	m := (float64(tc.t.Minute()) + s/60)
	h := (float64(tc.t.Hour()%12) + m/60)

	sa := s / 60 * 360
	ma := m / 60 * 360
	ha := h / 12 * 360
	return ha, ma, sa
}

func (tc *Clock) DrawSurface() {
	if !tc.ContentsChanged {
		return
	}
	tc.ContentsChanged = false

	tc.Rend.SetDrawColor(tc.colors[CO_bg][0], tc.colors[CO_bg][1], tc.colors[CO_bg][2], tc.colors[CO_bg][3])
	tc.Rend.Clear()

	h, m, s := tc.time2Angle()
	tc.drawFace()
	tc.drawDials()
	tc.drawHand(tc.colors[CO_hhand], h, 0.60)
	tc.drawHand(tc.colors[CO_mhand], m, 0.75)
	tc.drawHand(tc.colors[CO_shand], s, 0.90)

	tc.Rend.Present()
	tc.Win.AddUpdateControl(tc)
}
