package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println("Hello")

	const n_items = 5;
	
	item_id_to_col := map[string]int {"10": 0, "11": 1, "12": 2, "13": 3, "14": 4}
	col_to_item_id := map[int]string {}
	for k, v := range item_id_to_col {
		col_to_item_id[v] = k
	}
	similarity := make_sim_matrix(n_items)
	user_id_to_actions := make_user_id_to_actions()

	fmt.Println(item_id_to_col)
	fmt.Println(col_to_item_id)
	fmt.Println(similarity)
	fmt.Println(user_id_to_actions)

	greg_recs := make_recs(similarity, item_id_to_col, col_to_item_id, user_id_to_actions["greg"], 3)
	fmt.Println("Gregs recs: ", greg_recs)

	emily_recs := make_recs(similarity, item_id_to_col, col_to_item_id, user_id_to_actions["emily"], 3)
	fmt.Println("Emily recs: ", emily_recs)
	
}

func make_recs(similarity [][]float64, item_id_to_col map[string]int, col_to_item_id map[int]string, user_actions []string, n_recs int) []string {
	scores := make([]float64, len(item_id_to_col))
	for _, action := range user_actions {
		for i := 0; i < len(item_id_to_col); i++ {
			scores[i] += similarity[item_id_to_col[action]][i]
		}
	}

	fmt.Println(scores)
	item_rankings := argsort(scores)
	fmt.Println(item_rankings)

	recs := []string {}
	for i := 0; i < n_recs; i++ {
		recs = append(recs, col_to_item_id[item_rankings[i]])
	}

	return recs  
}

func make_sim_matrix(n_items int) [][]float64 {
	similarity := make([][]float64, n_items)
	for i := 0; i < n_items; i++ {
		similarity[i] = make([]float64, 5)
	}
	
	similarity[2][3] = 0.5
	similarity[3][2] = 0.5
	similarity[1][4] = 0.3
	similarity[4][1] = 0.3

	return similarity
}

func make_user_id_to_actions() map[string][]string {
	user_id_to_actions := make(map[string][]string)
	user_id_to_actions["greg"] = []string {"12",}
	user_id_to_actions["emily"] = []string {"12","14",}
	return user_id_to_actions
}

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

func argsort(scores []float64) []int {
	indices := make([]int, len(scores))
	for i := 0; i < len(scores); i++ {
		indices[i] = i
	}
	my_two_slices := TwoSlices{indices: indices, scores: scores}
	sort.Sort(SortByOther(my_two_slices))
	return my_two_slices.indices

}
