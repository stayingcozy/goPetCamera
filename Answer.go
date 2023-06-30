package main     

type Answer struct {
	Type string `firestore:"type"`
	SDP  string `firestore:"sdp"`
}