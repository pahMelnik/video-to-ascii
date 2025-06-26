package video

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// read frame as jpeg
func GetFrameAsJpeg(inFileName string, frameNum, width, height int, showLog bool) io.Reader {
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

// read all frames as jpeg
func GetAllFramesAsJpeg(inFileName string, width, height int, showLog bool) io.Reader {
	ffmpeg.LogCompiledCommand = showLog
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName, ffmpeg.KwArgs{"loglevel": "debug"}).
		Filter("scale", ffmpeg.Args{fmt.Sprint(width), fmt.Sprint(height)}).
		Output("pipe:", ffmpeg.KwArgs{"format": "image2pipe", "vcodec": "mjpeg", "q:v": "1"}).
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
		RFameRate     string `json:"r_frame_rate"` // 10000/3001
	} `json:"streams"`
}

func GetVideoInfo(inFileName string) (int, int, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-count_packets",
		"-show_entries", "stream",
		"-of", "json",
		inFileName,
	)
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get video frame count: %w", err)
	}
	var result FFProbeResult
	if err := json.Unmarshal(output, &result); err != nil {
		return 0, 0, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}
	if len(result.Streams) == 0 {
		return 0, 0, fmt.Errorf("no streams found in ffprobe output")
	}
	frameCount, _ := strconv.Atoi(result.Streams[0].NbReadPackets)
	count, _ := strconv.Atoi(strings.Split(result.Streams[0].RFameRate, "/")[0])
	duration, _ := strconv.Atoi(strings.Split(result.Streams[0].RFameRate, "/")[1])
	frameRate := int(count / duration)
	return frameCount, frameRate, nil
}
