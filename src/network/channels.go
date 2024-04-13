package network

import (
	"p2p/src/structs"
)

type Join struct {
    Id string
    Pos structs.Coords
    Name string
}

type StatusUpdate struct {
    Update string
}

type Movement struct {
	Id  string
	Pos structs.Coords
}

type Action struct {
    Pos structs.Coords
    Dig bool
}

type Chat struct {
	Txt string
}

type join struct {
	Id string
}

type Leave struct {
    Id string
}
