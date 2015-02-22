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
	NewApp().Run()
}

type App struct {
	Quit     bool
	SdlCh    chan interface{}
	Keys     sdlgui.KeyState
	Win      *sdlgui.Window
	Controls sdlgui.ControlIList

	Stat    *actionstat.ActionStat
	msgtext *sdlgui.TextBoxControl
	barctrl *sdlgui.TextControl
}

func NewApp() *App {
	app := App{
		SdlCh: make(chan interface{}, 1),
		Keys:  make(map[sdl.Scancode]bool),
		Win:   sdlgui.NewWindow("SDL GUI Example", 1024, 800, true),

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

type Clock struct {
	*sdlgui.Control
	font *sdlgui.Font
	bg   htmlcolors.RGBA
	fg   htmlcolors.RGBA
}

func (g *App) addControls() {
	g.msgtext = sdlgui.NewTextBoxControl(
		0, 0, 0,
		1024, 720, 60,
		sdlgui.LoadFont("DejaVuSerif.ttf", 12))
	g.msgtext.SetBG(htmlcolors.Gray.ToRGBA())
	g.AddControl(g.msgtext)

	g.barctrl = sdlgui.NewTextControl(
		0, 720, 0,
		1024, 80, "hello",
		sdlgui.LoadFont("DejaVuSerif.ttf", 36))
	g.barctrl.SetBG(htmlcolors.Pink.ToRGBA())
	g.AddControl(g.barctrl)

}

func (app *App) Run() {
	// need to co-exist sdl lib
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// start sdl event loop
	sdlgui.SDLEvent2Ch(app.SdlCh)

	timerInfoCh := time.Tick(time.Duration(1000) * time.Millisecond)
	timerDrawCh := time.Tick(time.Duration(1000/60) * time.Millisecond)
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
