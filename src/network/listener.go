package network

import (
	"fmt"
	"math/rand"
	"net/url"
	"p2p/src/constants"
	"p2p/src/mines"
	"p2p/src/structs"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mitchellh/mapstructure"
)

func listen(program *tea.Program, peer *Peer, name string, seed *string) {
	for {
		pkt := Packet{}
		err := peer.Socket.ReadJSON(&pkt)
		if err != nil {
            program.Send(Leave{peer.Ip})
			peer.Socket.Close()
			return
		}

		switch pkt.Channel {
		case constants.Chat:
			chat(program, peer, pkt.Data.(string))
		case constants.Hello:
			hello(program, peer, pkt.Data)
        case constants.Leave:
			leave(program, peer)
		case constants.Welcome:
			welcome(program, pkt.Data, name, seed)
        case constants.Move:
            move(program, peer, pkt.Data)
        case constants.Flag:
            Act(program, pkt.Data, false)
        case constants.Dig:
            Act(program, pkt.Data, true)
		}
	}
}

// Peer sends the network a chat message
func chat(program *tea.Program, peer *Peer, msg string) {
	output := fmt.Sprintf("%s - %s", peer.Name, msg)
    program.Send(Chat{Txt: output})
}

// Peer introduces themselves to you 
//TODO:: Look into combining hello(...) and welcome(...)
func hello(program *tea.Program, peer *Peer, data interface{}) {
    d := data.(map[string]interface{})

    peer.Name = d["name"].(string)

    c := structs.Coords{}.FromData(d["pos"])

    seed, hasSeed := d["seed"]
    if hasSeed {
        program.Send(StatusUpdate{seed.(string)})
        s, err := strconv.Atoi(seed.(string))
        if err != nil {
            // Uh-oh
        }
        program.Send(*rand.New(rand.NewSource(int64(s))))
    }

    program.Send(Join{
                    Id: peer.Ip, 
                    Pos: c,
                    Name: peer.Name,
                })
}

// Peer has told you they are leaving.
func leave(program *tea.Program, peer *Peer) {
	peer.Socket.Close()
	delete(Peers, peer.Ip) // Remove from network peer map

    // Notify to UI
    program.Send(Leave{peer.Ip})
}

// Peer Catches you up to speed
func welcome(program *tea.Program, data interface{}, name string, seed *string) {

    d := data.(map[string]interface{})

    field := mines.Field{}
    mapstructure.Decode(d["field"], &field)

    program.Send(field)

    // Connect to all discovered peers
	for _, ip := range d["others"].([]interface{}) {
        addr := fmt.Sprintf("%v", ip)
		url := url.URL{Scheme: "ws", Host: addr, Path: "/api/connect/"}
		Connect(program, url, name, seed)
	}
}

// Peer tells you where their cursor is
func move(program *tea.Program, peer *Peer, data interface{}) {
    c := structs.Coords{}.FromData(data)

    move := Movement{Id: peer.Ip, Pos: c}
    program.Send(move)
}

// Peer makes a move on the field (dig / flag)
func Act(program *tea.Program, data interface{}, dig bool){

    action := Action{
        Pos: structs.Coords{}.FromData(data),
        Dig: dig,
    }

    program.Send(action) 
}
