package utils


import(
	"fmt"
	"os/exec"
	"encoding/json"
	"bytes"
	"math"
)

func GetVideoAspectRatio(filePath string) (string, error){
	type FFProbeOutput struct {
		Streams []struct{
			Width int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error to get aspect ratio of the video: %w", err)
	}

	var probe FFProbeOutput
	if err := json.Unmarshal(out.Bytes(), &probe); err != nil {
		return "", fmt.Errorf("error to get ffprobe json: %w", err)
	}

	w := float64(probe.Streams[0].Width)
	h := float64(probe.Streams[0].Height)
	ratio := w/h
	if math.Abs(ratio-(9.0/16.0)) < 0.01{
		return "portrait", nil
	}else if math.Abs(ratio-(16.0/9.0)) < 0.01{
		return "landscape", nil
	}

	return "other", nil
	
}

func ProcessVideoForFastStart(filePath string) (string, error){
	outputFileName := filePath + ".processing"
	cmd := exec.Command("ffmpeg", "-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputFileName)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Error to fast start video: %w", err)
	}
	return outputFileName, nil
}

