# go-ruby-stringio documentation

**Ruby's StringIO — an in-memory IO over a String buffer — in pure Go, MRI-compatible, no cgo.**

`go-ruby-stringio/stringio` is a faithful, pure-Go (zero cgo) reimplementation of Ruby's StringIO in-memory IO,
matching reference Ruby (MRI) byte-for-byte. The module path is
`github.com/go-ruby-stringio/stringio`.

It is a **standalone, reusable library**: the module is importable by any Go
program, and it is the backend bound into
[go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo` as a
native module — just like [go-ruby-regexp](https://github.com/go-ruby-regexp)
and [go-ruby-erb](https://github.com/go-ruby-erb). The dependency runs the other
way: this library has **no dependency on the Ruby runtime**.

!!! success "Status: complete — MRI byte-exact"
    A faithful port of Ruby's StringIO: a read/write cursor over a String buffer with **mode gating**, **`read`** / **`gets`** (separator, byte limit, paragraph mode), character/byte iteration, **`write`** / **`puts`** / **`printf`** / **`putc`**, **positioning** (`pos`/`seek`/`rewind`), content ops, **`ungetc`** / **`ungetbyte`**, **seek-past-end NUL padding**, and the exact `IOError` / `EOFError` / `ArgumentError` raises. Validated by a **differential oracle** against the system `ruby` at 100% coverage, `gofmt` + `go vet` clean, CI green across the six 64-bit Go targets and three OSes.

## Quick taste

```go
// Write into an in-memory buffer.
w := stringio.NewString("")
w.Puts("hello", "world")
fmt.Printf("%q\n", w.String()) // "hello\nworld\n"

// Seek-past-end writes NUL-pad, exactly like MRI.
g := stringio.NewString("")
g.Write("abc")
g.Seek(6, stringio.SeekSet)
g.Write("z")
fmt.Printf("%q\n", g.String()) // "abc\x00\x00\x00z"
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`stringio`](https://github.com/go-ruby-stringio/stringio) | the library — Ruby's StringIO in-memory IO in pure Go |
| [`docs`](https://github.com/go-ruby-stringio/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-stringio.github.io`](https://github.com/go-ruby-stringio/go-ruby-stringio.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-stringio/brand) | logo and brand assets |


## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- **MRI byte-exact.** Output matches reference Ruby exactly, not approximately,
  validated by a differential oracle against the `ruby` binary.
- **Standalone & reusable.** No dependency on the Ruby runtime — the dependency
  runs the other way.
- **100% test coverage** is the target, enforced as a CI gate.

## Where to go next

- [Why pure Go](why.md) — why this slice of Ruby is deterministic enough to live
  as a standalone, interpreter-independent Go library.
- [Usage & API](api.md) — the public surface and worked examples.
- [Roadmap](roadmap.md) — what is done and what is downstream by design.

Source lives at [github.com/go-ruby-stringio/stringio](https://github.com/go-ruby-stringio/stringio).
