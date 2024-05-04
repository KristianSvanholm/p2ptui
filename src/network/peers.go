package network

import (
	"p2p/src/constants"
	"sync"

	"github.com/gorilla/websocket"
)

// List of peers in network
var Peers = make(map[string]*Peer)

type Peer struct {
	Socket *websocket.Conn
	mutex  sync.Mutex
	Name   string
	Ip     string
}

// Sends a packet to a peer.
func (p *Peer) Send(pkt interface{}) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.Socket.WriteJSON(pkt)
}

// Adds a new peer from a websocket connection to the peer list.
func addPeer(c *websocket.Conn, address string, name string, seed *string) *Peer {

	peer := Peer{
		Ip:     address,
		Socket: c,
	}

	// Already in list, close connection.
	_, found := Peers[address]
	if found {
		c.Close()
		return nil
	}

	// Information to share with new peer
	data := map[string]any{
		"name": name,
		"pos":  *Player,
	}

	// Optionally inform them of the new seed if you are new to network.
	if seed != nil {
		data["seed"] = *seed
	}

	Tell(&peer, data, constants.Hello)

	Peers[address] = &peer

	return &peer
}
