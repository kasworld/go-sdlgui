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
	"fmt"

	"github.com/kasworld/go-sdl2/sdl"
	"github.com/kasworld/go-sdl2/sdl_ttf"
	"github.com/kasworld/htmlcolors"
)

func init() {
	// println("init ttf")
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
		fmt.Printf("font open fail %v\n", err)
	}
	return &Font{
		Font:  font,
		Cache: make(map[FontKey]*sdl.Surface),
	}
}
func (f *Font) MakeSurface(co htmlcolors.RGBA, text string) (*sdl.Surface, sdl.Rect) {
	surface, err := f.Font.RenderUTF8_Blended(text, RGBA2SDL(co))
	if err != nil {
		fmt.Printf("%v %v %v", err, co, text)
		fmt.Printf("Failed to create surface: %s\n", sdl.GetError())
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
		fmt.Printf("Failed to create Texture: %s\n", err)
	}
	return t, srcRect
}
