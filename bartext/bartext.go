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

package bartext

import (
	"fmt"

	"github.com/kasworld/go-sdl2/sdl"
	"github.com/kasworld/go-sdlgui"
	"github.com/kasworld/htmlcolors"
)

type TextControl struct {
	*sdlgui.Control
	text     string
	font     *sdlgui.Font
	bg       htmlcolors.RGBA
	fg       htmlcolors.RGBA
	bar      float64
	aligndir uint8
}

func New(x, y, z int, wx, wy int, s string, font *sdlgui.Font) *TextControl {
	tc := TextControl{
		sdlgui.NewControl(x, y, z, wx, wy),
		s,
		font,
		htmlcolors.Black.ToRGBA(),
		htmlcolors.White.ToRGBA(),
		1.0,
		sdlgui.Align_Left,
	}
	tc.BorderSize = 2
	tc.BorderType = 1
	return &tc
}

func (tc *TextControl) DrawSurface() {
	if !tc.ContentsChanged {
		return
	}
	tc.ContentsChanged = false

	tc.Rend.SetDrawColor(127, 127, 127, 255)
	tc.Rend.Clear()
	tc.Rend.SetDrawColor(tc.bg[0], tc.bg[1], tc.bg[2], tc.bg[3])
	barrect := sdl.Rect{
		0, 0,
		int32(float64(tc.W) * tc.bar), int32(tc.H),
	}
	tc.Rend.FillRect(&barrect)

	surface, srcRect := tc.font.MakeSurface(tc.fg, tc.text)
	defer surface.Free()

	torect := tc.GetClientConnRect().ShrinkSym(int(tc.BorderSize))
	if srcRect.W < int32(torect.W) {
		dstRect := sdlgui.SdlInRectAlign(torect, srcRect, tc.aligndir)
		surface.Blit(&srcRect, tc.Suf, &dstRect)
	} else {
		tsrcRect := srcRect
		tsrcRect.W = int32(tc.W)
		dstRect := sdlgui.SdlInRectAlign(torect, tsrcRect, tc.aligndir)
		surface.BlitScaled(&srcRect, tc.Suf, &dstRect)
	}

	tc.Rend.Present()
	tc.Win.AddUpdateControl(tc)
}
func (tc *TextControl) MouseOver(x, y int, btnstate uint32) {
}
func (tc *TextControl) MouseButton(x, y int, btnnum uint8, btnstate uint8) {
	// log.Info("%v btn %v %v ", tc.GetID(), btnnum, btnstate)
	tc.bg, tc.fg = tc.fg, tc.bg
	tc.ContentsChanged = true
	tc.DrawSurface()
}
func (tc *TextControl) SetBG(bg htmlcolors.RGBA) {
	if bg == tc.bg {
		return
	}
	tc.ContentsChanged = true
	tc.bg = bg
}
func (tc *TextControl) SetFG(fg htmlcolors.RGBA) {
	if fg == tc.fg {
		return
	}
	tc.ContentsChanged = true
	tc.fg = fg
}
func (tc *TextControl) SetFGBG(fg, bg htmlcolors.RGBA) {
	tc.SetFG(fg)
	tc.SetBG(bg)
}
func (tc *TextControl) SetText(format string, a ...interface{}) {
	text := fmt.Sprintf(format, a...)
	if text == tc.text {
		return
	}
	tc.ContentsChanged = true
	tc.text = text
}
func (tc *TextControl) SetBar(bar float64) {
	if bar == tc.bar {
		return
	}
	tc.ContentsChanged = true
	tc.bar = bar
}
func (tc *TextControl) SetFGBGText(fg, bg htmlcolors.RGBA, format string, a ...interface{}) {
	tc.SetFG(fg)
	tc.SetBG(bg)
	tc.SetText(format, a...)
}
func (tc *TextControl) SetBarText(bar float64, format string, a ...interface{}) {
	tc.SetBar(bar)
	tc.SetText(format, a...)
}
func (tc *TextControl) SetAlign(align uint8) {
	if align == tc.aligndir {
		return
	}
	tc.aligndir = align
	tc.ContentsChanged = true
}
