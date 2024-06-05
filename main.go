// Consume a RTP stream video UDP, and then send to a WebRTC client.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"path/filepath"
	"math/rand"

	"github.com/pion/webrtc/v3"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {

	const RTP_IP string = "127.0.0.1"
	const RTP_PORT int = 5004
	var firstCheck bool = true

	//// Firebase init /////

	const uid = os.Getenv("UID") // hardcode my uid
	const serviceAccountKey = os.Getenv("SERVICE_ACCOUNT_KEY")
	const projectID = os.Getenv("PROJECT_ID")

	// Full path
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	path_to_serviceAccountKey := filepath.Join(homeDir, "goPetCamera", serviceAccountKey)

	// Initialize Cloud Firestore
	opt := option.WithCredentialsFile(path_to_serviceAccountKey)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Println("error initializing app: %w", err)
	}

	// Initialize the Firebase Admin SDK
	ctx := context.Background()

	//// WebRTC & RTP Server init ////

	// broadcast code //

	// split config values into own var for later func's
	peerConnectionConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// keeping out mediaEngine and Interceptors from broadcast as both use default except for the interval PLI added to interceptor default,
	// so our code will just be default, so just use the default which is just NewPeerConnection with above config var
	// if want video to be seekable (at a cost) then add it

	////

	// peerConnection, err := webrtc.NewPeerConnection(peerConnectionConfig)
	// if err != nil {
	// 	panic(err)
	// }

	// Open a UDP Listener for RTP Packets on port 5004
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP(RTP_IP), Port: RTP_PORT})
	if err != nil {
		panic(err)
	}

	// Increase the UDP receive buffer size
	// Default UDP buffer sizes vary on different operating systems
	bufferSize := 300000 // 300KB
	err = listener.SetReadBuffer(bufferSize)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = listener.Close(); err != nil {
			panic(err)
		}
	}()

	// Create a video track
	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: webrtc.MimeTypeVP8}, "video", "pion")
	if err != nil {
		panic(err)
	}
	// rtpSender, err := peerConnection.AddTrack(videoTrack)
	// if err != nil {
	// 	panic(err)
	// }

	// Read incoming RTCP packets
	// Before these packets are returned they are processed by interceptors. For things
	// like NACK this needs to be called.
	// go func() {
	// 	rtcpBuf := make([]byte, 1500)
	// 	for {
	// 		if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
	// 			return
	// 		}
	// 	}
	// }()

	//// FFMPEG Camera grab to WebRTC run ////
	// go RunFFMPEGtest()
	go RunFFMPEG()

	//// Start HTTP Post to ML server ////
	// go RunFFMPEGposttest()
	// go RunFFMPEGpost()

	//// Listen for RTP packets continously ////
	go ReadRTPtoWebRTC(listener, videoTrack)

	//// WebRTC Connection Status ////

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	// peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
	// 	fmt.Printf("Connection State has changed %s \n", connectionState.String())

	// 	if connectionState == webrtc.ICEConnectionStateFailed {
	// 		if closeErr := peerConnection.Close(); closeErr != nil {
	// 			panic(closeErr)
	// 		}
	// 	}
	// })

	//// Loop on new firebase docs ////

	w := os.Stdout

	fmt.Println("ListenChanges Started...")

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		fmt.Println("firestore.NewClient: %w", err)
	}
	defer client.Close()

	//// Update account with wifi connected ////
	// assumption - wifi checks from runPetCamera allow this code to run + wanted to keep firebase code in one spot
	// Initialize the random number generator with a seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := r.Intn(1000001)
	_, errU := client.Collection("users").Doc(uid).Update(ctx, []firestore.Update{
		{
			Path:  "wifiStatus",
			Value: randomNumber,
		},
	})
	if errU != nil {
		// Handle any errors in an appropriate way, such as returning them.
		log.Printf("An error while updating has occurred: %s", err)
	}


	uidFields := client.Collection("users").Doc(uid).Snapshots(ctx)
	for {
		snap, err := uidFields.Next()
		// DeadlineExceeded will be returned when ctx is cancelled.
		if status.Code(err) == codes.DeadlineExceeded {
			fmt.Println("DeadlineExceeded: %w", err)
		}
		if err != nil {
			fmt.Println("DeadlineExceeded: %w", err)
		}
		if !snap.Exists() {
			fmt.Fprintf(w, "Document no longer exists\n")
		}
		uidData := snap.Data()
		latestCall := uidData["latestCall"].(string)
		fmt.Println(latestCall)

		if !firstCheck {
			time.Sleep(500 * time.Millisecond)
			fmt.Printf("Answering document: %s\n", latestCall)
			// AnswerCall(client, uid, latestCall, ctx, peerConnection)
			AnswerSubscriber(client, uid, latestCall, ctx, peerConnectionConfig, videoTrack)
		}
		firstCheck = false
	}

	////

	// Get calls collection to iterate through
	// calls := client.Collection("users").Doc(uid).Collection("calls").Snapshots(ctx)

	// for {
	// 	snap, err := calls.Next()
	// 	// DeadlineExceeded will be returned when ctx is cancelled.
	// 	if status.Code(err) == codes.DeadlineExceeded {
	// 		fmt.Println("Deadline Exceeded")
	// 	}
	// 	if err != nil {
	// 		fmt.Println("Snapshots.Next: %w", err)
	// 	}
	// 	if snap != nil {
	// 		for _, change := range snap.Changes {
	// 			switch change.Kind {
	// 			case firestore.DocumentAdded:

	// 				fmt.Fprintf(w, "New call: %s\n", change.Doc.Ref.ID)

	// 				var callId = change.Doc.Ref.ID

	// 				// Function to handle the document addition
	// 				fmt.Printf("Answering document: %s\n", callId)

	// 				AnswerCall(client, uid, callId, ctx, peerConnection)
	// 			}
	// 		}
	// 	}
	// }


}
