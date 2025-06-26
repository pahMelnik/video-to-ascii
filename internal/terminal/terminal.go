package terminal

import (
	"fmt"
	"image"
	"image/color"
)

var lowerHalfBlock = "\u2584"
var backslash033 = "\033"
var backslashN = "\n"

// Получает цвет из картинки и возвращает его в виде 8-битного цвета
func get8bitcolor(color color.Color) (uint8, uint8, uint8) {
	r, g, b, _ := color.RGBA()
	// Каждый цвет записывается в 16 бит, смещением мы берем старшие 8 бит
	return uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)
}

// Возвращает строку для отрисовки картинки в терминале
func TerminalImage(img image.Image) string {
	var result string
	var r, g, b uint8
	for y := 0; y < img.Bounds().Dy(); y = y + 2 {
		var str string
		for x := range img.Bounds().Dx() {
			r, g, b = get8bitcolor(img.At(x, y))
			str += renderPixel(r, g, b, true)
			r, g, b = get8bitcolor(img.At(x, y+1))
			str += renderPixel(r, g, b, false)
		}
		str += fmt.Sprintf("%s[0m%s", backslash033, backslashN)
		result += str
	}
	return result
}

// Перемещает курсор вверх на количество строк и влево на количество колонок
func ClearArea(rows, cols int) {
	// Перемещает курсор влево на количество колонок
	fmt.Printf("%s[%dD", backslash033, cols)
	// Перемещает курсор вверх на количество строк
	fmt.Printf("%s[%dA", backslash033, rows)
}

// Возвращает строку для отрисовки пикселя в терминале
func renderPixel(r, g, b uint8, upper bool) string {
	var pixelStr string
	if upper {
		pixelStr = fmt.Sprintf(
			"%s[48;2;%d;%d;%dm",
			backslash033,
			r, g, b,
		)
	} else {
		pixelStr = fmt.Sprintf(
			"%s[38;2;%d;%d;%dm%s",
			backslash033,
			r, g, b,
			lowerHalfBlock,
		)
	}
	return pixelStr
}
