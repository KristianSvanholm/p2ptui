package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"sync"
)

var peers = make([]*Peer, 0)
var Name string

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	//Opens up the connection for websocket
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ConnectionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	connection, err := wsUpgrader.Upgrade(w, r, nil) // Upgrade to ws
	if err != nil {
		fmt.Println("FUCK")
		http.Error(w, "Couldn't convert request to websocket", http.StatusInternalServerError)
		return
	}

	addPeer(connection)
}

func Connect(url url.URL) {
	connection, errcode, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		log.Fatal("Could not connect to network. Bye.\n", err, errcode)
	}

	addPeer(connection)
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

func leave(peer *Peer) {
	peer.Socket.Close()
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

func addPeer(c *websocket.Conn) {

	peer := Peer{
		Ip:     c.RemoteAddr().String(),
		Name:   "not_disclosed",
		Socket: c,
	}
	peers = append(peers, &peer)

	Broadcast(Name, 1)
	go listen(&peer)
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
