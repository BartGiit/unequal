package unequal

import (
	"fmt"
	"image/color"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 405
	ScreenHeight = 600
	RectSize     = 80
	RectMargin   = 15
)

type Game struct {
	boardImage              *ebiten.Image
	InitValues              [][]int
	InitInequalitiesRows    [][]string
	InitInequalitiesColumns [][]string
	mouseX, mouseY          int
	ClickedCol              int
	ClickedRow              int
	ClickProcessed          bool
	InputNumberInvalid      bool
	InputNumber             int
	Finished                bool
	TimeElapsed             int
}

func (g *Game) Update() error {
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	// Check for mouse click
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !g.ClickProcessed {
		clickedCol, clickedRow := g.getClickedRect()
		if clickedCol != -1 && clickedRow != -1 {
			g.ClickedCol = clickedCol
			g.ClickedRow = clickedRow
			g.InputNumberInvalid = true
		} else {
			g.ClickedCol = -1
			g.ClickedRow = -1
			g.InputNumberInvalid = false
		}
		g.ClickProcessed = true
	} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.ClickProcessed = false
	}

	// Process keyboard input if a rectangle is clicked
	if g.ClickedCol != -1 && g.ClickedRow != -1 && g.InputNumberInvalid {
		g.InputNumber = g.keyPressed()
		if g.InputNumber != -1 {
			g.AddValue(g.ClickedCol, g.ClickedRow, g.InputNumber)
		}
	}

	//Checks for condition terminating the game
	_, inequalitiesRows := g.CorrectInequalitiesRows()
	_, inequalitiesColumns := g.CorrectInequalitiesColumns()
	_, valuesRows := g.CorrectValuesRows()
	_, valuesColumns := g.CorrectValuesColumns()
	filledCount := 0
	for _, value := range g.InitValues {
		for _, number := range value {
			if number != 0 {
				filledCount++
			}
		}
	}
	if inequalitiesRows && inequalitiesColumns && valuesRows && valuesColumns && filledCount == 16 {
		g.Finished = true
	}
	return nil
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{245, 230, 218, 255})

	if g.Finished {
		finishInfo := fmt.Sprintf("Game finished in: %02d:%02d", g.TimeElapsed/60, g.TimeElapsed%60)
		ebitenutil.DebugPrintAt(screen, finishInfo, 125, 220)

		restartMsg := "Click anywhere to restart"
		ebitenutil.DebugPrintAt(screen, restartMsg, 120, 250)

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			g.RestartGame()
		}
	} else {
		// Print time elapsed
		timeText := fmt.Sprintf("Time: %02d:%02d", g.TimeElapsed/60, g.TimeElapsed%60)
		ebitenutil.DebugPrintAt(screen, timeText, 160, 420)

		// Get the duplicate indexes from CorrectValuesRows and CorrectValuesColumns
		duplicateRows, _ := g.CorrectValuesRows()
		duplicateColumns, _ := g.CorrectValuesColumns()

		// Get the inequality indexes from CorrectInequalitiesRows and CorrectInequalitiesColumns
		inequalityRows, _ := g.CorrectInequalitiesRows()
		inequalityColumns, _ := g.CorrectInequalitiesColumns()

		// Draw rectangles and numbers
		for w := 0; w < 4; w++ {
			for h := 0; h < 4; h++ {
				rect := ebiten.NewImage(RectSize, RectSize)

				// Check duplicates in rows
				isDuplicateRow := false
				for _, index := range duplicateRows[h] {
					if index == w {
						isDuplicateRow = true
						break
					}
				}

				// Check duplicates in columns
				isDuplicateColumn := false
				for _, index := range duplicateColumns[w] {
					if index == h {
						isDuplicateColumn = true
						break
					}
				}

				// Check inequalities rows
				isInequalityRow := false
				if g.InitValues[w][h] != 0 {
					for _, index := range inequalityRows[h] {
						if index == w {
							isInequalityRow = true
							break
						}
					}
				}
				// Check inequalities columns
				isInequalityColumn := false
				if g.InitValues[w][h] != 0 {
					for _, index := range inequalityColumns[w] {
						if index == h {
							isInequalityColumn = true
							break
						}
					}
				}

				// Check if the mouse is over the current rectangle
				isMouseOver := g.isMouseOverRect(w, h)

				// Color based on duplicates, inequalities or mouse over
				if isDuplicateRow || isDuplicateColumn {
					if isMouseOver {
						rect.Fill(color.RGBA{255, 0, 0, 255}) // Duplicates with mouse over
					} else {
						rect.Fill(color.RGBA{255, 0, 0, 100}) // Duplicates
					}
				} else if isInequalityRow || isInequalityColumn {
					if isMouseOver {
						rect.Fill(color.RGBA{255, 0, 0, 255}) // Inequality with mouse over
					} else {
						rect.Fill(color.RGBA{255, 0, 0, 100}) // Inequality violation
					}
				} else if isMouseOver {
					rect.Fill(color.RGBA{0, 0, 0, 255}) // Mouse over
				} else {
					rect.Fill(color.RGBA{1, 1, 1, 100}) // Default
				}
				// Draw the number
				number := strconv.Itoa(g.InitValues[w][h])
				ebitenutil.DebugPrintAt(rect, number, 36, 36)
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(w*(RectSize+RectMargin)+RectMargin), float64(h*(RectSize+RectMargin)+RectMargin))
				screen.DrawImage(rect, op)
			}
		}
		// Draw inequalities
		for w := 0; w < 4; w++ {
			for h := 0; h < 4; h++ {
				// Inequality for rows
				if w < 3 && g.InitInequalitiesRows[h][w] != "" {
					arrowX := float64((w)*(RectSize+RectMargin) + RectMargin + RectSize + 2)
					arrowY := float64(h*(RectSize+RectMargin) + RectMargin + 32)
					if g.InitInequalitiesRows[h][w] == "left" {
						ebitenutil.DebugPrintAt(screen, ">", int(arrowX), int(arrowY))
					} else if g.InitInequalitiesRows[h][w] == "right" {
						ebitenutil.DebugPrintAt(screen, "<", int(arrowX), int(arrowY))
					}
				}

				// Inequality for columns
				if h < 3 && g.InitInequalitiesColumns[w][h] != "" {
					arrowX := float64(w*(RectSize+RectMargin) + RectMargin + 38)
					arrowY := float64((h)*(RectSize+RectMargin) + RectSize + RectMargin)
					if g.InitInequalitiesColumns[w][h] == "up" {
						ebitenutil.DebugPrintAt(screen, "v", int(arrowX), int(arrowY))
					} else if g.InitInequalitiesColumns[w][h] == "down" {
						ebitenutil.DebugPrintAt(screen, "^", int(arrowX), int(arrowY))
					}
				}
			}
		}
	}
}

func (g *Game) isMouseOverRect(w, h int) bool {
	rectX := w*(RectSize+RectMargin) + RectMargin
	rectY := h*(RectSize+RectMargin) + RectMargin
	return g.mouseX >= rectX && g.mouseX < rectX+RectSize && g.mouseY >= rectY && g.mouseY < rectY+RectSize
}

func (g *Game) getClickedRect() (int, int) {
	for w := 0; w < 4; w++ {
		for h := 0; h < 4; h++ {
			rectX := w*(RectSize+RectMargin) + RectMargin
			rectY := h*(RectSize+RectMargin) + RectMargin
			if g.mouseX >= rectX && g.mouseX < rectX+RectSize && g.mouseY >= rectY && g.mouseY < rectY+RectSize {
				return w, h
			}
		}
	}
	return -1, -1
}

func (g *Game) keyPressed() int {
	inputNumber := -1
	for _, k := range inpututil.PressedKeys() {
		if k >= ebiten.Key0 && k <= ebiten.Key9 {
			num := int(k - ebiten.Key0)
			if num >= 1 && num <= 4 {
				inputNumber = num
				break
			}
		}
	}
	return inputNumber
}

// Timer implementation using goroutine and channels
func (g *Game) StartTimer() {
	timer := time.NewTimer(time.Second)
	ticker := time.NewTicker(time.Second)
	quit := make(chan bool)

	go func() {
		for {
			select {
			case <-timer.C:
				// If game not finished
				if !g.Finished {
					g.updateTime()
					timer.Reset(time.Second)
				}

				// If game finished
				if g.Finished {
					quit <- true
					return
				}

			case <-ticker.C:
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func (g *Game) updateTime() {
	g.TimeElapsed++
}

func (g *Game) RestartGame() {
	result := GenerateResult()
	initValues, initInequalitiesRows, initInequalitiesColumns := GenerateInit(result)
	g.InitValues = initValues
	g.InitInequalitiesRows = initInequalitiesRows
	g.InitInequalitiesColumns = initInequalitiesColumns
	g.ClickedCol = -1
	g.ClickedRow = -1
	g.ClickProcessed = false
	g.InputNumberInvalid = false
	g.InputNumber = 0
	g.TimeElapsed = 0
	g.Finished = false

	g.StartTimer() // Start the timer
}
