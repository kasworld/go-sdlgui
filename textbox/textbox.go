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

package textbox

import (
	"fmt"

	"github.com/kasworld/go-sdl2/sdl"
	"github.com/kasworld/go-sdlgui"
	"github.com/kasworld/htmlcolors"
)

type TextBoxControl struct {
	*sdlgui.Control
	font      *sdlgui.Font
	bg        htmlcolors.RGBA
	fg        htmlcolors.RGBA
	texts     []string
	linecount int
}

func New(x, y, z int, wx, wy int, linecount int, font *sdlgui.Font) *TextBoxControl {
	tc := TextBoxControl{
		sdlgui.NewControl(x, y, z, wx, wy),
		font,
		htmlcolors.Black.ToRGBA(),
		htmlcolors.White.ToRGBA(),
		make([]string, 0),
		linecount,
	}
	tc.BorderSize = 2
	tc.BorderType = 1
	return &tc
}
func (tc *TextBoxControl) DrawSurface() {
	if !tc.ContentsChanged {
		return
	}
	tc.ContentsChanged = false

	tc.Rend.SetDrawColor(tc.bg[0], tc.bg[1], tc.bg[2], tc.bg[3])
	tc.Rend.Clear()
	st := 0
	if len(tc.texts) > tc.linecount {
		st = len(tc.texts) - tc.linecount
	}
	for i, text := range tc.texts[st:] {
		surface, srcRect := tc.font.MakeSurface(tc.fg, text)
		defer surface.Free()
		dstRect := sdl.Rect{
			tc.BorderSize, int32(i*int(tc.H)/tc.linecount) - tc.BorderSize,
			srcRect.W*int32(tc.H/tc.linecount)/srcRect.H - tc.BorderSize, int32(tc.H / tc.linecount),
		}
		if dstRect.W < int32(tc.W) {
			surface.Blit(&srcRect, tc.Suf, &dstRect)
		} else {
			dstRect.W = int32(tc.W)
			surface.BlitScaled(&srcRect, tc.Suf, &dstRect)
		}
	}
	tc.Rend.Present()
	tc.Win.AddUpdateControl(tc)
}
func (tc *TextBoxControl) MouseOver(x, y int, btnstate uint32) {
}
func (tc *TextBoxControl) MouseButton(x, y int, btnnum uint8, btnstate uint8) {
	tc.bg, tc.fg = tc.fg, tc.bg
	tc.ContentsChanged = true
	tc.DrawSurface()
}
func (tc *TextBoxControl) MouseWheel(x, y int, dx int32, dy int32, btnstate uint32) {
}
func (tc *TextBoxControl) SetBG(bg htmlcolors.RGBA) {
	if bg == tc.bg {
		return
	}
	tc.ContentsChanged = true
	tc.bg = bg
}
func (tc *TextBoxControl) SetFG(fg htmlcolors.RGBA) {
	if fg == tc.fg {
		return
	}
	tc.ContentsChanged = true
	tc.fg = fg
}
func (tc *TextBoxControl) AddText(format string, a ...interface{}) {
	tc.ContentsChanged = true
	text := fmt.Sprintf(format, a...)
	tc.texts = append(tc.texts, text)
}
func (tc *TextBoxControl) SetTexts(texts []string) {
	tc.ContentsChanged = true
	tc.texts = texts
}
