package main

import (
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"encoding/json"
	"log"
	"strconv"

	"github.com/Akavall/OnlineRecommender/utilities"

)

// possible usage with requests: 
// In [81]: r = requests.put("http://localhost:8000/update", data='{"3": ["12", "14"]}')

// In [83]: r = requests.get("http://localhost:8000/get_recs", params={"UserId": "3", "NRecs": 3})


type Recommender struct {
	item_id_to_col map[string]int
	col_to_item_id map[int]string 
	similarity [][]float64
	user_id_to_actions map[string][]string 
}

func (recommender *Recommender) update_user_id_to_actions (w http.ResponseWriter, r *http.Request) {


	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	log.Println("Printing")
	log.Println(body)

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	new_data := map[string][]string {}

	err = json.Unmarshal(body, &new_data)

	if err != nil {
		panic(err)
	}

	for k, v := range new_data {
		_, ok := recommender.user_id_to_actions[k]
		if ok == true {
			recommender.user_id_to_actions[k] = append(recommender.user_id_to_actions[k], v...)
		} else {
			recommender.user_id_to_actions[k] = v
		}
	}
}

func (recommender *Recommender) make_recs(w http.ResponseWriter, r *http.Request) {

	// parsing input 
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	log.Println(r.Form)

	user_id := r.Form["UserId"][0]
	n_recs_str := r.Form["NRecs"][0]

	n_recs, err := strconv.Atoi(n_recs_str)

	if err != nil {
		panic(err)
	}

	// done parsing input
	scores := make([]float64, len(recommender.item_id_to_col))
	user_actions, _ := recommender.user_id_to_actions[user_id]

	log.Println("User Actions")
	log.Println(user_actions)

	for _, action := range user_actions {
		for i := 0; i < len(recommender.item_id_to_col); i++ {
			scores[i] += recommender.similarity[recommender.item_id_to_col[action]][i]
		}
	}

	fmt.Println(scores)
	item_rankings := argsort.Argsort(scores)
	fmt.Println(item_rankings)

	recs := []string {}
	for i := 0; i < n_recs; i++ {
		recs = append(recs, recommender.col_to_item_id[item_rankings[i]])
	}

	recs_json, _ := json.Marshal(recs)

	log.Println(recs_json)

	
	fmt.Fprint(w, string(recs_json))

	
}

func main() {
	fmt.Println("Hello")

	const n_items = 5;

	recommender := Recommender{}
	recommender.item_id_to_col = map[string]int {"10": 0, "11": 1, "12": 2, "13": 3, "14": 4}
	col_to_item_id := map[int]string {}
	for k, v := range recommender.item_id_to_col {
		col_to_item_id[v] = k
	}
	recommender.col_to_item_id = col_to_item_id

	recommender.similarity = make_sim_matrix(n_items)
	recommender.user_id_to_actions = make_user_id_to_actions()

	http.HandleFunc("/update", recommender.update_user_id_to_actions)
	http.HandleFunc("/get_recs", recommender.make_recs)
	http.ListenAndServe("localhost:8000", nil)
	
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

// This should be in a different file
// type TwoSlices struct {
// 	indices  []int
// 	scores  []float64
// }

// type SortByOther TwoSlices

// func (sbo SortByOther) Len() int {
// 	return len(sbo.indices)
// }

// func (sbo SortByOther) Swap(i, j int) {
// 	sbo.indices[i], sbo.indices[j] = sbo.indices[j], sbo.indices[i]
// 	sbo.scores[i], sbo.scores[j] = sbo.scores[j], sbo.scores[i]
// }

// func (sbo SortByOther) Less(i, j int) bool {
// 	return sbo.scores[i] > sbo.scores[j] 
// }

// func argsort(scores []float64) []int {
// 	indices := make([]int, len(scores))
// 	for i := 0; i < len(scores); i++ {
// 		indices[i] = i
// 	}
// 	my_two_slices := TwoSlices{indices: indices, scores: scores}
// 	sort.Sort(SortByOther(my_two_slices))
// 	return my_two_slices.indices

// }

