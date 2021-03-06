package tournament

import (
	"./player"
	"math"
	"rand"
	"time"
	"fmt"
	"container/vector"
)

type Tournament interface {
	Run(player.Array, Match) []int
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

type Match interface {
	Play(i, j int, array player.Array, rand *rand.Rand) (int, int, int)
}

type BestOfMatch struct {
	Games int
}

var NumMatches int = 0

func (m BestOfMatch) Play(i, j int, array player.Array, rand *rand.Rand) (int, int, int) {
	NumMatches++
	iwins, jwins := 0, 0
	games := 0
	for {
		games++
		if winner, _ := array.Play(i,j,rand); winner == i {
			iwins++
		} else {
			jwins++
		}

		if iwins > m.Games / 2 {
			return i, j, games
		} else if jwins > m.Games / 2 {
			return j, i, games
		}
	}
	panic(nil)
}

type SingleElimination struct {
}

func (t SingleElimination) Run(array player.Array, match Match) []int {
	rand := rand.New(rand.NewSource(time.Seconds()))
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

			rounds[r][i], results[num_ranked], _ = match.Play(a, b, array, rand)
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
}

func (t DoubleElimination) Run(array player.Array, match Match) []int {
	rand := rand.New(rand.NewSource(time.Seconds()))
	results := make([]int, array.Len())
	num_ranked := len(results) - 1
	
	// Double elim tourney needs 3 brackets -- winners, losers
	// (losers play each other) and loservswinners (losers play
	// people coming down from winners bracket). Technically,
	// losers and loservswinners are both parts of the losers
	// bracket.
	winners := make([][]int, intlog(array.Len())+1)
	losers := make([][]int, intlog(array.Len())+1)
	losersvswinners := make([][]int, intlog(array.Len())+1)

	winners[0] = initialSeeds(array.Len())

	for r := 1; ; r++ {
		if len(winners[r-1])/2 == 0 {
			break
		}

		// Create next round
		winners[r] = make([]int, len(winners[r-1])/2)
		losers[r] = make([]int, len(winners[r-1])/2)
		losersvswinners[r] = make([]int, len(winners[r-1])/2)

		// Fill in this round
		for i := range winners[r] {
			// Play winners bracket, knocking someone down to losersvswinners bracket
			a := winners[r-1][2 * i]
			b := winners[r-1][2 * i + 1]

			winners[r][i], losersvswinners[r][i], _ = match.Play(a, b, array, rand)
		}

		// There is no losers bracket in the first round...
		if r == 1 {
			losers[r] = losersvswinners[r]
			continue
		}

		for i := range losers[r] {
			// Otherwise, we play the losers bracket and then put the losers against the winners
			a := losers[r-1][2 * i]
			b := losers[r-1][2 * i + 1]

			losers[r][i], results[num_ranked], _ = match.Play(a, b, array, rand)
			num_ranked--

			a = losers[r][i]
			b = losersvswinners[r][i]

			losers[r][i], results[num_ranked], _ = match.Play(a, b, array, rand)
			num_ranked--
		}
	}

	// Some special logic for finals -- losers champ has to beat winner twice
	winnerchamp := winners[len(winners)-1][0]
	loserchamp := losers[len(losers)-1][0]

	w1, l1, _ := match.Play(winnerchamp, loserchamp, array, rand)
	w2, _, _ := match.Play(w1, l1, array, rand)

	if w1 == loserchamp && w2 == loserchamp {
		winnerchamp, loserchamp = loserchamp, winnerchamp
	}
	results[num_ranked] = loserchamp
	num_ranked--
	results[num_ranked] = winnerchamp

	if num_ranked != 0 {
		panic("Didn't rank all the players!")
	}

	return results
}

type DoubleEliminationExtendedSeries struct {
}

type ExtendedSeriesMatch struct {
	match BestOfMatch
	gamehistory [][]int
}

var NumExtendedSeries int = 0
var NumCorrections int = 0
var NumUncorrections int = 0
var NumInjustices int = 0

func (m ExtendedSeriesMatch) Play(i, j int, array player.Array, rand *rand.Rand) (int, int, int) {
	var agames, bgames, ngames int

	towin := (m.match.Games + 1) / 2
	if m.gamehistory[i][j] != 0 {
		NumExtendedSeries++
		agames, bgames = towin, m.gamehistory[i][j] - towin
		ngames = m.match.Games * 2 + 1
	} else if m.gamehistory[j][i] != 0 {
		NumExtendedSeries++
		agames, bgames = m.gamehistory[j][i] - towin, towin
		ngames = m.match.Games * 2 + 1
	} else {
		agames, bgames = 0, 0
		ngames = m.match.Games
	}

	iwins, jwins := agames, bgames
	games := iwins + jwins
	for {
		games++
		if winner, _ := array.Play(i,j,rand); winner == i {
			iwins++
		} else {
			jwins++
		}

		if iwins > ngames / 2 {
			// Less(i,j) <=> i is better than j (confusing)
			if array.Less(i,j) && bgames > agames {
				NumCorrections++
			} else if array.Less(j,i) && agames > bgames {
				NumUncorrections++
			} else if array.Less(j,i) && bgames > agames {
				NumInjustices++
			}

			return i, j, games
		} else if jwins > ngames / 2 {
			// Less(i,j) <=> i is better than j (confusing)
			if array.Less(j,i) && agames > bgames {
				NumCorrections++
			} else if array.Less(i,j) && bgames > agames {
				NumUncorrections++
			} else if array.Less(i,j) && agames > bgames {
				NumInjustices++
			}

			return j, i, games
		}
	}
	panic(nil)
}

func (t DoubleEliminationExtendedSeries) Run(array player.Array, match Match) []int {
	rand := rand.New(rand.NewSource(time.Seconds()))
	results := make([]int, array.Len())
	num_ranked := len(results) - 1
	
	// Double elim tourney needs 3 brackets -- winners, losers
	// (losers play each other) and loservswinners (losers play
	// people coming down from winners bracket). Technically,
	// losers and loservswinners are both parts of the losers
	// bracket.
	winners := make([][]int, intlog(array.Len())+1)
	losers := make([][]int, intlog(array.Len())+1)
	losersvswinners := make([][]int, intlog(array.Len())+1)

	// State for each player on who has killed them and their score in that series...
	gamehistory := make([][]int, array.Len())
	for i := range gamehistory {
		gamehistory[i] = make([]int, array.Len())
		for j := range gamehistory[i] {
			gamehistory[i][j] = 0
		}
	}

	winners[0] = initialSeeds(array.Len())

	for r := 1; ; r++ {
		if len(winners[r-1])/2 == 0 {
			break
		}

		// Create next round
		winners[r] = make([]int, len(winners[r-1])/2)
		losers[r] = make([]int, len(winners[r-1])/2)
		losersvswinners[r] = make([]int, len(winners[r-1])/2)

		// Fill in this round
		for i := range winners[r] {
			// Play winners bracket, knocking someone down to losersvswinners bracket
			a := winners[r-1][2 * i]
			b := winners[r-1][2 * i + 1]

			winner, loser, games := match.Play(a, b, array, rand)
			
			winners[r][i], losersvswinners[r][i] = winner, loser
			gamehistory[winner][loser] = games
		}

		// There is no losers bracket in the first round...
		if r == 1 {
			losers[r] = losersvswinners[r]
			continue
		}

		for i := range losers[r] {
			// Otherwise, we play the losers bracket and then put the losers against the winners
			a := losers[r-1][2 * i]
			b := losers[r-1][2 * i + 1]

			// Here is where we need to employ extended series rule
			ext_match := ExtendedSeriesMatch{match.(BestOfMatch), gamehistory}
			losers[r][i], results[num_ranked], _ = ext_match.Play(a, b, array, rand)
			num_ranked--

			a = losers[r][i]
			b = losersvswinners[r][i]

			ext_match = ExtendedSeriesMatch{match.(BestOfMatch), gamehistory}
			losers[r][i], results[num_ranked], _ = ext_match.Play(a, b, array, rand)
			num_ranked--
		}
	}

	// Some special logic for finals -- losers champ has to beat winner twice
	winnerchamp := winners[len(winners)-1][0]
	loserchamp := losers[len(losers)-1][0]

	w1, l1, _ := match.Play(winnerchamp, loserchamp, array, rand)
	w2, _, _ := match.Play(w1, l1, array, rand)

	if w1 == loserchamp && w2 == loserchamp {
		winnerchamp, loserchamp = loserchamp, winnerchamp
	}
	results[num_ranked] = loserchamp
	num_ranked--
	results[num_ranked] = winnerchamp

	if num_ranked != 0 {
		panic("Didn't rank all the players!")
	}

	return results
}

type RoundRobin struct { }

func (t RoundRobin) Run(array player.Array, match Match) []int {
	rand := rand.New(rand.NewSource(time.Seconds()))
	results := make([]int, array.Len())
	wins := make([]int, array.Len())

	// Every player plays every other once, sort players by number of wins

	for a := 0; a < array.Len(); a++ {
		for b := a + 1; b < array.Len(); b++ {
			winner, _, _ := match.Play(a, b, array, rand)
			wins[winner]++
		}
	}

	for i := 0; i < len(results); i++ {
		best := 0

		for j := range wins {
			if wins[j] > wins[best] {
				best = j
			}
		}

		results[i] = best
		wins[best] = -1
	}

	return results
}

type SwissStyle struct { }

func (t SwissStyle) Run(array player.Array, match Match) []int {
	rand := rand.New(rand.NewSource(time.Seconds()))

	num_rounds := intlog(len(array)) * match.(BestOfMatch).Games
	match = BestOfMatch{1}
// 	num_rounds := intlog(len(array))

	// The number of wins for each player
	wins := make([]int, len(array))

	// The list of players each player has played so far (to avoid re-matches)
	history := make([]vector.IntVector, len(array))
	num_matches := 0
	num_repeats := 0

	// Now simulate all of the rounds
	for i := 1; i < num_rounds; i++ {

		// Pick the match-ups for the next round starting with
		// the best player, and selecting the next best that
		// he hasn't already played
		for j := 0; j < len(array) / 2; j++ {

			// Find the next player to rank -- hasn't
			// played this round, and has the most wins of
			// those who haven't played so far
			next := 0
			best := -1
			for k := range wins {
// 				fmt.Println(next, k, history[k], history[k].Len(), wins[k])
				if history[k].Len() < i && wins[k] > best {
					next = k
					best = wins[k]
				}
			}

			// Find his opponent...
			opp := 0
			best = -1
		opploop:
			for k := range wins {
// 				fmt.Println(opp, k, history[k], history[k].Len(), wins[k])
				if k == next {
					continue
				}
				// Guarantee that next hasn't played opp
				for _, p := range history[next] {
					if p == k {
						continue opploop
					}
				}
				if history[k].Len() < i && wins[k] > best {
					opp = k
					best = wins[k]
				}
			}

			if best == -1 {
				// We failed to find an opponent, so let's relax the 'double play' rule
				for k := range wins {
// 					fmt.Println(opp, k, history[k], history[k].Len(), wins[k])
					if k == next {
						continue
					}
					if history[k].Len() < i && wins[k] > best {
						opp = k
						best = wins[k]
					}
				}
				if best == -1 {
					panic("Couldn't find match, even after relaxing rules!")
				}
				num_repeats++
			}

			if next == opp || history[next].Len() >= i || history[opp].Len() >= i {
				panic(fmt.Sprintf("Couldn't find a good match... round : %d, next : %d, opp : %d", i, next, opp))
			}

			// Add to each's match history...
			history[next].Push(opp)
			history[opp].Push(next)

// 			fmt.Printf("Round %d: %d plays %d\n", i, next, opp)

			// Now play the match (a single game) and update the wins
			winner, _, _ := match.Play(next, opp, array, rand)
			wins[winner]++
			num_matches++
		}

// 		fmt.Printf("Round %v standings : %v\n", i, wins)

// 		min := i
// 		max := 0
// 		for _, w := range wins {
// 			if w < min { min = w }
// 			if w > max { max = w }
// 		}
// 		for w := min; w <= max; w++ {
// 			fmt.Printf("%4d | ", w)
// 			for i := range wins {
// 				if wins[i] == w {
// 					fmt.Printf("%d ", i)
// 				}
// 			}
// 			fmt.Println()
// 		}
	}

	results := make([]int, array.Len())

	for i := 0; i < len(results); i++ {
		best := 0

		for j := range wins {
			if wins[j] > wins[best] {
				best = j
			}
		
		}
		results[i] = best
		wins[best] = -1
	}

// 	fmt.Println("Num matches : ", num_matches, ", Num repeats : ", num_repeats)

	return results
}
