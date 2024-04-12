package structs

import (
	//"encoding/json"
	"p2p/src/constants"
	"sync"

	"github.com/gorilla/websocket"
)

type Movement struct {
	Id  string
	Pos Coords
}

type Action struct {
    Pos Coords
    Dig bool
}

type Chat struct {
	Txt string
}

type join struct {
	Id string
}

type Peer struct {
	Socket *websocket.Conn
	mutex  sync.Mutex
	Name   string
	Ip     string
}

func (p *Peer) Send(pkt interface{}) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.Socket.WriteJSON(pkt)
}

type Packet struct {
	Channel constants.WsEvent   `json:"channel"`
	Data    interface{}         `json:"data"`
}

type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (c Coords) New() *Coords {
    return &Coords{X:0, Y:0} 
}

func (c Coords) FromData(data interface{}) Coords {
    d := data.(map[string]interface{})
    c.X = int(d["x"].(float64))
    c.Y = int(d["y"].(float64))

    return c
}

func (c *Coords) Normalize() *Coords {
	if c.X < 0 {
		c.X = constants.Size - 1
	} else if c.X > constants.Size-1 {
		c.X = 0
	}

	if c.Y < 0 {
		c.Y = constants.Size - 1
	} else if c.Y > constants.Size-1 {
		c.Y = 0
	}

	return c
}

type Cell struct {
	Revealed  bool `json:"revealed"`
	Mine      bool `json:"mine"`
	Flagged   bool `json:"flagged"`
	Count     int  `json:"count"`
}

