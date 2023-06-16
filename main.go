package main

import (
	"log"

	"github.com/BartGiit/unequal/unequal"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	unequal.RandInit()
	result := unequal.GenerateResult()
	initValues, initInequalitiesRows, initInequalitiesColumns := unequal.GenerateInit(result)
	ebiten.SetWindowSize(unequal.ScreenWidth, unequal.ScreenHeight)
	ebiten.SetWindowTitle("Unequal")
	game := &unequal.Game{
		InitValues:              initValues,
		InitInequalitiesRows:    initInequalitiesRows,
		InitInequalitiesColumns: initInequalitiesColumns,
		ClickedCol:              -1,
		ClickedRow:              -1,
		ClickProcessed:          false,
		Finished:                false,
	}
	game.StartTimer()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
