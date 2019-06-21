# kraken

```
curl -s "http://localhost:9092/users" | jq '.'

curl -s "http://localhost:9091/scenes" | jq '.'

curl -s "http://localhost:9090/" | jq '.'

curl -s "http://localhost:9090/users" | jq '.'

curl -s "http://localhost:9090/scenes" | jq '.'

curl -s -u kraken:kraken "http://localhost:9090/users" | jq '.'

curl -s -u kraken:kraken "http://localhost:9090/scenes" | jq '.'

```
