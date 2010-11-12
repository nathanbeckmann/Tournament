package main

import (
	"fmt"
)

const PLAYERS=64
const ROUNDS=18

func main() {
	rounds := make([][]float, ROUNDS)

	rounds[0] = make([]float, 1)
	rounds[0][0] = PLAYERS
	fmt.Println(0, rounds[0])

	for i := 1; i < ROUNDS; i++ {
		rounds[i] = make([]float, len(rounds[i-1])+1)
		for j := range rounds[i] {
			if j == 0 {
				rounds[i][0] = rounds[i-1][0]/2
			} else if j == len(rounds[i])-1 {
				rounds[i][j] = rounds[i-1][j-1] / 2
			} else {
				rounds[i][j] = rounds[i-1][j-1] / 2 + rounds[i-1][j] / 2
			}
		}
		fmt.Println(i, rounds[i])
	}
}