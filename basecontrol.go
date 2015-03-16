package sdlgui

import (
	"github.com/kasworld/go-sdl2/sdl"
	"github.com/kasworld/idgen"
	"github.com/kasworld/log"
	"github.com/kasworld/rect"
)

const (
	Rmask = 0x000000ff
	Gmask = 0x0000ff00
	Bmask = 0x00ff0000
	Amask = 0xff000000
)

// QuadTreeObjI
func (c *Control) GetRect() rect.Rect {
	return rect.Rect{
		c.X, c.Y,
		c.W, c.H,
	}
}
func (c *Control) GetID() int64 {
	return c.ID
}

// ControlI
func (c *Control) UpdateContents() {
}
func (c *Control) DrawSurface() {
}
func (c *Control) UpdateToWindow() {
	t, err := c.Win.Rend.CreateTextureFromSurface(c.Suf)
	if err != nil {
		log.Fatalf("Failed to create Texture: %s\n", err)
	}
	defer t.Destroy()
	srcrect := sdl.Rect{}
	c.Suf.GetClipRect(&srcrect)
	dstrect := sdl.Rect{
		int32(c.X), int32(c.Y), srcrect.W, srcrect.H,
	}
	c.Win.Rend.SetClipRect(c.GetClipRect())
	c.Win.Rend.Copy(t, &srcrect, &dstrect)
	c.DrawBorder()

	c.Win.Rend.SetClipRect(nil)
}
func (c *Control) SetWindow(w *Window) {
	c.Win = w
}
func (c *Control) MouseOver(x, y int, btnstate uint32) {
}
func (c *Control) MouseButton(x, y int, btnnum uint8, btnstate uint8) {
}
func (c *Control) MouseWheel(x, y int, dx int32, dy int32, btnstate uint32) {
}
func (c *Control) GetZ() int {
	return c.Z
}
func (c *Control) Show(visible bool) {
	c.Visible = visible
	if visible {
		c.Win.AddUpdateControl(c)
	} else {
		c.Win.UpdateRect(c.GetRect())
	}
	// log.Info("show %v", c.GetRect())
}
func (c *Control) IsVisible() bool {
	return c.Visible
}
func (c *Control) IsTransparent() bool {
	return false
}

type Control struct {
	ID              int64
	Win             *Window
	X, Y, Z         int
	W, H            int
	Suf             *sdl.Surface
	Rend            *sdl.Renderer
	Visible         bool
	ContentsChanged bool
	BorderSize      int32
	BorderType      int
}

func NewControl(x, y, z int, wx, wy int) *Control {
	c := Control{
		ID:              <-idgen.GenCh(),
		X:               x,
		Y:               y,
		Z:               z,
		W:               wx,
		H:               wy,
		Visible:         true,
		ContentsChanged: true,
		BorderSize:      0,
		BorderType:      0,
	}
	var err error
	c.Suf, err = sdl.CreateRGBSurface(0, int32(c.W), int32(c.H), 32,
		Rmask, Gmask, Bmask, Amask)
	if err != nil {
		log.Fatalf("Failed to create surface: %s\n", err)
	}
	c.Rend, err = sdl.CreateSoftwareRenderer(c.Suf)
	if err != nil {
		log.Fatalf("Failed to create renderer: %s\n", err)
	}
	return &c
}
func (c *Control) GetClipRect() *sdl.Rect {
	clipRect := Rect2SdlRect(c.GetRect())
	// clipRect.X += int32(c.X)
	// clipRect.Y += int32(c.Y)
	return &clipRect
}
func (c *Control) GetClientConnRect() rect.Rect {
	return rect.Rect{
		0, 0,
		c.W, c.H,
	}
}
func (c *Control) Cleanup() {
	if c.Suf != nil {
		c.Suf.Free()
	}
	if c.Rend != nil {
		c.Rend.Destroy()
	}
}
func (c *Control) MakeTexture(surface *sdl.Surface) *sdl.Texture {
	t, err := c.Rend.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatalf("Failed to create Texture: %s\n", err)
	}
	return t
}
func (c *Control) DrawOverlayFilter(colors OverlayFilter) {
	c.Win.Rend.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	for i, v := range colors {
		if i >= c.W/2 || i >= c.H/2 {
			break
		}
		r := sdl.Rect{int32(c.X + i), int32(c.Y + i), int32(c.W - i*2), int32(c.H - i*2)}
		c.Win.Rend.SetDrawColor(v[0], v[1], v[2], v[3])
		c.Win.Rend.DrawRect(&r)
	}
}
func (c *Control) DrawBorder() {
	switch c.BorderType {
	case 0:
	case 1:
		c.DrawBoderBlackMild()
	case 2:
		c.DrawBorderBlackHard()
	case 3:
		c.DrawBorderWhiteHard()
	default:
		log.Error("unknown bordertype %v", c.BorderType)
	}
}
func (c *Control) DrawBorderBlackHard() {
	c.Win.Rend.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	for i := int32(0); i < c.BorderSize; i++ {
		r := sdl.Rect{int32(c.X) + i, int32(c.Y) + i, int32(c.W) - i*2, int32(c.H) - i*2}
		c.Win.Rend.SetDrawColor(0, 0, 0, uint8(255-i*256/c.BorderSize))
		c.Win.Rend.DrawRect(&r)
	}
}
func (c *Control) DrawBoderBlackMild() {
	c.Win.Rend.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	for i := int32(0); i < c.BorderSize; i++ {
		r := sdl.Rect{int32(c.X) + i, int32(c.Y) + i, int32(c.W) - i*2, int32(c.H) - i*2}
		c.Win.Rend.SetDrawColor(0, 0, 0, uint8(127-i*128/c.BorderSize))
		c.Win.Rend.DrawRect(&r)
	}
}
func (c *Control) DrawBorderWhiteHard() {
	c.Win.Rend.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	for i := int32(0); i < c.BorderSize; i++ {
		r := sdl.Rect{int32(c.X) + i, int32(c.Y) + i, int32(c.W) - i*2, int32(c.H) - i*2}
		c.Win.Rend.SetDrawColor(255, 255, 255, uint8(255-i*256/c.BorderSize))
		c.Win.Rend.DrawRect(&r)
	}
}
