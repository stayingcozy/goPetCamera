package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func RunFFMPEG() {
	s := "ffmpeg -re -i /dev/video0 -vcodec libvpx -cpu-used 5 -deadline 1 -g 10 -error-resilient 1 -auto-alt-ref 1 -f rtp 'rtp://127.0.0.1:5004?pkt_size=1200'"

	args := strings.Split(s, " ")

	// Load command
	cmd := exec.Command(args[0], args[1:]...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Run command
	err := cmd.Run()

	// Process errors if there are any
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
	}
}
