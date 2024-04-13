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

type Model struct {
	status    string
	player   *structs.Coords
    peers    map[string]*UIPeer
    chat     []string
    Seed     *string
    Rng      *rand.Rand
    Field    *mines.Field
	gameport viewport.Model
    peerport viewport.Model
	chatport viewport.Model
	textarea textarea.Model
}

func NewModel(field *mines.Field, rng rand.Rand, seed *string) *Model {
	ta := textarea.New()
	ta.Placeholder = "Enter search term"
	ta.Focus()
	ta.SetHeight(2)

	chat := make([]string, 0)

	vp := viewport.New(41, 23)
    pp := viewport.New(10,23)
	cp := viewport.New(41, 10)
	cp.SetContent("Welcome!")

	return &Model{
		status:    "Hello world",
        player:   network.Player,
        peers:    make(map[string]*UIPeer),
		chat:     chat,
        Seed:     seed,
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

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		gpCmd tea.Cmd
		cpCmd tea.Cmd
        ppCmd tea.Cmd
	)

	m.textarea, taCmd = m.textarea.Update(msg)
	m.gameport, gpCmd = m.gameport.Update(msg)
    m.chatport, cpCmd = m.chatport.Update(msg)
    m.peerport, ppCmd = m.peerport.Update(msg)

    // Previous player position
	var c = *m.player
    old_c := c

	switch msg := msg.(type) {
        case structs.Movement: m.peers[msg.Id].move(msg.Pos) // Peer movement
        case structs.Join: m.peers[msg.Id] = createUIPeer(msg) // Peer joining
        case structs.Leave: delete(m.peers, msg.Id) // Peer leaving
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
        case mines.Field: *m.Field = msg // Peer shares current MineField
        case structs.StatusUpdate: m.status = msg.Update // Headline is update (sys info)
        case rand.Rand: m.Rng = &msg // New peer updates seed
    	case tea.KeyMsg: // Local user key inputs

        	switch msg.Type {

                // Move cursor
        		case tea.KeyLeft:   c.X--
        		case tea.KeyRight:  c.X++
        		case tea.KeyUp:     c.Y--
        		case tea.KeyDown:   c.Y++
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
                    // Todo:: Inform the peers
        			return m, tea.Quit

            }
    }

    // Prevent out of field coordinates
	*m.player = *c.Normalize()

    board := generateBoard(m.Field)

    if old_c != c { // Broadcast move if not same position
        network.Broadcast(c, constants.Move)
    }

    // Update TUI contents
	m.gameport.SetContent(newTable(board, m).Render())
    m.chatport.SetContent(strings.Join(m.chat,"")) // No clue why a newline is not required here :P
    m.peerport.SetContent(peerList(m.peers))

	return m, tea.Batch(taCmd, gpCmd, cpCmd, ppCmd)
}

func (m *Model) View() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s",
        m.status,
        lipgloss.JoinHorizontal(lipgloss.Top, m.gameport.View(), m.peerport.View()),
		m.chatport.View(),
		m.textarea.View())
}
