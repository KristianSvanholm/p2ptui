package network

import (
	"net"
	"net/http"
	"p2p/src/mines"

	tea "github.com/charmbracelet/bubbletea"
)

// Serves simple webserver to facilitate connection for new peers.
func Serve(listener net.Listener, program *tea.Program, name string, seed *string, field *mines.Field) {
	http.HandleFunc("/api/connect/", func(w http.ResponseWriter, r *http.Request) {
		ConnectionHandler(w, r, program, name, seed, field)
	})

	go http.Serve(listener, nil)
}
