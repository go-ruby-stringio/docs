# Why pure Go

`go-ruby-stringio/stringio` reimplements Ruby's StringIO in **pure Go, with cgo disabled**. The
slice of Ruby it covers is **deterministic and interpreter-independent**: an
in-memory IO over a String buffer is pure compute — the cursor arithmetic, the
mode gating, the line/character/byte iteration, and the NUL-padding extension on
a seek-past-end write are all a pure function of their inputs, with no live
binding and no evaluation of arbitrary Ruby. That is exactly the part that can —
and should — live as a standalone Go library, separate from the interpreter.

## Extracted from rbgo, reusable by anyone

This library began life inside [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby)'s
`rbgo`. It has been **made a reusable standalone library** so that:

- any Go program can import `github.com/go-ruby-stringio/stringio` directly, with no Ruby runtime;
- the dependency runs the *other* way — `rbgo` binds this module as a native
  module (the same pattern as [go-ruby-regexp](https://github.com/go-ruby-regexp)
  and [go-ruby-erb](https://github.com/go-ruby-erb)), rather than this module
  depending on the interpreter;
- the behaviour is pinned by a **differential oracle** against the system
  `ruby`, independent of any one consumer.

Binding the type into a live Ruby object model — wiring `$stdout = StringIO.new`
or routing `Kernel#puts` through it — is the host's job; this library hands back
an idiomatic Go `StringIO` whose typed errors the host maps onto Ruby's exception
classes.

## Why pure Go matters here

Because the library is CGO-free and dependency-free, it:

- cross-compiles to every Go target with no C toolchain, and links into a single
  static binary;
- has **no dependency on the Ruby runtime** — the dependency runs the other way;
- can be differentially tested against the `ruby` binary wherever one is on
  `PATH`, while the cross-arch and Windows lanes (where `ruby` is absent) still
  validate the library itself.

See [Usage & API](api.md) for the surface and [Roadmap](roadmap.md) for what is
in scope.
