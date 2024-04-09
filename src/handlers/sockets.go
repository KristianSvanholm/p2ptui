package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"sync"
)

var peers = make(map[string]*Peer)
var Name string

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	//Opens up the connection for websocket
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Someone else connect to me
func ConnectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	connection, err := wsUpgrader.Upgrade(w, r, nil) // Upgrade to ws
	if err != nil {
		http.Error(w, "Couldn't convert request to websocket", http.StatusInternalServerError)
		return
	}

    p := addPeer(connection)
    if p == nil {
        return
    }

    //if host Share all of your peers with them?
    go listen(p)
}

// Me connect to someone else
func Connect(url url.URL) {
	connection, errcode, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatal("Could not connect to network. Bye.\n", err, errcode)
	}

    p := addPeer(connection)
    if p == nil {
        return
    }

    go listen(p)
}

func disconnect() {
	fmt.Println("disconnect")

}

func listen(peer *Peer) {
	for {
		pkt := Packet{}
		err := peer.Socket.ReadJSON(&pkt)
		if err != nil {
			//disconnect()
            peer.Socket.Close()
			return
			//continue //todo :: figure out what to do here
		}

		switch pkt.Channel {
		case 0:
			msg(peer, pkt.Data.(string))
			break
		case 1:
			name(peer, pkt.Data.(string))
			break
		case 2:
			leave(peer)
			break
		}
	}

}

// Peer has told you they are leaving.
// Remove them from your peers and close the connection
func leave(peer *Peer) {
	peer.Socket.Close()
    delete(peers,peer.Ip)
}

func name(peer *Peer, name string) {
	peer.Name = name
}

func msg(peer *Peer, msg string) {
	output := fmt.Sprintf("%s - %s", peer.Name, msg)
	fmt.Println(output)
}

func Broadcast(t any, channel int) {
	pkt := Packet{
		Channel: channel,
		Data:    t,
	}
	for _, p := range peers {
		p.Send(pkt)
	}
}

func Tell(p *Peer, t any, channel int) {
    pkt := Packet{
        Channel: channel,
        Data: t,
    }

    p.Send(pkt)
}

func addPeer(c *websocket.Conn) *Peer {

    addr := c.RemoteAddr().String()

	peer := Peer{
		Ip:     addr,
		Name:   "not_disclosed",
		Socket: c,
	}

    _, found := peers[addr]
    if found {
        c.Close()
        return nil
    }

    Tell(&peer, Name, 1)

    peers[addr] = &peer

    return &peer
}

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

type Packet struct {
	Channel int         `json:"channel"`
	Data    interface{} `json:"data"`
}

func handle(pkt *Packet) {
	fmt.Println(pkt.Data)
}
