# Rift – a minimal Redis-compatible server written in Go

> Phase 1 – MVP Protocol & Human Mode ✅  
> NOTE: This document will expand with each completed phase.

### What is Rift

Rift is a learning project that re-implements core Redis features from scratch in Go.

Phase 1 gives you a working TCP server that speaks RESP2 and a developer-friendly “human” REPL.

### Quick Start (Phase 1)

```bash
# Clone
git clone https://github.com/DarsenOP/Rift.git && cd Rift

# Build
make fmt                   # format the code
make lint                  # run a linter
make test                  # run the tests
make build                 # produces ./rift

# Run RESP mode (default)
./bin/rift                 # listens on :6380

# Run Human mode (telnet-friendly)
./bin/rift -human          # same port, plain-text commands

# Test
telnet localhost 6380
> PING
< +PONG
> PING hello
< $5
< hello

...

make clean                 # Clean the executable
```

### Features Delivered

| Feature      | Status | CLI Example               |
| ------------ | ------ | ------------------------- |
| RESP2 parser | ✅     | `*1\r\n$4\r\nPING\r\n`    |
| RESP2 writer | ✅     | `+PONG`, `$5\r\nhello`, … |
| Human parser | ✅     | `PING hello` → auto-RESP  |
| Commands     | ✅     | `PING`, `COMMAND`         |
| TCP server   | ✅     | goroutine-per-connection  |
| Tests (CI)   | ✅     | `make test` (100 % cov)   |


### Project Layout

```bash
cmd/rift          // server entry point
internal/resp     // parser + writer + human mode
internal/server   // command dispatcher
internal/version  // build-time version injection
```

### Development

See [CONTRIBUTING POLICY](CONTRIBUTING.md) for workflow & branching rules.

### License
MIT – see [LICENSE](LICENSE).
