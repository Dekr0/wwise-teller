package aio

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
)

const DownSampleRate = 2000

func FFMPEGDownsample(ctx context.Context, input string) (
	output string, r io.ReadSeekCloser, err error,
) {
	output = input + ".downsample.wav"
	cmd := exec.CommandContext(ctx, "ffmpeg", "-i", input, "-ar:a:0", strconv.FormatUint(DownSampleRate, 10), output)
	var result []byte
	result, err = cmd.CombinedOutput()
	if err != nil {
		slog.Error("Failed to down sampling input wave file to 2KHz", "error", err)
		for line := range bytes.SplitSeq(result, []byte("\n")) {
			slog.Info(string(line))
		}
		return "", nil, err
	}
	r, err = os.Open(output)
	if err != nil {
		return "", nil, err
	}
	return output, r, nil
}
