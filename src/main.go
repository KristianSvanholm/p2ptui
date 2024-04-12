package main

import (
	"fmt"
	"math/rand"
	"p2p/src/constants"
	"p2p/src/mines"
	"p2p/src/structs"
	"strconv"
	"time"

	//"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func main() {

    M = NewModel()

	program := tea.NewProgram(M)

	http.HandleFunc("/api/connect/", func (w http.ResponseWriter, r *http.Request) {
        ConnectionHandler(w, r, program)
    })

	var hostq string
	fmt.Print("Port: ")
	fmt.Scanln(&Port)
	fmt.Print("Name: ")
	fmt.Scanln(&Name)
	fmt.Print("Host? y/n: ")
	fmt.Scanln(&hostq)
	hostq = strings.ToLower(hostq)

	if hostq != "y" {
		var ntwrk string
		fmt.Print("Host port: ")
		fmt.Scanln(&ntwrk)
		url := url.URL{Scheme: "ws", Host: "0.0.0.0:" + ntwrk, Path: "/api/connect/"}
       
        // Set new seed for whole network
        seed := time.Now().UTC().UnixNano()
        *M.Rng = *rand.New(rand.NewSource(seed))
        M.Seed = fmt.Sprint(seed)
       
		Connect(program, url, true)
	}

	go program.Run()

	http.ListenAndServe(":"+Port, nil)

    /*
	if err != nil {
		log.Fatal(err)
	}*/
}

type Model struct {
	status    string
	player   *structs.Coords
    peers    map[string]structs.Coords
	chat     []string
    Seed     string
    Rng      *rand.Rand
    Field    mines.Field
	gameport viewport.Model
	chatport viewport.Model
	textarea textarea.Model
}

func NewModel() *Model {
	ta := textarea.New()
	ta.Placeholder = "Enter search term"
	ta.Focus()
	ta.SetHeight(10)

	chat := make([]string, 0)

	vp := viewport.New(50, 23)

	cp := viewport.New(50, 10)
	cp.SetContent("Welcome!")

	return &Model{
		status:    "hello world",
        player:   Player,
        peers:    make(map[string]structs.Coords),
        Rng:      rand.New(rand.NewSource(time.Now().UTC().UnixNano())),
        Field:    *mines.InitField(constants.Size),
		textarea: ta,
		gameport: vp,
		chatport: cp,
		chat:     chat,
	}
}

func newTable(board [][]string, m *Model) *table.Table {
	return table.New().
		Border(lipgloss.NormalBorder()).
		BorderRow(true).
		BorderColumn(true).
		Rows(board...).
        StyleFunc(func(row, col int) lipgloss.Style {

            style := lipgloss.NewStyle().Padding(0,1)/*.Background(lipgloss.ANSIColor(15)) */

            c := structs.Coords{X: col, Y: row-1}
            
            // Player cursor
            if *m.player == c {
                return style.Bold(true).Foreground(lipgloss.ANSIColor(160))
            }

            // Peer cursors
            for _, p := range m.peers {
                if p == c {
                    return style.Bold(true).Foreground(lipgloss.ANSIColor(199))
                }
            }
            
            return style/*.Foreground(lipgloss.ANSIColor(8))*/
        })
}

func generateBoard(f *mines.Field) [][]string {
    size := constants.Size

	board := make([][]string, size)
	for i := 0; i < size; i++ {
		row := make([]string, size)

		for j := 0; j < size; j++ {
            cell := f.Field[j][i]
            if cell.Revealed {
                if cell.Count != 0 {
                    row[j] = strconv.Itoa(cell.Count)
                } else {
                    row[j] = " " // NOTE: NOT A SPACE - SPECIAL SYMBOL
                }
            } else {
                if cell.Flagged {
                    row[j] = "▲"
                } else {
                    row[j] = "■"
                }
            }
		}

		board[i] = row
	}

	return board
}

func (m *Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		gpCmd tea.Cmd
		cpCmd tea.Cmd
	)

	m.textarea, taCmd = m.textarea.Update(msg)
	m.gameport, gpCmd = m.gameport.Update(msg)

	var c = *m.player
    old_c := c

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
            case tea.KeyShiftUp:
                m.Field.SetFlag(*m.player)
                Broadcast(*m.player, constants.Flag)
            case tea.KeyShiftDown:
                m.Field.Dig(m.player, m.Rng)
                Broadcast(*m.player, constants.Dig)
		    case tea.KeyEnter:
    			v := m.textarea.Value() // Get value of input
    			Broadcast(v, constants.Chat)
    			m.chat = append(m.chat, "you: "+v)
    			m.textarea.Reset()
    			m.chatport.GotoBottom()
    		case tea.KeyCtrlC, tea.KeyEsc:
    			return m, tea.Quit
    		case tea.KeyLeft:
    			c.X--
    		case tea.KeyRight:
    			c.X++
    		case tea.KeyUp:
    			c.Y--
    		case tea.KeyDown:
    			c.Y++
		}
    case structs.Movement: m.peers[msg.Id] = msg.Pos
    case structs.Chat: 
        m.chat = append(m.chat, msg.Txt)
        m.chatport.GotoBottom()
    case structs.Action:
        if msg.Dig {
            m.Field.Dig(&msg.Pos, m.Rng)
        } else {
            m.Field.SetFlag(msg.Pos)
        }	
    case mines.Field: m.Field = msg
    case structs.StatusUpdate: m.status = msg.Update
    }

	*m.player = *c.Normalize()

    board := generateBoard(&m.Field)

    if old_c != c { // Broadcast move if not same position
        Broadcast(c, constants.Move)
    }

    // Update tui contents
	m.gameport.SetContent(newTable(board, m).Render())
	m.chatport.SetContent(strings.Join(m.chat, "\n"))

	return m, tea.Batch(taCmd, gpCmd, cpCmd)
}

func (m *Model) View() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s",
        m.status,
		m.gameport.View(),
		m.chatport.View(),
		m.textarea.View())
}
