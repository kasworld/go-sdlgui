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
