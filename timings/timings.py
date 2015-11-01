import requests
from time import time
import ast

def time_get():
    n_iterations = 10000
    start = time()
    for i in xrange(n_iterations):
        if i % 1000 == 0:
            print "Time taken: {} after {} requests".format(time() - start, i)
        r = requests.get("http://localhost:8000/get_recs", params={"UserId": "3", "NRecs": 3})
        result = ast.literal_eval(r.content)
    print "Average time per request is: {}".format((time() - start) / n_iterations)

if __name__ == "__main__":
    time_get()
