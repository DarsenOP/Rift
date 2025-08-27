# Rift – a minimal Redis-compatible server written in Go

> Phase 0 – Foundation & Setup

> NOTE: This document will expand with each completed phase.

### What is Rift

Rift is a learning project that re-implements core Redis features from scratch in Go.
Today it is only scaffolding; tomorrow it will speak RESP, store data, and survive restarts.

### Goals

- 100 % Go implementation
- Redis-protocol (RESP) compatible
- Single-binary deployment
- Clean, idiomatic code with tests & CI

### Build & Run (Phase 0)

Copy the following code 

```bash
# Clone
git clone https://github.com/DarsenOP/Rift.git && cd Rift
# Build
make build            # produces ./rift binary
# Run (placeholder – starts minimal listener)
./rift
```

### Development

See [CONTRIBUTING POLICY](CONTRIBUTING.md) for workflow & branching rules.

### License
MIT – see [LICENSE](LICENSE).
