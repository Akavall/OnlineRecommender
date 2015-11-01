import requests
from time import time
import ast
import sys

def time_get(n_iterations, n_recs):
    # This is not testing under load,
    # since all calls happen in sequence 
    start = time()
    for i in xrange(n_iterations):
        if i % 1000 == 0:
            print "Time taken: {} after {} requests".format(time() - start, i)
        r = requests.get("http://localhost:8000/get_recs", params={"UserId": "1000010", "NRecs": n_recs})
        result = ast.literal_eval(r.content)
    print "Average time per request is: {}".format((time() - start) / n_iterations)

if __name__ == "__main__":
    n_iterations = int(sys.argv[1])
    n_recs = (sys.argv[2])
    time_get(n_iterations, n_recs)
