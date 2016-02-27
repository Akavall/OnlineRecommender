import json
import numpy as np
import random as rn 

N_ITEMS = 500

ITEM_IDS = map(str, xrange(1000, 1000+N_ITEMS))

def make_item_id_to_col():
    item_id_to_col = {item_id: col for col, item_id in enumerate(ITEM_IDS)}
    with open("item_id_to_col.json", "w") as f:
        json.dump(item_id_to_col, f)

def make_item_id_to_name():
    item_id_to_name = {item_id: "name_{}".format(item_id) for item_id in ITEM_IDS}
    with open("item_id_to_name.json", "w") as f:
        json.dump(item_id_to_name, f)
    

def make_sim_matrix():
    temp = np.random.uniform(0,1, (N_ITEMS, N_ITEMS))
    # creates symmetrical matrix with diagonal 0 of zeros 
    sim = np.tril(temp, -1) + np.tril(temp, -1).T
    with open("similarity.json", "w") as f:
        json.dump(sim.tolist(), f)

def make_user_id_to_actions():
    user_to_actions = {}
    for i in xrange(1000000, 1000100):
        user_to_actions[str(i)] = rn.sample(ITEM_IDS, 50)
    with open("user_id_to_actions.json", "w") as f:
        json.dump(user_to_actions, f)

if __name__ == "__main__":
    make_item_id_to_col()
    make_item_id_to_name()
    make_sim_matrix()
    make_user_id_to_actions()

