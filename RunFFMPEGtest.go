package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func RunFFMPEGtest() {
	
	// Set up the command and arguments
	cmd := exec.Command("ffmpeg", "-re", "-f", "lavfi", "-i", "testsrc=size=640x480:rate=30", "-vcodec", "libvpx", "-cpu-used", "5", "-deadline", "1", "-g", "10", "-error-resilient", "1", "-auto-alt-ref", "1", "-f", "rtp", "rtp://127.0.0.1:5004?pkt_size=1200")
	// for audio
	// cmd := exec.Command("ffmpeg", "-f", "lavfi", "-i", "sine=frequency=1000", "-c:a", "libopus", "-b:a", "48000", "-sample_fmt", "s16p", "-ssrc", "1", "-payload_type", "111", "-f", "rtp", "-max_delay", "0", "-application", "lowdelay", "rtp://127.0.0.1:5004?pkt_size=1200")
	//ffmpeg -f lavfi -i 'sine=frequency=1000' -c:a libopus -b:a 48000 -sample_fmt s16p -ssrc 1 -payload_type 111 -f rtp -max_delay 0 -application lowdelay 'rtp://127.0.0.1:5004?pkt_size=1200'

	// for video and audio FAILED
	// cmd := exec.Command("ffmpeg", "-re", "-f", "lavfi", "-i", "testsrc=size=640x480:rate=30", "-f", "lavfi", "-i", "sine=frequency=1000", "-vf", "format=yuv420p", "-vcodec", "libvpx", "-cpu-used", "5", "-deadline", "1", "-g", "10", "-error-resilient", "1", "-auto-alt-ref", "1", "-c:a", "libopus", "-b:a", "48000", "-sample_fmt", "s16p", "-ssrc", "1", "-payload_type", "111", "-f", "rtp", "-max_delay", "0", "-application", "lowdelay", "rtp://127.0.0.1:5004?pkt_size=1200")
	// ffmpeg -re -f lavfi -i testsrc=size=640x480:rate=30 -f lavfi -i 'sine=frequency=1000' -vf "format=yuv420p" -vcodec libvpx -cpu-used 5 -deadline 1 -g 10 -error-resilient 1 -auto-alt-ref 1 -c:a libopus -b:a 48000 -sample_fmt s16p -ssrc 1 -payload_type 111 -f rtp -max_delay 0 -application lowdelay 'rtp://127.0.0.1:5004?pkt_size=1200'

	// cmd := exec.Command("ffmpeg",
	// "-re",
	// "-f", "lavfi", "-i", "testsrc=size=640x480:rate=30",
	// "-f", "lavfi", "-i", "sine=frequency=1000",
	// "-map", "0:v", "-map", "1:a",
	// "-c:", "copy",
	// "-f", "rtp", "rtp://127.0.0.1:5004?pkt_size=1200",
	// )

	// ffmpeg -i INPUT_FILE.mp4 -i AUDIO.aac -c:v copy -c:a copy OUTPUT_FILE.mp4

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
