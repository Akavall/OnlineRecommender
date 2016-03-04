package utilities

func make_hotness_transparency(recommender *Recommender, item_id_to_hotness map[string]float64, recs []string, user_id string) map[string]float64 {
	rec_to_hotness_contrib := map[string]float64{}
	for _, item_id := range recs {
		_, in_map := recommender.User_id_to_actions_map[user_id][item_id]
		if in_map == false {
			rec_to_hotness_contrib[item_id] = item_id_to_hotness[item_id]
		} else {
			rec_to_hotness_contrib[item_id] = 0
		}
	}
	return rec_to_hotness_contrib
}

func make_item_contrib_transparecy(recommender *Recommender, recs []string, user_id string) map[string]map[string]float64 {
	rec_to_item_contrib := map[string]map[string]float64{}
	for _, item_id := range recs {
		item_id_to_contrib := map[string]float64{}
		for col, similarity := range recommender.Similarity[recommender.Item_id_to_col[item_id]] {
			lift_item_id, ok := recommender.Col_to_item_id[col]
			times_acted, has_acted := recommender.User_id_to_actions_map[user_id][lift_item_id]

			if has_acted && ok && similarity > 0.0 {
				item_id_to_contrib[lift_item_id] = similarity * float64(times_acted)
			}
		}
		rec_to_item_contrib[item_id] = item_id_to_contrib
	}
	return rec_to_item_contrib
}
