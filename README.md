Possible usage using Python's `requests`.

```
In [91]: import requests

# Generating recs for user with no actions
In [92]: r = requests.get("http://localhost:8000/get_recs", params={"UserId": "3", "NRecs": 3})


In [93]: r.content
Out[93]: '["10","11","12"]'

# Adding info that user "3" had action on item "11"
# Item "11" is similar to item "14" (by similarity matrix in source code)
In [94]: r = requests.put("http://localhost:8000/update", data='{"3": ["11"]}')

# Generating recs with updated data
In [95]: r = requests.get("http://localhost:8000/get_recs", params={"UserId": "3", "NRecs": 3})

In [96]: r.content
Out[96]: '["14","10","11"]'
```

