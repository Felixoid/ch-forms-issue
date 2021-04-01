# Issue reproduction

```
# run docker in one console
docker run --rm --net=host --name=clickhouse -v $(pwd)/schemes:/docker-entrypoint-initdb.d -v $(pwd)/configs/:/etc/clickhouse-server/config.d yandex/clickhouse-server:21.3.4

# Make requests from another
# works
go run .
# works
go run . -url "https://default@localhost:8443/" -lines 1
# works
go run . -url "https://default@localhost:8443/" -lines 1 -agent
# works
go run . -url "https://localhost:8443/?cancel_http_readonly_queries_on_client_close=1" -lines 1 -agent
# broken
go run . -url "https://default@localhost:8443/?cancel_http_readonly_queries_on_client_close=1" -lines 1 -agent
```
