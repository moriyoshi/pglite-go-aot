module github.com/moriyoshi/pglite-go-aot

go 1.25.5

require (
	github.com/moriyoshi/pglite-go v0.0.0-00010101000000-000000000000
	github.com/tetratelabs/wazero v1.11.0
)

require (
	github.com/bytecodealliance/wasmtime-go/v34 v34.0.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
)

// during development, build against the local main module
replace github.com/moriyoshi/pglite-go => ../pglite-go
