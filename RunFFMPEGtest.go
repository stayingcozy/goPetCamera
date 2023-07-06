package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func RunFFMPEGtest(ffmpegLock chan struct{}) {

	// Release the lock when the function exits
	defer func() { <-ffmpegLock }() 


	// Set up the command and arguments
	cmd := exec.Command("ffmpeg", "-re", "-f", "lavfi", "-i", "testsrc=size=640x480:rate=30", "-vcodec", "libvpx", "-cpu-used", "5", "-deadline", "1", "-g", "10", "-error-resilient", "1", "-auto-alt-ref", "1", "-f", "rtp", "rtp://127.0.0.1:5004?pkt_size=1200")

	// Set the output to the console
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("ffmpeg command completed successfully")
}
