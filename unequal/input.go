package unequal

func (g *Game) AddValue(w int, h int, value int) {
	g.InitValues[w][h] = value
}
