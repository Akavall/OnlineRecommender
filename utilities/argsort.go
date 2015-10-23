package argsort

import (
	"sort"
)

type TwoSlices struct {
	indices  []int
	scores  []float64
}

type SortByOther TwoSlices

func (sbo SortByOther) Len() int {
	return len(sbo.indices)
}

func (sbo SortByOther) Swap(i, j int) {
	sbo.indices[i], sbo.indices[j] = sbo.indices[j], sbo.indices[i]
	sbo.scores[i], sbo.scores[j] = sbo.scores[j], sbo.scores[i]
}

func (sbo SortByOther) Less(i, j int) bool {
	return sbo.scores[i] > sbo.scores[j] 
}

func Argsort(scores []float64) []int {
	indices := make([]int, len(scores))
	for i := 0; i < len(scores); i++ {
		indices[i] = i
	}
	my_two_slices := TwoSlices{indices: indices, scores: scores}
	sort.Sort(SortByOther(my_two_slices))
	return my_two_slices.indices

}
