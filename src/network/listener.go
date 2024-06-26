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

// Listen to peers WS connection.
func listen(program *tea.Program, peer *Peer, name string, seed *string) {

	for {

		// Read packet
		pkt := Packet{}
		err := peer.Socket.ReadJSON(&pkt)
		if err != nil {
			leave(program, peer)
			return // Exit loop.
		}

		switch pkt.Channel {
		case constants.Chat: // Peer sent a chat message
			chat(program, peer, pkt.Data.(string))
		case constants.Hello: // Peer introduced themselves
			hello(program, peer, pkt.Data)
		case constants.Leave: // Peer left
			leave(program, peer)
		case constants.Welcome: // Peer welcomes you with game state and other peers.
			welcome(program, pkt.Data, name, seed)
		case constants.Move: // Peer moves their cursor on board
			move(program, peer, pkt.Data)
		case constants.Flag: // Peer plants a flag
			Act(program, peer, pkt.Data, false)
		case constants.Dig: // Peer digs a cell
			Act(program, peer, pkt.Data, true)
		}
	}
}

// Peer sends the network a chat message
func chat(program *tea.Program, peer *Peer, msg string) {
	program.Send(structs.Chat{Id: peer.Ip, Txt: msg})
}

// Peer introduces themselves to you
// TODO:: Look into combining hello(...) and welcome(...)
func hello(program *tea.Program, peer *Peer, data interface{}) {
	d := data.(map[string]interface{})

	peer.Name = d["name"].(string) // Their name

	c := structs.Coords{}.FromData(d["pos"]) // Their location

	seed, hasSeed := d["seed"] // The new seed (optional)
	if hasSeed {
		program.Send(structs.StatusUpdate{fmt.Sprintf("New seed: %s", seed.(string))})
		s, err := strconv.Atoi(seed.(string))
		if err != nil {
			// Uh-oh :)
		}
		program.Send(*rand.New(rand.NewSource(int64(s)))) // Update game RNG.
	}

	// Inform TUI
	program.Send(structs.Join{
		Id:   peer.Ip,
		Pos:  c,
		Name: peer.Name,
	})
}

// Peer has told you they are leaving.
func leave(program *tea.Program, peer *Peer) {
	program.Send(structs.StatusUpdate{fmt.Sprintf("%s left the game", peer.Name)})

	peer.Socket.Close()
	delete(Peers, peer.Ip) // Remove from network peer map

	// Notify TUI
	program.Send(structs.Leave{peer.Ip})
}

// Peer Catches you up to speed
func welcome(program *tea.Program, data interface{}, name string, seed *string) {

	d := data.(map[string]interface{})

	field := mines.Field{}
	mapstructure.Decode(d["field"], &field)

	program.Send(field) // Update game field.

	// Connect to all discovered peers
	for _, ip := range d["others"].([]interface{}) {
		addr := fmt.Sprintf("%v", ip)
		url := url.URL{Scheme: "ws", Host: addr, Path: "/api/connect/"}
		Connect(program, url, name, seed, false)
	}
}

// Peer tells you where their cursor is
func move(program *tea.Program, peer *Peer, data interface{}) {
	c := structs.Coords{}.FromData(data)

	move := structs.Movement{Id: peer.Ip, Pos: c}
	program.Send(move)
}

// Peer makes a move on the field (dig / flag)
func Act(program *tea.Program, peer *Peer, data interface{}, dig bool) {

	action := structs.Action{
		Pos: structs.Coords{}.FromData(data),
		Dig: dig,
		Id:  peer.Ip,
	}

	program.Send(action)
}
