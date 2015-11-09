package main

import (
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"encoding/json"
	"log"
	"strconv"
	// "time"
	"runtime"
	"os"
	"math"
	"sync"

	"github.com/Akavall/OnlineRecommender/utilities"

)

const GREEN = "\033[0;32m"
const YELLOW = "\033[0;33m"
const NO_COLOR = "\033[0m"

// possible usage with requests: 
// In [81]: r = requests.put("http://localhost:8000/update", data='{"3": ["12", "14"]}')

// In [83]: r = requests.get("http://localhost:8000/get_recs", params={"UserId": "3", "NRecs": 3})


type Recommender struct {
	*sync.Mutex 
	item_id_to_col map[string]int
	col_to_item_id map[int]string 
	similarity [][]float64
	user_id_to_actions map[string][]string 
	user_id_to_actions_map map[string]map[string]int
	item_id_to_name map[string]string 
	item_id_to_n_rec map[string]int
	item_id_to_n_acted map[string]int 
}

func (recommender *Recommender) update_user_id_to_actions (w http.ResponseWriter, r *http.Request) {
	recommender.Lock()
	defer recommender.Unlock()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}

	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	new_data := map[string][]string {}

	err = json.Unmarshal(body, &new_data)

	log.Printf("%sPut info received:%s %v\n", YELLOW, NO_COLOR, new_data)

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

		_, ok = recommender.user_id_to_actions_map[k]
		if !ok {
			recommender.user_id_to_actions_map[k] = map[string]int {}
		}

		for _, item_id := range v {
			recommender.item_id_to_n_acted[item_id] += 1
			recommender.user_id_to_actions_map[k][item_id] += 1
		}
	}
	log.Printf("%suser_id_to_actions updated %s\n", GREEN, NO_COLOR)
}

func (recommender *Recommender) make_recs(w http.ResponseWriter, r *http.Request) {
	recommender.Lock()
	recommender.Unlock()
	// started parsing input 
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	log.Printf("%sRequest Parameters:%s %v", YELLOW, NO_COLOR, r.Form)

	user_id := r.Form["UserId"][0]
	n_recs_str := r.Form["NRecs"][0]

	n_recs, err := strconv.Atoi(n_recs_str)

	if err != nil {
		panic(err)
	}

	// done parsing input

	scores := make([]float64, len(recommender.item_id_to_col))
	user_actions, _ := recommender.user_id_to_actions[user_id]

	log.Printf("User Actions : %v\n", user_actions)

	user_actions_map, _ := recommender.user_id_to_actions_map[user_id]
	log.Printf("User Actions Set: %v\n", user_actions_map)

	for _, action := range user_actions {
		for i := 0; i < len(recommender.item_id_to_col); i++ {
			scores[i] += recommender.similarity[recommender.item_id_to_col[action]][i]
		}
	}

	log.Printf("%sItem id => times recommended:%s %v", YELLOW, NO_COLOR, recommender.item_id_to_n_rec)

	log.Printf("%sItem id => times acted:%s %v", YELLOW, NO_COLOR, recommender.item_id_to_n_acted)

	// Rebuilding this is probably not very efficient 
	item_id_to_hotness := calc_hotness_scores(recommender.item_id_to_n_acted, recommender.item_id_to_n_rec)

	for item_id, hotness_score := range item_id_to_hotness {
		_, ok := recommender.user_id_to_actions_map[user_id][item_id]
		if !ok {
			// User have not seen this item we can bump it up 
			scores[recommender.item_id_to_col[item_id]] += hotness_score
		}
	}

	log.Printf("scores: %v\n", scores)
	item_rankings := argsort.Argsort(scores)
	log.Printf("item_rankings: %v\n", item_rankings)

	recs := []string {}
	for i := 0; i < n_recs; i++ {
		recs = append(recs, recommender.col_to_item_id[item_rankings[i]])
	}

	recs_json, _ := json.Marshal(recs)

	log.Printf("%sRecommendations for user: %s%s: %v", YELLOW, user_id, NO_COLOR, recs)

	for _, item_id := range recs {
		log.Printf("%s", recommender.item_id_to_name[item_id])
		recommender.item_id_to_n_rec[item_id] += 1
	}

	fmt.Fprint(w, string(recs_json))

	log.Printf("%sRequest completed%s", GREEN, NO_COLOR)
	
}

func main() {
	runtime.GOMAXPROCS(4)
	folder := os.Args[1]

	log.Printf("Starting to load data from: %s\n", folder)

	recommender := Recommender{}
	recommender.item_id_to_col = load_item_id_to_col(fmt.Sprintf("./%s/item_id_to_col.json", folder))
	col_to_item_id := map[int]string {}
	for k, v := range recommender.item_id_to_col {
		col_to_item_id[v] = k
	}
	recommender.col_to_item_id = col_to_item_id

	recommender.similarity = load_sim_matrix(fmt.Sprintf("./%s/similarity.json", folder))

	fmt.Printf("%sSimilarity:%s\n", YELLOW, NO_COLOR)
	fmt.Printf("   ")
	for i := 0; i < len(recommender.col_to_item_id); i++ {
		fmt.Printf("%s ", recommender.col_to_item_id[i])
	}
	fmt.Println("")
	for i, row_sim := range recommender.similarity {
		fmt.Printf("%s %v\n", recommender.col_to_item_id[i], row_sim)
	}

	recommender.user_id_to_actions = load_user_id_to_actions(fmt.Sprintf("./%s/user_id_to_actions.json", folder))

	recommender.user_id_to_actions_map = map[string]map[string]int {}
	for user_id, actions := range recommender.user_id_to_actions {
		_, ok := recommender.user_id_to_actions_map[user_id]
		if !ok {
			recommender.user_id_to_actions_map[user_id] = map[string]int {}
		}
		for _, item_id := range actions {
			recommender.user_id_to_actions_map[user_id][item_id] += 1
		}
	}

	recommender.item_id_to_name = load_item_id_to_name(fmt.Sprintf("./%s/item_id_to_name.json", folder))

	recommender.item_id_to_n_rec = map[string]int {}
	recommender.item_id_to_n_acted = map[string]int {}

	log.Printf("%sitem_id => name%s\n", YELLOW, NO_COLOR)
	for k, v := range recommender.item_id_to_name {
		log.Printf("%s : %s\n", k, v)
	}

	log.Printf("%sDone loading the data, ready...%s\n", GREEN, NO_COLOR)

	http.HandleFunc("/update", recommender.update_user_id_to_actions)
	http.HandleFunc("/get_recs", recommender.make_recs)
	http.ListenAndServe("localhost:8000", nil)
}

func load_item_id_to_col(file_address string) map[string]int {
	item_id_to_col := map[string]int {}

	f, err := ioutil.ReadFile(file_address)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(f, &item_id_to_col)
	return item_id_to_col
}

func load_sim_matrix(file_address string) [][]float64 {
	sim_matrix := [][]float64 {}

	f, err := ioutil.ReadFile(file_address)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(f, &sim_matrix)
	return sim_matrix 
}

func load_user_id_to_actions(file_address string) map[string][]string {

	user_id_to_actions := map[string][]string {}

	f, err := ioutil.ReadFile(file_address)
	if err != nil {
		panic(err)
	}
	
	json.Unmarshal(f, &user_id_to_actions)
	return user_id_to_actions
}

func load_item_id_to_name(file_address string) map[string]string {

	item_id_to_name := map[string]string {}

	f, err := ioutil.ReadFile(file_address)
	if err != nil {
		panic(err)
	}
	
	json.Unmarshal(f, &item_id_to_name)
	return item_id_to_name
}

func calc_hotness_scores(item_id_to_n_acted, item_id_to_n_rec map[string]int) map[string]float64 {
	item_id_to_hotness := map[string]float64 {}
	for item_id, n_acted := range item_id_to_n_acted {
		n_rec, ok := item_id_to_n_rec[item_id]
		if !ok || n_rec == 0 {
			item_id_to_hotness[item_id] = 1.0;
		} else {
			score_1 := float64(n_acted) / float64(n_rec)
			score := math.Min(score_1, 1.0)
			item_id_to_hotness[item_id] = score
		}		
	}
	return item_id_to_hotness
}

