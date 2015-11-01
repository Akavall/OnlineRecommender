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

	// started parsing input 
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

	for _, action := range user_actions {
		for i := 0; i < len(recommender.item_id_to_col); i++ {
			scores[i] += recommender.similarity[recommender.item_id_to_col[action]][i]
		}
	}

	item_rankings := argsort.Argsort(scores)
	
	recs := []string {}
	for i := 0; i < n_recs; i++ {
		recs = append(recs, recommender.col_to_item_id[item_rankings[i]])
	}

	recs_json, _ := json.Marshal(recs)

	log.Printf("Recommendations for user: %d", user_id)
	log.Println(recs_json)

	
	fmt.Fprint(w, string(recs_json))
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
	recommender.user_id_to_actions = load_user_id_to_actions(fmt.Sprintf("./%s/user_id_to_actions.json", folder))

	log.Println("Done loading the data, ready...")

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
