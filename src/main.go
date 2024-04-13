package main

import (
	"fmt"
	"math/rand"
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

    rng, seed := rng()
    fmt.Println(seed)
    field := mines.InitField(constants.Size)
	program := tea.NewProgram(tui.NewModel(field, rng, &seed))

    // Config thingy
	var hostq, name string
	fmt.Print("Port: ")
	fmt.Scanln(&network.Port)
	fmt.Print("Name: ")
	fmt.Scanln(&name)
	fmt.Print("Host? y/n: ")
	fmt.Scanln(&hostq)
	hostq = strings.ToLower(hostq)

	if hostq != "y" {
		var ntwrk string
		fmt.Print("Host port: ")
		fmt.Scanln(&ntwrk)
		url := url.URL{Scheme: "ws", Host: "0.0.0.0:" + ntwrk, Path: "/api/connect/"}
       
		network.Connect(program, url, name, &seed)
	}

    network.Serve(program, name, &seed, field)

   	program.Run()

}

// Sets the seed for future rng
func rng() (rand.Rand, string) {
    seed := time.Now().UTC().UnixNano()
    return *rand.New(rand.NewSource(seed)), fmt.Sprint(seed) 
}
