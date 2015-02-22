package main

import (
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/kasworld/actionstat"
	"github.com/kasworld/go-sdlgui"
	"github.com/kasworld/go-sdlgui/analogueclock"
	// "github.com/kasworld/htmlcolors"
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

	Stat *actionstat.ActionStat
}

func NewApp() *App {
	app := App{
		SdlCh: make(chan interface{}, 1),
		Keys:  make(map[sdl.Scancode]bool),
		Win:   sdlgui.NewWindow("SDL GUI Example", 512, 512, true),

		Stat: actionstat.NewActionStat(),
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
	c := analogueclock.New()
	g.AddControl(c)

}

func (app *App) Run() {
	// need to co-exist sdl lib
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// start sdl event loop
	sdlgui.SDLEvent2Ch(app.SdlCh)

	timerInfoCh := time.Tick(time.Duration(1000) * time.Millisecond)
	timerDrawCh := time.Tick(time.Duration(1000/60) * time.Millisecond)

	for !app.Quit {
		select {
		case data := <-app.SdlCh:
			if app.Win.ProcessSDLMouseEvent(data) ||
				app.Keys.ProcessSDLKeyEvent(data) {
				app.Quit = true
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
