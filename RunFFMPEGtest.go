package main


import (
	"os/exec"
	"strings"
	"bytes"
	"fmt"
	// "math/rand"
	// "strconv"
)

func RunFFMPEGtest() {
	s := "ffmpeg -re -f lavfi -i testsrc=size=640x480:rate=30 -vcodec libvpx -cpu-used 5 -deadline 1 -g 10 -error-resilient 1 -auto-alt-ref 1 -f rtp 'rtp://127.0.0.1:5004?pkt_size=1200'"
	// s := "ffmpeg -re -i /dev/video0 -vcodec libvpx -cpu-used 5 -deadline 1 -g 10 -error-resilient 1 -auto-alt-ref 1 -f rtp 'rtp://127.0.0.1:5004?pkt_size=1200'"
	// rand_num := rand.Int()
	// s := "ffmpeg -f lavfi -i testsrc=size=640x480:rate=30 -t 10 -c:v libx264 -preset ultrafast -tune zerolatency output_" + strconv.Itoa(rand_num) + ".mp4"

	args := strings.Split(s, " ")
	// args := []string{
	// 	"ffmpeg","-re","-f lavfi","-i testsrc=size=640x480:rate=30",
	// 	"-vcodec libvpx","-cpu-used 5","-deadline 1","-g 10","-error-resilient 1",
	// 	"-auto-alt-ref 1","-f rtp","'rtp://127.0.0.1:5004?pkt_size=1200'"}

	// Load command
	cmd := exec.Command(args[0], args[1:]...)
	// cmd := exec.Command("ffmpeg","-re -f lavfi -i testsrc=size=640x480:rate=30 -vcodec libvpx -cpu-used 5 -deadline 1 -g 10 -error-resilient 1 -auto-alt-ref 1 -f rtp 'rtp://127.0.0.1:5004?pkt_size=1200'")
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