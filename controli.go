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
