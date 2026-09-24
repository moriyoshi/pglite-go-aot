//go:build aot

package pgliteaot

import (
	"fmt"
	"testing"

	pglite "github.com/moriyoshi/pglite-go"
)

// TestColdInit proves the companion module enables end-to-end pure-Go PGlite:
// being in package pgliteaot, register.go's init() has registered the AOT
// backend, so pglite.Open runs initdb (transpiled here) against an empty persist
// dir, bootstraps a cluster, and serves a query. WasmDir points at the core
// module's assets (the data bundle is still fetched as usual; only the code is
// transpiled).
func TestColdInit(t *testing.T) {
	dir := t.TempDir()
	db, err := pglite.Open(pglite.Config{
		WasmDir:    "../pglite-go/wasm",
		PersistDir: dir,
		Database:   "template1",
	})
	if err != nil {
		t.Fatalf("Open (cold init): %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT 42")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(rows.Rows) != 1 || len(rows.Rows[0]) != 1 || fmt.Sprint(rows.Rows[0][0]) != "42" {
		t.Fatalf("unexpected result: %v", rows.Rows)
	}
	t.Log("COLD INIT + QUERY via the pglite-go-aot companion module (pure Go)")
}
