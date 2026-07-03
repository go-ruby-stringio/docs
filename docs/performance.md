# Performance

`go-ruby-stringio/stringio` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's StringIO. This
page records a **comparative benchmark** of that module against the reference
Ruby runtimes, part of the ecosystem-wide per-module parity suite.

## What is measured

The **same** workload — a set of representative StringIO usage patterns (streaming
`write` + `string`, line iteration with `gets`/`each_line`, byte-chunk `read`
scanning, character-by-character `getc` scanning, and a `puts`/`printf` writer) —
is run through the pure-Go library's Go API and through each reference runtime's
own `StringIO`. Ruby's `StringIO` is a **C extension**, so the MRI / YJIT columns
are hand-written C doing the byte shuffling; the go column is
**this pure-Go library doing the work**. The comparison is therefore the
Ruby-visible operation, apples-to-apples across implementations. Every operation's
output is checked **byte-identical to MRI** (SHA-256 of the final buffer / the
reassembled read result) before any timing is trusted — the harness aborts on a
mismatch.

- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.
- **Runtimes:** `ruby` (MRI, the oracle) and `ruby --yjit`; `jruby` (on the JVM);
  `truffleruby` (GraalVM CE Native).

!!! note "rbgo end-to-end row"
    The whole-interpreter `rbgo`-vs-MRI row (single-shot `ruby file.rb` wall time,
    the format used by the other module pages) has not yet been captured for
    StringIO on a controlled host, so it is not shown here rather than printed as
    a fabricated figure. The library-level section below is the real, measured
    parity result for this module.

## Library-level benchmark (Go API vs runtimes) — 2026-07-03

This section measures the **pure-Go library directly, through its Go API** — not
the `rbgo` interpreter path. It isolates the library primitive from
Ruby-interpreter dispatch, answering the parity question head-on: *is the pure-Go
implementation as fast as the reference runtime's own `StringIO` C extension?* The
**same workload, same inputs, same iteration counts** run through the Go library
and through each reference runtime's stdlib; outputs were checked identical to MRI
before any timing.

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5.1 — **date 2026-07-03**.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` · MRI + YJIT · JRuby 10.1.0.0
  (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native).
- **Inputs:** a deterministic 256-line document (`"line %05d the quick brown fox
  jumps %d\n"`); the `puts` writer emits 128 `row`/`printf` pairs; `read` scans in
  64-byte chunks.
- **Method:** 3 untimed warm-up passes, then 25 timed passes of a fixed inner
  loop, monotonic clock, **best** pass reported as **ns/op**.

#### write

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 4487.6 | 0.35× |
| MRI | 12693.3 | 1.00× |
| MRI + YJIT | 8216.7 | 0.65× |
| JRuby | 7701.9 | 0.61× |
| TruffleRuby | 5504.9 | 0.43× |

#### read

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 5155.1 | 0.50× |
| MRI | 10310.0 | 1.00× |
| MRI + YJIT | 8340.0 | 0.81× |
| JRuby | 3451.1 | 0.33× |
| TruffleRuby | 6397.4 | 0.62× |

#### getc

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 51875.1 | 0.11× |
| MRI | 473253.3 | 1.00× |
| MRI + YJIT | 404340.0 | 0.85× |
| JRuby | 136266.0 | 0.29× |
| TruffleRuby | 262840.4 | 0.56× |

#### puts

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 12715.4 | 0.31× |
| MRI | 41253.3 | 1.00× |
| MRI + YJIT | 35060.0 | 0.85× |
| JRuby | 26284.7 | 0.64× |
| TruffleRuby | 16698.1 | 0.40× |

#### gets

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 102233.2 | 5.80× |
| MRI | 17620.0 | 1.00× |
| MRI + YJIT | 17066.7 | 0.97× |
| JRuby | 8284.9 | 0.47× |
| TruffleRuby | 11813.1 | 0.67× |

Mixed, and instructive. On the byte- and buffer-shuffling paths the pure-Go
library **beats MRI's C extension**: `write` 0.35×, `puts` 0.31×, `read` 0.50×,
and most strikingly `getc` **~9× faster** (0.11×) — MRI's per-character `getc`
pays a full method-dispatch + object-allocation round trip in C, where the Go
`EachChar` decode loop stays tight. The one **regression is `gets`/`each_line`
(5.80× slower)**: the library's line core calls `strings.Index(string(rest), sep)`
on the remaining buffer for every line, which reallocates and rescans a shrinking
copy each iteration (≈ O(n²) over the line count), whereas MRI's `gets` is a single
C `memchr` over the live buffer. That line-scan is the concrete optimization target
for this module — indexing on the `[]byte` directly (no `string(rest)` copy) and
advancing a saved offset removes both the allocation and the re-scan. Output stays
byte-identical to MRI throughout.

!!! note "Reproduce"
    The harness is committed under
    [`benchmarks/`](https://github.com/go-ruby-stringio/docs/tree/main/benchmarks):
    a self-contained Go driver (`go/`, pins the published library by
    pseudo-version in `go.mod`, no `replace`), the equivalent `ruby/stringio.rb`
    workload, and `run.sh`. Run `bash benchmarks/run.sh`; it first verifies the Go
    output is byte-identical to MRI (SHA-256 per op) and aborts on mismatch, then
    times. Env `OUTER`/`WARM` tune the pass budget and `RUBY`/`JRUBY`/`TRUFFLERUBY`
    select the runtime binaries.

!!! warning "Warm-up budget & noise — honest framing"
    Numbers reflect a **fixed warm-process budget** (3 warm-up + 25 timed passes
    in one process). The JVM/GraalVM JITs (JRuby, TruffleRuby) may need a larger
    warm-up to reach steady state, so their columns can **understate** peak
    throughput. Sub-microsecond-per-op work is not reached here (the smallest op,
    `write`, is ~4.5 µs), so per-row relative noise is modest, but treat
    differences under ~10% as noise. Every number here is a **real measured value**
    from the dated run above — nothing is fabricated, estimated, or cherry-picked.
    The go-ruby column is the pure-Go library; every other column is that
    interpreter's own `StringIO` (a C extension in MRI) doing the equivalent work.
