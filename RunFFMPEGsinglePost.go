package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os/exec"
)

func RunFFMPEGsinglePost() {

	url := "http://localhost:3030/upload" // Replace with your desired URL

	// Set up the command and arguments
	cmd := exec.Command("ffmpeg", "-re", "-f", "lavfi", "-i", "testsrc=size=640x480", "-ss", "00:00:01", "-vframes", "1", "-f", "image2pipe", "-")

	frameOutput, err := cmd.Output()
	if err != nil {
		fmt.Println("Error capturing video frame:", err)
		return
	}

	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a new form field and add the frame image
	fileField, err := writer.CreateFormFile("image", "frame.jpg")
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	// Copy the frame output to the form field
	_, err = io.Copy(fileField, bytes.NewReader(frameOutput))
	if err != nil {
		fmt.Println("Error copying frame output to form field:", err)
		return
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return
	}

	// Create a new HTTP POST request
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the content type of the request
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer response.Body.Close()

	// Process response message // 

	// Read the response body
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Check the response status code
	if response.StatusCode != http.StatusOK {
		fmt.Println("Server returned an error:", response.Status)
		return
	}

	// Convert the response body to a string
	responseMessage := string(responseBody)

	// Print the response message
	fmt.Println("Response:", responseMessage)

}
