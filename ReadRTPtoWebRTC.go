package main

import (
	"errors"
	"io"
	"fmt"
	"net"
	"github.com/pion/webrtc/v3"
	
)

func ReadRTPtoWebRTC(listener *net.UDPConn, videoTrack *webrtc.TrackLocalStaticRTP) {

	// Read RTP packets forever and send them to the WebRTC Client
	inboundRTPPacket := make([]byte, 1600) // UDP MTU
	for {
		n, _, err := listener.ReadFrom(inboundRTPPacket)
		if err != nil {
			panic(fmt.Sprintf("error during read: %s", err))
		}

		if _, err = videoTrack.Write(inboundRTPPacket[:n]); err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				// The peerConnection has been closed.
				return
			}

			panic(err)
		}
	}
}