package constants

import("github.com/charmbracelet/lipgloss")

const Size = 10

var PlayerStyle lipgloss.Style = lipgloss.NewStyle().
                                    Bold(true).
                                    Padding(0,1).
                                    Foreground(lipgloss.ANSIColor(160))

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
