Update the bearer:

https://www.strava.com/settings/api


```shell

curl -X GET \
https://www.strava.com/api/v3/athlete \
-H 'Authorization: Bearer 47534a526e0e815ab65d19fbeed8043203d144c7' | jq
```

Segment data
```shell
curl -s -X 'GET' \
  'https://www.strava.com/api/v3/segments/17275870' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer 11565fd054b1192df14c0c0a2db8128dac27cda1' | jq

```

Explore segment

```shell
curl -X 'GET' \
  'https://www.strava.com/api/v3/segments/explore?bounds=45,6,45.1,6.1&activity_type=riding&min_cat=5' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer 11565fd054b1192df14c0c0a2db8128dac27cda1' | jq
```


Segment stream
```shell
curl -s -X 'GET' \
  'https://www.strava.com/api/v3/segments/17275870/streams?keys=altitude&key_by_type=true' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer 11565fd054b1192df14c0c0a2db8128dac27cda1' | jq > stream.json

```

Search segment in algolia

```shell
export INDEX="profiles"
export API_KEY="cb6b989a18ef9a3070d5b5a54001a3da"
export APPLICATION_ID="1QMZVCS1V5"
export ALGOLIA_HOST="https://${APPLICATION_ID}.algolia.net"
```

```shell
curl -s -k -X POST --location "${ALGOLIA_HOST}/1/indexes/${INDEX}/query" \
    -H "x-algolia-application-id: ${APPLICATION_ID}" \
    -H "x-algolia-api-key: ${API_KEY}" \
    -d '
    {"attributesToRetrieve": [ "name", "objectID", "_tags"], "filters": "perso"}
    ' | jq .hits 
```