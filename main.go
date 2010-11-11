package main

import (
	"./player"
	"./tournament"
	"fmt"
	"rand"
	"time"
)

const ITERATIONS = 10000
const CHECKPOINT = 10000
const GAMES = 7

func simulate(array player.Array, tourney tournament.Tournament) {
	dist := make([]float64, 3)

	fmt.Println("Games , Win   , Depth , DepthExp")

	for g := 1; g <= GAMES; g += 2 {
		tourney.SetMatch(tournament.Match{g})

		for i := 0; i < ITERATIONS; i++ {
			result := tourney.Run(array)
			// 		fmt.Println(result)
			dist[0] += array.DistanceByFirst(result)
			dist[1] += array.DistanceByDepth(result)
			dist[2] += array.DistanceByDepthExponential(result)

			if i % CHECKPOINT == 0 && i > 0 {
				fmt.Println(i, dist)
			}
		}

		for i := range dist {
			dist[i] /= ITERATIONS
		}

		fmt.Printf("%5d", g)
		for _, d := range dist {
			fmt.Printf(" , %5.2f", d)
		}
		fmt.Println()
	}
}

func main() {
	rand.Seed(time.Seconds())

	array := player.NewArray(128)

	tourney := &tournament.SingleElimination{}

	simulate(array, tourney)
}