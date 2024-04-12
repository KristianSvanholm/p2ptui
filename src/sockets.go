package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"p2p/src/constants"
	"p2p/src/mines"
	"p2p/src/structs"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

var Peers = make(map[string]*structs.Peer)
var Name string
var Port string
var Player *structs.Coords = structs.Coords{}.New()
var M *Model

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

    ips := ips() // Get list before adding new peer (!)

	p := addPeer(connection, r.Header.Get("Origin"), false)
	if p == nil {
		return
	}

    sendWelcome(p, ips, r.Header.Get("Join") == "true")

	go listen(program, p)
}

func sendWelcome(p *structs.Peer, ips []string, join bool) {
    data := map[string]any{
        "others": ips,
        "field": M.Field, // This is probably bad
    }

    Tell(p, data, constants.Welcome)
}

// Me connect to someone else
func Connect(program *tea.Program, url url.URL, join bool) {

	header := http.Header{}
	header.Set("Origin", "0.0.0.0:"+Port)
    header.Set("Join", strconv.FormatBool(join))
	connection, errcode, err := websocket.DefaultDialer.Dial(url.String(), header)
	if err != nil {
		log.Fatal("Could not connect to network. Bye.\n", err, errcode)
	}

	p := addPeer(connection, connection.RemoteAddr().String(), true)
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

func addPeer(c *websocket.Conn, address string, seeded bool) *structs.Peer {

	peer := structs.Peer{
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
        "name": Name,
        "pos": *Player,
    }

    // Let new peer decide seed
    if seeded {
        data["seed"] = M.Seed
    }

	Tell(&peer, data, constants.Hello)

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
		case constants.Hello:
			hello(program, peer, pkt.Data)
        case constants.Leave:
			leave(peer)
		case constants.Welcome:
			welcome(program, pkt.Data)
        case constants.Move:
            move(program, peer, pkt.Data)
        case constants.Flag:
            Action(program, pkt.Data, false)
        case constants.Dig:
            Action(program, pkt.Data, true)
		}
	}
}

func Action(program *tea.Program, data interface{}, dig bool){

    action := structs.Action{
        Pos: structs.Coords{}.FromData(data),
        Dig: dig,
    }

    program.Send(action) 
}

func move(program *tea.Program, peer *structs.Peer, data interface{}) {
    c := structs.Coords{}.FromData(data)

    move := structs.Movement{Id: peer.Ip, Pos: c}
    program.Send(move)
}

// Peer has told you they are leaving.
// Remove them from your peers and close the connection
func leave(peer *structs.Peer) {
	peer.Socket.Close()
	delete(Peers, peer.Ip)
}

func welcome(program *tea.Program, data interface{}) {
    d := data.(map[string]interface{})

    ips := d["others"]
    f := d["field"]

    field := mines.Field{}
    mapstructure.Decode(f, &field)

    program.Send(field)

	for _, ip := range ips.([]interface{}) {
        addr := fmt.Sprintf("%v", ip)
		url := url.URL{Scheme: "ws", Host: addr, Path: "/api/connect/"}
		Connect(program, url, false)
	}
}

func hello(program *tea.Program, peer *structs.Peer, data interface{}) {
    d := data.(map[string]interface{})

    peer.Name = d["name"].(string)

    c := structs.Coords{}.FromData(d["pos"])

    seed, hasSeed := d["seed"]
    if hasSeed {
        s, err := strconv.Atoi(seed.(string))
        if err != nil {
            // Uh-oh
        }
        *M.Rng = *rand.New(rand.NewSource(int64(s)))
    }

    program.Send(structs.Movement{
                        Id: peer.Ip, 
                        Pos: c,
                })
}

func chat(program *tea.Program, peer *structs.Peer, msg string) {
	output := fmt.Sprintf("%s - %s", peer.Name, msg)
    program.Send(structs.Chat{Txt: output})
}
