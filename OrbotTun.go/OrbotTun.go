package IPtProxy

import (
	"github.com/eycorsican/go-tun2socks/core"
	"github.com/eycorsican/go-tun2socks/proxy/socks"
//	"github.com/eycorsican/go-tun2socks/proxy/dnsfallback"
)

var lwipStack core.LWIPStack

// PacketFlow should be implemented in Java/Kotlin.
type PacketFlow interface {
	// WritePacket should writes packets to the TUN fd.
	WritePacket(packet []byte)
}

// Write IP packets to the lwIP stack. Call this function in the main loop of
// the VpnService in Java/Kotlin, which should reads packets from the TUN fd.
func InputPacket(data []byte) {
	if lwipStack != nil {
		lwipStack.Write(data)
	}
}

// Start the lwIP stack
// 
// @param packetFlow an instance of the PacketFlow interface
//
// @param proxyHost host IP address of SOCKS proxy
// 
// @param proxyPort port for local SOCKS proxy
func StartSocks(packetFlow PacketFlow, proxyHost string, proxyPort int) {
	if packetFlow != nil && lwipStack == nil {
		lwipStack = core.NewLWIPStack()
		core.RegisterTCPConnHandler(socks.NewTCPHandler(proxyHost, uint16(proxyPort)))
		core.RegisterUDPConnHandler(NewUDPHandler())
		core.RegisterOutputFn(func(data []byte) (int, error) {
			packetFlow.WritePacket(data)
			return len(data), nil
		})
	}
}


// Stop the lwIP Stack if running
func StopSocks() {
	if lwipStack != nil {
		lwipStack.Close();
		lwipStack = nil
	}
}



