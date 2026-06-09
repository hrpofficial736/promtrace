# Portfolio Projects — Final 5

> **Stack:** Go (primary), Python (AI/ML where noted)
> **Format:** CLI-first, server-side, real problems only
> **Criteria:** Job-relevant · Real unsolved problem · Defensible in any interview · Math/algorithmic depth · Actually usable

---

## Build Order

| # | Project | Why This Order |
|---|---------|---------------|
| 1 | **promtrace** | Fastest to ship (2–3 weeks), immediately useful, hottest market signal right now |
| 2 | **Vectraflow** | Deepest technical project, builds on fresh momentum |
| 3 | **Flowpipe** | Postgres WAL depth — directly feeds ghostdb knowledge |
| 4 | **ghostdb** | Builds on Flowpipe's WAL understanding, most original idea |
| 5 | **delhibus** | Personal project, build whenever, use daily |

---

## What These Five Tell Together

- **AI infrastructure:** promtrace (LLM observability) + Vectraflow (retrieval infrastructure)
- **Backend systems depth:** Flowpipe (Postgres internals, streaming) + ghostdb (CoW branching)
- **Civic tech + personal utility:** delhibus (real open data, real daily use)
- **No overlap.** Five different hiring categories. Every project has a one-sentence answer to "what problem does it solve."

---
---

## Project 1 — promtrace

**Tagline:** A transparent LLM call interceptor. See exactly what prompts your app sends, what they cost, and when they silently change — with zero code modifications.

**One-line interview answer:** "LLM apps have zero visibility into what they're actually sending in production. promtrace is a local HTTP proxy that intercepts every call transparently, logs everything, and lets you diff and replay prompts — no SDK changes, no cloud account."

---

### The Problem

You're building an LLM-powered feature. Three weeks in, something starts behaving differently. You don't know if the prompt changed, the model changed, or the response format drifted. You add fmt.Println calls. You stare at logs. You still don't know.

Langfuse and LangSmith require you to instrument your code with their SDKs. Helicone requires rerouting your API keys through their cloud. Every existing solution requires either code changes or sending your prompts to someone else's server.

promtrace does neither. It sits between your app and the LLM API as a transparent local proxy. Your app doesn't know it's there. Every call gets logged. Nothing leaves your machine.

---

### How It's Different

| Tool | What's wrong with it |
|------|----------------------|
| Langfuse | Cloud-based, requires SDK instrumentation in your code |
| Helicone | Cloud proxy, API keys route through their servers |
| LangSmith | Tied entirely to the LangChain ecosystem |
| OpenAI dashboard | Shows billing, not prompt content or latency breakdowns |
| **promtrace** | Local MITM proxy, zero code changes, works with any SDK or language, fully offline |

---

### Exact Features

**1. Transparent process wrapping**
Starts a local proxy server, sets HTTP_PROXY / HTTPS_PROXY env vars, then launches your process as a subprocess. Your app's LLM SDK picks up the proxy automatically. No code changes needed. Works with Python, Node, Go, Ruby — anything that respects env-based proxy settings.

**2. Full call logging**
Every intercepted call is stored in a local SQLite database with: timestamp, model name, system prompt, user prompt, full response, latency in ms, input token count, output token count, estimated cost in USD (using hardcoded per-model pricing tables), HTTP status code, and a SHA-256 hash of the prompt content (used for diffing).

**3. Live TUI dashboard**
`promtrace watch` opens a bubbletea terminal UI showing a live scrolling feed of calls as they happen. Each row shows: time, model, latency, tokens, cost. Selecting a row expands it to show full prompt and response. Updates in real time.

**4. Structural prompt diff engine**
LLM prompts are not plain text diffs — they have system prompts, user turns, assistant turns, and injected variables. The diff engine understands this structure. `promtrace diff <id1> <id2>` shows which turns changed, which stayed the same, and highlights the specific variable injections that differ. Not just line diff — semantic diff of the message array.

**5. Call replay**
`promtrace replay <id>` re-sends a captured call exactly as it was. `--model gpt-4o` swaps the model. `--edit` opens the prompt in your $EDITOR before re-sending. The result of the replay is stored as a new trace linked to the original, so you can compare them.

**6. Session grouping**
Calls made within the same `promtrace wrap` session are grouped by session ID. `promtrace sessions` lists all sessions with aggregate stats: total calls, total cost, average latency, total tokens. Useful for comparing "before this prompt change" vs "after."

**7. Cost and token trend reporting**
`promtrace export --format jsonl` dumps all traces. `promtrace stats --last 7d` shows cost and token trends over time. Answers "how much did my LLM calls cost this week and is it growing?"

**8. SSE streaming support**
Correctly handles streaming responses (Server-Sent Events). Buffers chunks transparently, reconstructs the full response for storage, and forwards chunks to the client with no added latency. The user's streaming UI still works normally.

---

### Exact User Flow — Step by Step

**Scenario:** Harry is building a Go app that calls the OpenAI API. He wants to see what's happening.

```
Step 1 — Install
  $ go install github.com/harry/promtrace@latest

Step 2 — First run (one-time TLS setup)
  $ promtrace setup
  → Generates a local CA certificate
  → Installs it as a trusted cert in the system keychain
  → Creates ~/.promtrace/config.toml with defaults
  Output: "Setup complete. promtrace is ready."

Step 3 — Wrap your process
  $ promtrace wrap go run ./cmd/server
  → promtrace starts a proxy on localhost:9117
  → Sets HTTP_PROXY=http://localhost:9117 for the subprocess
  → Launches "go run ./cmd/server" as a child process
  → Your server starts normally — it has no idea it's being proxied
  Output: "[promtrace] proxy listening on :9117"
           "[promtrace] wrapping: go run ./cmd/server"
           <your normal server output follows>

Step 4 — Use your app normally
  → Harry sends a request to his server
  → His server calls OpenAI internally
  → promtrace intercepts the call, logs it, forwards it
  → The OpenAI response comes back normally
  → In another terminal: promtrace watch

Step 5 — Watch the live dashboard
  $ promtrace watch
  ┌─────────────────────────────────────────────────────────────┐
  │ promtrace — live trace feed                    session: a3f │
  ├──────────┬─────────────┬──────────┬────────┬───────────────┤
  │ time     │ model       │ latency  │ tokens │ cost          │
  ├──────────┼─────────────┼──────────┼────────┼───────────────┤
  │ 14:23:01 │ gpt-4o      │ 843ms    │ 1,204  │ $0.0024       │
  │ 14:23:04 │ gpt-4o      │ 1,102ms  │ 987    │ $0.0019       │
  │ 14:23:09 │ gpt-4o-mini │ 312ms    │ 453    │ $0.0001       │
  └──────────┴─────────────┴──────────┴────────┴───────────────┘
  Press ENTER on a row to expand. Press q to quit.

Step 6 — Inspect a specific call
  $ promtrace show a3f-002
  ┌─ Call a3f-002 ──────────────────────────────────────────────┐
  │ Model:    gpt-4o                                            │
  │ Time:     2026-05-09 14:23:04                               │
  │ Latency:  1,102ms                                           │
  │ Tokens:   987 in / 412 out                                  │
  │ Cost:     $0.0019                                           │
  ├─ System Prompt ─────────────────────────────────────────────┤
  │ You are a helpful assistant that answers questions about... │
  ├─ User Message ──────────────────────────────────────────────┤
  │ Summarise the following document: [2,400 chars of text]     │
  ├─ Response ──────────────────────────────────────────────────┤
  │ The document describes a proposal for...                    │
  └─────────────────────────────────────────────────────────────┘

Step 7 — Diff two calls to see what changed
  $ promtrace diff a3f-001 a3f-002
  system prompt:  identical
  user message:
  -  Summarise the following document in 3 bullet points: ...
  +  Summarise the following document: ...
  response length: 312 chars → 891 chars (+579)
  cost:    $0.0024 → $0.0019  (-21%)
  latency: 843ms   → 1,102ms  (+31%)

Step 8 — Replay a call with a different model
  $ promtrace replay a3f-001 --model gpt-4o-mini
  → Re-sends the exact same prompt to gpt-4o-mini
  → Stores result as a3f-001-replay-01
  → Shows side-by-side cost and latency comparison
  gpt-4o:      843ms, $0.0024, 312 chars response
  gpt-4o-mini: 201ms, $0.0001, 287 chars response

Step 9 — Export everything
  $ promtrace export --format jsonl > traces.jsonl
  → Exports all captured traces as newline-delimited JSON
  → Ready for further analysis, sharing, or archiving
```

---

### Algorithms / Technical Depth

- **HTTP MITM proxy** — certificate authority generation (crypto/x509), TLS interception, transparent TCP forwarding via net/http ReverseProxy with a custom Transport that logs before forwarding
- **SSE stream buffering** — intercepts chunked transfer encoding, accumulates `data:` frames, reconstructs full response, re-streams to client
- **Structural prompt diff** — diffs the messages array as a typed structure (role + content pairs), not raw text; uses Myers diff algorithm on the content strings within each role
- **Token estimation** — cl100k_base tokenizer approximation via a lightweight Go port; accurate to within 5% without needing tiktoken
- **Cost table** — hardcoded per-model input/output token pricing, updated at build time from a YAML config file

---

### Tech Stack

```
Language:     Go
Proxy:        net/http ReverseProxy + custom Transport for interception
TLS:          crypto/x509, crypto/tls — generate local CA, sign per-host certs on the fly
Storage:      SQLite via mattn/go-sqlite3
TUI:          bubbletea + lipgloss
CLI:          cobra + viper
Config:       ~/.promtrace/config.toml (TOML)
Diff:         Myers diff algorithm, implemented from scratch on message arrays
```

### Build Time
**2–3 weeks** solo

### Job Market Signal
Every company building LLM-powered products in 2026 is hiring for AI infrastructure. promtrace demonstrates proxy architecture, observability design, LLM system understanding, and local-first engineering — without calling a single paid API. Directly relevant to AI infra roles at YC startups, any HN-posting company building on LLM APIs, and teams where "how do we debug our prompts in production?" is an open question.

---
---

## Project 2 — Vectraflow

**Tagline:** HNSW vector search engine implemented from scratch in Go, with a full RAG pipeline on top. The storage layer every AI product is built on — understood from first principles.

**One-line interview answer:** "I wanted semantic search over private documents without sending embeddings to Pinecone. So I implemented the HNSW algorithm from scratch in Go — probabilistic layer assignment, greedy graph traversal, ef-parameter beam search — and built a RAG pipeline on top of it."

---

### The Problem

Every RAG system needs a vector store. Pinecone costs money and sends your data to their servers. Qdrant requires Docker and a separate process. pgvector requires a full Postgres instance. Chroma is Python-first and heavy.

For local use cases — a developer's personal notes, a company's internal docs, a sensitive codebase — none of these are acceptable. There is no single Go binary that does HNSW-backed semantic search locally with zero cloud dependencies and no Docker.

More importantly: nobody building on these tools understands how they actually work. Vectraflow fixes that by implementing the core algorithm from scratch.

---

### How It's Different

| Tool | What's wrong with it |
|------|----------------------|
| Pinecone | Cloud-only, sends your vectors externally, costs money |
| Qdrant | Requires Docker, separate Rust process, heavy ops |
| pgvector | Requires Postgres, no standalone mode, no RAG pipeline |
| Chroma | Python-first, not embeddable in Go services, heavy deps |
| **Vectraflow** | Single Go binary, HNSW from scratch, offline, private, RAG included |

---

### Exact Features

**1. HNSW index — built from scratch**
The full Hierarchical Navigable Small World algorithm: random level generation on insert, greedy layer-by-layer graph search, ef-parameter beam search at query time, bidirectional edge insertion with degree-cap pruning. Supports concurrent reads with a read-write mutex. Serialises and deserialises the graph to disk in a custom binary format.

**2. Multiple distance metrics**
Cosine similarity, euclidean distance, and dot product. Switchable per-index at creation time. Vectors are optionally L2-normalised on insert (when using cosine), reducing query-time computation to a dot product.

**3. Local embedding inference — no API key**
Uses ONNX Runtime Go bindings to run all-MiniLM-L6-v2 locally. The model is ~22MB. Converts any text string to a 384-dimensional float32 vector on-device, with no internet connection required.

**4. Document ingestion pipeline**
Reads a folder of files (markdown, plain text, code files). Splits each document into chunks using a sliding window with configurable size and overlap stride. Generates an embedding per chunk. Inserts into the HNSW index with metadata (source file, chunk offset, raw text) stored in BoltDB.

**5. Hybrid retrieval — HNSW + BM25**
At query time, runs both vector search (HNSW approximate nearest neighbours) and keyword search (BM25 over the inverted index). Combines scores with a tunable alpha blend: score = alpha x vector_score + (1-alpha) x BM25_score. This handles cases where the query uses different vocabulary than the document but means the same thing (vector wins), and cases where exact keyword matching matters (BM25 wins).

**6. RAG pipeline — retrieval + generation**
`vectraflow rag "your question"` retrieves the top-k most relevant chunks, assembles them into a context window, and sends to a configured LLM (Ollama local model or OpenAI API) to generate a grounded answer. The answer is streamed to stdout. Source chunks are printed below the answer with file path and chunk offset.

**7. gRPC + REST API server mode**
`vectraflow serve` exposes the index as a network service. Protobuf-defined API: Upsert, Query, Delete, Stats. A REST JSON wrapper is auto-generated. Other services can use Vectraflow as their vector database without knowing its internals.

**8. Recall benchmarking**
`vectraflow bench --k 10` computes exact brute-force nearest neighbours for a sample of query vectors, then compares against HNSW results. Reports Recall@K: what fraction of the true top-K neighbours did HNSW find? Shows how the ef parameter affects the speed/recall tradeoff across different settings.

**9. Incremental indexing with file watcher**
`vectraflow watch ./docs` monitors the folder for file changes using fsnotify. New or modified files are automatically re-chunked and re-indexed. Deleted files are removed from the index. The index stays in sync with the folder without manual re-ingestion.

---

### Exact User Flow — Step by Step

**Scenario:** Harry wants semantic search over his personal notes folder (500 markdown files).

```
Step 1 — Install and initialise
  $ vectraflow init ~/notes-index
  → Creates ~/notes-index/ directory
  → Writes index.config: {metric: cosine, M: 16, efConstruction: 200}
  → Downloads all-MiniLM-L6-v2 ONNX model to ~/.vectraflow/models/ (first run, ~22MB)
  Output: "Index initialised at ~/notes-index"
           "Model ready: all-MiniLM-L6-v2 (384 dims)"

Step 2 — Ingest documents
  $ vectraflow ingest ~/notes --index ~/notes-index --chunk 256 --overlap 32
  → Walks ~/notes recursively, finds 500 .md files
  → For each file: reads content, splits into chunks of ~256 tokens with
    32-token overlap, generates embedding per chunk via local ONNX model,
    inserts vector + metadata into HNSW index, stores raw chunk text +
    file path + offset in BoltDB
  → Progress bar shows files processed / total
  Output: "Ingesting 500 files..."
           "[████████████████████] 500/500 files"
           "Indexed 4,821 chunks in 47.3s"
           "Index size: 14.2 MB"

Step 3 — Query semantically
  $ vectraflow query "deployment failures in production" \
      --index ~/notes-index --top 5
  → Embeds the query string locally (no API call)
  → Runs HNSW beam search (ef=64 by default)
  → Retrieves top 5 chunks by cosine similarity
  Output:
    #1  score=0.847  ~/notes/2026-03-incidents.md  [chunk 3, offset 712]
        "...the deployment at 2am caused a cascade because the health check
         timeout was misconfigured..."

    #2  score=0.831  ~/notes/2025-retrospective.md  [chunk 7, offset 1840]
        "...three production failures in Q4 were traced back to..."

    #3  score=0.798  ~/notes/ops-runbook.md  [chunk 12, offset 3201]
        "...when a deployment fails, the first step is to check..."

Step 4 — RAG query (retrieval + generated answer)
  $ vectraflow rag "what caused our worst outage?" --index ~/notes-index
  → Retrieves top 8 relevant chunks
  → Sends to Ollama (llama3 running locally) with context
  → Streams generated answer to stdout
  Output:
    "Based on your notes, the worst outage was on March 14th 2026.
     The root cause was a misconfigured health check timeout (2s) that
     caused the load balancer to mark all pods as unhealthy during a
     rolling deployment. The fix was to increase the timeout to 10s
     and add a pre-deployment health check validation step.

     Sources:
     → ~/notes/2026-03-incidents.md (chunk 3)
     → ~/notes/ops-runbook.md (chunk 12)
     → ~/notes/2026-postmortems.md (chunk 2)"

Step 5 — Watch folder for changes
  $ vectraflow watch ~/notes --index ~/notes-index
  → Starts fsnotify watcher on ~/notes
  Output: "[vectraflow] watching ~/notes for changes..."
  <Harry creates a new note>
  Output: "[vectraflow] detected: new-note.md — indexing..."
           "[vectraflow] indexed 3 chunks from new-note.md"

Step 6 — Benchmark the index quality
  $ vectraflow bench --index ~/notes-index --k 10 --samples 100
  Output:
    ef=32:   Recall@10 = 0.91  avg_latency = 0.8ms
    ef=64:   Recall@10 = 0.97  avg_latency = 1.4ms
    ef=128:  Recall@10 = 0.99  avg_latency = 2.6ms
    Brute:   Recall@10 = 1.00  avg_latency = 18.3ms
    Recommendation: ef=64 gives best speed/recall tradeoff for this index.

Step 7 — Expose as a service
  $ vectraflow serve --index ~/notes-index --port 8080
  Output: "[vectraflow] gRPC server on :8080"
           "[vectraflow] REST server on :8081"

  $ curl localhost:8081/query -d '{"text": "auth bug", "top_k": 3}'
  → Returns JSON array of chunks with scores and metadata
```

---

### Algorithms / Math Involved

```
Layer assignment:   level = floor(−ln(uniform(0,1)) × mL)   where mL = 1 / ln(M)
                    Negative log of uniform → exponential decay
                    Most nodes at layer 0, few at higher layers — skip list geometry

Greedy search:      At each layer, start from entry point
                    Move to whichever neighbour minimises dist(neighbour, query)
                    Stop when no neighbour is closer → descend to next layer

Beam search (ef):   Maintain a max-heap of ef candidates + a visited set
                    At each step, pop nearest unvisited candidate, add its neighbours
                    Higher ef → better recall, slower search — tunable tradeoff

Cosine similarity:  cos(θ) = (A · B) / (||A|| × ||B||)
                    Normalise all vectors on insert → query reduces to dot product

BM25 retrieval:     score(q, d) = Σ IDF(t) × (tf × (k1+1)) / (tf + k1 × (1−b + b×|d|/avgdl))
                    IDF: log((N − df + 0.5) / (df + 0.5))
                    Implemented with an in-memory inverted index

Hybrid blend:       final = α × normalised_vector_score + (1−α) × normalised_BM25_score
                    α tunable per query or globally in config

Recall@K:           |HNSW_top_K ∩ bruteforce_top_K| / K
                    Benchmark metric — validates index quality vs exact search
```

### Tech Stack

```
Language:         Go
Embeddings:       onnxruntime-go (ONNX Runtime bindings) + all-MiniLM-L6-v2 model
HNSW:             implemented from scratch (no external library)
BM25:             implemented from scratch (inverted index + IDF table)
Storage:          custom binary format for HNSW graph + BoltDB for chunk metadata
API:              gRPC (google.golang.org/protobuf) + net/http REST wrapper
File watching:    fsnotify
CLI:              cobra
LLM (RAG):        Ollama local API or OpenAI API (configurable)
```

### Build Time
**4–5 weeks** solo

### Job Market Signal
Vector search is the infrastructure layer of every RAG system being built in 2026. Implementing HNSW yourself signals AI infrastructure engineering at the algorithm level — not calling Pinecone. Directly relevant to: AI/ML backend roles, search infrastructure teams, any company building RAG. Real HN April 2026 JDs: Chroma (Go + vector search), healthcare AI startups (RAG + gRPC + Go + Python).

---
---

## Project 3 — Flowpipe

**Tagline:** Change Data Capture pipeline engine — reads Postgres WAL in real time, streams changes to any sink. Debezium, without the JVM.

**One-line interview answer:** "Debezium is the right solution for CDC but it requires a JVM, Kafka, and 30 steps to set up. Flowpipe is a single Go binary that reads Postgres WAL directly and streams changes to Elasticsearch, S3, or any webhook — five-minute setup."

---

### The Problem

Your app writes to Postgres. You need those changes reflected in Elasticsearch (for search), S3 (for analytics), and a webhook (to notify downstream services). Your options: dual-write in application code (fragile, misses DB-level changes), Debezium (correct approach, but requires JVM + Kafka + Zookeeper — completely impractical for a small team), or write a custom polling loop (misses deletes, has lag, hammers your DB).

Flowpipe reads Postgres's logical replication stream directly — the same mechanism Postgres uses for replication — and streams every INSERT, UPDATE, and DELETE to any configured sink in real time, with at-least-once delivery.

---

### How It's Different

| Tool | What's wrong with it |
|------|----------------------|
| Debezium | Requires JVM, Kafka, Zookeeper — impossible overhead for small teams |
| pg_logical | Low-level library, not a complete pipeline tool, requires custom code |
| Airbyte | Cloud-first, heavy UI, not embeddable, overkill for this use case |
| Polling loops | Miss deletes, have lag, increase DB load, break under high write volume |
| **Flowpipe** | Single Go binary, reads WAL directly, pluggable sinks, 5-minute setup |

---

### Exact Features

**1. Postgres WAL reader**
Connects to Postgres as a logical replication client. Creates a replication slot (or connects to an existing one). Receives a stream of WAL records: every INSERT, UPDATE, DELETE on the configured tables, decoded to structured JSON with before/after row values.

**2. LSN checkpointing — crash-safe**
After each event is confirmed by all sinks, records the Log Sequence Number (LSN) to BoltDB. On restart, resumes from the last confirmed LSN. No events are missed or double-processed.

**3. Pluggable sink interface**
A Sink interface with Write(event) and Confirm(). Built-in sinks: Elasticsearch (bulk indexing), S3 (JSONL files partitioned by table and date), HTTP webhooks (POST with configurable headers), and stdout (for debugging). New sinks are a single Go file implementing the interface.

**4. Transform pipeline**
Between WAL read and sink write, events pass through a configurable transform chain: filter by table name, filter by operation type (INSERT/UPDATE/DELETE), rename fields, redact columns (replace sensitive values with [REDACTED] before they reach any sink), add computed fields.

**5. Backpressure system**
An in-memory channel buffers events between the WAL reader and sink writers. If the buffer exceeds high_water (default 8,000 events), WAL reading is paused. It resumes when the buffer drops below low_water (default 2,000). Prevents memory exhaustion if a sink is temporarily slow without dropping events.

**6. At-least-once delivery**
LSN is only advanced after all sinks confirm the event. If a sink fails, Flowpipe retries with exponential backoff. If retries are exhausted, the event is written to a dead-letter log (local JSONL file) and LSN advances. Dead-letter events can be replayed manually.

**7. Replay from LSN**
`flowpipe replay --from-lsn 0/1A2B3C4D` replays all WAL events from a given LSN position. Useful for recovering a sink that was down or re-populating after a wipe.

**8. Prometheus metrics**
Exposes /metrics endpoint: events_total (counter by table and operation), wal_lag_bytes, sink_errors_total (by sink name), buffer_depth. Plug directly into Grafana.

---

### Exact User Flow — Step by Step

**Scenario:** Harry has a Postgres DB with an orders table. He wants changes synced to Elasticsearch and S3.

```
Step 1 — Postgres setup (one-time, on the DB server)
  ALTER SYSTEM SET wal_level = logical;
  SELECT pg_reload_conf();
  -- Flowpipe will create the replication slot automatically on first run

Step 2 — Install Flowpipe
  $ go install github.com/harry/flowpipe@latest

Step 3 — Initialise config
  $ flowpipe init
  → Creates flowpipe.yaml in current directory with commented template
  Output: "Config written to flowpipe.yaml — edit to configure your source and sinks."

Step 4 — Edit flowpipe.yaml
  source:
    dsn: "postgres://harry:pass@localhost:5432/myapp"
    slot_name: "flowpipe_main"
    tables: [orders, users]

  transforms:
    - type: redact
      table: users
      columns: [password_hash, ssn]

  sinks:
    - name: elasticsearch
      type: elasticsearch
      url: "http://localhost:9200"
      index_template: "{{.Table}}-events"
      bulk_size: 100
    - name: data-lake
      type: s3
      bucket: "harry-data-lake"
      prefix: "cdc/{{.Table}}/{{.Date}}"
      format: jsonl

  backpressure:
    high_water: 8000
    low_water: 2000

  dead_letter:
    path: ./dead-letter.jsonl

Step 5 — Start the pipeline
  $ flowpipe start
  Output:
    [flowpipe] connecting to postgres://...@localhost/myapp
    [flowpipe] replication slot "flowpipe_main" created (LSN: 0/15A3C20)
    [flowpipe] sink "elasticsearch" ready
    [flowpipe] sink "data-lake" ready
    [flowpipe] pipeline running — streaming WAL changes

Step 6 — Make a change in Postgres
  INSERT INTO orders (id, user_id, total) VALUES (1001, 42, 99.99);

  Flowpipe output immediately:
    [flowpipe] INSERT orders id=1001 → elasticsearch ✓ (12ms)
    [flowpipe] INSERT orders id=1001 → data-lake ✓ (4ms)
    [flowpipe] LSN advanced to 0/15A3C88

Step 7 — Check pipeline status
  $ flowpipe status
  ┌─────────────────────────────────────────────────────────┐
  │ Flowpipe Status                                         │
  ├─────────────────────┬───────────────────────────────────┤
  │ Current LSN         │ 0/15A4B20                         │
  │ WAL lag             │ 0 bytes (real-time)               │
  │ Events processed    │ 4,821 total                       │
  │ Throughput          │ 142 events/sec (10s EMA)          │
  │ Buffer depth        │ 12 / 8000                         │
  ├─────────────────────┼───────────────────────────────────┤
  │ elasticsearch       │ ✓ healthy  4,821 written  0 err  │
  │ data-lake           │ ✓ healthy  4,821 written  0 err  │
  └─────────────────────┴───────────────────────────────────┘

Step 8 — Elasticsearch goes down, then comes back
  Flowpipe output:
    [flowpipe] sink "elasticsearch" error — retrying (attempt 1, backoff 1s)
    [flowpipe] sink "elasticsearch" error — retrying (attempt 2, backoff 2s)
    [flowpipe] sink "elasticsearch" max retries exceeded — writing to dead-letter
    [flowpipe] dead-letter: 3 events written to ./dead-letter.jsonl

  <Elasticsearch comes back>
  $ flowpipe replay-dead-letter
  Output: "[flowpipe] replayed 3 dead-letter events → elasticsearch ✓"

Step 9 — Replay from a past LSN (re-populate a wiped sink)
  $ flowpipe replay --from-lsn 0/15A3C20
  → Streams all WAL events from that LSN to all configured sinks
```

---

### Algorithms / Systems Depth

```
WAL decoding:       pgoutput protocol — logical replication messages decoded to
                    Relation, Begin, Commit, Insert, Update, Delete structs

LSN checkpoint:     After sink confirms: BoltDB.Put("last_lsn", current_lsn)
                    On restart: resume = BoltDB.Get("last_lsn")
                    Safe LSN = min(confirmed_lsn across all sinks)
                    Weakest-link guarantee — never advance past slowest confirmed sink

Backpressure:       channel := make(chan Event, buffer_size)
                    if len(channel) > high_water → pause WAL reader
                    if len(channel) < low_water  → resume WAL reader

Retry backoff:      delay = min(base × 2^(attempt−1), max_delay)
                    base=1s, max_delay=60s

WAL lag metric:     current_lsn − last_processed_lsn (bytes)
                    → estimated seconds via lag_bytes / avg_write_rate

Throughput EMA:     rate = EMA(events_count / Δt, α=0.1) — 10s rolling window
```

### Tech Stack

```
Language:       Go
WAL:            jackc/pglogrepl (Postgres logical replication protocol client)
Database:       PostgreSQL 14+ with wal_level=logical
Checkpointing:  BoltDB (bbolt) — embedded, crash-safe key-value store
Sinks:          olivere/elastic (Elasticsearch), aws-sdk-go-v2 (S3), net/http (webhooks)
Metrics:        prometheus/client_golang
Config:         YAML (gopkg.in/yaml.v3)
CLI:            cobra
```

### Build Time
**3–4 weeks** solo

### Job Market Signal
Event-driven architecture and CDC are top-requested patterns in 2026 backend JDs — fintech, data engineering, platform engineering, microservices. Flowpipe demonstrates Postgres WAL internals knowledge (LSN, logical replication, slot management) that almost no backend developer has. Directly relevant to: Cast AI (reporting pipelines), Voodoo (staff Go + event streaming), any team running Postgres at scale.

---
---

## Project 4 — ghostdb

**Tagline:** Branch your Postgres database like you branch your code. Instant copy-on-write snapshots, locally, with no data copying.

**One-line interview answer:** "The only tool that does local Postgres branching is Neon — a cloud service. ghostdb does the same thing locally using Postgres's WAL to create copy-on-write branches that diverge from a snapshot point without copying any data upfront."

---

### The Problem

You need to test a destructive migration. Or a risky query on production-scale data. Your options: run it on production (terrifying), restore a pg_dump locally (takes 20 minutes for a large DB, data is already stale), or maintain a staging environment (expensive, always out of date, never has real data).

What you actually want: branch the database like you branch code. Instantly. Test the migration on the branch. See what happens. Drop the branch. Production untouched.

The only product that does this is Neon — a cloud-only Postgres service you have to migrate your database to. Nothing does it for a database you already have running locally.

ghostdb fills that gap.

---

### How It's Different

| Tool | What's wrong with it |
|------|----------------------|
| pg_dump + restore | 10–30 minutes for large DBs, manual, stale data by the time it's ready |
| Neon | Cloud-only — you must migrate your DB to their platform |
| Docker volume snapshots | Filesystem-level, requires stopping Postgres, not portable |
| pg_basebackup | Full physical copy — copies all data, not copy-on-write |
| **ghostdb** | Local, near-instant, copy-on-write, works with your existing Postgres, zero migration |

---

### Exact Features

**1. Near-instant branch creation**
`ghostdb branch create my-migration` records the current WAL LSN as the branch point. Starts a secondary Postgres sidecar process seeded via a fast streaming basebackup. The branch reaches the exact state of your main DB at the moment of branching. For most development DBs this takes under 30 seconds.

**2. Copy-on-write semantics**
The branch starts sharing the main DB's data files via filesystem hardlinks. When you write to the branch, only the modified 8KB pages are written to a branch-specific overlay directory. Reads that hit unmodified pages are served from the shared files. The branch never copies data it hasn't touched.

**3. Instant connection string**
After `ghostdb branch create`, you get a connection string immediately: `postgres://localhost:5433/myapp`. Connect to it like any Postgres DB. It's a fully functional Postgres instance — run migrations, drop tables, do whatever you want.

**4. Schema and data diff**
`ghostdb diff my-branch main` shows: which tables/columns/indexes were added, modified, or dropped on the branch vs main, row count deltas per table, and a generated SQL migration that would apply the branch's schema changes to main.

**5. Named snapshots**
`ghostdb snapshot save v1.2-pre-migration` saves the current DB state as a named point-in-time snapshot — stored as just a LSN + schema hash (no data copy, instant). `ghostdb snapshot restore v1.2-pre-release` creates a new branch from that snapshot for inspection.

**6. Branch lifecycle management**
`ghostdb branch list` shows all branches: creation time, LSN divergence from main (bytes written since branch point), and overlay disk usage. `ghostdb branch drop my-branch` cleans up the overlay directory and stops the sidecar process, reclaiming disk.

**7. Auto-expire**
Branches created with `--ttl 2h` are automatically dropped by the ghostdb daemon after the TTL expires, reclaiming disk space without manual cleanup.

**8. Migration testing shortcut**
`ghostdb test-migration ./migrations/0042_add_index.sql` creates a branch, applies the migration, reports success/failure + execution time + any errors + lock duration, then drops the branch automatically. One command to safely test a migration against real data.

---

### Exact User Flow — Step by Step

**Scenario:** Harry is about to run a risky migration adding a non-nullable column to a 10M-row table. He wants to test it first against real data.

```
Step 1 — Install and initialise
  $ go install github.com/harry/ghostdb@latest
  $ ghostdb init --dsn "postgres://harry:pass@localhost/myapp"
  → Verifies Postgres connection
  → Checks wal_level = logical (or prints instructions to set it)
  → Creates ~/.ghostdb/config.toml
  → Starts ghostdb daemon in background
  Output: "ghostdb initialised. Main: myapp @ localhost:5432"
           "Daemon running (PID 48291)"

Step 2 — Create a branch
  $ ghostdb branch create test-migration
  → Records current LSN: 0/22A1B40 (the branch point)
  → Runs pg_basebackup to a temp directory (~15s for 2GB DB)
  → Starts a sidecar Postgres on port 5433
  → Replays WAL to exact branch point LSN
  Output: "Branch 'test-migration' ready in 18.4s"
           "Connection: postgres://localhost:5433/myapp"
           "Branched from main at LSN 0/22A1B40"
           "Divergence: 0 bytes"

Step 3 — Test the risky migration on the branch
  $ psql "postgres://localhost:5433/myapp"
  ALTER TABLE orders ADD COLUMN fulfilled_at TIMESTAMPTZ NOT NULL DEFAULT now();
  -- On a 10M row table this takes 34s and holds a lock
  -- On the branch: it runs, Harry sees exactly what would happen

  Or use the shortcut:
  $ ghostdb test-migration ./migrations/0042_add_column.sql
  Output:
    "Running on branch 'ghostdb-test-20260509-143201'..."
    "✓ Migration completed in 34.2s"
    "Table 'orders': 10,000,000 rows affected, 0 errors"
    "Lock held for: 34.1s  ⚠ WARNING: this blocks all production writes"
    "Recommendation: run during low-traffic window or use online migration"
    "Branch dropped. Main DB is untouched."

Step 4 — Diff the branch against main
  $ ghostdb diff test-migration main
  SCHEMA DIFF
  ┌──────────┬──────────────────────────────────────────────┐
  │ table    │ change                                        │
  ├──────────┼──────────────────────────────────────────────┤
  │ orders   │ + column: fulfilled_at TIMESTAMPTZ NOT NULL  │
  └──────────┴──────────────────────────────────────────────┘

  MIGRATION SQL (to apply branch changes to main):
  ALTER TABLE orders ADD COLUMN fulfilled_at TIMESTAMPTZ NOT NULL DEFAULT now();

Step 5 — Save a snapshot before a big release
  $ ghostdb snapshot save v2.0-pre-release
  → Stores: {lsn: "0/22A1B40", schema_hash: "sha256:...", timestamp: ...}
  → Just a tiny JSON file — no data copy at all
  Output: "Snapshot 'v2.0-pre-release' saved (instant)"

Step 6 — Restore a snapshot after something goes wrong
  $ ghostdb snapshot restore v2.0-pre-release
  → Creates a branch from the snapshot LSN
  → Harry gets a connection string to the pre-release DB state
  → Can query it, compare data, export rows, investigate what changed
  Output: "Snapshot restored as branch 'restore-v2.0-pre-release'"
           "Connection: postgres://localhost:5434/myapp"

Step 7 — List all branches and clean up
  $ ghostdb branch list
  ┌─────────────────────────┬──────────────┬──────────────┬──────────┐
  │ branch                  │ created      │ divergence   │ disk     │
  ├─────────────────────────┼──────────────┼──────────────┼──────────┤
  │ test-migration          │ 8 mins ago   │ 142 MB       │ 142 MB   │
  │ restore-v2.0-pre-release│ 1 min ago    │ 0 bytes      │ 0 bytes  │
  └─────────────────────────┴──────────────┴──────────────┴──────────┘

  $ ghostdb branch drop test-migration
  Output: "Branch 'test-migration' dropped. 142 MB reclaimed."
```

---

### Algorithms / Systems Depth

```
Branch point:       Record current_lsn = pg_current_wal_lsn() at branch creation time
                    The branch "starts" at this LSN — all pages before it are shared

WAL replay:         Sidecar Postgres starts from basebackup
                    Replays WAL up to branch_lsn, then enters standby/read-write mode

Copy-on-write:      Postgres data files are 8KB pages
                    On first write to any page: copy original page to overlay directory
                    Subsequent reads: check overlay first, fall back to shared main files
                    Implemented via filesystem hardlinks + per-branch overlay directory

Divergence metric:  pg_wal_lsn_diff(current_lsn, branch_lsn) → bytes
                    Proxy for "how much work has this branch done since branching"

Schema diff:        Query information_schema on both connections
                    Structural diff: tables, columns (name+type+nullable+default),
                    constraints, indexes — output as minimal ALTER TABLE SQL

Snapshot storage:   {lsn, schema_hash, timestamp} stored as JSON
                    Zero data copy — restore = create branch from stored LSN
                    Relies on WAL retention being configured (wal_keep_size)
```

### Tech Stack

```
Language:          Go
Postgres protocol: jackc/pgx (connections + queries) + jackc/pglogrepl (WAL/replication)
                   Direct reuse of Flowpipe's WAL knowledge — same internals
Sidecar Postgres:  os/exec to manage secondary postgres process lifecycle
Overlay storage:   Filesystem hardlinks + per-branch data directory
Schema diff:       information_schema queries + structural Go comparison
Snapshots:         JSON files in ~/.ghostdb/snapshots/
CLI:               cobra
Daemon:            Background goroutines + os/signal for graceful shutdown + TTL expiry
Config:            ~/.ghostdb/config.toml
```

### Build Time
**3–4 weeks** solo — significantly faster if Flowpipe is built first (WAL and LSN knowledge transfers directly)

### Job Market Signal
Deep Postgres internals knowledge is rare and extremely valued at senior levels. Understanding WAL, logical replication, LSN positions, and page-level copy-on-write at this depth signals an engineer who actually understands how their database works — not just how to write queries. Directly relevant to: platform engineering roles, DBA-adjacent backend positions, any company where "how do we safely test migrations against real data?" is an open and unsolved problem.

---
---

## Project 5 — delhibus

**Tagline:** The DTC/DIMTS bus intelligence CLI Delhi deserves. Real-time tracking, honest reliability scores from historical data, and trip planning that tells you how often the bus actually shows up.

**One-line interview answer:** "Google Maps shows you the schedule. delhibus tells you whether the bus will actually come. It archives the OTD GTFS-RT feed every 30 seconds, builds per-route reliability models from weeks of history, and answers 'what time should I really leave?' — which no app does."

---

### The Problem

Google Maps uses the Delhi OTD GTFS-RT data to show scheduled bus times. But the OTD documentation itself states that scheduled times are rough estimates based on constant speed assumptions. The real question every Delhi commuter needs answered is: "If I'm at Kashmere Gate at 8:45am wanting to reach Saket, which bus should I actually try to catch, and when should I leave — accounting for real historical delays?"

No app answers this. delhibus does.

---

### The Data (Real, Free, Already Used by Google)

The Delhi government's Open Transit Data portal (otd.delhi.gov.in), built by IIIT-Delhi:
- **Static GTFS**: all routes, stops, stop sequences, timetables — free download, no key required
- **Real-time GTFS-RT**: live vehicle positions updated every 10 seconds — free API key (apply at otd.delhi.gov.in)
- Used in production by Google Maps, MapMyIndia, Here Maps, TCS Research, and Harvard researchers

---

### How It's Different

| Tool | What's wrong with it |
|------|----------------------|
| Google Maps | Shows schedule, not reliability. Cannot tell you how often a route runs late. |
| DTC app | Frequently broken, no reliability data, no offline mode |
| DMRC app | Metro only, doesn't cover buses at all |
| **delhibus** | Reliability scores from archived history, Kalman ETA, RAPTOR trip planning, offline mode |

---

### Exact Features

**1. GTFS static data loader**
Downloads and parses the full Delhi GTFS static dataset: routes.txt, stops.txt, trips.txt, stop_times.txt. Imports everything into a local SQLite database. Gives you the complete route-stop graph of Delhi's bus network locally — no internet needed to query it.

**2. Real-time vehicle tracker**
Polls the OTD GTFS-RT VehiclePositions.pb endpoint every 10 seconds. Decodes the Protobuf feed (using google.golang.org/protobuf). Shows live positions — GPS coordinates, bearing, speed — for any route or stop.

**3. Historical feed archiver (background daemon)**
`delhibus daemon start` runs a background process that polls the GTFS-RT feed every 30 seconds and writes every vehicle position reading to the archive SQLite DB. After 1–2 weeks of running, the reliability model has enough data to be useful. After a month, confidence is high.

**4. Reliability scoring per route/stop/time**
For each (route, stop, 30-minute time bucket): computes from the archive — median delay (actual arrival vs scheduled), P90 delay, on-time rate (arrived within 5 min of schedule), and no-show rate (no vehicle within 20 minutes of scheduled time). This is data that exists nowhere else for Delhi buses.

**5. RAPTOR trip planner**
Implements the RAPTOR algorithm (Round-Based Public Transit Routing — the correct algorithm for multi-leg transit routing, used by OpenTripPlanner and Google Maps internally). Finds all viable routes between two stops including transfers. Overlays reliability scores: routes with poor on-time rates are flagged. Works fully offline against the static GTFS data.

**6. Honest next-bus estimate**
`delhibus next "Kashmere Gate" --route 543` combines three signals: the scheduled arrival time, the live GPS position of the nearest bus, and the historical delay distribution for this route at this stop at this hour. The output is an honest estimate with a confidence indicator — not just the timetable.

**7. Kalman filter ETA smoothing**
Raw GPS positions from the GTFS-RT feed are noisy — buses teleport between readings. A Kalman filter smooths the position stream using a constant-velocity motion model, giving more stable and accurate ETA predictions than naive position interpolation.

**8. Offline mode**
After a one-time sync, all static data lives in local SQLite. Trip planning, stop lookups, route info — all work without any internet connection. Only real-time tracking and daemon archiving require connectivity.

---

### Exact User Flow — Step by Step

**Scenario:** Harry needs to get from Connaught Place to Lajpat Nagar on a weekday morning.

```
Step 1 — Install and setup
  $ go install github.com/harry/delhibus@latest
  $ delhibus setup --api-key YOUR_OTD_KEY
  → Downloads Delhi GTFS static data (~12MB zip, extracts to ~80MB)
  → Parses and imports all routes, stops, trips, stop_times into SQLite
  → Verifies GTFS-RT API connection
  Output: "Downloaded Delhi GTFS static data"
           "Imported: 600 routes  |  10,484 stops  |  1.2M stop_times"
           "Real-time feed: ✓ connected (702 vehicles active)"
           ""
           "Run 'delhibus daemon start' to begin building reliability history."
           "After 1 week: basic reliability data available."
           "After 4 weeks: high-confidence reliability scores."

Step 2 — Start the background archiver (run once, leave running)
  $ delhibus daemon start
  → Runs in background, polls GTFS-RT every 30 seconds
  → Stores every vehicle position to ~/.delhibus/archive.db
  → Silently builds the reliability model over time
  Output: "delhibus daemon started (PID 48291)"
           "Archiving feed every 30s to ~/.delhibus/archive.db"

Step 3 — Plan a trip
  $ delhibus plan "Connaught Place" "Lajpat Nagar" --depart now
  Output:
  TRIP PLAN — Connaught Place → Lajpat Nagar
  Departing: ~10:30am  |  Current time: 10:27am

  ────────────────────────────────────────────────────────────────
  OPTION 1  ★ Recommended
  Route 419  Shivaji Stadium (CP) → Lajpat Nagar Terminal  [direct]
  Depart:    10:34am scheduled  (next bus ~7 min away, live GPS)
  Arrive:    ~11:12am  (±8 min based on 30-day history)
  Duration:  38 min
  Reliability this hour:  ████████░░  71% on-time  |  +3min median delay
  Confidence: HIGH — 3 buses on this route currently within 2km of stop
  → Leave in 3 minutes to reach Shivaji Stadium comfortably

  ────────────────────────────────────────────────────────────────
  OPTION 2
  Route 181  Rajiv Chowk → Andrews Ganj + 400m walk
  Depart:    10:31am scheduled
  Arrive:    ~11:24am  (±14 min)
  Duration:  53 min  |  1 transfer
  Reliability this hour:  ████░░░░░░  38% on-time  |  +11min median delay
  ⚠ HIGH VARIANCE — this route runs late frequently at this hour

  Recommendation: Take Option 1. Leave now.

Step 4 — Check live next buses at a specific stop
  $ delhibus next "Shivaji Stadium" --route 419
  Output:
  Next buses — Route 419 at Shivaji Stadium

  Bus 1  DL1PD4821
    Scheduled:  10:34am
    Live GPS:   Patel Chowk, 1.2km away, heading north at 18km/h
    ETA:        10:33am  (Kalman estimate — 1 min early)

  Bus 2  DL1PD5103
    Scheduled:  10:49am
    Live GPS:   Barakhamba Road, 3.8km away
    ETA:        10:51am  (2 min late — consistent with P50 delay for this route)

Step 5 — Check route reliability report
  $ delhibus reliability 419 --stop "Shivaji Stadium" --days 30
  Output:
  Route 419 — Reliability at Shivaji Stadium (last 30 days, weekdays only)

  Time window    On-time  Median delay  P90 delay  No-show
  ──────────────────────────────────────────────────────────
  6am – 7am      61%      +4 min        +12 min    5%
  7am – 8am      43%      +9 min        +22 min    12%
  8am – 9am      31%      +14 min       +31 min    18%   ← avoid
  9am – 10am     58%      +5 min        +14 min    6%
  10am – 11am    71%      +3 min        +9 min     3%    ← best window
  11am – 12pm    67%      +4 min        +11 min    4%
  ──────────────────────────────────────────────────────────
  Overall (all day):  58% on-time  |  +6min median  |  7% no-show

Step 6 — Look up all routes at a stop
  $ delhibus stop "Kashmere Gate ISBT"
  Output:
  Stop: Kashmere Gate ISBT  (ID: DL-1042)
  Coordinates: 28.6676°N, 77.2276°E

  Routes serving this stop:
  Route  Towards                Next bus   Reliability (this hour)
  ─────────────────────────────────────────────────────────────────
  543    Dwarka Sector 21       10:38am    ★★★★☆ 74% on-time
  212    Rohini Sector 3        10:41am    ★★★☆☆ 52% on-time
  780    Badarpur Terminal      10:44am    ★★☆☆☆ 38% on-time  ⚠

Step 7 — Use offline (no internet)
  <Harry is underground, no connection>
  $ delhibus plan "Rohini" "Nehru Place"
  Output:
  [OFFLINE — showing scheduled times only, no real-time or reliability data]

  Route 534  Rohini Sector 3 → Nehru Place
  Depart:  10:52am (scheduled)
  Arrive:  12:08pm (scheduled)
  Duration: 76 min  |  1 transfer at Kashmere Gate ISBT

  Note: Real-time ETA and reliability data unavailable offline.
        Reliability from last sync: Route 534 — 61% on-time, +5min median delay.
```

---

### Algorithms / Math Involved

```
RAPTOR routing:     Round-Based Public Transit Routing Algorithm
                    Each "round" = one additional transit leg (transfer)
                    Maintains earliest arrival time array, updated per round
                    Finds all Pareto-optimal journeys: min time AND min transfers
                    Correct for transit — handles frequency-based services naturally
                    Time complexity: O(p × |trips|) where p = max transfers allowed

Reliability model:  For each (route_id, stop_id, time_bucket_30min):
                      observations = all historical vehicle arrivals at this stop
                      delay_i = actual_arrival_i − scheduled_arrival_i  (seconds)
                      on_time_rate = |{i : delay_i < 300}| / |observations|
                      median_delay = percentile(delays, 50)
                      p90_delay    = percentile(delays, 90)
                      no_show_rate = |scheduled with no arrival within 1200s| / |scheduled|

Kalman filter ETA:  State vector: [position_km_along_route, velocity_km/s]
                    Motion model: x̂ₖ = F × xₖ₋₁  (constant velocity)
                    Update step:  xₖ = x̂ₖ + Kₖ × (zₖ − H × x̂ₖ)
                    Kₖ = Pₖ⁻ × Hᵀ × (H × Pₖ⁻ × Hᵀ + R)⁻¹  (Kalman gain)
                    R = GPS measurement noise covariance (tuned empirically)
                    Smooths noisy 10-second position reports into stable ETA

Confidence blend:   eta = w₁ × kalman_eta + w₂ × (scheduled + historical_median_delay)
                    w₁ = 1.0 if last GPS reading < 30s old, else decays exponentially
                    w₂ = 1 − w₁
                    Falls back gracefully to history-adjusted schedule when GPS is stale
```

### Tech Stack

```
Language:         Go
GTFS static:      Custom CSV parser for all GTFS files → SQLite (mattn/go-sqlite3)
GTFS-RT:          google.golang.org/protobuf (decodes binary VehiclePositions.pb feed)
Routing:          RAPTOR algorithm implemented from scratch
Reliability DB:   SQLite — vehicle position archive table + precomputed reliability scores
Kalman filter:    Implemented from scratch in pure Go (no ML library)
CLI:              cobra
Daemon:           Background goroutines + os/signal for graceful shutdown
Offline:          All static data + precomputed scores in local SQLite — fully offline
```

### Build Time
**3–4 weeks** solo. Suggested incremental order: (1) static GTFS loader + stop lookup, (2) real-time tracker, (3) RAPTOR trip planner, (4) daemon archiver, (5) reliability model + Kalman ETA.

### Why Build It
This one is personal. You live in Delhi. You take DTC buses. The data to answer "will the bus actually come?" exists — it's public, free, updated every 10 seconds — but no tool surfaces it usefully. delhibus is the tool you'll actually use every day.

It also makes a compelling portfolio story: "I used open government data to build a tool I use daily, and in the process discovered that Route 181 during 8–9am rush hour has an 18% no-show rate — which is invisible to every existing app." That kind of observation, grounded in real data from real infrastructure, is a story that stands out.

---
---

## Summary

| # | Project | Core Problem | Math / Depth | Build Time | Signal |
|---|---------|-------------|--------------|------------|--------|
| 1 | **promtrace** | Zero visibility into LLM calls | HTTP MITM, SSE, Myers diff | 2–3 wks | AI infra |
| 2 | **Vectraflow** | Private semantic search, HNSW from scratch | HNSW, BM25, ONNX | 4–5 wks | AI/ML infra |
| 3 | **Flowpipe** | CDC without JVM + Kafka | WAL, LSN, backpressure | 3–4 wks | Data eng, backend |
| 4 | **ghostdb** | Branch your DB like code, locally | CoW, WAL replay, schema diff | 3–4 wks | DB internals, platform |
| 5 | **delhibus** | Delhi buses have no honest reliability data | RAPTOR, Kalman filter, reliability model | 3–4 wks | Civic tech, personal |
