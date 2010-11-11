package main

import (
	"./player"
	"./tournament"
	"fmt"
 	"runtime"
	"container/vector"
	"sort"
)

const PLAYERS = 512
const ITERATIONS = 50000
const GAMES = 3

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

func measure_extended(double_extended tournament.Tournament) {
	tournament.NumMatches = 0
	tournament.NumExtendedSeries = 0
	for i := 0; i < ITERATIONS; i++ {
 		tournament.Tournament(double_extended).Run(array, tournament.BestOfMatch{3})
// 		fmt.Println(tournament.NumMatches)
	}
	fmt.Println("avg extended ", float(tournament.NumExtendedSeries) / ITERATIONS)
	fmt.Println("avg corrected ", float(tournament.NumCorrections) / ITERATIONS)
	fmt.Println("avg uncorrected ", float(tournament.NumUncorrections) / ITERATIONS)
	fmt.Println("avg injustice ", float(tournament.NumInjustices) / ITERATIONS)
	fmt.Println("fraction extended ", float(tournament.NumExtendedSeries) / ((2 * PLAYERS - 1) * ITERATIONS))
	fmt.Println("fraction corrected ", float(tournament.NumCorrections) / float(tournament.NumExtendedSeries))
	fmt.Println("fraction uncorrected ", float(tournament.NumUncorrections) / float(tournament.NumExtendedSeries))
	fmt.Println("fraction injustice ", float(tournament.NumInjustices) / float(tournament.NumExtendedSeries))
}

func main() {
  	runtime.GOMAXPROCS(1)
	array := player.NewArray(PLAYERS)
	ch := make(chan string)

	fmt.Println(PLAYERS, " players")
	fmt.Println(ITERATIONS, " iterations")

	single := &tournament.SingleElimination{}
  	go simulate(ch, array, single)

	double := &tournament.DoubleElimination{}
 	go simulate(ch, array, double)

 	double_extended := &tournament.DoubleEliminationExtendedSeries{}
 	go simulate(ch, array, double_extended)

	roundrobin := &tournament.RoundRobin{}
	go simulate(ch, array, roundrobin)

	for i := 0; i < 4; i++ {
		fmt.Print(<-ch)
	}
}