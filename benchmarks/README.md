<!-- SPDX-License-Identifier: BSD-3-Clause -->
# `go-ruby-stringio` library-level benchmark harness

Reproducible, cross-runtime benchmark of the **pure-Go `go-ruby-stringio` library**
against the reference Ruby runtimes (MRI, MRI + YJIT, JRuby, TruffleRuby). It
measures the **library primitive** through its Go API, isolated from the rbgo
interpreter, so the numbers answer: *is the pure-Go implementation as fast as the
reference runtime's own `StringIO` (a C extension)?*

## Layout

- `go/`            — self-contained Go driver; `go.mod` pins the published library
  by pseudo-version (no `replace`). `go/bench` (the built binary) is git-ignored.
- `ruby/stringio.rb` — the equivalent workload; `ruby/_harness.rb` is the shared
  timer + SHA-256 `check` helper.
- `run.sh`         — verifies the Go output is **byte-identical to MRI** (SHA-256
  of every op) and aborts on mismatch, then runs every available runtime and
  prints one Markdown table per sub-benchmark (ns/op + ratio vs MRI).

## Run

```sh
bash benchmarks/run.sh
```

Environment knobs: `OUTER` (timed passes, default 25), `WARM` (untimed warm-up
passes, default 3), and `RUBY`/`JRUBY`/`TRUFFLERUBY` to select runtime binaries.

## Method

Each process runs `WARM` untimed passes (to let the JVM/GraalVM JITs warm up),
then `OUTER` timed passes of a fixed inner loop, timed with a monotonic clock;
the **best** pass is reported as **ns/op**. Interpreter start-up is outside the
timed region. The Go driver and the Ruby script build **identical inputs** (the
same deterministic multi-line document) and their per-op outputs are checked
identical to MRI (SHA-256) before timing. Results are published, dated, in
`../docs/performance.md`.

## Operations

- **write** — open a write buffer, stream every line through `write`, read it
  back with `string`.
- **gets** — iterate the whole buffer line by line (`gets`/`each_line`).
- **read** — scan the buffer in fixed 64-byte `read(n)` chunks to EOF.
- **getc** — scan the buffer character by character (`getc`/`each_char`).
- **puts** — build a buffer with `puts` + `printf`.
