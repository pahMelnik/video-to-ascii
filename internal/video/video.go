package video

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

// Получает кадр в формате jpeg
func GetFrameAsJpeg(inFileName string, frameNum, width, height int, showLog bool) (io.Reader, error) {
	ffmpeg.LogCompiledCommand = showLog
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName, ffmpeg.KwArgs{"loglevel": "debug"}).
		Filter("select", ffmpeg.Args{fmt.Sprintf("eq(n,%d)", frameNum)}).
		Filter("scale", ffmpeg.Args{fmt.Sprint(width), fmt.Sprint(height)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf).
		Run()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// Получает все кадры в формате jpeg
func GetAllFramesAsJpeg(inFileName string, width, height int, showLog bool) (io.Reader, error) {
	ffmpeg.LogCompiledCommand = showLog
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName, ffmpeg.KwArgs{"loglevel": "debug"}).
		Filter("scale", ffmpeg.Args{fmt.Sprint(width), fmt.Sprint(height)}).
		Output("pipe:", ffmpeg.KwArgs{"format": "image2pipe", "vcodec": "mjpeg", "q:v": "1"}).
		WithOutput(buf).
		Run()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// Получает аудио поток из видео в формате mp3
func GetAudioFromVideo(inFileName string, showLog bool) (io.Reader, error) {
	ffmpeg.LogCompiledCommand = showLog
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName, ffmpeg.KwArgs{"loglevel": "debug"}).
		Output("pipe:", ffmpeg.KwArgs{"q:a": "0", "map": "a", "f": "mp3"}).
		WithOutput(buf).
		Run()
	if err != nil {
		return nil, err
	}
	return buf, nil
}

type FFProbeResult struct {
	Streams []struct {
		NbReadPackets string `json:"nb_read_packets"`
		RFameRate     string `json:"r_frame_rate"` // 10000/3001
		Width         int    `json:"width"`
		Height        int    `json:"height"`
	} `json:"streams"`
}

type VideoInfo struct {
	FrameCount int
	FPS        int
	Width      int
	Height     int
}

func GetVideoInfo(inFileName string) (VideoInfo, error) {
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
		return VideoInfo{}, fmt.Errorf("failed to get video frame count: %w", err)
	}
	var result FFProbeResult
	if err := json.Unmarshal(output, &result); err != nil {
		return VideoInfo{}, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}
	if len(result.Streams) == 0 {
		return VideoInfo{}, fmt.Errorf("no streams found in ffprobe output")
	}
	frameCount, _ := strconv.Atoi(result.Streams[0].NbReadPackets)
	count, _ := strconv.Atoi(strings.Split(result.Streams[0].RFameRate, "/")[0])
	duration, _ := strconv.Atoi(strings.Split(result.Streams[0].RFameRate, "/")[1])
	frameRate := int(count / duration)
	width := result.Streams[0].Width
	height := result.Streams[0].Height

	return VideoInfo{
			FrameCount: frameCount,
			FPS:        frameRate,
			Width:      width,
			Height:     height,
		},
		nil
}
