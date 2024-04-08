package main

import (
	"fmt"
	"net/http"
	"net/url"
	"p2p/src/handlers"
	"strings"
)

func main() {

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
