package sdlgui

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"

	"github.com/kasworld/htmlcolors"
	"github.com/kasworld/log"
)

func init() {
	ttf.Init()
}

type FontKey struct {
	co   htmlcolors.RGBA
	text string
}
type Font struct {
	Font  *ttf.Font
	Cache map[FontKey]*sdl.Surface
}

func LoadFont(filename string, size int) *Font {
	font, err := ttf.OpenFont(filename, size)
	if err != nil {
		log.Fatalf("font open fail %v\n", err)
	}
	return &Font{
		Font:  font,
		Cache: make(map[FontKey]*sdl.Surface),
	}
}
func (f *Font) MakeSurface(co htmlcolors.RGBA, text string) (*sdl.Surface, sdl.Rect) {
	surface := f.Font.RenderText_Blended(text, RGBA2SDL(co))
	// surface := f.Font.RenderText_Solid(text, co.SdlColor())
	// surface := f.Font.RenderText_Shaded(text, co.SdlColor(), co.Neg().SdlColor())
	if surface == nil {
		log.Printf("%v %v", co, text)
		log.Fatalf("Failed to create surface: %s\n", sdl.GetError())
	}
	srcRect := sdl.Rect{}
	surface.GetClipRect(&srcRect)
	return surface, srcRect
}
func (f *Font) MakeSurface2(fg, bg htmlcolors.RGBA, text string) (*sdl.Surface, sdl.Rect) {
	surface := f.Font.RenderText_Shaded(text, RGBA2SDL(fg), RGBA2SDL(bg))
	if surface == nil {
		// log.Printf("%v %v", co, text)
		log.Fatalf("Failed to create surface: %s\n", sdl.GetError())
	}
	srcRect := sdl.Rect{}
	surface.GetClipRect(&srcRect)
	return surface, srcRect
}
func (f *Font) GetSurfaceWithCache(co htmlcolors.RGBA, text string) (*sdl.Surface, sdl.Rect) {
	surface := f.Cache[FontKey{co, text}]
	if surface != nil {
		srcRect := sdl.Rect{}
		surface.GetClipRect(&srcRect)
		return surface, srcRect
	} else {
		surface, srcRect := f.MakeSurface(co, text)
		f.Cache[FontKey{co, text}] = surface
		return surface, srcRect
	}
}
func (f *Font) MakeTexture(co htmlcolors.RGBA, text string, rend *sdl.Renderer) (*sdl.Texture, sdl.Rect) {
	surface, srcRect := f.MakeSurface(co, text)
	defer surface.Free()
	t, err := rend.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatalf("Failed to create Texture: %s\n", err)
	}
	return t, srcRect
}
