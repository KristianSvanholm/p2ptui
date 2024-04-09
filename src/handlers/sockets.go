package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

var peers = make(map[string]*Peer)
var Name string
var Host bool = false
var Port string

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

    currPeers := ips()

    p := addPeer(connection, r.Header.Get("Origin"))
    if p == nil {
        return
    }

    if Host {
        Tell(p, currPeers, 3)
    }

    go listen(p)
}

// Me connect to someone else
func Connect(url url.URL) {

    header :=http.Header{}
    header.Set("Origin","0.0.0.0:"+Port)
	connection, errcode, err := websocket.DefaultDialer.Dial(url.String(), header)
	if err != nil { 
		log.Fatal("Could not connect to network. Bye.\n", err, errcode)
	}

    p := addPeer(connection, connection.RemoteAddr().String())
    if p == nil {
        return
    }

    go listen(p)
}

func ips() []string {
    keys := make([]string, len(peers))

    i := 0
    for k := range peers {
        keys[i]= k
        i++
    }
    return keys
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
        case 3:
            others(pkt.Data)
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

func others(data any) {
    var ips []interface{}
    switch data.(type) {
        case []interface{}: 
            ips = data.([]interface{})
        break;
    }
    
    fmt.Println(ips)
    for _, ip := range ips {
        addr := fmt.Sprintf("%v", ip)
        url := url.URL{Scheme: "ws", Host: addr, Path: "/api/connect/"}
        Connect(url)
    }
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

func addPeer(c *websocket.Conn, address string) *Peer {
    fmt.Println(c.RemoteAddr(), c.LocalAddr())

	peer := Peer{
		Ip:     address,
		Name:   "not_disclosed",
		Socket: c,
	}

    _, found := peers[address]
    if found {
        c.Close()
        return nil
    }

    Tell(&peer, Name, 1)

    peers[address] = &peer

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
