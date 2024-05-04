package structs

import (
	"p2p/src/constants"
)

type Coords struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Shorthand for creating new coords at 0,0
func (c Coords) New() *Coords {
	return &Coords{X: 0, Y: 0}
}

// returns coords from an interface of coords JSON.
func (c Coords) FromData(data interface{}) Coords {
	d := data.(map[string]interface{})
	c.X = int(d["x"].(float64))
	c.Y = int(d["y"].(float64))

	return c
}

// Limits a Coord to not go out of the game board bounds.
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

// A cell in a Mine Field
type Cell struct {
	Revealed bool `json:"revealed"`
	Mine     bool `json:"mine"`
	Flagged  bool `json:"flagged"`
	Count    int  `json:"count"`
}

// ====== networked program events which are handled in model.Update(...) ======

type Join struct {
	Id   string // Who
	Pos  Coords // Where
	Name string // What
}

type StatusUpdate struct {
	Update string // What
}

type Movement struct {
	Id  string // Who
	Pos Coords // Where
}

type Action struct {
	Pos Coords // Where
	Dig bool   // What
	Id  string // Who
}

type Chat struct {
	Id  string // Who
	Txt string // What
}

type join struct {
	Id string // Who
}

type Leave struct {
	Id string // Who
}
