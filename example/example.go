package main

import (
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/kasworld/actionstat"
	"github.com/kasworld/htmlcolors"
	"github.com/kasworld/log"

	"github.com/kasworld/go-sdlgui"
)

func main() {
	runtime.LockOSThread()

	app := App{
		Stat:  actionstat.NewActionStat(),
		SdlCh: make(chan interface{}, 1),
		Keys:  make(map[sdl.Scancode]bool),
	}
	app.Run()
	runtime.UnlockOSThread()
}

type App struct {
	Quit  bool
	SdlCh chan interface{}
	Keys  sdlgui.KeyState
	Win   *sdlgui.Window
	Stat  *actionstat.ActionStat

	controls []sdlgui.ControlI
	msgtexts *sdlgui.TextBoxControl
	barctrl  *sdlgui.TextControl
}

func (g *App) addControls() {
	g.Win = sdlgui.NewWindow("", 1024, 800, true)

	g.msgtexts = sdlgui.NewTextBoxControl(
		0, 0, 0,
		1024, 720, 60,
		sdlgui.LoadFont("DejaVuSerif.ttf", 12))
	g.msgtexts.SetBG(htmlcolors.Gray.ToRGBA())
	g.Win.AddControl(g.msgtexts)

	g.barctrl = sdlgui.NewTextControl(
		0, 720, 0,
		1024, 80, "hello",
		sdlgui.LoadFont("DejaVuSerif.ttf", 36))
	g.barctrl.SetBG(htmlcolors.Pink.ToRGBA())
	g.Win.AddControl(g.barctrl)

	g.Win.UpdateAll()
}

func (g *App) Run() {
	g.addControls()
	sdlgui.SDLEvent2Ch(g.SdlCh)
	timerInfoCh := time.Tick(time.Duration(1000) * time.Millisecond)
	timerDrawCh := time.Tick(time.Duration(1000/60) * time.Millisecond)
	barlen := 0.0
	for !g.Quit {
		select {
		case data := <-g.SdlCh:
			if g.Win.ProcessSDLMouseEvent(data) ||
				g.Keys.ProcessSDLKeyEvent(data) {
				g.Quit = true
			}
			g.msgtexts.AddText("data %v", data)
			g.barctrl.SetBar(barlen)
			barlen += 0.01
			if barlen > 1 {
				barlen = 0
			}
			g.Stat.Inc()

		case <-timerDrawCh:
			g.msgtexts.DrawSurface()
			g.barctrl.DrawSurface()
			g.Win.Update()

		case <-timerInfoCh:
			log.Info("stat %v", g.Stat)
			g.Stat.UpdateLap()
		}
	}
}
