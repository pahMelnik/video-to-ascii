package main

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/term"
	"github.com/disintegration/imaging"
	"github.com/pahMelnik/video-to-ascii/internal/terminal"
	"github.com/pahMelnik/video-to-ascii/internal/video"
)

func main() {
	fileName := "./sample-data/IMG_1135.MOV"
	terminalFD := os.Stdout.Fd()

	termWidth, termHeight, err := term.GetSize(terminalFD)
	if err != nil {
		log.Fatal(err)
	}

	for frameNum := 0; frameNum <= 200; frameNum += 5 {
		reader := video.ExampleReadFrameAsJpeg(fileName, frameNum, termWidth, (termHeight-1)*2, false)
		img, err := imaging.Decode(reader)
		if err != nil {
			log.Fatal(err)
		}
		terminal.RenderImage(img)
	}
}
