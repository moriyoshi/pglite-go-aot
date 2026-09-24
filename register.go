//go:build aot

// Package pgliteaot registers the ahead-of-time (goccy/wasm2go-transpiled)
// PGlite backend. It carries the ~200 MB of transpiled Go so the core
// github.com/moriyoshi/pglite-go module stays thin. Blank-import it and build
// with -tags aot to run PGlite as pure Go — no CGo, no wasm engine:
//
//	import (
//	    "github.com/moriyoshi/pglite-go"
//	    _ "github.com/moriyoshi/pglite-go-aot"
//	)
package pgliteaot

import (
	pglite "github.com/moriyoshi/pglite-go"
	"github.com/moriyoshi/pglite-go-aot/initdbwasm"
	"github.com/moriyoshi/pglite-go-aot/pgwasm"
)

func init() { pglite.RegisterAOT(pgwasm.NewAOT, initdbwasm.NewAOT) }
