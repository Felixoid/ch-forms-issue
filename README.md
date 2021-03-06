# Issue reproduction

This repo is related to ClickHouse issues:

- https://github.com/ClickHouse/ClickHouse/issues/22433
- https://github.com/ClickHouse/ClickHouse/issues/21953
- https://github.com/ClickHouse/ClickHouse/issues/21953
- https://github.com/lomik/graphite-clickhouse/issues/135

```
# run docker in one console
docker run --rm --net=host --name=clickhouse -v $(pwd)/schemes:/docker-entrypoint-initdb.d -v $(pwd)/configs/:/etc/clickhouse-server/config.d yandex/clickhouse-server:21.3.4

# Make requests from another
# works
go run .
# All following with TLS/SSL if another is not mentioned
# works: auth
go run . -url "https://default@localhost:8443/" -lines 1
# works: auth+agent
go run . -url "https://default@localhost:8443/" -lines 1 -agent
# works: agent+cancel_http_readonly_queries_on_client_close
go run . -url "https://localhost:8443/?cancel_http_readonly_queries_on_client_close=1" -lines 1 -agent
# works: agent+cancel_http_readonly_queries_on_client_close+http
go run . -url "http://default@localhost:8123/?cancel_http_readonly_queries_on_client_close=1" -lines 1 -agent
# broken: agent + cancel_http_readonly_queries_on_client_close + HTTPS + auth
go run . -url "https://default@localhost:8443/?cancel_http_readonly_queries_on_client_close=1" -lines 1 -agent
```

It works with `yandex/clickhouse-server:20.8` image.
