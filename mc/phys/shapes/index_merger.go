package shapes

import "math"

// IndexMerger mirrors net.minecraft.world.phys.shapes.IndexMerger. Java's
// fastutil DoubleList is represented as []float64.
type IndexMerger interface {
	GetList() []float64
	ForMergedIndexes(consumer func(firstIndex int, secondIndex int, resultIndex int) bool) bool
	Size() int
}

// IdenticalMerger mirrors IndexMerger.IdenticalMerger.
type IdenticalMerger struct {
	coords []float64
}

// NewIdenticalMerger mirrors IdenticalMerger(DoubleList).
func NewIdenticalMerger(coords []float64) (ret IdenticalMerger) {
	return IdenticalMerger{coords: coords}
}

// ForMergedIndexes mirrors IdenticalMerger.forMergedIndexes.
func (m IdenticalMerger) ForMergedIndexes(consumer func(int, int, int) bool) (ret bool) {
	size := len(m.coords) - 1
	for i := 0; i < size; i++ {
		if !consumer(i, i, i) {
			return false
		}
	}
	return true
}

// Size mirrors IdenticalMerger.size().
func (m IdenticalMerger) Size() (ret int) { return len(m.coords) }

// GetList mirrors IdenticalMerger.getList().
func (m IdenticalMerger) GetList() (ret []float64) { return m.coords }

// DiscreteCubeMerger mirrors IndexMerger.DiscreteCubeMerger.
type DiscreteCubeMerger struct {
	result    []float64
	firstDiv  int
	secondDiv int
}

// NewDiscreteCubeMerger mirrors DiscreteCubeMerger(int, int).
func NewDiscreteCubeMerger(firstSize int, secondSize int) (ret DiscreteCubeMerger) {
	gcd := gcdInt(firstSize, secondSize)
	return DiscreteCubeMerger{
		result:    cubePointRangeList(int(lcmInt(firstSize, secondSize))),
		firstDiv:  firstSize / gcd,
		secondDiv: secondSize / gcd,
	}
}

// ForMergedIndexes mirrors DiscreteCubeMerger.forMergedIndexes.
func (m DiscreteCubeMerger) ForMergedIndexes(consumer func(int, int, int) bool) (ret bool) {
	size := len(m.result) - 1
	for i := 0; i < size; i++ {
		if !consumer(i/m.secondDiv, i/m.firstDiv, i) {
			return false
		}
	}
	return true
}

// Size mirrors DiscreteCubeMerger.size().
func (m DiscreteCubeMerger) Size() (ret int) { return len(m.result) }

// GetList mirrors DiscreteCubeMerger.getList().
func (m DiscreteCubeMerger) GetList() (ret []float64) { return m.result }

// NonOverlappingMerger mirrors IndexMerger.NonOverlappingMerger.
type NonOverlappingMerger struct {
	lower []float64
	upper []float64
	swap  bool
}

// NewNonOverlappingMerger mirrors NonOverlappingMerger(DoubleList, DoubleList, boolean).
func NewNonOverlappingMerger(lower []float64, upper []float64, swap bool) (ret NonOverlappingMerger) {
	return NonOverlappingMerger{lower: lower, upper: upper, swap: swap}
}

// Size mirrors NonOverlappingMerger.size().
func (m NonOverlappingMerger) Size() (ret int) { return len(m.lower) + len(m.upper) }

// ForMergedIndexes mirrors NonOverlappingMerger.forMergedIndexes.
func (m NonOverlappingMerger) ForMergedIndexes(consumer func(int, int, int) bool) (ret bool) {
	if m.swap {
		return m.forNonSwappedIndexes(func(firstIndex int, secondIndex int, resultIndex int) bool {
			return consumer(secondIndex, firstIndex, resultIndex)
		})
	}
	return m.forNonSwappedIndexes(consumer)
}

func (m NonOverlappingMerger) forNonSwappedIndexes(consumer func(int, int, int) bool) (ret bool) {
	lowerSize := len(m.lower)
	for i := 0; i < lowerSize; i++ {
		if !consumer(i, -1, i) {
			return false
		}
	}
	upperSize := len(m.upper) - 1
	for i := 0; i < upperSize; i++ {
		if !consumer(lowerSize-1, i, lowerSize+i) {
			return false
		}
	}
	return true
}

// GetList mirrors NonOverlappingMerger.getList() (AbstractDoubleList view).
func (m NonOverlappingMerger) GetList() (ret []float64) {
	ret = make([]float64, 0, len(m.lower)+len(m.upper))
	ret = append(ret, m.lower...)
	ret = append(ret, m.upper...)
	return
}

// IndirectMerger mirrors IndexMerger.IndirectMerger.
type IndirectMerger struct {
	result        []float64
	firstIndices  []int
	secondIndices []int
	resultLength  int
}

// NewIndirectMerger mirrors IndirectMerger(DoubleList, DoubleList, boolean, boolean).
func NewIndirectMerger(first []float64, second []float64, firstOnlyMatters bool, secondOnlyMatters bool) (ret IndirectMerger) {
	lastValue := math.NaN()
	firstSize := len(first)
	secondSize := len(second)
	capacity := firstSize + secondSize
	result := make([]float64, capacity)
	firstIndices := make([]int, capacity)
	secondIndices := make([]int, capacity)
	canSkipFirst := !firstOnlyMatters
	canSkipSecond := !secondOnlyMatters
	resultIndex := 0
	firstIndex := 0
	secondIndex := 0

	for {
		ranOutOfFirst := firstIndex >= firstSize
		ranOutOfSecond := secondIndex >= secondSize
		if ranOutOfFirst && ranOutOfSecond {
			resultLength := resultIndex
			if resultLength < 1 {
				resultLength = 1
			}
			return IndirectMerger{result, firstIndices, secondIndices, resultLength}
		}

		choseFirst := !ranOutOfFirst && (ranOutOfSecond || first[firstIndex] < second[secondIndex]+1.0e-7)
		if choseFirst {
			firstIndex++
			if canSkipFirst && (secondIndex == 0 || ranOutOfSecond) {
				continue
			}
		} else {
			secondIndex++
			if canSkipSecond && (firstIndex == 0 || ranOutOfFirst) {
				continue
			}
		}

		currentFirstIndex := firstIndex - 1
		currentSecondIndex := secondIndex - 1
		nextValue := first[currentFirstIndex]
		if !choseFirst {
			nextValue = second[currentSecondIndex]
		}
		if !(lastValue >= nextValue-1.0e-7) {
			firstIndices[resultIndex] = currentFirstIndex
			secondIndices[resultIndex] = currentSecondIndex
			result[resultIndex] = nextValue
			resultIndex++
			lastValue = nextValue
		} else {
			firstIndices[resultIndex-1] = currentFirstIndex
			secondIndices[resultIndex-1] = currentSecondIndex
		}
	}
}

// ForMergedIndexes mirrors IndirectMerger.forMergedIndexes.
func (m IndirectMerger) ForMergedIndexes(consumer func(int, int, int) bool) (ret bool) {
	length := m.resultLength - 1
	for i := 0; i < length; i++ {
		if !consumer(m.firstIndices[i], m.secondIndices[i], i) {
			return false
		}
	}
	return true
}

// Size mirrors IndirectMerger.size().
func (m IndirectMerger) Size() (ret int) { return m.resultLength }

// GetList mirrors IndirectMerger.getList().
func (m IndirectMerger) GetList() (ret []float64) {
	if m.resultLength <= 1 {
		return []float64{0.0}
	}
	out := make([]float64, m.resultLength)
	copy(out, m.result[:m.resultLength])
	return out
}

// cubePointRangeList mirrors CubePointRange as a []float64.
func cubePointRangeList(parts int) (ret []float64) {
	if parts <= 0 {
		panic("Need at least 1 part")
	}
	ret = make([]float64, parts+1)
	for i := range ret {
		ret[i] = float64(i) / float64(parts)
	}
	return
}

// isCubePointRange reports whether list is exactly CubePointRange(parts) for
// some parts, mirroring Java's `instanceof CubePointRange` in createIndexMerger.
func isCubePointRange(list []float64) (parts int, ok bool) {
	size := len(list)
	if size < 2 {
		return 0, false
	}
	parts = size - 1
	for i := range list {
		if list[i] != float64(i)/float64(parts) {
			return 0, false
		}
	}
	return parts, true
}

func gcdInt(a int, b int) (ret int) {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func lcmInt(first int, second int) (ret int64) {
	return int64(first) * int64(second/gcdInt(first, second))
}
