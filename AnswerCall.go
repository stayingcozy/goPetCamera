package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pion/webrtc/v3"

	"cloud.google.com/go/firestore"
)

func AnswerCall(
	client *firestore.Client,
	uid string, callId string,
	ctx context.Context,
	peerConnection *webrtc.PeerConnection) {

	type Answer struct {
		Type string `firestore:"type"`
		SDP  string `firestore:"sdp"`
	}

	// get call doc
	var callRef = client.Collection("users").Doc(uid).Collection("calls").Doc(callId)
	fmt.Println("call ref \n")
	fmt.Println(callRef)
	fmt.Println("\n")

	// get answer candidates
	var answerCandidates = callRef.Collection("answerCandidates")
	// fmt.Println("answer ref \n")
	// fmt.Println(answerCandidates)
	// fmt.Println("\n")

	// pc.onicecandidate = (event) => {
	// 	event.candidate && setDoc(answerCandidates,event.candidate.toJSON());
	// };

	// Register the onICECandidate event handler
	peerConnection.OnICECandidate(func(event *webrtc.ICECandidate) {
		fmt.Println("OnCandidate triggered")
		// onICECandidateJS(peerConnection, event, ctx)
		if event != nil {
			candidateJSON := event.ToJSON()
			fmt.Println("\n CandidateJSON \n")
			fmt.Println(candidateJSON)
			fmt.Println("\n")
			// setDoc(answerCandidates, candidateJSON)
			newDocRef, _, err := answerCandidates.Add(ctx, candidateJSON)
			if err != nil {
				log.Fatalf("Failed to create document: %v", err)
			}
			log.Printf("Created document with ID: %s", newDocRef.ID)
		}
	})

	callSnap, err := callRef.Get(ctx)
	if err != nil {
		log.Fatalf("Failed to create document: %v", err)
	}
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

			if err = peerConnection.SetRemoteDescription(offer); err != nil {
				panic(err)
			}

			// Create answer
			answerDescription, err := peerConnection.CreateAnswer(nil)
			if err != nil {
				panic(err)
			}

			// Create channel that is blocked until ICE Gathering is complete
			gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

			// Sets the LocalDescription, and starts our UDP listeners
			if err = peerConnection.SetLocalDescription(answerDescription); err != nil {
				panic(err)
			}

			// Block until ICE Gathering is complete, disabling trickle ICE
			// we do this because we only can exchange one signaling message
			// in a production application you should exchange ICE Candidates via OnICECandidate
			<-gatherComplete

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

			// Add offerCandidates to RTCIceCandidate
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

							fmt.Println("\n Here is the sdpMLineIndex after import from fb")
							fmt.Println(sdpMLineIndex)
							fmt.Println("\n")

							ICEcandidate := webrtc.ICECandidateInit{
								Candidate:     candidate, // Assuming the "candidate" field is a string
								SDPMid:        &sdpMid,
								SDPMLineIndex: &sdpMLineIndex,
							}

							err := peerConnection.AddICECandidate(ICEcandidate)
							if err != nil {
								panic(err)
							}
						}
					}
				}
			}
		}
	}
}
