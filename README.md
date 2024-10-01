# README

## Running the thing

The go process should be run from the host, not in docker.

```bash
docker-compose up --build --remove-orphans
```

```bash
go run main.go
```

## Observed Error

```txt
2024/09/30 22:25:14 /Users/fd/code/ch-go-poc/main.go:38 code: 62, message: Syntax error: failed at position 323 ('.'): .`timestamp` DateTime64(9),`events`.`name` LowCardinality(String),`events`.`attributes` Map(LowCardinality(String), String)  ) ENGINE=MergeTree() ORDER BY tuple. Expected one of: NULL, NOT, DEFAULT, MATERIALIZED, EPHEMERAL, ALIAS, AUTO_INCREMENT, PRIMARY KEY, data type, identifier
```

## Generated SQL

```sql
CREATE TABLE `otel_log_records`(
    `id` STRING,
    `created_at` DateTime64(3),
    `updated_at` DateTime64(3),
    `timestamp` DateTime64(9) CODEC(Delta(8), ZSTD(1)),
    `exemplars` Nested(
        filtered_attributes Map(LowCardinality(STRING), STRING),
        time_unix DateTime64(9),
        value Float64,
        span_id STRING,
        trace_id STRING
    ) CODEC(ZSTD(1)),
    `events`.`timestamp` DateTime64(9),
    `events`.`name` LowCardinality(STRING),
    `events`.`attributes` Map(LowCardinality(STRING), STRING)
) ENGINE = MergeTree()
ORDER BY
    tuple()
```

### Working Version

```sql
CREATE TABLE `otel_log_records`(
    `id` STRING,
    `created_at` DateTime64(3),
    `updated_at` DateTime64(3),
    `timestamp` DateTime64(9) CODEC(Delta(8), ZSTD(1)),
    `exemplars` Nested(
        filtered_attributes Map(LowCardinality(STRING), STRING),
        time_unix DateTime64(9),
        value Float64,
        span_id STRING,
        trace_id STRING
    ) CODEC(ZSTD(1)),
    `events.timestamp` DateTime64(9),
    `events.name` LowCardinality(STRING),
    `events.attributes` Map(LowCardinality(STRING), STRING)
) ENGINE = MergeTree()
ORDER BY
    tuple()
```
