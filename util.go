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
	"github.com/kasworld/go-sdl2/sdl"
	"github.com/kasworld/go-sdl2/sdl_image"
	"github.com/kasworld/htmlcolors"
	"github.com/kasworld/log"
	"github.com/kasworld/rect"
)

func RGBA2SDL(c htmlcolors.RGBA) sdl.Color {
	return sdl.Color{c[0], c[1], c[2], c[3]}
}

func SdlRect2Rect(wr sdl.Rect) rect.Rect {
	return rect.Rect{
		int(wr.X), int(wr.Y),
		int(wr.W), int(wr.H),
	}
}
func Rect2SdlRect(rt rect.Rect) sdl.Rect {
	return sdl.Rect{
		int32(rt.X), int32(rt.Y),
		int32(rt.W), int32(rt.H),
	}
}

func LoadBMP(imageName string) *sdl.Surface {
	image, err := sdl.LoadBMP(imageName)
	if err != nil {
		log.Fatalf("Failed to load BMP: %s", err)
	}
	return image
}
func LoadImage(imageName string) *sdl.Surface {
	image, err := img.Load(imageName)
	if err != nil {
		log.Fatalf("Failed to load Image: %s", err)
	}
	return image
}
