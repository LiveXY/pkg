package ffprobe

import (
	"bytes"
	"encoding/json"
	"os/exec"
)

const ffprobeCmd = "ffprobe"

var videoInfoArgs = []string{"", "-of", "json", "-select_streams", "v:0", "-show_entries", "stream=height,width,bit_rate", "-select_streams", "v", "-show_entries", "format=size,duration", "-v", "quiet"}

// 查看数据
type ProbeData struct {
	Streams []Stream `json:"streams"`
	Format  Format   `json:"format"`
}
type Format struct {
	Duration   float64 `json:"duration"`    // 时长
	Size       int64   `json:"size"`        // 大小
	BitRate    int64   `json:"bit_rate"`    // 波特率
	FormatName int64   `json:"format_name"` // 格式名
	StartTime  float64 `json:"start_time"`  // 开始时间
}
type Stream struct {
	FrameRate    string `json:"r_frame_rate"`
	AvgFrameRate string `json:"avg_frame_rate"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	BitRate      int64  `json:"bit_rate"`
	Frames       int64  `json:"nb_frames"`
}

// 查看视频 音频 图片元数据
func Info(fpath string) (data ProbeData, err error) {
	videoInfoArgs[0] = fpath
	bytes, err := ExecCommandOutResult(ffprobeCmd, videoInfoArgs)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes.Bytes(), &data)
	return
}

// 执行命令返回结果
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
