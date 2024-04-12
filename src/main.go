package main

import (
	"fmt"
	"p2p/src/constants"
	"p2p/src/structs"

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

    m := NewModel()

	program := tea.NewProgram(m)

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

	if hostq == "y" {
		fmt.Println("u are host")
		Host = true
	} else {
		var ntwrk string
		fmt.Print("Host port: ")
		fmt.Scanln(&ntwrk)
		url := url.URL{Scheme: "ws", Host: "0.0.0.0:" + ntwrk, Path: "/api/connect/"}
		Connect(program, url)
	}

	go program.Run()

	http.ListenAndServe(":"+Port, nil)

    /*
	if err != nil {
		log.Fatal(err)
	}*/
}

type Model struct {
	title    string
	board    [][]string
	player   *structs.Coords
    peers    map[string]structs.Coords
	chat     []string
	gameport viewport.Model
	chatport viewport.Model
	textarea textarea.Model
}

func NewModel() Model {
	ta := textarea.New()
	ta.Placeholder = "Enter search term"
	ta.Focus()
	ta.SetHeight(10)

	board := generateBoard()

	chat := make([]string, 0)

	vp := viewport.New(50, 23)
	vp.SetContent(newTable(board).Render())

	cp := viewport.New(50, 10)
	cp.SetContent("Welcome!")

	return Model{
		title:    "hello world",
        player:   Player,
        peers:    make(map[string]structs.Coords),
		textarea: ta,
		gameport: vp,
		chatport: cp,
		chat:     chat,
	}
}

func newTable(board [][]string) *table.Table {
	return table.New().
		Border(lipgloss.NormalBorder()).
		BorderRow(true).
		BorderColumn(true).
		Rows(board...).
		StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Padding(0, 1)
		})
}

func generateBoard() [][]string {
    size := constants.Size

	board := make([][]string, size)
	for i := 0; i < size; i++ {
		row := make([]string, size)

		for j := 0; j < size; j++ {
			row[j] = ""
		}

		board[i] = row
	}

	return board
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		gpCmd tea.Cmd
		cpCmd tea.Cmd
	)

	m.textarea, taCmd = m.textarea.Update(msg)
	m.gameport, gpCmd = m.gameport.Update(msg)

	var c = *m.player
    old_c := c

    board := generateBoard()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
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
	}

	*m.player = *c.Normalize()
    board = peerMovements(m, board)

    if old_c != c {
        Broadcast(c, constants.Move)
    }

	board[c.Y][c.X] = "!"

	m.gameport.SetContent(newTable(board).Render())

	m.chatport.SetContent(strings.Join(m.chat, "\n"))

	return m, tea.Batch(taCmd, gpCmd, cpCmd)
}

func peerMovements(m Model, board [][]string) [][]string {
    for _, c := range m.peers {
        board[c.Y][c.X] = "?"
    }
    return board
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n%s\n%s",
		m.gameport.View(),
		m.chatport.View(),
		m.textarea.View())
}
