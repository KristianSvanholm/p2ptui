package tui

import (
	"fmt"
	"math/rand"
	"p2p/src/constants"
	"p2p/src/mines"
	"p2p/src/network"
	"p2p/src/structs"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// The central data structure for UI, controls and game state
type Model struct {
	status   string
	player   *structs.Coords
	peers    map[string]*UIPeer
	chat     []string
	Seed     *string
	Rng      *rand.Rand
	border   lipgloss.Border
	Field    *mines.Field
	gameport viewport.Model
	peerport viewport.Model
	chatport viewport.Model
	textarea textarea.Model
}

// Initiazes a model
func NewModel(field *mines.Field, rng rand.Rand, borders bool, seed *string) *Model {

	// Borders or not in table
	b := lipgloss.HiddenBorder()
	if borders {
		b = lipgloss.NormalBorder()
	}

	ta := textarea.New()
	ta.Placeholder = "Write your message..."
	ta.Focus()
	ta.SetHeight(2)

	chat := make([]string, 0)

	vp := viewport.New(41, 23)
	pp := viewport.New(40, 23)
	cp := viewport.New(41, 10)
	cp.SetContent("Welcome!")

	return &Model{
		status:   "Hello world",
		player:   network.Player,
		peers:    make(map[string]*UIPeer),
		chat:     chat,
		Seed:     seed,
		border:   b,
		Rng:      &rng,
		Field:    field,
		textarea: ta,
		gameport: vp,
		peerport: pp,
		chatport: cp,
	}
}

func (m *Model) Init() tea.Cmd {
	return textarea.Blink
}

// Main game loop.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		gpCmd tea.Cmd
		cpCmd tea.Cmd
		ppCmd tea.Cmd
	)

	// Run update function on all TUI components.
	m.textarea, taCmd = m.textarea.Update(msg)
	m.gameport, gpCmd = m.gameport.Update(msg)
	m.chatport, cpCmd = m.chatport.Update(msg)
	m.peerport, ppCmd = m.peerport.Update(msg)

	// Previous player position
	var c = *m.player
	old_c := c

	// Handle any incoming messages.
	switch msg := msg.(type) {
	case structs.Movement: // Peer movement
		m.peers[msg.Id].move(msg.Pos)
	case structs.Join: // Peer joining
		m.peers[msg.Id] = createUIPeer(msg)
	case structs.Leave: // Peer leaving
		delete(m.peers, msg.Id)
	case structs.Chat: // Peer Chatting
		p := m.peers[msg.Id]
		m.chat = append(m.chat, chatter(p.name, msg.Txt, p.style))
		m.chatport.GotoBottom()
	case structs.Action: // Peer plays
		if msg.Dig {
			handleDigEvent(m, m.Field.Dig(&msg.Pos, m.Rng), m.peers[msg.Id].name)
		} else {
			if m.Field.SetFlag(msg.Pos) {
				m.status = "Win!"
			}
		}
	case mines.Field: // Peer shares current MineField
		*m.Field = msg
	case structs.StatusUpdate: // Set status field text.
		m.status = msg.Update
	case rand.Rand: // New peer updates seed.
		m.Rng = &msg
	case tea.KeyMsg: // Local user key inputs

		switch msg.Type {

		// Move cursor
		case tea.KeyLeft:
			c.X--
		case tea.KeyRight:
			c.X++
		case tea.KeyUp:
			c.Y--
		case tea.KeyDown:
			c.Y++
		case tea.KeyShiftUp: // Plant flag "Flagpole is up"
			if m.Field.SetFlag(*m.player) {
				m.status = "Win!"
			}
			network.Broadcast(*m.player, constants.Flag)

		case tea.KeyShiftDown: // Dig ground "Shovel is down"
			handleDigEvent(m, m.Field.Dig(m.player, m.Rng), "You")
			network.Broadcast(*m.player, constants.Dig)

		case tea.KeyEnter: // Send message
			v := m.textarea.Value() // Get value of input
			if len(v) != 0 {
				network.Broadcast(v, constants.Chat)
				m.chat = append(m.chat, chatter("You", v, constants.PlayerStyle))
				m.textarea.Reset()
				m.chatport.GotoBottom()
			}

		case tea.KeyCtrlC, tea.KeyEsc: // Exit program
			return m, tea.Quit
		}
	}

	// Prevent out of field coordinates
	*m.player = *c.Normalize()

	board := generateBoard(m.Field) // Generate new board based on current game state

	if old_c != c { // Broadcast move if not same position
		network.Broadcast(c, constants.Move)
	}

	// Update TUI contents
	m.gameport.SetContent(newTable(board, m).Render())
	m.chatport.SetContent(strings.Join(m.chat, "")) // Seemingly no newline is required here as the messages come with an attached \n.
	m.peerport.SetContent(peerList(m.peers))

	return m, tea.Batch(taCmd, gpCmd, cpCmd, ppCmd)
}

// Formats the actual TUI through string operations.
func (m *Model) View() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s",
		m.status,
		lipgloss.JoinHorizontal(lipgloss.Top, m.gameport.View(), m.peerport.View()),
		m.chatport.View(),
		m.textarea.View())
}
