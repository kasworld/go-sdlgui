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

package main

import (
	"time"

	"github.com/kasworld/actionstat"
	"github.com/kasworld/go-sdl2/sdl"
	"github.com/kasworld/go-sdlgui"
	"github.com/kasworld/go-sdlgui/bartext"
	"github.com/kasworld/go-sdlgui/textbox"
	"github.com/kasworld/htmlcolors"
	"github.com/kasworld/log"
)

func main() {
	NewApp().Run()
}

type App struct {
	Quit     bool
	SdlCh    chan interface{}
	Keys     sdlgui.KeyState
	Win      *sdlgui.Window
	Controls sdlgui.ControlIList

	Stat    *actionstat.ActionStat
	msgtext *textbox.TextBoxControl
	barctrl *bartext.TextControl
}

func NewApp() *App {
	app := App{
		SdlCh: make(chan interface{}, 1),
		Keys:  make(map[sdl.Scancode]bool),
		Win:   sdlgui.NewWindow("SDL GUI Example", 1024, 800, true),

		Stat: actionstat.New(),
	}
	app.addControls()
	app.Win.UpdateAll()
	return &app
}

func (app *App) AddControl(c sdlgui.ControlI) {
	app.Controls = append(app.Controls, c)
	app.Win.AddControl(c)
}

// change as app's need

func (g *App) addControls() {
	g.msgtext = textbox.New(
		0, 0, 0,
		1024, 720, 60,
		sdlgui.LoadFont("NanumGothic.ttf", 12))
	g.msgtext.SetBG(htmlcolors.Gray.ToRGBA())
	g.AddControl(g.msgtext)

	dispstr := "안녕"
	println(dispstr)
	log.Printf("%v %x %+q", dispstr, dispstr, dispstr)

	g.barctrl = bartext.New(
		0, 720, 0,
		1024, 80, dispstr,
		sdlgui.LoadFont("NanumGothic.ttf", 36))
	g.barctrl.SetBG(htmlcolors.Pink.ToRGBA())
	g.AddControl(g.barctrl)

}

func (app *App) Run() {
	// need to co-exist sdl lib
	// runtime.LockOSThread()
	// defer runtime.UnlockOSThread()

	// start sdl event loop
	sdlgui.SDLEvent2Ch(app.SdlCh)

	timerInfoCh := time.Tick(time.Duration(1000) * time.Millisecond)
	timerDrawCh := time.Tick(time.Duration(1000/60) * time.Millisecond)
	// timerDrawCh := time.Tick(time.Microsecond)
	barlen := 0.0

	for !app.Quit {
		select {
		case data := <-app.SdlCh:
			if app.Win.ProcessSDLMouseEvent(data) ||
				app.Keys.ProcessSDLKeyEvent(data) {
				app.Quit = true
			}
			app.msgtext.AddText("data %v", data)
			app.barctrl.SetBar(barlen)
			barlen += 0.01
			if barlen > 1 {
				barlen = 0
			}
			app.Stat.Inc()

		case <-timerDrawCh:
			for _, v := range app.Controls {
				v.DrawSurface()
			}
			app.Win.Update()

		case <-timerInfoCh:
			log.Info("stat %v", app.Stat)
			app.Stat.UpdateLap()

		}
	}
}
