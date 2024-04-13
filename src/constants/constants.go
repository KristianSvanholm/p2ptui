package constants

import("github.com/charmbracelet/lipgloss")

const Density = 10 // Percent chance for a mine to be planted at cell.

const Size = 10 // Changing this will likely break UI. Viewport sizes are hardcoded to fit size 10

var PlayerStyle lipgloss.Style = lipgloss.NewStyle().
                                    Bold(true).
                                    Padding(0,1).
                                    Foreground(lipgloss.ANSIColor(160))

type DigEvent int 

const (
    Nothing DigEvent = iota
    Landmine
    Win
)

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
