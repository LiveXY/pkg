package ffprobe

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

const ffprobeCmd = "ffprobe"

var videoInfoArgs = []string{"", "-of", "json", "-select_streams", "v:0", "-show_entries", "stream=height,width,bit_rate", "-select_streams", "v", "-show_entries", "format=size,duration", "-v", "quiet"}

// ProbeData 多媒体文件元数据结果结构体
type ProbeData struct {
	Streams []Stream `json:"streams"`
	Format  Format   `json:"format"`
}

// Format 多媒体容器层信息
type Format struct {
	Duration   float64 `json:"duration"`
	Size       int64   `json:"size"`
	BitRate    int64   `json:"bit_rate"`
	FormatName int64   `json:"format_name"`
	StartTime  float64 `json:"start_time"`
}

// Stream 媒体流（如视频流）信息
type Stream struct {
	FrameRate    string `json:"r_frame_rate"`
	AvgFrameRate string `json:"avg_frame_rate"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	BitRate      int64  `json:"bit_rate"`
	Frames       int64  `json:"nb_frames"`
}

// Info 获取指定视频、音频或图片的元数据
func Info(fpath string) (data ProbeData, err error) {
	videoInfoArgs[0] = fpath
	res, err := ExecCommandOutResult(ffprobeCmd, videoInfoArgs)
	if err != nil {
		return
	}
	err = json.Unmarshal(res.Bytes(), &data)
	return
}

// ExecCommandOutResult 执行系统命令并捕获标准输出
func ExecCommandOutResult(c string, args []string) (out bytes.Buffer, err error) {
	cmd := exec.Command(c, args...)
	var stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &stderr
	err = cmd.Run()
	if err != nil {
		return
	}
	return
}
