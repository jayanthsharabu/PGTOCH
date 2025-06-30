# PGTOCH

ETL tool to stream data from PostgreSQL to ClickHouse

## Features

- Batching and parallel processing
- SQL injection protection
- Schema conversion with type inference
- Retry with exponential backoff
- Structured logging with Zap
- CLI-based interface
- Docker Compose setup
- YAML configuration support
- CDC polling
- UUID support
- CSV export

## Installation

```bash
git clone https://github.com/jayanthsharabu/pgtocg.git 
go build -o pgtoch
./pgtoch
```

## Docker Setup

```bash
docker compose up -d
```

This will set up PostgreSQL and ClickHouse on their respective ports.

## Usage

### Test Connections

```bash
./pgtoch connect --pg-url "postgres://chugger:secret@localhost:5432/pgtoch-db" \
                 --ch-url "localhost:9000"
```

### Ingest Data

```bash
./pgtoch ingest --pg-url "postgres://chugger:secret@localhost:5432/pgtoch-db" \
                --ch-url "localhost:9000" \
                --table "users" \
                --limit 10000 \
                --batch-size 1000
```

### Full Ingest Command

```bash
./pgtoch ingest --pg-url <postgres-connection-string> \
                --ch-url <clickhouse-connection-string> \
                --table <table-name> \
                [--limit <max-rows>] \
                [--batch-size <rows-per-batch>] \
                [--config <path-to-config-file>] \
                [--poll] \
                [--poll-delta <delta-column>] \
                [--poll-interval <seconds>]
```

### Generate Sample Configuration

```bash
./pgtoch sample-config
```

### Export Data

```bash
./pgtoch export --ch-url <clickhouse-connection-string> \
                --table <table-name> \
                --format csv \
                --out <output-directory>
```

## Architecture

### ETL Process

- **Extract**: Retrieve data from PostgreSQL with optional row limits
- **Transform**: Convert table schemas to ClickHouse-compatible format with data type mapping
- **Load**: Create target tables and load data in batches with retry logic

### Components

- **cmd/**: CLI interface definitions using Cobra, connection testing, ingestion commands, polling, root command
- **internal/db/**: PostgreSQL and ClickHouse connection handling
- **internal/etl/**: Core ETL functionality with retry mechanisms
- **internal/config/**: YAML configuration loading and parsing
- **internal/poller/**: CDC polling functionality
- **internal/log/**: Structured logging with Zap

## yet to implement

- [ ] Metrics 
- [ ] Data validation
- [ ] - [ ] Parquet format support
