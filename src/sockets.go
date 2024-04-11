package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"p2p/src/constants"
	"p2p/src/structs"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

var Peers = make(map[string]*structs.Peer)
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
func ConnectionHandler(w http.ResponseWriter, r *http.Request, program *tea.Program) {

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
		Tell(p, currPeers, constants.Others)
	}

	go listen(program, p)
}

// Me connect to someone else
func Connect(program *tea.Program, url url.URL) {

	header := http.Header{}
	header.Set("Origin", "0.0.0.0:"+Port)
	connection, errcode, err := websocket.DefaultDialer.Dial(url.String(), header)
	if err != nil {
		log.Fatal("Could not connect to network. Bye.\n", err, errcode)
	}

	p := addPeer(connection, connection.RemoteAddr().String())
	if p == nil {
		return
	}

	go listen(program, p)
}

func ips() []string {
	keys := make([]string, len(Peers))

	i := 0
	for k := range Peers {
		keys[i] = k
		i++
	}
	return keys
}


func Broadcast(t any, channel constants.WsEvent) {
	pkt := structs.Packet{
		Channel: channel,
		Data:    t,
	}
	for _, p := range Peers {
		p.Send(pkt)
	}
}

func Tell(p *structs.Peer, t any, channel constants.WsEvent) {
	pkt := structs.Packet{
		Channel: channel,
		Data:    t,
	}

	p.Send(pkt)
}

func addPeer(c *websocket.Conn, address string) *structs.Peer {

	peer := structs.Peer{
		Ip:     address,
		Name:   "not_disclosed",
        Pos:    structs.Coords{}.New(),
		Socket: c,
	}

	_, found := Peers[address]
	if found {
		c.Close()
		return nil
	}

	Tell(&peer, Name, constants.Hello)

	Peers[address] = &peer

	return &peer
}

func listen(program *tea.Program, peer *structs.Peer) {
	for {
		pkt := structs.Packet{}
		err := peer.Socket.ReadJSON(&pkt)
		if err != nil {
			//disconnect()
			peer.Socket.Close()
			return
		}

		switch pkt.Channel {
		case constants.Chat:
			chat(program, peer, pkt.Data.(string))
			break
		case constants.Hello:
			hello(peer, pkt.Data.(string))
			break
        case constants.Leave:
			leave(peer)
		case constants.Others:
			others(program, pkt.Data)
        case constants.Move:
            move(program, peer, pkt.Data)
			break
		}
	}
}

func move(program *tea.Program, peer *structs.Peer, data interface{}) {
    var c structs.Coords
    
    d := data.(map[string]interface{})
    c.X = int(d["x"].(float64))
    c.Y = int(d["y"].(float64))

    //fmt.Println(data, c)
    move := structs.Movement{Id: peer.Ip, Pos: c}
    program.Send(move)
}

// Peer has told you they are leaving.
// Remove them from your peers and close the connection
func leave(peer *structs.Peer) {
	peer.Socket.Close()
	delete(Peers, peer.Ip)
}

func others(program *tea.Program, data any) {
    var ips []string
    mapstructure.Decode(data, ips)

	for _, addr := range ips {
		url := url.URL{Scheme: "ws", Host: addr, Path: "/api/connect/"}
		Connect(program, url)
	}
}

func hello(peer *structs.Peer, name string) {
    peer.Name = name
}

func chat(program *tea.Program, peer *structs.Peer, msg string) {
	output := fmt.Sprintf("%s - %s", peer.Name, msg)
    program.Send(structs.Chat{Txt: output})
}
