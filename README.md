# p2p minesweeper

Little minesweeper for the terminal with p2p multiplayer capabilities.

## TUI

Uses the [BubbleTea](https://github.com/charmbracelet/bubbletea) framework for the TUI with [lipgloss](https://github.com/charmbracelet/lipgloss) for the colors & styling.

Player cursor position is indicated by a color styling on the specific cell location. If no contents in cell, player cursor will be invisible.  
Could potentially be improved by coloring the actual grid around the specific cell, which would allow for colored flags and numbers as well.

![image](https://github.com/KristianSvanholm/p2ptui/assets/61845965/23710705-9d40-46b9-9bf8-7686ad3c0827)

## Networking

- All peers are connected to all peers.
- Everyone runs the same game off the same seed individually and informs eachother of their actions within the game to keep the session in sync.
- The game state is only sent once when joining a session.
- Requests are sent using the [Gorilla Websockets](https://github.com/gorilla/websocket) package and sends as little information as possible on the line.
- Intense packet spam will likely break sync between peers, but I haven't been able to break sync myself.
- You can connect to any peer to join the session, there is no leader node.
  - Which also means no authority. Very easy to cheat :D

## Run

`go run src/main.go`

Add `-b` if you want borders on your grid.
Select `host` if you are first in network. This does not give your node any authority in the network.

## Controls

- `SHIFT+UP` - Plant flag
- `SHIFT+DOWN` - Dig cell
- `Arrow keys` - Move cursor around
- `ENTER` - Send chat message
- `CTRL+C` / `ESC` - Exit game

## Future ideas

- Adding a better setup UI for selecting port, name etc...
  - Potentially autoselecting port and then displaying that in UI for player to share with peers.
- Adding support for going across the internet. As of today it assumes you are playing with someone on your LAN.
  - Would require port forwarding by each peer. Potentially use network hole punching.
- Improved custom serialization of packets
- Voting system for:
  - Kicking peers,
  - Restarting games,
  - Adjusting difficulty.
- Tick system for synchronization purposes.
