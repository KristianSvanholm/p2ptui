package network

import (
	"p2p/src/constants"
)

type Packet struct {
	Channel constants.WsEvent   `json:"channel"`
	Data    interface{}         `json:"data"`
}

func Broadcast(t any, channel constants.WsEvent) {
	pkt := Packet{
		Channel: channel,
		Data:    t,
	}
	for _, p := range Peers {
		p.Send(pkt)
	}
}

func Tell(p *Peer, t any, channel constants.WsEvent) {
	pkt := Packet{
		Channel: channel,
		Data:    t,
	}

	p.Send(pkt)
}
