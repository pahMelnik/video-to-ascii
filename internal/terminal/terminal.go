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

// Перемещает курсор вверх на `rows` строк
func MoveCursorUp(rows int) {
	os.Stdout.WriteString(backslash033 + "[" + strconv.Itoa(rows) + "A")
}

// Перемещает курсор вниз на `rows` строк
func MoveCursorDown(rows int) {
	os.Stdout.WriteString(backslash033 + "[" + strconv.Itoa(rows) + "B")
}

// Перемещает курсор влево на `cols` колонок
func MoveCursorLeft(cols int) {
	os.Stdout.WriteString(backslash033 + "[" + strconv.Itoa(cols) + "D")
}

// Перемещает курсор вправо на `cols` колонок
func MoveCursorRight(cols int) {
	os.Stdout.WriteString(backslash033 + "[" + strconv.Itoa(cols) + "C")
}

// Перемещает курсор на позицию `row`, `col`
func MoveCursorTo(row, col int) {
	// Как будто вместо H можно использовать f
	os.Stdout.WriteString(backslash033 + "[" + strconv.Itoa(row) + ";" + strconv.Itoa(col) + "H")
}

// Перемещает курсор в координаты 0,0
func MoveCursorToHome() {
	os.Stdout.WriteString(backslash033 + "[H")
}

// Перемещает курсор в начало следующей строки после перемещения курсора вниз на `rows` строк
func MoveCursorToNextLineBegining(rows int) {
	os.Stdout.WriteString(backslash033 + "[" + strconv.Itoa(rows) + "E")
}

// Перемещает курсор в начало предыдущей строки после перемещения курсора вверх на `rows` строк
func MoveCursorToPreviousLineBegining(rows int) {
	os.Stdout.WriteString(backslash033 + "[" + strconv.Itoa(rows) + "F")
}
