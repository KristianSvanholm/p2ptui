package structs

import (
	"p2p/src/constants"
)

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

