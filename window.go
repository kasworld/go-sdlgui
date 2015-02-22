package sdlgui

import (
	"sort"

	"github.com/veandco/go-sdl2/sdl"

	"github.com/kasworld/idgen"
	"github.com/kasworld/log"
	"github.com/kasworld/quadtree"
	"github.com/kasworld/rect"
)

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
