package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/term"
	"github.com/disintegration/imaging"
	"github.com/pahMelnik/video-to-ascii/internal/terminal"
	"github.com/pahMelnik/video-to-ascii/internal/video"
	"github.com/schollz/progressbar/v3"
)

// TODO: parameterize file name
// TODO: parameterize log level
func main() {
	log.SetLevel(log.InfoLevel)
	fileName := "./sample-data/IMG_1135.MOV"
	terminalFD := os.Stdout.Fd()

	termWidth, termHeight, err := term.GetSize(terminalFD)
	if err != nil {
		log.Fatal(err)
	}

	videoFramesCount, err := video.GetVideoFrameCount(fileName)
	if err != nil {
		log.Fatal(err)
	}

	terminalFrames := make([]string, videoFramesCount)
	bar := progressbar.Default(int64(videoFramesCount))

	// PERF: lower speed at end as at start
	for frameNum := range videoFramesCount {
		// FIX: get all frames at once
		reader := video.ReadFrameAsJpeg(fileName, frameNum, termWidth, (termHeight-1)*2, false)
		img, err := imaging.Decode(reader)
		if err != nil {
			log.Fatal(err)
		}
		terminalFrames[frameNum] = terminal.TerminalImage(img)
		// increase progress bar
		bar.Add(1)
	}

	// render frames
	for frameNum, terminalFrame := range terminalFrames {
		// clear previous frame
		if frameNum > 0 {
			terminal.ClearArea(termHeight, termWidth)
		}
		fmt.Print(terminalFrame)
	}
}
