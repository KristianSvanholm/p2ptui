package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/url"
	"p2p/src/constants"
	"p2p/src/mines"
	"p2p/src/network"
	"p2p/src/tui"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	var borders bool
	flag.BoolVar(&borders, "b", false, "Adds borders to minefield")
	flag.Parse()

	// Automatically set port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}

	network.Port = fmt.Sprint(listener.Addr().(*net.TCPAddr).Port)

	// Init TUI
	rng, seed := rng()
	field := mines.InitField(constants.Size)
	program := tea.NewProgram(tui.NewModel(field, rng, borders, &seed), tea.WithAltScreen())

	// Simple setup wizard.
	var hostq, name string
	fmt.Print("Name: ")
	fmt.Scanln(&name)
	fmt.Print("Host? y/n: ") // Host only to indicate first node in network.
	fmt.Scanln(&hostq)
	hostq = strings.ToLower(hostq)

	if hostq != "y" {
		var ntwrk string
		fmt.Print("Host address: ")
		fmt.Scanln(&ntwrk)
		url := url.URL{Scheme: "ws", Host: ntwrk, Path: "/api/connect/"}

		// Connect to node
		network.Connect(program, url, name, &seed, true)
	}

	//Serve simple webserver
	network.Serve(listener, program, name, &seed, field)

	program.Run()
}

// Sets the seed for future rng
func rng() (rand.Rand, string) {
	seed := time.Now().UTC().UnixNano()
	return *rand.New(rand.NewSource(seed)), fmt.Sprint(seed)
}
