package network

import (
	"p2p/src/constants"
)

// The packet type informs receiver how to -
// handle Data field through thee Channel field
type Packet struct {
	Channel constants.WsEvent `json:"channel"`
	Data    interface{}       `json:"data"`
}

// Tell everyone
func Broadcast(t any, channel constants.WsEvent) {
	pkt := Packet{
		Channel: channel,
		Data:    t,
	}
	for _, p := range Peers {
		p.Send(pkt)
	}
}

// Tell someone
func Tell(p *Peer, t any, channel constants.WsEvent) {
	pkt := Packet{
		Channel: channel,
		Data:    t,
	}

	p.Send(pkt)
}
