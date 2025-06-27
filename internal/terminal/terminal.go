package terminal

import (
	"image"
	"image/color"
	"os"
	"strconv"
	"strings"
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
	var builder strings.Builder
	height := img.Bounds().Dy()
	width := img.Bounds().Dx()
	for y := 0; y+1 < height; y = y + 2 {
		for x := range width {
			r1, g1, b1 := get8bitcolor(img.At(x, y))
			r2, g2, b2 := get8bitcolor(img.At(x, y+1))

			// Вверхний пиксель
			builder.WriteString(backslash033)
			builder.WriteString("[48;2;")
			builder.WriteString(strconv.Itoa(int(r1)))
			builder.WriteString(";")
			builder.WriteString(strconv.Itoa(int(g1)))
			builder.WriteString(";")
			builder.WriteString(strconv.Itoa(int(b1)))
			builder.WriteString("m")

			// Нижний пиксель
			builder.WriteString(backslash033)
			builder.WriteString("[38;2;")
			builder.WriteString(strconv.Itoa(int(r2)))
			builder.WriteString(";")
			builder.WriteString(strconv.Itoa(int(g2)))
			builder.WriteString(";")
			builder.WriteString(strconv.Itoa(int(b2)))
			builder.WriteString("m")
			builder.WriteString(lowerHalfBlock)
		}
		// Перенос строки
		builder.WriteString(backslash033)
		builder.WriteString("[0m")
		builder.WriteString(backslashN)
	}
	return builder.String()
}

// Перемещает курсор вверх на количество строк и влево на количество колонок
func ClearArea(rows, cols int) {
	// Перемещает курсор влево на количество колонок
	os.Stdout.WriteString(backslash033 + "[" + strconv.Itoa(cols) + "D")
	// Перемещает курсор вверх на количество строк
	os.Stdout.WriteString(backslash033 + "[" + strconv.Itoa(rows) + "A")
}
