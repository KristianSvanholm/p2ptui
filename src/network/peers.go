package network

import (
	"p2p/src/constants"
	"sync"

	"github.com/gorilla/websocket"
)

var Peers = make(map[string]*Peer)

type Peer struct {
	Socket *websocket.Conn
	mutex  sync.Mutex
	Name   string
	Ip     string
}

func (p *Peer) Send(pkt interface{}) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.Socket.WriteJSON(pkt)
}

func addPeer(c *websocket.Conn, address string, name string, seed *string) *Peer {

	peer := Peer{
		Ip:     address,
		Name:   "not_disclosed",
		Socket: c,
	}

	_, found := Peers[address]
	if found {
		c.Close()
		return nil
	}

    data := map[string]any{
        "name": name,
        "pos": *Player,
    }

    if seed != nil {
        data["seed"] = *seed
    }

	Tell(&peer, data, constants.Hello)

	Peers[address] = &peer

	return &peer
}
