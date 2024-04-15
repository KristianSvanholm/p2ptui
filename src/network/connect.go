package network

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"p2p/src/constants"
	"p2p/src/mines"
	"p2p/src/structs"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)


var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

// Someone else connect to me
func ConnectionHandler(w http.ResponseWriter, r *http.Request, program *tea.Program, name string, seed *string, field *mines.Field) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	connection, err := wsUpgrader.Upgrade(w, r, nil) // Upgrade to ws

	if err != nil {
		http.Error(w, "Couldn't convert request to websocket", http.StatusInternalServerError)
		return
	}

    ips := ips() // Get list before adding new peer (!)

    ip := strings.Split(r.RemoteAddr, ":")[0]

    p := addPeer(connection, ip+":"+r.Header.Get("Origin"), name, nil)
	if p == nil {
		return
	}

    sendWelcome(p, ips, field)
        
    go program.Send(structs.StatusUpdate{"New peer joined the game!"})
	go listen(program, p, name, seed)
}


// Me connect to someone else
func Connect(program *tea.Program, url url.URL, name string, seed *string) {

	header := http.Header{}
    header.Set("Origin", Port)
	connection, errcode, err := websocket.DefaultDialer.Dial(url.String(), header)
	if err != nil {
		log.Fatal("Could not connect to network. Bye.\n", err, errcode)
	}

    go program.Send(structs.StatusUpdate{fmt.Sprintf("New seed: %s",*seed)})
	p := addPeer(connection, connection.RemoteAddr().String(), name, seed)
	if p == nil {
		return
	}

	go listen(program, p, name, seed)
}


// Gather list of peers' ips
func ips() []string {
	keys := make([]string, len(Peers))

	i := 0
	for k := range Peers {
		keys[i] = k
		i++
	}
	return keys
}


func sendWelcome(p *Peer, ips []string, field *mines.Field) {
    data := map[string]any{
        "others": ips,
        "field": *field,
    }
    
    Tell(p, data, constants.Welcome)
}
