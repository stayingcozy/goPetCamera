package main 

import (
	"github.com/pion/webrtc/v3"
)

func SetRemDes(peerConnection *webrtc.PeerConnection, offer webrtc.SessionDescription, resChan chan error) {
	resChan <- peerConnection.SetRemoteDescription(offer)
}