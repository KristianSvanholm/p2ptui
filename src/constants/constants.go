package constants

const Size = 10

type WsEvent int

const (
    Chat WsEvent = iota
    Hello
    Leave
    Others
    Move
)
