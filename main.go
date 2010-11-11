package main

import (
	"./player"
	"./tournament"
	"fmt"
// 	"runtime"
	"container/vector"
	"sort"
)

const ITERATIONS = 10000
const GAMES = 7

func simulate(retch chan string, array player.Array, tourney tournament.Tournament) {
	out := fmt.Sprintf("%T\n", tourney)
	out = fmt.Sprintln(out, "Games , Win   , Depth , DepthExp")
	subch := make(chan string)

	for g := 1; g <= GAMES; g += 2 {
		go func(g int) {
			dist := make([]float64, 3)

			match := tournament.BestOfMatch{g}
			out := ""

			for i := 0; i < ITERATIONS; i++ {
				result := tourney.Run(array, match)
				dist[0] += array.DistanceByFirst(result)
				dist[1] += array.DistanceByDepth(result)
				dist[2] += array.DistanceByDepthExponential(result)
			}

			for i := range dist {
				dist[i] /= ITERATIONS
			}

			out = fmt.Sprintf("%5d", g)
			for _, d := range dist {
				out = fmt.Sprintf("%s , %5.2f", out, d)
			}
			out = fmt.Sprintln(out)

			subch <- out
		}(g)
	}

	results := vector.StringVector{}

	for g := 1; g <= GAMES; g += 2 {
		results.Push(<-subch)
	}
	sort.Sort(&results)
	for _, s := range results {
		out = fmt.Sprint(out, s)
	}
	out = fmt.Sprintln(out)

	retch <- out
}

func main() {
// 	runtime.GOMAXPROCS(8)
	array := player.NewArray(512)
	ch := make(chan string)

	single := &tournament.SingleElimination{}
  	go simulate(ch, array, single)

	double := &tournament.DoubleElimination{}
 	go simulate(ch, array, double)

	double_extended := &tournament.DoubleEliminationExtendedSeries{}
 	go simulate(ch, array, double_extended)
// 	tournament.Tournament(double_extended).SetMatch(tournament.BestOfMatch{3})
// 	for i := 0; i < 1000; i++ {
// 		tournament.Tournament(double_extended).Run(array)
// 	}
// 	fmt.Println(tournament.NumExtendedSeries)

	for i := 0; i < 3; i++ {
		fmt.Print(<-ch)
	}
}