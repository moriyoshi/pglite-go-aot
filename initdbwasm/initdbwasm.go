package initdbwasm

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"math/bits"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"
)

type Wasi_snapshot_preview1Imports interface {
	Fd_close(m *Module, l0 int32) int32
	Fd_write(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	Fd_read(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	Fd_sync(m *Module, l0 int32) int32
	Environ_sizes_get(m *Module, l0 int32, l1 int32) int32
	Environ_get(m *Module, l0 int32, l1 int32) int32
	Fd_fdstat_get(m *Module, l0 int32, l1 int32) int32
	Fd_seek(m *Module, l0 int32, l1 int64, l2 int32, l3 int32) int32
	Proc_exit(m *Module, l0 int32)
}
type EnvImports interface {
	Getaddrinfo(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	Exit(m *Module, l0 int32)
	Invoke_iiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	Invoke_ii(m *Module, l0 int32, l1 int32) int32
	Invoke_vii(m *Module, l0 int32, l1 int32, l2 int32)
	X_tzset_js(m *Module, l0 int32, l1 int32, l2 int32, l3 int32)
	X_abort_js(m *Module)
	X__syscall_faccessat(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	X__syscall_chmod(m *Module, l0 int32, l1 int32) int32
	Emscripten_date_now(m *Module) float64
	X__syscall_openat(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	X__syscall_fcntl64(m *Module, l0 int32, l1 int32, l2 int32) int32
	X__syscall_ioctl(m *Module, l0 int32, l1 int32, l2 int32) int32
	X__syscall_dup3(m *Module, l0 int32, l1 int32, l2 int32) int32
	X__syscall_stat64(m *Module, l0 int32, l1 int32) int32
	X__syscall_newfstatat(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	X__syscall_lstat64(m *Module, l0 int32, l1 int32) int32
	X__syscall_getcwd(m *Module, l0 int32, l1 int32) int32
	Emscripten_get_now(m *Module) float64
	X__syscall_mkdirat(m *Module, l0 int32, l1 int32, l2 int32) int32
	X_mktime_js(m *Module, l0 int32) int64
	X_localtime_js(m *Module, l0 int64, l1 int32)
	X_munmap_js(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int64) int32
	X_mmap_js(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int64, l5 int32, l6 int32) int32
	X__syscall_fadvise64(m *Module, l0 int32, l1 int64, l2 int64, l3 int32) int32
	X_emscripten_runtime_keepalive_clear(m *Module)
	X__call_sighandler(m *Module, l0 int32, l1 int32)
	X__syscall_getdents64(m *Module, l0 int32, l1 int32, l2 int32) int32
	X__syscall_readlinkat(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	X__syscall_rmdir(m *Module, l0 int32) int32
	X_setitimer_js(m *Module, l0 int32, l1 float64) int32
	X__syscall_symlinkat(m *Module, l0 int32, l1 int32, l2 int32) int32
	X__syscall_unlinkat(m *Module, l0 int32, l1 int32, l2 int32) int32
	Emscripten_resize_heap(m *Module, l0 int32) int32
	X_emscripten_throw_longjmp(m *Module)
}
type Module struct {
	memory                 []byte
	maxMem                 uint64
	M                      unsafe.Pointer
	t0                     []any
	g0                     int32
	g1                     int32
	g2                     int32
	g3                     int32
	g4                     int32
	g5                     int32
	g6                     int32
	g7                     int32
	g8                     int32
	g9                     int32
	g10                    int32
	g11                    int32
	g12                    int32
	g13                    int32
	g14                    int32
	g15                    int32
	g16                    int32
	g17                    int32
	g18                    int32
	g19                    int32
	g20                    int32
	g21                    int32
	g22                    int32
	g23                    int32
	g24                    int32
	g25                    int32
	g26                    int32
	g27                    int32
	g28                    int32
	g29                    int32
	g30                    int32
	g31                    int32
	g32                    int32
	g33                    int32
	g34                    int32
	g35                    int32
	g36                    int32
	g37                    int32
	g38                    int32
	g39                    int32
	g40                    int32
	g41                    int32
	g42                    int32
	g43                    int32
	g44                    int32
	g45                    int32
	g46                    int32
	g47                    int32
	g48                    int32
	g49                    int32
	g50                    int32
	g51                    int32
	g52                    int32
	g53                    int32
	g54                    int32
	g55                    int32
	g56                    int32
	g57                    int32
	wasi_snapshot_preview1 Wasi_snapshot_preview1Imports
	env                    EnvImports
	memMu                  *sync.Mutex
	memSize                *atomic.Uint64
	dataEnd                uint32
	memShared              bool
	threads                *threadPool
	threadStart            func(*Module, int32, int32)
}

// WasiExitError is the sentinel that the recover layer of SafeInvokeExport
// promotes Proc_exit() panics into, so a wasm-level exit doesn't kill the
// host process and the caller can read the exit code instead.
type WasiExitError struct{ Code int32 }

func (e *WasiExitError) Error() string {
	return "wasi: proc_exit(" + itoa32(e.Code) + ")"
}

// itoa32 is a tiny dependency-free strconv replacement so this file
// doesn't drag in fmt for its sole error path.
func itoa32(v int32) string {
	if v == 0 {
		return "0"
	}
	neg := false
	if v < 0 {
		v = -v
		neg = true
	}
	var buf [12]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

// FS is the read/write filesystem backend the WASI host opens files through.
// It abstracts the default os-backed filesystem so an embedder can supply an
// alternative — an in-memory FS, an overlay, a read-only bundle, ... — and
// have every guest path operation (open, stat, mkdir, readdir, write, ...)
// routed to it. It is a write-capable superset of io/fs.FS.
//
// Names are GUEST paths relative to the preopen root: slash-separated, with no
// leading slash (e.g. "encodings/__init__.py", or "" for the root). Methods
// should return the standard fs errors (fs.ErrNotExist, fs.ErrExist,
// fs.ErrPermission) so the host maps them to the right wasi errno.
type FS interface {
	// OpenFile mirrors os.OpenFile: flag is O_RDONLY/O_WRONLY/O_RDWR optionally
	// OR'd with O_CREATE/O_EXCL/O_TRUNC/O_APPEND. The returned File must
	// support the operations the mode implies.
	OpenFile(name string, flag int, perm os.FileMode) (File, error)
	Mkdir(name string, perm os.FileMode) error
	Remove(name string) error
	Rename(oldName, newName string) error
	Stat(name string) (os.FileInfo, error)
	Lstat(name string) (os.FileInfo, error)
	Symlink(oldName, newName string) error
	Readlink(name string) (string, error)
	Link(oldName, newName string) error
}

// File is an open file handle returned by FS.OpenFile. *os.File satisfies it,
// so the default os backend needs no wrapper.
type File interface {
	Read(p []byte) (int, error)
	ReadAt(p []byte, off int64) (int, error)
	Write(p []byte) (int, error)
	WriteAt(p []byte, off int64) (int, error)
	Seek(offset int64, whence int) (int64, error)
	Close() error
	Stat() (os.FileInfo, error)
	ReadDir(n int) ([]os.DirEntry, error)
	Sync() error
	Truncate(size int64) error
	Name() string
}

// osFS is the default FS backend: a thin pass-through to the host filesystem,
// scoped to root (the preopen directory). root "" or "/" means no rewriting.
type osFS struct{ root string }

func (o osFS) join(name string) string {
	if o.root == "" || o.root == "/" {
		return "/" + name
	}
	return filepath.Join(o.root, name)
}
func (o osFS) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	f, err := os.OpenFile(o.join(name), flag, perm)
	if err != nil {
		return nil, err
	}
	return f, nil
}
func (o osFS) Mkdir(name string, perm os.FileMode) error { return os.Mkdir(o.join(name), perm) }
func (o osFS) Chmod(name string, mode os.FileMode) error { return os.Chmod(o.join(name), mode) }
func (o osFS) Remove(name string) error                  { return os.Remove(o.join(name)) }
func (o osFS) Rename(a, b string) error                  { return os.Rename(o.join(a), o.join(b)) }
func (o osFS) Stat(name string) (os.FileInfo, error)     { return os.Stat(o.join(name)) }
func (o osFS) Lstat(name string) (os.FileInfo, error)    { return os.Lstat(o.join(name)) }
func (o osFS) Symlink(target, name string) error         { return os.Symlink(target, o.join(name)) }
func (o osFS) Readlink(name string) (string, error)      { return os.Readlink(o.join(name)) }
func (o osFS) Link(a, b string) error                    { return os.Link(o.join(a), o.join(b)) }

// MemFS is an in-memory read/write FS. Each value is an independent tree, so
// two interpreters given separate MemFS values cannot observe each other's
// files (full per-interpreter filesystem isolation, no disk). Build one with
// NewMemFS. Safe for concurrent use.
type MemFS struct {
	mu   sync.Mutex
	root *memNode
}

// NewMemFS returns an empty in-memory filesystem with a root directory.
func NewMemFS() *MemFS {
	return &MemFS{root: &memNode{dir: true, mode: os.ModeDir | 0o755, modTime: time.Unix(0, 0), children: map[string]*memNode{}}}
}

// memNode is a file or directory in a MemFS tree.
type memNode struct {
	name     string
	dir      bool
	mode     os.FileMode
	modTime  time.Time
	data     []byte
	children map[string]*memNode
}

func memSplit(name string) []string {
	name = strings.Trim(name, "/")
	if name == "" {
		return nil
	}
	raw := strings.Split(name, "/")
	out := make([]string, 0, len(raw))
	for _, p := range raw {
		switch p {
		case "", ".":
		case "..":
			if len(out) > 0 {
				out = out[:len(out)-1]
			}
		default:
			out = append(out, p)
		}
	}
	return out
}

// lookup resolves name to a node. Caller holds fsys.mu.
func (fsys *MemFS) lookup(name string) (*memNode, error) {
	n := fsys.root
	for _, part := range memSplit(name) {
		if !n.dir {
			return nil, fs.ErrNotExist
		}
		c, ok := n.children[part]
		if !ok {
			return nil, fs.ErrNotExist
		}
		n = c
	}
	return n, nil
}

// lookupParent resolves the parent dir of name. Caller holds fsys.mu.
func (fsys *MemFS) lookupParent(name string) (*memNode, string, error) {
	parts := memSplit(name)
	if len(parts) == 0 {
		return nil, "", fs.ErrInvalid
	}
	n := fsys.root
	for _, part := range parts[:len(parts)-1] {
		c, ok := n.children[part]
		if !ok || !c.dir {
			return nil, "", fs.ErrNotExist
		}
		n = c
	}
	return n, parts[len(parts)-1], nil
}

func (fsys *MemFS) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	fsys.mu.Lock()
	defer fsys.mu.Unlock()
	node, err := fsys.lookup(name)
	if err != nil {
		if flag&os.O_CREATE == 0 {
			return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
		}
		parent, base, perr := fsys.lookupParent(name)
		if perr != nil {
			return nil, &fs.PathError{Op: "open", Path: name, Err: perr}
		}
		node = &memNode{name: base, mode: perm & 0o777, modTime: time.Now()}
		parent.children[base] = node
	} else {
		if flag&os.O_EXCL != 0 && flag&os.O_CREATE != 0 {
			return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrExist}
		}
		if flag&os.O_TRUNC != 0 && !node.dir {
			node.data = node.data[:0]
			node.modTime = time.Now()
		}
	}
	f := &memFile{fsys: fsys, node: node}
	if flag&os.O_APPEND != 0 {
		f.off = int64(len(node.data))
	}
	return f, nil
}

func (fsys *MemFS) Mkdir(name string, perm os.FileMode) error {
	fsys.mu.Lock()
	defer fsys.mu.Unlock()
	parent, base, err := fsys.lookupParent(name)
	if err != nil {
		return &fs.PathError{Op: "mkdir", Path: name, Err: err}
	}
	if _, ok := parent.children[base]; ok {
		return &fs.PathError{Op: "mkdir", Path: name, Err: fs.ErrExist}
	}
	parent.children[base] = &memNode{name: base, dir: true, mode: os.ModeDir | (perm & 0o777), modTime: time.Now(), children: map[string]*memNode{}}
	return nil
}

func (fsys *MemFS) Remove(name string) error {
	fsys.mu.Lock()
	defer fsys.mu.Unlock()
	parent, base, err := fsys.lookupParent(name)
	if err != nil {
		return &fs.PathError{Op: "remove", Path: name, Err: err}
	}
	n, ok := parent.children[base]
	if !ok {
		return &fs.PathError{Op: "remove", Path: name, Err: fs.ErrNotExist}
	}
	if n.dir && len(n.children) > 0 {
		return &fs.PathError{Op: "remove", Path: name, Err: fs.ErrInvalid}
	}
	delete(parent.children, base)
	return nil
}

func (fsys *MemFS) Rename(oldName, newName string) error {
	fsys.mu.Lock()
	defer fsys.mu.Unlock()
	op, ob, err := fsys.lookupParent(oldName)
	if err != nil {
		return &fs.PathError{Op: "rename", Path: oldName, Err: err}
	}
	node, ok := op.children[ob]
	if !ok {
		return &fs.PathError{Op: "rename", Path: oldName, Err: fs.ErrNotExist}
	}
	np, nb, err := fsys.lookupParent(newName)
	if err != nil {
		return &fs.PathError{Op: "rename", Path: newName, Err: err}
	}
	delete(op.children, ob)
	node.name = nb
	np.children[nb] = node
	return nil
}

func (fsys *MemFS) Stat(name string) (os.FileInfo, error) {
	fsys.mu.Lock()
	defer fsys.mu.Unlock()
	n, err := fsys.lookup(name)
	if err != nil {
		return nil, &fs.PathError{Op: "stat", Path: name, Err: err}
	}
	return n.info(), nil
}

func (fsys *MemFS) Lstat(name string) (os.FileInfo, error) { return fsys.Stat(name) }

// memfs has no symlinks/hardlinks.
func (fsys *MemFS) Symlink(_, _ string) error { return fs.ErrPermission }
func (fsys *MemFS) Link(_, _ string) error    { return fs.ErrPermission }
func (fsys *MemFS) Readlink(name string) (string, error) {
	return "", &fs.PathError{Op: "readlink", Path: name, Err: fs.ErrInvalid}
}

// Chtimes implements the optional chtimesFS capability.
func (fsys *MemFS) Chtimes(name string, _ time.Time, mtime time.Time) error {
	fsys.mu.Lock()
	defer fsys.mu.Unlock()
	n, err := fsys.lookup(name)
	if err != nil {
		return &fs.PathError{Op: "chtimes", Path: name, Err: err}
	}
	n.modTime = mtime
	return nil
}

// MkdirAll creates name and any missing parents. Exposed so embedders can
// populate the FS (e.g. unpack a stdlib bundle) before handing it to a module.
func (fsys *MemFS) MkdirAll(name string, perm os.FileMode) error {
	fsys.mu.Lock()
	defer fsys.mu.Unlock()
	n := fsys.root
	for _, part := range memSplit(name) {
		c, ok := n.children[part]
		if !ok {
			c = &memNode{name: part, dir: true, mode: os.ModeDir | (perm & 0o777), modTime: time.Now(), children: map[string]*memNode{}}
			n.children[part] = c
		} else if !c.dir {
			return &fs.PathError{Op: "mkdir", Path: name, Err: fs.ErrExist}
		}
		n = c
	}
	return nil
}

// WriteFile creates (or overwrites) a file with data, making parent dirs as
// needed. Exposed for pre-populating the FS.
func (fsys *MemFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	if parts := memSplit(name); len(parts) > 1 {
		if err := fsys.MkdirAll(strings.Join(parts[:len(parts)-1], "/"), 0o755); err != nil {
			return err
		}
	}
	fsys.mu.Lock()
	defer fsys.mu.Unlock()
	parent, base, err := fsys.lookupParent(name)
	if err != nil {
		return &fs.PathError{Op: "writefile", Path: name, Err: err}
	}
	cp := make([]byte, len(data))
	copy(cp, data)
	parent.children[base] = &memNode{name: base, mode: perm & 0o777, modTime: time.Now(), data: cp}
	return nil
}

func (n *memNode) info() os.FileInfo {
	if n.dir {
		return memFileInfo{name: n.name, mode: os.ModeDir | (n.mode & 0o777), modTime: n.modTime}
	}
	return memFileInfo{name: n.name, size: int64(len(n.data)), mode: n.mode & 0o777, modTime: n.modTime}
}

type memFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi memFileInfo) Name() string       { return fi.name }
func (fi memFileInfo) Size() int64        { return fi.size }
func (fi memFileInfo) Mode() os.FileMode  { return fi.mode }
func (fi memFileInfo) ModTime() time.Time { return fi.modTime }
func (fi memFileInfo) IsDir() bool        { return fi.mode.IsDir() }
func (fi memFileInfo) Sys() any           { return nil }

type memDirEntry struct{ n *memNode }

func (e memDirEntry) Name() string { return e.n.name }
func (e memDirEntry) IsDir() bool  { return e.n.dir }
func (e memDirEntry) Type() os.FileMode {
	if e.n.dir {
		return os.ModeDir
	}
	return 0
}
func (e memDirEntry) Info() (os.FileInfo, error) { return e.n.info(), nil }

// memFile is an open handle into a MemFS node.
type memFile struct {
	fsys   *MemFS
	node   *memNode
	off    int64
	dirOff int
}

func (f *memFile) Name() string { return f.node.name }
func (f *memFile) Close() error { return nil }
func (f *memFile) Sync() error  { return nil }

func (f *memFile) Stat() (os.FileInfo, error) {
	f.fsys.mu.Lock()
	defer f.fsys.mu.Unlock()
	return f.node.info(), nil
}

func (f *memFile) Read(p []byte) (int, error) {
	f.fsys.mu.Lock()
	defer f.fsys.mu.Unlock()
	if f.node.dir {
		return 0, &fs.PathError{Op: "read", Path: f.node.name, Err: fs.ErrInvalid}
	}
	if f.off >= int64(len(f.node.data)) {
		return 0, io.EOF
	}
	n := copy(p, f.node.data[f.off:])
	f.off += int64(n)
	return n, nil
}

func (f *memFile) ReadAt(p []byte, off int64) (int, error) {
	f.fsys.mu.Lock()
	defer f.fsys.mu.Unlock()
	if off >= int64(len(f.node.data)) {
		return 0, io.EOF
	}
	n := copy(p, f.node.data[off:])
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

// writeAt grows node.data as needed and writes p at off. Caller holds the lock.
func (f *memFile) writeAt(p []byte, off int64) int {
	end := off + int64(len(p))
	if end > int64(len(f.node.data)) {
		grown := make([]byte, end)
		copy(grown, f.node.data)
		f.node.data = grown
	}
	copy(f.node.data[off:], p)
	f.node.modTime = time.Now()
	return len(p)
}

func (f *memFile) Write(p []byte) (int, error) {
	f.fsys.mu.Lock()
	defer f.fsys.mu.Unlock()
	n := f.writeAt(p, f.off)
	f.off += int64(n)
	return n, nil
}

func (f *memFile) WriteAt(p []byte, off int64) (int, error) {
	f.fsys.mu.Lock()
	defer f.fsys.mu.Unlock()
	return f.writeAt(p, off), nil
}

func (f *memFile) Seek(offset int64, whence int) (int64, error) {
	f.fsys.mu.Lock()
	defer f.fsys.mu.Unlock()
	switch whence {
	case io.SeekStart:
		f.off = offset
	case io.SeekCurrent:
		f.off += offset
	case io.SeekEnd:
		f.off = int64(len(f.node.data)) + offset
	}
	return f.off, nil
}

func (f *memFile) Truncate(size int64) error {
	f.fsys.mu.Lock()
	defer f.fsys.mu.Unlock()
	if size <= int64(len(f.node.data)) {
		f.node.data = f.node.data[:size]
	} else {
		grown := make([]byte, size)
		copy(grown, f.node.data)
		f.node.data = grown
	}
	f.node.modTime = time.Now()
	return nil
}

func (f *memFile) ReadDir(n int) ([]os.DirEntry, error) {
	f.fsys.mu.Lock()
	defer f.fsys.mu.Unlock()
	if !f.node.dir {
		return nil, &fs.PathError{Op: "readdir", Path: f.node.name, Err: fs.ErrInvalid}
	}
	names := make([]string, 0, len(f.node.children))
	for name := range f.node.children {
		names = append(names, name)
	}
	sort.Strings(names)
	if f.dirOff >= len(names) {
		if n <= 0 {
			return nil, nil
		}
		return nil, io.EOF
	}
	end := len(names)
	if n > 0 && f.dirOff+n < end {
		end = f.dirOff + n
	}
	out := make([]os.DirEntry, 0, end-f.dirOff)
	for _, name := range names[f.dirOff:end] {
		out = append(out, memDirEntry{f.node.children[name]})
	}
	f.dirOff = end
	return out, nil
}

// wasiOpen is one entry in WasiStubs' fd table. Stdio entries are nil-file
// markers (writes go to the OS handles directly via the WasiStubs fields).
// The conn arm carries a net.Conn for sockets opened via Sock_accept.
type wasiOpen struct {
	f        File
	conn     net.Conn
	listener net.Listener
	isDir    bool
	isSocket bool   // created by Sock_socket, may not yet have a conn
	path     string // guest path relative to the preopen root
	fdflags  int32  // last fdflags set via Path_open or Fd_fdstat_set_flags
	dirCache []os.DirEntry
	// stdio marks an alias of an interpreter stream (1/2/3 = the
	// configured stdin/stdout/stderr; 0 = not an alias). Fd_dup of a bare
	// fd 0/1/2 creates one; closing it never touches the real stream.
	stdio int8
	// refs counts EXTRA table slots sharing this entry (dup/dup2):
	// closeWasiOpen only closes the descriptor when it reaches zero.
	refs int32
}

// WasiStubs is the default Go-native implementation of wasi_snapshot_preview1.
// State is owned per-Module via NewWithWASI / DefaultWASI.
type WasiStubs struct {
	mu sync.Mutex

	// stdin/stdout/stderr back guest fds 0/1/2. They default to the host
	// os.Std* (DefaultWASI) but can be redirected to any io.Reader/io.Writer
	// (an in-process buffer, pipe, ...) via SetStdin/SetStdout/SetStderr, so an
	// embedder can feed input and capture/stream output without touching the
	// host process stdio.
	stdin          io.Reader
	stdout, stderr io.Writer
	fdTable        map[int32]*wasiOpen
	nextFD         int32
	args, env      []string
	monoStart      time.Time
	// preopenDir is the host directory mapped to wasi preopen fd 3.
	// Defaults to "/" (i.e. no rewriting) — the legacy behaviour. Tests
	// can set this via SetPreopenDir to scope filesystem ops to a
	// temporary directory.
	preopenDir string
	// fsHook, when non-nil, is consulted before every filesystem access
	// (Path_open, Path_create_directory, Path_unlink_file). It receives
	// the guest-supplied path (relative to the preopen, e.g. "a.txt" or
	// "sub/a.txt") and whether the access is a write. Returning false
	// denies the operation, which surfaces to the guest as EACCES. This
	// is the host-controlled whitelist hook: the policy itself lives in
	// the embedding application, OUTSIDE the generated runtime.
	fsHook func(path string, write bool) bool
	// netHook, when non-nil, is consulted before every socket operation
	// (Sock_accept, Sock_recv, Sock_send). op is "accept"/"recv"/"send".
	// Returning false denies the operation (EACCES). The same
	// host-controlled-whitelist intent as fsHook, for the network surface.
	netHook func(op string) bool
	// dialHook, when non-nil, is consulted before an OUTBOUND connect
	// (Sock_connect) with the resolved network ("tcp"), the HOST the guest
	// resolved to reach this address (from the preceding Sock_getaddrinfo, or ""
	// if the guest dialed a literal IP), the dotted-quad IP, and the port.
	// Returning false denies the connection (EACCES). Passing the host lets the
	// policy match host+port jointly, which a port-scoped rule needs — the IP
	// alone cannot be tied back to the rule that authorized the name.
	dialHook func(network, host, ip string, port int) bool
	// resolveHook, when non-nil, is consulted before a name lookup
	// (Sock_getaddrinfo) with the requested host. Returning false denies the
	// resolution (the guest sees a gaierror). This is the hostname-level
	// whitelist control point (e.g. block "example.com" by name).
	resolveHook func(host string) bool
	// resolvedHosts maps a resolved dotted-quad IP back to the host name the
	// guest looked it up under (populated by Sock_getaddrinfo, read by
	// Sock_connect), so the dial hook can be given the host. Guarded by mu.
	resolvedHosts map[string]string
	// fsys is the filesystem backend every guest path operation is routed
	// through. Defaults to an osFS scoped to preopenDir (the host filesystem);
	// SetFS swaps in an alternative (e.g. an in-memory FS) so each module can
	// see a private, arbitrary filesystem.
	fsys FS
	// procs tracks host processes spawned via Proc_spawn, keyed by the pid
	// handed back to the guest. nextPID is the handle counter (kept distinct
	// from real OS pids — the guest only ever sees these tokens).
	procs   map[int32]*wasiProc
	nextPID int32
	// execHook, when non-nil, gates every Proc_spawn with the resolved
	// executable path and argv; returning false denies the spawn (EACCES).
	// This is the outbound-process whitelist control point — the analogue of
	// dialHook for sockets. Spawning runs a HOST binary, so a sandbox that
	// enables host processes should always install this.
	execHook func(path string, argv []string) bool
}

// wasiProc is a host process spawned by Proc_spawn. A background goroutine
// Waits on the command and publishes the encoded POSIX status, so Proc_wait
// can support both the blocking (options 0) and non-blocking (WNOHANG) forms
// without holding the WasiStubs lock across the child's lifetime.
type wasiProc struct {
	cmd    *exec.Cmd
	done   chan struct{}
	status int32 // POSIX wait status, valid once done is closed
}

// DefaultWASI returns a WasiStubs configured for typical CLI use: real
// stdio, os.Args, os.Environ(), wall + monotonic clocks. Consumers who
// want a sandboxed setup should construct their own WasiStubs (or any
// Wasi_snapshot_preview1Imports implementation) and pass it to
// NewWithWASI.
func DefaultWASI() *WasiStubs {
	return &WasiStubs{
		stdin:      os.Stdin,
		stdout:     os.Stdout,
		stderr:     os.Stderr,
		fdTable:    map[int32]*wasiOpen{},
		nextFD:     4,
		args:       os.Args,
		env:        os.Environ(),
		monoStart:  time.Now(),
		preopenDir: "/",
		fsys:       osFS{root: "/"},
		procs:      map[int32]*wasiProc{},
		nextPID:    1000,
	}
}

// SetPreopenDir scopes the default (os-backed) filesystem to a host directory.
// Empty string restores the default ("/"), i.e. no rewriting. Tests use this
// to run filesystem syscalls against t.TempDir(). Has no effect once SetFS has
// installed a non-os backend.
func (w *WasiStubs) SetPreopenDir(dir string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if dir == "" {
		dir = "/"
	}
	w.preopenDir = dir
	w.fsys = osFS{root: dir}
}

// SetFS installs a custom filesystem backend. Every guest path operation
// (open, stat, mkdir, readdir, read, write, ...) is then routed to fsys, so a
// caller can give a module a private, arbitrary filesystem — for example an
// in-memory FS so writes never touch disk and are invisible to other modules.
// Pass nil to restore the default os-backed filesystem.
func (w *WasiStubs) SetFS(fsys FS) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if fsys == nil {
		fsys = osFS{root: w.preopenDir}
	}
	w.fsys = fsys
}

// SetFSAccessHook installs a host-controlled filesystem access policy.
// hook is called with the guest path (relative to the preopen) and a
// write flag before each open/create/unlink; returning false denies the
// operation (the guest sees EACCES). Pass nil to clear the policy
// (unrestricted, the default). The hook runs without w.mu held, so it
// may itself call back into the host freely.
func (w *WasiStubs) SetFSAccessHook(hook func(path string, write bool) bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.fsHook = hook
}

// SetNetAccessHook installs a host-controlled network access policy.
// hook is called with the operation name ("accept"/"recv"/"send")
// before each socket operation; returning false denies it (EACCES).
// Pass nil to clear (unrestricted, the default).
//
// NOTE: WASI preview1 has no outbound connect or name resolution, so a
// guest cannot initiate connections regardless of this hook; it governs
// the accept/recv/send surface that preview1 does expose (host-preopened
// listening sockets). Full outbound control requires a host connect
// import, which this runtime does not yet provide.
func (w *WasiStubs) SetNetAccessHook(hook func(op string) bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.netHook = hook
}

// SetDialHook installs a host-controlled OUTBOUND-connection policy. hook is
// called with ("tcp", host, dotted-quad-IP, port) before each Sock_connect,
// where host is the name the guest resolved to reach the IP (from the preceding
// Sock_getaddrinfo) or "" for a literal-IP dial; returning false denies the
// connection (the guest sees a connect EACCES). Pass nil to clear (all outbound
// allowed, the default once outbound is wired).
func (w *WasiStubs) SetDialHook(hook func(network, host, ip string, port int) bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.dialHook = hook
}

// SetResolveHook installs a host-controlled name-resolution policy. hook is
// called with the host being resolved (Sock_getaddrinfo) before the lookup;
// returning false denies it (the guest sees a name-resolution error). Pass nil
// to clear (all lookups allowed). This is where a hostname whitelist such as
// "block example.com" is enforced.
func (w *WasiStubs) SetResolveHook(hook func(host string) bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.resolveHook = hook
}

// SetExecHook installs the process-spawn whitelist consulted by Proc_spawn
// with the executable path and full argv. Returning false denies the spawn
// (the guest's posix_spawn sees EACCES). Spawning runs a HOST binary, so a
// sandbox enabling host processes should always set this.
func (w *WasiStubs) SetExecHook(hook func(path string, argv []string) bool) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.execHook = hook
}

// readCStr reads a NUL-terminated C string at ptr from linear memory. A nil
// ptr (0) yields "". ok is false on an out-of-bounds or unterminated read.
func (w *WasiStubs) readCStr(m *Module, ptr int32) (s string, ok bool) {
	if ptr == 0 {
		return "", true
	}
	mem := m.memory
	lo := uint64(uint32(ptr))
	if lo > uint64(len(mem)) {
		return "", false
	}
	rest := mem[lo:]
	for i := 0; i < len(rest); i++ {
		if rest[i] == 0 {
			return string(rest[:i]), true
		}
	}
	return "", false
}

// readCStrArray reads a NULL-terminated array of C-string pointers (a char**)
// at ptr. A nil ptr (0) yields a nil slice. ok is false on a bad read.
func (w *WasiStubs) readCStrArray(m *Module, ptr int32) (out []string, ok bool) {
	if ptr == 0 {
		return nil, true
	}
	for off := ptr; ; off += 4 {
		b := w.memSlice(m, off, 4)
		if b == nil {
			return nil, false
		}
		p := int32(binary.LittleEndian.Uint32(b))
		if p == 0 {
			break
		}
		s, sok := w.readCStr(m, p)
		if !sok {
			return nil, false
		}
		out = append(out, s)
	}
	return out, true
}

// Proc_spawn is a NON-STANDARD host import (module wasi_snapshot_preview1,
// name "proc_spawn") backing the bridge's posix_spawn(). It spawns a HOST
// process: path is the executable, argv/envp are NUL-terminated char** in
// linear memory. The child inherits the interpreter's stdin/stdout/stderr.
// The new pid token is written at pidOutPtr. Returns 0 or a negative errno.
//
// Only stdio inheritance is supported today (no fd remapping / pipes), which
// covers subprocess.run/call with default streams; capture_output via host
// pipes is a follow-up.
func (w *WasiStubs) Proc_spawn(m *Module, pathPtr, argvPtr, envpPtr, stdinFd, stdoutFd, stderrFd, cwdPtr, pidOutPtr int32) int32 {
	path, ok := w.readCStr(m, pathPtr)
	if !ok || path == "" {
		return -_wasiEINVAL
	}
	argv, ok := w.readCStrArray(m, argvPtr)
	if !ok {
		return -_wasiEFAULT
	}
	env, ok := w.readCStrArray(m, envpPtr)
	if !ok {
		return -_wasiEFAULT
	}

	cwd, ok := w.readCStr(m, cwdPtr)
	if !ok {
		return -_wasiEFAULT
	}
	out := w.memSlice(m, pidOutPtr, 4)
	if out == nil {
		return -_wasiEFAULT
	}

	w.mu.Lock()
	hook := w.execHook
	cin := w.childReaderLocked(stdinFd)
	cout := w.childWriterLocked(stdoutFd, w.stdout)
	cerr := w.childWriterLocked(stderrFd, w.stderr)
	w.mu.Unlock()
	if hook != nil && !hook(path, argv) {
		return -_wasiEACCES
	}

	cmd := exec.Command(path)
	if len(argv) > 0 {
		cmd.Args = argv
	} else {
		cmd.Args = []string{path}
	}
	if cwd != "" {
		cmd.Dir = cwd
	}

	cmd.Env = env
	if cmd.Env == nil {
		cmd.Env = []string{}
	}

	cmd.Stdin, cmd.Stdout, cmd.Stderr = cin, cout, cerr
	if err := cmd.Start(); err != nil {
		return -mapExecError(err)
	}

	proc := &wasiProc{cmd: cmd, done: make(chan struct{})}
	go func() {
		werr := cmd.Wait()
		proc.status = encodeWaitStatus(cmd.ProcessState)
		// A non-zero exit or signal surfaces as *exec.ExitError and is the
		// normal path (status already encoded from ProcessState above). Any
		// OTHER error means the wait itself failed; report a 127 exit.
		var exitErr *exec.ExitError
		if werr != nil && !errors.As(werr, &exitErr) {
			proc.status = int32(127) << 8
		}
		close(proc.done)
	}()

	w.mu.Lock()
	if w.procs == nil {
		w.procs = map[int32]*wasiProc{}
	}
	if w.nextPID == 0 {
		w.nextPID = 1000
	}
	pid := w.nextPID
	w.nextPID++
	w.procs[pid] = proc
	w.mu.Unlock()

	binary.LittleEndian.PutUint32(out, uint32(pid))
	return _wasiESUCCESS
}

// childReaderLocked resolves a child stdin source fd. A guest fd whose
// table entry carries a real file (a pipe end, or a guest stdio fd the
// program re-opened onto a file) is used directly; everything else —
// including -1 and an unredirected fd 0 — inherits the interpreter's
// stdin. Caller holds w.mu.
func (w *WasiStubs) childReaderLocked(fd int32) io.Reader {
	if fd >= 0 {
		if op := w.fdTable[fd]; op != nil && op.f != nil {
			return op.f
		}
	}
	return w.stdin
}

// childWriterLocked is the stdout/stderr counterpart of childReaderLocked.
func (w *WasiStubs) childWriterLocked(fd int32, deflt io.Writer) io.Writer {
	if fd >= 0 {
		if op := w.fdTable[fd]; op != nil && op.f != nil {
			return op.f
		}
	}
	return deflt
}

// Pipe is a NON-STANDARD host import (module wasi_snapshot_preview1, name
// "pipe") backing the bridge's pipe()/pipe2(). It creates a host OS pipe and
// registers both ends as guest fds, writing [readFd, writeFd] (two i32) at
// fdsOutPtr. The guest reads the read end via Fd_read; the write end is given
// to a child as its stdout/stderr via Proc_spawn, so subprocess.run can
// capture output. Returns 0 or a negative errno.
func (w *WasiStubs) Pipe(m *Module, fdsOutPtr int32) int32 {
	out := w.memSlice(m, fdsOutPtr, 8)
	if out == nil {
		return -_wasiEFAULT
	}
	r, wr, err := os.Pipe()
	if err != nil {
		return -mapOSError(err)
	}
	w.mu.Lock()
	if w.fdTable == nil {
		w.fdTable = map[int32]*wasiOpen{}
	}
	if w.nextFD < 4 {
		w.nextFD = 4
	}
	rfd := w.nextFD
	w.nextFD++
	wfd := w.nextFD
	w.nextFD++
	w.fdTable[rfd] = &wasiOpen{f: r}
	w.fdTable[wfd] = &wasiOpen{f: wr}
	w.mu.Unlock()
	binary.LittleEndian.PutUint32(out[0:], uint32(rfd))
	binary.LittleEndian.PutUint32(out[4:], uint32(wfd))
	return _wasiESUCCESS
}

// Proc_wait is a NON-STANDARD host import (name "proc_wait") backing the
// bridge's waitpid(). It waits for the process token pid and writes the POSIX
// wait status at statusOutPtr. options is the waitpid() options mask; bit 0
// (WNOHANG) makes it return without blocking when the child is still running
// (the guest sees the documented "0 means no child ready" result, signalled
// by writing pid 0 — encoded by returning EAGAIN). Returns 0, or a negative
// errno (ECHILD for an unknown pid).
func (w *WasiStubs) Proc_wait(m *Module, pid, options, statusOutPtr int32) int32 {
	out := w.memSlice(m, statusOutPtr, 4)
	if out == nil {
		return -_wasiEFAULT
	}
	w.mu.Lock()
	proc := w.procs[pid]
	w.mu.Unlock()
	if proc == nil {
		return -_wasiECHILD
	}
	const wnohang = 1
	if options&wnohang != 0 {
		select {
		case <-proc.done:
		default:

			return -_wasiEAGAIN
		}
	} else {
		<-proc.done
	}
	w.mu.Lock()
	delete(w.procs, pid)
	w.mu.Unlock()
	binary.LittleEndian.PutUint32(out, uint32(proc.status))
	return _wasiESUCCESS
}

// encodeWaitStatus turns a Go ProcessState into a POSIX wait status int (the
// raw value os.waitstatus_to_exitcode decodes): a normal exit N becomes
// (N&0xff)<<8 (WIFEXITED); a signal becomes the low-7-bits signal number
// (WIFSIGNALED).
func encodeWaitStatus(st *os.ProcessState) int32 {
	if st == nil {
		return 0
	}
	if ws, ok := st.Sys().(syscall.WaitStatus); ok {
		if ws.Signaled() {
			return int32(ws.Signal()) & 0x7f
		}
		return int32(ws.ExitStatus()&0xff) << 8
	}
	code := st.ExitCode()
	if code < 0 {
		code = 127
	}
	return int32(code&0xff) << 8
}

// mapExecError maps an exec.Command Start() failure to a wasi errno.
func mapExecError(err error) int32 {
	switch {
	case errors.Is(err, exec.ErrNotFound), errors.Is(err, fs.ErrNotExist):
		return _wasiENOENT
	case errors.Is(err, fs.ErrPermission):
		return _wasiEACCES
	default:
		return _wasiENOENT
	}
}

// Sock_getaddrinfo is a NON-STANDARD host import (module wasi_snapshot_preview1,
// name "sock_getaddrinfo") backing the bridge's getaddrinfo(). It reads the
// host string at (nodePtr,nodeLen), consults the resolve whitelist, resolves it
// to an IPv4 address via Go's resolver (numeric IPs pass through), and writes
// the 4-byte network-order address at outPtr. Returns 0 on success or a
// negative POSIX-ish errno (the bridge maps it to an EAI_* code).
func (w *WasiStubs) Sock_getaddrinfo(m *Module, nodePtr, nodeLen, outPtr int32) int32 {
	host := ""
	if nodeLen > 0 {
		b := w.memSlice(m, nodePtr, nodeLen)
		if b == nil {
			return -_wasiEFAULT
		}
		host = string(b)
	}
	out := w.memSlice(m, outPtr, 4)
	if out == nil {
		return -_wasiEFAULT
	}
	if host == "" {
		binary.LittleEndian.PutUint32(out, 0)
		return _wasiESUCCESS
	}
	w.mu.Lock()
	hook := w.resolveHook
	w.mu.Unlock()
	if hook != nil && !hook(host) {
		return -_wasiEACCES
	}

	if ip := net.ParseIP(host); ip != nil {
		if v4 := ip.To4(); v4 != nil {
			out[0], out[1], out[2], out[3] = v4[0], v4[1], v4[2], v4[3]
			w.recordResolvedHost(v4, host)
			return _wasiESUCCESS
		}
		return -_wasiEAFNOSUPPORT
	}
	ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip4", host)
	if err != nil || len(ips) == 0 {
		return -_wasiENOENT
	}
	v4 := ips[0].To4()
	if v4 == nil {
		return -_wasiEAFNOSUPPORT
	}
	out[0], out[1], out[2], out[3] = v4[0], v4[1], v4[2], v4[3]
	w.recordResolvedHost(v4, host)
	return _wasiESUCCESS
}

// recordResolvedHost remembers that host resolved to v4, so a later Sock_connect
// to that IP can hand the dial hook the host name it was looked up under.
func (w *WasiStubs) recordResolvedHost(v4 net.IP, host string) {
	ip := fmt.Sprintf("%d.%d.%d.%d", v4[0], v4[1], v4[2], v4[3])
	w.mu.Lock()
	if w.resolvedHosts == nil {
		w.resolvedHosts = make(map[string]string)
	}
	w.resolvedHosts[ip] = host
	w.mu.Unlock()
}

// SetEnv overrides the environment the guest sees via environ_get /
// environ_sizes_get. By default DefaultWASI leaks the host process
// os.Environ(); a sandboxed embedding should call SetEnv with an
// explicit (possibly empty) slice of "KEY=VALUE" strings.
func (w *WasiStubs) SetEnv(env []string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.env = append([]string(nil), env...)
}

// SetArgs overrides os.Args as seen by the guest (argv). Mirrors SetEnv.
func (w *WasiStubs) SetArgs(args []string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.args = append([]string(nil), args...)
}

// SetStdin redirects guest fd 0 to r. A nil r leaves the current source.
// Use this to feed input() / sys.stdin from an in-process io.Reader instead
// of the host process stdin.
func (w *WasiStubs) SetStdin(r io.Reader) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if r != nil {
		w.stdin = r
	}
}

// SetStdout redirects guest fd 1 to wr. A nil wr leaves the current sink.
func (w *WasiStubs) SetStdout(wr io.Writer) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if wr != nil {
		w.stdout = wr
	}
}

// SetStderr redirects guest fd 2 to wr. A nil wr leaves the current sink.
func (w *WasiStubs) SetStderr(wr io.Writer) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if wr != nil {
		w.stderr = wr
	}
}

// checkFS consults the FS policy hook (if any). Returns true when the
// access is permitted. Callers must NOT hold w.mu.
func (w *WasiStubs) checkFS(path string, write bool) bool {
	w.mu.Lock()
	hook := w.fsHook
	w.mu.Unlock()
	if hook == nil {
		return true
	}
	return hook(path, write)
}

// checkNet consults the network policy hook (if any). Returns true when
// the operation is permitted. Callers must NOT hold w.mu.
func (w *WasiStubs) checkNet(op string) bool {
	w.mu.Lock()
	hook := w.netHook
	w.mu.Unlock()
	if hook == nil {
		return true
	}
	return hook(op)
}

// memSlice returns m.memory[off : off+n]. Callers must hold any locks
// they need on the wasm side; WasiStubs.mu is independent. Returns an
// empty slice on out-of-range (the wasi function should then return
// EFAULT / EINVAL).
func (w *WasiStubs) memSlice(m *Module, off, n int32) []byte {
	mem := m.memory
	lo := uint64(uint32(off))
	hi := lo + uint64(uint32(n))
	if hi > uint64(len(mem)) {
		return nil
	}
	return mem[lo:hi]
}

// errno values used below (subset; see wasi-libc errno.h).
const (
	_wasiESUCCESS     int32 = 0
	_wasiE2BIG        int32 = 1
	_wasiEACCES       int32 = 2
	_wasiEAFNOSUPPORT int32 = 5
	_wasiEAGAIN       int32 = 6
	_wasiEBADF        int32 = 8
	_wasiECHILD       int32 = 12
	_wasiECONNREFUSED int32 = 14
	_wasiEISCONN      int32 = 33
	_wasiEBUSY        int32 = 10
	_wasiEEXIST       int32 = 20
	_wasiEFAULT       int32 = 21
	_wasiEINVAL       int32 = 28
	_wasiEIO          int32 = 29
	_wasiEISDIR       int32 = 31
	_wasiENOENT       int32 = 44
	_wasiENOTDIR      int32 = 54
	_wasiENOTSOCK     int32 = 57
	_wasiENOTSUP      int32 = 58
	_wasiENOSYS       int32 = 52
	_wasiEPERM        int32 = 63
	_wasiEPIPE        int32 = 64
)

// mapOSError turns an os/filesystem error into a wasi errno. Used by the
// path-based syscalls so any os.PathError surfaces as the appropriate
// guest-visible code instead of a coarse EIO.
func mapOSError(err error) int32 {
	if err == nil {
		return _wasiESUCCESS
	}
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return _wasiENOENT
	case errors.Is(err, fs.ErrExist):
		return _wasiEEXIST
	case errors.Is(err, fs.ErrPermission):
		return _wasiEACCES
	case errors.Is(err, syscall.ENOTDIR):
		return _wasiENOTDIR
	case errors.Is(err, syscall.EISDIR):
		return _wasiEISDIR
	case errors.Is(err, syscall.EINVAL):
		return _wasiEINVAL
	case errors.Is(err, syscall.EBADF):
		return _wasiEBADF
	case errors.Is(err, syscall.EAGAIN):
		return _wasiEAGAIN
	case errors.Is(err, syscall.EPIPE):
		return _wasiEPIPE
	}
	return _wasiEIO
}

// totalBytes sums len(s)+1 over s in a uint64 and reports whether the
// total fits in an int32 (i.e. is representable as a wasm-side i32
// length). Callers route the result through memSlice and an OOB on a
// pathologically long arg list surfaces as EFAULT to the guest rather
// than a host-side panic via a wrapped-int32 length.
func totalBytesPlusNul(ss []string) (int32, bool) {
	var total uint64
	for _, s := range ss {
		total += uint64(len(s)) + 1
		if total > 0x7fffffff {
			return 0, false
		}
	}
	return int32(total), true
}

func (w *WasiStubs) Args_get(m *Module, argv, argvBuf int32) int32 {
	w.mu.Lock()
	defer w.mu.Unlock()

	argvBytes64 := uint64(len(w.args)) * 4
	if argvBytes64 > 0x7fffffff {
		return _wasiEFAULT
	}
	argvSlice := w.memSlice(m, argv, int32(argvBytes64))
	if argvSlice == nil {
		return _wasiEFAULT
	}
	total, ok := totalBytesPlusNul(w.args)
	if !ok {
		return _wasiEFAULT
	}
	argvBufSlice := w.memSlice(m, argvBuf, total)
	if argvBufSlice == nil {
		return _wasiEFAULT
	}
	bufOff := uint32(0)
	for i, a := range w.args {
		binary.LittleEndian.PutUint32(argvSlice[i*4:], uint32(argvBuf)+bufOff)
		n := copy(argvBufSlice[bufOff:], a)
		if n < len(a) {
			return _wasiEFAULT
		}
		bufOff += uint32(n)
		argvBufSlice[bufOff] = 0
		bufOff++
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Args_sizes_get(m *Module, argcPtr, argvBufLenPtr int32) int32 {
	w.mu.Lock()
	defer w.mu.Unlock()
	argcSlice := w.memSlice(m, argcPtr, 4)
	bufLenSlice := w.memSlice(m, argvBufLenPtr, 4)
	if argcSlice == nil || bufLenSlice == nil {
		return _wasiEFAULT
	}
	total, ok := totalBytesPlusNul(w.args)
	if !ok {
		return _wasiEFAULT
	}
	binary.LittleEndian.PutUint32(argcSlice, uint32(len(w.args)))
	binary.LittleEndian.PutUint32(bufLenSlice, uint32(total))
	return _wasiESUCCESS
}

func (w *WasiStubs) Environ_get(m *Module, envv, envBuf int32) int32 {
	w.mu.Lock()
	defer w.mu.Unlock()
	envvBytes64 := uint64(len(w.env)) * 4
	if envvBytes64 > 0x7fffffff {
		return _wasiEFAULT
	}
	envvSlice := w.memSlice(m, envv, int32(envvBytes64))
	if envvSlice == nil {
		return _wasiEFAULT
	}
	total, ok := totalBytesPlusNul(w.env)
	if !ok {
		return _wasiEFAULT
	}
	envBufSlice := w.memSlice(m, envBuf, total)
	if envBufSlice == nil {
		return _wasiEFAULT
	}
	bufOff := uint32(0)
	for i, e := range w.env {
		binary.LittleEndian.PutUint32(envvSlice[i*4:], uint32(envBuf)+bufOff)
		n := copy(envBufSlice[bufOff:], e)
		if n < len(e) {
			return _wasiEFAULT
		}
		bufOff += uint32(n)
		envBufSlice[bufOff] = 0
		bufOff++
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Environ_sizes_get(m *Module, envcPtr, envBufLenPtr int32) int32 {
	w.mu.Lock()
	defer w.mu.Unlock()
	envcSlice := w.memSlice(m, envcPtr, 4)
	bufLenSlice := w.memSlice(m, envBufLenPtr, 4)
	if envcSlice == nil || bufLenSlice == nil {
		return _wasiEFAULT
	}
	total, ok := totalBytesPlusNul(w.env)
	if !ok {
		return _wasiEFAULT
	}
	binary.LittleEndian.PutUint32(envcSlice, uint32(len(w.env)))
	binary.LittleEndian.PutUint32(bufLenSlice, uint32(total))
	return _wasiESUCCESS
}

func (w *WasiStubs) Clock_res_get(m *Module, clockID int32, resPtr int32) int32 {

	out := w.memSlice(m, resPtr, 8)
	if out == nil {
		return _wasiEFAULT
	}
	binary.LittleEndian.PutUint64(out, 1)
	return _wasiESUCCESS
}

func (w *WasiStubs) Clock_time_get(m *Module, clockID int32, precision int64, timePtr int32) int32 {
	out := w.memSlice(m, timePtr, 8)
	if out == nil {
		return _wasiEFAULT
	}
	nanos, errno := w.clockNanos(clockID)
	if errno != _wasiESUCCESS {
		return errno
	}
	binary.LittleEndian.PutUint64(out, nanos)
	return _wasiESUCCESS
}

// clockNanos is the layout-independent body of clock_time_get, shared
// by the wasm32 and wasm64 bindings.
func (w *WasiStubs) clockNanos(clockID int32) (uint64, int32) {
	switch clockID {
	case 0:
		return uint64(time.Now().UnixNano()), _wasiESUCCESS
	case 1:
		w.mu.Lock()
		nanos := uint64(time.Since(w.monoStart).Nanoseconds())
		w.mu.Unlock()
		return nanos, _wasiESUCCESS
	default:
		return 0, _wasiEINVAL
	}
}

// closeWasiOpen releases every underlying handle held by op and
// joins any Close errors so callers can map them to a wasi errno
// instead of silently dropping the failure.
func closeWasiOpen(op *wasiOpen) error {
	if op.refs > 0 {

		op.refs--
		return nil
	}
	var err error
	if op.f != nil {
		err = errors.Join(err, op.f.Close())
	}
	if op.conn != nil {
		err = errors.Join(err, op.conn.Close())
	}
	if op.listener != nil {
		err = errors.Join(err, op.listener.Close())
	}
	return err
}

func (w *WasiStubs) Fd_close(m *Module, fd int32) int32 {
	w.mu.Lock()
	defer w.mu.Unlock()
	op := w.fdTable[fd]
	if op == nil {
		return _wasiEBADF
	}
	closeErr := closeWasiOpen(op)
	delete(w.fdTable, fd)
	if closeErr != nil {
		return mapOSError(closeErr)
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_fdstat_get(m *Module, fd, ptr int32) int32 {

	out := w.memSlice(m, ptr, 24)
	if out == nil {
		return _wasiEFAULT
	}
	return w.fdstatFill(fd, out)
}

// fdstatFill writes the 24-byte fdstat for fd into out — the shared
// body of the 32- and 64-bit Fd_fdstat_get bindings (the struct holds
// no pointers, so the layout is width-independent).
func (w *WasiStubs) fdstatFill(fd int32, out []byte) int32 {
	w.mu.Lock()
	defer w.mu.Unlock()
	var ftype byte = 4 // regular file
	var fdflags uint16
	if fd >= 0 && fd <= 2 {
		ftype = 2
	} else if op := w.fdTable[fd]; op != nil {
		if op.isDir {
			ftype = 3
		} else if op.conn != nil {
			ftype = 6
		} else if op.listener != nil {
			ftype = 6
		}
		fdflags = uint16(op.fdflags)
	} else if fd == 3 {
		ftype = 3
	} else if fd >= 4 {
		return _wasiEBADF
	}
	out[0] = ftype
	out[1] = 0
	binary.LittleEndian.PutUint16(out[2:], fdflags)

	binary.LittleEndian.PutUint64(out[8:], ^uint64(0))
	binary.LittleEndian.PutUint64(out[16:], ^uint64(0))
	return _wasiESUCCESS
}

// Fd_fdstat_set_flags maps WASI fdflags to OS file-status flags via the
// per-platform Fcntl wrapper. The flags are also cached on the wasiOpen
// so a subsequent Fd_fdstat_get reflects what the guest set. Stdio fds
// store the requested flags but otherwise no-op; sockets/listeners take
// only the cache update because Go's net layer manages blocking mode
// internally.
func (w *WasiStubs) Fd_fdstat_set_flags(m *Module, fd, flags int32) int32 {
	w.mu.Lock()
	op := w.fdTable[fd]
	if op == nil && fd > 2 {
		w.mu.Unlock()
		return _wasiEBADF
	}
	if op != nil {
		op.fdflags = flags
	}
	w.mu.Unlock()

	_ = op
	_ = flags
	return _wasiESUCCESS
}

// Fd_fdstat_set_rights stores the requested rights on the wasiOpen but
// does not enforce them — the host process is the trust boundary. WASI
// programs that succeed with maximal rights (per Fd_fdstat_get) get the
// same ESUCCESS here.
func (w *WasiStubs) Fd_fdstat_set_rights(m *Module, fd int32, rightsBase, rightsInherit int64) int32 {
	w.mu.Lock()
	defer w.mu.Unlock()
	if fd >= 0 && fd <= 2 {
		return _wasiESUCCESS
	}
	if w.fdTable[fd] == nil {
		return _wasiEBADF
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_filestat_get(m *Module, fd, ptr int32) int32 {

	out := w.memSlice(m, ptr, 64)
	if out == nil {
		return _wasiEFAULT
	}
	w.mu.Lock()
	op := w.fdTable[fd]
	w.mu.Unlock()
	if op == nil || op.f == nil {

		for i := range out {
			out[i] = 0
		}

		switch fd {
		case 0, 1, 2:
			out[16] = 2
		case 3:
			out[16] = 3
		}
		return _wasiESUCCESS
	}
	st, err := op.f.Stat()
	if err != nil {
		return mapOSError(err)
	}
	writeFilestat(out, st)
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_filestat_set_size(m *Module, fd int32, size int64) int32 {
	w.mu.Lock()
	op := w.fdTable[fd]
	w.mu.Unlock()
	if op == nil || op.f == nil {
		return _wasiEBADF
	}
	if err := op.f.Truncate(size); err != nil {
		return mapOSError(err)
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_filestat_set_times(m *Module, fd int32, atim, mtim int64, fstFlags int32) int32 {

	w.mu.Lock()
	op := w.fdTable[fd]
	fsys := w.fsys
	w.mu.Unlock()
	if op == nil || op.f == nil {
		return _wasiEBADF
	}
	atime, mtime, err := resolveFiletimes(uint64(atim), uint64(mtim), fstFlags, op.f)
	if err != nil {
		return mapOSError(err)
	}

	if cf, ok := fsys.(chtimesFS); ok {
		if err := cf.Chtimes(op.path, atime, mtime); err != nil {
			return mapOSError(err)
		}
	}
	return _wasiESUCCESS
}

// combine64 reconstructs an unsigned 64-bit time value from a pair of
// 32-bit args. WASI signature uses two i32s for the nanosecond timestamp
// in fd_filestat_set_times.
func combine64(hi, lo int32) uint64 {
	return (uint64(uint32(hi)) << 32) | uint64(uint32(lo))
}

// resolveFiletimes decides the (atime, mtime) pair to apply given a
// fstFlags bitmask. Bits 0x2 (ATIME_NOW) and 0x8 (MTIME_NOW) override the
// explicit values with time.Now(). Unset ATIME/MTIME bits keep the
// existing on-disk time, so f.Stat must succeed when those bits are
// unset; the error is returned so the caller can surface it as a wasi
// errno rather than silently writing epoch.
func resolveFiletimes(atimNs, mtimNs uint64, fstFlags int32, f File) (time.Time, time.Time, error) {
	now := time.Now()
	var atime, mtime time.Time

	needCurrent := fstFlags&(0x1|0x2) == 0 || fstFlags&(0x4|0x8) == 0
	if needCurrent {
		st, err := f.Stat()
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
		atime = st.ModTime()
		mtime = st.ModTime()
	}
	if fstFlags&0x1 != 0 {
		atime = time.Unix(0, int64(atimNs))
	}
	if fstFlags&0x2 != 0 {
		atime = now
	}
	if fstFlags&0x4 != 0 {
		mtime = time.Unix(0, int64(mtimNs))
	}
	if fstFlags&0x8 != 0 {
		mtime = now
	}
	return atime, mtime, nil
}

func (w *WasiStubs) Fd_prestat_get(m *Module, fd, ptr int32) int32 {

	if fd != 3 {
		return _wasiEBADF
	}

	out := w.memSlice(m, ptr, 8)
	if out == nil {
		return _wasiEFAULT
	}
	out[0] = 0
	binary.LittleEndian.PutUint32(out[4:], 1)
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_prestat_dir_name(m *Module, fd, buf, buflen int32) int32 {
	if fd != 3 {
		return _wasiEBADF
	}
	if buflen < 1 {
		return _wasiESUCCESS
	}
	out := w.memSlice(m, buf, buflen)
	if out == nil {
		return _wasiEFAULT
	}
	out[0] = '/'
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_read(m *Module, fd, iovs, iovsLen, nreadPtr int32) int32 {
	w.mu.Lock()
	src, op := w.fdSrcLocked(fd)
	w.mu.Unlock()
	if src == nil {
		return _wasiEBADF
	}
	bufs, ok := w.iovecSlices(m, iovs, iovsLen)
	nreadSlice := w.memSlice(m, nreadPtr, 4)
	if !ok || nreadSlice == nil {
		return _wasiEFAULT
	}
	_ = op
	binary.LittleEndian.PutUint32(nreadSlice, uint32(readVec(src, bufs)))
	return _wasiESUCCESS
}

// iovecSlices resolves a wasm32 ciovec/iovec array ({u32 ptr, u32 len}
// entries at iovs) into the backing memory windows. Every entry is
// validated before any I/O happens, so a bad iovec faults the whole
// call instead of after a partial transfer.
func (w *WasiStubs) iovecSlices(m *Module, iovs, iovsLen int32) ([][]byte, bool) {

	iovBytes := uint64(uint32(iovsLen)) * 8
	if iovBytes > 0x7fffffff {
		return nil, false
	}
	iovecs := w.memSlice(m, iovs, int32(iovBytes))
	if iovecs == nil {
		return nil, false
	}
	bufs := make([][]byte, 0, iovsLen)
	for i := int32(0); i < iovsLen; i++ {
		bufPtr := binary.LittleEndian.Uint32(iovecs[i*8:])
		bufLen := binary.LittleEndian.Uint32(iovecs[i*8+4:])
		buf := w.memSlice(m, int32(bufPtr), int32(bufLen))
		if buf == nil {
			return nil, false
		}
		bufs = append(bufs, buf)
	}
	return bufs, true
}

// readVec fills bufs from src in order, stopping at the first error
// (EOF included) or short read; returns the bytes read. Shared by the
// wasm32 and wasm64 fd_read bindings — only the iovec layout differs.
func readVec(src io.Reader, bufs [][]byte) uint64 {
	var total uint64
	for _, buf := range bufs {
		n, err := src.Read(buf)
		total += uint64(n)
		if err != nil || n < len(buf) {
			break
		}
	}
	return total
}

// writeVec drains bufs into dst in order, stopping at the first failed
// write; returns the bytes written. Shared like readVec.
func writeVec(dst io.Writer, bufs [][]byte) uint64 {
	var total uint64
	for _, buf := range bufs {
		n, err := dst.Write(buf)
		total += uint64(n)
		if err != nil {
			break
		}
	}
	return total
}

// fdSrcLocked returns the io.Reader for fd and (when applicable) the
// wasiOpen it came from, or nil if fd is invalid. Caller must hold w.mu.
func (w *WasiStubs) fdSrcLocked(fd int32) (io.Reader, *wasiOpen) {

	op := w.fdTable[fd]
	if op == nil {
		if fd == 0 {
			return w.stdin, nil
		}
		return nil, nil
	}
	if op.stdio == 1 {
		return w.stdin, op
	}
	if op.f != nil {
		return op.f, op
	}
	if op.conn != nil {
		return op.conn, op
	}
	return nil, op
}

// fdDstLocked returns the io.Writer for fd or nil if fd is invalid.
// Caller must hold w.mu.
func (w *WasiStubs) fdDstLocked(fd int32) (io.Writer, *wasiOpen) {
	op := w.fdTable[fd]
	if op == nil {
		switch fd {
		case 1:
			return w.stdout, nil
		case 2:
			return w.stderr, nil
		}
		return nil, nil
	}
	switch op.stdio {
	case 2:
		return w.stdout, op
	case 3:
		return w.stderr, op
	}
	if op.f != nil {
		return op.f, op
	}
	if op.conn != nil {
		return op.conn, op
	}
	return nil, op
}

func (w *WasiStubs) Fd_pread(m *Module, fd, iovs, iovsLen int32, offset int64, nreadPtr int32) int32 {
	w.mu.Lock()
	op := w.fdTable[fd]
	w.mu.Unlock()
	if op == nil || op.f == nil {
		return _wasiEBADF
	}
	iovBytes := uint64(uint32(iovsLen)) * 8
	if iovBytes > 0x7fffffff {
		return _wasiEFAULT
	}
	iovecs := w.memSlice(m, iovs, int32(iovBytes))
	nreadSlice := w.memSlice(m, nreadPtr, 4)
	if iovecs == nil || nreadSlice == nil {
		return _wasiEFAULT
	}
	var total uint32
	curOff := offset
	for i := int32(0); i < iovsLen; i++ {
		bufPtr := binary.LittleEndian.Uint32(iovecs[i*8:])
		bufLen := binary.LittleEndian.Uint32(iovecs[i*8+4:])
		buf := w.memSlice(m, int32(bufPtr), int32(bufLen))
		if buf == nil {
			return _wasiEFAULT
		}
		n, err := op.f.ReadAt(buf, curOff)
		total += uint32(n)
		curOff += int64(n)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			break
		}
		if n < int(bufLen) {
			break
		}
	}
	binary.LittleEndian.PutUint32(nreadSlice, total)
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_pwrite(m *Module, fd, iovs, iovsLen int32, offset int64, nwrittenPtr int32) int32 {
	w.mu.Lock()
	op := w.fdTable[fd]
	w.mu.Unlock()
	if op == nil || op.f == nil {
		return _wasiEBADF
	}
	iovBytes := uint64(uint32(iovsLen)) * 8
	if iovBytes > 0x7fffffff {
		return _wasiEFAULT
	}
	iovecs := w.memSlice(m, iovs, int32(iovBytes))
	nwSlice := w.memSlice(m, nwrittenPtr, 4)
	if iovecs == nil || nwSlice == nil {
		return _wasiEFAULT
	}
	var total uint32
	curOff := offset
	for i := int32(0); i < iovsLen; i++ {
		bufPtr := binary.LittleEndian.Uint32(iovecs[i*8:])
		bufLen := binary.LittleEndian.Uint32(iovecs[i*8+4:])
		buf := w.memSlice(m, int32(bufPtr), int32(bufLen))
		if buf == nil {
			return _wasiEFAULT
		}
		n, err := op.f.WriteAt(buf, curOff)
		total += uint32(n)
		curOff += int64(n)
		if err != nil {
			break
		}
	}
	binary.LittleEndian.PutUint32(nwSlice, total)
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_seek(m *Module, fd int32, offset int64, whence, newOffPtr int32) int32 {
	out := w.memSlice(m, newOffPtr, 8)
	if out == nil {
		return _wasiEFAULT
	}
	n, errno := w.fdSeek(fd, offset, int(whence))
	if errno != _wasiESUCCESS {
		return errno
	}
	binary.LittleEndian.PutUint64(out, uint64(n))
	return _wasiESUCCESS
}

// fdSeek is the layout-independent body of fd_seek, shared by the
// wasm32 and wasm64 bindings.
func (w *WasiStubs) fdSeek(fd int32, offset int64, whence int) (int64, int32) {
	w.mu.Lock()
	op := w.fdTable[fd]
	w.mu.Unlock()
	if op == nil || op.f == nil {
		return 0, _wasiEBADF
	}
	n, err := op.f.Seek(offset, whence)
	if err != nil {
		return 0, _wasiEINVAL
	}
	return n, _wasiESUCCESS
}

func (w *WasiStubs) Fd_tell(m *Module, fd, offsetPtr int32) int32 {
	out := w.memSlice(m, offsetPtr, 8)
	if out == nil {
		return _wasiEFAULT
	}
	w.mu.Lock()
	op := w.fdTable[fd]
	w.mu.Unlock()
	if op == nil || op.f == nil {
		return _wasiEBADF
	}
	n, err := op.f.Seek(0, 1)
	if err != nil {
		return _wasiEIO
	}
	binary.LittleEndian.PutUint64(out, uint64(n))
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_write(m *Module, fd, iovs, iovsLen, nwrittenPtr int32) int32 {
	w.mu.Lock()
	dst, _ := w.fdDstLocked(fd)
	w.mu.Unlock()
	bufs, ok := w.iovecSlices(m, iovs, iovsLen)
	nwrittenSlice := w.memSlice(m, nwrittenPtr, 4)
	if !ok || nwrittenSlice == nil {
		return _wasiEFAULT
	}
	if dst == nil {
		binary.LittleEndian.PutUint32(nwrittenSlice, 0)
		return _wasiEBADF
	}
	binary.LittleEndian.PutUint32(nwrittenSlice, uint32(writeVec(dst, bufs)))
	return _wasiESUCCESS
}
func (w *WasiStubs) Fd_sync(m *Module, fd int32) int32 {
	w.mu.Lock()
	op := w.fdTable[fd]
	w.mu.Unlock()
	if op == nil || op.f == nil {
		return _wasiEBADF
	}
	if err := op.f.Sync(); err != nil {
		return mapOSError(err)
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_datasync(m *Module, fd int32) int32 {
	w.mu.Lock()
	op := w.fdTable[fd]
	w.mu.Unlock()
	if op == nil || op.f == nil {
		return _wasiEBADF
	}

	if err := op.f.Sync(); err != nil {
		return mapOSError(err)
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_advise(m *Module, fd int32, offset, length int64, advice int32) int32 {
	w.mu.Lock()
	op := w.fdTable[fd]
	w.mu.Unlock()
	if op == nil || op.f == nil {
		return _wasiEBADF
	}

	_, _, _ = offset, length, advice
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_allocate(m *Module, fd int32, offset, length int64) int32 {
	w.mu.Lock()
	op := w.fdTable[fd]
	w.mu.Unlock()
	if op == nil || op.f == nil {
		return _wasiEBADF
	}

	if err := op.f.Truncate(offset + length); err != nil {
		return mapOSError(err)
	}
	return _wasiESUCCESS
}

// Path_chmod is a NON-STANDARD host import (module wasi_snapshot_preview1,
// name "path_chmod") backing a bridge-provided chmod(): WASI preview1 has
// no way to change file modes. The path at (pathPtr,pathLen) is
// preopen-relative, like path_open's. Backends without chmod support
// (MemFS keeps no modes) report ENOSYS. Returns 0 or a negative errno.
func (w *WasiStubs) Path_chmod(m *Module, pathPtr, pathLen, mode int32) int32 {
	pathSlice := w.memSlice(m, pathPtr, pathLen)
	if pathSlice == nil {
		return -_wasiEFAULT
	}
	w.mu.Lock()
	fsys := w.fsys
	w.mu.Unlock()
	ch, ok := fsys.(interface {
		Chmod(string, os.FileMode) error
	})
	if !ok {
		return -_wasiENOSYS
	}
	if err := ch.Chmod(string(pathSlice), os.FileMode(uint32(mode)&0o7777)); err != nil {
		return -mapOSError(err)
	}
	return _wasiESUCCESS
}

// Path_filestat_mode is a NON-STANDARD host import (module
// wasi_snapshot_preview1, name "path_filestat_mode") backing a
// bridge-provided stat/lstat: WASI's filestat carries no permission
// bits, so the bridge merges the real mode in from here. The path is
// preopen-relative; follow selects stat vs lstat semantics. Writes the
// unix permission bits at modeOutPtr; returns 0 or a negative errno.
func (w *WasiStubs) Path_filestat_mode(m *Module, pathPtr, pathLen, follow, modeOutPtr int32) int32 {
	pathSlice := w.memSlice(m, pathPtr, pathLen)
	out := w.memSlice(m, modeOutPtr, 4)
	if pathSlice == nil || out == nil {
		return -_wasiEFAULT
	}
	w.mu.Lock()
	fsys := w.fsys
	w.mu.Unlock()
	var fi os.FileInfo
	var err error
	if follow != 0 {
		fi, err = fsys.Stat(string(pathSlice))
	} else {
		fi, err = fsys.Lstat(string(pathSlice))
	}
	if err != nil {
		return -mapOSError(err)
	}
	mode := fi.Mode()
	bits := uint32(mode.Perm())
	if mode&os.ModeSetuid != 0 {
		bits |= 0o4000
	}
	if mode&os.ModeSetgid != 0 {
		bits |= 0o2000
	}
	if mode&os.ModeSticky != 0 {
		bits |= 0o1000
	}
	binary.LittleEndian.PutUint32(out, bits)
	return _wasiESUCCESS
}

// dupSourceLocked resolves the entry a dup of fd should share: the
// existing table entry, or a fresh alias for a bare interpreter stdio fd.
// Caller holds w.mu.
func (w *WasiStubs) dupSourceLocked(fd int32) *wasiOpen {
	if op := w.fdTable[fd]; op != nil {
		return op
	}
	if fd >= 0 && fd <= 2 {
		op := &wasiOpen{stdio: int8(fd + 1)}
		w.fdTable[fd] = op
		return op
	}
	return nil
}

// Fd_dup is a NON-STANDARD host import (module wasi_snapshot_preview1,
// name "fd_dup") backing the bridge's dup(): the new fd shares the same
// open descriptor (offset included), and the underlying file closes only
// when the last sharing fd does. Writes the new fd at outPtr.
func (w *WasiStubs) Fd_dup(m *Module, fd, outPtr int32) int32 {
	out := w.memSlice(m, outPtr, 4)
	if out == nil {
		return _wasiEFAULT
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	op := w.dupSourceLocked(fd)
	if op == nil {
		return _wasiEBADF
	}
	nfd := w.nextFD
	w.nextFD++
	w.fdTable[nfd] = op
	op.refs++
	binary.LittleEndian.PutUint32(out, uint32(nfd))
	return _wasiESUCCESS
}

// Fd_dup2 is a NON-STANDARD host import (module wasi_snapshot_preview1,
// name "fd_dup2") backing the bridge's dup2(): to becomes another
// reference to from's descriptor, closing whatever to previously held.
func (w *WasiStubs) Fd_dup2(m *Module, from, to int32) int32 {
	w.mu.Lock()
	defer w.mu.Unlock()
	src := w.dupSourceLocked(from)
	if src == nil {
		return _wasiEBADF
	}
	if from == to {
		return _wasiESUCCESS
	}
	var closeErr error
	if dst := w.fdTable[to]; dst != nil {
		if dst == src {
			return _wasiESUCCESS
		}
		closeErr = closeWasiOpen(dst)
	}
	w.fdTable[to] = src
	src.refs++
	if closeErr != nil {
		return mapOSError(closeErr)
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_renumber(m *Module, from, to int32) int32 {
	w.mu.Lock()
	defer w.mu.Unlock()
	if from == to {

		if _, ok := w.fdTable[from]; ok {
			return _wasiESUCCESS
		}
		return _wasiEBADF
	}
	src, ok := w.fdTable[from]
	if !ok {
		return _wasiEBADF
	}
	var closeErr error
	if dst, ok2 := w.fdTable[to]; ok2 {
		closeErr = closeWasiOpen(dst)
	}
	w.fdTable[to] = src
	delete(w.fdTable, from)
	if closeErr != nil {
		return mapOSError(closeErr)
	}
	return _wasiESUCCESS
}

// readDirCached lazily caches the directory listing on first
// Fd_readdir, so paged reads (cookie-driven) walk the same snapshot.
func (op *wasiOpen) readDirCached() ([]os.DirEntry, error) {
	if op.dirCache != nil {
		return op.dirCache, nil
	}
	if op.f == nil {
		return nil, syscall.EBADF
	}
	if _, err := op.f.Seek(0, 0); err != nil {
		return nil, err
	}
	entries, err := op.f.ReadDir(-1)
	if err != nil {
		return nil, err
	}

	out := make([]os.DirEntry, 0, len(entries)+2)
	out = append(out, dotEntry(op.path, "."), dotEntry(op.path, ".."))
	out = append(out, entries...)

	sort.SliceStable(out[2:], func(i, j int) bool {
		return out[2+i].Name() < out[2+j].Name()
	})
	op.dirCache = out
	return out, nil
}

// dotEntry produces a synthetic os.DirEntry for "." and "..". Its
// Info() returns the stat of the parent directory (good enough for
// guest-side d_type detection).
func dotEntry(parent, name string) os.DirEntry {
	return &dotDirEntry{name: name, parent: parent}
}

type dotDirEntry struct {
	name, parent string
}

func (d *dotDirEntry) Name() string { return d.name }
func (d *dotDirEntry) IsDir() bool  { return true }
func (d *dotDirEntry) Type() os.FileMode {
	return os.ModeDir
}
func (d *dotDirEntry) Info() (os.FileInfo, error) {
	if d.name == "." {
		return os.Stat(d.parent)
	}
	return os.Stat(filepath.Dir(d.parent))
}

func (w *WasiStubs) Fd_readdir(m *Module, fd, buf, buflen int32, cookie int64, bufusedPtr int32) int32 {
	bufSlice := w.memSlice(m, buf, buflen)
	bufusedSlice := w.memSlice(m, bufusedPtr, 4)
	if bufSlice == nil || bufusedSlice == nil {
		return _wasiEFAULT
	}
	written, errno := w.fdReaddir(fd, bufSlice, cookie)
	if errno != _wasiESUCCESS {
		return errno
	}
	binary.LittleEndian.PutUint32(bufusedSlice, uint32(written))
	return _wasiESUCCESS
}

// fdReaddir is the layout-independent body of fd_readdir: it packs
// dirents into bufSlice starting at the cookie'th entry and returns the
// byte count used. The dirent wire format has no pointer-width fields,
// so wasm32 and wasm64 share it; only the bufused out-pointer differs.
func (w *WasiStubs) fdReaddir(fd int32, bufSlice []byte, cookie int64) (int, int32) {
	w.mu.Lock()
	op := w.fdTable[fd]
	w.mu.Unlock()
	if op == nil || op.f == nil || !op.isDir {
		return 0, _wasiEBADF
	}

	if cookie == 0 {
		op.dirCache = nil
	}
	entries, err := op.readDirCached()
	if err != nil {
		return 0, mapOSError(err)
	}
	startIdx := int(cookie)
	if startIdx < 0 {
		startIdx = 0
	}
	written := 0
	for i := startIdx; i < len(entries); i++ {
		e := entries[i]
		nameBytes := []byte(e.Name())
		// dirent header: d_next u64 + d_ino u64 + d_namlen u32 + d_type u8 + 3 pad = 24 bytes.
		const headerLen = 24
		// os.FileInfo does not expose inode portably; report 0.
		var dtype byte = 4 // regular file
		if e.IsDir() {
			dtype = 3
		} else if e.Type()&os.ModeSymlink != 0 {
			dtype = 7
		} else if e.Type()&os.ModeNamedPipe != 0 {
			dtype = 6
		} else if e.Type()&os.ModeSocket != 0 {
			dtype = 6
		}
		// Assemble the fixed header, then copy header+name into the buffer.
		// When a record does not fully fit we copy as much as fits so that
		// bufused == buflen, which is the wasi-libc signal for "more entries
		// available; call again with the last returned cookie". We must NOT
		// zero-fill the leftover: a zeroed dirent (d_namlen=0, d_next=0) is
		// misread by wasi-libc as end-of-directory and silently truncates the
		// listing (e.g. makes a guest's importer miss standard-library packages).
		var hdr [headerLen]byte
		binary.LittleEndian.PutUint64(hdr[0:], uint64(i+1))

		binary.LittleEndian.PutUint64(hdr[8:], uint64(i)+1)
		binary.LittleEndian.PutUint32(hdr[16:], uint32(len(nameBytes)))
		hdr[20] = dtype
		n := copy(bufSlice[written:], hdr[:])
		written += n
		if n < len(hdr) {
			written = len(bufSlice)
			break
		}
		n = copy(bufSlice[written:], nameBytes)
		written += n
		if n < len(nameBytes) {
			written = len(bufSlice)
			break
		}
	}
	return written, _wasiESUCCESS
}

// Path_open opens a wasm-supplied path and registers it in the fd
// table. The path is resolved against the host filesystem with the same
// rights the host Go process has — wasm2go's default WASI is a thin
// passthrough, not a sandbox. The dirFd == 3 special case keeps the
// "preopen /" convention that wasi-libc requires for its directory
// enumeration, but the path itself is opened verbatim (joined to "/")
// using os.OpenFile. Callers that need a sandbox should provide their
// own Wasi_snapshot_preview1Imports implementation via NewWithWASI.
func (w *WasiStubs) Path_open(m *Module, dirFd, dirflags, pathPtr, pathLen, oflags int32, fsRightsBase, fsRightsInherit int64, fdflags, openedFdPtr int32) int32 {
	if dirFd != 3 {
		return _wasiEBADF
	}
	pathSlice := w.memSlice(m, pathPtr, pathLen)
	outSlice := w.memSlice(m, openedFdPtr, 4)
	if pathSlice == nil || outSlice == nil {
		return _wasiEFAULT
	}
	fd, errno := w.pathOpen(string(pathSlice), dirflags, oflags, fsRightsBase, fdflags)
	if errno != _wasiESUCCESS {
		return errno
	}
	binary.LittleEndian.PutUint32(outSlice, uint32(fd))
	return _wasiESUCCESS
}

// pathOpen is the layout-independent body of path_open: it resolves and
// opens rel, registers the fd, and returns it. Callers own reading the
// path and writing the opened fd at their ABI's pointer width.
func (w *WasiStubs) pathOpen(rel string, dirflags, oflags int32, fsRightsBase int64, fdflags int32) (int32, int32) {
	w.mu.Lock()
	fsys := w.fsys
	w.mu.Unlock()

	canRead := fsRightsBase&(1<<1) != 0
	canWrite := fsRightsBase&(1<<6) != 0
	var flag int
	switch {
	case canRead && canWrite:
		flag = os.O_RDWR
	case canWrite && !canRead:
		flag = os.O_WRONLY
	default:
		flag = os.O_RDONLY
	}

	if oflags&0x1 != 0 {
		flag |= os.O_CREATE
	}
	if oflags&0x4 != 0 {
		flag |= os.O_EXCL
	}
	if oflags&0x8 != 0 {
		flag |= os.O_TRUNC
	}

	if fdflags&0x1 != 0 {
		flag |= os.O_APPEND
	}
	if fdflags&(0x2|0x8|0x10) != 0 {
		flag |= os.O_SYNC
	}

	writeAccess := flag&(os.O_WRONLY|os.O_RDWR) != 0 || flag&(os.O_CREATE|os.O_TRUNC) != 0
	if !w.checkFS(rel, writeAccess) {
		return -1, _wasiEACCES
	}

	requireDir := oflags&0x2 != 0
	noFollow := dirflags&0x1 == 0

	if requireDir {

		flag = os.O_RDONLY
	}

	if noFollow {
		if li, lerr := fsys.Lstat(rel); lerr == nil && (li.Mode()&os.ModeSymlink) != 0 {
			return -1, _wasiENOENT
		}
	}
	f, err := fsys.OpenFile(rel, flag, 0o644)
	if err != nil {
		return -1, mapOSError(err)
	}
	st, statErr := f.Stat()
	if statErr != nil {
		return -1, mapOSError(errors.Join(statErr, f.Close()))
	}
	isDir := st.IsDir()
	if requireDir && !isDir {
		if cerr := f.Close(); cerr != nil {
			return -1, mapOSError(cerr)
		}
		return -1, _wasiENOTDIR
	}
	w.mu.Lock()
	fd := w.nextFD
	w.nextFD++
	w.fdTable[fd] = &wasiOpen{f: f, isDir: isDir, path: rel, fdflags: fdflags}
	w.mu.Unlock()
	return fd, _wasiESUCCESS
}

func (w *WasiStubs) Path_create_directory(m *Module, dirFd, pathPtr, pathLen int32) int32 {
	if dirFd != 3 {
		return _wasiEBADF
	}
	pathSlice := w.memSlice(m, pathPtr, pathLen)
	if pathSlice == nil {
		return _wasiEFAULT
	}
	if !w.checkFS(string(pathSlice), true) {
		return _wasiEACCES
	}
	w.mu.Lock()
	fsys := w.fsys
	w.mu.Unlock()
	if err := fsys.Mkdir(string(pathSlice), 0o755); err != nil {
		return mapOSError(err)
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Path_unlink_file(m *Module, dirFd, pathPtr, pathLen int32) int32 {
	if dirFd != 3 {
		return _wasiEBADF
	}
	pathSlice := w.memSlice(m, pathPtr, pathLen)
	if pathSlice == nil {
		return _wasiEFAULT
	}
	if !w.checkFS(string(pathSlice), true) {
		return _wasiEACCES
	}
	w.mu.Lock()
	fsys := w.fsys
	w.mu.Unlock()
	rel := string(pathSlice)
	st, err := fsys.Lstat(rel)
	if err != nil {
		return mapOSError(err)
	}
	if st.IsDir() {
		return _wasiEISDIR
	}
	if err := fsys.Remove(rel); err != nil {
		return mapOSError(err)
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Path_remove_directory(m *Module, dirFd, pathPtr, pathLen int32) int32 {
	if dirFd != 3 {
		return _wasiEBADF
	}
	pathSlice := w.memSlice(m, pathPtr, pathLen)
	if pathSlice == nil {
		return _wasiEFAULT
	}
	w.mu.Lock()
	fsys := w.fsys
	w.mu.Unlock()
	rel := string(pathSlice)
	st, err := fsys.Lstat(rel)
	if err != nil {
		return mapOSError(err)
	}
	if !st.IsDir() {
		return _wasiENOTDIR
	}
	if err := fsys.Remove(rel); err != nil {
		return mapOSError(err)
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Path_rename(m *Module, oldFd, oldPathPtr, oldPathLen, newFd, newPathPtr, newPathLen int32) int32 {
	if oldFd != 3 || newFd != 3 {
		return _wasiEBADF
	}
	oldSlice := w.memSlice(m, oldPathPtr, oldPathLen)
	newSlice := w.memSlice(m, newPathPtr, newPathLen)
	if oldSlice == nil || newSlice == nil {
		return _wasiEFAULT
	}
	w.mu.Lock()
	fsys := w.fsys
	w.mu.Unlock()
	if err := fsys.Rename(string(oldSlice), string(newSlice)); err != nil {
		return mapOSError(err)
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Path_filestat_get(m *Module, dirFd, flags, pathPtr, pathLen, outPtr int32) int32 {
	if dirFd != 3 {
		return _wasiEBADF
	}
	pathSlice := w.memSlice(m, pathPtr, pathLen)
	out := w.memSlice(m, outPtr, 64)
	if pathSlice == nil || out == nil {
		return _wasiEFAULT
	}
	w.mu.Lock()
	fsys := w.fsys
	w.mu.Unlock()
	rel := string(pathSlice)
	var st os.FileInfo
	var err error
	if flags&0x1 != 0 {
		st, err = fsys.Stat(rel)
	} else {
		st, err = fsys.Lstat(rel)
	}
	if err != nil {
		return mapOSError(err)
	}
	writeFilestat(out, st)
	return _wasiESUCCESS
}

func (w *WasiStubs) Path_filestat_set_times(m *Module, dirFd, flags, pathPtr, pathLen int32, atim, mtim int64, fstFlags int32) int32 {
	if dirFd != 3 {
		return _wasiEBADF
	}
	pathSlice := w.memSlice(m, pathPtr, pathLen)
	if pathSlice == nil {
		return _wasiEFAULT
	}
	w.mu.Lock()
	fsys := w.fsys
	w.mu.Unlock()
	rel := string(pathSlice)
	follow := flags&0x1 != 0
	now := time.Now()
	var st os.FileInfo
	var statErr error
	if follow {
		st, statErr = fsys.Stat(rel)
	} else {
		st, statErr = fsys.Lstat(rel)
	}
	if statErr != nil {
		return mapOSError(statErr)
	}
	atime := st.ModTime()
	mtime := st.ModTime()
	if fstFlags&0x1 != 0 {
		atime = time.Unix(0, int64(atim))
	}
	if fstFlags&0x2 != 0 {
		atime = now
	}
	if fstFlags&0x4 != 0 {
		mtime = time.Unix(0, int64(mtim))
	}
	if fstFlags&0x8 != 0 {
		mtime = now
	}

	if cf, ok := fsys.(chtimesFS); ok {
		if err := cf.Chtimes(rel, atime, mtime); err != nil {
			return mapOSError(err)
		}
	}
	return _wasiESUCCESS
}

// chtimesFS is an optional FS capability for backends that track timestamps.
type chtimesFS interface {
	Chtimes(name string, atime, mtime time.Time) error
}

func (o osFS) Chtimes(name string, atime, mtime time.Time) error {
	return os.Chtimes(o.join(name), atime, mtime)
}

func (w *WasiStubs) Path_link(m *Module, oldFd, oldFlags, oldPathPtr, oldPathLen, newFd, newPathPtr, newPathLen int32) int32 {
	if oldFd != 3 || newFd != 3 {
		return _wasiEBADF
	}
	oldSlice := w.memSlice(m, oldPathPtr, oldPathLen)
	newSlice := w.memSlice(m, newPathPtr, newPathLen)
	if oldSlice == nil || newSlice == nil {
		return _wasiEFAULT
	}
	w.mu.Lock()
	fsys := w.fsys
	w.mu.Unlock()
	if err := fsys.Link(string(oldSlice), string(newSlice)); err != nil {
		return mapOSError(err)
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Path_symlink(m *Module, targetPtr, targetLen, dirFd, linkPathPtr, linkPathLen int32) int32 {
	if dirFd != 3 {
		return _wasiEBADF
	}
	targetSlice := w.memSlice(m, targetPtr, targetLen)
	linkSlice := w.memSlice(m, linkPathPtr, linkPathLen)
	if targetSlice == nil || linkSlice == nil {
		return _wasiEFAULT
	}
	w.mu.Lock()
	fsys := w.fsys
	w.mu.Unlock()
	if err := fsys.Symlink(string(targetSlice), string(linkSlice)); err != nil {
		return mapOSError(err)
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Path_readlink(m *Module, dirFd, pathPtr, pathLen, buf, buflen, bufusedPtr int32) int32 {
	if dirFd != 3 {
		return _wasiEBADF
	}
	pathSlice := w.memSlice(m, pathPtr, pathLen)
	bufSlice := w.memSlice(m, buf, buflen)
	bufused := w.memSlice(m, bufusedPtr, 4)
	if pathSlice == nil || bufSlice == nil || bufused == nil {
		return _wasiEFAULT
	}
	w.mu.Lock()
	fsys := w.fsys
	w.mu.Unlock()
	target, err := fsys.Readlink(string(pathSlice))
	if err != nil {
		return mapOSError(err)
	}
	n := copy(bufSlice, target)
	binary.LittleEndian.PutUint32(bufused, uint32(n))
	return _wasiESUCCESS
}

func (w *WasiStubs) Random_get(m *Module, buf, bufLen int32) int32 {
	slice := w.memSlice(m, buf, bufLen)
	if slice == nil {
		return _wasiEFAULT
	}
	_, err := rand.Read(slice)
	if err != nil {
		return _wasiEIO
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Sched_yield(m *Module) int32 {
	runtime.Gosched()
	return _wasiESUCCESS
}

// Poll_oneoff decodes the WASI subscription_u records and reproduces the
// requested events.
//
// Each subscription is 48 bytes:
//
//	u64 userdata
//	u8  eventtype  (0=clock, 1=fd_read, 2=fd_write)
//	... per-type payload starting at offset 16
//
// For clock subscriptions, payload at offset 16 is: u32 clock_id, u64
// timeout, u64 precision, u16 sub_clock_flags (bit0=ABSTIME). We sleep
// for `timeout` ns (relative timer) or the diff to `timeout` (absolute
// timer). For fd_read / fd_write subscriptions, payload at offset 16 is
// a u32 fd; we call into the platform Poll helper to wait for
// readiness.
//
// Each emitted event is 32 bytes: u64 userdata, u16 errno, u16
// eventtype, u64 fd_readwrite_nbytes (filled for fd events), u16
// flags, then 6 bytes of padding.
func (w *WasiStubs) Poll_oneoff(m *Module, inPtr, outPtr, nsubs, neventsPtr int32) int32 {
	subsTotal := uint64(uint32(nsubs)) * 48
	if subsTotal > 0x7fffffff {
		return _wasiEFAULT
	}
	subs := w.memSlice(m, inPtr, int32(subsTotal))
	evTotal := uint64(uint32(nsubs)) * 32
	if evTotal > 0x7fffffff {
		return _wasiEFAULT
	}
	events := w.memSlice(m, outPtr, int32(evTotal))
	nev := w.memSlice(m, neventsPtr, 4)
	if subs == nil || events == nil || nev == nil {
		return _wasiEFAULT
	}

	type pollItem struct {
		userdata uint64
		etype    byte
		fd       int32
		isRead   bool
	}
	var minClockNs int64 = -1
	var clockEvents []pollItem
	var fdEvents []pollItem
	for i := int32(0); i < nsubs; i++ {
		base := i * 48
		userdata := binary.LittleEndian.Uint64(subs[base:])
		etype := subs[base+8]
		switch etype {
		case 0:
			timeout := int64(binary.LittleEndian.Uint64(subs[base+24:]))
			flags := binary.LittleEndian.Uint16(subs[base+40:])
			ns := timeout
			if flags&0x1 != 0 {

				ns = timeout - time.Now().UnixNano()
				if ns < 0 {
					ns = 0
				}
			}
			if minClockNs < 0 || ns < minClockNs {
				minClockNs = ns
			}
			clockEvents = append(clockEvents, pollItem{userdata: userdata, etype: 0})
		case 1, 2:
			fd := int32(binary.LittleEndian.Uint32(subs[base+16:]))
			fdEvents = append(fdEvents, pollItem{userdata: userdata, etype: etype, fd: fd, isRead: etype == 1})
		default:

			clockEvents = append(clockEvents, pollItem{userdata: userdata, etype: etype})
		}
	}

	if minClockNs > 0 && len(fdEvents) == 0 {
		time.Sleep(time.Duration(minClockNs))
	}

	written := int32(0)
	for _, ev := range clockEvents {
		if ev.etype == 0 && len(fdEvents) > 0 {
			continue
		}
		writeEvent(events[written:written+32], ev.userdata, ev.etype, 0, 0)
		written += 32
	}
	for _, ev := range fdEvents {
		w.mu.Lock()
		op := w.fdTable[ev.fd]
		w.mu.Unlock()
		var errno int32
		var nbytes uint64
		if op == nil {
			errno = _wasiEBADF
		} else if op.f != nil {

			if ev.isRead {
				if st, err := op.f.Stat(); err == nil {

					if cur, err := op.f.Seek(0, 1); err == nil && st.Size() > cur {
						nbytes = uint64(st.Size() - cur)
					}
				}
			}
		} else if op.conn != nil {

			_ = minClockNs
		}
		writeEvent(events[written:written+32], ev.userdata, ev.etype, uint16(errno), nbytes)
		written += 32
	}

	binary.LittleEndian.PutUint32(nev, uint32(written/32))
	return _wasiESUCCESS
}

func writeEvent(dst []byte, userdata uint64, etype byte, errno uint16, nbytes uint64) {
	for i := range dst {
		dst[i] = 0
	}
	binary.LittleEndian.PutUint64(dst[0:], userdata)
	binary.LittleEndian.PutUint16(dst[8:], errno)
	binary.LittleEndian.PutUint16(dst[10:], uint16(etype))
	binary.LittleEndian.PutUint64(dst[16:], nbytes)
}

func (w *WasiStubs) Proc_exit(m *Module, code int32) {

	panic(&WasiExitError{Code: code})
}

func (w *WasiStubs) Proc_raise(m *Module, sig int32) int32 {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return mapOSError(err)
	}
	if err := p.Signal(syscall.Signal(sig)); err != nil {
		return mapOSError(err)
	}
	return _wasiESUCCESS
}

// Sock_socket is a NON-STANDARD host import (module wasi_snapshot_preview1,
// name "sock_socket") that backs a libc socket() call wrapped via
// -Wl,--wrap=socket in the guest. WASI preview1 has no way to create an
// outbound socket; this gives the guest a host-managed fd whose connection is
// established later by Sock_connect. domain/type follow the POSIX socket()
// args (AF_INET / SOCK_STREAM); only TCP over IPv4 is supported. Returns the
// new fd, or a negative errno on failure.
func (w *WasiStubs) Sock_socket(m *Module, domain, typ int32) int32 {

	_ = domain
	_ = typ
	w.mu.Lock()
	defer w.mu.Unlock()
	fd := w.nextFD
	w.nextFD++
	w.fdTable[fd] = &wasiOpen{isSocket: true}
	return fd
}

// Sock_connect is a NON-STANDARD host import (module wasi_snapshot_preview1,
// name "sock_connect") backing a libc connect() wrapped via
// -Wl,--wrap=connect. ipBE carries the IPv4 address in network byte order
// exactly as it sat in sockaddr_in.sin_addr.s_addr (so the low byte is the
// first octet); port is host byte order. It consults the dial whitelist,
// dials via Go's net, and attaches the resulting conn to the socket fd so the
// existing Sock_send / Sock_recv / Fd_close paths drive it. Returns 0 or a
// negative errno.
func (w *WasiStubs) Sock_connect(m *Module, fd, ipBE, port int32) int32 {
	u := uint32(ipBE)
	ip := fmt.Sprintf("%d.%d.%d.%d", u&0xff, (u>>8)&0xff, (u>>16)&0xff, (u>>24)&0xff)
	w.mu.Lock()
	op := w.fdTable[fd]
	hook := w.dialHook
	host := w.resolvedHosts[ip]
	w.mu.Unlock()
	if op == nil || !op.isSocket {
		return -_wasiENOTSOCK
	}
	if op.conn != nil {
		return -_wasiEISCONN
	}
	p := int(uint16(port))
	if hook != nil && !hook("tcp", host, ip, p) {
		return -_wasiEACCES
	}
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, strconv.Itoa(p)), 30*time.Second)
	if err != nil {
		return -_wasiECONNREFUSED
	}
	w.mu.Lock()

	if cur := w.fdTable[fd]; cur == op {
		op.conn = conn
		w.mu.Unlock()
		return _wasiESUCCESS
	}
	w.mu.Unlock()
	if cerr := conn.Close(); cerr != nil {
		return mapOSError(cerr)
	}
	return -_wasiEBADF
}

// Sock_accept accepts the next incoming TCP/Unix connection on the
// listener associated with fd, registers it as a new wasiOpen with a
// conn arm, and writes the new fd at fdOutPtr. Returns ENOTSOCK if fd
// isn't a listener.
func (w *WasiStubs) Sock_accept(m *Module, fd, flags, fdOutPtr int32) int32 {
	if !w.checkNet("accept") {
		return _wasiEACCES
	}
	out := w.memSlice(m, fdOutPtr, 4)
	if out == nil {
		return _wasiEFAULT
	}
	w.mu.Lock()
	op := w.fdTable[fd]
	w.mu.Unlock()
	if op == nil || op.listener == nil {
		return _wasiENOTSOCK
	}
	conn, err := op.listener.Accept()
	if err != nil {
		return mapOSError(err)
	}
	w.mu.Lock()
	newFD := w.nextFD
	w.nextFD++
	w.fdTable[newFD] = &wasiOpen{conn: conn}
	w.mu.Unlock()
	binary.LittleEndian.PutUint32(out, uint32(newFD))
	return _wasiESUCCESS
}

func (w *WasiStubs) Sock_recv(m *Module, fd, riData, riDataLen, riFlags, roDataLenPtr, roFlagsPtr int32) int32 {
	if !w.checkNet("recv") {
		return _wasiEACCES
	}
	w.mu.Lock()
	op := w.fdTable[fd]
	w.mu.Unlock()
	if op == nil || op.conn == nil {
		return _wasiENOTSOCK
	}
	iovBytes := uint64(uint32(riDataLen)) * 8
	if iovBytes > 0x7fffffff {
		return _wasiEFAULT
	}
	iovecs := w.memSlice(m, riData, int32(iovBytes))
	lenOut := w.memSlice(m, roDataLenPtr, 4)

	flagsOut := w.memSlice(m, roFlagsPtr, 2)
	if iovecs == nil || lenOut == nil || flagsOut == nil {
		return _wasiEFAULT
	}
	var total uint32
	for i := int32(0); i < riDataLen; i++ {
		bufPtr := binary.LittleEndian.Uint32(iovecs[i*8:])
		bufLen := binary.LittleEndian.Uint32(iovecs[i*8+4:])
		buf := w.memSlice(m, int32(bufPtr), int32(bufLen))
		if buf == nil {
			return _wasiEFAULT
		}
		n, err := op.conn.Read(buf)
		total += uint32(n)
		if err != nil {
			break
		}
	}
	binary.LittleEndian.PutUint16(flagsOut, 0)
	binary.LittleEndian.PutUint32(lenOut, total)
	return _wasiESUCCESS
}

func (w *WasiStubs) Sock_send(m *Module, fd, siData, siDataLen, siFlags, soDataLenPtr int32) int32 {
	if !w.checkNet("send") {
		return _wasiEACCES
	}
	w.mu.Lock()
	op := w.fdTable[fd]
	w.mu.Unlock()
	if op == nil || op.conn == nil {
		return _wasiENOTSOCK
	}
	iovBytes := uint64(uint32(siDataLen)) * 8
	if iovBytes > 0x7fffffff {
		return _wasiEFAULT
	}
	iovecs := w.memSlice(m, siData, int32(iovBytes))
	lenOut := w.memSlice(m, soDataLenPtr, 4)
	if iovecs == nil || lenOut == nil {
		return _wasiEFAULT
	}
	var total uint32
	for i := int32(0); i < siDataLen; i++ {
		bufPtr := binary.LittleEndian.Uint32(iovecs[i*8:])
		bufLen := binary.LittleEndian.Uint32(iovecs[i*8+4:])
		buf := w.memSlice(m, int32(bufPtr), int32(bufLen))
		if buf == nil {
			return _wasiEFAULT
		}
		n, err := op.conn.Write(buf)
		total += uint32(n)
		if err != nil {
			break
		}
	}
	binary.LittleEndian.PutUint32(lenOut, total)
	return _wasiESUCCESS
}

func (w *WasiStubs) Sock_shutdown(m *Module, fd, how int32) int32 {
	w.mu.Lock()
	op := w.fdTable[fd]
	w.mu.Unlock()
	if op == nil || op.conn == nil {
		return _wasiENOTSOCK
	}
	type shutdowner interface {
		CloseRead() error
		CloseWrite() error
	}
	sh, ok := op.conn.(shutdowner)
	if !ok {

		if err := op.conn.Close(); err != nil {
			return mapOSError(err)
		}
		return _wasiESUCCESS
	}
	var shErr error
	if how&0x1 != 0 {
		shErr = errors.Join(shErr, sh.CloseRead())
	}
	if how&0x2 != 0 {
		shErr = errors.Join(shErr, sh.CloseWrite())
	}
	if shErr != nil {
		return mapOSError(shErr)
	}
	return _wasiESUCCESS
}

// writeFilestat populates the 64-byte WASI filestat structure from a
// host os.FileInfo. The dev/ino fields come from the per-platform
// wasiPlatformStatSys helper (unix returns Stat_t.Dev/.Ino; Windows
// returns zeros).
func writeFilestat(out []byte, st os.FileInfo) {

	binary.LittleEndian.PutUint64(out[0:], 0)
	binary.LittleEndian.PutUint64(out[8:], 0)
	var ftype byte = 4
	mode := st.Mode()
	switch {
	case mode.IsDir():
		ftype = 3
	case mode&os.ModeSymlink != 0:
		ftype = 7
	case mode&os.ModeNamedPipe != 0:
		ftype = 6
	case mode&os.ModeSocket != 0:
		ftype = 6
	case mode&os.ModeDevice != 0:
		ftype = 1
	case mode&os.ModeCharDevice != 0:
		ftype = 2
	}
	out[16] = ftype
	binary.LittleEndian.PutUint64(out[24:], 1)
	binary.LittleEndian.PutUint64(out[32:], uint64(st.Size()))
	nanos := uint64(st.ModTime().UnixNano())
	binary.LittleEndian.PutUint64(out[40:], nanos)
	binary.LittleEndian.PutUint64(out[48:], nanos)
	binary.LittleEndian.PutUint64(out[56:], nanos)
}

// memSlice64 is memSlice for full-range 64-bit guest pointers.
func (w *WasiStubs) memSlice64(m *Module, off int64, n int64) []byte {
	mem := m.memory
	lo := uint64(off)
	hi := lo + uint64(n)
	if n < 0 || hi < lo || hi > uint64(len(mem)) {
		return nil
	}
	return mem[lo:hi]
}

func (w *WasiStubs) Clock_time_get64(m *Module, clockID int64, precision int64, timePtr int64) int32 {
	out := w.memSlice64(m, timePtr, 8)
	if out == nil {
		return _wasiEFAULT
	}
	nanos, errno := w.clockNanos(int32(clockID))
	if errno != _wasiESUCCESS {
		return errno
	}
	binary.LittleEndian.PutUint64(out, nanos)
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_close64(m *Module, fd int64) int32 {
	return w.Fd_close(m, int32(fd))
}

func (w *WasiStubs) Sched_yield64(m *Module) int32 {

	return w.Sched_yield(m)
}

func (w *WasiStubs) Fd_fdstat_get64(m *Module, fd int64, ptr int64) int32 {

	out := w.memSlice64(m, ptr, 24)
	if out == nil {
		return _wasiEFAULT
	}
	return w.fdstatFill(int32(fd), out)
}

func (w *WasiStubs) Fd_seek64(m *Module, fd int64, offset int64, whence int64, newOffPtr int64) int32 {
	out := w.memSlice64(m, newOffPtr, 8)
	if out == nil {
		return _wasiEFAULT
	}
	n, errno := w.fdSeek(int32(fd), offset, int(whence))
	if errno != _wasiESUCCESS {
		return errno
	}
	binary.LittleEndian.PutUint64(out, uint64(n))
	return _wasiESUCCESS
}

// iovecSlices64 is iovecSlices for the LP64 iovec layout: {u64 buf,
// u64 len}, 16 bytes per entry.
func (w *WasiStubs) iovecSlices64(m *Module, iovs, iovsLen int64) ([][]byte, bool) {
	if iovsLen < 0 || iovsLen > 1<<20 {
		return nil, false
	}
	iovecs := w.memSlice64(m, iovs, iovsLen*16)
	if iovecs == nil {
		return nil, false
	}
	bufs := make([][]byte, 0, iovsLen)
	for i := int64(0); i < iovsLen; i++ {
		bufPtr := binary.LittleEndian.Uint64(iovecs[i*16:])
		bufLen := binary.LittleEndian.Uint64(iovecs[i*16+8:])
		buf := w.memSlice64(m, int64(bufPtr), int64(bufLen))
		if buf == nil {
			return nil, false
		}
		bufs = append(bufs, buf)
	}
	return bufs, true
}

func (w *WasiStubs) Fd_write64(m *Module, fd int64, iovs int64, iovsLen int64, nwrittenPtr int64) int32 {
	w.mu.Lock()
	dst, _ := w.fdDstLocked(int32(fd))
	w.mu.Unlock()
	bufs, ok := w.iovecSlices64(m, iovs, iovsLen)

	nwrittenSlice := w.memSlice64(m, nwrittenPtr, 8)
	if !ok || nwrittenSlice == nil {
		return _wasiEFAULT
	}
	if dst == nil {
		binary.LittleEndian.PutUint64(nwrittenSlice, 0)
		return _wasiEBADF
	}
	binary.LittleEndian.PutUint64(nwrittenSlice, writeVec(dst, bufs))
	return _wasiESUCCESS
}

func (w *WasiStubs) Proc_exit64(m *Module, code int64) {
	panic(&WasiExitError{Code: int32(code)})
}

// putStrVec64 packs ss as an LP64 char** table (8-byte guest pointers
// at vec) plus NUL-terminated bodies (at buf, guest address bufBase).
// Both slices must already be sized: len(ss)*8 and totalBytesPlusNul.
func putStrVec64(vec, buf []byte, bufBase uint64, ss []string) int32 {
	bufOff := uint64(0)
	for i, s := range ss {
		binary.LittleEndian.PutUint64(vec[i*8:], bufBase+bufOff)
		n := copy(buf[bufOff:], s)
		if n < len(s) {
			return _wasiEFAULT
		}
		bufOff += uint64(n)
		buf[bufOff] = 0
		bufOff++
	}
	return _wasiESUCCESS
}

func (w *WasiStubs) Args_get64(m *Module, argv, argvBuf int64) int32 {
	w.mu.Lock()
	defer w.mu.Unlock()
	argvSlice := w.memSlice64(m, argv, int64(len(w.args))*8)
	if argvSlice == nil {
		return _wasiEFAULT
	}
	total, ok := totalBytesPlusNul(w.args)
	if !ok {
		return _wasiEFAULT
	}
	argvBufSlice := w.memSlice64(m, argvBuf, int64(total))
	if argvBufSlice == nil {
		return _wasiEFAULT
	}
	return putStrVec64(argvSlice, argvBufSlice, uint64(argvBuf), w.args)
}

func (w *WasiStubs) Args_sizes_get64(m *Module, argcPtr, argvBufLenPtr int64) int32 {
	w.mu.Lock()
	defer w.mu.Unlock()
	argcSlice := w.memSlice64(m, argcPtr, 8)
	bufLenSlice := w.memSlice64(m, argvBufLenPtr, 8)
	if argcSlice == nil || bufLenSlice == nil {
		return _wasiEFAULT
	}
	total, ok := totalBytesPlusNul(w.args)
	if !ok {
		return _wasiEFAULT
	}
	binary.LittleEndian.PutUint64(argcSlice, uint64(len(w.args)))
	binary.LittleEndian.PutUint64(bufLenSlice, uint64(total))
	return _wasiESUCCESS
}

func (w *WasiStubs) Environ_get64(m *Module, envv, envBuf int64) int32 {
	w.mu.Lock()
	defer w.mu.Unlock()
	envvSlice := w.memSlice64(m, envv, int64(len(w.env))*8)
	if envvSlice == nil {
		return _wasiEFAULT
	}
	total, ok := totalBytesPlusNul(w.env)
	if !ok {
		return _wasiEFAULT
	}
	envBufSlice := w.memSlice64(m, envBuf, int64(total))
	if envBufSlice == nil {
		return _wasiEFAULT
	}
	return putStrVec64(envvSlice, envBufSlice, uint64(envBuf), w.env)
}

func (w *WasiStubs) Environ_sizes_get64(m *Module, envcPtr, envBufLenPtr int64) int32 {
	w.mu.Lock()
	defer w.mu.Unlock()
	envcSlice := w.memSlice64(m, envcPtr, 8)
	bufLenSlice := w.memSlice64(m, envBufLenPtr, 8)
	if envcSlice == nil || bufLenSlice == nil {
		return _wasiEFAULT
	}
	total, ok := totalBytesPlusNul(w.env)
	if !ok {
		return _wasiEFAULT
	}
	binary.LittleEndian.PutUint64(envcSlice, uint64(len(w.env)))
	binary.LittleEndian.PutUint64(bufLenSlice, uint64(total))
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_fdstat_set_flags64(m *Module, fd, flags int64) int32 {
	return w.Fd_fdstat_set_flags(m, int32(fd), int32(flags))
}

func (w *WasiStubs) Fd_prestat_get64(m *Module, fd, ptr int64) int32 {
	if int32(fd) != 3 {
		return _wasiEBADF
	}

	out := w.memSlice64(m, ptr, 16)
	if out == nil {
		return _wasiEFAULT
	}
	out[0] = 0
	binary.LittleEndian.PutUint64(out[8:], 1)
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_prestat_dir_name64(m *Module, fd, buf, buflen int64) int32 {
	if int32(fd) != 3 {
		return _wasiEBADF
	}
	if buflen < 1 {
		return _wasiESUCCESS
	}
	out := w.memSlice64(m, buf, buflen)
	if out == nil {
		return _wasiEFAULT
	}
	out[0] = '/'
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_read64(m *Module, fd, iovs, iovsLen, nreadPtr int64) int32 {
	w.mu.Lock()
	src, _ := w.fdSrcLocked(int32(fd))
	w.mu.Unlock()
	if src == nil {
		return _wasiEBADF
	}
	bufs, ok := w.iovecSlices64(m, iovs, iovsLen)

	nreadSlice := w.memSlice64(m, nreadPtr, 8)
	if !ok || nreadSlice == nil {
		return _wasiEFAULT
	}
	binary.LittleEndian.PutUint64(nreadSlice, readVec(src, bufs))
	return _wasiESUCCESS
}

func (w *WasiStubs) Fd_readdir64(m *Module, fd, buf, buflen, cookie, bufusedPtr int64) int32 {

	bufSlice := w.memSlice64(m, buf, buflen)
	bufusedSlice := w.memSlice64(m, bufusedPtr, 8)
	if bufSlice == nil || bufusedSlice == nil {
		return _wasiEFAULT
	}
	written, errno := w.fdReaddir(int32(fd), bufSlice, cookie)
	if errno != _wasiESUCCESS {
		return errno
	}
	binary.LittleEndian.PutUint64(bufusedSlice, uint64(written))
	return _wasiESUCCESS
}

func (w *WasiStubs) Path_open64(m *Module, dirFd, dirflags, pathPtr, pathLen, oflags, fsRightsBase, fsRightsInherit, fdflags, openedFdPtr int64) int32 {
	if int32(dirFd) != 3 {
		return _wasiEBADF
	}
	pathSlice := w.memSlice64(m, pathPtr, pathLen)

	outSlice := w.memSlice64(m, openedFdPtr, 4)
	if pathSlice == nil || outSlice == nil {
		return _wasiEFAULT
	}
	fd, errno := w.pathOpen(string(pathSlice), int32(dirflags), int32(oflags), fsRightsBase, int32(fdflags))
	if errno != _wasiESUCCESS {
		return errno
	}
	binary.LittleEndian.PutUint32(outSlice, uint32(fd))
	return _wasiESUCCESS
}

func (w *WasiStubs) Path_filestat_get64(m *Module, dirFd, flags, pathPtr, pathLen, outPtr int64) int32 {
	if int32(dirFd) != 3 {
		return _wasiEBADF
	}
	pathSlice := w.memSlice64(m, pathPtr, pathLen)

	out := w.memSlice64(m, outPtr, 64)
	if pathSlice == nil || out == nil {
		return _wasiEFAULT
	}
	w.mu.Lock()
	fsys := w.fsys
	w.mu.Unlock()
	rel := string(pathSlice)
	var st os.FileInfo
	var err error
	if flags&0x1 != 0 {
		st, err = fsys.Stat(rel)
	} else {
		st, err = fsys.Lstat(rel)
	}
	if err != nil {
		return mapOSError(err)
	}
	writeFilestat(out, st)
	return _wasiESUCCESS
}

func (w *WasiStubs) Random_get64(m *Module, buf, bufLen int64) int32 {
	slice := w.memSlice64(m, buf, bufLen)
	if slice == nil {
		return _wasiEFAULT
	}
	if _, err := rand.Read(slice); err != nil {
		return _wasiEIO
	}
	return _wasiESUCCESS
}
func NewWithWASIReserve(wasi_snapshot_preview1 Wasi_snapshot_preview1Imports, env EnvImports, reserveBytes int) *Module {
	m := &Module{wasi_snapshot_preview1: wasi_snapshot_preview1, env: env}
	__memcap := reserveBytes
	if __memcap < 67108864 {
		__memcap = 67108864
	}
	m.memory = make([]byte, 67108864, __memcap)
	m.memMu = &sync.Mutex{}
	m.memSize = &atomic.Uint64{}
	m.threads = &threadPool{}
	m.memSize.Store(67108864)
	m.M = unsafe.Pointer(unsafe.SliceData(m.memory))
	m.maxMem = 2147483648
	m.t0 = make([]any, 143)
	m.g0 = int32(10937088)
	m.g1 = int32(0)
	m.g2 = int32(0)
	m.g3 = int32(10937088)
	m.g4 = int32(0)
	m.g5 = int32(97288)
	m.g6 = int32(97592)
	m.g7 = int32(95896)
	m.g8 = int32(95904)
	m.g9 = int32(95900)
	m.g10 = int32(126752)
	m.g11 = int32(96808)
	m.g12 = int32(85164)
	m.g13 = int32(85240)
	m.g14 = int32(85296)
	m.g15 = int32(85336)
	m.g16 = int32(86180)
	m.g17 = int32(22)
	m.g18 = int32(23)
	m.g19 = int32(29)
	m.g20 = int32(30)
	m.g21 = int32(31)
	m.g22 = int32(85368)
	m.g23 = int32(86104)
	m.g24 = int32(31497)
	m.g25 = int32(51)
	m.g26 = int32(52)
	m.g27 = int32(93920)
	m.g28 = int32(94024)
	m.g29 = int32(31840)
	m.g30 = int32(37824)
	m.g31 = int32(125244)
	m.g32 = int32(62)
	m.g33 = int32(63)
	m.g34 = int32(125428)
	m.g35 = int32(79120)
	m.g36 = int32(94108)
	m.g37 = int32(94828)
	m.g38 = int32(97440)
	m.g39 = int32(139240)
	m.g40 = int32(115)
	m.g41 = int32(139244)
	m.g42 = int32(116)
	m.g43 = int32(117)
	m.g44 = int32(118)
	m.g45 = int32(126616)
	m.g46 = int32(126612)
	m.g47 = int32(126620)
	m.g48 = int32(126740)
	m.g49 = int32(126744)
	m.g50 = int32(126748)
	m.g51 = int32(126756)
	m.g52 = int32(96812)
	m.g53 = int32(126620)
	m.g54 = int32(127200)
	m.g55 = int32(127552)
	m.g56 = int32(97616)
	m.g57 = int32(97604)
	initElem0_0(m)
	m.dataEnd = 139320
	initData_0(m)
	fn46(m)
	return m
}

// NewWithWASI constructs a *Module with a custom
// wasi_snapshot_preview1 implementation and a default initial
// linear-memory reservation. Use NewWithWASIReserve to pre-size
// the reservation (e.g. to cover an interpreter's whole boot and
// avoid reallocating/copying linear memory on the first grow).
func NewWithWASI(wasi_snapshot_preview1 Wasi_snapshot_preview1Imports, env EnvImports) *Module {
	return NewWithWASIReserve(wasi_snapshot_preview1, env, 83886080)
}

// New constructs a *Module using DefaultWASI() for the
// wasi_snapshot_preview1 import. Use NewWithWASI to plug in a
// custom implementation (sandboxed FS, captured stdout, ...).
func New(env EnvImports) *Module {
	return NewWithWASI(DefaultWASI(), env)
}

const InitialMemoryBytes = 67108864

func NewWithMemory(wasi_snapshot_preview1 Wasi_snapshot_preview1Imports, env EnvImports, memory []byte, memSize uint64) *Module {
	m := &Module{wasi_snapshot_preview1: wasi_snapshot_preview1, env: env}
	m.memory = memory
	m.memMu = &sync.Mutex{}
	m.memSize = &atomic.Uint64{}
	m.threads = &threadPool{}
	if memSize > 4294836224 {
		panic("wasm2go: memory size exceeds the implementation limit (4294836224 bytes)")
	}
	m.memSize.Store(memSize)
	m.M = unsafe.Pointer(unsafe.SliceData(m.memory))
	m.maxMem = uint64(len(memory))
	m.t0 = make([]any, 143)
	m.g0 = int32(10937088)
	m.g1 = int32(0)
	m.g2 = int32(0)
	m.g3 = int32(10937088)
	m.g4 = int32(0)
	m.g5 = int32(97288)
	m.g6 = int32(97592)
	m.g7 = int32(95896)
	m.g8 = int32(95904)
	m.g9 = int32(95900)
	m.g10 = int32(126752)
	m.g11 = int32(96808)
	m.g12 = int32(85164)
	m.g13 = int32(85240)
	m.g14 = int32(85296)
	m.g15 = int32(85336)
	m.g16 = int32(86180)
	m.g17 = int32(22)
	m.g18 = int32(23)
	m.g19 = int32(29)
	m.g20 = int32(30)
	m.g21 = int32(31)
	m.g22 = int32(85368)
	m.g23 = int32(86104)
	m.g24 = int32(31497)
	m.g25 = int32(51)
	m.g26 = int32(52)
	m.g27 = int32(93920)
	m.g28 = int32(94024)
	m.g29 = int32(31840)
	m.g30 = int32(37824)
	m.g31 = int32(125244)
	m.g32 = int32(62)
	m.g33 = int32(63)
	m.g34 = int32(125428)
	m.g35 = int32(79120)
	m.g36 = int32(94108)
	m.g37 = int32(94828)
	m.g38 = int32(97440)
	m.g39 = int32(139240)
	m.g40 = int32(115)
	m.g41 = int32(139244)
	m.g42 = int32(116)
	m.g43 = int32(117)
	m.g44 = int32(118)
	m.g45 = int32(126616)
	m.g46 = int32(126612)
	m.g47 = int32(126620)
	m.g48 = int32(126740)
	m.g49 = int32(126744)
	m.g50 = int32(126748)
	m.g51 = int32(126756)
	m.g52 = int32(96812)
	m.g53 = int32(126620)
	m.g54 = int32(127200)
	m.g55 = int32(127552)
	m.g56 = int32(97616)
	m.g57 = int32(97604)
	initElem0_0(m)
	m.dataEnd = 139320
	fn46(m)
	return m
}
func NewFromSnapshot(wasi_snapshot_preview1 Wasi_snapshot_preview1Imports, env EnvImports, memory []byte, memSize uint64, globals []uint64) *Module {
	m := &Module{wasi_snapshot_preview1: wasi_snapshot_preview1, env: env}
	m.memory = memory
	m.memMu = &sync.Mutex{}
	m.memSize = &atomic.Uint64{}
	m.threads = &threadPool{}
	if memSize > 4294836224 {
		panic("wasm2go: memory size exceeds the implementation limit (4294836224 bytes)")
	}
	m.memSize.Store(memSize)
	m.M = unsafe.Pointer(unsafe.SliceData(m.memory))
	m.maxMem = uint64(len(memory))
	m.t0 = make([]any, 143)
	m.g0 = int32(10937088)
	m.g1 = int32(0)
	m.g2 = int32(0)
	m.g3 = int32(10937088)
	m.g4 = int32(0)
	m.g5 = int32(97288)
	m.g6 = int32(97592)
	m.g7 = int32(95896)
	m.g8 = int32(95904)
	m.g9 = int32(95900)
	m.g10 = int32(126752)
	m.g11 = int32(96808)
	m.g12 = int32(85164)
	m.g13 = int32(85240)
	m.g14 = int32(85296)
	m.g15 = int32(85336)
	m.g16 = int32(86180)
	m.g17 = int32(22)
	m.g18 = int32(23)
	m.g19 = int32(29)
	m.g20 = int32(30)
	m.g21 = int32(31)
	m.g22 = int32(85368)
	m.g23 = int32(86104)
	m.g24 = int32(31497)
	m.g25 = int32(51)
	m.g26 = int32(52)
	m.g27 = int32(93920)
	m.g28 = int32(94024)
	m.g29 = int32(31840)
	m.g30 = int32(37824)
	m.g31 = int32(125244)
	m.g32 = int32(62)
	m.g33 = int32(63)
	m.g34 = int32(125428)
	m.g35 = int32(79120)
	m.g36 = int32(94108)
	m.g37 = int32(94828)
	m.g38 = int32(97440)
	m.g39 = int32(139240)
	m.g40 = int32(115)
	m.g41 = int32(139244)
	m.g42 = int32(116)
	m.g43 = int32(117)
	m.g44 = int32(118)
	m.g45 = int32(126616)
	m.g46 = int32(126612)
	m.g47 = int32(126620)
	m.g48 = int32(126740)
	m.g49 = int32(126744)
	m.g50 = int32(126748)
	m.g51 = int32(126756)
	m.g52 = int32(96812)
	m.g53 = int32(126620)
	m.g54 = int32(127200)
	m.g55 = int32(127552)
	m.g56 = int32(97616)
	m.g57 = int32(97604)
	initElem0_0(m)
	m.dataEnd = 139320
	restoreGlobals(m, globals)
	return m
}
func initElem0_0(m *Module) {
	m.t0[0] = fn54
	m.t0[1] = fn66
	m.t0[2] = fn85
	m.t0[3] = fn86
	m.t0[4] = fn87
	m.t0[5] = fn93
	m.t0[6] = fn85
	m.t0[7] = fn86
	m.t0[8] = fn277
	m.t0[9] = fn94
	m.t0[10] = fn85
	m.t0[11] = fn86
	m.t0[12] = fn111
	m.t0[13] = fn112
	m.t0[14] = fn87
	m.t0[15] = fn110
	m.t0[16] = fn85
	m.t0[17] = fn86
	m.t0[18] = fn114
	m.t0[19] = fn115
	m.t0[20] = fn117
	m.t0[21] = fn118
	m.t0[22] = fn188
	m.t0[23] = fn185
	m.t0[24] = fn123
	m.t0[25] = fn124
	m.t0[26] = fn152
	m.t0[27] = fn136
	m.t0[28] = fn137
	m.t0[29] = fn190
	m.t0[30] = fn187
	m.t0[31] = fn184
	m.t0[32] = fn139
	m.t0[33] = fn140
	m.t0[34] = fn139
	m.t0[35] = fn139
	m.t0[36] = fn129
	m.t0[37] = fn130
	m.t0[38] = fn128
	m.t0[39] = fn143
	m.t0[40] = fn144
	m.t0[41] = fn146
	m.t0[42] = fn147
	m.t0[43] = fn148
	m.t0[44] = fn149
	m.t0[45] = fn353
	m.t0[46] = fn150
	m.t0[47] = fn151
	m.t0[48] = fn355
	m.t0[49] = fn356
	m.t0[50] = fn86
	m.t0[51] = fn278
	m.t0[52] = fn123
	m.t0[53] = fn158
	m.t0[54] = fn167
	m.t0[55] = fn203
	m.t0[56] = fn219
	m.t0[57] = fn220
	m.t0[58] = fn199
	m.t0[59] = fn221
	m.t0[60] = fn222
	m.t0[61] = fn224
	m.t0[62] = fn189
	m.t0[63] = fn186
	m.t0[64] = fn247
	m.t0[65] = fn245
	m.t0[66] = fn246
	m.t0[67] = fn264
	m.t0[68] = fn265
	m.t0[69] = fn252
	m.t0[70] = fn272
	m.t0[71] = fn275
	m.t0[72] = fn276
	m.t0[73] = fn277
	m.t0[74] = fn263
	m.t0[75] = fn269
	m.t0[76] = fn270
	m.t0[77] = fn271
	m.t0[78] = fn298
	m.t0[79] = fn299
	m.t0[80] = fn300
	m.t0[81] = fn306
	m.t0[82] = fn307
	m.t0[83] = fn308
	m.t0[84] = fn307
	m.t0[85] = fn317
	m.t0[86] = fn321
	m.t0[87] = fn323
	m.t0[88] = fn85
	m.t0[89] = fn86
	m.t0[90] = fn327
	m.t0[91] = fn328
	m.t0[92] = fn329
	m.t0[93] = fn330
	m.t0[94] = fn331
	m.t0[95] = fn332
	m.t0[96] = fn333
	m.t0[97] = fn334
	m.t0[98] = fn335
	m.t0[99] = fn336
	m.t0[100] = fn337
	m.t0[101] = fn339
	m.t0[102] = fn340
	m.t0[103] = fn350
	m.t0[104] = fn351
	m.t0[105] = fn277
	m.t0[106] = fn566
	m.t0[107] = fn352
	m.t0[108] = fn354
	m.t0[109] = fn374
	m.t0[110] = fn365
	m.t0[111] = fn366
	m.t0[112] = fn363
	m.t0[113] = fn404
	m.t0[114] = fn406
	m.t0[115] = fn462
	m.t0[116] = fn524
	m.t0[117] = fn415
	m.t0[118] = fn460
	m.t0[119] = fn431
	m.t0[120] = fn465
	m.t0[121] = fn466
	m.t0[122] = fn467
	m.t0[123] = fn468
	m.t0[124] = fn498
	m.t0[125] = fn497
	m.t0[126] = fn277
	m.t0[127] = fn510
	m.t0[128] = fn548
	m.t0[129] = fn549
	m.t0[130] = fn551
	m.t0[131] = fn85
	m.t0[132] = fn567
	m.t0[133] = fn87
	m.t0[134] = fn87
	m.t0[135] = fn569
	m.t0[136] = fn578
	m.t0[137] = fn576
	m.t0[138] = fn571
	m.t0[139] = fn567
	m.t0[140] = fn577
	m.t0[141] = fn575
	m.t0[142] = fn572
}
func initData_0(m *Module) {
	copy(m.memory[0:], wasm2goData_data_bin[0:139320])
}

var _consts = [3224]uintptr{127032, 126960, 83904, 83908, 83912, 83916, 83920, 83924, 83928, 83932, 83936, 83940, 83944, 83948, 83952, 83956, 83960, 83964, 83968, 83972, 83976, 83980, 83984, 83988, 83992, 83996, 84000, 84016, 84032, 84048, 84064, 84080, 84096, 84112, 84128, 84144, 84160, 84176, 84192, 84208, 84224, 84240, 84256, 84272, 84288, 84304, 84320, 84336, 84352, 84368, 84384, 84400, 84416, 84432, 84448, 84464, 84480, 84496, 84512, 84528, 84544, 84560, 84576, 84592, 84608, 84624, 84672, 84676, 84680, 84684, 84688, 84692, 84696, 84700, 84704, 84708, 84712, 84716, 84720, 84724, 84728, 84732, 84736, 84740, 84744, 84748, 84752, 84756, 84760, 84764, 84768, 84772, 84776, 84780, 84784, 84788, 84792, 84796, 84800, 84804, 84808, 84812, 84816, 84820, 84824, 84828, 84832, 84836, 84840, 84844, 84848, 84852, 84856, 84860, 84864, 84868, 84872, 84876, 84880, 84884, 84888, 84892, 84896, 84900, 84904, 84908, 84912, 84916, 84920, 84924, 84928, 84932, 84936, 84940, 84944, 84948, 84952, 84956, 84960, 84964, 84968, 84972, 84976, 84980, 84984, 84988, 84992, 84996, 85000, 85004, 85008, 85012, 85016, 85020, 85024, 85028, 85032, 85036, 85040, 85044, 85048, 85052, 85056, 85060, 85064, 85068, 85072, 85076, 85080, 85084, 85088, 85092, 85096, 85100, 85104, 85108, 85112, 85116, 85120, 85124, 85128, 85132, 85136, 85140, 85168, 85172, 85176, 85180, 85184, 85188, 85192, 85196, 85200, 85208, 85212, 85216, 85220, 85224, 85228, 85232, 85236, 85244, 85248, 85252, 85256, 85260, 85264, 85268, 85272, 85276, 85280, 85284, 85288, 85292, 85300, 85304, 85308, 85312, 85316, 85320, 85324, 85328, 85332, 85340, 85344, 85348, 85352, 85356, 85360, 85364, 85372, 85376, 85380, 85384, 85392, 85396, 85400, 85404, 85408, 85412, 85416, 85420, 85424, 85428, 85432, 85436, 85440, 85444, 85448, 85452, 85456, 85460, 85464, 85468, 85472, 85476, 85480, 85484, 85488, 85492, 85496, 85500, 85504, 85508, 85512, 85516, 85520, 85524, 85528, 85532, 85536, 85540, 85544, 85548, 85552, 85556, 85560, 85564, 85568, 85572, 85576, 85580, 85584, 85588, 85592, 85596, 85600, 85604, 85608, 85612, 85616, 85620, 85624, 85628, 85632, 85636, 85640, 85644, 85648, 85652, 85656, 85660, 85664, 85668, 85672, 85676, 85680, 85684, 85688, 85692, 85696, 85700, 85704, 85708, 85712, 85716, 85720, 85724, 85728, 85732, 85736, 85740, 85744, 85748, 85752, 85756, 85760, 85764, 85768, 85772, 85776, 85780, 85784, 85788, 85792, 85796, 85800, 85804, 85808, 85812, 85816, 85820, 85824, 85828, 85832, 85836, 85840, 85844, 85848, 85852, 85856, 85860, 85864, 85868, 85872, 85876, 85880, 85884, 85888, 85892, 85896, 85900, 85904, 85908, 85912, 85916, 85920, 85924, 85928, 85932, 85936, 85940, 85944, 85948, 85952, 85956, 85960, 85964, 85968, 85972, 85976, 85980, 85984, 85988, 85992, 85996, 86000, 86004, 86008, 86012, 86016, 86020, 86024, 86028, 86032, 86036, 86040, 86044, 86048, 86052, 86056, 86060, 86064, 86068, 86072, 86076, 86080, 86084, 86088, 86092, 86096, 86100, 86108, 86112, 86116, 86120, 86124, 86128, 86132, 86136, 86140, 86144, 86148, 86152, 86156, 86160, 86164, 86168, 86172, 86176, 86184, 86188, 86192, 86196, 86200, 86204, 86208, 86212, 86216, 86224, 86228, 86232, 86236, 86240, 86244, 86248, 86252, 86256, 86260, 86264, 86268, 86272, 86276, 86280, 86284, 86288, 86292, 86296, 86300, 86304, 86308, 86312, 86316, 86320, 86324, 86328, 86332, 86336, 86340, 86352, 86356, 86360, 86364, 86368, 86384, 86388, 86392, 86396, 86400, 86404, 86408, 86412, 86416, 86420, 86424, 86428, 86432, 86436, 86440, 86444, 86448, 86452, 86456, 86460, 86464, 86468, 86472, 86476, 86480, 86484, 86488, 86492, 86496, 86500, 86504, 86508, 86512, 86516, 86520, 86524, 86528, 86532, 86536, 86540, 86544, 86548, 86552, 86556, 86560, 86564, 86568, 86572, 86576, 86580, 86584, 86588, 86592, 86596, 86600, 86604, 86608, 86612, 86616, 86620, 86624, 86628, 86632, 86636, 86640, 86644, 86648, 86652, 86656, 86660, 86664, 86668, 86672, 86676, 86680, 86684, 86688, 86692, 86696, 86700, 86704, 86708, 86712, 86716, 86720, 86724, 86728, 86732, 86736, 86740, 86744, 86748, 86752, 86756, 86760, 86764, 86768, 86772, 86776, 86780, 86784, 86788, 86792, 86796, 86800, 86804, 86808, 86812, 86816, 86820, 86824, 86828, 86832, 86836, 86840, 86844, 86848, 86852, 86856, 86860, 86864, 86868, 86872, 86876, 86880, 86884, 86888, 86892, 86896, 86900, 86904, 86908, 86912, 86916, 86920, 86924, 86928, 86932, 86936, 86940, 86944, 86948, 86952, 86956, 86960, 86964, 86968, 86972, 86976, 86980, 86984, 86988, 86992, 86996, 87000, 87004, 87008, 87012, 87016, 87020, 87024, 87028, 87032, 87036, 87040, 87044, 87048, 87052, 87056, 87060, 87064, 87068, 87072, 87076, 87080, 87084, 87088, 87092, 87096, 87100, 87104, 87108, 87112, 87116, 87120, 87124, 87128, 87132, 87136, 87140, 87144, 87148, 87152, 87156, 87160, 87164, 87168, 87172, 87176, 87180, 87184, 87188, 87192, 87196, 87200, 87204, 87208, 87212, 87216, 87220, 87224, 87228, 87232, 87236, 87240, 87244, 87248, 87252, 87256, 87260, 87264, 87268, 87272, 87276, 87280, 87284, 87288, 87292, 87296, 87300, 87304, 87308, 87312, 87316, 87320, 87324, 87328, 87332, 87336, 87340, 87344, 87348, 87352, 87356, 87360, 87364, 87368, 87372, 87376, 87380, 87384, 87388, 87392, 87396, 87400, 87404, 87408, 87412, 87416, 87420, 87424, 87428, 87432, 87436, 87440, 87444, 87448, 87452, 87456, 87460, 87464, 87468, 87472, 87476, 87480, 87484, 87488, 87492, 87496, 87500, 87504, 87508, 87512, 87516, 87520, 87524, 87528, 87532, 87536, 87540, 87544, 87548, 87552, 87556, 87560, 87564, 87568, 87572, 87576, 87580, 87584, 87588, 87592, 87596, 87600, 87604, 87608, 87612, 87616, 87620, 87624, 87628, 87632, 87636, 87640, 87644, 87648, 87652, 87656, 87660, 87664, 87668, 87672, 87676, 87680, 87684, 87688, 87692, 87696, 87700, 87704, 87708, 87712, 87716, 87720, 87724, 87728, 87732, 87736, 87740, 87744, 87748, 87752, 87756, 87760, 87764, 87768, 87772, 87776, 87780, 87784, 87788, 87792, 87796, 87800, 87804, 87808, 87812, 87816, 87820, 87824, 87828, 87832, 87836, 87840, 87844, 87848, 87852, 87856, 87860, 87864, 87868, 87872, 87876, 87880, 87884, 87888, 87892, 87896, 87900, 87904, 87908, 87912, 87916, 87920, 87924, 87928, 87932, 87936, 87940, 87944, 87948, 87952, 87956, 87960, 87964, 87968, 87972, 87976, 87980, 87984, 87988, 87992, 87996, 88000, 88004, 88008, 88012, 88016, 88020, 88024, 88028, 88032, 88036, 88040, 88044, 88048, 88052, 88056, 88060, 88064, 88068, 88072, 88076, 88080, 88084, 88088, 88092, 88096, 88100, 88104, 88108, 88112, 88116, 88120, 88124, 88128, 88132, 88136, 88140, 88144, 88148, 88152, 88156, 88160, 88164, 88168, 88172, 88176, 88180, 88184, 88188, 88192, 88196, 88200, 88204, 88208, 88212, 88216, 88220, 88224, 88228, 88232, 88236, 88240, 88244, 88248, 88252, 88256, 88260, 88264, 88268, 88272, 88276, 88280, 88284, 88288, 88292, 88296, 88300, 88304, 88308, 88312, 88316, 88320, 88324, 88328, 88332, 88336, 88340, 88344, 88348, 88352, 88356, 88360, 88364, 88368, 88372, 88376, 88380, 88384, 88388, 88392, 88396, 88400, 88404, 88408, 88412, 88416, 88420, 88424, 88428, 88432, 88436, 88440, 88444, 88448, 88452, 88456, 88460, 88464, 88468, 88472, 88476, 88480, 88484, 88488, 88492, 88496, 88500, 88504, 88508, 88512, 88516, 88520, 88524, 88528, 88532, 88536, 88540, 88544, 88548, 88552, 88556, 88560, 88564, 88568, 88572, 88576, 88580, 88584, 88588, 88592, 88596, 88600, 88604, 88608, 88612, 88616, 88620, 88624, 88628, 88632, 88636, 88640, 88644, 88648, 88652, 88656, 88660, 88664, 88668, 88672, 88676, 88680, 88684, 88688, 88692, 88696, 88700, 88704, 88708, 88712, 88716, 88720, 88724, 88728, 88732, 88736, 88740, 88744, 88748, 88752, 88756, 88760, 88764, 88768, 88772, 88776, 88780, 88784, 88792, 88796, 88800, 88804, 88808, 88812, 88816, 88820, 88832, 88836, 88840, 88844, 88848, 88852, 88856, 88860, 88864, 88868, 88872, 88876, 88880, 88884, 88888, 88892, 88896, 88900, 88904, 88908, 88912, 88916, 88920, 88924, 88928, 88932, 88936, 88940, 88944, 88948, 88952, 88956, 88960, 88964, 88968, 88972, 88976, 88980, 88984, 88988, 88992, 88996, 89000, 89004, 89008, 89012, 89016, 89020, 89024, 89028, 89032, 89036, 89040, 89044, 89048, 89052, 89056, 89060, 89064, 89068, 89072, 89076, 89080, 89084, 89088, 89092, 89096, 89100, 89104, 89108, 89112, 89116, 89120, 89124, 89128, 89132, 89136, 89140, 89144, 89148, 89152, 89156, 89160, 89164, 89168, 89172, 89176, 89180, 89184, 89188, 89192, 89196, 89200, 89204, 89208, 89212, 89216, 89220, 89224, 89228, 89232, 89236, 89240, 89244, 89248, 89252, 89256, 89260, 89264, 89268, 89272, 89276, 89280, 89284, 89288, 89292, 89296, 89300, 89304, 89308, 89312, 89316, 89320, 89324, 89328, 89332, 89336, 89340, 89344, 89348, 89352, 89356, 89360, 89364, 89368, 89372, 89376, 89380, 89384, 89388, 89392, 89396, 89400, 89404, 89408, 89412, 89416, 89420, 89424, 89428, 89432, 89436, 89440, 89444, 89448, 89452, 89456, 89460, 89464, 89468, 89472, 89476, 89480, 89484, 89488, 89492, 89496, 89500, 89504, 89508, 89512, 89516, 89520, 89524, 89528, 89532, 89536, 89540, 89544, 89548, 89552, 89556, 89560, 89564, 89568, 89572, 89576, 89580, 89584, 89588, 89592, 89596, 89600, 89604, 89608, 89612, 89616, 89620, 89624, 89628, 89632, 89636, 89640, 89644, 89648, 89652, 89656, 89660, 89664, 89668, 89672, 89676, 89680, 89684, 89688, 89692, 89696, 89700, 89704, 89708, 89712, 89716, 89720, 89724, 89728, 89732, 89736, 89740, 89744, 89748, 89752, 89756, 89760, 89764, 89768, 89772, 89776, 89780, 89784, 89788, 89792, 89796, 89800, 89804, 89808, 89812, 89816, 89820, 89824, 89828, 89832, 89836, 89840, 89844, 89848, 89852, 89856, 89860, 89864, 89868, 89872, 89876, 89880, 89884, 89888, 89892, 89896, 89900, 89904, 89908, 89912, 89916, 89920, 89924, 89928, 89932, 89936, 89940, 89944, 89948, 89952, 89956, 89960, 89964, 89968, 89972, 89976, 89980, 89984, 89988, 89992, 89996, 90000, 90004, 90008, 90012, 90016, 90020, 90024, 90028, 90032, 90036, 90040, 90044, 90048, 90052, 90056, 90060, 90064, 90068, 90072, 90076, 90080, 90084, 90088, 90092, 90096, 90100, 90104, 90108, 90112, 90116, 90120, 90124, 90128, 90132, 90136, 90140, 90144, 90148, 90152, 90156, 90160, 90164, 90168, 90172, 90176, 90180, 90184, 90188, 90192, 90196, 90200, 90204, 90208, 90212, 90216, 90220, 90224, 90228, 90232, 90236, 90240, 90244, 90248, 90252, 90256, 90260, 90264, 90268, 90272, 90276, 90280, 90284, 90288, 90292, 90296, 90300, 90304, 90308, 90312, 90316, 90320, 90324, 90328, 90332, 90336, 90340, 90344, 90348, 90352, 90356, 90360, 90364, 90368, 90372, 90376, 90380, 90384, 90388, 90392, 90396, 90400, 90404, 90408, 90412, 90416, 90420, 90424, 90428, 90432, 90436, 90440, 90444, 90448, 90452, 90456, 90460, 90464, 90468, 90472, 90476, 90480, 90484, 90488, 90492, 90496, 90500, 90504, 90508, 90512, 90516, 90520, 90524, 90528, 90532, 90536, 90540, 90544, 90548, 90552, 90556, 90560, 90564, 90568, 90572, 90576, 90580, 90584, 90588, 90592, 90596, 90600, 90604, 90608, 90612, 90616, 90620, 90624, 90628, 90632, 90636, 90640, 90644, 90648, 90652, 90656, 90660, 90664, 90668, 90672, 90676, 90680, 90684, 90688, 90692, 90696, 90700, 90704, 90708, 90712, 90716, 90720, 90724, 90728, 90732, 90736, 90740, 90744, 90748, 90752, 90756, 90760, 90764, 90768, 90772, 90776, 90780, 90784, 90788, 90792, 90796, 90800, 90804, 90808, 90812, 90816, 90820, 90824, 90828, 90832, 90836, 90840, 90844, 90848, 90852, 90856, 90860, 90864, 90868, 90872, 90876, 90880, 90884, 90888, 90892, 90896, 90900, 90904, 90908, 90912, 90916, 90920, 90924, 90928, 90932, 90936, 90940, 90944, 90948, 90952, 90956, 90960, 90964, 90968, 90972, 90976, 90980, 90984, 90988, 90992, 90996, 91000, 91004, 91008, 91012, 91016, 91020, 91024, 91028, 91032, 91036, 91040, 91044, 91048, 91052, 91056, 91060, 91064, 91068, 91072, 91076, 91080, 91084, 91088, 91092, 91096, 91100, 91104, 91108, 91112, 91116, 91120, 91124, 91128, 91132, 91136, 91140, 91144, 91148, 91152, 91156, 91160, 91164, 91168, 91172, 91176, 91180, 91184, 91188, 91192, 91196, 91200, 91204, 91208, 91212, 91216, 91220, 91224, 91228, 91232, 91240, 91244, 91248, 91252, 91256, 91260, 91264, 91268, 91280, 91284, 91288, 91292, 91296, 91300, 91304, 91308, 91312, 91316, 91320, 91324, 91328, 91332, 91336, 91340, 91344, 91348, 91352, 91356, 91360, 91364, 91368, 91372, 91376, 91380, 91384, 91388, 91392, 91396, 91400, 91404, 91408, 91412, 91416, 91420, 91424, 91428, 91432, 91436, 91440, 91444, 91448, 91452, 91456, 91460, 91464, 91468, 91472, 91476, 91480, 91484, 91488, 91492, 91496, 91500, 91504, 91508, 91512, 91516, 91520, 91524, 91528, 91532, 91536, 91540, 91544, 91548, 91552, 91556, 91560, 91564, 91568, 91572, 91576, 91580, 91584, 91588, 91592, 91596, 91600, 91604, 91608, 91612, 91616, 91620, 91624, 91628, 91632, 91636, 91640, 91644, 91648, 91652, 91656, 91660, 91664, 91668, 91672, 91676, 91680, 91684, 91688, 91692, 91696, 91700, 91704, 91708, 91712, 91716, 91720, 91724, 91728, 91732, 91736, 91740, 91744, 91748, 91752, 91756, 91760, 91764, 91768, 91772, 91776, 91780, 91784, 91788, 91792, 91796, 91800, 91804, 91808, 91812, 91816, 91820, 91824, 91828, 91832, 91836, 91840, 91844, 91848, 91852, 91856, 91860, 91864, 91868, 91872, 91876, 91880, 91884, 91888, 91892, 91896, 91900, 91904, 91908, 91912, 91916, 91920, 91924, 91928, 91932, 91936, 91940, 91944, 91948, 91952, 91956, 91960, 91964, 91968, 91972, 91976, 91980, 91984, 91988, 91992, 91996, 92000, 92004, 92008, 92012, 92016, 92020, 92024, 92028, 92032, 92036, 92040, 92044, 92048, 92052, 92056, 92060, 92064, 92068, 92072, 92076, 92080, 92084, 92088, 92092, 92096, 92100, 92104, 92108, 92112, 92116, 92120, 92124, 92128, 92132, 92136, 92140, 92144, 92148, 92152, 92156, 92160, 92164, 92168, 92172, 92176, 92180, 92184, 92188, 92192, 92196, 92200, 92204, 92208, 92212, 92216, 92220, 92224, 92228, 92232, 92236, 92240, 92244, 92248, 92252, 92256, 92260, 92264, 92268, 92272, 92276, 92280, 92284, 92288, 92292, 92300, 92304, 92308, 92312, 92316, 92320, 92324, 92328, 92332, 92336, 92352, 92356, 92360, 92364, 92368, 92372, 92376, 92380, 92384, 92388, 92392, 92396, 92400, 92404, 92408, 92412, 92416, 92420, 92424, 92428, 92432, 92436, 92440, 92444, 92448, 92452, 92456, 92460, 92464, 92468, 92472, 92476, 92480, 92484, 92488, 92492, 92496, 92500, 92504, 92508, 92512, 92516, 92520, 92524, 92528, 92532, 92536, 92540, 92544, 92548, 92552, 92556, 92560, 92564, 92568, 92572, 92576, 92580, 92584, 92588, 92592, 92596, 92600, 92604, 92608, 92612, 92616, 92620, 92624, 92628, 92632, 92636, 92640, 92644, 92648, 92652, 92656, 92660, 92664, 92668, 92672, 92676, 92680, 92684, 92688, 92692, 92696, 92700, 92704, 92708, 92712, 92716, 92720, 92724, 92728, 92732, 92736, 92740, 92744, 92748, 92752, 92756, 92760, 92764, 92768, 92772, 92776, 92780, 92784, 92788, 92792, 92796, 92800, 92804, 92808, 92812, 92816, 92820, 92824, 92828, 92832, 92836, 92840, 92844, 92848, 92852, 92856, 92860, 92864, 92868, 92872, 92876, 92880, 92884, 92888, 92892, 92896, 92900, 92904, 92908, 92912, 92916, 92920, 92924, 92928, 92932, 92936, 92940, 92944, 92948, 92952, 92956, 92960, 92964, 92968, 92972, 92976, 92980, 92984, 92988, 92992, 92996, 93000, 93004, 93008, 93012, 93016, 93020, 93024, 93028, 93032, 93036, 93040, 93044, 93048, 93052, 93056, 93060, 93064, 93068, 93072, 93076, 93080, 93084, 93088, 93092, 93096, 93100, 93104, 93108, 93112, 93116, 93120, 93124, 93128, 93132, 93136, 93140, 93144, 93148, 93152, 93156, 93160, 93164, 93168, 93172, 93176, 93180, 93184, 93188, 93192, 93196, 93200, 93204, 93208, 93212, 93216, 93220, 93224, 93228, 93232, 93236, 93240, 93244, 93248, 93252, 93256, 93260, 93264, 93268, 93272, 93276, 93280, 93284, 93288, 93292, 93296, 93300, 93304, 93308, 93312, 93316, 93320, 93324, 93328, 93332, 93336, 93340, 93344, 93348, 93352, 93356, 93360, 93364, 93372, 93376, 93380, 93384, 93388, 93392, 93396, 93400, 93404, 93408, 93424, 93428, 93432, 93436, 93440, 93444, 93448, 93452, 93456, 93460, 93464, 93468, 93472, 93476, 93480, 93484, 93488, 93492, 93496, 93500, 93504, 93508, 93512, 93516, 93520, 93524, 93528, 93532, 93536, 93540, 93544, 93548, 93552, 93556, 93560, 93564, 93568, 93572, 93576, 93580, 93584, 93588, 93592, 93596, 93600, 93604, 93608, 93612, 93616, 93620, 93624, 93628, 93632, 93636, 93640, 93644, 93648, 93652, 93656, 93660, 93664, 93668, 93672, 93676, 93680, 93684, 93688, 93692, 93696, 93700, 93712, 93716, 93720, 93724, 93728, 93732, 93736, 93740, 93744, 93748, 93752, 93756, 93760, 93764, 93768, 93772, 93776, 93780, 93784, 93788, 93792, 93796, 93800, 93804, 93808, 93812, 93816, 93820, 93824, 93828, 93832, 93836, 93840, 93844, 93848, 93852, 93856, 93860, 93864, 93868, 93872, 93876, 93880, 93884, 93888, 93892, 93896, 93900, 93904, 93908, 93912, 93916, 93924, 93928, 93932, 93936, 93940, 93944, 93948, 93952, 93956, 93960, 93964, 93968, 93972, 93976, 93980, 93984, 93988, 93992, 93996, 94000, 94004, 94008, 94012, 94016, 94020, 94028, 94032, 94036, 94040, 94044, 94048, 94052, 94056, 94060, 94064, 94068, 94072, 94076, 94080, 94084, 94088, 94092, 94096, 94100, 94104, 94112, 94116, 94120, 94124, 94128, 94132, 94136, 94140, 94144, 94148, 94152, 94156, 94160, 94164, 94168, 94176, 94180, 94184, 94188, 94192, 94196, 94200, 94204, 94208, 94212, 94216, 94220, 94224, 94228, 94232, 94236, 94240, 94244, 94248, 94252, 94256, 94260, 94264, 94268, 94272, 94276, 94280, 94284, 94288, 94292, 94296, 94300, 94304, 94308, 94312, 94316, 94320, 94324, 94328, 94332, 94336, 94352, 94356, 94360, 94364, 94368, 94372, 94376, 94380, 94384, 94388, 94392, 94396, 94400, 94404, 94408, 94412, 94416, 94420, 94424, 94428, 94432, 94436, 94440, 94444, 94448, 94452, 94456, 94460, 94464, 94468, 94472, 94476, 94480, 94484, 94488, 94496, 94500, 94504, 94508, 94512, 94516, 94520, 94524, 94528, 94532, 94536, 94540, 94544, 94548, 94552, 94556, 94560, 94564, 94568, 94572, 94576, 94580, 94584, 94588, 94592, 94596, 94600, 94604, 94608, 94612, 94616, 94620, 94624, 94628, 94632, 94636, 94640, 94644, 94648, 94652, 94656, 94660, 94664, 94668, 94672, 94676, 94688, 94692, 94696, 94700, 94704, 94708, 94712, 94716, 94720, 94724, 94728, 94732, 94736, 94740, 94744, 94748, 94752, 94756, 94760, 94764, 94768, 94772, 94784, 94788, 94792, 94796, 94800, 94804, 94808, 94812, 94816, 94820, 94824, 94832, 94836, 94840, 94844, 94848, 94852, 94856, 94912, 94920, 94928, 94936, 94944, 94952, 94960, 94968, 94976, 94984, 94992, 95000, 95008, 95016, 95024, 95032, 95040, 95048, 95056, 95064, 95072, 95080, 95088, 95096, 95104, 95112, 95120, 95128, 95136, 95144, 95152, 95160, 95168, 95176, 95184, 95192, 95200, 95208, 95216, 95224, 95232, 95240, 95248, 95256, 95264, 95272, 95280, 95288, 95296, 95304, 95312, 95320, 95328, 95336, 95344, 95352, 95360, 95368, 95376, 95384, 95392, 95400, 95408, 95416, 95424, 95432, 95440, 95448, 95456, 95464, 95472, 95480, 95488, 95496, 95504, 95512, 95520, 95528, 95536, 95544, 95552, 95560, 95568, 95576, 95584, 95592, 95600, 95608, 95616, 95624, 95632, 95640, 95648, 95656, 95664, 95672, 95680, 95688, 95696, 95704, 95712, 95720, 95728, 95736, 95744, 95752, 95760, 95768, 95776, 95784, 95792, 95800, 95808, 95816, 95824, 95832, 95840, 95848, 95856, 95864, 95872, 95880, 95888, 95924, 95932, 95940, 95948, 95956, 95964, 95972, 95980, 95988, 95996, 96004, 96012, 96020, 96028, 96036, 96044, 96052, 96060, 96068, 96076, 96084, 96092, 96100, 96108, 96116, 96124, 96132, 96140, 96148, 96156, 96164, 96172, 96180, 96188, 96196, 96204, 96212, 96220, 96228, 96236, 96244, 96252, 96260, 96268, 96276, 96284, 96292, 96300, 96308, 96316, 96324, 96332, 96340, 96348, 96356, 96364, 96372, 96380, 96388, 96396, 96404, 96412, 96420, 96428, 96436, 96444, 96452, 96460, 96468, 96476, 96484, 96492, 96500, 96508, 96516, 96524, 96532, 96540, 96548, 96556, 96564, 96572, 96580, 96588, 96596, 96604, 96612, 96620, 96628, 96636, 96644, 96652, 96660, 96668, 96676, 96684, 96692, 96700, 96708, 96716, 96724, 96732, 96740, 96748, 96756, 96764, 96772, 96780, 96788, 96796, 96820, 96856, 96884, 96888, 96892, 96896, 96900, 96904, 96908, 96912, 96916, 96920, 96924, 96928, 96932, 96936, 96940, 96976, 96980, 96984, 96988, 96996, 97004, 97156, 97180, 97184, 97188, 97288, 97292, 97308, 97328, 97336, 97340, 97440, 97460, 97484, 97488, 97492, 97592, 97596, 97600, 97604, 97608, 97612, 97616, 97620, 97624, 97632, 97636, 97640, 97644, 97648, 97652, 97656, 97660, 97664, 97672, 97676, 97680, 97684, 97688, 97692, 97696, 97700, 97704, 97708, 97712, 97716, 97720, 97724, 9901, 139240, 126608, 15800, 122448, 15796, 7656, 138760, 138764, 138744, 138752, 138748, 138768, 138756, 138776, 122704, 122708, 127328, 126792, 79904, 79888, 79896, 126808, 126800, 123744, 124935, 124936, 124940, 124944, 124896, 16375, 16376, 16427, 9057, 16428, 9058, 126040, 126036, 126048, 126044, 126052, 126056, 16342, 16343, 6950, 7375, 16596, 16608, 124916, 124932, 124908, 8586, 4109, 95900, 95896, 95904, 124920, 124864, 124912, 124907, 124906, 124905, 124865, 124904, 16251, 124892, 124888, 124900, 124816, 124868, 124784, 124792, 124788, 124804, 124800, 124796, 124808, 124872, 124880, 84656, 83888, 8668, 124921, 124812, 124922, 124924, 124928, 125876, 124933, 16885, 17011, 8041, 122720, 24426, 125872, 125868, 125864, 125860, 125856, 125852, 125848, 124824, 124828, 124832, 124836, 124840, 124844, 124848, 124852, 124856, 124860, 16618, 16537, 9748, 125124, 125140, 125144, 31050, 125132, 124820, 124952, 126088, 126140, 126084, 126092, 126132, 96816, 126204, 85152, 85156, 85160, 124948, 43480, 43520, 43504, 43496, 43528, 43512, 43488, 43508, 43500, 43472, 8103, 8111, 124934, 97528, 97468, 97464, 124876, 124884, 124885, 124877, 54976, 54972, 16024, 23432, 18040, 18024, 18048, 18032, 18036, 18028, 18052, 18044, 125104, 125100, 125080, 125088, 125068, 125084, 125064, 125060, 125096, 125092, 125072, 125076, 16276, 16277, 125108, 125112, 125116, 125120, 125436, 125328, 125424, 125136, 125128, 7365, 125396, 125276, 125260, 126812, 9316, 9310, 9318, 6110, 9032, 9034, 125264, 125376, 125148, 125152, 125304, 125160, 125220, 10151, 125212, 125216, 125312, 125208, 125200, 125196, 125204, 125308, 125236, 125228, 125240, 31051, 31052, 31053, 125232, 8091, 8572, 6808, 125320, 125252, 125268, 125256, 125392, 125536, 125512, 125528, 125520, 125540, 62192, 62176, 62180, 62184, 62188, 62196, 62200, 62204, 62208, 62212, 62216, 62220, 62224, 62228, 62232, 62236, 62240, 62244, 62248, 62252, 62256, 62260, 62264, 62268, 62272, 62276, 62280, 86344, 125400, 125404, 125300, 125408, 125412, 125416, 125420, 7638, 10219, 15270, 14849, 62340, 62996, 63002, 7290, 63015, 62331, 125440, 125464, 125432, 125456, 125448, 125472, 125480, 125488, 125496, 125504, 15196, 15998, 125580, 125584, 125588, 125388, 17226, 125592, 126864, 125600, 125372, 125604, 125608, 125612, 125616, 125620, 125624, 125628, 125632, 125636, 125640, 125644, 125652, 125648, 125836, 125840, 125844, 126016, 126020, 126064, 126068, 126060, 126872, 126924, 126876, 126900, 8066, 126656, 126932, 126761, 97220, 97224, 97164, 97160, 126848, 126856, 4098, 4096, 127336, 127332, 83215, 139216, 139224, 139236, 139228, 139220, 139188, 139184, 139176, 139180, 139204, 139196, 139192, 138780, 139232, 138772, 139200, 139048}

func (m *Module) WasmCallCtors() {
	fn44(m)
}
func (m *Module) Fflush(l0 int32) int32 {
	return fn461(m, l0)
}
func (m *Module) Fopen(l0 int32, l1 int32) int32 {
	return fn469(m, l0, l1)
}
func (m *Module) Fclose(l0 int32) int32 {
	return fn459(m, l0)
}
func (m *Module) MainArgcArgv(l0 int32, l1 int32) int32 {
	return fn65(m, l0, l1)
}
func (m *Module) Malloc(l0 int32) int32 {
	return fn554(m, l0)
}
func (m *Module) PglSetSystemFn(l0 int32) {
	fn377(m, l0)
}
func (m *Module) PglSetPopenFn(l0 int32) {
	fn379(m, l0)
}
func (m *Module) PglSetPcloseFn(l0 int32) {
	fn381(m, l0)
}
func (m *Module) EmscriptenBuiltinMemalign(l0 int32, l1 int32) int32 {
	return fn557(m, l0, l1)
}
func (m *Module) EmscriptenStackAlloc(l0 int32) int32 {
	return fn580(m, l0)
}
func (m *Module) WasmApplyDataRelocs() {
	fn45(m)
}
func (m *Module) Memory() []byte {
	return m.memory
}

// ui32 / ui64 reinterpret a signed integer as its unsigned bit
// equivalent at runtime. Used for the operands of wasm unsigned
// comparisons (i32.lt_u etc.) — emitting `uint32(int32(-N))` directly
// fails Go's compile-time constant rule because the negative typed
// constant isn't representable in uint32; routing through these
// function-call boundaries forces runtime conversion.
func ui32(x int32) uint32 { return uint32(x) }

func ui64(x int64) uint64 { return uint64(x) }

// b2i32 materialises a wasm comparison result — an i32 that is 0 or 1 — from
// the Go bool the comparison expression evaluates to.
//
// It exists as a named helper rather than an inline `func() int32 { ... }()`
// because the gcasm backend requires every direct call left in the compiled
// output to be either a package-local FnN or something the Go inliner removed.
// A func literal is normally inlined at its call site, but the inliner gives up
// once the ENCLOSING function grows past its budget — and a single wasm function
// can translate to tens of thousands of lines of Go, as an interpreter's
// bytecode dispatch loop does. The literal is then outlined into a real closure
// symbol (FnN.funcA.funcB), which reaches the assembler as a direct call gcasm
// cannot marshal. A named helper this small is always inlined, and if it ever
// were not, it would fail loudly at its own symbol rather than as a nested
// closure.
func b2i32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

//go:noinline
func wasm_trap_div_zero() { panic("wasm: integer divide by zero") }

//go:noinline
func wasm_trap_int_overflow() { panic("wasm: integer overflow") }

//go:noinline
func wasm_trap_invalid_conv() { panic("wasm: invalid conversion to integer") }

//go:noinline
func wasm_trap_unreachable() { panic("wasm: unreachable") }

//go:noinline
func wasm_trap_memfill_oob() { panic("wasm: memory.fill out of bounds") }

//go:noinline
func wasm_trap_memcopy_oob() { panic("wasm: memory.copy out of bounds") }

func i32_div_s(x, y int32) int32 {
	if y == -1 && x == math.MinInt32 {
		wasm_trap_int_overflow()
	}
	if y == 0 {
		wasm_trap_div_zero()
	}
	return x / y
}

func i64_div_s(x, y int64) int64 {
	if y == -1 && x == math.MinInt64 {
		wasm_trap_int_overflow()
	}
	if y == 0 {
		wasm_trap_div_zero()
	}
	return x / y
}

func i32_div_u(x, y uint32) uint32 {
	if y == 0 {
		wasm_trap_div_zero()
	}
	return x / y
}

func i64_div_u(x, y uint64) uint64 {
	if y == 0 {
		wasm_trap_div_zero()
	}
	return x / y
}

func i32_rem_s(x, y int32) int32 {
	if y == 0 {
		wasm_trap_div_zero()
	}
	if y == -1 {

		return 0
	}
	return x % y
}

func i64_rem_s(x, y int64) int64 {
	if y == 0 {
		wasm_trap_div_zero()
	}
	if y == -1 {
		return 0
	}
	return x % y
}

func i32_rem_u(x, y uint32) uint32 {
	if y == 0 {
		wasm_trap_div_zero()
	}
	return x % y
}

func i32_rotl(x, y int32) int32 { return int32(bits.RotateLeft32(uint32(x), int(y&31))) }

func i64_rotl(x, y int64) int64 { return int64(bits.RotateLeft64(uint64(x), int(y&63))) }

func f32_abs(x float32) float32 {
	return math.Float32frombits(math.Float32bits(x) &^ (1 << 31))
}

func f64_abs(x float64) float64 {
	return math.Float64frombits(math.Float64bits(x) &^ (1 << 63))
}

func f64_neg(x float64) float64 {
	return math.Float64frombits(math.Float64bits(x) ^ (1 << 63))
}

func i32_trunc_f32_s(x float32) int32 {
	if x != x {
		wasm_trap_invalid_conv()
	}

	if !(x > -2147483904.0 && x < 2147483648.0) {
		wasm_trap_int_overflow()
	}
	return int32(x)
}

func i32_trunc_f64_s(x float64) int32 {
	if x != x {
		wasm_trap_invalid_conv()
	}

	if !(x > -2147483649.0 && x < 2147483648.0) {
		wasm_trap_int_overflow()
	}
	return int32(x)
}

func i32_trunc_f64_u(x float64) int32 {
	if x != x {
		wasm_trap_invalid_conv()
	}
	if !(x > -1.0 && x < 4294967296.0) {
		wasm_trap_int_overflow()
	}
	return int32(uint32(x))
}

func i64_trunc_f64_s(x float64) int64 {
	if x != x {
		wasm_trap_invalid_conv()
	}
	if !(x >= -9223372036854775808.0 && x < 9223372036854775808.0) {
		wasm_trap_int_overflow()
	}
	return int64(x)
}

// memorySize returns the current size of m.memory in wasm pages (each
// page is 64 KiB).
func memorySize(m *Module) int32 {
	return int32(m.memSize.Load() >> 16)
}

// accessMemory runs f with the module's current linear memory while
// holding the same lock memoryGrow takes to mutate the memory slice
// header or relocate its backing array. It is the ONE safe way to
// touch linear memory from OUTSIDE the module's execution goroutine —
// e.g. a watchdog goroutine raising CPython's eval-breaker bit while
// an evaluation is running. For the duration of f the memory can
// neither be resliced nor relocated, so f's writes land in the array
// the guest observes; a grow that raced in just before blocks until f
// returns and then copies f's writes forward with the rest of the
// contents. Determinism notes for callers:
//
//   - f MUST NOT call back into the module or into memoryGrow — that
//     would self-deadlock.
//   - f should be short: a running guest blocks inside memory.grow
//     until f returns (ordinary guest loads/stores do not block).
//   - Bytes the guest reads or writes concurrently with f (that is
//     the point of an eval-breaker-style flag) are exchanged with
//     plain single-word accesses; keep such shared words
//     word-aligned and word-sized.
func accessMemory(m *Module, f func(mem []byte)) {
	m.memMu.Lock()
	defer m.memMu.Unlock()
	f(m.memory)
}

func i32_div_u_s(x, y int32) int32 { return int32(i32_div_u(uint32(x), uint32(y))) }
func i32_rem_u_s(x, y int32) int32 { return int32(i32_rem_u(uint32(x), uint32(y))) }
func i64_div_u_s(x, y int64) int64 { return int64(i64_div_u(uint64(x), uint64(y))) }

func f32_mul(x, y float32) float32 { return float32(x * y) }

func f64_add(x, y float64) float64 { return float64(x + y) }
func f64_sub(x, y float64) float64 { return float64(x - y) }
func f64_mul(x, y float64) float64 { return float64(x * y) }
func f64_div(x, y float64) float64 { return float64(x / y) }

func i32_clz(x int32) int32    { return int32(bits.LeadingZeros32(uint32(x))) }
func i32_ctz(x int32) int32    { return int32(bits.TrailingZeros32(uint32(x))) }
func i32_popcnt(x int32) int32 { return int32(bits.OnesCount32(uint32(x))) }

func f32_lt(x, y float32) int32 {
	if x < y {
		return 1
	}
	return 0
}

func f64_eq(x, y float64) int32 {
	if x == y {
		return 1
	}
	return 0
}
func f64_ne(x, y float64) int32 {
	if x != y {
		return 1
	}
	return 0
}
func f64_lt(x, y float64) int32 {
	if x < y {
		return 1
	}
	return 0
}

func f64_ge(x, y float64) int32 {
	if x >= y {
		return 1
	}
	return 0
}

func i32_wrap_i64(x int64) int32     { return int32(x) }
func i64_extend_i32_s(x int32) int64 { return int64(x) }
func i64_extend_i32_u(x int32) int64 { return int64(uint32(x)) }

func f32_convert_i32_s(x int32) float32 { return float32(x) }

func f64_convert_i32_s(x int32) float64 { return float64(x) }
func f64_convert_i32_u(x int32) float64 { return float64(uint32(x)) }

func i64_reinterpret_f64(x float64) int64 { return int64(math.Float64bits(x)) }

func f64_reinterpret_i64(x int64) float64 { return math.Float64frombits(uint64(x)) }

func i32_extend8_s(x int32) int32  { return int32(int8(x)) }
func i32_extend16_s(x int32) int32 { return int32(int16(x)) }

func memoryFill(m *Module, dst int32, val int32, n int32) {
	if n == 0 {
		return
	}
	end := uint64(uint32(dst)) + uint64(uint32(n))
	if end > m.memSize.Load() {
		wasm_trap_memfill_oob()
	}
	b := m.memory[uint32(dst):uint32(end)]
	v := byte(val)

	if v == 0 {
		for k := range b {
			b[k] = 0
		}
		return
	}
	b[0] = v
	for filled := 1; filled < len(b); filled *= 2 {
		copy(b[filled:], b[:filled])
	}
}

func memoryCopy(m *Module, dst int32, src int32, n int32) {
	if n == 0 {
		return
	}
	srcEnd := uint64(uint32(src)) + uint64(uint32(n))
	dstEnd := uint64(uint32(dst)) + uint64(uint32(n))
	if size := m.memSize.Load(); srcEnd > size || dstEnd > size {
		wasm_trap_memcopy_oob()
	}
	copy(m.memory[uint32(dst):uint32(dstEnd)], m.memory[uint32(src):uint32(srcEnd)])
}

var spinAgents int32
var spinOversubscribed uint32

type threadPool struct {
	nextTID atomic.Int32
	wg      sync.WaitGroup

	parkMu sync.Mutex
	parked map[uint64][]chan struct{}
}

// wake releases up to count waiters on ea and reports how many it woke.
func (p *threadPool) wake(ea uint64, count int32) int32 {
	p.parkMu.Lock()
	defer p.parkMu.Unlock()
	waiters := p.parked[ea]
	n := int32(len(waiters))
	if count >= 0 && count < n {
		n = count
	}
	for _, ch := range waiters[:n] {
		close(ch)
	}
	if int(n) == len(waiters) {
		delete(p.parked, ea)
	} else {
		p.parked[ea] = waiters[n:]
	}
	return n
}

// saveGlobals returns the module's mutable globals, in a form that can be handed back
// to restoreGlobals. It is how a snapshot of an instance captures the state that does not
// live in linear memory.
func saveGlobals(m *Module) []uint64 {
	g := make([]uint64, 56)
	g[0] = uint64(uint32(m.g0))
	g[1] = uint64(uint32(m.g3))
	g[2] = uint64(uint32(m.g4))
	g[3] = uint64(uint32(m.g5))
	g[4] = uint64(uint32(m.g6))
	g[5] = uint64(uint32(m.g7))
	g[6] = uint64(uint32(m.g8))
	g[7] = uint64(uint32(m.g9))
	g[8] = uint64(uint32(m.g10))
	g[9] = uint64(uint32(m.g11))
	g[10] = uint64(uint32(m.g12))
	g[11] = uint64(uint32(m.g13))
	g[12] = uint64(uint32(m.g14))
	g[13] = uint64(uint32(m.g15))
	g[14] = uint64(uint32(m.g16))
	g[15] = uint64(uint32(m.g17))
	g[16] = uint64(uint32(m.g18))
	g[17] = uint64(uint32(m.g19))
	g[18] = uint64(uint32(m.g20))
	g[19] = uint64(uint32(m.g21))
	g[20] = uint64(uint32(m.g22))
	g[21] = uint64(uint32(m.g23))
	g[22] = uint64(uint32(m.g24))
	g[23] = uint64(uint32(m.g25))
	g[24] = uint64(uint32(m.g26))
	g[25] = uint64(uint32(m.g27))
	g[26] = uint64(uint32(m.g28))
	g[27] = uint64(uint32(m.g29))
	g[28] = uint64(uint32(m.g30))
	g[29] = uint64(uint32(m.g31))
	g[30] = uint64(uint32(m.g32))
	g[31] = uint64(uint32(m.g33))
	g[32] = uint64(uint32(m.g34))
	g[33] = uint64(uint32(m.g35))
	g[34] = uint64(uint32(m.g36))
	g[35] = uint64(uint32(m.g37))
	g[36] = uint64(uint32(m.g38))
	g[37] = uint64(uint32(m.g39))
	g[38] = uint64(uint32(m.g40))
	g[39] = uint64(uint32(m.g41))
	g[40] = uint64(uint32(m.g42))
	g[41] = uint64(uint32(m.g43))
	g[42] = uint64(uint32(m.g44))
	g[43] = uint64(uint32(m.g45))
	g[44] = uint64(uint32(m.g46))
	g[45] = uint64(uint32(m.g47))
	g[46] = uint64(uint32(m.g48))
	g[47] = uint64(uint32(m.g49))
	g[48] = uint64(uint32(m.g50))
	g[49] = uint64(uint32(m.g51))
	g[50] = uint64(uint32(m.g52))
	g[51] = uint64(uint32(m.g53))
	g[52] = uint64(uint32(m.g54))
	g[53] = uint64(uint32(m.g55))
	g[54] = uint64(uint32(m.g56))
	g[55] = uint64(uint32(m.g57))
	return g
}

// restoreGlobals puts a snapshot's globals back. A snapshot from a different module (or a
// different build of the same one) has a different global count; rather than
// index out of bounds, take what fits and leave the rest at their declared
// initializers.
func restoreGlobals(m *Module, g []uint64) {
	if len(g) != 56 {
		return
	}
	m.g0 = int32(uint32(g[0]))
	m.g3 = int32(uint32(g[1]))
	m.g4 = int32(uint32(g[2]))
	m.g5 = int32(uint32(g[3]))
	m.g6 = int32(uint32(g[4]))
	m.g7 = int32(uint32(g[5]))
	m.g8 = int32(uint32(g[6]))
	m.g9 = int32(uint32(g[7]))
	m.g10 = int32(uint32(g[8]))
	m.g11 = int32(uint32(g[9]))
	m.g12 = int32(uint32(g[10]))
	m.g13 = int32(uint32(g[11]))
	m.g14 = int32(uint32(g[12]))
	m.g15 = int32(uint32(g[13]))
	m.g16 = int32(uint32(g[14]))
	m.g17 = int32(uint32(g[15]))
	m.g18 = int32(uint32(g[16]))
	m.g19 = int32(uint32(g[17]))
	m.g20 = int32(uint32(g[18]))
	m.g21 = int32(uint32(g[19]))
	m.g22 = int32(uint32(g[20]))
	m.g23 = int32(uint32(g[21]))
	m.g24 = int32(uint32(g[22]))
	m.g25 = int32(uint32(g[23]))
	m.g26 = int32(uint32(g[24]))
	m.g27 = int32(uint32(g[25]))
	m.g28 = int32(uint32(g[26]))
	m.g29 = int32(uint32(g[27]))
	m.g30 = int32(uint32(g[28]))
	m.g31 = int32(uint32(g[29]))
	m.g32 = int32(uint32(g[30]))
	m.g33 = int32(uint32(g[31]))
	m.g34 = int32(uint32(g[32]))
	m.g35 = int32(uint32(g[33]))
	m.g36 = int32(uint32(g[34]))
	m.g37 = int32(uint32(g[35]))
	m.g38 = int32(uint32(g[36]))
	m.g39 = int32(uint32(g[37]))
	m.g40 = int32(uint32(g[38]))
	m.g41 = int32(uint32(g[39]))
	m.g42 = int32(uint32(g[40]))
	m.g43 = int32(uint32(g[41]))
	m.g44 = int32(uint32(g[42]))
	m.g45 = int32(uint32(g[43]))
	m.g46 = int32(uint32(g[44]))
	m.g47 = int32(uint32(g[45]))
	m.g48 = int32(uint32(g[46]))
	m.g49 = int32(uint32(g[47]))
	m.g50 = int32(uint32(g[48]))
	m.g51 = int32(uint32(g[49]))
	m.g52 = int32(uint32(g[50]))
	m.g53 = int32(uint32(g[51]))
	m.g54 = int32(uint32(g[52]))
	m.g55 = int32(uint32(g[53]))
	m.g56 = int32(uint32(g[54]))
	m.g57 = int32(uint32(g[55]))
}

//go:embed data.bin
var wasm2goData_data_bin []byte
