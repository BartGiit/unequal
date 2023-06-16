package unequal

// Checks uniqe values in rows

func (g *Game) CorrectValuesRows() ([][]int, bool) {
	result := make([][]int, 4)
	correctValues := true

	for r := 0; r < 4; r++ {
		rowValues := make([]int, 0)
		duplicateIndexes := make(map[int][]int)
		seen := make(map[int]bool)

		for c := range g.InitValues {
			num := g.InitValues[c][r]
			rowValues = append(rowValues, num)

			if num != 0 {
				if seen[num] {
					duplicateIndexes[num] = append(duplicateIndexes[num], c)
				} else {
					seen[num] = true
					duplicateIndexes[num] = []int{c}
				}
			}
		}

		var indexes []int

		for _, dups := range duplicateIndexes {
			if len(dups) >= 2 {
				indexes = append(indexes, dups...)
			}
		}

		if len(indexes) == 0 {
			indexes = append(indexes, -1)
		}

		result[r] = indexes
	}

	for _, value := range result {
		if len(value) > 1 {
			correctValues = false
		}
	}

	return result, correctValues
}

// Checks uniqe values in columns

func (g *Game) CorrectValuesColumns() ([][]int, bool) {
	result := make([][]int, 4)
	correctValues := true

	for c := 0; c < 4; c++ {
		columnValues := make([]int, 0)
		duplicateIndexes := make(map[int][]int)
		seen := make(map[int]bool)

		for r := range g.InitValues[c] {
			num := g.InitValues[c][r]
			columnValues = append(columnValues, num)

			if num != 0 {
				if seen[num] {
					duplicateIndexes[num] = append(duplicateIndexes[num], r)
				} else {
					seen[num] = true
					duplicateIndexes[num] = []int{r}
				}
			}
		}

		var indexes []int

		for _, dups := range duplicateIndexes {
			if len(dups) >= 2 {
				indexes = append(indexes, dups...)
			}
		}

		if len(indexes) == 0 {
			indexes = append(indexes, -1)
		}

		result[c] = indexes
	}

	for _, value := range result {
		if len(value) > 1 {
			correctValues = false
		}
	}

	return result, correctValues
}

// Checks inequalities in rows

func (g *Game) CorrectInequalitiesRows() ([][]int, bool) {
	result := make([][]int, 0)
	correctValues := true

	for r := 0; r < 4; r++ {
		indexes := make([]int, 0)
		for c := 0; c < 3; c++ {
			if g.InitInequalitiesRows[r][c] != "" && g.InitValues[c][r] != 0 && g.InitValues[c+1][r] != 0 {
				if g.InitInequalitiesRows[r][c] == "left" && g.InitValues[c][r] < g.InitValues[c+1][r] {
					indexes = append(indexes, c, c+1)
				}
				if g.InitInequalitiesRows[r][c] == "right" && g.InitValues[c][r] > g.InitValues[c+1][r] {
					indexes = append(indexes, c, c+1)
				}
			}
		}
		result = append(result, indexes)
		if len(indexes) > 0 {
			correctValues = false
		}
	}

	return result, correctValues
}

// Checks inequalities in columns

func (g *Game) CorrectInequalitiesColumns() ([][]int, bool) {
	result := make([][]int, 0)
	correctValues := true

	for r := 0; r < 4; r++ {
		indexes := make([]int, 0)
		for c := 0; c < 3; c++ {
			if g.InitInequalitiesColumns[r][c] != "" && g.InitValues[r][c] != 0 && g.InitValues[r][c+1] != 0 {
				if g.InitInequalitiesColumns[r][c] == "up" && g.InitValues[r][c] < g.InitValues[r][c+1] {
					indexes = append(indexes, c, c+1)
				}
				if g.InitInequalitiesColumns[r][c] == "down" && g.InitValues[r][c] > g.InitValues[r][c+1] {
					indexes = append(indexes, c, c+1)
				}
			}
		}
		result = append(result, indexes)
		if len(indexes) > 0 {
			correctValues = false
		}
	}

	return result, correctValues
}
