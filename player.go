package player

import (
	"sort"
	"rand"
	"math"
)

type Player struct {
	mean float64
	dev float64
}

const PlayerMeanRange = 2.0
const PlayerDevRange = 1.0

func RandPlayer() *Player {
	return &Player{ PlayerMeanRange * rand.Float64(),
		PlayerDevRange }
}

func (p Player) performance() float64 {
	return p.mean + (rand.Float64() * 2.0 - 1.0) * p.dev
}

func (p Player) Defeats(r Player) bool {
	return p.performance() > r.performance()
}

type Array []*Player

func NewArray(size int) Array {
	a := make([]*Player, size)

	for i := 0; i < size; i++ {
		a[i] = RandPlayer()
	}

	sort.Sort(Array(a))

	return Array(a)
}

func (a Array) Len() int {
	return len(a)
}

// i "less than" j if i defeats j
func (a Array) Less(i, j int) bool {
	return a[i].mean > a[j].mean
}

func (a Array) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Array) Play(i, j int) (int, int) {
	if a[i].Defeats(*a[j]) {
		return i, j
	} else {
		return j, i
	}
	panic(nil)
}

// Distance by first place finisher
func (a Array) DistanceByFirst(ranking []int) float64 {
	if ranking[0] == 0 {
		return 0
	} else {
		return 1
	}
	panic(nil)
}

// Distance by the depth each player reached
func depth(i, size int) float64 {
	i++
	var d uint
	for d = 0; ; d++ {
		if (1 << d) >= i {
			return math.Log2(float64(size)) - float64(d)
		}
	}
	panic(nil)
}

func (a Array) DistanceByDepth(ranking []int) float64 {
	dist := float64(0.0)

	// For each player, measure the depth they 'should have gone'
	// in tournament against where they ended up
	for i, j := range ranking {
		// This means player j was ranked at position i
		ideal_ranking := depth(j, len(ranking))
		actual_ranking := depth(i, len(ranking))
		dist += math.Fabs(ideal_ranking - actual_ranking)
	}

	return dist
}

func (a Array) DistanceByDepthExponential(ranking []int) float64 {
	dist := float64(0.0)

	// For each player, measure the depth they 'should have gone'
	// in tournament against where they ended up
	for i, j := range ranking {
		// This means player j was ranked at position i
		ideal_ranking := depth(j, len(ranking))
		actual_ranking := depth(i, len(ranking))
		dist += math.Exp2(math.Fabs(ideal_ranking - actual_ranking))
	}

	return dist
}