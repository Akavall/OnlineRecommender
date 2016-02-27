
package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"log"
	"sync"
	"runtime"
	"flag"

	"github.com/Akavall/OnlineRecommender/utilities"

)

const GREEN = "\033[0;32m"
const YELLOW = "\033[0;33m"
const NO_COLOR = "\033[0m"

// possible usage with requests: 
// In [81]: r = requests.put("http://localhost:8000/update", data='{"3": ["12", "14"]}')

// In [83]: r = requests.get("http://localhost:8000/get_recs", params={"UserId": "3", "NRecs": 3})


func main() {
	runtime.GOMAXPROCS(4)
	folder_string_ptr := flag.String("folder", "", "a string")
	verbose_bool_ptr := flag.Bool("verbose", false, "a bool")

	flag.Parse()

	folder := *folder_string_ptr
	verbose := *verbose_bool_ptr

	log.Printf("Starting to load data from: %s\n", folder)

	recommender := utilities.Recommender{}
	recommender.Mutex = sync.Mutex{}
	recommender.Item_id_to_col = load_item_id_to_col(fmt.Sprintf("./%s/item_id_to_col.json", folder))
	col_to_item_id := map[int]string {}
	for k, v := range recommender.Item_id_to_col {
		col_to_item_id[v] = k
	}
	recommender.Col_to_item_id = col_to_item_id

	recommender.Similarity = load_sim_matrix(fmt.Sprintf("./%s/similarity.json", folder))

	if verbose {
		fmt.Printf("%sSimilarity:%s\n", YELLOW, NO_COLOR)
		fmt.Printf("   ")
		for i := 0; i < len(recommender.Col_to_item_id); i++ {
			fmt.Printf("%s ", recommender.Col_to_item_id[i])
		}
		fmt.Println("")
		for i, row_sim := range recommender.Similarity {
			fmt.Printf("%s %v\n", recommender.Col_to_item_id[i], row_sim)
		}
	}

	recommender.User_id_to_actions = load_user_id_to_actions(fmt.Sprintf("./%s/user_id_to_actions.json", folder))

	recommender.User_id_to_actions_map = map[string]map[string]int {}
	for user_id, actions := range recommender.User_id_to_actions {
		_, ok := recommender.User_id_to_actions_map[user_id]
		if !ok {
			recommender.User_id_to_actions_map[user_id] = map[string]int {}
		}
		for _, item_id := range actions {
			recommender.User_id_to_actions_map[user_id][item_id] += 1
		}
	}

	recommender.Item_id_to_name = load_item_id_to_name(fmt.Sprintf("./%s/item_id_to_name.json", folder))

	recommender.Item_id_to_n_rec = map[string]int {}
	recommender.Item_id_to_n_acted = map[string]int {}

	if verbose {
		log.Printf("%sitem_id => col%s\n", YELLOW, NO_COLOR)
		for k, v := range recommender.Item_id_to_col {
			log.Printf("%s : %d\n", k, v)
		}

		log.Printf("%sitem_id => name%s\n", YELLOW, NO_COLOR)
		for k, v := range recommender.Item_id_to_name {
			log.Printf("%s : %s\n", k, v)
		}
	}

	log.Printf("%sDone loading the data, ready...%s\n", GREEN, NO_COLOR)

	http.HandleFunc("/update", recommender.Update_user_id_to_actions)
	http.HandleFunc("/get_recs", recommender.Make_recs)
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
