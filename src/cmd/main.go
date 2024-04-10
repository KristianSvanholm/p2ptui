package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"p2p/src/handlers"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var size = 10

func main() {

    m := NewModel()

    p := tea.NewProgram(m)

    _, err := p.Run()
    if err != nil {
        log.Fatal(err)
    } 

	http.HandleFunc("/api/connect/", handlers.ConnectionHandler)

	var hostq string
	fmt.Print("Port: ")
	fmt.Scanln(&handlers.Port)
	fmt.Print("Name: ")
	fmt.Scanln(&handlers.Name)
	fmt.Print("Host? y/n: ")
	fmt.Scanln(&hostq)
	hostq = strings.ToLower(hostq)

	if hostq == "y" {
		fmt.Println("u are host")
        handlers.Host = true
	} else {
		var ntwrk string
		fmt.Print("Host port: ")
		fmt.Scanln(&ntwrk)
		url := url.URL{Scheme: "ws", Host: "0.0.0.0:" + ntwrk, Path: "/api/connect/"}
		handlers.Connect(url)
	}

    go http.ListenAndServe(":"+handlers.Port, nil)

	writeMsg()
}

func writeMsg(){
    for {
		var txt string

		fmt.Scanln(&txt)

		handlers.Broadcast(txt, 0)
	}
}

type Model struct {
    title string
    board [][]string
    players []coords
    chat []string
    gameport viewport.Model
    chatport viewport.Model
    textarea textarea.Model
}

func NewModel() Model {
    ta := textarea.New()
    ta.Placeholder = "Enter search term"
    ta.Focus()
    ta.SetHeight(10)

    players := make([]coords, 1)
    players[0] = coords{0,0}

    board := generateBoard(size)

    chat := make([]string, 0)
    
    vp := viewport.New(33, 23)
    vp.SetContent(newTable(board).Render())

    cp := viewport.New(33, 10)
    cp.SetContent("Welcome!")

    return Model{
        title: "hello world",
        players: players,
        textarea: ta,
        gameport: vp,
        chatport: cp,
        chat: chat,
    }
}

func newTable(board [][]string) *table.Table {
    return table.New().
                Border(lipgloss.NormalBorder()).
                BorderRow(true).
                BorderColumn(true).
                Rows(board...).
                StyleFunc(func(row, col int) lipgloss.Style {
                    return lipgloss.NewStyle().Padding(0,1)
                })
}

func generateBoard(size int) [][]string {
    board := make([][]string, size)
    for i:= 0; i<size; i++ {
        row := make([]string, size)

        for j := 0; j<size; j++ {
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

    var c = m.players[0]

    switch msg := msg.(type) {
        case tea.KeyMsg: 
            switch msg.Type {
                case tea.KeyEnter:
                    v := m.textarea.Value() // Get value of input
                    handlers.Broadcast(v, 0)
                    m.chat = append(m.chat, "you: " + v)
                    m.textarea.Reset()
                    m.chatport.GotoBottom()
                case tea.KeyCtrlC, tea.KeyEsc:
                    return m, tea.Quit
                case tea.KeyLeft: c.x--
                case tea.KeyRight: c.x++
                case tea.KeyUp: c.y--
                case tea.KeyDown: c.y++
            }
    }

    m.players[0] = c.normalize()

    board := generateBoard(size)
    board[c.y][c.x] = "?"

    //m.table.Rows(board...)

    m.gameport.SetContent(newTable(board).Render())

    m.chatport.SetContent(strings.Join(m.chat,"\n"))

    return m, tea.Batch(taCmd, gpCmd, cpCmd)
}

func (m Model) View() string {
    return fmt.Sprintf("%s\n%s\n%s",
        m.gameport.View(),
        m.chatport.View(),
        m.textarea.View())
}

type coords struct {
    x int
    y int
}

func (c *coords) normalize() coords {
    if c.x < 0 {
        c.x = size-1
    } else if c.x > size-1 {
        c.x = 0
    }

    if c.y <0 {
        c.y = size-1
    } else if c.y > size-1 {
        c.y = 0
    }

    return *c
}
