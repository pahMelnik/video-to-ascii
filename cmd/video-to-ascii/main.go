package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/x/term"
	"github.com/pahMelnik/video-to-ascii/internal/decode"
	"github.com/pahMelnik/video-to-ascii/internal/terminal"
	"github.com/pahMelnik/video-to-ascii/internal/video"
	"github.com/schollz/progressbar/v3"
)

// TODO: parameterize file name
// TODO: parameterize log level
func main() {
	var debug bool
	var fileName string
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.StringVar(&fileName, "file", "", "File name")
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	terminalFD := os.Stdout.Fd()

	termWidth, termHeight, err := term.GetSize(terminalFD)
	if err != nil {
		log.Fatal(err)
	}

	videoFramesCount, videoFrameRate, err := video.GetVideoInfo(fileName)
	if err != nil {
		log.Fatal(err)
	}

	terminalFrames := make([]string, videoFramesCount)

	reader := video.GetAllFramesAsJpeg(fileName, termWidth, (termHeight-1)*2, debug)
	images, err := decode.ExtractJPEGsFromMJPEG(reader, videoFramesCount)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("Extracted %d/%d images", len(images), videoFramesCount)

	if debug {
		for i, img := range images {
			//save images to files
			fileName := fmt.Sprintf("./sample-data/frame-%d.jpg", i)
			f, err := os.Create(fileName)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			err = jpeg.Encode(f, img, nil)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	terminalFrameBar := progressbar.Default(int64(videoFramesCount), "Rendering frames")
	for frameNum := range len(images) {
		terminalFrames[frameNum] = terminal.TerminalImage(images[frameNum])
		// increase progress bar
		terminalFrameBar.Add(1)
	}

	// render frames
	// TODO: limit framerate to be same as in original video
	// TODO: add function to get original framerate
	msPerFrame := int64(1000 / videoFrameRate)
	for frameNum, terminalFrame := range terminalFrames {
		start := time.Now()
		// clear previous frame
		if frameNum > 0 {
			terminal.ClearArea(termHeight, termWidth)
		}
		fmt.Print(terminalFrame)
		if time.Since(start).Milliseconds() > msPerFrame {
			time.Sleep(time.Duration(1000 * (msPerFrame - time.Since(start).Milliseconds())))
		}
	}
}
