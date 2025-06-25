package video

import (
	"bytes"
	"fmt"
	"io"

	"github.com/charmbracelet/log"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ExampleReadFrameAsJpeg(inFileName string, frameNum, width, height int, showLog bool) io.Reader {
	ffmpeg.LogCompiledCommand = showLog
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName, ffmpeg.KwArgs{"loglevel": "debug"}).
		Filter("select", ffmpeg.Args{fmt.Sprintf("eq(n,%d)", frameNum)}).
		Filter("scale", ffmpeg.Args{fmt.Sprint(width), fmt.Sprint(height)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf).
		Run()
	if err != nil {
		log.Fatal(err, err.Error())
	}
	return buf
}
