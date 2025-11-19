package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type Recorder struct {
	storage []any
}

func NewRecorder() *Recorder {
	return &Recorder{}
}

func (r *Recorder) getMonitorSource() (string, error) {
	listCmd := exec.Command("pactl", "list", "sinks")
	grepCmd := exec.Command("grep", ".monitor")

	pr, pw := io.Pipe()
	defer pr.Close()

	listCmd.Stdout = pw
	grepCmd.Stdin = pr

	var outBuf bytes.Buffer
	grepCmd.Stdout = &outBuf

	if err := grepCmd.Start(); err != nil {
		pw.Close()
		return "", err
	}

	if err := listCmd.Run(); err != nil {
		pw.Close()
		grepCmd.Wait()
		return "", err
	}

	// close writer to signal EOF to grep
	pw.Close()

	if err := grepCmd.Wait(); err != nil {
		return "", err
	}

	output := strings.TrimSpace(outBuf.String())
	if output == "" {
		return "", errors.New("no monitor sink found")
	}

	// use the first matching line
	line := output
	if i := strings.IndexByte(output, '\n'); i >= 0 {
		line = output[:i]
	}

	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("failed to parse pactl output: %q", line)
	}

	monitorSource := strings.TrimSpace(parts[1])
	if monitorSource == "" {
		return "", errors.New("no monitor sink found")
	}

	return monitorSource, nil
}

func (r *Recorder) Record() error {
	monitorSource, err := r.getMonitorSource()
	if err != nil {
		return err
	}

	cmd := exec.Command("ffmpeg",
		"-f", "pulse",
		"-i", monitorSource,
		"-t", "10",
		"output.wav",
	)
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}
