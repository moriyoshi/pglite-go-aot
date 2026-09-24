# pglite-go-aot

Ahead-of-time (goccy/wasm2go-transpiled) backend for
[pglite-go](https://github.com/moriyoshi/pglite-go): PGlite / PostgreSQL as
**pure Go** — no CGo, no wasm engine.

This module carries the ~200 MB of transpiled Go so the core module stays thin.
Enable the AOT backend with a `database/sql`-style blank import:

```go
import (
    "github.com/moriyoshi/pglite-go"
    _ "github.com/moriyoshi/pglite-go-aot" // registers the AOT backend
)
// build with -tags aot
```

The transpiled code is regenerated and re-tagged per PGlite release using the
tools in the core module (`internal/wasmpass`, `internal/genaot`, `cmd/wasmpass`).
