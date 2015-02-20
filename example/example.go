package main

import (
	"runtime"
	"time"

	// "github.com/veandco/go-sdl2/sdl"

	"github.com/kasworld/actionstat"
	"github.com/kasworld/htmlcolors"
	"github.com/kasworld/log"

	"github.com/kasworld/go-sdlgui"
)

func main() {
	runtime.LockOSThread()

	app := App{
		Stat:  actionstat.NewActionStat(),
		sdlCh: make(chan interface{}, 1),
	}
	app.Run()
	runtime.UnlockOSThread()
}

type App struct {
	Quit  bool
	sdlCh chan interface{}
	Stat  *actionstat.ActionStat

	window   *sdlgui.Window
	controls []sdlgui.ControlI
	msgtexts *sdlgui.TextBoxControl
	barctrl  *sdlgui.TextControl
}

func (g *App) addControls() {
	g.window = sdlgui.NewWindow("", 1024, 800, true)

	g.msgtexts = sdlgui.NewTextBoxControl(
		0, 0, 0,
		1024, 720, 60,
		sdlgui.LoadFont("DejaVuSerif.ttf", 12))
	g.msgtexts.SetBG(htmlcolors.Gray.ToRGBA())
	g.window.AddControl(g.msgtexts)

	g.barctrl = sdlgui.NewTextControl(
		0, 720, 0,
		1024, 80, "hello",
		sdlgui.LoadFont("DejaVuSerif.ttf", 36))
	g.barctrl.SetBG(htmlcolors.Pink.ToRGBA())
	g.window.AddControl(g.barctrl)

	g.window.UpdateAll()
}

func (g *App) Run() {
	g.addControls()
	sdlgui.SDLEvent2Ch(g.sdlCh)
	timerInfoCh := time.Tick(time.Duration(1000) * time.Millisecond)
	timerDrawCh := time.Tick(time.Duration(1000/60) * time.Millisecond)
	barlen := 0.0
	for !g.Quit {
		select {
		case data := <-g.sdlCh:
			g.msgtexts.AddText("data %v", data)
			g.barctrl.SetBar(barlen)
			barlen += 0.01
			if barlen > 1 {
				barlen = 0
			}

		case <-timerDrawCh:
			g.msgtexts.DrawSurface()
			g.barctrl.DrawSurface()
			g.window.Update()

		case <-timerInfoCh:
			log.Info("stat %v", g.Stat)
			g.Stat.UpdateLap()
			g.Stat.Inc()
		}
	}
}
