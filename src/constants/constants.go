package constants

import (
	"github.com/charmbracelet/lipgloss"
)

const Density = 25 // Percent chance for a mine to be planted at cell.

const Size = 10 // Changing this will likely break UI. Viewport sizes are hardcoded to fit size 10

var PlayerStyle lipgloss.Style = lipgloss.NewStyle().
	Bold(true).
	Padding(0, 1).
	Foreground(lipgloss.ANSIColor(160))

// Enum for switching on results of digging up a cell
type DigEvent int

const (
	Nothing DigEvent = iota
	Win
	Landmine
)

// Enum for switching on any network request.
type WsEvent int

const (
	Chat WsEvent = iota
	Hello
	Leave
	Welcome
	Move
	Flag
	Dig
)
