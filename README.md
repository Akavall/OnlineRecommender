To start a run

```
go run online_recommender.go -folder=./sample_data -verbose=true
```

Possible usage using Python's `requests`.

```
In [4]: r = requests.get("http://localhost:8000/get
_recs", params={"UserId": "Bob", "NRecs": 3})      

In [5]: r.content
Out[5]: '{"Recs":["10","11","12"],"RecToItemContrib":{"10":{},"11":{},"12":{}},"RecToHotnessContrib":{"10":0,"11":0,"12":0}}'

In [6]: r = requests.put("http://localhost:8000/update", data='{"Bob": ["11"]}')

In [7]: r = requests.get("http://localhost:8000/get_recs", params={"UserId": "Bob", "NRecs": 3})

In [8]: r.content
Out[8]: '{"Recs":["14","10","11"],"RecToItemContrib":{"10":{},"11":{},"14":{"11":0.3}},"RecToHotnessContrib":{"10":0,"11":0,"14":0}}'

```

To run unit test, go to /utilities:

```
go test
```

Very naive performance test:

```
wrk -t12 -c400 -d30s http://localhost:8000/get_recs?"UserId=Bon&NRecs=3"
```

