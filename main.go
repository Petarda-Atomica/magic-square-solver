package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"strings"
)

// Boolean algebra
func xor(a, b bool) bool {
	return (a && !b) || (!a && b)
}

// Conversions
func sprintSquare(square [][]int, n int) [][]string {
	var output [][]string
	for y := 0; y < n; y++ {
		this := []string{}
		for x := 0; x < n; x++ {
			this = append(this, fmt.Sprint(square[y][x]))
		}
		output = append(output, this)
	}

	return output
}

func trimSquare(square [][]int, n int) [][]int {
	var output [][]int
	for y := 0; y < n; y++ {
		this := []int{}
		for x := 0; x < n; x++ {
			this = append(this, square[y][x])
		}
		output = append(output, this)
	}

	return output
}

func swap(square [][]int, a image.Point, b image.Point) {
	aux := square[a.Y][a.X]
	square[a.Y][a.X] = square[b.Y][b.X]
	square[b.Y][b.X] = aux
}

// Draw operations
func buildTable(values [][]string, additive int) string {
	var output string
	// Variables
	var rows, cols []int
	for i := 0; i < len(values); i++ {
		rows = append(rows, 0)
	}
	for i := 0; i < len(values[0]); i++ {
		cols = append(cols, 0)
	}

	// Find variables
	for x := range values {
		for y := range values[x] {
			// Get widths
			lines := strings.Split(values[x][y], "\n")
			for _, val := range lines {
				cols[y] = max(cols[y], len(val))
			}
			rows[x] = max(rows[x], len(lines))
		}
	}

	// Draw table head
	output += "\u250c"
	for i := range cols {
		for j := 0; j < cols[i]+2*additive; j++ {
			output += "\u2500"
		}
		if i != len(cols)-1 {
			output += "\u252c"
		}
	}
	output += "\u2510\n"

	// Draw table
	for y := range rows {
		// Draw data
		output += "\u2502"
		for x := range cols {
			// Draw start additive
			for i := 0; i < additive; i++ {
				output += " "
			}

			// Draw start compensation
			for i := 0; i < (cols[x]-len(values[y][x]))/2; i++ {
				output += " "
			}

			// Draw value
			output += values[y][x]

			// Draw end compensation
			for i := 0; i < (cols[x]-len(values[y][x]))-(cols[x]-len(values[y][x]))/2; i++ {
				output += " "
			}

			// Draw end additive
			for i := 0; i < additive; i++ {
				output += " "
			}

			// Draw separator
			output += "\u2502"

		}
		output += "\n"

		// Draw separator
		if y == len(rows)-1 {
			continue
		}
		output += "\u2502"
		for i := range cols {
			for j := 0; j < cols[i]+2*additive; j++ {
				output += "\u2500"
			}
			if i != len(cols)-1 {
				output += "\u253c"
			}
		}
		output += "\u2502\n"
	}

	// Draw table footer
	output += "\u2514"
	for i := range cols {
		for j := 0; j < cols[i]+2*additive; j++ {
			output += "\u2500"
		}
		if i != len(cols)-1 {
			output += "\u2534"
		}
	}
	output += "\u2518\n"

	return output
}

// Square algorithms
func oddSquare(n int) [][]int {
	// Check if solvable
	if n%2 == 0 {
		return nil
	}

	// Define square
	var square [][]int
	for y := 0; y <= n; y++ {
		this := []int{}
		for x := 0; x <= n; x++ {
			this = append(this, 0)
		}
		square = append(square, this)
	}

	// Solve square
	pointer := image.Point{X: (n - 1) / 2, Y: 0}
	for i := 1; i <= n*n; i++ {
		// Place number
		square[pointer.Y][pointer.X] = i

		// Move pointer
		pointer.X++
		pointer.Y--

		// Check collisions
		if pointer.X >= n && pointer.Y < 0 {
			pointer.X--
			pointer.Y++
			pointer.Y++
		} else if pointer.Y < 0 {
			pointer.Y = n - 1
		} else if pointer.X >= n {
			pointer.X = 0
		} else if square[pointer.Y][pointer.X] != 0 {
			pointer.X--
			pointer.Y++
			pointer.Y++
		}
	}

	// Convert square

	return trimSquare(square, n)
}

func doublyEvenSquare(n int) [][]int {
	// Check if solvable
	if n%4 != 0 {
		return nil
	}

	// Define square
	var square [][]int
	for y := 0; y <= n; y++ {
		this := []int{}
		for x := 0; x <= n; x++ {
			this = append(this, 0)
		}
		square = append(square, this)
	}

	// Define variables
	regionSize := n / 4

	// Solve square
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			if xor((x/regionSize == 0 || x/regionSize == 3) && (y/regionSize == 0 || y/regionSize == 3), (x/regionSize == 1 || x/regionSize == 2) && (y/regionSize == 1 || y/regionSize == 2)) {
				square[y][x] = y*n + x + 1
			} else {
				square[y][x] = n*n - (y*n + x)
			}
		}
	}

	return trimSquare(square, n)
}

func singlyEvenSquare(n int) [][]int {
	// Check if solvable
	if !(n%2 == 0 && n%4 != 0) {
		return nil
	}

	// Define square
	var square [][]int
	for y := 0; y <= n; y++ {
		this := []int{}
		for x := 0; x <= n; x++ {
			this = append(this, 0)
		}
		square = append(square, this)
	}

	// Solve the odd part of the square
	quadron := oddSquare(n / 2)
	step := (n * n / 4)

	/*Quadron distribution:
	┌─────┬─────┐
	│  1  │  2  │
	│─────┼─────│
	│  3  │  4  │
	└─────┴─────┘
	*/

	// Pseudo-solve first quadron
	for y := 0; y < n/2; y++ {
		for x := 0; x < n/2; x++ {
			square[y][x] = quadron[y][x]
		}
	}

	// Pseudo-solve second quadron
	for y := 0; y < n/2; y++ {
		for x := n / 2; x < n; x++ {
			square[y][x] = quadron[y][x-n/2] + 2*step
		}
	}

	// Pseudo-solve third quadron
	for y := n / 2; y < n; y++ {
		for x := 0; x < n/2; x++ {
			square[y][x] = quadron[y-n/2][x] + 3*step
		}
	}

	// Pseudo-solve fourth quadron
	for y := n / 2; y < n; y++ {
		for x := n / 2; x < n; x++ {
			square[y][x] = quadron[y-n/2][x-n/2] + step
		}
	}

	// Swap right collumn
	for x := n - 1; x >= n-(n-6)/4; x-- {
		for y := 0; y < n; y++ {
			if y < n/2 {
				square[y][x] -= step
			} else {
				square[y][x] += step
			}
		}
	}

	// Swap left collumn
	for x := 0; x < (n-2)/4; x++ {
		for y := 0; y < n/2; y++ {
			if y == (n/2-1)/2 && x == 0 {
				continue
			}

			swap(square, image.Pt(x, y), image.Pt(x, y+n/2))
		}
	}
	swap(square, image.Pt((n-2)/4, (n/2-1)/2), image.Pt((n-2)/4, (n/2-1)/2+n/2))

	return trimSquare(square, n)
}

// Checks
func checkMagic(square [][]int) {
	// Integrity tests
	fatalDamage := false
	for y := 0; y < len(square); y++ {
		if len(square) != len(square[y]) {
			log.Printf("Failed integrity check at row %d. Required: %d Got: %d\n", y, len(square), len(square[y]))
			if len(square[y]) < len(square) {
				fatalDamage = true
				log.Println("Integrity compromised beyond repair")
			}
		}
	}

	// Check if we can still continue
	if fatalDamage {
		log.Println("Continuing test from here on now is considered unsafe. Exiting...")
		return
	}

	// Magic constant check row
	n := len(square)
	magicConstant := n * (n*n + 1) / 2
	for y := 0; y < n; y++ {
		this := 0
		for x := 0; x < n; x++ {
			this += square[y][x]
		}
		if this != magicConstant {
			log.Printf("Mismatched row-sum on row %d. Required: %d Got: %d\n", y, magicConstant, this)
		}
	}

	// Magic constant check collumn
	for x := 0; x < n; x++ {
		this := 0
		for y := 0; y < n; y++ {
			this += square[y][x]
		}
		if this != magicConstant {
			log.Printf("Mismatched collumn-sum on collumn %d. Required: %d Got: %d\n", x, magicConstant, this)
		}
	}

	// Magic constant check diagonal - 1
	sum := 0
	for i := 0; i < n; i++ {
		sum += square[i][i]
	}
	if sum != magicConstant {
		log.Printf("Mismatched diagonal-sum on diagonal 1. Required: %d Got: %d\n", magicConstant, sum)
	}

	// Magic constant check diagonal - 2
	sum = 0
	for i := 0; i < n; i++ {
		sum += square[n-i-1][i]
	}
	if sum != magicConstant {
		log.Printf("Mismatched diagonal-sum on diagonal 2. Required: %d Got: %d\n", magicConstant, sum)
	}

}

// Main function
func main() {
	// Get CLI arguments
	size := flag.Int("size", 3, "The size of the magic square to generate.")
	format := flag.String("format", "table", "What format to use to display the magic square: table OR json.")
	output := flag.String("output", "", "The name of the output file. Leave empty in order to print to CLI.")
	flag.Parse()

	// Make sure the square exist
	if *size < 3 {
		log.Fatal(fmt.Errorf("the value for size can't be %d. the value must be greater or equal to 3", *size))
	}

	// Calculate square
	var square [][]int
	if *size%2 == 1 {
		square = oddSquare(*size)
	} else if *size%4 == 0 {
		square = doublyEvenSquare(*size)
	} else {
		square = singlyEvenSquare(*size)
	}

	// Perform checks
	checkMagic(square)

	// Convert to desired format
	var data []byte
	if *format == "table" {
		data = []byte(buildTable(sprintSquare(square, *size), 2))
	} else {
		raw, err := json.Marshal(square)
		if err != nil {
			log.Fatal(err)
		}
		data = raw
	}

	// Save data
	if *output == "" {
		fmt.Println(string(data))
	} else {
		err := os.WriteFile(*output, data, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}
