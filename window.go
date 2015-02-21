package sdlgui

import (
	"sort"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/kasworld/idgen"
	"github.com/kasworld/log"
	"github.com/kasworld/quadtree"
	"github.com/kasworld/rect"
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
		c.MouseOver(x, y, btnstate)
	case *sdl.MouseButtonEvent:
		c, x, y, n, s := w.ProcessMouseButtonEvent(t)
		c.MouseButton(x, y, n, s)
	case *sdl.MouseWheelEvent:
		c, x, y, dx, dy, btnstate := w.ProcessMouseWheelEvent(t)
		c.MouseWheel(x, y, dx, dy, btnstate)
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

type Window struct {
	ID               int64
	Win              *sdl.Window
	Rend             *sdl.Renderer
	Controls         *quadtree.QuadTree
	controlsToUpdate ControlIList
}

func NewWindow(title string, wx, wy int, show bool) *Window {
	w := Window{
		ID: <-idgen.GenCh(),
	}
	if show {
		w.Win = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			wx, wy, sdl.WINDOW_SHOWN)
	} else {
		w.Win = sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
			wx, wy, sdl.WINDOW_HIDDEN)
	}
	if w.Win == nil {
		log.Fatalf("Failed to create window: %s\n", sdl.GetError())
	}
	// w.Rend = sdl.CreateRenderer(w.Win, -1, sdl.RENDERER_SOFTWARE)
	w.Rend = sdl.CreateRenderer(w.Win, -1, sdl.RENDERER_ACCELERATED)
	if w.Rend == nil {
		log.Fatalf("Failed to create renderer: %s\n", sdl.GetError())
	}

	wr := sdl.Rect{}
	w.Rend.GetViewport(&wr)
	w.Controls = quadtree.NewQuadTree(SdlRect2Rect(wr))
	return &w
}
func (w *Window) Cleanup() {
	if w.Rend != nil {
		w.Rend.Destroy()
	}
	if w.Win != nil {
		w.Win.Destroy()
	}
}
func (w *Window) SetTitle(title string) {
	w.Win.SetTitle(title)
}
func (w *Window) Show(show bool) {
	if show {
		w.Win.Show()
	} else {
		w.Win.Hide()
	}
}
func (w *Window) AddControl(c ControlI) {
	w.Controls.Insert(c)
	c.SetWindow(w)
}
func (w *Window) DelControl(c ControlI) {
	w.Controls.Remove(c)
}
func (w *Window) UpdateAll() {
	wr := sdl.Rect{}
	w.Rend.GetViewport(&wr)
	controlsToUpdate := w.GetControlListByRect(SdlRect2Rect(wr))

	w.UpdateControls(controlsToUpdate)
	w.controlsToUpdate = w.controlsToUpdate[:0]
	w.Rend.Present()
}
func (w *Window) Update() {
	w.UpdateControls(w.controlsToUpdate)
	w.controlsToUpdate = w.controlsToUpdate[:0]
	w.Rend.Present()
}

func (w *Window) UpdateControls(controlsToUpdate ControlIList) {
	toupdatecontrols := w.listVisibleControls(controlsToUpdate)
	for _, v := range toupdatecontrols {
		v.UpdateToWindow()
	}
}

func (w *Window) GetControlListByRect(wr rect.Rect) ControlIList {
	controlsToUpdate := ControlIList{}
	fn := func(qi quadtree.QuadTreeObjI) bool {
		c := qi.(ControlI)
		if c.IsVisible() {
			controlsToUpdate = append(controlsToUpdate, c)
		}
		return false
	}
	w.Controls.QueryByRect(fn, wr)
	return controlsToUpdate
}

func (w *Window) listVisibleControls(in ControlIList) ControlIList {
	sort.Sort(in)
	rtn := ControlIList{}
loop:
	for _, v := range in {
		if !v.IsVisible() {
			continue loop
		}

		cinrect := w.GetControlListByRect(v.GetRect())
		sort.Sort(cinrect)

		for _, w := range cinrect {
			if w.IsTransparent() || !w.IsVisible() || v.GetID() == w.GetID() {
				continue
			}
			if v.GetZ() < w.GetZ() && v.GetRect().IsIn(w.GetRect()) {
				continue loop
			}
		}
		rtn = append(rtn, v)
	}
	return rtn
}

func (w *Window) UpdateRect(wr rect.Rect) {
	controlsToUpdate := w.GetControlListByRect(wr)
	w.UpdateControls(controlsToUpdate)
}

func (w *Window) AddUpdateControl(c ControlI) {
	w.controlsToUpdate = append(w.controlsToUpdate, c)
}
func (w *Window) FindControl(x, y int) ControlI {
	var rtn ControlI
	fn := func(qi quadtree.QuadTreeObjI) bool {
		c := qi.(ControlI)
		if !c.IsVisible() {
			return false
		}
		if rtn == nil {
			rtn = c
		} else {
			if rtn.GetZ() < c.GetZ() {
				rtn = c
			}
		}
		return false
	}
	w.Controls.QueryByPos(fn, [2]int{x, y})
	// log.Info("find control %v %v %v", rtn, x, y)
	return rtn
}

func (w *Window) ProcessMouseEvent() (ControlI, int, int, uint32) {
	mx, my, btnstate := sdl.GetMouseState()
	c := w.FindControl(int(mx), int(my))
	cx, cy := int(mx)-c.GetRect().X, int(my)-c.GetRect().Y
	return c, cx, cy, btnstate
}
func (w *Window) ProcessMouseMotionEvent(t *sdl.MouseMotionEvent) (ControlI, int, int, uint32) {
	_, _, btnstate := sdl.GetMouseState()
	c := w.FindControl(int(t.X), int(t.Y))
	cx, cy := int(t.X)-c.GetRect().X, int(t.Y)-c.GetRect().Y
	return c, cx, cy, btnstate
}
func (w *Window) ProcessMouseButtonEvent(t *sdl.MouseButtonEvent) (ControlI, int, int, uint8, uint8) {
	c := w.FindControl(int(t.X), int(t.Y))
	cx, cy := int(t.X)-c.GetRect().X, int(t.Y)-c.GetRect().Y
	return c, cx, cy, t.Button, t.State
}
func (w *Window) ProcessMouseWheelEvent(t *sdl.MouseWheelEvent) (ControlI, int, int, int32, int32, uint32) {
	mx, my, btnstate := sdl.GetMouseState()
	c := w.FindControl(mx, my)
	cx, cy := mx-c.GetRect().X, my-c.GetRect().Y
	return c, cx, cy, t.X, t.Y, btnstate
}
