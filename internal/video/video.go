package video

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"

	"github.com/charmbracelet/log"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func ReadFrameAsJpeg(inFileName string, frameNum, width, height int, showLog bool) io.Reader {
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

type FFProbeResult struct {
	Streams []struct {
		NbReadPackets string `json:"nb_read_packets"`
	} `json:"streams"`
}

func GetVideoFrameCount(inFileName string) (int, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-count_packets",
		"-show_entries", "stream=nb_read_packets",
		"-of", "json",
		inFileName,
	)
	output, err := cmd.Output()
	log.Debugf("ffprobe output: %s", output)
	if err != nil {
		return 0, fmt.Errorf("failed to get video frame count: %w", err)
	}
	var result FFProbeResult
	if err := json.Unmarshal(output, &result); err != nil {
		return 0, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}
	if len(result.Streams) == 0 {
		return 0, fmt.Errorf("no streams found in ffprobe output")
	}
	frameCount, err := strconv.Atoi(result.Streams[0].NbReadPackets)
	if err != nil {
		return 0, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}
	return frameCount, nil
}
