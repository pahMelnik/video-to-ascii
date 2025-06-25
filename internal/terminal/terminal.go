package terminal

import (
	"fmt"
	"image"
)

var lowerHalfBlock = "\u2584"
var backslash033 = "\033"
var backslashN = "\n"

func RenderImage(img image.Image) {
	fmt.Printf("\x1b[%dD", img.Bounds().Dx())
	fmt.Printf("\x1b[%dA", img.Bounds().Dy()/2)

	var result string
	for y := 0; y < img.Bounds().Dy(); y = y + 2 {
		var str string
		for x := range img.Bounds().Dx() {
			r, g, b, _ := img.At(x, y).RGBA()
			// log.Infof("32rgb %d %d %d, 8rgb %d %d %d\n", r, g, b, color32ToColor8(r), color32ToColor8(g), color32ToColor8(b))
			str += RenderPixel(
				color16ToColor8(r),
				color16ToColor8(g),
				color16ToColor8(b),
				true,
			)
			r, g, b, _ = img.At(x, y+1).RGBA()
			str += RenderPixel(
				color16ToColor8(r),
				color16ToColor8(g),
				color16ToColor8(b),
				false,
			)
		}
		str += fmt.Sprintf("%s[0m%s", backslash033, backslashN)
		result += str
	}
	fmt.Print("\033[2J")
	fmt.Print(result)
}

func RenderPixel(r, g, b uint8, upper bool) string {
	var renderStr string
	if upper {
		renderStr = fmt.Sprintf(
			"%s[48;2;%d;%d;%dm",
			backslash033,
			r, g, b,
		)
	} else {
		renderStr = fmt.Sprintf(
			"%s[38;2;%d;%d;%dm%s",
			backslash033,
			r, g, b,
			lowerHalfBlock,
		)
	}
	return renderStr
}

func color16ToColor8(i uint32) uint8 {
	return uint8(i >> 8)
}
