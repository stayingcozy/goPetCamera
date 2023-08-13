package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/pion/webrtc/v3"

	"cloud.google.com/go/firestore"
)

func AnswerSubscriber(
	client *firestore.Client,
	uid string, callId string,
	ctx context.Context,
	peerConnectionConfig webrtc.Configuration,
	localTrack *webrtc.TrackLocalStaticRTP) { // callsId string

	peerConnection, err := webrtc.NewPeerConnection(peerConnectionConfig)
	if err != nil {
		panic(err)
	}

	rtpSender, err := peerConnection.AddTrack(localTrack)
	if err != nil {
		panic(err)
	}

	// Read incoming RTCP packets
	// Before these packets are returned they are processed by interceptors. For things
	// like NACK this needs to be called.
	go func() {
		rtcpBuf := make([]byte, 1500)
		for {
			if _, _, rtcpErr := rtpSender.Read(rtcpBuf); rtcpErr != nil {
				return
			}
		}
	}()

	// get call doc
	var callRef = client.Collection("users").Doc(uid).Collection("calls").Doc(callId)
	// var callRef = client.Collection("users").Doc(uid).Collection("calls").Doc(callsId).Collection("call").Doc(callId)
	fmt.Println("call ref id: ")
	fmt.Println(callId)

	// get answer candidates
	// var answerCandidates = callRef.Collection("answerCandidates")
	answerCandidates := callRef.Collection("answerCandidates").NewDoc()

	// fmt.Println("answer ref \n")
	// fmt.Println(answerCandidates)
	// fmt.Println("\n")

	// pc.onicecandidate = (event) => {
	// 	event.candidate && setDoc(answerCandidates,event.candidate.toJSON());
	// };

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())

		// TODO might add one for disconnect

		if connectionState == webrtc.ICEConnectionStateFailed {
			if closeErr := peerConnection.Close(); closeErr != nil {
				panic(closeErr)
			}
		}
	})

	// Register the onICECandidate event handler
	peerConnection.OnICECandidate(func(event *webrtc.ICECandidate) {
		// onICECandidateJS(peerConnection, event, ctx)
		if event != nil {
			candidateJSON := event.ToJSON()
			// fmt.Println("CandidateJSON")
			// fmt.Println(candidateJSON)

			// Accessing individual fields
			// fmt.Println("Candidate:", candidateJSON.Candidate)
			// fmt.Println("SDPMid:", *candidateJSON.SDPMid)
			// fmt.Println("SDPMLineIndex:", *candidateJSON.SDPMLineIndex)

			// answercandidateFB := AnswerCandidateFB{
			// 	Candidate: candidateJSON.Candidate,
			// 	SDPMLineIndex: int64(*candidateJSON.SDPMLineIndex),
			// 	SDPMid: "0",
			// }
			// fmt.Println("Candidate:", answercandidateFB.Candidate)
			// fmt.Println("SDPMid:", answercandidateFB.SDPMid)
			// fmt.Println("SDPMLineIndex:", answercandidateFB.SDPMLineIndex)

			// Convert to readable format on Firebase for User on WebRTC js side to process
			answercandidateFB := AnswerCandidateFB{
				Candidate:     candidateJSON.Candidate,
				SDPMLineIndex: int64(*candidateJSON.SDPMLineIndex),
				SDPMid:        strconv.Itoa(int(*candidateJSON.SDPMLineIndex)), // value is blank, so use SDPMLineIndex for now
			}
			// fmt.Println("Candidate:", answercandidateFB.Candidate)
			// fmt.Println("SDPMid:", answercandidateFB.SDPMLineIndex)
			// fmt.Println("SDPMLineIndex:", answercandidateFB.SDPMid)

			// setDoc(answerCandidates, candidateJSON)

			// newDocRef, _, err := answerCandidates.Add(ctx, candidateJSON)
			// if err != nil {
			// 	log.Fatalf("Failed to create document: %v", err)
			// }

			_, err := answerCandidates.Set(ctx, answercandidateFB)
			if err != nil {
				// Handle any errors in an appropriate way, such as returning them.
				log.Printf("An error has occurred: %s", err)
			}

			// log.Printf("Created document with ID: %s", newDocRef.ID)

		}
	})

	// Stop code until doc snap is taken of fb doc
	callSnapChan := make(chan *firestore.DocumentSnapshot)
	go GetDoc(callRef, ctx, callSnapChan)
	callSnap := <-callSnapChan

	callData := callSnap.Data()
	// fmt.Printf("Document data: %#v\n", callData)

	offerDescription := callData["offer"]
	// fmt.Println(offerDescription)

	// Break down interface into sdp
	if value, ok := offerDescription.(map[string]interface{}); ok {
		fieldValue := value["sdp"]
		if sdpStrValue, ok := fieldValue.(string); ok {

			// fmt.Println("\n Printing string value: \n ")
			// fmt.Println(sdpStrValue)

			// Set the remote SessionDescription
			offer := webrtc.SessionDescription{
				Type: webrtc.SDPTypeOffer,
				SDP:  sdpStrValue,
			}

			// Create channel that is blocked until ICE Gathering is complete
			fmt.Println("Starting setRemoteDescription")
			// gatherCompleteRem := webrtc.GatheringCompletePromise(peerConnection)
			var err error
			if err = peerConnection.SetRemoteDescription(offer); err != nil {
				panic(err)
			}
			// <-gatherCompleteRem
			fmt.Println("Ending setRemoteDescription")

			/////////

			// Create answer
			answerDescription, err := peerConnection.CreateAnswer(nil)
			if err != nil {
				panic(err)
			}

			fmt.Println("Starting SetLocalDescription Promise")
			// Create channel that is blocked until ICE Gathering is complete
			// gatherCompleteLocal := webrtc.GatheringCompletePromise(peerConnection)

			// Sets the LocalDescription, and starts our UDP listeners
			if err = peerConnection.SetLocalDescription(answerDescription); err != nil {
				panic(err)
			}

			// Block until ICE Gathering is complete, disabling trickle ICE
			// we do this because we only can exchange one signaling message
			// in a production application you should exchange ICE Candidates via OnICECandidate
			// <-gatherCompleteLocal
			fmt.Println("Ending SetLocalDescription Promise")

			answer := Answer{
				Type: "answer",
				SDP:  answerDescription.SDP,
			}

			// callRef.Update(ctx, []firestore.Update{ { answer } })
			_, Uerr := callRef.Update(ctx, []firestore.Update{
				{
					Path:  "answer",
					Value: answer,
				},
			})
			if Uerr != nil {
				// Handle any errors in an appropriate way, such as returning them.
				log.Printf("An error has occurred: %s", Uerr)
			}

			/////////
			// Add offerCandidates to RTCIceCandidate
			fmt.Println("\nChecking docs in offerCandidates....")
			oCit := callRef.Collection("offerCandidates").Snapshots(ctx)
			for {
				snap, err := oCit.Next()
				if err != nil {
					panic(err)
				}
				if snap != nil {
					for _, change := range snap.Changes {
						switch change.Kind {
						case firestore.DocumentAdded:
							oCdata := change.Doc.Data()
							var candidate string = oCdata["candidate"].(string)
							var sdpMid string = oCdata["sdpMid"].(string)
							var sdpMLineIndex uint16 = uint16(oCdata["sdpMLineIndex"].(int64))

							fmt.Println("\n Here is the candidate after import from fb")
							fmt.Println(candidate)
							fmt.Println("\n Here is the sdpMid after import from fb")
							fmt.Println(sdpMid)
							fmt.Println("\n Here is the sdpMLineIndex after import from fb")
							fmt.Println(sdpMLineIndex)

							// If candidate && sdpMid && sdpMLineIndex exist then add ICE
							if candidate != "" {
								fmt.Println("adding ICE")
								ICEcandidate := webrtc.ICECandidateInit{
									Candidate:     candidate, // Assuming the "candidate" field is a string
									SDPMid:        &sdpMid,
									SDPMLineIndex: &sdpMLineIndex,
								}

								err := peerConnection.AddICECandidate(ICEcandidate)
								if err != nil {
									panic(err)
								}
								fmt.Println("post ICE added")
								break
							}
							candidate = ""
						}
					}
					break
				}
			}
			/////

			// // Create answer
			// answerDescription, err := peerConnection.CreateAnswer(nil)
			// if err != nil {
			// 	panic(err)
			// }

			// fmt.Println("Starting SetLocalDescription Promise")
			// // Create channel that is blocked until ICE Gathering is complete
			// // gatherCompleteLocal := webrtc.GatheringCompletePromise(peerConnection)

			// // Sets the LocalDescription, and starts our UDP listeners
			// if err = peerConnection.SetLocalDescription(answerDescription); err != nil {
			// 	panic(err)
			// }

			// // Block until ICE Gathering is complete, disabling trickle ICE
			// // we do this because we only can exchange one signaling message
			// // in a production application you should exchange ICE Candidates via OnICECandidate
			// // <-gatherCompleteLocal
			// fmt.Println("Ending SetLocalDescription Promise")

			// answer := Answer{
			// 	Type: "answer",
			// 	SDP:  answerDescription.SDP,
			// }

			// // callRef.Update(ctx, []firestore.Update{ { answer } })
			// _, Uerr := callRef.Update(ctx, []firestore.Update{
			// 	{
			// 		Path:  "answer",
			// 		Value: answer,
			// 	},
			// })
			// if Uerr != nil {
			// 	// Handle any errors in an appropriate way, such as returning them.
			// 	log.Printf("An error has occurred: %s", Uerr)
			// }
		}
	}
}
