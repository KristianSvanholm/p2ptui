package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"p2p/src/handlers"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

    m := NewModel()

    p := tea.NewProgram(m)

    _, err := p.Run()
    if err != nil {
        log.Fatal(err)
    }


	http.HandleFunc("/api/connect/", handlers.ConnectionHandler)

	var port, hostq string
	fmt.Print("Port: ")
	fmt.Scanln(&port)
	fmt.Print("Name: ")
	fmt.Scanln(&handlers.Name)
	fmt.Print("Host? y/n: ")
	fmt.Scanln(&hostq)
	hostq = strings.ToLower(hostq)

	if hostq == "y" {
		fmt.Println("u are host")
	} else {
		var ntwrk string
		fmt.Print("Host port: ")
		fmt.Scanln(&ntwrk)
		url := url.URL{Scheme: "ws", Host: "0.0.0.0:" + ntwrk, Path: "/api/connect/"}
		handlers.Connect(url)
	}

	go http.ListenAndServe(":"+port, nil)

	writeMsg()
}

func writeMsg() {
	for {
		var txt string

		fmt.Scanln(&txt)

		handlers.Broadcast(txt, 0)
	}
}

type Model struct {
    title string

}

func NewModel() Model {
    return Model{
        title: "hello world",
    }
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return m, nil
}

func (m Model) View() string {
    return m.title
}


