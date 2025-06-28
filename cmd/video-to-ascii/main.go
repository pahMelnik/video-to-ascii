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
	"github.com/pahMelnik/video-to-ascii/package/utils"
	"github.com/schollz/progressbar/v3"
)

// TODO: tui file selector
func main() {
	var debug bool
	var saveFrames bool
	var noShowVideo bool
	var fileName string
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.BoolVar(&saveFrames, "save-frames", false, "Save frames to files")
	flag.BoolVar(&noShowVideo, "no-show-video", false, "Do not show video")
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

	videoInfo, err := video.GetVideoInfo(fileName)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("Video resolution: %dx%d", videoInfo.Width, videoInfo.Height)
	d := utils.Gcd(videoInfo.Width, videoInfo.Height)
	log.Debugf("Video aspect ratio: %d:%d", videoInfo.Width/d, videoInfo.Height/d)

	var videoOutputWidth, videoOutputHeight int
	isVideoHorizontal := videoInfo.Width > videoInfo.Height
	if isVideoHorizontal {
		log.Debugf("Video is horizontal")
		videoOutputWidth = termWidth
		videoOutputHeight = videoOutputWidth * videoInfo.Height / videoInfo.Width
		if videoOutputHeight > termHeight {
			videoOutputHeight = termHeight
			videoOutputWidth = videoOutputHeight * videoInfo.Width / videoInfo.Height
		}
	} else {
		log.Debugf("Video is vertical and terminal is horizontal")
		videoOutputHeight = termHeight
		videoOutputWidth = videoOutputHeight * videoInfo.Width / videoInfo.Height
		if videoOutputWidth > termWidth {
			videoOutputWidth = termWidth
			videoOutputHeight = videoOutputWidth * videoInfo.Height / videoInfo.Width
		}
	}
	log.Debugf("Output resolution: %dx%d", videoOutputWidth, videoOutputHeight)
	d = utils.Gcd(videoOutputWidth, videoOutputHeight)
	log.Debugf("Output aspect ratio: %d:%d", videoOutputWidth/d, videoOutputHeight/d)

	terminalFrames := make([]string, videoInfo.FrameCount)

	videoOutputHeight = (videoOutputHeight - 1) * 2
	reader := video.GetAllFramesAsJpeg(fileName, videoOutputWidth, videoOutputHeight, debug)
	images, err := decode.ExtractJPEGsFromMJPEG(reader, videoInfo.FrameCount)
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("Extracted %d/%d images", len(images), videoInfo.FrameCount)

	if saveFrames {
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

	terminalFrameBar := progressbar.Default(int64(videoInfo.FrameCount), "Rendering frames")
	for frameNum := range len(images) {
		terminalFrames[frameNum] = terminal.TerminalImage(images[frameNum])
		// increase progress bar
		terminalFrameBar.Add(1)
	}

	if !noShowVideo {
		// render frames
		msPerFrame := time.Duration(1000/videoInfo.FPS) * time.Millisecond
		for frameNum, terminalFrame := range terminalFrames {
			start := time.Now()
			// clear previous frame
			if frameNum > 0 {
				terminal.ClearArea(termHeight, termWidth)
			}
			fmt.Print(terminalFrame)
			elapsed := time.Since(start)
			if elapsed < msPerFrame {
				time.Sleep(msPerFrame - elapsed)
			}
		}
	}
}
