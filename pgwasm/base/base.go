package base

import (
	"context"
	"crypto/rand"
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
	Environ_sizes_get(m *Module, l0 int32, l1 int32) int32
	Environ_get(m *Module, l0 int32, l1 int32) int32
	Proc_exit(m *Module, l0 int32)
	Clock_time_get(m *Module, l0 int32, l1 int64, l2 int32) int32
	Fd_close(m *Module, l0 int32) int32
	Fd_write(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	Fd_read(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	Fd_sync(m *Module, l0 int32) int32
	Fd_fdstat_get(m *Module, l0 int32, l1 int32) int32
	Fd_seek(m *Module, l0 int32, l1 int64, l2 int32, l3 int32) int32
	Fd_pread(m *Module, l0 int32, l1 int32, l2 int32, l3 int64, l4 int32) int32
	Fd_pwrite(m *Module, l0 int32, l1 int32, l2 int32, l3 int64, l4 int32) int32
	Random_get(m *Module, l0 int32, l1 int32) int32
}
type EnvImports interface {
	Invoke_ii(m *Module, l0 int32, l1 int32) int32
	Invoke_vii(m *Module, l0 int32, l1 int32, l2 int32)
	Invoke_viiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32)
	Invoke_vi(m *Module, l0 int32, l1 int32)
	Invoke_v(m *Module, l0 int32)
	Invoke_iii(m *Module, l0 int32, l1 int32, l2 int32) int32
	Invoke_viii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32)
	Invoke_iiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	Invoke_jii(m *Module, l0 int32, l1 int32, l2 int32) int64
	Invoke_iiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) int32
	Invoke_i(m *Module, l0 int32) int32
	Invoke_ji(m *Module, l0 int32, l1 int32) int64
	Invoke_jiiiiiiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32) int64
	Invoke_jiiiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) int64
	Invoke_viiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32)
	Invoke_iiiiiiiiiiiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32, l10 int32, l11 int32, l12 int32, l13 int32) int32
	Invoke_vji(m *Module, l0 int32, l1 int64, l2 int32)
	Invoke_viiji(m *Module, l0 int32, l1 int32, l2 int32, l3 int64, l4 int32)
	Invoke_iiiij(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int64) int32
	Invoke_vijiji(m *Module, l0 int32, l1 int32, l2 int64, l3 int32, l4 int64, l5 int32)
	Invoke_viji(m *Module, l0 int32, l1 int32, l2 int64, l3 int32)
	Invoke_iiiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32
	Invoke_iiiiiiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32) int32
	Invoke_iiiiiiiiiiiiiiiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32, l10 int32, l11 int32, l12 int32, l13 int32, l14 int32, l15 int32, l16 int32, l17 int32) int32
	Invoke_iiiiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) int32
	Invoke_vj(m *Module, l0 int32, l1 int64)
	Invoke_iiiiiiiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32) int32
	Invoke_viij(m *Module, l0 int32, l1 int32, l2 int32, l3 int64)
	Invoke_viiiiiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32)
	Invoke_viiiiiiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32)
	Invoke_vij(m *Module, l0 int32, l1 int32, l2 int64)
	Invoke_viiiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32)
	Invoke_iiji(m *Module, l0 int32, l1 int32, l2 int64, l3 int32) int32
	Invoke_ij(m *Module, l0 int32, l1 int64) int32
	Invoke_viiiiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32)
	Invoke_viiiji(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int64, l5 int32)
	Invoke_iiiiiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32) int32
	Invoke_iiij(m *Module, l0 int32, l1 int32, l2 int32, l3 int64) int32
	Invoke_vid(m *Module, l0 int32, l1 int32, l2 float64)
	Getaddrinfo(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	Invoke_j(m *Module, l0 int32) int64
	Invoke_ijji(m *Module, l0 int32, l1 int64, l2 int64, l3 int32) int32
	Invoke_iijj(m *Module, l0 int32, l1 int32, l2 int64, l3 int64) int32
	Invoke_jiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int64
	Invoke_jij(m *Module, l0 int32, l1 int32, l2 int64) int64
	Invoke_ijiiiiii(m *Module, l0 int32, l1 int64, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32) int32
	Invoke_viijii(m *Module, l0 int32, l1 int32, l2 int32, l3 int64, l4 int32, l5 int32)
	Invoke_iiiiiji(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int64, l6 int32) int32
	Invoke_viijiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int64, l4 int32, l5 int32, l6 int32, l7 int32)
	Invoke_vijjii(m *Module, l0 int32, l1 int32, l2 int64, l3 int64, l4 int32, l5 int32)
	Invoke_vjii(m *Module, l0 int32, l1 int64, l2 int32, l3 int32)
	Invoke_jiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) int64
	Invoke_viiiiiiiiiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32, l10 int32, l11 int32, l12 int32)
	X__assert_fail(m *Module, l0 int32, l1 int32, l2 int32, l3 int32)
	Invoke_di(m *Module, l0 int32, l1 int32) float64
	Invoke_id(m *Module, l0 int32, l1 float64) int32
	Invoke_ijiiiii(m *Module, l0 int32, l1 int64, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) int32
	Invoke_iiiiiiiiiii(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32, l7 int32, l8 int32, l9 int32, l10 int32) int32
	Getnameinfo(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32, l6 int32) int32
	Exit(m *Module, l0 int32)
	X_tzset_js(m *Module, l0 int32, l1 int32, l2 int32, l3 int32)
	X_abort_js(m *Module)
	X__syscall_faccessat(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	X__syscall_chdir(m *Module, l0 int32) int32
	X__syscall_chmod(m *Module, l0 int32, l1 int32) int32
	X__syscall_fchownat(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) int32
	X__syscall_dup(m *Module, l0 int32) int32
	X__syscall_dup3(m *Module, l0 int32, l1 int32, l2 int32) int32
	X_dlopen_js(m *Module, l0 int32) int32
	X_dlsym_js(m *Module, l0 int32, l1 int32, l2 int32) int32
	Emscripten_date_now(m *Module) float64
	X__syscall_fdatasync(m *Module, l0 int32) int32
	X__syscall_openat(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	X__syscall_fcntl64(m *Module, l0 int32, l1 int32, l2 int32) int32
	X__syscall_ioctl(m *Module, l0 int32, l1 int32, l2 int32) int32
	X__syscall_fstat64(m *Module, l0 int32, l1 int32) int32
	X__syscall_stat64(m *Module, l0 int32, l1 int32) int32
	X__syscall_newfstatat(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	X__syscall_lstat64(m *Module, l0 int32, l1 int32) int32
	X__syscall_ftruncate64(m *Module, l0 int32, l1 int64) int32
	X__syscall_getcwd(m *Module, l0 int32, l1 int32) int32
	Emscripten_get_now(m *Module) float64
	X__syscall_mkdirat(m *Module, l0 int32, l1 int32, l2 int32) int32
	X_localtime_js(m *Module, l0 int64, l1 int32)
	X_gmtime_js(m *Module, l0 int64, l1 int32)
	X_munmap_js(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int64) int32
	X_mmap_js(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int64, l5 int32, l6 int32) int32
	X__syscall_pipe(m *Module, l0 int32) int32
	X__syscall_fadvise64(m *Module, l0 int32, l1 int64, l2 int64, l3 int32) int32
	X__syscall_fallocate(m *Module, l0 int32, l1 int32, l2 int64, l3 int64) int32
	X_emscripten_runtime_keepalive_clear(m *Module)
	X__call_sighandler(m *Module, l0 int32, l1 int32)
	X__syscall_getdents64(m *Module, l0 int32, l1 int32, l2 int32) int32
	X__syscall_readlinkat(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	X__syscall_renameat(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	X__syscall_rmdir(m *Module, l0 int32) int32
	X__syscall__newselect(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32) int32
	X_setitimer_js(m *Module, l0 int32, l1 float64) int32
	X__syscall_symlinkat(m *Module, l0 int32, l1 int32, l2 int32) int32
	Emscripten_get_heap_max(m *Module) int32
	X__syscall_truncate64(m *Module, l0 int32, l1 int64) int32
	X__syscall_unlinkat(m *Module, l0 int32, l1 int32, l2 int32) int32
	X__syscall_utimensat(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	Emscripten_resize_heap(m *Module, l0 int32) int32
	X_emscripten_throw_longjmp(m *Module)
	X__syscall_accept4(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32
	X__syscall_bind(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32
	X__syscall_connect(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32
	X__syscall_listen(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32
	X__syscall_recvfrom(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32
	X__syscall_sendto(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32
	X__syscall_socket(m *Module, l0 int32, l1 int32, l2 int32, l3 int32, l4 int32, l5 int32) int32
	X__syscall_fchmod(m *Module, l0 int32, l1 int32) int32
	X__syscall_fchmodat2(m *Module, l0 int32, l1 int32, l2 int32, l3 int32) int32
	X__syscall_fchown32(m *Module, l0 int32, l1 int32, l2 int32) int32
	X__syscall_statfs64(m *Module, l0 int32, l1 int32, l2 int32) int32
}
type Module struct {
	Memory                 []byte
	MaxMem                 uint64
	M                      unsafe.Pointer
	T0                     []any
	G0                     int32
	G1                     int32
	G2                     int32
	G3                     int32
	G4                     int32
	G5                     int32
	G6                     int32
	G7                     int32
	G8                     int32
	G9                     int32
	G10                    int32
	G11                    int32
	G12                    int32
	G13                    int32
	G14                    int32
	G15                    int32
	G16                    int32
	G17                    int32
	G18                    int32
	G19                    int32
	G20                    int32
	G21                    int32
	G22                    int32
	G23                    int32
	G24                    int32
	G25                    int32
	G26                    int32
	G27                    int32
	G28                    int32
	G29                    int32
	G30                    int32
	G31                    int32
	G32                    int32
	G33                    int32
	G34                    int32
	G35                    int32
	G36                    int32
	G37                    int32
	G38                    int32
	G39                    int32
	G40                    int32
	G41                    int32
	G42                    int32
	G43                    int32
	G44                    int32
	G45                    int32
	G46                    int32
	G47                    int32
	G48                    int32
	G49                    int32
	G50                    int32
	G51                    int32
	G52                    int32
	G53                    int32
	G54                    int32
	G55                    int32
	G56                    int32
	G57                    int32
	G58                    int32
	G59                    int32
	G60                    int32
	G61                    int32
	G62                    int32
	G63                    int32
	G64                    int32
	G65                    int32
	G66                    int32
	G67                    int32
	G68                    int32
	G69                    int32
	G70                    int32
	G71                    int32
	G72                    int32
	G73                    int32
	G74                    int32
	G75                    int32
	G76                    int32
	G77                    int32
	G78                    int32
	G79                    int32
	G80                    int32
	G81                    int32
	G82                    int32
	G83                    int32
	G84                    int32
	G85                    int32
	G86                    int32
	G87                    int32
	G88                    int32
	G89                    int32
	G90                    int32
	G91                    int32
	G92                    int32
	G93                    int32
	G94                    int32
	G95                    int32
	G96                    int32
	G97                    int32
	G98                    int32
	G99                    int32
	G100                   int32
	G101                   int32
	G102                   int32
	G103                   int32
	G104                   int32
	G105                   int32
	G106                   int32
	G107                   int32
	G108                   int32
	G109                   int32
	G110                   int32
	G111                   int32
	G112                   int32
	G113                   int32
	G114                   int32
	G115                   int32
	G116                   int32
	G117                   int32
	G118                   int32
	G119                   int32
	G120                   int32
	G121                   int32
	G122                   int32
	G123                   int32
	G124                   int32
	G125                   int32
	G126                   int32
	G127                   int32
	G128                   int32
	G129                   int32
	G130                   int32
	G131                   int32
	G132                   int32
	G133                   int32
	G134                   int32
	G135                   int32
	G136                   int32
	G137                   int32
	G138                   int32
	G139                   int32
	G140                   int32
	G141                   int32
	G142                   int32
	G143                   int32
	G144                   int32
	G145                   int32
	G146                   int32
	G147                   int32
	G148                   int32
	G149                   int32
	G150                   int32
	G151                   int32
	G152                   int32
	G153                   int32
	G154                   int32
	G155                   int32
	G156                   int32
	G157                   int32
	G158                   int32
	G159                   int32
	G160                   int32
	G161                   int32
	G162                   int32
	G163                   int32
	G164                   int32
	G165                   int32
	G166                   int32
	G167                   int32
	G168                   int32
	G169                   int32
	G170                   int32
	G171                   int32
	G172                   int32
	G173                   int32
	G174                   int32
	G175                   int32
	G176                   int32
	G177                   int32
	G178                   int32
	G179                   int32
	G180                   int32
	G181                   int32
	G182                   int32
	G183                   int32
	G184                   int32
	G185                   int32
	G186                   int32
	G187                   int32
	G188                   int32
	G189                   int32
	G190                   int32
	G191                   int32
	G192                   int32
	G193                   int32
	G194                   int32
	G195                   int32
	G196                   int32
	G197                   int32
	G198                   int32
	G199                   int32
	G200                   int32
	G201                   int32
	G202                   int32
	G203                   int32
	G204                   int32
	G205                   int32
	G206                   int32
	G207                   int32
	G208                   int32
	G209                   int32
	G210                   int32
	G211                   int32
	G212                   int32
	G213                   int32
	G214                   int32
	G215                   int32
	G216                   int32
	G217                   int32
	G218                   int32
	G219                   int32
	G220                   int32
	G221                   int32
	G222                   int32
	G223                   int32
	G224                   int32
	G225                   int32
	G226                   int32
	G227                   int32
	G228                   int32
	G229                   int32
	G230                   int32
	G231                   int32
	G232                   int32
	G233                   int32
	G234                   int32
	G235                   int32
	G236                   int32
	G237                   int32
	G238                   int32
	G239                   int32
	G240                   int32
	G241                   int32
	G242                   int32
	G243                   int32
	G244                   int32
	G245                   int32
	G246                   int32
	G247                   int32
	G248                   int32
	G249                   int32
	G250                   int32
	G251                   int32
	G252                   int32
	G253                   int32
	G254                   int32
	G255                   int32
	G256                   int32
	G257                   int32
	G258                   int32
	G259                   int32
	G260                   int32
	G261                   int32
	G262                   int32
	G263                   int32
	G264                   int32
	G265                   int32
	G266                   int32
	G267                   int32
	G268                   int32
	G269                   int32
	G270                   int32
	G271                   int32
	G272                   int32
	G273                   int32
	G274                   int32
	G275                   int32
	G276                   int32
	G277                   int32
	G278                   int32
	G279                   int32
	G280                   int32
	G281                   int32
	G282                   int32
	G283                   int32
	G284                   int32
	G285                   int32
	G286                   int32
	G287                   int32
	G288                   int32
	G289                   int32
	G290                   int32
	G291                   int32
	G292                   int32
	G293                   int32
	G294                   int32
	G295                   int32
	G296                   int32
	G297                   int32
	G298                   int32
	G299                   int32
	G300                   int32
	G301                   int32
	G302                   int32
	G303                   int32
	G304                   int32
	G305                   int32
	G306                   int32
	G307                   int32
	G308                   int32
	G309                   int32
	G310                   int32
	G311                   int32
	G312                   int32
	G313                   int32
	G314                   int32
	G315                   int32
	G316                   int32
	G317                   int32
	G318                   int32
	G319                   int32
	G320                   int32
	G321                   int32
	G322                   int32
	G323                   int32
	G324                   int32
	G325                   int32
	G326                   int32
	G327                   int32
	G328                   int32
	G329                   int32
	G330                   int32
	G331                   int32
	G332                   int32
	G333                   int32
	G334                   int32
	G335                   int32
	G336                   int32
	G337                   int32
	G338                   int32
	G339                   int32
	G340                   int32
	G341                   int32
	G342                   int32
	G343                   int32
	G344                   int32
	G345                   int32
	G346                   int32
	G347                   int32
	G348                   int32
	G349                   int32
	G350                   int32
	G351                   int32
	G352                   int32
	G353                   int32
	G354                   int32
	G355                   int32
	G356                   int32
	G357                   int32
	G358                   int32
	G359                   int32
	G360                   int32
	G361                   int32
	G362                   int32
	G363                   int32
	G364                   int32
	G365                   int32
	G366                   int32
	G367                   int32
	G368                   int32
	G369                   int32
	G370                   int32
	G371                   int32
	G372                   int32
	G373                   int32
	G374                   int32
	G375                   int32
	G376                   int32
	G377                   int32
	G378                   int32
	G379                   int32
	G380                   int32
	G381                   int32
	G382                   int32
	G383                   int32
	G384                   int32
	G385                   int32
	G386                   int32
	G387                   int32
	G388                   int32
	G389                   int32
	G390                   int32
	G391                   int32
	G392                   int32
	G393                   int32
	G394                   int32
	G395                   int32
	G396                   int32
	G397                   int32
	G398                   int32
	G399                   int32
	G400                   int32
	G401                   int32
	G402                   int32
	G403                   int32
	G404                   int32
	G405                   int32
	G406                   int32
	G407                   int32
	G408                   int32
	G409                   int32
	G410                   int32
	G411                   int32
	G412                   int32
	G413                   int32
	G414                   int32
	G415                   int32
	G416                   int32
	G417                   int32
	G418                   int32
	G419                   int32
	G420                   int32
	G421                   int32
	G422                   int32
	G423                   int32
	G424                   int32
	G425                   int32
	G426                   int32
	G427                   int32
	G428                   int32
	G429                   int32
	G430                   int32
	G431                   int32
	G432                   int32
	G433                   int32
	G434                   int32
	G435                   int32
	G436                   int32
	G437                   int32
	G438                   int32
	G439                   int32
	G440                   int32
	G441                   int32
	G442                   int32
	G443                   int32
	G444                   int32
	G445                   int32
	G446                   int32
	G447                   int32
	G448                   int32
	G449                   int32
	G450                   int32
	G451                   int32
	G452                   int32
	G453                   int32
	G454                   int32
	G455                   int32
	G456                   int32
	G457                   int32
	G458                   int32
	G459                   int32
	G460                   int32
	G461                   int32
	G462                   int32
	G463                   int32
	G464                   int32
	G465                   int32
	G466                   int32
	G467                   int32
	G468                   int32
	G469                   int32
	G470                   int32
	G471                   int32
	G472                   int32
	G473                   int32
	G474                   int32
	G475                   int32
	G476                   int32
	G477                   int32
	G478                   int32
	G479                   int32
	G480                   int32
	G481                   int32
	G482                   int32
	G483                   int32
	G484                   int32
	G485                   int32
	G486                   int32
	G487                   int32
	G488                   int32
	G489                   int32
	G490                   int32
	G491                   int32
	G492                   int32
	G493                   int32
	G494                   int32
	G495                   int32
	G496                   int32
	G497                   int32
	G498                   int32
	G499                   int32
	G500                   int32
	G501                   int32
	G502                   int32
	G503                   int32
	G504                   int32
	G505                   int32
	G506                   int32
	G507                   int32
	G508                   int32
	G509                   int32
	G510                   int32
	G511                   int32
	G512                   int32
	G513                   int32
	G514                   int32
	G515                   int32
	G516                   int32
	G517                   int32
	G518                   int32
	G519                   int32
	G520                   int32
	G521                   int32
	G522                   int32
	G523                   int32
	G524                   int32
	G525                   int32
	G526                   int32
	G527                   int32
	G528                   int32
	G529                   int32
	G530                   int32
	G531                   int32
	G532                   int32
	G533                   int32
	G534                   int32
	G535                   int32
	G536                   int32
	G537                   int32
	G538                   int32
	G539                   int32
	G540                   int32
	G541                   int32
	G542                   int32
	G543                   int32
	G544                   int32
	G545                   int32
	G546                   int32
	G547                   int32
	G548                   int32
	G549                   int32
	G550                   int32
	G551                   int32
	G552                   int32
	G553                   int32
	G554                   int32
	G555                   int32
	G556                   int32
	G557                   int32
	G558                   int32
	G559                   int32
	G560                   int32
	G561                   int32
	G562                   int32
	G563                   int32
	G564                   int32
	G565                   int32
	G566                   int32
	G567                   int32
	G568                   int32
	G569                   int32
	G570                   int32
	G571                   int32
	G572                   int32
	G573                   int32
	G574                   int32
	G575                   int32
	G576                   int32
	G577                   int32
	G578                   int32
	G579                   int32
	G580                   int32
	G581                   int32
	G582                   int32
	G583                   int32
	G584                   int32
	G585                   int32
	G586                   int32
	G587                   int32
	G588                   int32
	G589                   int32
	G590                   int32
	G591                   int32
	G592                   int32
	G593                   int32
	G594                   int32
	G595                   int32
	G596                   int32
	G597                   int32
	G598                   int32
	G599                   int32
	G600                   int32
	G601                   int32
	G602                   int32
	G603                   int32
	G604                   int32
	G605                   int32
	G606                   int32
	G607                   int32
	G608                   int32
	G609                   int32
	G610                   int32
	G611                   int32
	G612                   int32
	G613                   int32
	G614                   int32
	G615                   int32
	G616                   int32
	G617                   int32
	G618                   int32
	G619                   int32
	G620                   int32
	G621                   int32
	G622                   int32
	G623                   int32
	G624                   int32
	G625                   int32
	G626                   int32
	G627                   int32
	G628                   int32
	G629                   int32
	G630                   int32
	G631                   int32
	G632                   int32
	G633                   int32
	G634                   int32
	G635                   int32
	G636                   int32
	G637                   int32
	G638                   int32
	G639                   int32
	G640                   int32
	G641                   int32
	G642                   int32
	G643                   int32
	G644                   int32
	G645                   int32
	G646                   int32
	G647                   int32
	G648                   int32
	G649                   int32
	G650                   int32
	G651                   int32
	G652                   int32
	G653                   int32
	G654                   int32
	G655                   int32
	G656                   int32
	G657                   int32
	G658                   int32
	G659                   int32
	G660                   int32
	G661                   int32
	G662                   int32
	G663                   int32
	G664                   int32
	G665                   int32
	G666                   int32
	G667                   int32
	G668                   int32
	G669                   int32
	G670                   int32
	G671                   int32
	G672                   int32
	G673                   int32
	G674                   int32
	G675                   int32
	G676                   int32
	G677                   int32
	G678                   int32
	G679                   int32
	G680                   int32
	G681                   int32
	G682                   int32
	G683                   int32
	G684                   int32
	G685                   int32
	G686                   int32
	G687                   int32
	G688                   int32
	G689                   int32
	G690                   int32
	G691                   int32
	G692                   int32
	G693                   int32
	G694                   int32
	G695                   int32
	G696                   int32
	G697                   int32
	G698                   int32
	G699                   int32
	G700                   int32
	G701                   int32
	G702                   int32
	G703                   int32
	G704                   int32
	G705                   int32
	G706                   int32
	G707                   int32
	G708                   int32
	G709                   int32
	G710                   int32
	G711                   int32
	G712                   int32
	G713                   int32
	G714                   int32
	G715                   int32
	G716                   int32
	G717                   int32
	G718                   int32
	G719                   int32
	G720                   int32
	G721                   int32
	G722                   int32
	G723                   int32
	G724                   int32
	G725                   int32
	G726                   int32
	G727                   int32
	G728                   int32
	G729                   int32
	G730                   int32
	G731                   int32
	G732                   int32
	G733                   int32
	G734                   int32
	G735                   int32
	G736                   int32
	G737                   int32
	G738                   int32
	G739                   int32
	G740                   int32
	G741                   int32
	G742                   int32
	G743                   int32
	G744                   int32
	G745                   int32
	G746                   int32
	G747                   int32
	G748                   int32
	G749                   int32
	G750                   int32
	G751                   int32
	G752                   int32
	G753                   int32
	G754                   int32
	G755                   int32
	G756                   int32
	G757                   int32
	G758                   int32
	G759                   int32
	G760                   int32
	G761                   int32
	G762                   int32
	G763                   int32
	G764                   int32
	G765                   int32
	G766                   int32
	G767                   int32
	G768                   int32
	G769                   int32
	G770                   int32
	G771                   int32
	G772                   int32
	G773                   int32
	G774                   int32
	G775                   int32
	G776                   int32
	G777                   int32
	G778                   int32
	G779                   int32
	G780                   int32
	G781                   int32
	G782                   int32
	G783                   int32
	G784                   int32
	G785                   int32
	G786                   int32
	G787                   int32
	G788                   int32
	G789                   int32
	G790                   int32
	G791                   int32
	G792                   int32
	G793                   int32
	G794                   int32
	G795                   int32
	G796                   int32
	G797                   int32
	G798                   int32
	G799                   int32
	G800                   int32
	G801                   int32
	G802                   int32
	G803                   int32
	G804                   int32
	G805                   int32
	G806                   int32
	G807                   int32
	G808                   int32
	G809                   int32
	G810                   int32
	G811                   int32
	G812                   int32
	G813                   int32
	G814                   int32
	G815                   int32
	G816                   int32
	G817                   int32
	G818                   int32
	G819                   int32
	G820                   int32
	G821                   int32
	G822                   int32
	G823                   int32
	G824                   int32
	G825                   int32
	G826                   int32
	G827                   int32
	G828                   int32
	G829                   int32
	G830                   int32
	G831                   int32
	G832                   int32
	G833                   int32
	G834                   int32
	G835                   int32
	G836                   int32
	G837                   int32
	G838                   int32
	G839                   int32
	G840                   int32
	G841                   int32
	G842                   int32
	G843                   int32
	G844                   int32
	G845                   int32
	G846                   int32
	G847                   int32
	G848                   int32
	G849                   int32
	G850                   int32
	G851                   int32
	G852                   int32
	G853                   int32
	G854                   int32
	G855                   int32
	G856                   int32
	G857                   int32
	G858                   int32
	G859                   int32
	G860                   int32
	G861                   int32
	G862                   int32
	G863                   int32
	G864                   int32
	G865                   int32
	G866                   int32
	G867                   int32
	G868                   int32
	G869                   int32
	G870                   int32
	G871                   int32
	G872                   int32
	G873                   int32
	G874                   int32
	G875                   int32
	G876                   int32
	G877                   int32
	G878                   int32
	G879                   int32
	G880                   int32
	G881                   int32
	G882                   int32
	G883                   int32
	G884                   int32
	G885                   int32
	G886                   int32
	G887                   int32
	G888                   int32
	G889                   int32
	G890                   int32
	G891                   int32
	G892                   int32
	G893                   int32
	G894                   int32
	G895                   int32
	G896                   int32
	G897                   int32
	G898                   int32
	G899                   int32
	G900                   int32
	G901                   int32
	G902                   int32
	G903                   int32
	G904                   int32
	G905                   int32
	G906                   int32
	G907                   int32
	G908                   int32
	G909                   int32
	G910                   int32
	G911                   int32
	G912                   int32
	G913                   int32
	G914                   int32
	G915                   int32
	G916                   int32
	G917                   int32
	G918                   int32
	G919                   int32
	G920                   int32
	G921                   int32
	G922                   int32
	G923                   int32
	G924                   int32
	G925                   int32
	G926                   int32
	G927                   int32
	G928                   int32
	G929                   int32
	G930                   int32
	G931                   int32
	G932                   int32
	G933                   int32
	G934                   int32
	G935                   int32
	G936                   int32
	G937                   int32
	G938                   int32
	G939                   int32
	G940                   int32
	G941                   int32
	G942                   int32
	G943                   int32
	G944                   int32
	G945                   int32
	G946                   int32
	G947                   int32
	G948                   int32
	G949                   int32
	G950                   int32
	G951                   int32
	G952                   int32
	G953                   int32
	G954                   int32
	G955                   int32
	G956                   int32
	G957                   int32
	G958                   int32
	G959                   int32
	G960                   int32
	G961                   int32
	G962                   int32
	G963                   int32
	G964                   int32
	G965                   int32
	G966                   int32
	G967                   int32
	G968                   int32
	G969                   int32
	G970                   int32
	G971                   int32
	G972                   int32
	G973                   int32
	G974                   int32
	G975                   int32
	G976                   int32
	G977                   int32
	G978                   int32
	G979                   int32
	G980                   int32
	G981                   int32
	G982                   int32
	G983                   int32
	G984                   int32
	G985                   int32
	G986                   int32
	G987                   int32
	G988                   int32
	G989                   int32
	G990                   int32
	G991                   int32
	G992                   int32
	G993                   int32
	G994                   int32
	G995                   int32
	G996                   int32
	G997                   int32
	G998                   int32
	G999                   int32
	G1000                  int32
	G1001                  int32
	G1002                  int32
	G1003                  int32
	G1004                  int32
	G1005                  int32
	G1006                  int32
	G1007                  int32
	G1008                  int32
	G1009                  int32
	G1010                  int32
	G1011                  int32
	G1012                  int32
	G1013                  int32
	G1014                  int32
	G1015                  int32
	G1016                  int32
	G1017                  int32
	G1018                  int32
	G1019                  int32
	G1020                  int32
	G1021                  int32
	G1022                  int32
	G1023                  int32
	G1024                  int32
	G1025                  int32
	G1026                  int32
	G1027                  int32
	G1028                  int32
	G1029                  int32
	G1030                  int32
	G1031                  int32
	G1032                  int32
	G1033                  int32
	G1034                  int32
	G1035                  int32
	G1036                  int32
	G1037                  int32
	G1038                  int32
	G1039                  int32
	G1040                  int32
	G1041                  int32
	G1042                  int32
	G1043                  int32
	G1044                  int32
	G1045                  int32
	G1046                  int32
	G1047                  int32
	G1048                  int32
	G1049                  int32
	G1050                  int32
	G1051                  int32
	G1052                  int32
	G1053                  int32
	G1054                  int32
	G1055                  int32
	G1056                  int32
	G1057                  int32
	G1058                  int32
	G1059                  int32
	G1060                  int32
	G1061                  int32
	G1062                  int32
	G1063                  int32
	G1064                  int32
	G1065                  int32
	G1066                  int32
	G1067                  int32
	G1068                  int32
	G1069                  int32
	G1070                  int32
	G1071                  int32
	G1072                  int32
	G1073                  int32
	G1074                  int32
	G1075                  int32
	G1076                  int32
	G1077                  int32
	G1078                  int32
	G1079                  int32
	G1080                  int32
	G1081                  int32
	G1082                  int32
	G1083                  int32
	G1084                  int32
	G1085                  int32
	G1086                  int32
	G1087                  int32
	G1088                  int32
	G1089                  int32
	G1090                  int32
	G1091                  int32
	G1092                  int32
	G1093                  int32
	G1094                  int32
	G1095                  int32
	G1096                  int32
	G1097                  int32
	G1098                  int32
	G1099                  int32
	G1100                  int32
	G1101                  int32
	G1102                  int32
	G1103                  int32
	G1104                  int32
	G1105                  int32
	G1106                  int32
	G1107                  int32
	G1108                  int32
	G1109                  int32
	G1110                  int32
	G1111                  int32
	G1112                  int32
	G1113                  int32
	G1114                  int32
	G1115                  int32
	G1116                  int32
	G1117                  int32
	G1118                  int32
	G1119                  int32
	G1120                  int32
	G1121                  int32
	G1122                  int32
	G1123                  int32
	G1124                  int32
	G1125                  int32
	G1126                  int32
	G1127                  int32
	G1128                  int32
	G1129                  int32
	G1130                  int32
	G1131                  int32
	G1132                  int32
	G1133                  int32
	G1134                  int32
	G1135                  int32
	G1136                  int32
	G1137                  int32
	G1138                  int32
	G1139                  int32
	G1140                  int32
	G1141                  int32
	G1142                  int32
	G1143                  int32
	G1144                  int32
	G1145                  int32
	G1146                  int32
	G1147                  int32
	G1148                  int32
	G1149                  int32
	G1150                  int32
	G1151                  int32
	G1152                  int32
	G1153                  int32
	G1154                  int32
	G1155                  int32
	G1156                  int32
	G1157                  int32
	G1158                  int32
	G1159                  int32
	G1160                  int32
	G1161                  int32
	G1162                  int32
	G1163                  int32
	G1164                  int32
	G1165                  int32
	G1166                  int32
	G1167                  int32
	G1168                  int32
	G1169                  int32
	G1170                  int32
	G1171                  int32
	G1172                  int32
	G1173                  int32
	G1174                  int32
	G1175                  int32
	G1176                  int32
	G1177                  int32
	G1178                  int32
	G1179                  int32
	G1180                  int32
	G1181                  int32
	G1182                  int32
	G1183                  int32
	G1184                  int32
	G1185                  int32
	G1186                  int32
	G1187                  int32
	G1188                  int32
	G1189                  int32
	G1190                  int32
	G1191                  int32
	G1192                  int32
	G1193                  int32
	G1194                  int32
	G1195                  int32
	G1196                  int32
	G1197                  int32
	G1198                  int32
	G1199                  int32
	G1200                  int32
	G1201                  int32
	G1202                  int32
	G1203                  int32
	G1204                  int32
	G1205                  int32
	G1206                  int32
	G1207                  int32
	G1208                  int32
	G1209                  int32
	G1210                  int32
	G1211                  int32
	G1212                  int32
	G1213                  int32
	G1214                  int32
	G1215                  int32
	G1216                  int32
	G1217                  int32
	G1218                  int32
	G1219                  int32
	G1220                  int32
	G1221                  int32
	G1222                  int32
	G1223                  int32
	G1224                  int32
	G1225                  int32
	G1226                  int32
	G1227                  int32
	G1228                  int32
	G1229                  int32
	G1230                  int32
	G1231                  int32
	G1232                  int32
	G1233                  int32
	G1234                  int32
	G1235                  int32
	G1236                  int32
	G1237                  int32
	G1238                  int32
	G1239                  int32
	G1240                  int32
	G1241                  int32
	G1242                  int32
	G1243                  int32
	G1244                  int32
	G1245                  int32
	G1246                  int32
	G1247                  int32
	G1248                  int32
	G1249                  int32
	G1250                  int32
	G1251                  int32
	G1252                  int32
	G1253                  int32
	G1254                  int32
	G1255                  int32
	G1256                  int32
	G1257                  int32
	G1258                  int32
	G1259                  int32
	G1260                  int32
	G1261                  int32
	G1262                  int32
	G1263                  int32
	G1264                  int32
	G1265                  int32
	G1266                  int32
	G1267                  int32
	G1268                  int32
	G1269                  int32
	G1270                  int32
	G1271                  int32
	G1272                  int32
	G1273                  int32
	G1274                  int32
	G1275                  int32
	G1276                  int32
	G1277                  int32
	G1278                  int32
	G1279                  int32
	G1280                  int32
	G1281                  int32
	G1282                  int32
	G1283                  int32
	G1284                  int32
	G1285                  int32
	G1286                  int32
	G1287                  int32
	G1288                  int32
	G1289                  int32
	G1290                  int32
	G1291                  int32
	G1292                  int32
	G1293                  int32
	G1294                  int32
	G1295                  int32
	G1296                  int32
	G1297                  int32
	G1298                  int32
	G1299                  int32
	G1300                  int32
	G1301                  int32
	G1302                  int32
	G1303                  int32
	G1304                  int32
	G1305                  int32
	G1306                  int32
	G1307                  int32
	G1308                  int32
	G1309                  int32
	G1310                  int32
	G1311                  int32
	G1312                  int32
	G1313                  int32
	G1314                  int32
	G1315                  int32
	G1316                  int32
	G1317                  int32
	G1318                  int32
	G1319                  int32
	G1320                  int32
	G1321                  int32
	G1322                  int32
	G1323                  int32
	G1324                  int32
	G1325                  int32
	G1326                  int32
	G1327                  int32
	G1328                  int32
	G1329                  int32
	G1330                  int32
	G1331                  int32
	G1332                  int32
	G1333                  int32
	G1334                  int32
	G1335                  int32
	G1336                  int32
	G1337                  int32
	G1338                  int32
	G1339                  int32
	G1340                  int32
	G1341                  int32
	G1342                  int32
	G1343                  int32
	G1344                  int32
	G1345                  int32
	G1346                  int32
	G1347                  int32
	G1348                  int32
	G1349                  int32
	G1350                  int32
	G1351                  int32
	G1352                  int32
	G1353                  int32
	G1354                  int32
	G1355                  int32
	G1356                  int32
	G1357                  int32
	G1358                  int32
	G1359                  int32
	G1360                  int32
	G1361                  int32
	G1362                  int32
	G1363                  int32
	G1364                  int32
	G1365                  int32
	G1366                  int32
	G1367                  int32
	G1368                  int32
	G1369                  int32
	G1370                  int32
	G1371                  int32
	G1372                  int32
	G1373                  int32
	G1374                  int32
	G1375                  int32
	G1376                  int32
	G1377                  int32
	G1378                  int32
	G1379                  int32
	G1380                  int32
	G1381                  int32
	G1382                  int32
	G1383                  int32
	G1384                  int32
	G1385                  int32
	G1386                  int32
	G1387                  int32
	G1388                  int32
	G1389                  int32
	G1390                  int32
	G1391                  int32
	G1392                  int32
	G1393                  int32
	G1394                  int32
	G1395                  int32
	G1396                  int32
	G1397                  int32
	G1398                  int32
	G1399                  int32
	G1400                  int32
	G1401                  int32
	G1402                  int32
	G1403                  int32
	G1404                  int32
	G1405                  int32
	G1406                  int32
	G1407                  int32
	G1408                  int32
	G1409                  int32
	G1410                  int32
	G1411                  int32
	G1412                  int32
	G1413                  int32
	G1414                  int32
	G1415                  int32
	G1416                  int32
	G1417                  int32
	G1418                  int32
	G1419                  int32
	G1420                  int32
	G1421                  int32
	G1422                  int32
	G1423                  int32
	G1424                  int32
	G1425                  int32
	G1426                  int32
	G1427                  int32
	G1428                  int32
	G1429                  int32
	G1430                  int32
	G1431                  int32
	G1432                  int32
	G1433                  int32
	G1434                  int32
	G1435                  int32
	G1436                  int32
	G1437                  int32
	G1438                  int32
	G1439                  int32
	G1440                  int32
	G1441                  int32
	G1442                  int32
	G1443                  int32
	G1444                  int32
	G1445                  int32
	G1446                  int32
	G1447                  int32
	G1448                  int32
	G1449                  int32
	G1450                  int32
	G1451                  int32
	G1452                  int32
	G1453                  int32
	G1454                  int32
	G1455                  int32
	G1456                  int32
	G1457                  int32
	G1458                  int32
	G1459                  int32
	G1460                  int32
	G1461                  int32
	G1462                  int32
	G1463                  int32
	G1464                  int32
	G1465                  int32
	G1466                  int32
	G1467                  int32
	G1468                  int32
	G1469                  int32
	G1470                  int32
	G1471                  int32
	G1472                  int32
	G1473                  int32
	G1474                  int32
	G1475                  int32
	G1476                  int32
	G1477                  int32
	G1478                  int32
	G1479                  int32
	G1480                  int32
	G1481                  int32
	G1482                  int32
	G1483                  int32
	G1484                  int32
	G1485                  int32
	G1486                  int32
	G1487                  int32
	G1488                  int32
	G1489                  int32
	G1490                  int32
	G1491                  int32
	G1492                  int32
	G1493                  int32
	G1494                  int32
	G1495                  int32
	G1496                  int32
	G1497                  int32
	G1498                  int32
	G1499                  int32
	G1500                  int32
	G1501                  int32
	G1502                  int32
	G1503                  int32
	G1504                  int32
	G1505                  int32
	G1506                  int32
	G1507                  int32
	G1508                  int32
	G1509                  int32
	G1510                  int32
	G1511                  int32
	G1512                  int32
	G1513                  int32
	G1514                  int32
	G1515                  int32
	G1516                  int32
	G1517                  int32
	G1518                  int32
	G1519                  int32
	G1520                  int32
	G1521                  int32
	G1522                  int32
	G1523                  int32
	G1524                  int32
	G1525                  int32
	G1526                  int32
	G1527                  int32
	G1528                  int32
	G1529                  int32
	G1530                  int32
	G1531                  int32
	G1532                  int32
	G1533                  int32
	G1534                  int32
	G1535                  int32
	G1536                  int32
	G1537                  int32
	G1538                  int32
	G1539                  int32
	G1540                  int32
	G1541                  int32
	G1542                  int32
	G1543                  int32
	G1544                  int32
	G1545                  int32
	G1546                  int32
	G1547                  int32
	G1548                  int32
	G1549                  int32
	G1550                  int32
	G1551                  int32
	G1552                  int32
	G1553                  int32
	G1554                  int32
	G1555                  int32
	G1556                  int32
	G1557                  int32
	G1558                  int32
	G1559                  int32
	G1560                  int32
	G1561                  int32
	G1562                  int32
	G1563                  int32
	G1564                  int32
	G1565                  int32
	G1566                  int32
	G1567                  int32
	G1568                  int32
	G1569                  int32
	G1570                  int32
	G1571                  int32
	G1572                  int32
	G1573                  int32
	G1574                  int32
	G1575                  int32
	G1576                  int32
	G1577                  int32
	G1578                  int32
	G1579                  int32
	G1580                  int32
	G1581                  int32
	G1582                  int32
	G1583                  int32
	G1584                  int32
	G1585                  int32
	G1586                  int32
	G1587                  int32
	G1588                  int32
	G1589                  int32
	G1590                  int32
	G1591                  int32
	G1592                  int32
	G1593                  int32
	G1594                  int32
	G1595                  int32
	G1596                  int32
	G1597                  int32
	G1598                  int32
	G1599                  int32
	G1600                  int32
	G1601                  int32
	G1602                  int32
	G1603                  int32
	G1604                  int32
	G1605                  int32
	G1606                  int32
	G1607                  int32
	G1608                  int32
	G1609                  int32
	G1610                  int32
	G1611                  int32
	G1612                  int32
	G1613                  int32
	G1614                  int32
	G1615                  int32
	G1616                  int32
	G1617                  int32
	G1618                  int32
	G1619                  int32
	G1620                  int32
	G1621                  int32
	G1622                  int32
	G1623                  int32
	G1624                  int32
	G1625                  int32
	G1626                  int32
	G1627                  int32
	G1628                  int32
	G1629                  int32
	G1630                  int32
	G1631                  int32
	G1632                  int32
	G1633                  int32
	G1634                  int32
	G1635                  int32
	G1636                  int32
	G1637                  int32
	G1638                  int32
	G1639                  int32
	G1640                  int32
	G1641                  int32
	G1642                  int32
	G1643                  int32
	G1644                  int32
	G1645                  int32
	G1646                  int32
	G1647                  int32
	G1648                  int32
	G1649                  int32
	G1650                  int32
	G1651                  int32
	G1652                  int32
	G1653                  int32
	G1654                  int32
	G1655                  int32
	G1656                  int32
	G1657                  int32
	G1658                  int32
	G1659                  int32
	G1660                  int32
	G1661                  int32
	G1662                  int32
	G1663                  int32
	G1664                  int32
	G1665                  int32
	G1666                  int32
	G1667                  int32
	G1668                  int32
	G1669                  int32
	G1670                  int32
	G1671                  int32
	G1672                  int32
	G1673                  int32
	G1674                  int32
	G1675                  int32
	G1676                  int32
	G1677                  int32
	G1678                  int32
	G1679                  int32
	G1680                  int32
	G1681                  int32
	G1682                  int32
	G1683                  int32
	G1684                  int32
	G1685                  int32
	G1686                  int32
	G1687                  int32
	G1688                  int32
	G1689                  int32
	G1690                  int32
	G1691                  int32
	G1692                  int32
	G1693                  int32
	G1694                  int32
	G1695                  int32
	G1696                  int32
	G1697                  int32
	G1698                  int32
	G1699                  int32
	G1700                  int32
	G1701                  int32
	G1702                  int32
	G1703                  int32
	G1704                  int32
	G1705                  int32
	G1706                  int32
	G1707                  int32
	G1708                  int32
	G1709                  int32
	G1710                  int32
	G1711                  int32
	G1712                  int32
	G1713                  int32
	G1714                  int32
	G1715                  int32
	G1716                  int32
	G1717                  int32
	G1718                  int32
	G1719                  int32
	G1720                  int32
	G1721                  int32
	G1722                  int32
	G1723                  int32
	G1724                  int32
	G1725                  int32
	G1726                  int32
	G1727                  int32
	G1728                  int32
	G1729                  int32
	G1730                  int32
	G1731                  int32
	G1732                  int32
	G1733                  int32
	G1734                  int32
	G1735                  int32
	G1736                  int32
	G1737                  int32
	G1738                  int32
	G1739                  int32
	G1740                  int32
	G1741                  int32
	G1742                  int32
	G1743                  int32
	G1744                  int32
	G1745                  int32
	G1746                  int32
	G1747                  int32
	G1748                  int32
	G1749                  int32
	G1750                  int32
	G1751                  int32
	G1752                  int32
	G1753                  int32
	G1754                  int32
	G1755                  int32
	G1756                  int32
	G1757                  int32
	G1758                  int32
	G1759                  int32
	G1760                  int32
	G1761                  int32
	G1762                  int32
	G1763                  int32
	G1764                  int32
	G1765                  int32
	G1766                  int32
	G1767                  int32
	G1768                  int32
	G1769                  int32
	G1770                  int32
	G1771                  int32
	G1772                  int32
	G1773                  int32
	G1774                  int32
	G1775                  int32
	G1776                  int32
	G1777                  int32
	G1778                  int32
	G1779                  int32
	G1780                  int32
	G1781                  int32
	G1782                  int32
	G1783                  int32
	G1784                  int32
	G1785                  int32
	G1786                  int32
	G1787                  int32
	G1788                  int32
	G1789                  int32
	G1790                  int32
	G1791                  int32
	G1792                  int32
	G1793                  int32
	G1794                  int32
	G1795                  int32
	G1796                  int32
	G1797                  int32
	G1798                  int32
	G1799                  int32
	G1800                  int32
	G1801                  int32
	G1802                  int32
	G1803                  int32
	G1804                  int32
	G1805                  int32
	G1806                  int32
	G1807                  int32
	G1808                  int32
	G1809                  int32
	G1810                  int32
	G1811                  int32
	G1812                  int32
	G1813                  int32
	G1814                  int32
	G1815                  int32
	G1816                  int32
	G1817                  int32
	G1818                  int32
	G1819                  int32
	G1820                  int32
	G1821                  int32
	G1822                  int32
	G1823                  int32
	G1824                  int32
	G1825                  int32
	G1826                  int32
	G1827                  int32
	G1828                  int32
	G1829                  int32
	G1830                  int32
	G1831                  int32
	G1832                  int32
	G1833                  int32
	G1834                  int32
	G1835                  int32
	G1836                  int32
	G1837                  int32
	G1838                  int32
	G1839                  int32
	G1840                  int32
	G1841                  int32
	G1842                  int32
	G1843                  int32
	G1844                  int32
	G1845                  int32
	G1846                  int32
	G1847                  int32
	G1848                  int32
	G1849                  int32
	G1850                  int32
	G1851                  int32
	G1852                  int32
	G1853                  int32
	G1854                  int32
	G1855                  int32
	G1856                  int32
	G1857                  int32
	G1858                  int32
	G1859                  int32
	G1860                  int32
	G1861                  int32
	G1862                  int32
	G1863                  int32
	G1864                  int32
	G1865                  int32
	G1866                  int32
	G1867                  int32
	G1868                  int32
	G1869                  int32
	G1870                  int32
	G1871                  int32
	G1872                  int32
	G1873                  int32
	G1874                  int32
	G1875                  int32
	G1876                  int32
	G1877                  int32
	G1878                  int32
	G1879                  int32
	G1880                  int32
	G1881                  int32
	G1882                  int32
	G1883                  int32
	G1884                  int32
	G1885                  int32
	G1886                  int32
	G1887                  int32
	G1888                  int32
	G1889                  int32
	G1890                  int32
	G1891                  int32
	G1892                  int32
	G1893                  int32
	G1894                  int32
	G1895                  int32
	G1896                  int32
	G1897                  int32
	G1898                  int32
	G1899                  int32
	G1900                  int32
	G1901                  int32
	G1902                  int32
	G1903                  int32
	G1904                  int32
	G1905                  int32
	G1906                  int32
	G1907                  int32
	G1908                  int32
	G1909                  int32
	G1910                  int32
	G1911                  int32
	G1912                  int32
	G1913                  int32
	G1914                  int32
	G1915                  int32
	G1916                  int32
	G1917                  int32
	G1918                  int32
	G1919                  int32
	G1920                  int32
	G1921                  int32
	G1922                  int32
	G1923                  int32
	G1924                  int32
	G1925                  int32
	G1926                  int32
	G1927                  int32
	G1928                  int32
	G1929                  int32
	G1930                  int32
	G1931                  int32
	G1932                  int32
	G1933                  int32
	G1934                  int32
	G1935                  int32
	G1936                  int32
	G1937                  int32
	G1938                  int32
	G1939                  int32
	G1940                  int32
	G1941                  int32
	G1942                  int32
	G1943                  int32
	G1944                  int32
	G1945                  int32
	G1946                  int32
	G1947                  int32
	G1948                  int32
	G1949                  int32
	G1950                  int32
	G1951                  int32
	G1952                  int32
	G1953                  int32
	G1954                  int32
	G1955                  int32
	G1956                  int32
	G1957                  int32
	G1958                  int32
	G1959                  int32
	G1960                  int32
	G1961                  int32
	G1962                  int32
	G1963                  int32
	G1964                  int32
	G1965                  int32
	G1966                  int32
	G1967                  int32
	G1968                  int32
	G1969                  int32
	G1970                  int32
	G1971                  int32
	G1972                  int32
	G1973                  int32
	G1974                  int32
	G1975                  int32
	G1976                  int32
	G1977                  int32
	G1978                  int32
	G1979                  int32
	G1980                  int32
	G1981                  int32
	G1982                  int32
	G1983                  int32
	G1984                  int32
	G1985                  int32
	G1986                  int32
	G1987                  int32
	G1988                  int32
	G1989                  int32
	G1990                  int32
	G1991                  int32
	G1992                  int32
	G1993                  int32
	G1994                  int32
	G1995                  int32
	G1996                  int32
	G1997                  int32
	G1998                  int32
	G1999                  int32
	G2000                  int32
	G2001                  int32
	G2002                  int32
	G2003                  int32
	G2004                  int32
	G2005                  int32
	G2006                  int32
	G2007                  int32
	G2008                  int32
	G2009                  int32
	G2010                  int32
	G2011                  int32
	G2012                  int32
	G2013                  int32
	G2014                  int32
	G2015                  int32
	G2016                  int32
	G2017                  int32
	G2018                  int32
	G2019                  int32
	G2020                  int32
	G2021                  int32
	G2022                  int32
	G2023                  int32
	G2024                  int32
	G2025                  int32
	G2026                  int32
	G2027                  int32
	G2028                  int32
	G2029                  int32
	G2030                  int32
	G2031                  int32
	G2032                  int32
	G2033                  int32
	G2034                  int32
	G2035                  int32
	G2036                  int32
	G2037                  int32
	G2038                  int32
	G2039                  int32
	G2040                  int32
	G2041                  int32
	G2042                  int32
	G2043                  int32
	G2044                  int32
	Wasi_snapshot_preview1 Wasi_snapshot_preview1Imports
	Env                    EnvImports
	MemMu                  *sync.Mutex
	MemSize                *atomic.Uint64
	DataEnd                uint32
	MemShared              bool
	Threads                *ThreadPool
	ThreadStart            func(*Module, int32, int32)
}

func I32(x int32) int32 { return x }

func I64(x int64) int64 { return x }

// ui32 / ui64 reinterpret a signed integer as its unsigned bit
// equivalent at runtime. Used for the operands of wasm unsigned
// comparisons (i32.lt_u etc.) — emitting `uint32(int32(-N))` directly
// fails Go's compile-time constant rule because the negative typed
// constant isn't representable in uint32; routing through these
// function-call boundaries forces runtime conversion.
func Ui32(x int32) uint32 { return uint32(x) }

func Ui64(x int64) uint64 { return uint64(x) }

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
func B2i32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func F32(x float32) float32 { runtime.KeepAlive(&x); return x }

func F64(x float64) float64 { runtime.KeepAlive(&x); return x }

//go:noinline
func Wasm_trap_div_zero() { panic("wasm: integer divide by zero") }

//go:noinline
func Wasm_trap_int_overflow() { panic("wasm: integer overflow") }

//go:noinline
func Wasm_trap_invalid_conv() { panic("wasm: invalid conversion to integer") }

//go:noinline
func Wasm_trap_unreachable() { panic("wasm: unreachable") }

//go:noinline
func Wasm_trap_memfill_oob() { panic("wasm: memory.fill out of bounds") }

//go:noinline
func Wasm_trap_memcopy_oob() { panic("wasm: memory.copy out of bounds") }

func I32_div_s(x, y int32) int32 {
	if y == -1 && x == math.MinInt32 {
		Wasm_trap_int_overflow()
	}
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x / y
}

func I64_div_s(x, y int64) int64 {
	if y == -1 && x == math.MinInt64 {
		Wasm_trap_int_overflow()
	}
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x / y
}

func I32_div_u(x, y uint32) uint32 {
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x / y
}

func I64_div_u(x, y uint64) uint64 {
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x / y
}

func I32_rem_s(x, y int32) int32 {
	if y == 0 {
		Wasm_trap_div_zero()
	}
	if y == -1 {

		return 0
	}
	return x % y
}

func I64_rem_s(x, y int64) int64 {
	if y == 0 {
		Wasm_trap_div_zero()
	}
	if y == -1 {
		return 0
	}
	return x % y
}

func I32_rem_u(x, y uint32) uint32 {
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x % y
}

func I64_rem_u(x, y uint64) uint64 {
	if y == 0 {
		Wasm_trap_div_zero()
	}
	return x % y
}

func I32_rotl(x, y int32) int32 { return int32(bits.RotateLeft32(uint32(x), int(y&31))) }

func I32_rotr(x, y int32) int32 { return int32(bits.RotateLeft32(uint32(x), -int(y&31))) }

func I64_rotl(x, y int64) int64 { return int64(bits.RotateLeft64(uint64(x), int(y&63))) }

func I64_rotr(x, y int64) int64 { return int64(bits.RotateLeft64(uint64(x), -int(y&63))) }

func F64_max(x, y float64) float64 {
	if x != x || y != y {
		return math.NaN()
	}
	if x > y {
		return x
	}
	if y > x {
		return y
	}
	if x == 0 {
		if math.Signbit(x) {
			return y
		}
		return x
	}
	return x
}

func F32_abs(x float32) float32 {
	return math.Float32frombits(math.Float32bits(x) &^ (1 << 31))
}

func F64_abs(x float64) float64 {
	return math.Float64frombits(math.Float64bits(x) &^ (1 << 63))
}

func F32_neg(x float32) float32 {
	return math.Float32frombits(math.Float32bits(x) ^ (1 << 31))
}

func F64_neg(x float64) float64 {
	return math.Float64frombits(math.Float64bits(x) ^ (1 << 63))
}

func F32_copysign(x, y float32) float32 {
	return float32(math.Copysign(float64(x), float64(y)))
}

func F64_copysign(x, y float64) float64 { return math.Copysign(x, y) }

func F32_nearest(x float32) float32 { return float32(math.RoundToEven(float64(x))) }

func F64_nearest(x float64) float64 { return math.RoundToEven(x) }

func I32_trunc_f32_s(x float32) int32 {
	if x != x {
		Wasm_trap_invalid_conv()
	}

	if !(x > -2147483904.0 && x < 2147483648.0) {
		Wasm_trap_int_overflow()
	}
	return int32(x)
}

func I32_trunc_f32_u(x float32) int32 {
	if x != x {
		Wasm_trap_invalid_conv()
	}
	if !(x > -1.0 && x < 4294967296.0) {
		Wasm_trap_int_overflow()
	}
	return int32(uint32(x))
}

func I32_trunc_f64_s(x float64) int32 {
	if x != x {
		Wasm_trap_invalid_conv()
	}

	if !(x > -2147483649.0 && x < 2147483648.0) {
		Wasm_trap_int_overflow()
	}
	return int32(x)
}

func I32_trunc_f64_u(x float64) int32 {
	if x != x {
		Wasm_trap_invalid_conv()
	}
	if !(x > -1.0 && x < 4294967296.0) {
		Wasm_trap_int_overflow()
	}
	return int32(uint32(x))
}

func I64_trunc_f32_s(x float32) int64 {
	if x != x {
		Wasm_trap_invalid_conv()
	}

	if !(float64(x) > -9223373136366403584.0 && float64(x) < 9223372036854775808.0) {
		Wasm_trap_int_overflow()
	}
	return int64(x)
}

func I64_trunc_f64_s(x float64) int64 {
	if x != x {
		Wasm_trap_invalid_conv()
	}
	if !(x >= -9223372036854775808.0 && x < 9223372036854775808.0) {
		Wasm_trap_int_overflow()
	}
	return int64(x)
}

func I64_trunc_f64_u(x float64) int64 {
	if x != x {
		Wasm_trap_invalid_conv()
	}
	if !(x > -1.0 && x < 18446744073709551616.0) {
		Wasm_trap_int_overflow()
	}
	return int64(uint64(x))
}

// memorySize returns the current size of m.memory in wasm pages (each
// page is 64 KiB).
func MemorySize(m *Module) int32 {
	return int32(m.MemSize.Load() >> 16)
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
func AccessMemory(m *Module, f func(mem []byte)) {
	m.MemMu.Lock()
	defer m.MemMu.Unlock()
	f(m.Memory)
}

func I32_div_u_s(x, y int32) int32 { return int32(I32_div_u(uint32(x), uint32(y))) }
func I32_rem_u_s(x, y int32) int32 { return int32(I32_rem_u(uint32(x), uint32(y))) }
func I64_div_u_s(x, y int64) int64 { return int64(I64_div_u(uint64(x), uint64(y))) }
func I64_rem_u_s(x, y int64) int64 { return int64(I64_rem_u(uint64(x), uint64(y))) }

// The explicit same-type conversions are NOT redundant: they are
// rounding points. Once these helpers inline, gc is free to fuse a
// multiply feeding an add into a single FMA — legal Go, but wasm
// requires every operation individually rounded, and a fused result
// diverges from every wasm runtime (bitwise, and observably in greedy
// sampling). A float conversion forces the intermediate rounding and
// forbids the fusion (spec: Conversions, "rounds to the precision of
// the target type"; the same rule math.FMA documents).
func F32_add(x, y float32) float32 { return float32(x + y) }
func F32_sub(x, y float32) float32 { return float32(x - y) }
func F32_mul(x, y float32) float32 { return float32(x * y) }
func F32_div(x, y float32) float32 { return float32(x / y) }
func F64_add(x, y float64) float64 { return float64(x + y) }
func F64_sub(x, y float64) float64 { return float64(x - y) }
func F64_mul(x, y float64) float64 { return float64(x * y) }
func F64_div(x, y float64) float64 { return float64(x / y) }

func I32_clz(x int32) int32    { return int32(bits.LeadingZeros32(uint32(x))) }
func I32_ctz(x int32) int32    { return int32(bits.TrailingZeros32(uint32(x))) }
func I32_popcnt(x int32) int32 { return int32(bits.OnesCount32(uint32(x))) }

func I64_clz(x int64) int64    { return int64(bits.LeadingZeros64(uint64(x))) }
func I64_ctz(x int64) int64    { return int64(bits.TrailingZeros64(uint64(x))) }
func I64_popcnt(x int64) int64 { return int64(bits.OnesCount64(uint64(x))) }

func F64_ceil(x float64) float64 { return math.Ceil(x) }

func F64_floor(x float64) float64 { return math.Floor(x) }

func F32_sqrt(x float32) float32 { return float32(math.Sqrt(float64(x))) }
func F64_sqrt(x float64) float64 { return math.Sqrt(x) }

func F32_eq(x, y float32) int32 {
	if x == y {
		return 1
	}
	return 0
}
func F32_ne(x, y float32) int32 {
	if x != y {
		return 1
	}
	return 0
}
func F32_lt(x, y float32) int32 {
	if x < y {
		return 1
	}
	return 0
}
func F32_gt(x, y float32) int32 {
	if x > y {
		return 1
	}
	return 0
}
func F32_le(x, y float32) int32 {
	if x <= y {
		return 1
	}
	return 0
}
func F32_ge(x, y float32) int32 {
	if x >= y {
		return 1
	}
	return 0
}

func F64_eq(x, y float64) int32 {
	if x == y {
		return 1
	}
	return 0
}
func F64_ne(x, y float64) int32 {
	if x != y {
		return 1
	}
	return 0
}
func F64_lt(x, y float64) int32 {
	if x < y {
		return 1
	}
	return 0
}
func F64_gt(x, y float64) int32 {
	if x > y {
		return 1
	}
	return 0
}
func F64_le(x, y float64) int32 {
	if x <= y {
		return 1
	}
	return 0
}
func F64_ge(x, y float64) int32 {
	if x >= y {
		return 1
	}
	return 0
}

func I32_wrap_i64(x int64) int32       { return int32(x) }
func I64_extend_i32_s(x int32) int64   { return int64(x) }
func I64_extend_i32_u(x int32) int64   { return int64(uint32(x)) }
func F32_demote_f64(x float64) float32 { return float32(x) }
func F64_promote_f32(x float32) float64 {

	if math.IsNaN(float64(x)) {

		return float64(x)
	}
	return float64(x)
}

func F32_convert_i32_s(x int32) float32 { return float32(x) }
func F32_convert_i32_u(x int32) float32 { return float32(uint32(x)) }
func F32_convert_i64_s(x int64) float32 { return float32(x) }

func F64_convert_i32_s(x int32) float64 { return float64(x) }
func F64_convert_i32_u(x int32) float64 { return float64(uint32(x)) }
func F64_convert_i64_s(x int64) float64 { return float64(x) }
func F64_convert_i64_u(x int64) float64 { return float64(uint64(x)) }

func I32_reinterpret_f32(x float32) int32 { return int32(math.Float32bits(x)) }
func I64_reinterpret_f64(x float64) int64 { return int64(math.Float64bits(x)) }
func F32_reinterpret_i32(x int32) float32 { return math.Float32frombits(uint32(x)) }
func F64_reinterpret_i64(x int64) float64 { return math.Float64frombits(uint64(x)) }

func I32_extend8_s(x int32) int32  { return int32(int8(x)) }
func I32_extend16_s(x int32) int32 { return int32(int16(x)) }
func I64_extend8_s(x int64) int64  { return int64(int8(x)) }
func I64_extend16_s(x int64) int64 { return int64(int16(x)) }
func I64_extend32_s(x int64) int64 { return int64(int32(x)) }

func MemoryFill(m *Module, dst int32, val int32, n int32) {
	if n == 0 {
		return
	}
	end := uint64(uint32(dst)) + uint64(uint32(n))
	if end > m.MemSize.Load() {
		Wasm_trap_memfill_oob()
	}
	b := m.Memory[uint32(dst):uint32(end)]
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

func MemoryCopy(m *Module, dst int32, src int32, n int32) {
	if n == 0 {
		return
	}
	srcEnd := uint64(uint32(src)) + uint64(uint32(n))
	dstEnd := uint64(uint32(dst)) + uint64(uint32(n))
	if size := m.MemSize.Load(); srcEnd > size || dstEnd > size {
		Wasm_trap_memcopy_oob()
	}
	copy(m.Memory[uint32(dst):uint32(dstEnd)], m.Memory[uint32(src):uint32(srcEnd)])
}

var spinAgents int32
var spinOversubscribed uint32

type ThreadPool struct {
	nextTID atomic.Int32
	wg      sync.WaitGroup

	parkMu sync.Mutex
	parked map[uint64][]chan struct{}
}

// wake releases up to count waiters on ea and reports how many it woke.
func (p *ThreadPool) wake(ea uint64, count int32) int32 {
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

// SaveGlobals returns the module's mutable globals, in a form that can be handed back
// to RestoreGlobals. It is how a snapshot of an instance captures the state that does not
// live in linear memory.
func SaveGlobals(m *Module) []uint64 {
	g := make([]uint64, 1881)
	g[0] = uint64(uint32(m.G1))
	g[1] = uint64(uint32(m.G3))
	g[2] = uint64(uint32(m.G4))
	g[3] = uint64(uint32(m.G5))
	g[4] = uint64(uint32(m.G6))
	g[5] = uint64(uint32(m.G7))
	g[6] = uint64(uint32(m.G8))
	g[7] = uint64(uint32(m.G9))
	g[8] = uint64(uint32(m.G10))
	g[9] = uint64(uint32(m.G11))
	g[10] = uint64(uint32(m.G12))
	g[11] = uint64(uint32(m.G13))
	g[12] = uint64(uint32(m.G14))
	g[13] = uint64(uint32(m.G15))
	g[14] = uint64(uint32(m.G16))
	g[15] = uint64(uint32(m.G17))
	g[16] = uint64(uint32(m.G18))
	g[17] = uint64(uint32(m.G19))
	g[18] = uint64(uint32(m.G20))
	g[19] = uint64(uint32(m.G21))
	g[20] = uint64(uint32(m.G22))
	g[21] = uint64(uint32(m.G23))
	g[22] = uint64(uint32(m.G24))
	g[23] = uint64(uint32(m.G25))
	g[24] = uint64(uint32(m.G26))
	g[25] = uint64(uint32(m.G27))
	g[26] = uint64(uint32(m.G28))
	g[27] = uint64(uint32(m.G29))
	g[28] = uint64(uint32(m.G30))
	g[29] = uint64(uint32(m.G31))
	g[30] = uint64(uint32(m.G32))
	g[31] = uint64(uint32(m.G33))
	g[32] = uint64(uint32(m.G34))
	g[33] = uint64(uint32(m.G35))
	g[34] = uint64(uint32(m.G36))
	g[35] = uint64(uint32(m.G37))
	g[36] = uint64(uint32(m.G38))
	g[37] = uint64(uint32(m.G39))
	g[38] = uint64(uint32(m.G40))
	g[39] = uint64(uint32(m.G41))
	g[40] = uint64(uint32(m.G42))
	g[41] = uint64(uint32(m.G43))
	g[42] = uint64(uint32(m.G44))
	g[43] = uint64(uint32(m.G45))
	g[44] = uint64(uint32(m.G46))
	g[45] = uint64(uint32(m.G47))
	g[46] = uint64(uint32(m.G48))
	g[47] = uint64(uint32(m.G49))
	g[48] = uint64(uint32(m.G50))
	g[49] = uint64(uint32(m.G51))
	g[50] = uint64(uint32(m.G52))
	g[51] = uint64(uint32(m.G53))
	g[52] = uint64(uint32(m.G54))
	g[53] = uint64(uint32(m.G55))
	g[54] = uint64(uint32(m.G56))
	g[55] = uint64(uint32(m.G57))
	g[56] = uint64(uint32(m.G58))
	g[57] = uint64(uint32(m.G59))
	g[58] = uint64(uint32(m.G60))
	g[59] = uint64(uint32(m.G61))
	g[60] = uint64(uint32(m.G62))
	g[61] = uint64(uint32(m.G63))
	g[62] = uint64(uint32(m.G64))
	g[63] = uint64(uint32(m.G65))
	g[64] = uint64(uint32(m.G66))
	g[65] = uint64(uint32(m.G67))
	g[66] = uint64(uint32(m.G68))
	g[67] = uint64(uint32(m.G69))
	g[68] = uint64(uint32(m.G70))
	g[69] = uint64(uint32(m.G71))
	g[70] = uint64(uint32(m.G72))
	g[71] = uint64(uint32(m.G73))
	g[72] = uint64(uint32(m.G74))
	g[73] = uint64(uint32(m.G75))
	g[74] = uint64(uint32(m.G76))
	g[75] = uint64(uint32(m.G77))
	g[76] = uint64(uint32(m.G78))
	g[77] = uint64(uint32(m.G79))
	g[78] = uint64(uint32(m.G80))
	g[79] = uint64(uint32(m.G81))
	g[80] = uint64(uint32(m.G82))
	g[81] = uint64(uint32(m.G83))
	g[82] = uint64(uint32(m.G84))
	g[83] = uint64(uint32(m.G85))
	g[84] = uint64(uint32(m.G86))
	g[85] = uint64(uint32(m.G87))
	g[86] = uint64(uint32(m.G88))
	g[87] = uint64(uint32(m.G89))
	g[88] = uint64(uint32(m.G90))
	g[89] = uint64(uint32(m.G91))
	g[90] = uint64(uint32(m.G92))
	g[91] = uint64(uint32(m.G93))
	g[92] = uint64(uint32(m.G94))
	g[93] = uint64(uint32(m.G95))
	g[94] = uint64(uint32(m.G96))
	g[95] = uint64(uint32(m.G97))
	g[96] = uint64(uint32(m.G98))
	g[97] = uint64(uint32(m.G99))
	g[98] = uint64(uint32(m.G100))
	g[99] = uint64(uint32(m.G101))
	g[100] = uint64(uint32(m.G102))
	g[101] = uint64(uint32(m.G103))
	g[102] = uint64(uint32(m.G104))
	g[103] = uint64(uint32(m.G105))
	g[104] = uint64(uint32(m.G106))
	g[105] = uint64(uint32(m.G107))
	g[106] = uint64(uint32(m.G108))
	g[107] = uint64(uint32(m.G109))
	g[108] = uint64(uint32(m.G110))
	g[109] = uint64(uint32(m.G111))
	g[110] = uint64(uint32(m.G112))
	g[111] = uint64(uint32(m.G113))
	g[112] = uint64(uint32(m.G114))
	g[113] = uint64(uint32(m.G115))
	g[114] = uint64(uint32(m.G116))
	g[115] = uint64(uint32(m.G117))
	g[116] = uint64(uint32(m.G118))
	g[117] = uint64(uint32(m.G119))
	g[118] = uint64(uint32(m.G120))
	g[119] = uint64(uint32(m.G121))
	g[120] = uint64(uint32(m.G122))
	g[121] = uint64(uint32(m.G123))
	g[122] = uint64(uint32(m.G124))
	g[123] = uint64(uint32(m.G125))
	g[124] = uint64(uint32(m.G126))
	g[125] = uint64(uint32(m.G127))
	g[126] = uint64(uint32(m.G128))
	g[127] = uint64(uint32(m.G129))
	g[128] = uint64(uint32(m.G130))
	g[129] = uint64(uint32(m.G131))
	g[130] = uint64(uint32(m.G132))
	g[131] = uint64(uint32(m.G133))
	g[132] = uint64(uint32(m.G134))
	g[133] = uint64(uint32(m.G135))
	g[134] = uint64(uint32(m.G136))
	g[135] = uint64(uint32(m.G137))
	g[136] = uint64(uint32(m.G138))
	g[137] = uint64(uint32(m.G139))
	g[138] = uint64(uint32(m.G140))
	g[139] = uint64(uint32(m.G141))
	g[140] = uint64(uint32(m.G142))
	g[141] = uint64(uint32(m.G143))
	g[142] = uint64(uint32(m.G144))
	g[143] = uint64(uint32(m.G145))
	g[144] = uint64(uint32(m.G146))
	g[145] = uint64(uint32(m.G147))
	g[146] = uint64(uint32(m.G148))
	g[147] = uint64(uint32(m.G149))
	g[148] = uint64(uint32(m.G150))
	g[149] = uint64(uint32(m.G151))
	g[150] = uint64(uint32(m.G152))
	g[151] = uint64(uint32(m.G153))
	g[152] = uint64(uint32(m.G154))
	g[153] = uint64(uint32(m.G155))
	g[154] = uint64(uint32(m.G156))
	g[155] = uint64(uint32(m.G157))
	g[156] = uint64(uint32(m.G158))
	g[157] = uint64(uint32(m.G159))
	g[158] = uint64(uint32(m.G160))
	g[159] = uint64(uint32(m.G161))
	g[160] = uint64(uint32(m.G162))
	g[161] = uint64(uint32(m.G163))
	g[162] = uint64(uint32(m.G164))
	g[163] = uint64(uint32(m.G165))
	g[164] = uint64(uint32(m.G166))
	g[165] = uint64(uint32(m.G167))
	g[166] = uint64(uint32(m.G168))
	g[167] = uint64(uint32(m.G169))
	g[168] = uint64(uint32(m.G170))
	g[169] = uint64(uint32(m.G171))
	g[170] = uint64(uint32(m.G172))
	g[171] = uint64(uint32(m.G173))
	g[172] = uint64(uint32(m.G174))
	g[173] = uint64(uint32(m.G175))
	g[174] = uint64(uint32(m.G176))
	g[175] = uint64(uint32(m.G177))
	g[176] = uint64(uint32(m.G178))
	g[177] = uint64(uint32(m.G179))
	g[178] = uint64(uint32(m.G180))
	g[179] = uint64(uint32(m.G181))
	g[180] = uint64(uint32(m.G182))
	g[181] = uint64(uint32(m.G183))
	g[182] = uint64(uint32(m.G184))
	g[183] = uint64(uint32(m.G185))
	g[184] = uint64(uint32(m.G186))
	g[185] = uint64(uint32(m.G187))
	g[186] = uint64(uint32(m.G188))
	g[187] = uint64(uint32(m.G189))
	g[188] = uint64(uint32(m.G190))
	g[189] = uint64(uint32(m.G191))
	g[190] = uint64(uint32(m.G192))
	g[191] = uint64(uint32(m.G193))
	g[192] = uint64(uint32(m.G194))
	g[193] = uint64(uint32(m.G195))
	g[194] = uint64(uint32(m.G196))
	g[195] = uint64(uint32(m.G197))
	g[196] = uint64(uint32(m.G198))
	g[197] = uint64(uint32(m.G199))
	g[198] = uint64(uint32(m.G200))
	g[199] = uint64(uint32(m.G201))
	g[200] = uint64(uint32(m.G202))
	g[201] = uint64(uint32(m.G203))
	g[202] = uint64(uint32(m.G204))
	g[203] = uint64(uint32(m.G205))
	g[204] = uint64(uint32(m.G206))
	g[205] = uint64(uint32(m.G207))
	g[206] = uint64(uint32(m.G208))
	g[207] = uint64(uint32(m.G209))
	g[208] = uint64(uint32(m.G210))
	g[209] = uint64(uint32(m.G211))
	g[210] = uint64(uint32(m.G212))
	g[211] = uint64(uint32(m.G213))
	g[212] = uint64(uint32(m.G214))
	g[213] = uint64(uint32(m.G215))
	g[214] = uint64(uint32(m.G216))
	g[215] = uint64(uint32(m.G217))
	g[216] = uint64(uint32(m.G218))
	g[217] = uint64(uint32(m.G219))
	g[218] = uint64(uint32(m.G220))
	g[219] = uint64(uint32(m.G221))
	g[220] = uint64(uint32(m.G222))
	g[221] = uint64(uint32(m.G223))
	g[222] = uint64(uint32(m.G224))
	g[223] = uint64(uint32(m.G225))
	g[224] = uint64(uint32(m.G226))
	g[225] = uint64(uint32(m.G227))
	g[226] = uint64(uint32(m.G228))
	g[227] = uint64(uint32(m.G229))
	g[228] = uint64(uint32(m.G230))
	g[229] = uint64(uint32(m.G231))
	g[230] = uint64(uint32(m.G232))
	g[231] = uint64(uint32(m.G233))
	g[232] = uint64(uint32(m.G234))
	g[233] = uint64(uint32(m.G235))
	g[234] = uint64(uint32(m.G236))
	g[235] = uint64(uint32(m.G237))
	g[236] = uint64(uint32(m.G238))
	g[237] = uint64(uint32(m.G239))
	g[238] = uint64(uint32(m.G240))
	g[239] = uint64(uint32(m.G241))
	g[240] = uint64(uint32(m.G242))
	g[241] = uint64(uint32(m.G243))
	g[242] = uint64(uint32(m.G244))
	g[243] = uint64(uint32(m.G245))
	g[244] = uint64(uint32(m.G246))
	g[245] = uint64(uint32(m.G247))
	g[246] = uint64(uint32(m.G248))
	g[247] = uint64(uint32(m.G249))
	g[248] = uint64(uint32(m.G250))
	g[249] = uint64(uint32(m.G251))
	g[250] = uint64(uint32(m.G252))
	g[251] = uint64(uint32(m.G253))
	g[252] = uint64(uint32(m.G254))
	g[253] = uint64(uint32(m.G255))
	g[254] = uint64(uint32(m.G256))
	g[255] = uint64(uint32(m.G257))
	g[256] = uint64(uint32(m.G258))
	g[257] = uint64(uint32(m.G259))
	g[258] = uint64(uint32(m.G260))
	g[259] = uint64(uint32(m.G261))
	g[260] = uint64(uint32(m.G262))
	g[261] = uint64(uint32(m.G263))
	g[262] = uint64(uint32(m.G264))
	g[263] = uint64(uint32(m.G265))
	g[264] = uint64(uint32(m.G266))
	g[265] = uint64(uint32(m.G267))
	g[266] = uint64(uint32(m.G268))
	g[267] = uint64(uint32(m.G269))
	g[268] = uint64(uint32(m.G270))
	g[269] = uint64(uint32(m.G271))
	g[270] = uint64(uint32(m.G272))
	g[271] = uint64(uint32(m.G273))
	g[272] = uint64(uint32(m.G274))
	g[273] = uint64(uint32(m.G275))
	g[274] = uint64(uint32(m.G276))
	g[275] = uint64(uint32(m.G277))
	g[276] = uint64(uint32(m.G278))
	g[277] = uint64(uint32(m.G279))
	g[278] = uint64(uint32(m.G280))
	g[279] = uint64(uint32(m.G281))
	g[280] = uint64(uint32(m.G282))
	g[281] = uint64(uint32(m.G283))
	g[282] = uint64(uint32(m.G284))
	g[283] = uint64(uint32(m.G285))
	g[284] = uint64(uint32(m.G286))
	g[285] = uint64(uint32(m.G287))
	g[286] = uint64(uint32(m.G288))
	g[287] = uint64(uint32(m.G289))
	g[288] = uint64(uint32(m.G290))
	g[289] = uint64(uint32(m.G291))
	g[290] = uint64(uint32(m.G292))
	g[291] = uint64(uint32(m.G293))
	g[292] = uint64(uint32(m.G294))
	g[293] = uint64(uint32(m.G295))
	g[294] = uint64(uint32(m.G296))
	g[295] = uint64(uint32(m.G297))
	g[296] = uint64(uint32(m.G298))
	g[297] = uint64(uint32(m.G299))
	g[298] = uint64(uint32(m.G300))
	g[299] = uint64(uint32(m.G301))
	g[300] = uint64(uint32(m.G302))
	g[301] = uint64(uint32(m.G303))
	g[302] = uint64(uint32(m.G304))
	g[303] = uint64(uint32(m.G305))
	g[304] = uint64(uint32(m.G306))
	g[305] = uint64(uint32(m.G307))
	g[306] = uint64(uint32(m.G308))
	g[307] = uint64(uint32(m.G309))
	g[308] = uint64(uint32(m.G310))
	g[309] = uint64(uint32(m.G311))
	g[310] = uint64(uint32(m.G312))
	g[311] = uint64(uint32(m.G313))
	g[312] = uint64(uint32(m.G314))
	g[313] = uint64(uint32(m.G315))
	g[314] = uint64(uint32(m.G316))
	g[315] = uint64(uint32(m.G317))
	g[316] = uint64(uint32(m.G318))
	g[317] = uint64(uint32(m.G319))
	g[318] = uint64(uint32(m.G320))
	g[319] = uint64(uint32(m.G321))
	g[320] = uint64(uint32(m.G322))
	g[321] = uint64(uint32(m.G323))
	g[322] = uint64(uint32(m.G324))
	g[323] = uint64(uint32(m.G325))
	g[324] = uint64(uint32(m.G326))
	g[325] = uint64(uint32(m.G327))
	g[326] = uint64(uint32(m.G328))
	g[327] = uint64(uint32(m.G329))
	g[328] = uint64(uint32(m.G330))
	g[329] = uint64(uint32(m.G331))
	g[330] = uint64(uint32(m.G332))
	g[331] = uint64(uint32(m.G333))
	g[332] = uint64(uint32(m.G334))
	g[333] = uint64(uint32(m.G335))
	g[334] = uint64(uint32(m.G336))
	g[335] = uint64(uint32(m.G337))
	g[336] = uint64(uint32(m.G338))
	g[337] = uint64(uint32(m.G339))
	g[338] = uint64(uint32(m.G340))
	g[339] = uint64(uint32(m.G341))
	g[340] = uint64(uint32(m.G342))
	g[341] = uint64(uint32(m.G343))
	g[342] = uint64(uint32(m.G344))
	g[343] = uint64(uint32(m.G345))
	g[344] = uint64(uint32(m.G346))
	g[345] = uint64(uint32(m.G347))
	g[346] = uint64(uint32(m.G348))
	g[347] = uint64(uint32(m.G349))
	g[348] = uint64(uint32(m.G350))
	g[349] = uint64(uint32(m.G351))
	g[350] = uint64(uint32(m.G352))
	g[351] = uint64(uint32(m.G353))
	g[352] = uint64(uint32(m.G354))
	g[353] = uint64(uint32(m.G355))
	g[354] = uint64(uint32(m.G356))
	g[355] = uint64(uint32(m.G357))
	g[356] = uint64(uint32(m.G358))
	g[357] = uint64(uint32(m.G359))
	g[358] = uint64(uint32(m.G360))
	g[359] = uint64(uint32(m.G361))
	g[360] = uint64(uint32(m.G362))
	g[361] = uint64(uint32(m.G363))
	g[362] = uint64(uint32(m.G364))
	g[363] = uint64(uint32(m.G365))
	g[364] = uint64(uint32(m.G366))
	g[365] = uint64(uint32(m.G367))
	g[366] = uint64(uint32(m.G368))
	g[367] = uint64(uint32(m.G369))
	g[368] = uint64(uint32(m.G370))
	g[369] = uint64(uint32(m.G371))
	g[370] = uint64(uint32(m.G372))
	g[371] = uint64(uint32(m.G373))
	g[372] = uint64(uint32(m.G374))
	g[373] = uint64(uint32(m.G375))
	g[374] = uint64(uint32(m.G376))
	g[375] = uint64(uint32(m.G377))
	g[376] = uint64(uint32(m.G378))
	g[377] = uint64(uint32(m.G379))
	g[378] = uint64(uint32(m.G380))
	g[379] = uint64(uint32(m.G381))
	g[380] = uint64(uint32(m.G382))
	g[381] = uint64(uint32(m.G383))
	g[382] = uint64(uint32(m.G384))
	g[383] = uint64(uint32(m.G385))
	g[384] = uint64(uint32(m.G386))
	g[385] = uint64(uint32(m.G387))
	g[386] = uint64(uint32(m.G388))
	g[387] = uint64(uint32(m.G389))
	g[388] = uint64(uint32(m.G390))
	g[389] = uint64(uint32(m.G391))
	g[390] = uint64(uint32(m.G392))
	g[391] = uint64(uint32(m.G393))
	g[392] = uint64(uint32(m.G394))
	g[393] = uint64(uint32(m.G395))
	g[394] = uint64(uint32(m.G396))
	g[395] = uint64(uint32(m.G397))
	g[396] = uint64(uint32(m.G398))
	g[397] = uint64(uint32(m.G399))
	g[398] = uint64(uint32(m.G400))
	g[399] = uint64(uint32(m.G401))
	g[400] = uint64(uint32(m.G402))
	g[401] = uint64(uint32(m.G403))
	g[402] = uint64(uint32(m.G404))
	g[403] = uint64(uint32(m.G405))
	g[404] = uint64(uint32(m.G406))
	g[405] = uint64(uint32(m.G407))
	g[406] = uint64(uint32(m.G408))
	g[407] = uint64(uint32(m.G409))
	g[408] = uint64(uint32(m.G410))
	g[409] = uint64(uint32(m.G411))
	g[410] = uint64(uint32(m.G412))
	g[411] = uint64(uint32(m.G413))
	g[412] = uint64(uint32(m.G414))
	g[413] = uint64(uint32(m.G415))
	g[414] = uint64(uint32(m.G416))
	g[415] = uint64(uint32(m.G417))
	g[416] = uint64(uint32(m.G418))
	g[417] = uint64(uint32(m.G419))
	g[418] = uint64(uint32(m.G420))
	g[419] = uint64(uint32(m.G421))
	g[420] = uint64(uint32(m.G422))
	g[421] = uint64(uint32(m.G423))
	g[422] = uint64(uint32(m.G424))
	g[423] = uint64(uint32(m.G425))
	g[424] = uint64(uint32(m.G426))
	g[425] = uint64(uint32(m.G427))
	g[426] = uint64(uint32(m.G428))
	g[427] = uint64(uint32(m.G429))
	g[428] = uint64(uint32(m.G430))
	g[429] = uint64(uint32(m.G431))
	g[430] = uint64(uint32(m.G432))
	g[431] = uint64(uint32(m.G433))
	g[432] = uint64(uint32(m.G434))
	g[433] = uint64(uint32(m.G435))
	g[434] = uint64(uint32(m.G436))
	g[435] = uint64(uint32(m.G437))
	g[436] = uint64(uint32(m.G438))
	g[437] = uint64(uint32(m.G439))
	g[438] = uint64(uint32(m.G440))
	g[439] = uint64(uint32(m.G441))
	g[440] = uint64(uint32(m.G442))
	g[441] = uint64(uint32(m.G443))
	g[442] = uint64(uint32(m.G444))
	g[443] = uint64(uint32(m.G445))
	g[444] = uint64(uint32(m.G446))
	g[445] = uint64(uint32(m.G447))
	g[446] = uint64(uint32(m.G448))
	g[447] = uint64(uint32(m.G449))
	g[448] = uint64(uint32(m.G450))
	g[449] = uint64(uint32(m.G451))
	g[450] = uint64(uint32(m.G452))
	g[451] = uint64(uint32(m.G453))
	g[452] = uint64(uint32(m.G454))
	g[453] = uint64(uint32(m.G455))
	g[454] = uint64(uint32(m.G456))
	g[455] = uint64(uint32(m.G457))
	g[456] = uint64(uint32(m.G458))
	g[457] = uint64(uint32(m.G459))
	g[458] = uint64(uint32(m.G460))
	g[459] = uint64(uint32(m.G461))
	g[460] = uint64(uint32(m.G462))
	g[461] = uint64(uint32(m.G463))
	g[462] = uint64(uint32(m.G464))
	g[463] = uint64(uint32(m.G465))
	g[464] = uint64(uint32(m.G466))
	g[465] = uint64(uint32(m.G467))
	g[466] = uint64(uint32(m.G468))
	g[467] = uint64(uint32(m.G469))
	g[468] = uint64(uint32(m.G470))
	g[469] = uint64(uint32(m.G471))
	g[470] = uint64(uint32(m.G472))
	g[471] = uint64(uint32(m.G473))
	g[472] = uint64(uint32(m.G474))
	g[473] = uint64(uint32(m.G475))
	g[474] = uint64(uint32(m.G476))
	g[475] = uint64(uint32(m.G477))
	g[476] = uint64(uint32(m.G478))
	g[477] = uint64(uint32(m.G479))
	g[478] = uint64(uint32(m.G480))
	g[479] = uint64(uint32(m.G481))
	g[480] = uint64(uint32(m.G482))
	g[481] = uint64(uint32(m.G483))
	g[482] = uint64(uint32(m.G484))
	g[483] = uint64(uint32(m.G485))
	g[484] = uint64(uint32(m.G486))
	g[485] = uint64(uint32(m.G487))
	g[486] = uint64(uint32(m.G488))
	g[487] = uint64(uint32(m.G489))
	g[488] = uint64(uint32(m.G490))
	g[489] = uint64(uint32(m.G491))
	g[490] = uint64(uint32(m.G492))
	g[491] = uint64(uint32(m.G493))
	g[492] = uint64(uint32(m.G494))
	g[493] = uint64(uint32(m.G495))
	g[494] = uint64(uint32(m.G496))
	g[495] = uint64(uint32(m.G497))
	g[496] = uint64(uint32(m.G498))
	g[497] = uint64(uint32(m.G499))
	g[498] = uint64(uint32(m.G500))
	g[499] = uint64(uint32(m.G501))
	g[500] = uint64(uint32(m.G502))
	g[501] = uint64(uint32(m.G503))
	g[502] = uint64(uint32(m.G504))
	g[503] = uint64(uint32(m.G505))
	g[504] = uint64(uint32(m.G506))
	g[505] = uint64(uint32(m.G507))
	g[506] = uint64(uint32(m.G508))
	g[507] = uint64(uint32(m.G509))
	g[508] = uint64(uint32(m.G510))
	g[509] = uint64(uint32(m.G511))
	g[510] = uint64(uint32(m.G512))
	g[511] = uint64(uint32(m.G513))
	g[512] = uint64(uint32(m.G514))
	g[513] = uint64(uint32(m.G515))
	g[514] = uint64(uint32(m.G516))
	g[515] = uint64(uint32(m.G517))
	g[516] = uint64(uint32(m.G518))
	g[517] = uint64(uint32(m.G519))
	g[518] = uint64(uint32(m.G520))
	g[519] = uint64(uint32(m.G521))
	g[520] = uint64(uint32(m.G522))
	g[521] = uint64(uint32(m.G523))
	g[522] = uint64(uint32(m.G524))
	g[523] = uint64(uint32(m.G525))
	g[524] = uint64(uint32(m.G526))
	g[525] = uint64(uint32(m.G527))
	g[526] = uint64(uint32(m.G528))
	g[527] = uint64(uint32(m.G529))
	g[528] = uint64(uint32(m.G530))
	g[529] = uint64(uint32(m.G531))
	g[530] = uint64(uint32(m.G532))
	g[531] = uint64(uint32(m.G533))
	g[532] = uint64(uint32(m.G534))
	g[533] = uint64(uint32(m.G535))
	g[534] = uint64(uint32(m.G536))
	g[535] = uint64(uint32(m.G537))
	g[536] = uint64(uint32(m.G538))
	g[537] = uint64(uint32(m.G539))
	g[538] = uint64(uint32(m.G540))
	g[539] = uint64(uint32(m.G541))
	g[540] = uint64(uint32(m.G542))
	g[541] = uint64(uint32(m.G543))
	g[542] = uint64(uint32(m.G544))
	g[543] = uint64(uint32(m.G545))
	g[544] = uint64(uint32(m.G546))
	g[545] = uint64(uint32(m.G547))
	g[546] = uint64(uint32(m.G548))
	g[547] = uint64(uint32(m.G549))
	g[548] = uint64(uint32(m.G550))
	g[549] = uint64(uint32(m.G551))
	g[550] = uint64(uint32(m.G552))
	g[551] = uint64(uint32(m.G553))
	g[552] = uint64(uint32(m.G554))
	g[553] = uint64(uint32(m.G555))
	g[554] = uint64(uint32(m.G556))
	g[555] = uint64(uint32(m.G557))
	g[556] = uint64(uint32(m.G558))
	g[557] = uint64(uint32(m.G559))
	g[558] = uint64(uint32(m.G560))
	g[559] = uint64(uint32(m.G561))
	g[560] = uint64(uint32(m.G562))
	g[561] = uint64(uint32(m.G563))
	g[562] = uint64(uint32(m.G564))
	g[563] = uint64(uint32(m.G565))
	g[564] = uint64(uint32(m.G566))
	g[565] = uint64(uint32(m.G567))
	g[566] = uint64(uint32(m.G568))
	g[567] = uint64(uint32(m.G569))
	g[568] = uint64(uint32(m.G570))
	g[569] = uint64(uint32(m.G571))
	g[570] = uint64(uint32(m.G572))
	g[571] = uint64(uint32(m.G573))
	g[572] = uint64(uint32(m.G574))
	g[573] = uint64(uint32(m.G575))
	g[574] = uint64(uint32(m.G576))
	g[575] = uint64(uint32(m.G577))
	g[576] = uint64(uint32(m.G578))
	g[577] = uint64(uint32(m.G579))
	g[578] = uint64(uint32(m.G580))
	g[579] = uint64(uint32(m.G581))
	g[580] = uint64(uint32(m.G582))
	g[581] = uint64(uint32(m.G583))
	g[582] = uint64(uint32(m.G584))
	g[583] = uint64(uint32(m.G585))
	g[584] = uint64(uint32(m.G586))
	g[585] = uint64(uint32(m.G587))
	g[586] = uint64(uint32(m.G588))
	g[587] = uint64(uint32(m.G589))
	g[588] = uint64(uint32(m.G590))
	g[589] = uint64(uint32(m.G591))
	g[590] = uint64(uint32(m.G592))
	g[591] = uint64(uint32(m.G593))
	g[592] = uint64(uint32(m.G594))
	g[593] = uint64(uint32(m.G595))
	g[594] = uint64(uint32(m.G596))
	g[595] = uint64(uint32(m.G597))
	g[596] = uint64(uint32(m.G598))
	g[597] = uint64(uint32(m.G599))
	g[598] = uint64(uint32(m.G600))
	g[599] = uint64(uint32(m.G601))
	g[600] = uint64(uint32(m.G602))
	g[601] = uint64(uint32(m.G603))
	g[602] = uint64(uint32(m.G604))
	g[603] = uint64(uint32(m.G605))
	g[604] = uint64(uint32(m.G606))
	g[605] = uint64(uint32(m.G607))
	g[606] = uint64(uint32(m.G608))
	g[607] = uint64(uint32(m.G609))
	g[608] = uint64(uint32(m.G610))
	g[609] = uint64(uint32(m.G611))
	g[610] = uint64(uint32(m.G612))
	g[611] = uint64(uint32(m.G613))
	g[612] = uint64(uint32(m.G614))
	g[613] = uint64(uint32(m.G615))
	g[614] = uint64(uint32(m.G616))
	g[615] = uint64(uint32(m.G617))
	g[616] = uint64(uint32(m.G618))
	g[617] = uint64(uint32(m.G619))
	g[618] = uint64(uint32(m.G620))
	g[619] = uint64(uint32(m.G621))
	g[620] = uint64(uint32(m.G622))
	g[621] = uint64(uint32(m.G623))
	g[622] = uint64(uint32(m.G624))
	g[623] = uint64(uint32(m.G625))
	g[624] = uint64(uint32(m.G626))
	g[625] = uint64(uint32(m.G627))
	g[626] = uint64(uint32(m.G628))
	g[627] = uint64(uint32(m.G629))
	g[628] = uint64(uint32(m.G630))
	g[629] = uint64(uint32(m.G631))
	g[630] = uint64(uint32(m.G632))
	g[631] = uint64(uint32(m.G633))
	g[632] = uint64(uint32(m.G634))
	g[633] = uint64(uint32(m.G635))
	g[634] = uint64(uint32(m.G636))
	g[635] = uint64(uint32(m.G637))
	g[636] = uint64(uint32(m.G638))
	g[637] = uint64(uint32(m.G639))
	g[638] = uint64(uint32(m.G640))
	g[639] = uint64(uint32(m.G641))
	g[640] = uint64(uint32(m.G642))
	g[641] = uint64(uint32(m.G643))
	g[642] = uint64(uint32(m.G644))
	g[643] = uint64(uint32(m.G645))
	g[644] = uint64(uint32(m.G646))
	g[645] = uint64(uint32(m.G647))
	g[646] = uint64(uint32(m.G648))
	g[647] = uint64(uint32(m.G649))
	g[648] = uint64(uint32(m.G650))
	g[649] = uint64(uint32(m.G651))
	g[650] = uint64(uint32(m.G652))
	g[651] = uint64(uint32(m.G653))
	g[652] = uint64(uint32(m.G654))
	g[653] = uint64(uint32(m.G655))
	g[654] = uint64(uint32(m.G656))
	g[655] = uint64(uint32(m.G657))
	g[656] = uint64(uint32(m.G658))
	g[657] = uint64(uint32(m.G659))
	g[658] = uint64(uint32(m.G660))
	g[659] = uint64(uint32(m.G661))
	g[660] = uint64(uint32(m.G662))
	g[661] = uint64(uint32(m.G663))
	g[662] = uint64(uint32(m.G664))
	g[663] = uint64(uint32(m.G665))
	g[664] = uint64(uint32(m.G666))
	g[665] = uint64(uint32(m.G667))
	g[666] = uint64(uint32(m.G668))
	g[667] = uint64(uint32(m.G669))
	g[668] = uint64(uint32(m.G670))
	g[669] = uint64(uint32(m.G671))
	g[670] = uint64(uint32(m.G672))
	g[671] = uint64(uint32(m.G673))
	g[672] = uint64(uint32(m.G674))
	g[673] = uint64(uint32(m.G675))
	g[674] = uint64(uint32(m.G676))
	g[675] = uint64(uint32(m.G677))
	g[676] = uint64(uint32(m.G678))
	g[677] = uint64(uint32(m.G679))
	g[678] = uint64(uint32(m.G680))
	g[679] = uint64(uint32(m.G681))
	g[680] = uint64(uint32(m.G682))
	g[681] = uint64(uint32(m.G683))
	g[682] = uint64(uint32(m.G684))
	g[683] = uint64(uint32(m.G685))
	g[684] = uint64(uint32(m.G686))
	g[685] = uint64(uint32(m.G687))
	g[686] = uint64(uint32(m.G688))
	g[687] = uint64(uint32(m.G689))
	g[688] = uint64(uint32(m.G690))
	g[689] = uint64(uint32(m.G691))
	g[690] = uint64(uint32(m.G692))
	g[691] = uint64(uint32(m.G693))
	g[692] = uint64(uint32(m.G694))
	g[693] = uint64(uint32(m.G695))
	g[694] = uint64(uint32(m.G696))
	g[695] = uint64(uint32(m.G697))
	g[696] = uint64(uint32(m.G698))
	g[697] = uint64(uint32(m.G699))
	g[698] = uint64(uint32(m.G700))
	g[699] = uint64(uint32(m.G701))
	g[700] = uint64(uint32(m.G702))
	g[701] = uint64(uint32(m.G703))
	g[702] = uint64(uint32(m.G704))
	g[703] = uint64(uint32(m.G705))
	g[704] = uint64(uint32(m.G706))
	g[705] = uint64(uint32(m.G707))
	g[706] = uint64(uint32(m.G708))
	g[707] = uint64(uint32(m.G709))
	g[708] = uint64(uint32(m.G710))
	g[709] = uint64(uint32(m.G711))
	g[710] = uint64(uint32(m.G712))
	g[711] = uint64(uint32(m.G713))
	g[712] = uint64(uint32(m.G714))
	g[713] = uint64(uint32(m.G715))
	g[714] = uint64(uint32(m.G716))
	g[715] = uint64(uint32(m.G717))
	g[716] = uint64(uint32(m.G718))
	g[717] = uint64(uint32(m.G719))
	g[718] = uint64(uint32(m.G720))
	g[719] = uint64(uint32(m.G721))
	g[720] = uint64(uint32(m.G722))
	g[721] = uint64(uint32(m.G723))
	g[722] = uint64(uint32(m.G724))
	g[723] = uint64(uint32(m.G725))
	g[724] = uint64(uint32(m.G726))
	g[725] = uint64(uint32(m.G727))
	g[726] = uint64(uint32(m.G728))
	g[727] = uint64(uint32(m.G729))
	g[728] = uint64(uint32(m.G730))
	g[729] = uint64(uint32(m.G731))
	g[730] = uint64(uint32(m.G732))
	g[731] = uint64(uint32(m.G733))
	g[732] = uint64(uint32(m.G734))
	g[733] = uint64(uint32(m.G735))
	g[734] = uint64(uint32(m.G736))
	g[735] = uint64(uint32(m.G737))
	g[736] = uint64(uint32(m.G738))
	g[737] = uint64(uint32(m.G739))
	g[738] = uint64(uint32(m.G740))
	g[739] = uint64(uint32(m.G741))
	g[740] = uint64(uint32(m.G742))
	g[741] = uint64(uint32(m.G743))
	g[742] = uint64(uint32(m.G744))
	g[743] = uint64(uint32(m.G745))
	g[744] = uint64(uint32(m.G746))
	g[745] = uint64(uint32(m.G747))
	g[746] = uint64(uint32(m.G748))
	g[747] = uint64(uint32(m.G749))
	g[748] = uint64(uint32(m.G750))
	g[749] = uint64(uint32(m.G751))
	g[750] = uint64(uint32(m.G752))
	g[751] = uint64(uint32(m.G753))
	g[752] = uint64(uint32(m.G754))
	g[753] = uint64(uint32(m.G755))
	g[754] = uint64(uint32(m.G756))
	g[755] = uint64(uint32(m.G757))
	g[756] = uint64(uint32(m.G758))
	g[757] = uint64(uint32(m.G759))
	g[758] = uint64(uint32(m.G760))
	g[759] = uint64(uint32(m.G761))
	g[760] = uint64(uint32(m.G762))
	g[761] = uint64(uint32(m.G763))
	g[762] = uint64(uint32(m.G764))
	g[763] = uint64(uint32(m.G765))
	g[764] = uint64(uint32(m.G766))
	g[765] = uint64(uint32(m.G767))
	g[766] = uint64(uint32(m.G768))
	g[767] = uint64(uint32(m.G769))
	g[768] = uint64(uint32(m.G770))
	g[769] = uint64(uint32(m.G771))
	g[770] = uint64(uint32(m.G772))
	g[771] = uint64(uint32(m.G773))
	g[772] = uint64(uint32(m.G774))
	g[773] = uint64(uint32(m.G775))
	g[774] = uint64(uint32(m.G776))
	g[775] = uint64(uint32(m.G777))
	g[776] = uint64(uint32(m.G778))
	g[777] = uint64(uint32(m.G779))
	g[778] = uint64(uint32(m.G780))
	g[779] = uint64(uint32(m.G781))
	g[780] = uint64(uint32(m.G782))
	g[781] = uint64(uint32(m.G783))
	g[782] = uint64(uint32(m.G784))
	g[783] = uint64(uint32(m.G785))
	g[784] = uint64(uint32(m.G786))
	g[785] = uint64(uint32(m.G787))
	g[786] = uint64(uint32(m.G788))
	g[787] = uint64(uint32(m.G789))
	g[788] = uint64(uint32(m.G790))
	g[789] = uint64(uint32(m.G791))
	g[790] = uint64(uint32(m.G792))
	g[791] = uint64(uint32(m.G793))
	g[792] = uint64(uint32(m.G794))
	g[793] = uint64(uint32(m.G795))
	g[794] = uint64(uint32(m.G796))
	g[795] = uint64(uint32(m.G797))
	g[796] = uint64(uint32(m.G798))
	g[797] = uint64(uint32(m.G799))
	g[798] = uint64(uint32(m.G800))
	g[799] = uint64(uint32(m.G801))
	g[800] = uint64(uint32(m.G802))
	g[801] = uint64(uint32(m.G803))
	g[802] = uint64(uint32(m.G804))
	g[803] = uint64(uint32(m.G805))
	g[804] = uint64(uint32(m.G806))
	g[805] = uint64(uint32(m.G807))
	g[806] = uint64(uint32(m.G808))
	g[807] = uint64(uint32(m.G809))
	g[808] = uint64(uint32(m.G810))
	g[809] = uint64(uint32(m.G811))
	g[810] = uint64(uint32(m.G812))
	g[811] = uint64(uint32(m.G813))
	g[812] = uint64(uint32(m.G814))
	g[813] = uint64(uint32(m.G815))
	g[814] = uint64(uint32(m.G816))
	g[815] = uint64(uint32(m.G817))
	g[816] = uint64(uint32(m.G818))
	g[817] = uint64(uint32(m.G819))
	g[818] = uint64(uint32(m.G820))
	g[819] = uint64(uint32(m.G821))
	g[820] = uint64(uint32(m.G822))
	g[821] = uint64(uint32(m.G823))
	g[822] = uint64(uint32(m.G824))
	g[823] = uint64(uint32(m.G825))
	g[824] = uint64(uint32(m.G826))
	g[825] = uint64(uint32(m.G827))
	g[826] = uint64(uint32(m.G828))
	g[827] = uint64(uint32(m.G829))
	g[828] = uint64(uint32(m.G830))
	g[829] = uint64(uint32(m.G831))
	g[830] = uint64(uint32(m.G832))
	g[831] = uint64(uint32(m.G833))
	g[832] = uint64(uint32(m.G834))
	g[833] = uint64(uint32(m.G835))
	g[834] = uint64(uint32(m.G836))
	g[835] = uint64(uint32(m.G837))
	g[836] = uint64(uint32(m.G838))
	g[837] = uint64(uint32(m.G839))
	g[838] = uint64(uint32(m.G840))
	g[839] = uint64(uint32(m.G841))
	g[840] = uint64(uint32(m.G842))
	g[841] = uint64(uint32(m.G843))
	g[842] = uint64(uint32(m.G844))
	g[843] = uint64(uint32(m.G845))
	g[844] = uint64(uint32(m.G846))
	g[845] = uint64(uint32(m.G847))
	g[846] = uint64(uint32(m.G848))
	g[847] = uint64(uint32(m.G849))
	g[848] = uint64(uint32(m.G850))
	g[849] = uint64(uint32(m.G851))
	g[850] = uint64(uint32(m.G852))
	g[851] = uint64(uint32(m.G853))
	g[852] = uint64(uint32(m.G854))
	g[853] = uint64(uint32(m.G855))
	g[854] = uint64(uint32(m.G856))
	g[855] = uint64(uint32(m.G857))
	g[856] = uint64(uint32(m.G858))
	g[857] = uint64(uint32(m.G859))
	g[858] = uint64(uint32(m.G860))
	g[859] = uint64(uint32(m.G861))
	g[860] = uint64(uint32(m.G862))
	g[861] = uint64(uint32(m.G863))
	g[862] = uint64(uint32(m.G864))
	g[863] = uint64(uint32(m.G865))
	g[864] = uint64(uint32(m.G866))
	g[865] = uint64(uint32(m.G867))
	g[866] = uint64(uint32(m.G868))
	g[867] = uint64(uint32(m.G869))
	g[868] = uint64(uint32(m.G870))
	g[869] = uint64(uint32(m.G871))
	g[870] = uint64(uint32(m.G872))
	g[871] = uint64(uint32(m.G873))
	g[872] = uint64(uint32(m.G874))
	g[873] = uint64(uint32(m.G875))
	g[874] = uint64(uint32(m.G876))
	g[875] = uint64(uint32(m.G877))
	g[876] = uint64(uint32(m.G878))
	g[877] = uint64(uint32(m.G879))
	g[878] = uint64(uint32(m.G880))
	g[879] = uint64(uint32(m.G881))
	g[880] = uint64(uint32(m.G882))
	g[881] = uint64(uint32(m.G883))
	g[882] = uint64(uint32(m.G884))
	g[883] = uint64(uint32(m.G885))
	g[884] = uint64(uint32(m.G886))
	g[885] = uint64(uint32(m.G887))
	g[886] = uint64(uint32(m.G888))
	g[887] = uint64(uint32(m.G889))
	g[888] = uint64(uint32(m.G890))
	g[889] = uint64(uint32(m.G891))
	g[890] = uint64(uint32(m.G892))
	g[891] = uint64(uint32(m.G893))
	g[892] = uint64(uint32(m.G894))
	g[893] = uint64(uint32(m.G895))
	g[894] = uint64(uint32(m.G896))
	g[895] = uint64(uint32(m.G897))
	g[896] = uint64(uint32(m.G898))
	g[897] = uint64(uint32(m.G899))
	g[898] = uint64(uint32(m.G900))
	g[899] = uint64(uint32(m.G901))
	g[900] = uint64(uint32(m.G902))
	g[901] = uint64(uint32(m.G903))
	g[902] = uint64(uint32(m.G904))
	g[903] = uint64(uint32(m.G905))
	g[904] = uint64(uint32(m.G906))
	g[905] = uint64(uint32(m.G907))
	g[906] = uint64(uint32(m.G908))
	g[907] = uint64(uint32(m.G909))
	g[908] = uint64(uint32(m.G910))
	g[909] = uint64(uint32(m.G911))
	g[910] = uint64(uint32(m.G912))
	g[911] = uint64(uint32(m.G913))
	g[912] = uint64(uint32(m.G914))
	g[913] = uint64(uint32(m.G915))
	g[914] = uint64(uint32(m.G916))
	g[915] = uint64(uint32(m.G917))
	g[916] = uint64(uint32(m.G918))
	g[917] = uint64(uint32(m.G919))
	g[918] = uint64(uint32(m.G920))
	g[919] = uint64(uint32(m.G921))
	g[920] = uint64(uint32(m.G922))
	g[921] = uint64(uint32(m.G923))
	g[922] = uint64(uint32(m.G924))
	g[923] = uint64(uint32(m.G925))
	g[924] = uint64(uint32(m.G926))
	g[925] = uint64(uint32(m.G927))
	g[926] = uint64(uint32(m.G928))
	g[927] = uint64(uint32(m.G929))
	g[928] = uint64(uint32(m.G930))
	g[929] = uint64(uint32(m.G931))
	g[930] = uint64(uint32(m.G932))
	g[931] = uint64(uint32(m.G933))
	g[932] = uint64(uint32(m.G934))
	g[933] = uint64(uint32(m.G935))
	g[934] = uint64(uint32(m.G936))
	g[935] = uint64(uint32(m.G937))
	g[936] = uint64(uint32(m.G938))
	g[937] = uint64(uint32(m.G939))
	g[938] = uint64(uint32(m.G940))
	g[939] = uint64(uint32(m.G941))
	g[940] = uint64(uint32(m.G942))
	g[941] = uint64(uint32(m.G943))
	g[942] = uint64(uint32(m.G944))
	g[943] = uint64(uint32(m.G945))
	g[944] = uint64(uint32(m.G946))
	g[945] = uint64(uint32(m.G947))
	g[946] = uint64(uint32(m.G948))
	g[947] = uint64(uint32(m.G949))
	g[948] = uint64(uint32(m.G950))
	g[949] = uint64(uint32(m.G951))
	g[950] = uint64(uint32(m.G952))
	g[951] = uint64(uint32(m.G953))
	g[952] = uint64(uint32(m.G954))
	g[953] = uint64(uint32(m.G955))
	g[954] = uint64(uint32(m.G956))
	g[955] = uint64(uint32(m.G957))
	g[956] = uint64(uint32(m.G958))
	g[957] = uint64(uint32(m.G959))
	g[958] = uint64(uint32(m.G960))
	g[959] = uint64(uint32(m.G961))
	g[960] = uint64(uint32(m.G962))
	g[961] = uint64(uint32(m.G963))
	g[962] = uint64(uint32(m.G964))
	g[963] = uint64(uint32(m.G965))
	g[964] = uint64(uint32(m.G966))
	g[965] = uint64(uint32(m.G967))
	g[966] = uint64(uint32(m.G968))
	g[967] = uint64(uint32(m.G969))
	g[968] = uint64(uint32(m.G970))
	g[969] = uint64(uint32(m.G971))
	g[970] = uint64(uint32(m.G972))
	g[971] = uint64(uint32(m.G973))
	g[972] = uint64(uint32(m.G974))
	g[973] = uint64(uint32(m.G975))
	g[974] = uint64(uint32(m.G976))
	g[975] = uint64(uint32(m.G977))
	g[976] = uint64(uint32(m.G978))
	g[977] = uint64(uint32(m.G979))
	g[978] = uint64(uint32(m.G980))
	g[979] = uint64(uint32(m.G981))
	g[980] = uint64(uint32(m.G982))
	g[981] = uint64(uint32(m.G983))
	g[982] = uint64(uint32(m.G984))
	g[983] = uint64(uint32(m.G985))
	g[984] = uint64(uint32(m.G986))
	g[985] = uint64(uint32(m.G987))
	g[986] = uint64(uint32(m.G988))
	g[987] = uint64(uint32(m.G989))
	g[988] = uint64(uint32(m.G990))
	g[989] = uint64(uint32(m.G991))
	g[990] = uint64(uint32(m.G992))
	g[991] = uint64(uint32(m.G993))
	g[992] = uint64(uint32(m.G994))
	g[993] = uint64(uint32(m.G995))
	g[994] = uint64(uint32(m.G996))
	g[995] = uint64(uint32(m.G997))
	g[996] = uint64(uint32(m.G998))
	g[997] = uint64(uint32(m.G999))
	g[998] = uint64(uint32(m.G1000))
	g[999] = uint64(uint32(m.G1001))
	g[1000] = uint64(uint32(m.G1002))
	g[1001] = uint64(uint32(m.G1003))
	g[1002] = uint64(uint32(m.G1004))
	g[1003] = uint64(uint32(m.G1005))
	g[1004] = uint64(uint32(m.G1006))
	g[1005] = uint64(uint32(m.G1007))
	g[1006] = uint64(uint32(m.G1008))
	g[1007] = uint64(uint32(m.G1009))
	g[1008] = uint64(uint32(m.G1010))
	g[1009] = uint64(uint32(m.G1011))
	g[1010] = uint64(uint32(m.G1012))
	g[1011] = uint64(uint32(m.G1013))
	g[1012] = uint64(uint32(m.G1014))
	g[1013] = uint64(uint32(m.G1015))
	g[1014] = uint64(uint32(m.G1016))
	g[1015] = uint64(uint32(m.G1017))
	g[1016] = uint64(uint32(m.G1018))
	g[1017] = uint64(uint32(m.G1019))
	g[1018] = uint64(uint32(m.G1020))
	g[1019] = uint64(uint32(m.G1021))
	g[1020] = uint64(uint32(m.G1022))
	g[1021] = uint64(uint32(m.G1023))
	g[1022] = uint64(uint32(m.G1024))
	g[1023] = uint64(uint32(m.G1025))
	g[1024] = uint64(uint32(m.G1026))
	g[1025] = uint64(uint32(m.G1027))
	g[1026] = uint64(uint32(m.G1028))
	g[1027] = uint64(uint32(m.G1029))
	g[1028] = uint64(uint32(m.G1030))
	g[1029] = uint64(uint32(m.G1031))
	g[1030] = uint64(uint32(m.G1032))
	g[1031] = uint64(uint32(m.G1033))
	g[1032] = uint64(uint32(m.G1034))
	g[1033] = uint64(uint32(m.G1035))
	g[1034] = uint64(uint32(m.G1036))
	g[1035] = uint64(uint32(m.G1037))
	g[1036] = uint64(uint32(m.G1038))
	g[1037] = uint64(uint32(m.G1039))
	g[1038] = uint64(uint32(m.G1040))
	g[1039] = uint64(uint32(m.G1041))
	g[1040] = uint64(uint32(m.G1042))
	g[1041] = uint64(uint32(m.G1043))
	g[1042] = uint64(uint32(m.G1044))
	g[1043] = uint64(uint32(m.G1045))
	g[1044] = uint64(uint32(m.G1046))
	g[1045] = uint64(uint32(m.G1047))
	g[1046] = uint64(uint32(m.G1048))
	g[1047] = uint64(uint32(m.G1049))
	g[1048] = uint64(uint32(m.G1050))
	g[1049] = uint64(uint32(m.G1051))
	g[1050] = uint64(uint32(m.G1052))
	g[1051] = uint64(uint32(m.G1053))
	g[1052] = uint64(uint32(m.G1054))
	g[1053] = uint64(uint32(m.G1055))
	g[1054] = uint64(uint32(m.G1056))
	g[1055] = uint64(uint32(m.G1057))
	g[1056] = uint64(uint32(m.G1058))
	g[1057] = uint64(uint32(m.G1059))
	g[1058] = uint64(uint32(m.G1060))
	g[1059] = uint64(uint32(m.G1061))
	g[1060] = uint64(uint32(m.G1062))
	g[1061] = uint64(uint32(m.G1063))
	g[1062] = uint64(uint32(m.G1064))
	g[1063] = uint64(uint32(m.G1065))
	g[1064] = uint64(uint32(m.G1066))
	g[1065] = uint64(uint32(m.G1067))
	g[1066] = uint64(uint32(m.G1068))
	g[1067] = uint64(uint32(m.G1069))
	g[1068] = uint64(uint32(m.G1070))
	g[1069] = uint64(uint32(m.G1071))
	g[1070] = uint64(uint32(m.G1072))
	g[1071] = uint64(uint32(m.G1073))
	g[1072] = uint64(uint32(m.G1074))
	g[1073] = uint64(uint32(m.G1075))
	g[1074] = uint64(uint32(m.G1076))
	g[1075] = uint64(uint32(m.G1077))
	g[1076] = uint64(uint32(m.G1078))
	g[1077] = uint64(uint32(m.G1079))
	g[1078] = uint64(uint32(m.G1080))
	g[1079] = uint64(uint32(m.G1081))
	g[1080] = uint64(uint32(m.G1082))
	g[1081] = uint64(uint32(m.G1083))
	g[1082] = uint64(uint32(m.G1084))
	g[1083] = uint64(uint32(m.G1085))
	g[1084] = uint64(uint32(m.G1086))
	g[1085] = uint64(uint32(m.G1087))
	g[1086] = uint64(uint32(m.G1088))
	g[1087] = uint64(uint32(m.G1089))
	g[1088] = uint64(uint32(m.G1090))
	g[1089] = uint64(uint32(m.G1091))
	g[1090] = uint64(uint32(m.G1092))
	g[1091] = uint64(uint32(m.G1093))
	g[1092] = uint64(uint32(m.G1094))
	g[1093] = uint64(uint32(m.G1095))
	g[1094] = uint64(uint32(m.G1096))
	g[1095] = uint64(uint32(m.G1097))
	g[1096] = uint64(uint32(m.G1098))
	g[1097] = uint64(uint32(m.G1099))
	g[1098] = uint64(uint32(m.G1100))
	g[1099] = uint64(uint32(m.G1101))
	g[1100] = uint64(uint32(m.G1102))
	g[1101] = uint64(uint32(m.G1103))
	g[1102] = uint64(uint32(m.G1104))
	g[1103] = uint64(uint32(m.G1105))
	g[1104] = uint64(uint32(m.G1106))
	g[1105] = uint64(uint32(m.G1107))
	g[1106] = uint64(uint32(m.G1108))
	g[1107] = uint64(uint32(m.G1109))
	g[1108] = uint64(uint32(m.G1110))
	g[1109] = uint64(uint32(m.G1111))
	g[1110] = uint64(uint32(m.G1112))
	g[1111] = uint64(uint32(m.G1113))
	g[1112] = uint64(uint32(m.G1114))
	g[1113] = uint64(uint32(m.G1115))
	g[1114] = uint64(uint32(m.G1116))
	g[1115] = uint64(uint32(m.G1117))
	g[1116] = uint64(uint32(m.G1118))
	g[1117] = uint64(uint32(m.G1119))
	g[1118] = uint64(uint32(m.G1120))
	g[1119] = uint64(uint32(m.G1121))
	g[1120] = uint64(uint32(m.G1122))
	g[1121] = uint64(uint32(m.G1123))
	g[1122] = uint64(uint32(m.G1124))
	g[1123] = uint64(uint32(m.G1125))
	g[1124] = uint64(uint32(m.G1126))
	g[1125] = uint64(uint32(m.G1127))
	g[1126] = uint64(uint32(m.G1128))
	g[1127] = uint64(uint32(m.G1129))
	g[1128] = uint64(uint32(m.G1130))
	g[1129] = uint64(uint32(m.G1131))
	g[1130] = uint64(uint32(m.G1132))
	g[1131] = uint64(uint32(m.G1133))
	g[1132] = uint64(uint32(m.G1134))
	g[1133] = uint64(uint32(m.G1135))
	g[1134] = uint64(uint32(m.G1136))
	g[1135] = uint64(uint32(m.G1137))
	g[1136] = uint64(uint32(m.G1138))
	g[1137] = uint64(uint32(m.G1139))
	g[1138] = uint64(uint32(m.G1140))
	g[1139] = uint64(uint32(m.G1141))
	g[1140] = uint64(uint32(m.G1142))
	g[1141] = uint64(uint32(m.G1143))
	g[1142] = uint64(uint32(m.G1144))
	g[1143] = uint64(uint32(m.G1145))
	g[1144] = uint64(uint32(m.G1146))
	g[1145] = uint64(uint32(m.G1147))
	g[1146] = uint64(uint32(m.G1148))
	g[1147] = uint64(uint32(m.G1149))
	g[1148] = uint64(uint32(m.G1150))
	g[1149] = uint64(uint32(m.G1151))
	g[1150] = uint64(uint32(m.G1152))
	g[1151] = uint64(uint32(m.G1153))
	g[1152] = uint64(uint32(m.G1154))
	g[1153] = uint64(uint32(m.G1155))
	g[1154] = uint64(uint32(m.G1156))
	g[1155] = uint64(uint32(m.G1157))
	g[1156] = uint64(uint32(m.G1158))
	g[1157] = uint64(uint32(m.G1159))
	g[1158] = uint64(uint32(m.G1160))
	g[1159] = uint64(uint32(m.G1161))
	g[1160] = uint64(uint32(m.G1162))
	g[1161] = uint64(uint32(m.G1163))
	g[1162] = uint64(uint32(m.G1164))
	g[1163] = uint64(uint32(m.G1165))
	g[1164] = uint64(uint32(m.G1166))
	g[1165] = uint64(uint32(m.G1167))
	g[1166] = uint64(uint32(m.G1168))
	g[1167] = uint64(uint32(m.G1169))
	g[1168] = uint64(uint32(m.G1170))
	g[1169] = uint64(uint32(m.G1171))
	g[1170] = uint64(uint32(m.G1172))
	g[1171] = uint64(uint32(m.G1173))
	g[1172] = uint64(uint32(m.G1174))
	g[1173] = uint64(uint32(m.G1175))
	g[1174] = uint64(uint32(m.G1176))
	g[1175] = uint64(uint32(m.G1177))
	g[1176] = uint64(uint32(m.G1178))
	g[1177] = uint64(uint32(m.G1179))
	g[1178] = uint64(uint32(m.G1180))
	g[1179] = uint64(uint32(m.G1181))
	g[1180] = uint64(uint32(m.G1182))
	g[1181] = uint64(uint32(m.G1183))
	g[1182] = uint64(uint32(m.G1184))
	g[1183] = uint64(uint32(m.G1185))
	g[1184] = uint64(uint32(m.G1186))
	g[1185] = uint64(uint32(m.G1187))
	g[1186] = uint64(uint32(m.G1188))
	g[1187] = uint64(uint32(m.G1189))
	g[1188] = uint64(uint32(m.G1190))
	g[1189] = uint64(uint32(m.G1191))
	g[1190] = uint64(uint32(m.G1192))
	g[1191] = uint64(uint32(m.G1193))
	g[1192] = uint64(uint32(m.G1194))
	g[1193] = uint64(uint32(m.G1195))
	g[1194] = uint64(uint32(m.G1196))
	g[1195] = uint64(uint32(m.G1197))
	g[1196] = uint64(uint32(m.G1198))
	g[1197] = uint64(uint32(m.G1199))
	g[1198] = uint64(uint32(m.G1200))
	g[1199] = uint64(uint32(m.G1201))
	g[1200] = uint64(uint32(m.G1202))
	g[1201] = uint64(uint32(m.G1203))
	g[1202] = uint64(uint32(m.G1204))
	g[1203] = uint64(uint32(m.G1205))
	g[1204] = uint64(uint32(m.G1206))
	g[1205] = uint64(uint32(m.G1207))
	g[1206] = uint64(uint32(m.G1208))
	g[1207] = uint64(uint32(m.G1209))
	g[1208] = uint64(uint32(m.G1210))
	g[1209] = uint64(uint32(m.G1211))
	g[1210] = uint64(uint32(m.G1212))
	g[1211] = uint64(uint32(m.G1213))
	g[1212] = uint64(uint32(m.G1214))
	g[1213] = uint64(uint32(m.G1215))
	g[1214] = uint64(uint32(m.G1216))
	g[1215] = uint64(uint32(m.G1217))
	g[1216] = uint64(uint32(m.G1218))
	g[1217] = uint64(uint32(m.G1219))
	g[1218] = uint64(uint32(m.G1220))
	g[1219] = uint64(uint32(m.G1221))
	g[1220] = uint64(uint32(m.G1222))
	g[1221] = uint64(uint32(m.G1223))
	g[1222] = uint64(uint32(m.G1224))
	g[1223] = uint64(uint32(m.G1225))
	g[1224] = uint64(uint32(m.G1226))
	g[1225] = uint64(uint32(m.G1227))
	g[1226] = uint64(uint32(m.G1228))
	g[1227] = uint64(uint32(m.G1229))
	g[1228] = uint64(uint32(m.G1230))
	g[1229] = uint64(uint32(m.G1231))
	g[1230] = uint64(uint32(m.G1232))
	g[1231] = uint64(uint32(m.G1233))
	g[1232] = uint64(uint32(m.G1234))
	g[1233] = uint64(uint32(m.G1235))
	g[1234] = uint64(uint32(m.G1236))
	g[1235] = uint64(uint32(m.G1237))
	g[1236] = uint64(uint32(m.G1238))
	g[1237] = uint64(uint32(m.G1239))
	g[1238] = uint64(uint32(m.G1240))
	g[1239] = uint64(uint32(m.G1241))
	g[1240] = uint64(uint32(m.G1242))
	g[1241] = uint64(uint32(m.G1243))
	g[1242] = uint64(uint32(m.G1244))
	g[1243] = uint64(uint32(m.G1245))
	g[1244] = uint64(uint32(m.G1246))
	g[1245] = uint64(uint32(m.G1247))
	g[1246] = uint64(uint32(m.G1248))
	g[1247] = uint64(uint32(m.G1249))
	g[1248] = uint64(uint32(m.G1250))
	g[1249] = uint64(uint32(m.G1251))
	g[1250] = uint64(uint32(m.G1252))
	g[1251] = uint64(uint32(m.G1253))
	g[1252] = uint64(uint32(m.G1254))
	g[1253] = uint64(uint32(m.G1255))
	g[1254] = uint64(uint32(m.G1256))
	g[1255] = uint64(uint32(m.G1257))
	g[1256] = uint64(uint32(m.G1258))
	g[1257] = uint64(uint32(m.G1259))
	g[1258] = uint64(uint32(m.G1260))
	g[1259] = uint64(uint32(m.G1261))
	g[1260] = uint64(uint32(m.G1262))
	g[1261] = uint64(uint32(m.G1263))
	g[1262] = uint64(uint32(m.G1264))
	g[1263] = uint64(uint32(m.G1265))
	g[1264] = uint64(uint32(m.G1266))
	g[1265] = uint64(uint32(m.G1267))
	g[1266] = uint64(uint32(m.G1268))
	g[1267] = uint64(uint32(m.G1269))
	g[1268] = uint64(uint32(m.G1270))
	g[1269] = uint64(uint32(m.G1271))
	g[1270] = uint64(uint32(m.G1272))
	g[1271] = uint64(uint32(m.G1273))
	g[1272] = uint64(uint32(m.G1274))
	g[1273] = uint64(uint32(m.G1275))
	g[1274] = uint64(uint32(m.G1276))
	g[1275] = uint64(uint32(m.G1277))
	g[1276] = uint64(uint32(m.G1278))
	g[1277] = uint64(uint32(m.G1279))
	g[1278] = uint64(uint32(m.G1280))
	g[1279] = uint64(uint32(m.G1281))
	g[1280] = uint64(uint32(m.G1282))
	g[1281] = uint64(uint32(m.G1283))
	g[1282] = uint64(uint32(m.G1284))
	g[1283] = uint64(uint32(m.G1285))
	g[1284] = uint64(uint32(m.G1286))
	g[1285] = uint64(uint32(m.G1287))
	g[1286] = uint64(uint32(m.G1288))
	g[1287] = uint64(uint32(m.G1289))
	g[1288] = uint64(uint32(m.G1290))
	g[1289] = uint64(uint32(m.G1291))
	g[1290] = uint64(uint32(m.G1292))
	g[1291] = uint64(uint32(m.G1293))
	g[1292] = uint64(uint32(m.G1294))
	g[1293] = uint64(uint32(m.G1295))
	g[1294] = uint64(uint32(m.G1296))
	g[1295] = uint64(uint32(m.G1297))
	g[1296] = uint64(uint32(m.G1298))
	g[1297] = uint64(uint32(m.G1299))
	g[1298] = uint64(uint32(m.G1300))
	g[1299] = uint64(uint32(m.G1301))
	g[1300] = uint64(uint32(m.G1302))
	g[1301] = uint64(uint32(m.G1303))
	g[1302] = uint64(uint32(m.G1304))
	g[1303] = uint64(uint32(m.G1305))
	g[1304] = uint64(uint32(m.G1306))
	g[1305] = uint64(uint32(m.G1307))
	g[1306] = uint64(uint32(m.G1308))
	g[1307] = uint64(uint32(m.G1309))
	g[1308] = uint64(uint32(m.G1310))
	g[1309] = uint64(uint32(m.G1311))
	g[1310] = uint64(uint32(m.G1312))
	g[1311] = uint64(uint32(m.G1313))
	g[1312] = uint64(uint32(m.G1314))
	g[1313] = uint64(uint32(m.G1315))
	g[1314] = uint64(uint32(m.G1316))
	g[1315] = uint64(uint32(m.G1317))
	g[1316] = uint64(uint32(m.G1318))
	g[1317] = uint64(uint32(m.G1319))
	g[1318] = uint64(uint32(m.G1320))
	g[1319] = uint64(uint32(m.G1321))
	g[1320] = uint64(uint32(m.G1322))
	g[1321] = uint64(uint32(m.G1323))
	g[1322] = uint64(uint32(m.G1324))
	g[1323] = uint64(uint32(m.G1325))
	g[1324] = uint64(uint32(m.G1326))
	g[1325] = uint64(uint32(m.G1327))
	g[1326] = uint64(uint32(m.G1328))
	g[1327] = uint64(uint32(m.G1329))
	g[1328] = uint64(uint32(m.G1330))
	g[1329] = uint64(uint32(m.G1331))
	g[1330] = uint64(uint32(m.G1332))
	g[1331] = uint64(uint32(m.G1333))
	g[1332] = uint64(uint32(m.G1334))
	g[1333] = uint64(uint32(m.G1335))
	g[1334] = uint64(uint32(m.G1336))
	g[1335] = uint64(uint32(m.G1337))
	g[1336] = uint64(uint32(m.G1338))
	g[1337] = uint64(uint32(m.G1339))
	g[1338] = uint64(uint32(m.G1340))
	g[1339] = uint64(uint32(m.G1341))
	g[1340] = uint64(uint32(m.G1342))
	g[1341] = uint64(uint32(m.G1343))
	g[1342] = uint64(uint32(m.G1344))
	g[1343] = uint64(uint32(m.G1345))
	g[1344] = uint64(uint32(m.G1346))
	g[1345] = uint64(uint32(m.G1347))
	g[1346] = uint64(uint32(m.G1348))
	g[1347] = uint64(uint32(m.G1349))
	g[1348] = uint64(uint32(m.G1350))
	g[1349] = uint64(uint32(m.G1351))
	g[1350] = uint64(uint32(m.G1352))
	g[1351] = uint64(uint32(m.G1353))
	g[1352] = uint64(uint32(m.G1354))
	g[1353] = uint64(uint32(m.G1355))
	g[1354] = uint64(uint32(m.G1356))
	g[1355] = uint64(uint32(m.G1357))
	g[1356] = uint64(uint32(m.G1358))
	g[1357] = uint64(uint32(m.G1359))
	g[1358] = uint64(uint32(m.G1360))
	g[1359] = uint64(uint32(m.G1361))
	g[1360] = uint64(uint32(m.G1362))
	g[1361] = uint64(uint32(m.G1363))
	g[1362] = uint64(uint32(m.G1364))
	g[1363] = uint64(uint32(m.G1365))
	g[1364] = uint64(uint32(m.G1366))
	g[1365] = uint64(uint32(m.G1367))
	g[1366] = uint64(uint32(m.G1368))
	g[1367] = uint64(uint32(m.G1369))
	g[1368] = uint64(uint32(m.G1370))
	g[1369] = uint64(uint32(m.G1371))
	g[1370] = uint64(uint32(m.G1372))
	g[1371] = uint64(uint32(m.G1373))
	g[1372] = uint64(uint32(m.G1374))
	g[1373] = uint64(uint32(m.G1375))
	g[1374] = uint64(uint32(m.G1376))
	g[1375] = uint64(uint32(m.G1377))
	g[1376] = uint64(uint32(m.G1378))
	g[1377] = uint64(uint32(m.G1379))
	g[1378] = uint64(uint32(m.G1380))
	g[1379] = uint64(uint32(m.G1381))
	g[1380] = uint64(uint32(m.G1382))
	g[1381] = uint64(uint32(m.G1383))
	g[1382] = uint64(uint32(m.G1384))
	g[1383] = uint64(uint32(m.G1385))
	g[1384] = uint64(uint32(m.G1386))
	g[1385] = uint64(uint32(m.G1387))
	g[1386] = uint64(uint32(m.G1388))
	g[1387] = uint64(uint32(m.G1389))
	g[1388] = uint64(uint32(m.G1390))
	g[1389] = uint64(uint32(m.G1391))
	g[1390] = uint64(uint32(m.G1392))
	g[1391] = uint64(uint32(m.G1393))
	g[1392] = uint64(uint32(m.G1394))
	g[1393] = uint64(uint32(m.G1395))
	g[1394] = uint64(uint32(m.G1396))
	g[1395] = uint64(uint32(m.G1397))
	g[1396] = uint64(uint32(m.G1398))
	g[1397] = uint64(uint32(m.G1399))
	g[1398] = uint64(uint32(m.G1400))
	g[1399] = uint64(uint32(m.G1401))
	g[1400] = uint64(uint32(m.G1402))
	g[1401] = uint64(uint32(m.G1403))
	g[1402] = uint64(uint32(m.G1404))
	g[1403] = uint64(uint32(m.G1405))
	g[1404] = uint64(uint32(m.G1406))
	g[1405] = uint64(uint32(m.G1407))
	g[1406] = uint64(uint32(m.G1408))
	g[1407] = uint64(uint32(m.G1409))
	g[1408] = uint64(uint32(m.G1410))
	g[1409] = uint64(uint32(m.G1411))
	g[1410] = uint64(uint32(m.G1412))
	g[1411] = uint64(uint32(m.G1413))
	g[1412] = uint64(uint32(m.G1414))
	g[1413] = uint64(uint32(m.G1415))
	g[1414] = uint64(uint32(m.G1416))
	g[1415] = uint64(uint32(m.G1417))
	g[1416] = uint64(uint32(m.G1418))
	g[1417] = uint64(uint32(m.G1419))
	g[1418] = uint64(uint32(m.G1420))
	g[1419] = uint64(uint32(m.G1421))
	g[1420] = uint64(uint32(m.G1422))
	g[1421] = uint64(uint32(m.G1423))
	g[1422] = uint64(uint32(m.G1424))
	g[1423] = uint64(uint32(m.G1425))
	g[1424] = uint64(uint32(m.G1426))
	g[1425] = uint64(uint32(m.G1427))
	g[1426] = uint64(uint32(m.G1428))
	g[1427] = uint64(uint32(m.G1429))
	g[1428] = uint64(uint32(m.G1430))
	g[1429] = uint64(uint32(m.G1431))
	g[1430] = uint64(uint32(m.G1432))
	g[1431] = uint64(uint32(m.G1433))
	g[1432] = uint64(uint32(m.G1434))
	g[1433] = uint64(uint32(m.G1435))
	g[1434] = uint64(uint32(m.G1436))
	g[1435] = uint64(uint32(m.G1437))
	g[1436] = uint64(uint32(m.G1438))
	g[1437] = uint64(uint32(m.G1439))
	g[1438] = uint64(uint32(m.G1440))
	g[1439] = uint64(uint32(m.G1441))
	g[1440] = uint64(uint32(m.G1442))
	g[1441] = uint64(uint32(m.G1443))
	g[1442] = uint64(uint32(m.G1444))
	g[1443] = uint64(uint32(m.G1445))
	g[1444] = uint64(uint32(m.G1446))
	g[1445] = uint64(uint32(m.G1447))
	g[1446] = uint64(uint32(m.G1448))
	g[1447] = uint64(uint32(m.G1449))
	g[1448] = uint64(uint32(m.G1450))
	g[1449] = uint64(uint32(m.G1451))
	g[1450] = uint64(uint32(m.G1452))
	g[1451] = uint64(uint32(m.G1453))
	g[1452] = uint64(uint32(m.G1454))
	g[1453] = uint64(uint32(m.G1455))
	g[1454] = uint64(uint32(m.G1456))
	g[1455] = uint64(uint32(m.G1457))
	g[1456] = uint64(uint32(m.G1458))
	g[1457] = uint64(uint32(m.G1459))
	g[1458] = uint64(uint32(m.G1460))
	g[1459] = uint64(uint32(m.G1461))
	g[1460] = uint64(uint32(m.G1462))
	g[1461] = uint64(uint32(m.G1463))
	g[1462] = uint64(uint32(m.G1464))
	g[1463] = uint64(uint32(m.G1465))
	g[1464] = uint64(uint32(m.G1466))
	g[1465] = uint64(uint32(m.G1467))
	g[1466] = uint64(uint32(m.G1468))
	g[1467] = uint64(uint32(m.G1469))
	g[1468] = uint64(uint32(m.G1470))
	g[1469] = uint64(uint32(m.G1471))
	g[1470] = uint64(uint32(m.G1472))
	g[1471] = uint64(uint32(m.G1473))
	g[1472] = uint64(uint32(m.G1474))
	g[1473] = uint64(uint32(m.G1475))
	g[1474] = uint64(uint32(m.G1476))
	g[1475] = uint64(uint32(m.G1477))
	g[1476] = uint64(uint32(m.G1478))
	g[1477] = uint64(uint32(m.G1479))
	g[1478] = uint64(uint32(m.G1480))
	g[1479] = uint64(uint32(m.G1481))
	g[1480] = uint64(uint32(m.G1482))
	g[1481] = uint64(uint32(m.G1483))
	g[1482] = uint64(uint32(m.G1484))
	g[1483] = uint64(uint32(m.G1485))
	g[1484] = uint64(uint32(m.G1486))
	g[1485] = uint64(uint32(m.G1487))
	g[1486] = uint64(uint32(m.G1488))
	g[1487] = uint64(uint32(m.G1489))
	g[1488] = uint64(uint32(m.G1490))
	g[1489] = uint64(uint32(m.G1491))
	g[1490] = uint64(uint32(m.G1492))
	g[1491] = uint64(uint32(m.G1493))
	g[1492] = uint64(uint32(m.G1494))
	g[1493] = uint64(uint32(m.G1495))
	g[1494] = uint64(uint32(m.G1496))
	g[1495] = uint64(uint32(m.G1497))
	g[1496] = uint64(uint32(m.G1498))
	g[1497] = uint64(uint32(m.G1499))
	g[1498] = uint64(uint32(m.G1500))
	g[1499] = uint64(uint32(m.G1501))
	g[1500] = uint64(uint32(m.G1502))
	g[1501] = uint64(uint32(m.G1503))
	g[1502] = uint64(uint32(m.G1504))
	g[1503] = uint64(uint32(m.G1505))
	g[1504] = uint64(uint32(m.G1506))
	g[1505] = uint64(uint32(m.G1507))
	g[1506] = uint64(uint32(m.G1508))
	g[1507] = uint64(uint32(m.G1509))
	g[1508] = uint64(uint32(m.G1510))
	g[1509] = uint64(uint32(m.G1511))
	g[1510] = uint64(uint32(m.G1512))
	g[1511] = uint64(uint32(m.G1513))
	g[1512] = uint64(uint32(m.G1514))
	g[1513] = uint64(uint32(m.G1515))
	g[1514] = uint64(uint32(m.G1516))
	g[1515] = uint64(uint32(m.G1517))
	g[1516] = uint64(uint32(m.G1518))
	g[1517] = uint64(uint32(m.G1519))
	g[1518] = uint64(uint32(m.G1520))
	g[1519] = uint64(uint32(m.G1521))
	g[1520] = uint64(uint32(m.G1522))
	g[1521] = uint64(uint32(m.G1523))
	g[1522] = uint64(uint32(m.G1524))
	g[1523] = uint64(uint32(m.G1525))
	g[1524] = uint64(uint32(m.G1526))
	g[1525] = uint64(uint32(m.G1527))
	g[1526] = uint64(uint32(m.G1528))
	g[1527] = uint64(uint32(m.G1529))
	g[1528] = uint64(uint32(m.G1530))
	g[1529] = uint64(uint32(m.G1531))
	g[1530] = uint64(uint32(m.G1532))
	g[1531] = uint64(uint32(m.G1533))
	g[1532] = uint64(uint32(m.G1534))
	g[1533] = uint64(uint32(m.G1535))
	g[1534] = uint64(uint32(m.G1536))
	g[1535] = uint64(uint32(m.G1537))
	g[1536] = uint64(uint32(m.G1538))
	g[1537] = uint64(uint32(m.G1539))
	g[1538] = uint64(uint32(m.G1540))
	g[1539] = uint64(uint32(m.G1541))
	g[1540] = uint64(uint32(m.G1542))
	g[1541] = uint64(uint32(m.G1543))
	g[1542] = uint64(uint32(m.G1544))
	g[1543] = uint64(uint32(m.G1545))
	g[1544] = uint64(uint32(m.G1546))
	g[1545] = uint64(uint32(m.G1547))
	g[1546] = uint64(uint32(m.G1548))
	g[1547] = uint64(uint32(m.G1549))
	g[1548] = uint64(uint32(m.G1550))
	g[1549] = uint64(uint32(m.G1551))
	g[1550] = uint64(uint32(m.G1552))
	g[1551] = uint64(uint32(m.G1553))
	g[1552] = uint64(uint32(m.G1554))
	g[1553] = uint64(uint32(m.G1555))
	g[1554] = uint64(uint32(m.G1556))
	g[1555] = uint64(uint32(m.G1557))
	g[1556] = uint64(uint32(m.G1558))
	g[1557] = uint64(uint32(m.G1559))
	g[1558] = uint64(uint32(m.G1560))
	g[1559] = uint64(uint32(m.G1561))
	g[1560] = uint64(uint32(m.G1562))
	g[1561] = uint64(uint32(m.G1563))
	g[1562] = uint64(uint32(m.G1564))
	g[1563] = uint64(uint32(m.G1565))
	g[1564] = uint64(uint32(m.G1566))
	g[1565] = uint64(uint32(m.G1567))
	g[1566] = uint64(uint32(m.G1568))
	g[1567] = uint64(uint32(m.G1569))
	g[1568] = uint64(uint32(m.G1570))
	g[1569] = uint64(uint32(m.G1571))
	g[1570] = uint64(uint32(m.G1572))
	g[1571] = uint64(uint32(m.G1573))
	g[1572] = uint64(uint32(m.G1574))
	g[1573] = uint64(uint32(m.G1575))
	g[1574] = uint64(uint32(m.G1576))
	g[1575] = uint64(uint32(m.G1577))
	g[1576] = uint64(uint32(m.G1578))
	g[1577] = uint64(uint32(m.G1579))
	g[1578] = uint64(uint32(m.G1580))
	g[1579] = uint64(uint32(m.G1581))
	g[1580] = uint64(uint32(m.G1582))
	g[1581] = uint64(uint32(m.G1583))
	g[1582] = uint64(uint32(m.G1584))
	g[1583] = uint64(uint32(m.G1585))
	g[1584] = uint64(uint32(m.G1586))
	g[1585] = uint64(uint32(m.G1587))
	g[1586] = uint64(uint32(m.G1588))
	g[1587] = uint64(uint32(m.G1589))
	g[1588] = uint64(uint32(m.G1590))
	g[1589] = uint64(uint32(m.G1591))
	g[1590] = uint64(uint32(m.G1592))
	g[1591] = uint64(uint32(m.G1593))
	g[1592] = uint64(uint32(m.G1594))
	g[1593] = uint64(uint32(m.G1595))
	g[1594] = uint64(uint32(m.G1596))
	g[1595] = uint64(uint32(m.G1597))
	g[1596] = uint64(uint32(m.G1598))
	g[1597] = uint64(uint32(m.G1599))
	g[1598] = uint64(uint32(m.G1600))
	g[1599] = uint64(uint32(m.G1601))
	g[1600] = uint64(uint32(m.G1602))
	g[1601] = uint64(uint32(m.G1603))
	g[1602] = uint64(uint32(m.G1604))
	g[1603] = uint64(uint32(m.G1605))
	g[1604] = uint64(uint32(m.G1606))
	g[1605] = uint64(uint32(m.G1607))
	g[1606] = uint64(uint32(m.G1608))
	g[1607] = uint64(uint32(m.G1609))
	g[1608] = uint64(uint32(m.G1610))
	g[1609] = uint64(uint32(m.G1611))
	g[1610] = uint64(uint32(m.G1612))
	g[1611] = uint64(uint32(m.G1613))
	g[1612] = uint64(uint32(m.G1614))
	g[1613] = uint64(uint32(m.G1615))
	g[1614] = uint64(uint32(m.G1616))
	g[1615] = uint64(uint32(m.G1617))
	g[1616] = uint64(uint32(m.G1618))
	g[1617] = uint64(uint32(m.G1619))
	g[1618] = uint64(uint32(m.G1620))
	g[1619] = uint64(uint32(m.G1621))
	g[1620] = uint64(uint32(m.G1622))
	g[1621] = uint64(uint32(m.G1623))
	g[1622] = uint64(uint32(m.G1624))
	g[1623] = uint64(uint32(m.G1625))
	g[1624] = uint64(uint32(m.G1626))
	g[1625] = uint64(uint32(m.G1627))
	g[1626] = uint64(uint32(m.G1628))
	g[1627] = uint64(uint32(m.G1629))
	g[1628] = uint64(uint32(m.G1630))
	g[1629] = uint64(uint32(m.G1631))
	g[1630] = uint64(uint32(m.G1632))
	g[1631] = uint64(uint32(m.G1633))
	g[1632] = uint64(uint32(m.G1634))
	g[1633] = uint64(uint32(m.G1635))
	g[1634] = uint64(uint32(m.G1636))
	g[1635] = uint64(uint32(m.G1637))
	g[1636] = uint64(uint32(m.G1638))
	g[1637] = uint64(uint32(m.G1639))
	g[1638] = uint64(uint32(m.G1640))
	g[1639] = uint64(uint32(m.G1641))
	g[1640] = uint64(uint32(m.G1642))
	g[1641] = uint64(uint32(m.G1643))
	g[1642] = uint64(uint32(m.G1644))
	g[1643] = uint64(uint32(m.G1645))
	g[1644] = uint64(uint32(m.G1646))
	g[1645] = uint64(uint32(m.G1647))
	g[1646] = uint64(uint32(m.G1648))
	g[1647] = uint64(uint32(m.G1649))
	g[1648] = uint64(uint32(m.G1650))
	g[1649] = uint64(uint32(m.G1651))
	g[1650] = uint64(uint32(m.G1652))
	g[1651] = uint64(uint32(m.G1653))
	g[1652] = uint64(uint32(m.G1654))
	g[1653] = uint64(uint32(m.G1655))
	g[1654] = uint64(uint32(m.G1656))
	g[1655] = uint64(uint32(m.G1657))
	g[1656] = uint64(uint32(m.G1658))
	g[1657] = uint64(uint32(m.G1659))
	g[1658] = uint64(uint32(m.G1660))
	g[1659] = uint64(uint32(m.G1661))
	g[1660] = uint64(uint32(m.G1662))
	g[1661] = uint64(uint32(m.G1663))
	g[1662] = uint64(uint32(m.G1664))
	g[1663] = uint64(uint32(m.G1665))
	g[1664] = uint64(uint32(m.G1666))
	g[1665] = uint64(uint32(m.G1667))
	g[1666] = uint64(uint32(m.G1668))
	g[1667] = uint64(uint32(m.G1669))
	g[1668] = uint64(uint32(m.G1670))
	g[1669] = uint64(uint32(m.G1671))
	g[1670] = uint64(uint32(m.G1672))
	g[1671] = uint64(uint32(m.G1673))
	g[1672] = uint64(uint32(m.G1674))
	g[1673] = uint64(uint32(m.G1675))
	g[1674] = uint64(uint32(m.G1676))
	g[1675] = uint64(uint32(m.G1677))
	g[1676] = uint64(uint32(m.G1678))
	g[1677] = uint64(uint32(m.G1679))
	g[1678] = uint64(uint32(m.G1680))
	g[1679] = uint64(uint32(m.G1681))
	g[1680] = uint64(uint32(m.G1682))
	g[1681] = uint64(uint32(m.G1683))
	g[1682] = uint64(uint32(m.G1684))
	g[1683] = uint64(uint32(m.G1685))
	g[1684] = uint64(uint32(m.G1686))
	g[1685] = uint64(uint32(m.G1687))
	g[1686] = uint64(uint32(m.G1688))
	g[1687] = uint64(uint32(m.G1689))
	g[1688] = uint64(uint32(m.G1690))
	g[1689] = uint64(uint32(m.G1691))
	g[1690] = uint64(uint32(m.G1692))
	g[1691] = uint64(uint32(m.G1693))
	g[1692] = uint64(uint32(m.G1694))
	g[1693] = uint64(uint32(m.G1695))
	g[1694] = uint64(uint32(m.G1696))
	g[1695] = uint64(uint32(m.G1697))
	g[1696] = uint64(uint32(m.G1698))
	g[1697] = uint64(uint32(m.G1699))
	g[1698] = uint64(uint32(m.G1700))
	g[1699] = uint64(uint32(m.G1701))
	g[1700] = uint64(uint32(m.G1702))
	g[1701] = uint64(uint32(m.G1703))
	g[1702] = uint64(uint32(m.G1704))
	g[1703] = uint64(uint32(m.G1705))
	g[1704] = uint64(uint32(m.G1706))
	g[1705] = uint64(uint32(m.G1707))
	g[1706] = uint64(uint32(m.G1708))
	g[1707] = uint64(uint32(m.G1709))
	g[1708] = uint64(uint32(m.G1710))
	g[1709] = uint64(uint32(m.G1711))
	g[1710] = uint64(uint32(m.G1712))
	g[1711] = uint64(uint32(m.G1713))
	g[1712] = uint64(uint32(m.G1714))
	g[1713] = uint64(uint32(m.G1715))
	g[1714] = uint64(uint32(m.G1716))
	g[1715] = uint64(uint32(m.G1717))
	g[1716] = uint64(uint32(m.G1718))
	g[1717] = uint64(uint32(m.G1719))
	g[1718] = uint64(uint32(m.G1720))
	g[1719] = uint64(uint32(m.G1721))
	g[1720] = uint64(uint32(m.G1722))
	g[1721] = uint64(uint32(m.G1723))
	g[1722] = uint64(uint32(m.G1724))
	g[1723] = uint64(uint32(m.G1725))
	g[1724] = uint64(uint32(m.G1726))
	g[1725] = uint64(uint32(m.G1727))
	g[1726] = uint64(uint32(m.G1728))
	g[1727] = uint64(uint32(m.G1729))
	g[1728] = uint64(uint32(m.G1730))
	g[1729] = uint64(uint32(m.G1731))
	g[1730] = uint64(uint32(m.G1732))
	g[1731] = uint64(uint32(m.G1733))
	g[1732] = uint64(uint32(m.G1734))
	g[1733] = uint64(uint32(m.G1735))
	g[1734] = uint64(uint32(m.G1736))
	g[1735] = uint64(uint32(m.G1737))
	g[1736] = uint64(uint32(m.G1738))
	g[1737] = uint64(uint32(m.G1739))
	g[1738] = uint64(uint32(m.G1740))
	g[1739] = uint64(uint32(m.G1741))
	g[1740] = uint64(uint32(m.G1742))
	g[1741] = uint64(uint32(m.G1743))
	g[1742] = uint64(uint32(m.G1744))
	g[1743] = uint64(uint32(m.G1745))
	g[1744] = uint64(uint32(m.G1746))
	g[1745] = uint64(uint32(m.G1747))
	g[1746] = uint64(uint32(m.G1748))
	g[1747] = uint64(uint32(m.G1749))
	g[1748] = uint64(uint32(m.G1750))
	g[1749] = uint64(uint32(m.G1751))
	g[1750] = uint64(uint32(m.G1752))
	g[1751] = uint64(uint32(m.G1753))
	g[1752] = uint64(uint32(m.G1754))
	g[1753] = uint64(uint32(m.G1755))
	g[1754] = uint64(uint32(m.G1756))
	g[1755] = uint64(uint32(m.G1757))
	g[1756] = uint64(uint32(m.G1758))
	g[1757] = uint64(uint32(m.G1759))
	g[1758] = uint64(uint32(m.G1760))
	g[1759] = uint64(uint32(m.G1761))
	g[1760] = uint64(uint32(m.G1762))
	g[1761] = uint64(uint32(m.G1763))
	g[1762] = uint64(uint32(m.G1764))
	g[1763] = uint64(uint32(m.G1765))
	g[1764] = uint64(uint32(m.G1766))
	g[1765] = uint64(uint32(m.G1767))
	g[1766] = uint64(uint32(m.G1768))
	g[1767] = uint64(uint32(m.G1769))
	g[1768] = uint64(uint32(m.G1770))
	g[1769] = uint64(uint32(m.G1771))
	g[1770] = uint64(uint32(m.G1772))
	g[1771] = uint64(uint32(m.G1773))
	g[1772] = uint64(uint32(m.G1774))
	g[1773] = uint64(uint32(m.G1775))
	g[1774] = uint64(uint32(m.G1776))
	g[1775] = uint64(uint32(m.G1777))
	g[1776] = uint64(uint32(m.G1778))
	g[1777] = uint64(uint32(m.G1779))
	g[1778] = uint64(uint32(m.G1780))
	g[1779] = uint64(uint32(m.G1781))
	g[1780] = uint64(uint32(m.G1782))
	g[1781] = uint64(uint32(m.G1783))
	g[1782] = uint64(uint32(m.G1784))
	g[1783] = uint64(uint32(m.G1785))
	g[1784] = uint64(uint32(m.G1786))
	g[1785] = uint64(uint32(m.G1787))
	g[1786] = uint64(uint32(m.G1788))
	g[1787] = uint64(uint32(m.G1789))
	g[1788] = uint64(uint32(m.G1790))
	g[1789] = uint64(uint32(m.G1791))
	g[1790] = uint64(uint32(m.G1792))
	g[1791] = uint64(uint32(m.G1793))
	g[1792] = uint64(uint32(m.G1794))
	g[1793] = uint64(uint32(m.G1795))
	g[1794] = uint64(uint32(m.G1796))
	g[1795] = uint64(uint32(m.G1797))
	g[1796] = uint64(uint32(m.G1798))
	g[1797] = uint64(uint32(m.G1799))
	g[1798] = uint64(uint32(m.G1800))
	g[1799] = uint64(uint32(m.G1801))
	g[1800] = uint64(uint32(m.G1802))
	g[1801] = uint64(uint32(m.G1803))
	g[1802] = uint64(uint32(m.G1804))
	g[1803] = uint64(uint32(m.G1805))
	g[1804] = uint64(uint32(m.G1806))
	g[1805] = uint64(uint32(m.G1807))
	g[1806] = uint64(uint32(m.G1808))
	g[1807] = uint64(uint32(m.G1809))
	g[1808] = uint64(uint32(m.G1810))
	g[1809] = uint64(uint32(m.G1811))
	g[1810] = uint64(uint32(m.G1812))
	g[1811] = uint64(uint32(m.G1813))
	g[1812] = uint64(uint32(m.G1814))
	g[1813] = uint64(uint32(m.G1815))
	g[1814] = uint64(uint32(m.G1816))
	g[1815] = uint64(uint32(m.G1817))
	g[1816] = uint64(uint32(m.G1818))
	g[1817] = uint64(uint32(m.G1819))
	g[1818] = uint64(uint32(m.G1820))
	g[1819] = uint64(uint32(m.G1821))
	g[1820] = uint64(uint32(m.G1822))
	g[1821] = uint64(uint32(m.G1823))
	g[1822] = uint64(uint32(m.G1824))
	g[1823] = uint64(uint32(m.G1825))
	g[1824] = uint64(uint32(m.G1826))
	g[1825] = uint64(uint32(m.G1827))
	g[1826] = uint64(uint32(m.G1828))
	g[1827] = uint64(uint32(m.G1829))
	g[1828] = uint64(uint32(m.G1830))
	g[1829] = uint64(uint32(m.G1831))
	g[1830] = uint64(uint32(m.G1832))
	g[1831] = uint64(uint32(m.G1833))
	g[1832] = uint64(uint32(m.G1834))
	g[1833] = uint64(uint32(m.G1835))
	g[1834] = uint64(uint32(m.G1836))
	g[1835] = uint64(uint32(m.G1837))
	g[1836] = uint64(uint32(m.G1838))
	g[1837] = uint64(uint32(m.G1839))
	g[1838] = uint64(uint32(m.G1840))
	g[1839] = uint64(uint32(m.G1841))
	g[1840] = uint64(uint32(m.G1842))
	g[1841] = uint64(uint32(m.G1843))
	g[1842] = uint64(uint32(m.G1844))
	g[1843] = uint64(uint32(m.G1845))
	g[1844] = uint64(uint32(m.G1846))
	g[1845] = uint64(uint32(m.G1847))
	g[1846] = uint64(uint32(m.G1848))
	g[1847] = uint64(uint32(m.G1849))
	g[1848] = uint64(uint32(m.G1850))
	g[1849] = uint64(uint32(m.G1851))
	g[1850] = uint64(uint32(m.G1852))
	g[1851] = uint64(uint32(m.G1853))
	g[1852] = uint64(uint32(m.G1854))
	g[1853] = uint64(uint32(m.G1855))
	g[1854] = uint64(uint32(m.G1856))
	g[1855] = uint64(uint32(m.G1857))
	g[1856] = uint64(uint32(m.G1858))
	g[1857] = uint64(uint32(m.G1859))
	g[1858] = uint64(uint32(m.G1860))
	g[1859] = uint64(uint32(m.G1861))
	g[1860] = uint64(uint32(m.G1862))
	g[1861] = uint64(uint32(m.G1863))
	g[1862] = uint64(uint32(m.G1864))
	g[1863] = uint64(uint32(m.G1865))
	g[1864] = uint64(uint32(m.G1866))
	g[1865] = uint64(uint32(m.G1867))
	g[1866] = uint64(uint32(m.G1868))
	g[1867] = uint64(uint32(m.G1869))
	g[1868] = uint64(uint32(m.G1870))
	g[1869] = uint64(uint32(m.G1871))
	g[1870] = uint64(uint32(m.G1872))
	g[1871] = uint64(uint32(m.G1873))
	g[1872] = uint64(uint32(m.G1874))
	g[1873] = uint64(uint32(m.G1875))
	g[1874] = uint64(uint32(m.G1876))
	g[1875] = uint64(uint32(m.G1877))
	g[1876] = uint64(uint32(m.G1878))
	g[1877] = uint64(uint32(m.G1879))
	g[1878] = uint64(uint32(m.G1880))
	g[1879] = uint64(uint32(m.G1881))
	g[1880] = uint64(uint32(m.G1882))
	return g
}

// RestoreGlobals puts a snapshot's globals back. A snapshot from a different module (or a
// different build of the same one) has a different global count; rather than
// index out of bounds, take what fits and leave the rest at their declared
// initializers.
func RestoreGlobals(m *Module, g []uint64) {
	if len(g) != 1881 {
		return
	}
	m.G1 = int32(uint32(g[0]))
	m.G3 = int32(uint32(g[1]))
	m.G4 = int32(uint32(g[2]))
	m.G5 = int32(uint32(g[3]))
	m.G6 = int32(uint32(g[4]))
	m.G7 = int32(uint32(g[5]))
	m.G8 = int32(uint32(g[6]))
	m.G9 = int32(uint32(g[7]))
	m.G10 = int32(uint32(g[8]))
	m.G11 = int32(uint32(g[9]))
	m.G12 = int32(uint32(g[10]))
	m.G13 = int32(uint32(g[11]))
	m.G14 = int32(uint32(g[12]))
	m.G15 = int32(uint32(g[13]))
	m.G16 = int32(uint32(g[14]))
	m.G17 = int32(uint32(g[15]))
	m.G18 = int32(uint32(g[16]))
	m.G19 = int32(uint32(g[17]))
	m.G20 = int32(uint32(g[18]))
	m.G21 = int32(uint32(g[19]))
	m.G22 = int32(uint32(g[20]))
	m.G23 = int32(uint32(g[21]))
	m.G24 = int32(uint32(g[22]))
	m.G25 = int32(uint32(g[23]))
	m.G26 = int32(uint32(g[24]))
	m.G27 = int32(uint32(g[25]))
	m.G28 = int32(uint32(g[26]))
	m.G29 = int32(uint32(g[27]))
	m.G30 = int32(uint32(g[28]))
	m.G31 = int32(uint32(g[29]))
	m.G32 = int32(uint32(g[30]))
	m.G33 = int32(uint32(g[31]))
	m.G34 = int32(uint32(g[32]))
	m.G35 = int32(uint32(g[33]))
	m.G36 = int32(uint32(g[34]))
	m.G37 = int32(uint32(g[35]))
	m.G38 = int32(uint32(g[36]))
	m.G39 = int32(uint32(g[37]))
	m.G40 = int32(uint32(g[38]))
	m.G41 = int32(uint32(g[39]))
	m.G42 = int32(uint32(g[40]))
	m.G43 = int32(uint32(g[41]))
	m.G44 = int32(uint32(g[42]))
	m.G45 = int32(uint32(g[43]))
	m.G46 = int32(uint32(g[44]))
	m.G47 = int32(uint32(g[45]))
	m.G48 = int32(uint32(g[46]))
	m.G49 = int32(uint32(g[47]))
	m.G50 = int32(uint32(g[48]))
	m.G51 = int32(uint32(g[49]))
	m.G52 = int32(uint32(g[50]))
	m.G53 = int32(uint32(g[51]))
	m.G54 = int32(uint32(g[52]))
	m.G55 = int32(uint32(g[53]))
	m.G56 = int32(uint32(g[54]))
	m.G57 = int32(uint32(g[55]))
	m.G58 = int32(uint32(g[56]))
	m.G59 = int32(uint32(g[57]))
	m.G60 = int32(uint32(g[58]))
	m.G61 = int32(uint32(g[59]))
	m.G62 = int32(uint32(g[60]))
	m.G63 = int32(uint32(g[61]))
	m.G64 = int32(uint32(g[62]))
	m.G65 = int32(uint32(g[63]))
	m.G66 = int32(uint32(g[64]))
	m.G67 = int32(uint32(g[65]))
	m.G68 = int32(uint32(g[66]))
	m.G69 = int32(uint32(g[67]))
	m.G70 = int32(uint32(g[68]))
	m.G71 = int32(uint32(g[69]))
	m.G72 = int32(uint32(g[70]))
	m.G73 = int32(uint32(g[71]))
	m.G74 = int32(uint32(g[72]))
	m.G75 = int32(uint32(g[73]))
	m.G76 = int32(uint32(g[74]))
	m.G77 = int32(uint32(g[75]))
	m.G78 = int32(uint32(g[76]))
	m.G79 = int32(uint32(g[77]))
	m.G80 = int32(uint32(g[78]))
	m.G81 = int32(uint32(g[79]))
	m.G82 = int32(uint32(g[80]))
	m.G83 = int32(uint32(g[81]))
	m.G84 = int32(uint32(g[82]))
	m.G85 = int32(uint32(g[83]))
	m.G86 = int32(uint32(g[84]))
	m.G87 = int32(uint32(g[85]))
	m.G88 = int32(uint32(g[86]))
	m.G89 = int32(uint32(g[87]))
	m.G90 = int32(uint32(g[88]))
	m.G91 = int32(uint32(g[89]))
	m.G92 = int32(uint32(g[90]))
	m.G93 = int32(uint32(g[91]))
	m.G94 = int32(uint32(g[92]))
	m.G95 = int32(uint32(g[93]))
	m.G96 = int32(uint32(g[94]))
	m.G97 = int32(uint32(g[95]))
	m.G98 = int32(uint32(g[96]))
	m.G99 = int32(uint32(g[97]))
	m.G100 = int32(uint32(g[98]))
	m.G101 = int32(uint32(g[99]))
	m.G102 = int32(uint32(g[100]))
	m.G103 = int32(uint32(g[101]))
	m.G104 = int32(uint32(g[102]))
	m.G105 = int32(uint32(g[103]))
	m.G106 = int32(uint32(g[104]))
	m.G107 = int32(uint32(g[105]))
	m.G108 = int32(uint32(g[106]))
	m.G109 = int32(uint32(g[107]))
	m.G110 = int32(uint32(g[108]))
	m.G111 = int32(uint32(g[109]))
	m.G112 = int32(uint32(g[110]))
	m.G113 = int32(uint32(g[111]))
	m.G114 = int32(uint32(g[112]))
	m.G115 = int32(uint32(g[113]))
	m.G116 = int32(uint32(g[114]))
	m.G117 = int32(uint32(g[115]))
	m.G118 = int32(uint32(g[116]))
	m.G119 = int32(uint32(g[117]))
	m.G120 = int32(uint32(g[118]))
	m.G121 = int32(uint32(g[119]))
	m.G122 = int32(uint32(g[120]))
	m.G123 = int32(uint32(g[121]))
	m.G124 = int32(uint32(g[122]))
	m.G125 = int32(uint32(g[123]))
	m.G126 = int32(uint32(g[124]))
	m.G127 = int32(uint32(g[125]))
	m.G128 = int32(uint32(g[126]))
	m.G129 = int32(uint32(g[127]))
	m.G130 = int32(uint32(g[128]))
	m.G131 = int32(uint32(g[129]))
	m.G132 = int32(uint32(g[130]))
	m.G133 = int32(uint32(g[131]))
	m.G134 = int32(uint32(g[132]))
	m.G135 = int32(uint32(g[133]))
	m.G136 = int32(uint32(g[134]))
	m.G137 = int32(uint32(g[135]))
	m.G138 = int32(uint32(g[136]))
	m.G139 = int32(uint32(g[137]))
	m.G140 = int32(uint32(g[138]))
	m.G141 = int32(uint32(g[139]))
	m.G142 = int32(uint32(g[140]))
	m.G143 = int32(uint32(g[141]))
	m.G144 = int32(uint32(g[142]))
	m.G145 = int32(uint32(g[143]))
	m.G146 = int32(uint32(g[144]))
	m.G147 = int32(uint32(g[145]))
	m.G148 = int32(uint32(g[146]))
	m.G149 = int32(uint32(g[147]))
	m.G150 = int32(uint32(g[148]))
	m.G151 = int32(uint32(g[149]))
	m.G152 = int32(uint32(g[150]))
	m.G153 = int32(uint32(g[151]))
	m.G154 = int32(uint32(g[152]))
	m.G155 = int32(uint32(g[153]))
	m.G156 = int32(uint32(g[154]))
	m.G157 = int32(uint32(g[155]))
	m.G158 = int32(uint32(g[156]))
	m.G159 = int32(uint32(g[157]))
	m.G160 = int32(uint32(g[158]))
	m.G161 = int32(uint32(g[159]))
	m.G162 = int32(uint32(g[160]))
	m.G163 = int32(uint32(g[161]))
	m.G164 = int32(uint32(g[162]))
	m.G165 = int32(uint32(g[163]))
	m.G166 = int32(uint32(g[164]))
	m.G167 = int32(uint32(g[165]))
	m.G168 = int32(uint32(g[166]))
	m.G169 = int32(uint32(g[167]))
	m.G170 = int32(uint32(g[168]))
	m.G171 = int32(uint32(g[169]))
	m.G172 = int32(uint32(g[170]))
	m.G173 = int32(uint32(g[171]))
	m.G174 = int32(uint32(g[172]))
	m.G175 = int32(uint32(g[173]))
	m.G176 = int32(uint32(g[174]))
	m.G177 = int32(uint32(g[175]))
	m.G178 = int32(uint32(g[176]))
	m.G179 = int32(uint32(g[177]))
	m.G180 = int32(uint32(g[178]))
	m.G181 = int32(uint32(g[179]))
	m.G182 = int32(uint32(g[180]))
	m.G183 = int32(uint32(g[181]))
	m.G184 = int32(uint32(g[182]))
	m.G185 = int32(uint32(g[183]))
	m.G186 = int32(uint32(g[184]))
	m.G187 = int32(uint32(g[185]))
	m.G188 = int32(uint32(g[186]))
	m.G189 = int32(uint32(g[187]))
	m.G190 = int32(uint32(g[188]))
	m.G191 = int32(uint32(g[189]))
	m.G192 = int32(uint32(g[190]))
	m.G193 = int32(uint32(g[191]))
	m.G194 = int32(uint32(g[192]))
	m.G195 = int32(uint32(g[193]))
	m.G196 = int32(uint32(g[194]))
	m.G197 = int32(uint32(g[195]))
	m.G198 = int32(uint32(g[196]))
	m.G199 = int32(uint32(g[197]))
	m.G200 = int32(uint32(g[198]))
	m.G201 = int32(uint32(g[199]))
	m.G202 = int32(uint32(g[200]))
	m.G203 = int32(uint32(g[201]))
	m.G204 = int32(uint32(g[202]))
	m.G205 = int32(uint32(g[203]))
	m.G206 = int32(uint32(g[204]))
	m.G207 = int32(uint32(g[205]))
	m.G208 = int32(uint32(g[206]))
	m.G209 = int32(uint32(g[207]))
	m.G210 = int32(uint32(g[208]))
	m.G211 = int32(uint32(g[209]))
	m.G212 = int32(uint32(g[210]))
	m.G213 = int32(uint32(g[211]))
	m.G214 = int32(uint32(g[212]))
	m.G215 = int32(uint32(g[213]))
	m.G216 = int32(uint32(g[214]))
	m.G217 = int32(uint32(g[215]))
	m.G218 = int32(uint32(g[216]))
	m.G219 = int32(uint32(g[217]))
	m.G220 = int32(uint32(g[218]))
	m.G221 = int32(uint32(g[219]))
	m.G222 = int32(uint32(g[220]))
	m.G223 = int32(uint32(g[221]))
	m.G224 = int32(uint32(g[222]))
	m.G225 = int32(uint32(g[223]))
	m.G226 = int32(uint32(g[224]))
	m.G227 = int32(uint32(g[225]))
	m.G228 = int32(uint32(g[226]))
	m.G229 = int32(uint32(g[227]))
	m.G230 = int32(uint32(g[228]))
	m.G231 = int32(uint32(g[229]))
	m.G232 = int32(uint32(g[230]))
	m.G233 = int32(uint32(g[231]))
	m.G234 = int32(uint32(g[232]))
	m.G235 = int32(uint32(g[233]))
	m.G236 = int32(uint32(g[234]))
	m.G237 = int32(uint32(g[235]))
	m.G238 = int32(uint32(g[236]))
	m.G239 = int32(uint32(g[237]))
	m.G240 = int32(uint32(g[238]))
	m.G241 = int32(uint32(g[239]))
	m.G242 = int32(uint32(g[240]))
	m.G243 = int32(uint32(g[241]))
	m.G244 = int32(uint32(g[242]))
	m.G245 = int32(uint32(g[243]))
	m.G246 = int32(uint32(g[244]))
	m.G247 = int32(uint32(g[245]))
	m.G248 = int32(uint32(g[246]))
	m.G249 = int32(uint32(g[247]))
	m.G250 = int32(uint32(g[248]))
	m.G251 = int32(uint32(g[249]))
	m.G252 = int32(uint32(g[250]))
	m.G253 = int32(uint32(g[251]))
	m.G254 = int32(uint32(g[252]))
	m.G255 = int32(uint32(g[253]))
	m.G256 = int32(uint32(g[254]))
	m.G257 = int32(uint32(g[255]))
	m.G258 = int32(uint32(g[256]))
	m.G259 = int32(uint32(g[257]))
	m.G260 = int32(uint32(g[258]))
	m.G261 = int32(uint32(g[259]))
	m.G262 = int32(uint32(g[260]))
	m.G263 = int32(uint32(g[261]))
	m.G264 = int32(uint32(g[262]))
	m.G265 = int32(uint32(g[263]))
	m.G266 = int32(uint32(g[264]))
	m.G267 = int32(uint32(g[265]))
	m.G268 = int32(uint32(g[266]))
	m.G269 = int32(uint32(g[267]))
	m.G270 = int32(uint32(g[268]))
	m.G271 = int32(uint32(g[269]))
	m.G272 = int32(uint32(g[270]))
	m.G273 = int32(uint32(g[271]))
	m.G274 = int32(uint32(g[272]))
	m.G275 = int32(uint32(g[273]))
	m.G276 = int32(uint32(g[274]))
	m.G277 = int32(uint32(g[275]))
	m.G278 = int32(uint32(g[276]))
	m.G279 = int32(uint32(g[277]))
	m.G280 = int32(uint32(g[278]))
	m.G281 = int32(uint32(g[279]))
	m.G282 = int32(uint32(g[280]))
	m.G283 = int32(uint32(g[281]))
	m.G284 = int32(uint32(g[282]))
	m.G285 = int32(uint32(g[283]))
	m.G286 = int32(uint32(g[284]))
	m.G287 = int32(uint32(g[285]))
	m.G288 = int32(uint32(g[286]))
	m.G289 = int32(uint32(g[287]))
	m.G290 = int32(uint32(g[288]))
	m.G291 = int32(uint32(g[289]))
	m.G292 = int32(uint32(g[290]))
	m.G293 = int32(uint32(g[291]))
	m.G294 = int32(uint32(g[292]))
	m.G295 = int32(uint32(g[293]))
	m.G296 = int32(uint32(g[294]))
	m.G297 = int32(uint32(g[295]))
	m.G298 = int32(uint32(g[296]))
	m.G299 = int32(uint32(g[297]))
	m.G300 = int32(uint32(g[298]))
	m.G301 = int32(uint32(g[299]))
	m.G302 = int32(uint32(g[300]))
	m.G303 = int32(uint32(g[301]))
	m.G304 = int32(uint32(g[302]))
	m.G305 = int32(uint32(g[303]))
	m.G306 = int32(uint32(g[304]))
	m.G307 = int32(uint32(g[305]))
	m.G308 = int32(uint32(g[306]))
	m.G309 = int32(uint32(g[307]))
	m.G310 = int32(uint32(g[308]))
	m.G311 = int32(uint32(g[309]))
	m.G312 = int32(uint32(g[310]))
	m.G313 = int32(uint32(g[311]))
	m.G314 = int32(uint32(g[312]))
	m.G315 = int32(uint32(g[313]))
	m.G316 = int32(uint32(g[314]))
	m.G317 = int32(uint32(g[315]))
	m.G318 = int32(uint32(g[316]))
	m.G319 = int32(uint32(g[317]))
	m.G320 = int32(uint32(g[318]))
	m.G321 = int32(uint32(g[319]))
	m.G322 = int32(uint32(g[320]))
	m.G323 = int32(uint32(g[321]))
	m.G324 = int32(uint32(g[322]))
	m.G325 = int32(uint32(g[323]))
	m.G326 = int32(uint32(g[324]))
	m.G327 = int32(uint32(g[325]))
	m.G328 = int32(uint32(g[326]))
	m.G329 = int32(uint32(g[327]))
	m.G330 = int32(uint32(g[328]))
	m.G331 = int32(uint32(g[329]))
	m.G332 = int32(uint32(g[330]))
	m.G333 = int32(uint32(g[331]))
	m.G334 = int32(uint32(g[332]))
	m.G335 = int32(uint32(g[333]))
	m.G336 = int32(uint32(g[334]))
	m.G337 = int32(uint32(g[335]))
	m.G338 = int32(uint32(g[336]))
	m.G339 = int32(uint32(g[337]))
	m.G340 = int32(uint32(g[338]))
	m.G341 = int32(uint32(g[339]))
	m.G342 = int32(uint32(g[340]))
	m.G343 = int32(uint32(g[341]))
	m.G344 = int32(uint32(g[342]))
	m.G345 = int32(uint32(g[343]))
	m.G346 = int32(uint32(g[344]))
	m.G347 = int32(uint32(g[345]))
	m.G348 = int32(uint32(g[346]))
	m.G349 = int32(uint32(g[347]))
	m.G350 = int32(uint32(g[348]))
	m.G351 = int32(uint32(g[349]))
	m.G352 = int32(uint32(g[350]))
	m.G353 = int32(uint32(g[351]))
	m.G354 = int32(uint32(g[352]))
	m.G355 = int32(uint32(g[353]))
	m.G356 = int32(uint32(g[354]))
	m.G357 = int32(uint32(g[355]))
	m.G358 = int32(uint32(g[356]))
	m.G359 = int32(uint32(g[357]))
	m.G360 = int32(uint32(g[358]))
	m.G361 = int32(uint32(g[359]))
	m.G362 = int32(uint32(g[360]))
	m.G363 = int32(uint32(g[361]))
	m.G364 = int32(uint32(g[362]))
	m.G365 = int32(uint32(g[363]))
	m.G366 = int32(uint32(g[364]))
	m.G367 = int32(uint32(g[365]))
	m.G368 = int32(uint32(g[366]))
	m.G369 = int32(uint32(g[367]))
	m.G370 = int32(uint32(g[368]))
	m.G371 = int32(uint32(g[369]))
	m.G372 = int32(uint32(g[370]))
	m.G373 = int32(uint32(g[371]))
	m.G374 = int32(uint32(g[372]))
	m.G375 = int32(uint32(g[373]))
	m.G376 = int32(uint32(g[374]))
	m.G377 = int32(uint32(g[375]))
	m.G378 = int32(uint32(g[376]))
	m.G379 = int32(uint32(g[377]))
	m.G380 = int32(uint32(g[378]))
	m.G381 = int32(uint32(g[379]))
	m.G382 = int32(uint32(g[380]))
	m.G383 = int32(uint32(g[381]))
	m.G384 = int32(uint32(g[382]))
	m.G385 = int32(uint32(g[383]))
	m.G386 = int32(uint32(g[384]))
	m.G387 = int32(uint32(g[385]))
	m.G388 = int32(uint32(g[386]))
	m.G389 = int32(uint32(g[387]))
	m.G390 = int32(uint32(g[388]))
	m.G391 = int32(uint32(g[389]))
	m.G392 = int32(uint32(g[390]))
	m.G393 = int32(uint32(g[391]))
	m.G394 = int32(uint32(g[392]))
	m.G395 = int32(uint32(g[393]))
	m.G396 = int32(uint32(g[394]))
	m.G397 = int32(uint32(g[395]))
	m.G398 = int32(uint32(g[396]))
	m.G399 = int32(uint32(g[397]))
	m.G400 = int32(uint32(g[398]))
	m.G401 = int32(uint32(g[399]))
	m.G402 = int32(uint32(g[400]))
	m.G403 = int32(uint32(g[401]))
	m.G404 = int32(uint32(g[402]))
	m.G405 = int32(uint32(g[403]))
	m.G406 = int32(uint32(g[404]))
	m.G407 = int32(uint32(g[405]))
	m.G408 = int32(uint32(g[406]))
	m.G409 = int32(uint32(g[407]))
	m.G410 = int32(uint32(g[408]))
	m.G411 = int32(uint32(g[409]))
	m.G412 = int32(uint32(g[410]))
	m.G413 = int32(uint32(g[411]))
	m.G414 = int32(uint32(g[412]))
	m.G415 = int32(uint32(g[413]))
	m.G416 = int32(uint32(g[414]))
	m.G417 = int32(uint32(g[415]))
	m.G418 = int32(uint32(g[416]))
	m.G419 = int32(uint32(g[417]))
	m.G420 = int32(uint32(g[418]))
	m.G421 = int32(uint32(g[419]))
	m.G422 = int32(uint32(g[420]))
	m.G423 = int32(uint32(g[421]))
	m.G424 = int32(uint32(g[422]))
	m.G425 = int32(uint32(g[423]))
	m.G426 = int32(uint32(g[424]))
	m.G427 = int32(uint32(g[425]))
	m.G428 = int32(uint32(g[426]))
	m.G429 = int32(uint32(g[427]))
	m.G430 = int32(uint32(g[428]))
	m.G431 = int32(uint32(g[429]))
	m.G432 = int32(uint32(g[430]))
	m.G433 = int32(uint32(g[431]))
	m.G434 = int32(uint32(g[432]))
	m.G435 = int32(uint32(g[433]))
	m.G436 = int32(uint32(g[434]))
	m.G437 = int32(uint32(g[435]))
	m.G438 = int32(uint32(g[436]))
	m.G439 = int32(uint32(g[437]))
	m.G440 = int32(uint32(g[438]))
	m.G441 = int32(uint32(g[439]))
	m.G442 = int32(uint32(g[440]))
	m.G443 = int32(uint32(g[441]))
	m.G444 = int32(uint32(g[442]))
	m.G445 = int32(uint32(g[443]))
	m.G446 = int32(uint32(g[444]))
	m.G447 = int32(uint32(g[445]))
	m.G448 = int32(uint32(g[446]))
	m.G449 = int32(uint32(g[447]))
	m.G450 = int32(uint32(g[448]))
	m.G451 = int32(uint32(g[449]))
	m.G452 = int32(uint32(g[450]))
	m.G453 = int32(uint32(g[451]))
	m.G454 = int32(uint32(g[452]))
	m.G455 = int32(uint32(g[453]))
	m.G456 = int32(uint32(g[454]))
	m.G457 = int32(uint32(g[455]))
	m.G458 = int32(uint32(g[456]))
	m.G459 = int32(uint32(g[457]))
	m.G460 = int32(uint32(g[458]))
	m.G461 = int32(uint32(g[459]))
	m.G462 = int32(uint32(g[460]))
	m.G463 = int32(uint32(g[461]))
	m.G464 = int32(uint32(g[462]))
	m.G465 = int32(uint32(g[463]))
	m.G466 = int32(uint32(g[464]))
	m.G467 = int32(uint32(g[465]))
	m.G468 = int32(uint32(g[466]))
	m.G469 = int32(uint32(g[467]))
	m.G470 = int32(uint32(g[468]))
	m.G471 = int32(uint32(g[469]))
	m.G472 = int32(uint32(g[470]))
	m.G473 = int32(uint32(g[471]))
	m.G474 = int32(uint32(g[472]))
	m.G475 = int32(uint32(g[473]))
	m.G476 = int32(uint32(g[474]))
	m.G477 = int32(uint32(g[475]))
	m.G478 = int32(uint32(g[476]))
	m.G479 = int32(uint32(g[477]))
	m.G480 = int32(uint32(g[478]))
	m.G481 = int32(uint32(g[479]))
	m.G482 = int32(uint32(g[480]))
	m.G483 = int32(uint32(g[481]))
	m.G484 = int32(uint32(g[482]))
	m.G485 = int32(uint32(g[483]))
	m.G486 = int32(uint32(g[484]))
	m.G487 = int32(uint32(g[485]))
	m.G488 = int32(uint32(g[486]))
	m.G489 = int32(uint32(g[487]))
	m.G490 = int32(uint32(g[488]))
	m.G491 = int32(uint32(g[489]))
	m.G492 = int32(uint32(g[490]))
	m.G493 = int32(uint32(g[491]))
	m.G494 = int32(uint32(g[492]))
	m.G495 = int32(uint32(g[493]))
	m.G496 = int32(uint32(g[494]))
	m.G497 = int32(uint32(g[495]))
	m.G498 = int32(uint32(g[496]))
	m.G499 = int32(uint32(g[497]))
	m.G500 = int32(uint32(g[498]))
	m.G501 = int32(uint32(g[499]))
	m.G502 = int32(uint32(g[500]))
	m.G503 = int32(uint32(g[501]))
	m.G504 = int32(uint32(g[502]))
	m.G505 = int32(uint32(g[503]))
	m.G506 = int32(uint32(g[504]))
	m.G507 = int32(uint32(g[505]))
	m.G508 = int32(uint32(g[506]))
	m.G509 = int32(uint32(g[507]))
	m.G510 = int32(uint32(g[508]))
	m.G511 = int32(uint32(g[509]))
	m.G512 = int32(uint32(g[510]))
	m.G513 = int32(uint32(g[511]))
	m.G514 = int32(uint32(g[512]))
	m.G515 = int32(uint32(g[513]))
	m.G516 = int32(uint32(g[514]))
	m.G517 = int32(uint32(g[515]))
	m.G518 = int32(uint32(g[516]))
	m.G519 = int32(uint32(g[517]))
	m.G520 = int32(uint32(g[518]))
	m.G521 = int32(uint32(g[519]))
	m.G522 = int32(uint32(g[520]))
	m.G523 = int32(uint32(g[521]))
	m.G524 = int32(uint32(g[522]))
	m.G525 = int32(uint32(g[523]))
	m.G526 = int32(uint32(g[524]))
	m.G527 = int32(uint32(g[525]))
	m.G528 = int32(uint32(g[526]))
	m.G529 = int32(uint32(g[527]))
	m.G530 = int32(uint32(g[528]))
	m.G531 = int32(uint32(g[529]))
	m.G532 = int32(uint32(g[530]))
	m.G533 = int32(uint32(g[531]))
	m.G534 = int32(uint32(g[532]))
	m.G535 = int32(uint32(g[533]))
	m.G536 = int32(uint32(g[534]))
	m.G537 = int32(uint32(g[535]))
	m.G538 = int32(uint32(g[536]))
	m.G539 = int32(uint32(g[537]))
	m.G540 = int32(uint32(g[538]))
	m.G541 = int32(uint32(g[539]))
	m.G542 = int32(uint32(g[540]))
	m.G543 = int32(uint32(g[541]))
	m.G544 = int32(uint32(g[542]))
	m.G545 = int32(uint32(g[543]))
	m.G546 = int32(uint32(g[544]))
	m.G547 = int32(uint32(g[545]))
	m.G548 = int32(uint32(g[546]))
	m.G549 = int32(uint32(g[547]))
	m.G550 = int32(uint32(g[548]))
	m.G551 = int32(uint32(g[549]))
	m.G552 = int32(uint32(g[550]))
	m.G553 = int32(uint32(g[551]))
	m.G554 = int32(uint32(g[552]))
	m.G555 = int32(uint32(g[553]))
	m.G556 = int32(uint32(g[554]))
	m.G557 = int32(uint32(g[555]))
	m.G558 = int32(uint32(g[556]))
	m.G559 = int32(uint32(g[557]))
	m.G560 = int32(uint32(g[558]))
	m.G561 = int32(uint32(g[559]))
	m.G562 = int32(uint32(g[560]))
	m.G563 = int32(uint32(g[561]))
	m.G564 = int32(uint32(g[562]))
	m.G565 = int32(uint32(g[563]))
	m.G566 = int32(uint32(g[564]))
	m.G567 = int32(uint32(g[565]))
	m.G568 = int32(uint32(g[566]))
	m.G569 = int32(uint32(g[567]))
	m.G570 = int32(uint32(g[568]))
	m.G571 = int32(uint32(g[569]))
	m.G572 = int32(uint32(g[570]))
	m.G573 = int32(uint32(g[571]))
	m.G574 = int32(uint32(g[572]))
	m.G575 = int32(uint32(g[573]))
	m.G576 = int32(uint32(g[574]))
	m.G577 = int32(uint32(g[575]))
	m.G578 = int32(uint32(g[576]))
	m.G579 = int32(uint32(g[577]))
	m.G580 = int32(uint32(g[578]))
	m.G581 = int32(uint32(g[579]))
	m.G582 = int32(uint32(g[580]))
	m.G583 = int32(uint32(g[581]))
	m.G584 = int32(uint32(g[582]))
	m.G585 = int32(uint32(g[583]))
	m.G586 = int32(uint32(g[584]))
	m.G587 = int32(uint32(g[585]))
	m.G588 = int32(uint32(g[586]))
	m.G589 = int32(uint32(g[587]))
	m.G590 = int32(uint32(g[588]))
	m.G591 = int32(uint32(g[589]))
	m.G592 = int32(uint32(g[590]))
	m.G593 = int32(uint32(g[591]))
	m.G594 = int32(uint32(g[592]))
	m.G595 = int32(uint32(g[593]))
	m.G596 = int32(uint32(g[594]))
	m.G597 = int32(uint32(g[595]))
	m.G598 = int32(uint32(g[596]))
	m.G599 = int32(uint32(g[597]))
	m.G600 = int32(uint32(g[598]))
	m.G601 = int32(uint32(g[599]))
	m.G602 = int32(uint32(g[600]))
	m.G603 = int32(uint32(g[601]))
	m.G604 = int32(uint32(g[602]))
	m.G605 = int32(uint32(g[603]))
	m.G606 = int32(uint32(g[604]))
	m.G607 = int32(uint32(g[605]))
	m.G608 = int32(uint32(g[606]))
	m.G609 = int32(uint32(g[607]))
	m.G610 = int32(uint32(g[608]))
	m.G611 = int32(uint32(g[609]))
	m.G612 = int32(uint32(g[610]))
	m.G613 = int32(uint32(g[611]))
	m.G614 = int32(uint32(g[612]))
	m.G615 = int32(uint32(g[613]))
	m.G616 = int32(uint32(g[614]))
	m.G617 = int32(uint32(g[615]))
	m.G618 = int32(uint32(g[616]))
	m.G619 = int32(uint32(g[617]))
	m.G620 = int32(uint32(g[618]))
	m.G621 = int32(uint32(g[619]))
	m.G622 = int32(uint32(g[620]))
	m.G623 = int32(uint32(g[621]))
	m.G624 = int32(uint32(g[622]))
	m.G625 = int32(uint32(g[623]))
	m.G626 = int32(uint32(g[624]))
	m.G627 = int32(uint32(g[625]))
	m.G628 = int32(uint32(g[626]))
	m.G629 = int32(uint32(g[627]))
	m.G630 = int32(uint32(g[628]))
	m.G631 = int32(uint32(g[629]))
	m.G632 = int32(uint32(g[630]))
	m.G633 = int32(uint32(g[631]))
	m.G634 = int32(uint32(g[632]))
	m.G635 = int32(uint32(g[633]))
	m.G636 = int32(uint32(g[634]))
	m.G637 = int32(uint32(g[635]))
	m.G638 = int32(uint32(g[636]))
	m.G639 = int32(uint32(g[637]))
	m.G640 = int32(uint32(g[638]))
	m.G641 = int32(uint32(g[639]))
	m.G642 = int32(uint32(g[640]))
	m.G643 = int32(uint32(g[641]))
	m.G644 = int32(uint32(g[642]))
	m.G645 = int32(uint32(g[643]))
	m.G646 = int32(uint32(g[644]))
	m.G647 = int32(uint32(g[645]))
	m.G648 = int32(uint32(g[646]))
	m.G649 = int32(uint32(g[647]))
	m.G650 = int32(uint32(g[648]))
	m.G651 = int32(uint32(g[649]))
	m.G652 = int32(uint32(g[650]))
	m.G653 = int32(uint32(g[651]))
	m.G654 = int32(uint32(g[652]))
	m.G655 = int32(uint32(g[653]))
	m.G656 = int32(uint32(g[654]))
	m.G657 = int32(uint32(g[655]))
	m.G658 = int32(uint32(g[656]))
	m.G659 = int32(uint32(g[657]))
	m.G660 = int32(uint32(g[658]))
	m.G661 = int32(uint32(g[659]))
	m.G662 = int32(uint32(g[660]))
	m.G663 = int32(uint32(g[661]))
	m.G664 = int32(uint32(g[662]))
	m.G665 = int32(uint32(g[663]))
	m.G666 = int32(uint32(g[664]))
	m.G667 = int32(uint32(g[665]))
	m.G668 = int32(uint32(g[666]))
	m.G669 = int32(uint32(g[667]))
	m.G670 = int32(uint32(g[668]))
	m.G671 = int32(uint32(g[669]))
	m.G672 = int32(uint32(g[670]))
	m.G673 = int32(uint32(g[671]))
	m.G674 = int32(uint32(g[672]))
	m.G675 = int32(uint32(g[673]))
	m.G676 = int32(uint32(g[674]))
	m.G677 = int32(uint32(g[675]))
	m.G678 = int32(uint32(g[676]))
	m.G679 = int32(uint32(g[677]))
	m.G680 = int32(uint32(g[678]))
	m.G681 = int32(uint32(g[679]))
	m.G682 = int32(uint32(g[680]))
	m.G683 = int32(uint32(g[681]))
	m.G684 = int32(uint32(g[682]))
	m.G685 = int32(uint32(g[683]))
	m.G686 = int32(uint32(g[684]))
	m.G687 = int32(uint32(g[685]))
	m.G688 = int32(uint32(g[686]))
	m.G689 = int32(uint32(g[687]))
	m.G690 = int32(uint32(g[688]))
	m.G691 = int32(uint32(g[689]))
	m.G692 = int32(uint32(g[690]))
	m.G693 = int32(uint32(g[691]))
	m.G694 = int32(uint32(g[692]))
	m.G695 = int32(uint32(g[693]))
	m.G696 = int32(uint32(g[694]))
	m.G697 = int32(uint32(g[695]))
	m.G698 = int32(uint32(g[696]))
	m.G699 = int32(uint32(g[697]))
	m.G700 = int32(uint32(g[698]))
	m.G701 = int32(uint32(g[699]))
	m.G702 = int32(uint32(g[700]))
	m.G703 = int32(uint32(g[701]))
	m.G704 = int32(uint32(g[702]))
	m.G705 = int32(uint32(g[703]))
	m.G706 = int32(uint32(g[704]))
	m.G707 = int32(uint32(g[705]))
	m.G708 = int32(uint32(g[706]))
	m.G709 = int32(uint32(g[707]))
	m.G710 = int32(uint32(g[708]))
	m.G711 = int32(uint32(g[709]))
	m.G712 = int32(uint32(g[710]))
	m.G713 = int32(uint32(g[711]))
	m.G714 = int32(uint32(g[712]))
	m.G715 = int32(uint32(g[713]))
	m.G716 = int32(uint32(g[714]))
	m.G717 = int32(uint32(g[715]))
	m.G718 = int32(uint32(g[716]))
	m.G719 = int32(uint32(g[717]))
	m.G720 = int32(uint32(g[718]))
	m.G721 = int32(uint32(g[719]))
	m.G722 = int32(uint32(g[720]))
	m.G723 = int32(uint32(g[721]))
	m.G724 = int32(uint32(g[722]))
	m.G725 = int32(uint32(g[723]))
	m.G726 = int32(uint32(g[724]))
	m.G727 = int32(uint32(g[725]))
	m.G728 = int32(uint32(g[726]))
	m.G729 = int32(uint32(g[727]))
	m.G730 = int32(uint32(g[728]))
	m.G731 = int32(uint32(g[729]))
	m.G732 = int32(uint32(g[730]))
	m.G733 = int32(uint32(g[731]))
	m.G734 = int32(uint32(g[732]))
	m.G735 = int32(uint32(g[733]))
	m.G736 = int32(uint32(g[734]))
	m.G737 = int32(uint32(g[735]))
	m.G738 = int32(uint32(g[736]))
	m.G739 = int32(uint32(g[737]))
	m.G740 = int32(uint32(g[738]))
	m.G741 = int32(uint32(g[739]))
	m.G742 = int32(uint32(g[740]))
	m.G743 = int32(uint32(g[741]))
	m.G744 = int32(uint32(g[742]))
	m.G745 = int32(uint32(g[743]))
	m.G746 = int32(uint32(g[744]))
	m.G747 = int32(uint32(g[745]))
	m.G748 = int32(uint32(g[746]))
	m.G749 = int32(uint32(g[747]))
	m.G750 = int32(uint32(g[748]))
	m.G751 = int32(uint32(g[749]))
	m.G752 = int32(uint32(g[750]))
	m.G753 = int32(uint32(g[751]))
	m.G754 = int32(uint32(g[752]))
	m.G755 = int32(uint32(g[753]))
	m.G756 = int32(uint32(g[754]))
	m.G757 = int32(uint32(g[755]))
	m.G758 = int32(uint32(g[756]))
	m.G759 = int32(uint32(g[757]))
	m.G760 = int32(uint32(g[758]))
	m.G761 = int32(uint32(g[759]))
	m.G762 = int32(uint32(g[760]))
	m.G763 = int32(uint32(g[761]))
	m.G764 = int32(uint32(g[762]))
	m.G765 = int32(uint32(g[763]))
	m.G766 = int32(uint32(g[764]))
	m.G767 = int32(uint32(g[765]))
	m.G768 = int32(uint32(g[766]))
	m.G769 = int32(uint32(g[767]))
	m.G770 = int32(uint32(g[768]))
	m.G771 = int32(uint32(g[769]))
	m.G772 = int32(uint32(g[770]))
	m.G773 = int32(uint32(g[771]))
	m.G774 = int32(uint32(g[772]))
	m.G775 = int32(uint32(g[773]))
	m.G776 = int32(uint32(g[774]))
	m.G777 = int32(uint32(g[775]))
	m.G778 = int32(uint32(g[776]))
	m.G779 = int32(uint32(g[777]))
	m.G780 = int32(uint32(g[778]))
	m.G781 = int32(uint32(g[779]))
	m.G782 = int32(uint32(g[780]))
	m.G783 = int32(uint32(g[781]))
	m.G784 = int32(uint32(g[782]))
	m.G785 = int32(uint32(g[783]))
	m.G786 = int32(uint32(g[784]))
	m.G787 = int32(uint32(g[785]))
	m.G788 = int32(uint32(g[786]))
	m.G789 = int32(uint32(g[787]))
	m.G790 = int32(uint32(g[788]))
	m.G791 = int32(uint32(g[789]))
	m.G792 = int32(uint32(g[790]))
	m.G793 = int32(uint32(g[791]))
	m.G794 = int32(uint32(g[792]))
	m.G795 = int32(uint32(g[793]))
	m.G796 = int32(uint32(g[794]))
	m.G797 = int32(uint32(g[795]))
	m.G798 = int32(uint32(g[796]))
	m.G799 = int32(uint32(g[797]))
	m.G800 = int32(uint32(g[798]))
	m.G801 = int32(uint32(g[799]))
	m.G802 = int32(uint32(g[800]))
	m.G803 = int32(uint32(g[801]))
	m.G804 = int32(uint32(g[802]))
	m.G805 = int32(uint32(g[803]))
	m.G806 = int32(uint32(g[804]))
	m.G807 = int32(uint32(g[805]))
	m.G808 = int32(uint32(g[806]))
	m.G809 = int32(uint32(g[807]))
	m.G810 = int32(uint32(g[808]))
	m.G811 = int32(uint32(g[809]))
	m.G812 = int32(uint32(g[810]))
	m.G813 = int32(uint32(g[811]))
	m.G814 = int32(uint32(g[812]))
	m.G815 = int32(uint32(g[813]))
	m.G816 = int32(uint32(g[814]))
	m.G817 = int32(uint32(g[815]))
	m.G818 = int32(uint32(g[816]))
	m.G819 = int32(uint32(g[817]))
	m.G820 = int32(uint32(g[818]))
	m.G821 = int32(uint32(g[819]))
	m.G822 = int32(uint32(g[820]))
	m.G823 = int32(uint32(g[821]))
	m.G824 = int32(uint32(g[822]))
	m.G825 = int32(uint32(g[823]))
	m.G826 = int32(uint32(g[824]))
	m.G827 = int32(uint32(g[825]))
	m.G828 = int32(uint32(g[826]))
	m.G829 = int32(uint32(g[827]))
	m.G830 = int32(uint32(g[828]))
	m.G831 = int32(uint32(g[829]))
	m.G832 = int32(uint32(g[830]))
	m.G833 = int32(uint32(g[831]))
	m.G834 = int32(uint32(g[832]))
	m.G835 = int32(uint32(g[833]))
	m.G836 = int32(uint32(g[834]))
	m.G837 = int32(uint32(g[835]))
	m.G838 = int32(uint32(g[836]))
	m.G839 = int32(uint32(g[837]))
	m.G840 = int32(uint32(g[838]))
	m.G841 = int32(uint32(g[839]))
	m.G842 = int32(uint32(g[840]))
	m.G843 = int32(uint32(g[841]))
	m.G844 = int32(uint32(g[842]))
	m.G845 = int32(uint32(g[843]))
	m.G846 = int32(uint32(g[844]))
	m.G847 = int32(uint32(g[845]))
	m.G848 = int32(uint32(g[846]))
	m.G849 = int32(uint32(g[847]))
	m.G850 = int32(uint32(g[848]))
	m.G851 = int32(uint32(g[849]))
	m.G852 = int32(uint32(g[850]))
	m.G853 = int32(uint32(g[851]))
	m.G854 = int32(uint32(g[852]))
	m.G855 = int32(uint32(g[853]))
	m.G856 = int32(uint32(g[854]))
	m.G857 = int32(uint32(g[855]))
	m.G858 = int32(uint32(g[856]))
	m.G859 = int32(uint32(g[857]))
	m.G860 = int32(uint32(g[858]))
	m.G861 = int32(uint32(g[859]))
	m.G862 = int32(uint32(g[860]))
	m.G863 = int32(uint32(g[861]))
	m.G864 = int32(uint32(g[862]))
	m.G865 = int32(uint32(g[863]))
	m.G866 = int32(uint32(g[864]))
	m.G867 = int32(uint32(g[865]))
	m.G868 = int32(uint32(g[866]))
	m.G869 = int32(uint32(g[867]))
	m.G870 = int32(uint32(g[868]))
	m.G871 = int32(uint32(g[869]))
	m.G872 = int32(uint32(g[870]))
	m.G873 = int32(uint32(g[871]))
	m.G874 = int32(uint32(g[872]))
	m.G875 = int32(uint32(g[873]))
	m.G876 = int32(uint32(g[874]))
	m.G877 = int32(uint32(g[875]))
	m.G878 = int32(uint32(g[876]))
	m.G879 = int32(uint32(g[877]))
	m.G880 = int32(uint32(g[878]))
	m.G881 = int32(uint32(g[879]))
	m.G882 = int32(uint32(g[880]))
	m.G883 = int32(uint32(g[881]))
	m.G884 = int32(uint32(g[882]))
	m.G885 = int32(uint32(g[883]))
	m.G886 = int32(uint32(g[884]))
	m.G887 = int32(uint32(g[885]))
	m.G888 = int32(uint32(g[886]))
	m.G889 = int32(uint32(g[887]))
	m.G890 = int32(uint32(g[888]))
	m.G891 = int32(uint32(g[889]))
	m.G892 = int32(uint32(g[890]))
	m.G893 = int32(uint32(g[891]))
	m.G894 = int32(uint32(g[892]))
	m.G895 = int32(uint32(g[893]))
	m.G896 = int32(uint32(g[894]))
	m.G897 = int32(uint32(g[895]))
	m.G898 = int32(uint32(g[896]))
	m.G899 = int32(uint32(g[897]))
	m.G900 = int32(uint32(g[898]))
	m.G901 = int32(uint32(g[899]))
	m.G902 = int32(uint32(g[900]))
	m.G903 = int32(uint32(g[901]))
	m.G904 = int32(uint32(g[902]))
	m.G905 = int32(uint32(g[903]))
	m.G906 = int32(uint32(g[904]))
	m.G907 = int32(uint32(g[905]))
	m.G908 = int32(uint32(g[906]))
	m.G909 = int32(uint32(g[907]))
	m.G910 = int32(uint32(g[908]))
	m.G911 = int32(uint32(g[909]))
	m.G912 = int32(uint32(g[910]))
	m.G913 = int32(uint32(g[911]))
	m.G914 = int32(uint32(g[912]))
	m.G915 = int32(uint32(g[913]))
	m.G916 = int32(uint32(g[914]))
	m.G917 = int32(uint32(g[915]))
	m.G918 = int32(uint32(g[916]))
	m.G919 = int32(uint32(g[917]))
	m.G920 = int32(uint32(g[918]))
	m.G921 = int32(uint32(g[919]))
	m.G922 = int32(uint32(g[920]))
	m.G923 = int32(uint32(g[921]))
	m.G924 = int32(uint32(g[922]))
	m.G925 = int32(uint32(g[923]))
	m.G926 = int32(uint32(g[924]))
	m.G927 = int32(uint32(g[925]))
	m.G928 = int32(uint32(g[926]))
	m.G929 = int32(uint32(g[927]))
	m.G930 = int32(uint32(g[928]))
	m.G931 = int32(uint32(g[929]))
	m.G932 = int32(uint32(g[930]))
	m.G933 = int32(uint32(g[931]))
	m.G934 = int32(uint32(g[932]))
	m.G935 = int32(uint32(g[933]))
	m.G936 = int32(uint32(g[934]))
	m.G937 = int32(uint32(g[935]))
	m.G938 = int32(uint32(g[936]))
	m.G939 = int32(uint32(g[937]))
	m.G940 = int32(uint32(g[938]))
	m.G941 = int32(uint32(g[939]))
	m.G942 = int32(uint32(g[940]))
	m.G943 = int32(uint32(g[941]))
	m.G944 = int32(uint32(g[942]))
	m.G945 = int32(uint32(g[943]))
	m.G946 = int32(uint32(g[944]))
	m.G947 = int32(uint32(g[945]))
	m.G948 = int32(uint32(g[946]))
	m.G949 = int32(uint32(g[947]))
	m.G950 = int32(uint32(g[948]))
	m.G951 = int32(uint32(g[949]))
	m.G952 = int32(uint32(g[950]))
	m.G953 = int32(uint32(g[951]))
	m.G954 = int32(uint32(g[952]))
	m.G955 = int32(uint32(g[953]))
	m.G956 = int32(uint32(g[954]))
	m.G957 = int32(uint32(g[955]))
	m.G958 = int32(uint32(g[956]))
	m.G959 = int32(uint32(g[957]))
	m.G960 = int32(uint32(g[958]))
	m.G961 = int32(uint32(g[959]))
	m.G962 = int32(uint32(g[960]))
	m.G963 = int32(uint32(g[961]))
	m.G964 = int32(uint32(g[962]))
	m.G965 = int32(uint32(g[963]))
	m.G966 = int32(uint32(g[964]))
	m.G967 = int32(uint32(g[965]))
	m.G968 = int32(uint32(g[966]))
	m.G969 = int32(uint32(g[967]))
	m.G970 = int32(uint32(g[968]))
	m.G971 = int32(uint32(g[969]))
	m.G972 = int32(uint32(g[970]))
	m.G973 = int32(uint32(g[971]))
	m.G974 = int32(uint32(g[972]))
	m.G975 = int32(uint32(g[973]))
	m.G976 = int32(uint32(g[974]))
	m.G977 = int32(uint32(g[975]))
	m.G978 = int32(uint32(g[976]))
	m.G979 = int32(uint32(g[977]))
	m.G980 = int32(uint32(g[978]))
	m.G981 = int32(uint32(g[979]))
	m.G982 = int32(uint32(g[980]))
	m.G983 = int32(uint32(g[981]))
	m.G984 = int32(uint32(g[982]))
	m.G985 = int32(uint32(g[983]))
	m.G986 = int32(uint32(g[984]))
	m.G987 = int32(uint32(g[985]))
	m.G988 = int32(uint32(g[986]))
	m.G989 = int32(uint32(g[987]))
	m.G990 = int32(uint32(g[988]))
	m.G991 = int32(uint32(g[989]))
	m.G992 = int32(uint32(g[990]))
	m.G993 = int32(uint32(g[991]))
	m.G994 = int32(uint32(g[992]))
	m.G995 = int32(uint32(g[993]))
	m.G996 = int32(uint32(g[994]))
	m.G997 = int32(uint32(g[995]))
	m.G998 = int32(uint32(g[996]))
	m.G999 = int32(uint32(g[997]))
	m.G1000 = int32(uint32(g[998]))
	m.G1001 = int32(uint32(g[999]))
	m.G1002 = int32(uint32(g[1000]))
	m.G1003 = int32(uint32(g[1001]))
	m.G1004 = int32(uint32(g[1002]))
	m.G1005 = int32(uint32(g[1003]))
	m.G1006 = int32(uint32(g[1004]))
	m.G1007 = int32(uint32(g[1005]))
	m.G1008 = int32(uint32(g[1006]))
	m.G1009 = int32(uint32(g[1007]))
	m.G1010 = int32(uint32(g[1008]))
	m.G1011 = int32(uint32(g[1009]))
	m.G1012 = int32(uint32(g[1010]))
	m.G1013 = int32(uint32(g[1011]))
	m.G1014 = int32(uint32(g[1012]))
	m.G1015 = int32(uint32(g[1013]))
	m.G1016 = int32(uint32(g[1014]))
	m.G1017 = int32(uint32(g[1015]))
	m.G1018 = int32(uint32(g[1016]))
	m.G1019 = int32(uint32(g[1017]))
	m.G1020 = int32(uint32(g[1018]))
	m.G1021 = int32(uint32(g[1019]))
	m.G1022 = int32(uint32(g[1020]))
	m.G1023 = int32(uint32(g[1021]))
	m.G1024 = int32(uint32(g[1022]))
	m.G1025 = int32(uint32(g[1023]))
	m.G1026 = int32(uint32(g[1024]))
	m.G1027 = int32(uint32(g[1025]))
	m.G1028 = int32(uint32(g[1026]))
	m.G1029 = int32(uint32(g[1027]))
	m.G1030 = int32(uint32(g[1028]))
	m.G1031 = int32(uint32(g[1029]))
	m.G1032 = int32(uint32(g[1030]))
	m.G1033 = int32(uint32(g[1031]))
	m.G1034 = int32(uint32(g[1032]))
	m.G1035 = int32(uint32(g[1033]))
	m.G1036 = int32(uint32(g[1034]))
	m.G1037 = int32(uint32(g[1035]))
	m.G1038 = int32(uint32(g[1036]))
	m.G1039 = int32(uint32(g[1037]))
	m.G1040 = int32(uint32(g[1038]))
	m.G1041 = int32(uint32(g[1039]))
	m.G1042 = int32(uint32(g[1040]))
	m.G1043 = int32(uint32(g[1041]))
	m.G1044 = int32(uint32(g[1042]))
	m.G1045 = int32(uint32(g[1043]))
	m.G1046 = int32(uint32(g[1044]))
	m.G1047 = int32(uint32(g[1045]))
	m.G1048 = int32(uint32(g[1046]))
	m.G1049 = int32(uint32(g[1047]))
	m.G1050 = int32(uint32(g[1048]))
	m.G1051 = int32(uint32(g[1049]))
	m.G1052 = int32(uint32(g[1050]))
	m.G1053 = int32(uint32(g[1051]))
	m.G1054 = int32(uint32(g[1052]))
	m.G1055 = int32(uint32(g[1053]))
	m.G1056 = int32(uint32(g[1054]))
	m.G1057 = int32(uint32(g[1055]))
	m.G1058 = int32(uint32(g[1056]))
	m.G1059 = int32(uint32(g[1057]))
	m.G1060 = int32(uint32(g[1058]))
	m.G1061 = int32(uint32(g[1059]))
	m.G1062 = int32(uint32(g[1060]))
	m.G1063 = int32(uint32(g[1061]))
	m.G1064 = int32(uint32(g[1062]))
	m.G1065 = int32(uint32(g[1063]))
	m.G1066 = int32(uint32(g[1064]))
	m.G1067 = int32(uint32(g[1065]))
	m.G1068 = int32(uint32(g[1066]))
	m.G1069 = int32(uint32(g[1067]))
	m.G1070 = int32(uint32(g[1068]))
	m.G1071 = int32(uint32(g[1069]))
	m.G1072 = int32(uint32(g[1070]))
	m.G1073 = int32(uint32(g[1071]))
	m.G1074 = int32(uint32(g[1072]))
	m.G1075 = int32(uint32(g[1073]))
	m.G1076 = int32(uint32(g[1074]))
	m.G1077 = int32(uint32(g[1075]))
	m.G1078 = int32(uint32(g[1076]))
	m.G1079 = int32(uint32(g[1077]))
	m.G1080 = int32(uint32(g[1078]))
	m.G1081 = int32(uint32(g[1079]))
	m.G1082 = int32(uint32(g[1080]))
	m.G1083 = int32(uint32(g[1081]))
	m.G1084 = int32(uint32(g[1082]))
	m.G1085 = int32(uint32(g[1083]))
	m.G1086 = int32(uint32(g[1084]))
	m.G1087 = int32(uint32(g[1085]))
	m.G1088 = int32(uint32(g[1086]))
	m.G1089 = int32(uint32(g[1087]))
	m.G1090 = int32(uint32(g[1088]))
	m.G1091 = int32(uint32(g[1089]))
	m.G1092 = int32(uint32(g[1090]))
	m.G1093 = int32(uint32(g[1091]))
	m.G1094 = int32(uint32(g[1092]))
	m.G1095 = int32(uint32(g[1093]))
	m.G1096 = int32(uint32(g[1094]))
	m.G1097 = int32(uint32(g[1095]))
	m.G1098 = int32(uint32(g[1096]))
	m.G1099 = int32(uint32(g[1097]))
	m.G1100 = int32(uint32(g[1098]))
	m.G1101 = int32(uint32(g[1099]))
	m.G1102 = int32(uint32(g[1100]))
	m.G1103 = int32(uint32(g[1101]))
	m.G1104 = int32(uint32(g[1102]))
	m.G1105 = int32(uint32(g[1103]))
	m.G1106 = int32(uint32(g[1104]))
	m.G1107 = int32(uint32(g[1105]))
	m.G1108 = int32(uint32(g[1106]))
	m.G1109 = int32(uint32(g[1107]))
	m.G1110 = int32(uint32(g[1108]))
	m.G1111 = int32(uint32(g[1109]))
	m.G1112 = int32(uint32(g[1110]))
	m.G1113 = int32(uint32(g[1111]))
	m.G1114 = int32(uint32(g[1112]))
	m.G1115 = int32(uint32(g[1113]))
	m.G1116 = int32(uint32(g[1114]))
	m.G1117 = int32(uint32(g[1115]))
	m.G1118 = int32(uint32(g[1116]))
	m.G1119 = int32(uint32(g[1117]))
	m.G1120 = int32(uint32(g[1118]))
	m.G1121 = int32(uint32(g[1119]))
	m.G1122 = int32(uint32(g[1120]))
	m.G1123 = int32(uint32(g[1121]))
	m.G1124 = int32(uint32(g[1122]))
	m.G1125 = int32(uint32(g[1123]))
	m.G1126 = int32(uint32(g[1124]))
	m.G1127 = int32(uint32(g[1125]))
	m.G1128 = int32(uint32(g[1126]))
	m.G1129 = int32(uint32(g[1127]))
	m.G1130 = int32(uint32(g[1128]))
	m.G1131 = int32(uint32(g[1129]))
	m.G1132 = int32(uint32(g[1130]))
	m.G1133 = int32(uint32(g[1131]))
	m.G1134 = int32(uint32(g[1132]))
	m.G1135 = int32(uint32(g[1133]))
	m.G1136 = int32(uint32(g[1134]))
	m.G1137 = int32(uint32(g[1135]))
	m.G1138 = int32(uint32(g[1136]))
	m.G1139 = int32(uint32(g[1137]))
	m.G1140 = int32(uint32(g[1138]))
	m.G1141 = int32(uint32(g[1139]))
	m.G1142 = int32(uint32(g[1140]))
	m.G1143 = int32(uint32(g[1141]))
	m.G1144 = int32(uint32(g[1142]))
	m.G1145 = int32(uint32(g[1143]))
	m.G1146 = int32(uint32(g[1144]))
	m.G1147 = int32(uint32(g[1145]))
	m.G1148 = int32(uint32(g[1146]))
	m.G1149 = int32(uint32(g[1147]))
	m.G1150 = int32(uint32(g[1148]))
	m.G1151 = int32(uint32(g[1149]))
	m.G1152 = int32(uint32(g[1150]))
	m.G1153 = int32(uint32(g[1151]))
	m.G1154 = int32(uint32(g[1152]))
	m.G1155 = int32(uint32(g[1153]))
	m.G1156 = int32(uint32(g[1154]))
	m.G1157 = int32(uint32(g[1155]))
	m.G1158 = int32(uint32(g[1156]))
	m.G1159 = int32(uint32(g[1157]))
	m.G1160 = int32(uint32(g[1158]))
	m.G1161 = int32(uint32(g[1159]))
	m.G1162 = int32(uint32(g[1160]))
	m.G1163 = int32(uint32(g[1161]))
	m.G1164 = int32(uint32(g[1162]))
	m.G1165 = int32(uint32(g[1163]))
	m.G1166 = int32(uint32(g[1164]))
	m.G1167 = int32(uint32(g[1165]))
	m.G1168 = int32(uint32(g[1166]))
	m.G1169 = int32(uint32(g[1167]))
	m.G1170 = int32(uint32(g[1168]))
	m.G1171 = int32(uint32(g[1169]))
	m.G1172 = int32(uint32(g[1170]))
	m.G1173 = int32(uint32(g[1171]))
	m.G1174 = int32(uint32(g[1172]))
	m.G1175 = int32(uint32(g[1173]))
	m.G1176 = int32(uint32(g[1174]))
	m.G1177 = int32(uint32(g[1175]))
	m.G1178 = int32(uint32(g[1176]))
	m.G1179 = int32(uint32(g[1177]))
	m.G1180 = int32(uint32(g[1178]))
	m.G1181 = int32(uint32(g[1179]))
	m.G1182 = int32(uint32(g[1180]))
	m.G1183 = int32(uint32(g[1181]))
	m.G1184 = int32(uint32(g[1182]))
	m.G1185 = int32(uint32(g[1183]))
	m.G1186 = int32(uint32(g[1184]))
	m.G1187 = int32(uint32(g[1185]))
	m.G1188 = int32(uint32(g[1186]))
	m.G1189 = int32(uint32(g[1187]))
	m.G1190 = int32(uint32(g[1188]))
	m.G1191 = int32(uint32(g[1189]))
	m.G1192 = int32(uint32(g[1190]))
	m.G1193 = int32(uint32(g[1191]))
	m.G1194 = int32(uint32(g[1192]))
	m.G1195 = int32(uint32(g[1193]))
	m.G1196 = int32(uint32(g[1194]))
	m.G1197 = int32(uint32(g[1195]))
	m.G1198 = int32(uint32(g[1196]))
	m.G1199 = int32(uint32(g[1197]))
	m.G1200 = int32(uint32(g[1198]))
	m.G1201 = int32(uint32(g[1199]))
	m.G1202 = int32(uint32(g[1200]))
	m.G1203 = int32(uint32(g[1201]))
	m.G1204 = int32(uint32(g[1202]))
	m.G1205 = int32(uint32(g[1203]))
	m.G1206 = int32(uint32(g[1204]))
	m.G1207 = int32(uint32(g[1205]))
	m.G1208 = int32(uint32(g[1206]))
	m.G1209 = int32(uint32(g[1207]))
	m.G1210 = int32(uint32(g[1208]))
	m.G1211 = int32(uint32(g[1209]))
	m.G1212 = int32(uint32(g[1210]))
	m.G1213 = int32(uint32(g[1211]))
	m.G1214 = int32(uint32(g[1212]))
	m.G1215 = int32(uint32(g[1213]))
	m.G1216 = int32(uint32(g[1214]))
	m.G1217 = int32(uint32(g[1215]))
	m.G1218 = int32(uint32(g[1216]))
	m.G1219 = int32(uint32(g[1217]))
	m.G1220 = int32(uint32(g[1218]))
	m.G1221 = int32(uint32(g[1219]))
	m.G1222 = int32(uint32(g[1220]))
	m.G1223 = int32(uint32(g[1221]))
	m.G1224 = int32(uint32(g[1222]))
	m.G1225 = int32(uint32(g[1223]))
	m.G1226 = int32(uint32(g[1224]))
	m.G1227 = int32(uint32(g[1225]))
	m.G1228 = int32(uint32(g[1226]))
	m.G1229 = int32(uint32(g[1227]))
	m.G1230 = int32(uint32(g[1228]))
	m.G1231 = int32(uint32(g[1229]))
	m.G1232 = int32(uint32(g[1230]))
	m.G1233 = int32(uint32(g[1231]))
	m.G1234 = int32(uint32(g[1232]))
	m.G1235 = int32(uint32(g[1233]))
	m.G1236 = int32(uint32(g[1234]))
	m.G1237 = int32(uint32(g[1235]))
	m.G1238 = int32(uint32(g[1236]))
	m.G1239 = int32(uint32(g[1237]))
	m.G1240 = int32(uint32(g[1238]))
	m.G1241 = int32(uint32(g[1239]))
	m.G1242 = int32(uint32(g[1240]))
	m.G1243 = int32(uint32(g[1241]))
	m.G1244 = int32(uint32(g[1242]))
	m.G1245 = int32(uint32(g[1243]))
	m.G1246 = int32(uint32(g[1244]))
	m.G1247 = int32(uint32(g[1245]))
	m.G1248 = int32(uint32(g[1246]))
	m.G1249 = int32(uint32(g[1247]))
	m.G1250 = int32(uint32(g[1248]))
	m.G1251 = int32(uint32(g[1249]))
	m.G1252 = int32(uint32(g[1250]))
	m.G1253 = int32(uint32(g[1251]))
	m.G1254 = int32(uint32(g[1252]))
	m.G1255 = int32(uint32(g[1253]))
	m.G1256 = int32(uint32(g[1254]))
	m.G1257 = int32(uint32(g[1255]))
	m.G1258 = int32(uint32(g[1256]))
	m.G1259 = int32(uint32(g[1257]))
	m.G1260 = int32(uint32(g[1258]))
	m.G1261 = int32(uint32(g[1259]))
	m.G1262 = int32(uint32(g[1260]))
	m.G1263 = int32(uint32(g[1261]))
	m.G1264 = int32(uint32(g[1262]))
	m.G1265 = int32(uint32(g[1263]))
	m.G1266 = int32(uint32(g[1264]))
	m.G1267 = int32(uint32(g[1265]))
	m.G1268 = int32(uint32(g[1266]))
	m.G1269 = int32(uint32(g[1267]))
	m.G1270 = int32(uint32(g[1268]))
	m.G1271 = int32(uint32(g[1269]))
	m.G1272 = int32(uint32(g[1270]))
	m.G1273 = int32(uint32(g[1271]))
	m.G1274 = int32(uint32(g[1272]))
	m.G1275 = int32(uint32(g[1273]))
	m.G1276 = int32(uint32(g[1274]))
	m.G1277 = int32(uint32(g[1275]))
	m.G1278 = int32(uint32(g[1276]))
	m.G1279 = int32(uint32(g[1277]))
	m.G1280 = int32(uint32(g[1278]))
	m.G1281 = int32(uint32(g[1279]))
	m.G1282 = int32(uint32(g[1280]))
	m.G1283 = int32(uint32(g[1281]))
	m.G1284 = int32(uint32(g[1282]))
	m.G1285 = int32(uint32(g[1283]))
	m.G1286 = int32(uint32(g[1284]))
	m.G1287 = int32(uint32(g[1285]))
	m.G1288 = int32(uint32(g[1286]))
	m.G1289 = int32(uint32(g[1287]))
	m.G1290 = int32(uint32(g[1288]))
	m.G1291 = int32(uint32(g[1289]))
	m.G1292 = int32(uint32(g[1290]))
	m.G1293 = int32(uint32(g[1291]))
	m.G1294 = int32(uint32(g[1292]))
	m.G1295 = int32(uint32(g[1293]))
	m.G1296 = int32(uint32(g[1294]))
	m.G1297 = int32(uint32(g[1295]))
	m.G1298 = int32(uint32(g[1296]))
	m.G1299 = int32(uint32(g[1297]))
	m.G1300 = int32(uint32(g[1298]))
	m.G1301 = int32(uint32(g[1299]))
	m.G1302 = int32(uint32(g[1300]))
	m.G1303 = int32(uint32(g[1301]))
	m.G1304 = int32(uint32(g[1302]))
	m.G1305 = int32(uint32(g[1303]))
	m.G1306 = int32(uint32(g[1304]))
	m.G1307 = int32(uint32(g[1305]))
	m.G1308 = int32(uint32(g[1306]))
	m.G1309 = int32(uint32(g[1307]))
	m.G1310 = int32(uint32(g[1308]))
	m.G1311 = int32(uint32(g[1309]))
	m.G1312 = int32(uint32(g[1310]))
	m.G1313 = int32(uint32(g[1311]))
	m.G1314 = int32(uint32(g[1312]))
	m.G1315 = int32(uint32(g[1313]))
	m.G1316 = int32(uint32(g[1314]))
	m.G1317 = int32(uint32(g[1315]))
	m.G1318 = int32(uint32(g[1316]))
	m.G1319 = int32(uint32(g[1317]))
	m.G1320 = int32(uint32(g[1318]))
	m.G1321 = int32(uint32(g[1319]))
	m.G1322 = int32(uint32(g[1320]))
	m.G1323 = int32(uint32(g[1321]))
	m.G1324 = int32(uint32(g[1322]))
	m.G1325 = int32(uint32(g[1323]))
	m.G1326 = int32(uint32(g[1324]))
	m.G1327 = int32(uint32(g[1325]))
	m.G1328 = int32(uint32(g[1326]))
	m.G1329 = int32(uint32(g[1327]))
	m.G1330 = int32(uint32(g[1328]))
	m.G1331 = int32(uint32(g[1329]))
	m.G1332 = int32(uint32(g[1330]))
	m.G1333 = int32(uint32(g[1331]))
	m.G1334 = int32(uint32(g[1332]))
	m.G1335 = int32(uint32(g[1333]))
	m.G1336 = int32(uint32(g[1334]))
	m.G1337 = int32(uint32(g[1335]))
	m.G1338 = int32(uint32(g[1336]))
	m.G1339 = int32(uint32(g[1337]))
	m.G1340 = int32(uint32(g[1338]))
	m.G1341 = int32(uint32(g[1339]))
	m.G1342 = int32(uint32(g[1340]))
	m.G1343 = int32(uint32(g[1341]))
	m.G1344 = int32(uint32(g[1342]))
	m.G1345 = int32(uint32(g[1343]))
	m.G1346 = int32(uint32(g[1344]))
	m.G1347 = int32(uint32(g[1345]))
	m.G1348 = int32(uint32(g[1346]))
	m.G1349 = int32(uint32(g[1347]))
	m.G1350 = int32(uint32(g[1348]))
	m.G1351 = int32(uint32(g[1349]))
	m.G1352 = int32(uint32(g[1350]))
	m.G1353 = int32(uint32(g[1351]))
	m.G1354 = int32(uint32(g[1352]))
	m.G1355 = int32(uint32(g[1353]))
	m.G1356 = int32(uint32(g[1354]))
	m.G1357 = int32(uint32(g[1355]))
	m.G1358 = int32(uint32(g[1356]))
	m.G1359 = int32(uint32(g[1357]))
	m.G1360 = int32(uint32(g[1358]))
	m.G1361 = int32(uint32(g[1359]))
	m.G1362 = int32(uint32(g[1360]))
	m.G1363 = int32(uint32(g[1361]))
	m.G1364 = int32(uint32(g[1362]))
	m.G1365 = int32(uint32(g[1363]))
	m.G1366 = int32(uint32(g[1364]))
	m.G1367 = int32(uint32(g[1365]))
	m.G1368 = int32(uint32(g[1366]))
	m.G1369 = int32(uint32(g[1367]))
	m.G1370 = int32(uint32(g[1368]))
	m.G1371 = int32(uint32(g[1369]))
	m.G1372 = int32(uint32(g[1370]))
	m.G1373 = int32(uint32(g[1371]))
	m.G1374 = int32(uint32(g[1372]))
	m.G1375 = int32(uint32(g[1373]))
	m.G1376 = int32(uint32(g[1374]))
	m.G1377 = int32(uint32(g[1375]))
	m.G1378 = int32(uint32(g[1376]))
	m.G1379 = int32(uint32(g[1377]))
	m.G1380 = int32(uint32(g[1378]))
	m.G1381 = int32(uint32(g[1379]))
	m.G1382 = int32(uint32(g[1380]))
	m.G1383 = int32(uint32(g[1381]))
	m.G1384 = int32(uint32(g[1382]))
	m.G1385 = int32(uint32(g[1383]))
	m.G1386 = int32(uint32(g[1384]))
	m.G1387 = int32(uint32(g[1385]))
	m.G1388 = int32(uint32(g[1386]))
	m.G1389 = int32(uint32(g[1387]))
	m.G1390 = int32(uint32(g[1388]))
	m.G1391 = int32(uint32(g[1389]))
	m.G1392 = int32(uint32(g[1390]))
	m.G1393 = int32(uint32(g[1391]))
	m.G1394 = int32(uint32(g[1392]))
	m.G1395 = int32(uint32(g[1393]))
	m.G1396 = int32(uint32(g[1394]))
	m.G1397 = int32(uint32(g[1395]))
	m.G1398 = int32(uint32(g[1396]))
	m.G1399 = int32(uint32(g[1397]))
	m.G1400 = int32(uint32(g[1398]))
	m.G1401 = int32(uint32(g[1399]))
	m.G1402 = int32(uint32(g[1400]))
	m.G1403 = int32(uint32(g[1401]))
	m.G1404 = int32(uint32(g[1402]))
	m.G1405 = int32(uint32(g[1403]))
	m.G1406 = int32(uint32(g[1404]))
	m.G1407 = int32(uint32(g[1405]))
	m.G1408 = int32(uint32(g[1406]))
	m.G1409 = int32(uint32(g[1407]))
	m.G1410 = int32(uint32(g[1408]))
	m.G1411 = int32(uint32(g[1409]))
	m.G1412 = int32(uint32(g[1410]))
	m.G1413 = int32(uint32(g[1411]))
	m.G1414 = int32(uint32(g[1412]))
	m.G1415 = int32(uint32(g[1413]))
	m.G1416 = int32(uint32(g[1414]))
	m.G1417 = int32(uint32(g[1415]))
	m.G1418 = int32(uint32(g[1416]))
	m.G1419 = int32(uint32(g[1417]))
	m.G1420 = int32(uint32(g[1418]))
	m.G1421 = int32(uint32(g[1419]))
	m.G1422 = int32(uint32(g[1420]))
	m.G1423 = int32(uint32(g[1421]))
	m.G1424 = int32(uint32(g[1422]))
	m.G1425 = int32(uint32(g[1423]))
	m.G1426 = int32(uint32(g[1424]))
	m.G1427 = int32(uint32(g[1425]))
	m.G1428 = int32(uint32(g[1426]))
	m.G1429 = int32(uint32(g[1427]))
	m.G1430 = int32(uint32(g[1428]))
	m.G1431 = int32(uint32(g[1429]))
	m.G1432 = int32(uint32(g[1430]))
	m.G1433 = int32(uint32(g[1431]))
	m.G1434 = int32(uint32(g[1432]))
	m.G1435 = int32(uint32(g[1433]))
	m.G1436 = int32(uint32(g[1434]))
	m.G1437 = int32(uint32(g[1435]))
	m.G1438 = int32(uint32(g[1436]))
	m.G1439 = int32(uint32(g[1437]))
	m.G1440 = int32(uint32(g[1438]))
	m.G1441 = int32(uint32(g[1439]))
	m.G1442 = int32(uint32(g[1440]))
	m.G1443 = int32(uint32(g[1441]))
	m.G1444 = int32(uint32(g[1442]))
	m.G1445 = int32(uint32(g[1443]))
	m.G1446 = int32(uint32(g[1444]))
	m.G1447 = int32(uint32(g[1445]))
	m.G1448 = int32(uint32(g[1446]))
	m.G1449 = int32(uint32(g[1447]))
	m.G1450 = int32(uint32(g[1448]))
	m.G1451 = int32(uint32(g[1449]))
	m.G1452 = int32(uint32(g[1450]))
	m.G1453 = int32(uint32(g[1451]))
	m.G1454 = int32(uint32(g[1452]))
	m.G1455 = int32(uint32(g[1453]))
	m.G1456 = int32(uint32(g[1454]))
	m.G1457 = int32(uint32(g[1455]))
	m.G1458 = int32(uint32(g[1456]))
	m.G1459 = int32(uint32(g[1457]))
	m.G1460 = int32(uint32(g[1458]))
	m.G1461 = int32(uint32(g[1459]))
	m.G1462 = int32(uint32(g[1460]))
	m.G1463 = int32(uint32(g[1461]))
	m.G1464 = int32(uint32(g[1462]))
	m.G1465 = int32(uint32(g[1463]))
	m.G1466 = int32(uint32(g[1464]))
	m.G1467 = int32(uint32(g[1465]))
	m.G1468 = int32(uint32(g[1466]))
	m.G1469 = int32(uint32(g[1467]))
	m.G1470 = int32(uint32(g[1468]))
	m.G1471 = int32(uint32(g[1469]))
	m.G1472 = int32(uint32(g[1470]))
	m.G1473 = int32(uint32(g[1471]))
	m.G1474 = int32(uint32(g[1472]))
	m.G1475 = int32(uint32(g[1473]))
	m.G1476 = int32(uint32(g[1474]))
	m.G1477 = int32(uint32(g[1475]))
	m.G1478 = int32(uint32(g[1476]))
	m.G1479 = int32(uint32(g[1477]))
	m.G1480 = int32(uint32(g[1478]))
	m.G1481 = int32(uint32(g[1479]))
	m.G1482 = int32(uint32(g[1480]))
	m.G1483 = int32(uint32(g[1481]))
	m.G1484 = int32(uint32(g[1482]))
	m.G1485 = int32(uint32(g[1483]))
	m.G1486 = int32(uint32(g[1484]))
	m.G1487 = int32(uint32(g[1485]))
	m.G1488 = int32(uint32(g[1486]))
	m.G1489 = int32(uint32(g[1487]))
	m.G1490 = int32(uint32(g[1488]))
	m.G1491 = int32(uint32(g[1489]))
	m.G1492 = int32(uint32(g[1490]))
	m.G1493 = int32(uint32(g[1491]))
	m.G1494 = int32(uint32(g[1492]))
	m.G1495 = int32(uint32(g[1493]))
	m.G1496 = int32(uint32(g[1494]))
	m.G1497 = int32(uint32(g[1495]))
	m.G1498 = int32(uint32(g[1496]))
	m.G1499 = int32(uint32(g[1497]))
	m.G1500 = int32(uint32(g[1498]))
	m.G1501 = int32(uint32(g[1499]))
	m.G1502 = int32(uint32(g[1500]))
	m.G1503 = int32(uint32(g[1501]))
	m.G1504 = int32(uint32(g[1502]))
	m.G1505 = int32(uint32(g[1503]))
	m.G1506 = int32(uint32(g[1504]))
	m.G1507 = int32(uint32(g[1505]))
	m.G1508 = int32(uint32(g[1506]))
	m.G1509 = int32(uint32(g[1507]))
	m.G1510 = int32(uint32(g[1508]))
	m.G1511 = int32(uint32(g[1509]))
	m.G1512 = int32(uint32(g[1510]))
	m.G1513 = int32(uint32(g[1511]))
	m.G1514 = int32(uint32(g[1512]))
	m.G1515 = int32(uint32(g[1513]))
	m.G1516 = int32(uint32(g[1514]))
	m.G1517 = int32(uint32(g[1515]))
	m.G1518 = int32(uint32(g[1516]))
	m.G1519 = int32(uint32(g[1517]))
	m.G1520 = int32(uint32(g[1518]))
	m.G1521 = int32(uint32(g[1519]))
	m.G1522 = int32(uint32(g[1520]))
	m.G1523 = int32(uint32(g[1521]))
	m.G1524 = int32(uint32(g[1522]))
	m.G1525 = int32(uint32(g[1523]))
	m.G1526 = int32(uint32(g[1524]))
	m.G1527 = int32(uint32(g[1525]))
	m.G1528 = int32(uint32(g[1526]))
	m.G1529 = int32(uint32(g[1527]))
	m.G1530 = int32(uint32(g[1528]))
	m.G1531 = int32(uint32(g[1529]))
	m.G1532 = int32(uint32(g[1530]))
	m.G1533 = int32(uint32(g[1531]))
	m.G1534 = int32(uint32(g[1532]))
	m.G1535 = int32(uint32(g[1533]))
	m.G1536 = int32(uint32(g[1534]))
	m.G1537 = int32(uint32(g[1535]))
	m.G1538 = int32(uint32(g[1536]))
	m.G1539 = int32(uint32(g[1537]))
	m.G1540 = int32(uint32(g[1538]))
	m.G1541 = int32(uint32(g[1539]))
	m.G1542 = int32(uint32(g[1540]))
	m.G1543 = int32(uint32(g[1541]))
	m.G1544 = int32(uint32(g[1542]))
	m.G1545 = int32(uint32(g[1543]))
	m.G1546 = int32(uint32(g[1544]))
	m.G1547 = int32(uint32(g[1545]))
	m.G1548 = int32(uint32(g[1546]))
	m.G1549 = int32(uint32(g[1547]))
	m.G1550 = int32(uint32(g[1548]))
	m.G1551 = int32(uint32(g[1549]))
	m.G1552 = int32(uint32(g[1550]))
	m.G1553 = int32(uint32(g[1551]))
	m.G1554 = int32(uint32(g[1552]))
	m.G1555 = int32(uint32(g[1553]))
	m.G1556 = int32(uint32(g[1554]))
	m.G1557 = int32(uint32(g[1555]))
	m.G1558 = int32(uint32(g[1556]))
	m.G1559 = int32(uint32(g[1557]))
	m.G1560 = int32(uint32(g[1558]))
	m.G1561 = int32(uint32(g[1559]))
	m.G1562 = int32(uint32(g[1560]))
	m.G1563 = int32(uint32(g[1561]))
	m.G1564 = int32(uint32(g[1562]))
	m.G1565 = int32(uint32(g[1563]))
	m.G1566 = int32(uint32(g[1564]))
	m.G1567 = int32(uint32(g[1565]))
	m.G1568 = int32(uint32(g[1566]))
	m.G1569 = int32(uint32(g[1567]))
	m.G1570 = int32(uint32(g[1568]))
	m.G1571 = int32(uint32(g[1569]))
	m.G1572 = int32(uint32(g[1570]))
	m.G1573 = int32(uint32(g[1571]))
	m.G1574 = int32(uint32(g[1572]))
	m.G1575 = int32(uint32(g[1573]))
	m.G1576 = int32(uint32(g[1574]))
	m.G1577 = int32(uint32(g[1575]))
	m.G1578 = int32(uint32(g[1576]))
	m.G1579 = int32(uint32(g[1577]))
	m.G1580 = int32(uint32(g[1578]))
	m.G1581 = int32(uint32(g[1579]))
	m.G1582 = int32(uint32(g[1580]))
	m.G1583 = int32(uint32(g[1581]))
	m.G1584 = int32(uint32(g[1582]))
	m.G1585 = int32(uint32(g[1583]))
	m.G1586 = int32(uint32(g[1584]))
	m.G1587 = int32(uint32(g[1585]))
	m.G1588 = int32(uint32(g[1586]))
	m.G1589 = int32(uint32(g[1587]))
	m.G1590 = int32(uint32(g[1588]))
	m.G1591 = int32(uint32(g[1589]))
	m.G1592 = int32(uint32(g[1590]))
	m.G1593 = int32(uint32(g[1591]))
	m.G1594 = int32(uint32(g[1592]))
	m.G1595 = int32(uint32(g[1593]))
	m.G1596 = int32(uint32(g[1594]))
	m.G1597 = int32(uint32(g[1595]))
	m.G1598 = int32(uint32(g[1596]))
	m.G1599 = int32(uint32(g[1597]))
	m.G1600 = int32(uint32(g[1598]))
	m.G1601 = int32(uint32(g[1599]))
	m.G1602 = int32(uint32(g[1600]))
	m.G1603 = int32(uint32(g[1601]))
	m.G1604 = int32(uint32(g[1602]))
	m.G1605 = int32(uint32(g[1603]))
	m.G1606 = int32(uint32(g[1604]))
	m.G1607 = int32(uint32(g[1605]))
	m.G1608 = int32(uint32(g[1606]))
	m.G1609 = int32(uint32(g[1607]))
	m.G1610 = int32(uint32(g[1608]))
	m.G1611 = int32(uint32(g[1609]))
	m.G1612 = int32(uint32(g[1610]))
	m.G1613 = int32(uint32(g[1611]))
	m.G1614 = int32(uint32(g[1612]))
	m.G1615 = int32(uint32(g[1613]))
	m.G1616 = int32(uint32(g[1614]))
	m.G1617 = int32(uint32(g[1615]))
	m.G1618 = int32(uint32(g[1616]))
	m.G1619 = int32(uint32(g[1617]))
	m.G1620 = int32(uint32(g[1618]))
	m.G1621 = int32(uint32(g[1619]))
	m.G1622 = int32(uint32(g[1620]))
	m.G1623 = int32(uint32(g[1621]))
	m.G1624 = int32(uint32(g[1622]))
	m.G1625 = int32(uint32(g[1623]))
	m.G1626 = int32(uint32(g[1624]))
	m.G1627 = int32(uint32(g[1625]))
	m.G1628 = int32(uint32(g[1626]))
	m.G1629 = int32(uint32(g[1627]))
	m.G1630 = int32(uint32(g[1628]))
	m.G1631 = int32(uint32(g[1629]))
	m.G1632 = int32(uint32(g[1630]))
	m.G1633 = int32(uint32(g[1631]))
	m.G1634 = int32(uint32(g[1632]))
	m.G1635 = int32(uint32(g[1633]))
	m.G1636 = int32(uint32(g[1634]))
	m.G1637 = int32(uint32(g[1635]))
	m.G1638 = int32(uint32(g[1636]))
	m.G1639 = int32(uint32(g[1637]))
	m.G1640 = int32(uint32(g[1638]))
	m.G1641 = int32(uint32(g[1639]))
	m.G1642 = int32(uint32(g[1640]))
	m.G1643 = int32(uint32(g[1641]))
	m.G1644 = int32(uint32(g[1642]))
	m.G1645 = int32(uint32(g[1643]))
	m.G1646 = int32(uint32(g[1644]))
	m.G1647 = int32(uint32(g[1645]))
	m.G1648 = int32(uint32(g[1646]))
	m.G1649 = int32(uint32(g[1647]))
	m.G1650 = int32(uint32(g[1648]))
	m.G1651 = int32(uint32(g[1649]))
	m.G1652 = int32(uint32(g[1650]))
	m.G1653 = int32(uint32(g[1651]))
	m.G1654 = int32(uint32(g[1652]))
	m.G1655 = int32(uint32(g[1653]))
	m.G1656 = int32(uint32(g[1654]))
	m.G1657 = int32(uint32(g[1655]))
	m.G1658 = int32(uint32(g[1656]))
	m.G1659 = int32(uint32(g[1657]))
	m.G1660 = int32(uint32(g[1658]))
	m.G1661 = int32(uint32(g[1659]))
	m.G1662 = int32(uint32(g[1660]))
	m.G1663 = int32(uint32(g[1661]))
	m.G1664 = int32(uint32(g[1662]))
	m.G1665 = int32(uint32(g[1663]))
	m.G1666 = int32(uint32(g[1664]))
	m.G1667 = int32(uint32(g[1665]))
	m.G1668 = int32(uint32(g[1666]))
	m.G1669 = int32(uint32(g[1667]))
	m.G1670 = int32(uint32(g[1668]))
	m.G1671 = int32(uint32(g[1669]))
	m.G1672 = int32(uint32(g[1670]))
	m.G1673 = int32(uint32(g[1671]))
	m.G1674 = int32(uint32(g[1672]))
	m.G1675 = int32(uint32(g[1673]))
	m.G1676 = int32(uint32(g[1674]))
	m.G1677 = int32(uint32(g[1675]))
	m.G1678 = int32(uint32(g[1676]))
	m.G1679 = int32(uint32(g[1677]))
	m.G1680 = int32(uint32(g[1678]))
	m.G1681 = int32(uint32(g[1679]))
	m.G1682 = int32(uint32(g[1680]))
	m.G1683 = int32(uint32(g[1681]))
	m.G1684 = int32(uint32(g[1682]))
	m.G1685 = int32(uint32(g[1683]))
	m.G1686 = int32(uint32(g[1684]))
	m.G1687 = int32(uint32(g[1685]))
	m.G1688 = int32(uint32(g[1686]))
	m.G1689 = int32(uint32(g[1687]))
	m.G1690 = int32(uint32(g[1688]))
	m.G1691 = int32(uint32(g[1689]))
	m.G1692 = int32(uint32(g[1690]))
	m.G1693 = int32(uint32(g[1691]))
	m.G1694 = int32(uint32(g[1692]))
	m.G1695 = int32(uint32(g[1693]))
	m.G1696 = int32(uint32(g[1694]))
	m.G1697 = int32(uint32(g[1695]))
	m.G1698 = int32(uint32(g[1696]))
	m.G1699 = int32(uint32(g[1697]))
	m.G1700 = int32(uint32(g[1698]))
	m.G1701 = int32(uint32(g[1699]))
	m.G1702 = int32(uint32(g[1700]))
	m.G1703 = int32(uint32(g[1701]))
	m.G1704 = int32(uint32(g[1702]))
	m.G1705 = int32(uint32(g[1703]))
	m.G1706 = int32(uint32(g[1704]))
	m.G1707 = int32(uint32(g[1705]))
	m.G1708 = int32(uint32(g[1706]))
	m.G1709 = int32(uint32(g[1707]))
	m.G1710 = int32(uint32(g[1708]))
	m.G1711 = int32(uint32(g[1709]))
	m.G1712 = int32(uint32(g[1710]))
	m.G1713 = int32(uint32(g[1711]))
	m.G1714 = int32(uint32(g[1712]))
	m.G1715 = int32(uint32(g[1713]))
	m.G1716 = int32(uint32(g[1714]))
	m.G1717 = int32(uint32(g[1715]))
	m.G1718 = int32(uint32(g[1716]))
	m.G1719 = int32(uint32(g[1717]))
	m.G1720 = int32(uint32(g[1718]))
	m.G1721 = int32(uint32(g[1719]))
	m.G1722 = int32(uint32(g[1720]))
	m.G1723 = int32(uint32(g[1721]))
	m.G1724 = int32(uint32(g[1722]))
	m.G1725 = int32(uint32(g[1723]))
	m.G1726 = int32(uint32(g[1724]))
	m.G1727 = int32(uint32(g[1725]))
	m.G1728 = int32(uint32(g[1726]))
	m.G1729 = int32(uint32(g[1727]))
	m.G1730 = int32(uint32(g[1728]))
	m.G1731 = int32(uint32(g[1729]))
	m.G1732 = int32(uint32(g[1730]))
	m.G1733 = int32(uint32(g[1731]))
	m.G1734 = int32(uint32(g[1732]))
	m.G1735 = int32(uint32(g[1733]))
	m.G1736 = int32(uint32(g[1734]))
	m.G1737 = int32(uint32(g[1735]))
	m.G1738 = int32(uint32(g[1736]))
	m.G1739 = int32(uint32(g[1737]))
	m.G1740 = int32(uint32(g[1738]))
	m.G1741 = int32(uint32(g[1739]))
	m.G1742 = int32(uint32(g[1740]))
	m.G1743 = int32(uint32(g[1741]))
	m.G1744 = int32(uint32(g[1742]))
	m.G1745 = int32(uint32(g[1743]))
	m.G1746 = int32(uint32(g[1744]))
	m.G1747 = int32(uint32(g[1745]))
	m.G1748 = int32(uint32(g[1746]))
	m.G1749 = int32(uint32(g[1747]))
	m.G1750 = int32(uint32(g[1748]))
	m.G1751 = int32(uint32(g[1749]))
	m.G1752 = int32(uint32(g[1750]))
	m.G1753 = int32(uint32(g[1751]))
	m.G1754 = int32(uint32(g[1752]))
	m.G1755 = int32(uint32(g[1753]))
	m.G1756 = int32(uint32(g[1754]))
	m.G1757 = int32(uint32(g[1755]))
	m.G1758 = int32(uint32(g[1756]))
	m.G1759 = int32(uint32(g[1757]))
	m.G1760 = int32(uint32(g[1758]))
	m.G1761 = int32(uint32(g[1759]))
	m.G1762 = int32(uint32(g[1760]))
	m.G1763 = int32(uint32(g[1761]))
	m.G1764 = int32(uint32(g[1762]))
	m.G1765 = int32(uint32(g[1763]))
	m.G1766 = int32(uint32(g[1764]))
	m.G1767 = int32(uint32(g[1765]))
	m.G1768 = int32(uint32(g[1766]))
	m.G1769 = int32(uint32(g[1767]))
	m.G1770 = int32(uint32(g[1768]))
	m.G1771 = int32(uint32(g[1769]))
	m.G1772 = int32(uint32(g[1770]))
	m.G1773 = int32(uint32(g[1771]))
	m.G1774 = int32(uint32(g[1772]))
	m.G1775 = int32(uint32(g[1773]))
	m.G1776 = int32(uint32(g[1774]))
	m.G1777 = int32(uint32(g[1775]))
	m.G1778 = int32(uint32(g[1776]))
	m.G1779 = int32(uint32(g[1777]))
	m.G1780 = int32(uint32(g[1778]))
	m.G1781 = int32(uint32(g[1779]))
	m.G1782 = int32(uint32(g[1780]))
	m.G1783 = int32(uint32(g[1781]))
	m.G1784 = int32(uint32(g[1782]))
	m.G1785 = int32(uint32(g[1783]))
	m.G1786 = int32(uint32(g[1784]))
	m.G1787 = int32(uint32(g[1785]))
	m.G1788 = int32(uint32(g[1786]))
	m.G1789 = int32(uint32(g[1787]))
	m.G1790 = int32(uint32(g[1788]))
	m.G1791 = int32(uint32(g[1789]))
	m.G1792 = int32(uint32(g[1790]))
	m.G1793 = int32(uint32(g[1791]))
	m.G1794 = int32(uint32(g[1792]))
	m.G1795 = int32(uint32(g[1793]))
	m.G1796 = int32(uint32(g[1794]))
	m.G1797 = int32(uint32(g[1795]))
	m.G1798 = int32(uint32(g[1796]))
	m.G1799 = int32(uint32(g[1797]))
	m.G1800 = int32(uint32(g[1798]))
	m.G1801 = int32(uint32(g[1799]))
	m.G1802 = int32(uint32(g[1800]))
	m.G1803 = int32(uint32(g[1801]))
	m.G1804 = int32(uint32(g[1802]))
	m.G1805 = int32(uint32(g[1803]))
	m.G1806 = int32(uint32(g[1804]))
	m.G1807 = int32(uint32(g[1805]))
	m.G1808 = int32(uint32(g[1806]))
	m.G1809 = int32(uint32(g[1807]))
	m.G1810 = int32(uint32(g[1808]))
	m.G1811 = int32(uint32(g[1809]))
	m.G1812 = int32(uint32(g[1810]))
	m.G1813 = int32(uint32(g[1811]))
	m.G1814 = int32(uint32(g[1812]))
	m.G1815 = int32(uint32(g[1813]))
	m.G1816 = int32(uint32(g[1814]))
	m.G1817 = int32(uint32(g[1815]))
	m.G1818 = int32(uint32(g[1816]))
	m.G1819 = int32(uint32(g[1817]))
	m.G1820 = int32(uint32(g[1818]))
	m.G1821 = int32(uint32(g[1819]))
	m.G1822 = int32(uint32(g[1820]))
	m.G1823 = int32(uint32(g[1821]))
	m.G1824 = int32(uint32(g[1822]))
	m.G1825 = int32(uint32(g[1823]))
	m.G1826 = int32(uint32(g[1824]))
	m.G1827 = int32(uint32(g[1825]))
	m.G1828 = int32(uint32(g[1826]))
	m.G1829 = int32(uint32(g[1827]))
	m.G1830 = int32(uint32(g[1828]))
	m.G1831 = int32(uint32(g[1829]))
	m.G1832 = int32(uint32(g[1830]))
	m.G1833 = int32(uint32(g[1831]))
	m.G1834 = int32(uint32(g[1832]))
	m.G1835 = int32(uint32(g[1833]))
	m.G1836 = int32(uint32(g[1834]))
	m.G1837 = int32(uint32(g[1835]))
	m.G1838 = int32(uint32(g[1836]))
	m.G1839 = int32(uint32(g[1837]))
	m.G1840 = int32(uint32(g[1838]))
	m.G1841 = int32(uint32(g[1839]))
	m.G1842 = int32(uint32(g[1840]))
	m.G1843 = int32(uint32(g[1841]))
	m.G1844 = int32(uint32(g[1842]))
	m.G1845 = int32(uint32(g[1843]))
	m.G1846 = int32(uint32(g[1844]))
	m.G1847 = int32(uint32(g[1845]))
	m.G1848 = int32(uint32(g[1846]))
	m.G1849 = int32(uint32(g[1847]))
	m.G1850 = int32(uint32(g[1848]))
	m.G1851 = int32(uint32(g[1849]))
	m.G1852 = int32(uint32(g[1850]))
	m.G1853 = int32(uint32(g[1851]))
	m.G1854 = int32(uint32(g[1852]))
	m.G1855 = int32(uint32(g[1853]))
	m.G1856 = int32(uint32(g[1854]))
	m.G1857 = int32(uint32(g[1855]))
	m.G1858 = int32(uint32(g[1856]))
	m.G1859 = int32(uint32(g[1857]))
	m.G1860 = int32(uint32(g[1858]))
	m.G1861 = int32(uint32(g[1859]))
	m.G1862 = int32(uint32(g[1860]))
	m.G1863 = int32(uint32(g[1861]))
	m.G1864 = int32(uint32(g[1862]))
	m.G1865 = int32(uint32(g[1863]))
	m.G1866 = int32(uint32(g[1864]))
	m.G1867 = int32(uint32(g[1865]))
	m.G1868 = int32(uint32(g[1866]))
	m.G1869 = int32(uint32(g[1867]))
	m.G1870 = int32(uint32(g[1868]))
	m.G1871 = int32(uint32(g[1869]))
	m.G1872 = int32(uint32(g[1870]))
	m.G1873 = int32(uint32(g[1871]))
	m.G1874 = int32(uint32(g[1872]))
	m.G1875 = int32(uint32(g[1873]))
	m.G1876 = int32(uint32(g[1874]))
	m.G1877 = int32(uint32(g[1875]))
	m.G1878 = int32(uint32(g[1876]))
	m.G1879 = int32(uint32(g[1877]))
	m.G1880 = int32(uint32(g[1878]))
	m.G1881 = int32(uint32(g[1879]))
	m.G1882 = int32(uint32(g[1880]))
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
	mem := m.Memory
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
	mem := m.Memory
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
	mem := m.Memory
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
