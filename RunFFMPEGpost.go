package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os/exec"
)

func RunFFMPEGpost() {
	url := "http://localhost:3030/upload" // Replace with the URL of your Node.js server
	const imageBytes int = 640 * 480      // image H x W x Color Channels = # of bytes in image

	// Execute the ffmpeg command to capture images
	cmd := exec.Command("ffmpeg", "-re", "-f", "lavfi", "-i", "testsrc=size=640x480", "-vf", "fps=2", "-f", "image2pipe", "-")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}

	// Start the ffmpeg command
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting ffmpeg command:", err)
		return
	}
	defer cmd.Wait()

	// Read and send the captured images
	for {
		frameOutput := make([]byte, imageBytes) // Adjust the buffer size as needed
		n, err := stdout.Read(frameOutput)
		if err != nil {
			if err == io.EOF {
				break // End of output
			}
			fmt.Println("Error reading frame output:", err)
			break
		}

		// fmt.Println("Length of image in bytes")
		// fmt.Println(n)

		// Create a new request with the frame image data
		req, err := http.NewRequest("POST", url, bytes.NewReader(frameOutput[:n]))
		if err != nil {
			fmt.Println("Error creating request:", err)
			break
		}
		req.Header.Set("Content-Type", "image/jpeg")

		// Send the request with the frame image
		client := http.DefaultClient
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending HTTP POST request:", err)
			break
		}
		defer resp.Body.Close()

		// Process the response if needed
		// ...

		// fmt.Println("Image uploaded successfully!")
	}
}
