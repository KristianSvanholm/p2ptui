package mines

import (
	"math/rand"
	"p2p/src/constants"
	"p2p/src/structs"
	"sync"
)

type Field struct {
    FirstMove bool
    TotalMines int
    TotalFlags int
    TotalRevealed int
    TotalCells int
    mutex sync.Mutex
    Field [][]structs.Cell 
}

func InitField(size int) Field {

    f := Field{
        FirstMove: true,
        TotalMines: 0,
        TotalFlags: 0,
        TotalRevealed: 0,
        TotalCells: size*size,
        Field: make([][]structs.Cell, size),
    }

	for i := range f.Field {
		f.Field[i] = make([]structs.Cell, size)
	}

    return f
}

func (f *Field) setFlag(c structs.Coords) {

    f.mutex.Lock()
    defer f.mutex.Unlock()

	cell := &f.Field[c.X][c.Y]

	if cell.Revealed {
		return
	}

	cell.Flagged = !cell.Flagged

	// Update flag count
	if cell.Flagged {
		f.TotalFlags++
	} else {
		f.TotalFlags--
	}

	f.checkWin()
	//sendChanges()
}

func (f *Field) openCell(c *structs.Coords, name string) {

    f.mutex.Lock()
    defer f.mutex.Unlock()

	cell := &f.Field[c.X][c.Y]
	if cell.Revealed {
		return
	}

	if f.FirstMove {
        f.FirstMove = false
		f.PlantMines(c)
		f.CalculateCells()
		f.flip(c)
	} else {
		if cell.Mine {
			f.explode(name)
		} else {
			f.flip(c)
		}
	}

    f.checkWin()
	//sendChanges()
}

func (f *Field) checkWin() {
	if (f.TotalCells == f.TotalRevealed+f.TotalMines) && f.TotalMines == f.TotalFlags {
		//SystemMessage("Win!")
		InitField(constants.Size)
	}
}

/*
func sendChanges() {
	msg := structs.ClientMsg{MsgType: "Update", MsgData: Field}
	sendToAll(&msg)
}
*/

func (f *Field) explode(name string) {
	//SystemMessage(fmt.Sprintf("Loss! - %s fucked up", name))
	*f = InitField(constants.Size)
}

func (f *Field) flip(c *structs.Coords) {
    
    f.mutex.Lock()
    defer f.mutex.Unlock()

	size := len(f.Field)

	cell := &f.Field[c.X][c.Y]
	cell.Count = cell.TrueCount
	cell.Revealed = true
	cell.Flagged = false
	f.TotalRevealed++

	if cell.TrueCount != 0 {
		return
	}

	mutX := []int{-1, 0, 1}
	mutY := []int{-1, 0, 1}

	for _, mX := range mutX {
		for _, mY := range mutY {
			mutC := structs.Coords{X: c.X + mX, Y: c.Y + mY}
			if validCell(&mutC, size) && !f.Field[mutC.X][mutC.Y].Revealed {
				f.flip(&mutC)
			}
		}
	}
}

func validCell(c *structs.Coords, size int) bool {
	return !(c.X < 0 || c.Y < 0 || c.X == size || c.Y == size)
}

func (f * Field) PlantMines(c *structs.Coords) {

    f.mutex.Lock()
    defer f.mutex.Unlock()

	for x, row := range f.Field {
		for y := range row {
			if x == c.X && y == c.Y || surroundsCell(c, x, y) {
				continue
			}
    
			if rand.Intn(10) == 1 {
				f.Field[x][y].Mine = true
				f.TotalMines++
			}
		}
	}
}

func (f * Field) CalculateCells() {

    f.mutex.Lock()
    defer f.mutex.Unlock()

	for x, row := range f.Field {
		for y, cell := range row {
			if !cell.Mine {
				f.Field[x][y].TrueCount = f.cellTotal(x, y)
			}
		}
	}
}

func surroundsCell(cell *structs.Coords, x int, y int) bool {
	mutx := []int{-1, 0, 1}
	muty := []int{-1, 0, 1}
	for _, mX := range mutx {
		for _, mY := range muty {
			tempx := cell.X + mX
			tempy := cell.Y + mY
			if x == tempx && y == tempy {
				return true
			}
		}
	}
	return false
}

func (f *Field) cellTotal(x int, y int) int {
	size := len(f.Field)
	total := 0
	mutx := []int{-1, 0, 1}
	muty := []int{-1, 0, 1}

	for _, mX := range mutx {
		for _, mY := range muty {
			tempx := x + mX
			tempy := y + mY
			if tempx < 0 || tempy < 0 || tempy == size || tempx == size || (mX == 0 && mY == 0) {
				continue
			}
			if f.Field[tempx][tempy].Mine {
				total++
			}
		}
	}

	return total
}
