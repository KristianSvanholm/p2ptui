# p2p minesweeper
Little minesweeper for the terminal with p2p multiplayer capabilities.
## TUI
Uses the [BubbleTea](https://github.com/charmbracelet/bubbletea) framework for the TUI with [lipgloss](https://github.com/charmbracelet/lipgloss) for the colors & styling.

![image](https://github.com/KristianSvanholm/p2ptui/assets/61845965/27a61a4c-2b47-4bac-aca9-5e9d7707949d)

## Networking
* All peers are connected to all peers.  
* Everyone runs the same game off the same seed individually and informs eachother of their actions within the game to keep the session in sync.
* The game state is only sent once when joining a session.
* Requests are sent using the [Gorilla Websockets](https://github.com/gorilla/websocket) package and sends as little information as possible on the line.
* Packet spam will likely break sync between peers, but I haven't been able to break sync myself.
* You can connect to any peer to join the session, there is no leader node.
  * Which also means no authority. Very easy to cheat :D

## Run
`go run src/main.go`

## Controls
* SHIFT+UP - Plant flag
* SHIFT+DOWN - Dig cell
* Arrow keys - Move cursor around
* ENTER - Send chat message
* CTRL+C / ESC - Exit game
