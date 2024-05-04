package mines

import (
	"math/rand"
	"p2p/src/constants"
	"p2p/src/structs"
	"sync"
)

type Field struct {
	FirstMove     bool             `json:"firstMove"`
	TotalMines    int              `json:"totalMines"`
	TotalFlags    int              `json:"totalFlags"`
	TotalRevealed int              `json:"totalRevealed"`
	TotalCells    int              `json:"totalCells"`
	mutex         sync.Mutex       `json:"-"`
	Field         [][]structs.Cell `json:"field"`
}

// Generate the initial blank field data structure.
func InitField(size int) *Field {

	f := Field{
		FirstMove:     true,
		TotalMines:    0,
		TotalFlags:    0,
		TotalRevealed: 0,
		TotalCells:    size * size,
		Field:         make([][]structs.Cell, size),
	}

	for i := range f.Field {
		f.Field[i] = make([]structs.Cell, size)
	}

	return &f
}

// Any user plants a flag in the field
func (f *Field) SetFlag(c structs.Coords) bool {

	// Disallow flags in first move.
	if f.FirstMove {
		return false
	}

	f.mutex.Lock()
	defer f.mutex.Unlock()

	cell := &f.Field[c.X][c.Y]

	// Check if cell has already been revealed.
	if cell.Revealed {
		return false
	}

	cell.Flagged = !cell.Flagged

	// Update flag count
	if cell.Flagged {
		f.TotalFlags++
	} else {
		f.TotalFlags--
	}

	return f.checkWin()
}

// Any user digs a cell.
func (f *Field) Dig(c *structs.Coords, rng *rand.Rand) constants.DigEvent {

	var result constants.DigEvent = constants.Nothing // Deault to nothing happening.

	cell := &f.Field[c.X][c.Y]

	// Do a chunck dig around the cell
	if cell.Revealed || cell.Flagged {
		return f.digAround(c, rng)
	}

	// If it is the first move, plant the mines and add numbers to the field.
	if f.FirstMove {
		f.FirstMove = false
		f.PlantMines(c, rng)
		f.CalculateCells()
		f.flip(c) // Flip the cell
	} else { // If not, check if a mine was hit.

		if cell.Mine {
			go f.explode()
			result = constants.Landmine
		} else {
			f.flip(c) // If no mine was hit, flip the cell
		}
	}

	// Check if the game is won.
	if f.checkWin() {
		result = constants.Win
	}

	return result
}

// Digs all cells around a cell if they are not already revealed or flagged.
func (f *Field) digAround(c *structs.Coords, rng *rand.Rand) constants.DigEvent {
	// "Neighbor-matrix"
	mutX := []int{-1, 0, 1}
	mutY := []int{-1, 0, 1}

	result := constants.Nothing // Default to nothing happening.

	// Loop over neighbour cells
	for _, mX := range mutX {
		for _, mY := range mutY {
			mutC := structs.Coords{X: c.X + mX, Y: c.Y + mY}
			if !(validCell(&mutC) && !f.Field[mutC.X][mutC.Y].Revealed && !f.Field[mutC.X][mutC.Y].Flagged) {
				continue
			}

			currResult := f.Dig(&mutC, rng)
			if currResult > result { // Nothing < Win < Landmine
				result = currResult
			}
		}
	}
	return result
}

// Simply checks if the game is won through some math & if so, resets the field.
func (f *Field) checkWin() bool {
	if (f.TotalCells == f.TotalRevealed+f.TotalMines) && f.TotalMines == f.TotalFlags {
		go func() { *f = *InitField(constants.Size) }() // This needs to happen concurrently
		return true
	}
	return false
}

// Lose the game & reset the field.
func (f *Field) explode() {
	*f = *InitField(constants.Size)
}

// Recursively flip a cell.
// If this cell has 0 mines as neighbors and any neighbour cell is not revealed, flip that one as well.
func (f *Field) flip(c *structs.Coords) {

	// Flip current cell.
	cell := &f.Field[c.X][c.Y]
	cell.Revealed = true
	cell.Flagged = false
	f.TotalRevealed++

	// If current cell not equal to 0, abort flipping neighbor cells.
	if cell.Count != 0 {
		return
	}

	// "Neighbor-matrix"
	mutX := []int{-1, 0, 1}
	mutY := []int{-1, 0, 1}

	// Loop over neighbors.
	for _, mX := range mutX {
		for _, mY := range mutY {
			mutC := structs.Coords{X: c.X + mX, Y: c.Y + mY}
			if validCell(&mutC) && !f.Field[mutC.X][mutC.Y].Revealed {
				f.flip(&mutC) // Flip neighbor cell.
			}
		}
	}
}

// Check if a cell exists within the bounds of the field.
func validCell(c *structs.Coords) bool {
	return !(c.X < 0 || c.Y < 0 || c.X == constants.Size || c.Y == constants.Size)
}

// Plant the mines in the field except for the given opening-cell and its neighbors.
func (f *Field) PlantMines(c *structs.Coords, rng *rand.Rand) {

	// Loop over entire field.
	for x, row := range f.Field {
		for y := range row {

			// Ignore opening-cell and its neighbors.
			if x == c.X && y == c.Y || surroundsCell(c, x, y) {
				continue
			}

			// Randomly plant a mine here.
			if rng.Intn(100) <= constants.Density {
				f.Field[x][y].Mine = true
				f.TotalMines++
			}
		}
	}
}

// Calculate the total amount of mines around all cells.
func (f *Field) CalculateCells() {

	for x, row := range f.Field {
		for y, cell := range row {
			if !cell.Mine {
				f.Field[x][y].Count = f.cellTotal(x, y)
			}
		}
	}
}

// Check if coordinates are neighbors to cell.
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

// Calculate total mines around coordinates.
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
