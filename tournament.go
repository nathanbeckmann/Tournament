package tournament

import (
	"./player"
	"math"
)

type Tournament interface {
	Run(player.Array) []int
	SetMatch(Match)
}

type tournamentBase struct {
	match Match
}

func (t *tournamentBase) SetMatch(m Match) {
	t.match = m
}

func intlog(a int) int {
	return int(math.Log2(float64(a))+0.5)
}

func initialSeeds(len int) []int {
	array := make([]int, len)

	for j := 0; j < (len / 2); j++ {
		pos := 0
		k := j
		for i := 0; i < intlog(len)-1; i++ {
			pos |= k & 1
			pos <<= 1
			k >>= 1
		}

		array[pos] = j
		array[pos+1] = len-j-1
	}

	return array
}

type Match struct {
	Games int
}

func (m Match) Play(i, j int, array player.Array) (int, int) {
	iwins, jwins := 0, 0
	for {
		if winner, _ := array.Play(i,j); winner == i {
			iwins++
		} else {
			jwins++
		}

		if iwins > m.Games / 2 {
			return i, j
		} else if jwins > m.Games / 2 {
			return j, i
		}
	}
	panic(nil)
}

type SingleElimination struct {
	tournamentBase
}

func (t SingleElimination) Run(array player.Array) []int {
	results := make([]int, array.Len())
	num_ranked := len(results) - 1

	rounds := make([][]int, intlog(array.Len())+1)
	rounds[0] = initialSeeds(array.Len())

	for r := 1; ; r++ {
		if len(rounds[r-1])/2 == 0 {
			break
		}

		rounds[r] = make([]int, len(rounds[r-1])/2)

		for i := range rounds[r] {
			a := rounds[r-1][2 * i]
			b := rounds[r-1][2 * i + 1]

			rounds[r][i], results[num_ranked] = t.match.Play(a, b, array)
			num_ranked--
		}
	}

	results[num_ranked] = rounds[len(rounds)-1][0]

	if num_ranked != 0 {
		panic("Didn't rank all the players!")
	}

	return results
}

type DoubleElimination struct {
	tournamentBase
}

func (t DoubleElimination) Run(array player.Array) []int {
	results := make([]int, array.Len())
	num_ranked := len(results) - 1
	
	winners := make([][]int, intlog(array.Len())+1)
	winners[0] = initialSeeds(array.Len())

	num_ranked--

	return results
}
