package unequal

import (
	"math/rand"
	"time"
)

func RandInit() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateResult() [][]int {
	matrix := make([][]int, 4)
	for i := 0; i < 4; i++ {
		matrix[i] = make([]int, 4)
	}

	availableValues := []int{1, 2, 3, 4}

	// Fill the matrix using backtracking
	fillMatrix(matrix, 0, availableValues)

	return matrix
}

func fillMatrix(matrix [][]int, index int, availableValues []int) bool {
	if index == 16 {
		return true
	}

	row := index / 4
	col := index % 4

	// Shuffle the available values to add randomness
	shuffle(availableValues)

	for _, value := range availableValues {
		if isSafe(matrix, row, col, value) {
			matrix[row][col] = value
			if fillMatrix(matrix, index+1, availableValues) {
				return true
			}
			matrix[row][col] = 0 // Backtrack
		}
	}

	return false
}

func isSafe(matrix [][]int, row, col, value int) bool {
	// Check if the value already exists in the row or column
	for i := 0; i < 4; i++ {
		if matrix[row][i] == value || matrix[i][col] == value {
			return false
		}
	}

	return true
}

// shuffle using Fisher–Yates algorithm
func shuffle(arr []int) {
	n := len(arr)
	for i := n - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

// Generate init values
func GenerateInit(result [][]int) ([][]int, [][]string, [][]string) {
	boardSize := len(result)

	// Create new slices without elements visible
	initValues := make([][]int, boardSize)
	for i := 0; i < boardSize; i++ {
		initValues[i] = []int{0, 0, 0, 0}
	}

	initInequalitiesRows := make([][]string, boardSize)
	for i := 0; i < boardSize; i++ {
		initInequalitiesRows[i] = []string{"", "", ""}
	}

	initInequalitiesColumns := make([][]string, boardSize)
	for i := 0; i < boardSize; i++ {
		initInequalitiesColumns[i] = []string{"", "", ""}
	}

	// Fill init

	// Fill board
	for r := 0; r < boardSize; r++ {
		index := rand.Intn(boardSize)
		initValues[r][index] = result[r][index]
	}

	//Fill inequalities rows
	for r := 0; r < boardSize; r++ {
		index := rand.Intn(boardSize - 1)
		if result[index][r] > result[index+1][r] {
			initInequalitiesRows[r][index] = "left"
		} else {
			initInequalitiesRows[r][index] = "right"
		}
	}

	//Fill inequalities columns
	for c := 0; c < boardSize; c++ {
		index := rand.Intn(boardSize - 1)
		if result[c][index] > result[c][index+1] {
			initInequalitiesColumns[c][index] = "up"
		} else {
			initInequalitiesColumns[c][index] = "down"
		}
	}

	return initValues, initInequalitiesRows, initInequalitiesColumns
}
