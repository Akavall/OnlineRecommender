package utilities

import (
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"encoding/json"
	"log"
	"strconv"
	"sync"
	"math"
)

const GREEN = "\033[0;32m"
const YELLOW = "\033[0;33m"
const NO_COLOR = "\033[0m"

type Recommender struct {
	Mutex sync.Mutex 
	Item_id_to_col map[string]int
	Col_to_item_id map[int]string 
	Similarity [][]float64
	User_id_to_actions map[string][]string 
	User_id_to_actions_map map[string]map[string]int
	Item_id_to_name map[string]string 
	Item_id_to_n_rec map[string]int
	Item_id_to_n_acted map[string]int 
}

type RecResponse struct {
	Recs []string 
	RecToItemContrib map[string]map[string]float64
	RecToHotnessContrib map[string] float64
}

func (recommender *Recommender) Update_user_id_to_actions (w http.ResponseWriter, r *http.Request) {

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

	recommender.Mutex.Lock()
	defer recommender.Mutex.Unlock()

	for k, v := range new_data {
		_, ok := recommender.User_id_to_actions[k]
		if ok == true {
			recommender.User_id_to_actions[k] = append(recommender.User_id_to_actions[k], v...)
		} else {
			recommender.User_id_to_actions[k] = v
		}

		_, ok = recommender.User_id_to_actions_map[k]
		if !ok {
			recommender.User_id_to_actions_map[k] = map[string]int {}
		}

		for _, item_id := range v {
			recommender.Item_id_to_n_acted[item_id] += 1
			recommender.User_id_to_actions_map[k][item_id] += 1
		}
	}
	log.Printf("%suser_id_to_actions updated %s\n", GREEN, NO_COLOR)
}

func (recommender *Recommender) Make_recs(w http.ResponseWriter, r *http.Request) {

	log.Printf("Set up mutexes")
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

	scores := make([]float64, len(recommender.Item_id_to_col))

	recommender.Mutex.Lock()
	user_actions, _ := recommender.User_id_to_actions[user_id]

	log.Printf("User Actions : %v\n", user_actions)

	user_actions_map, _ := recommender.User_id_to_actions_map[user_id]
	log.Printf("User Actions Set: %v\n", user_actions_map)
	recommender.Mutex.Unlock()


	for _, action := range user_actions {
		for i := 0; i < len(recommender.Item_id_to_col); i++ {
			scores[i] += recommender.Similarity[recommender.Item_id_to_col[action]][i]
		}
	}

	log.Printf("%sItem id => times recommended:%s %v", YELLOW, NO_COLOR, recommender.Item_id_to_n_rec)

	log.Printf("%sItem id => times acted:%s %v", YELLOW, NO_COLOR, recommender.Item_id_to_n_acted)

	item_id_to_hotness := calc_hotness_scores(recommender.Item_id_to_n_acted, recommender.Item_id_to_n_rec)

	for item_id, hotness_score := range item_id_to_hotness {
		_, ok := user_actions_map[item_id]
		if !ok {
			// User have not seen this item we can bump it up 
			scores[recommender.Item_id_to_col[item_id]] += hotness_score
		}
	}

	log.Printf("scores: %v\n", scores)
	item_rankings := Argsort(scores)
	log.Printf("item_rankings: %v\n", item_rankings)

	recs := []string {}
	for i := 0; i < n_recs; i++ {
		recs = append(recs, recommender.Col_to_item_id[item_rankings[i]])
	}

	log.Printf("%sRecommendations for user: %s%s: %v", YELLOW, user_id, NO_COLOR, recs)

	for _, item_id := range recs {
		log.Printf("%s", recommender.Item_id_to_name[item_id])
		recommender.Item_id_to_n_rec[item_id] += 1
	}

	rec_to_item_contrib := make_item_contrib_transparecy(recommender, recs, user_id)

	rec_to_item_hotness := make_hotness_transparency(recommender, item_id_to_hotness, recs, user_id)
	
	rec_response := RecResponse {}
	rec_response.Recs = recs
	rec_response.RecToItemContrib = rec_to_item_contrib
	rec_response.RecToHotnessContrib = rec_to_item_hotness 

	rec_response_json, _ := json.Marshal(rec_response)

	fmt.Fprint(w, string(rec_response_json))
	
	log.Printf("%sRequest completed%s", GREEN, NO_COLOR)
	
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
