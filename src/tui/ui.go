package tui

import (
	"p2p/src/constants"
	"p2p/src/mines"
	"p2p/src/network"
	"p2p/src/structs"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type UIPeer struct {
    style lipgloss.Style
    name  string
    pos   structs.Coords
}

func (p *UIPeer) move(pos structs.Coords){
    p.pos = pos
}

func newTable(board [][]string, m *Model) *table.Table {
	return table.New().
		Border(lipgloss.NormalBorder()).
		BorderRow(true).
		BorderColumn(true).
		Rows(board...).
        StyleFunc(func(row, col int) lipgloss.Style {

            c := structs.Coords{X: col, Y: row-1}
            
            // Player cursor
            if *m.player == c {
                return constants.PlayerStyle
            }

            // Peer cursors
            for _, p := range m.peers {
                if p.pos == c {
                    return p.style.Padding(0,1)
                }
            }
            
            return lipgloss.NewStyle().Padding(0,1)
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
                    row[j] = "⚑"
                } else {
                    row[j] = "■"
                }
            }
		}

		board[i] = row
	}

	return board
}

func createUIPeer(peer network.Join) *UIPeer{
    return &UIPeer{
        pos: peer.Pos,
        style: lipgloss.NewStyle().
                        Bold(true).
                        Padding(0,1).
                        Foreground(lipgloss.ANSIColor(199)),
        name: peer.Name,
    }
}

func peerList(peers map[string]*UIPeer) string {

    arr := make([]*UIPeer, len(peers))
    i := 0

    // map don't always evaluate to the same order, 
    // so we need an array ...
    for _, peer := range peers {
        arr[i] = peer
        i++
    }

    // ... Which we can then sort ...
    sort.Slice(arr, func(i2, j int) bool {
        return arr[i2].name > arr[j].name
    })

    var sb strings.Builder
    sb.WriteString(constants.PlayerStyle.Render("\nYou"))

    // ... To then build a string out of
    for _, peer := range arr {
        sb.WriteString("\n")
		sb.WriteString(peer.style.Render(peer.name))
    }

    return sb.String()
}
