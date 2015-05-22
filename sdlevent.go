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
	"github.com/kasworld/go-sdl2/sdl"
)

type KeyState map[sdl.Scancode]bool

func (ks KeyState) ProcessSDLKeyEvent(data interface{}) (quit bool) {
	quit = false
	switch t := data.(type) {
	case *sdl.KeyUpEvent:
		delete(ks, t.Keysym.Scancode)
	case *sdl.KeyDownEvent:
		if t.Repeat == 0 {
			ks[t.Keysym.Scancode] = true
		}
		switch t.Keysym.Scancode {
		case sdl.SCANCODE_ESCAPE:
			quit = true
		}
	}
	return
}

func (w *Window) ProcessSDLMouseEvent(data interface{}) (quit bool) {
	quit = false
	switch t := data.(type) {
	case *sdl.MouseMotionEvent:
		c, x, y, btnstate := w.ProcessMouseMotionEvent(t)
		if c != nil {
			c.MouseOver(x, y, btnstate)
		}
	case *sdl.MouseButtonEvent:
		c, x, y, n, s := w.ProcessMouseButtonEvent(t)
		if c != nil {
			c.MouseButton(x, y, n, s)
		}
	case *sdl.MouseWheelEvent:
		c, x, y, dx, dy, btnstate := w.ProcessMouseWheelEvent(t)
		if c != nil {
			c.MouseWheel(x, y, dx, dy, btnstate)
		}
	case *sdl.QuitEvent:
		// log.Printf("quit %v\n", t)
		quit = true
	}
	return
}

func SDLEvent2Ch(ch chan<- interface{}) {
	go func() {
		for {
			event := sdl.WaitEvent()
			if event == nil {
				continue
			}
			ch <- event
		}
	}()
}
