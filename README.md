Possible usage using Python's `requests`.

```
In [91]: import requests


In [92]: r = requests.get("http://localhost:8000/get_recs", params={"UserId": "3", "NRecs": 3})


In [93]: r.content
Out[93]: '["10","11","12"]'

In [94]: r = requests.put("http://localhost:8000/update", data='{"3": ["11"]}')
In [95]: r = requests.get("http://localhost:8000/get_recs", params={"UserId": "3", "NRecs": 3})

In [96]: r.content
Out[96]: '["14","10","11"]'
```