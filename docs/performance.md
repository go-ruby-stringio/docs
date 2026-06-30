# Performance

`go-ruby-stringio/stringio` is the pure-Go library that
[`rbgo`](https://github.com/go-embedded-ruby/ruby) binds for Ruby's StringIO. This
page records the **methodology** for a comparative benchmark of that module
against the reference Ruby runtimes, part of the ecosystem-wide per-module
parity suite.

## What is measured

The **same** Ruby script — a StringIO write/read round-trip over a representative buffer — is run under every runtime. `rbgo`'s
number reflects **this pure-Go library doing the work**; every other column is
that interpreter's own stdlib. So the comparison is the **Ruby-visible
operation**, apples-to-apples across interpreters. The script prints a
deterministic checksum and its output is checked **byte-identical to MRI** before
timing.

- **Method:** best-of-5 wall time (best, not mean, to suppress scheduler noise);
  single-shot processes, no warm-up beyond the script's own loop.
- **Runtimes:** `ruby` (MRI, the oracle) and `ruby --yjit`; `jruby` (OpenJDK);
  `truffleruby` (GraalVM CE Native).
- The benchmark script and harness live in rbgo's repo under
  [`bench/modules/`](https://github.com/go-embedded-ruby/ruby/tree/main/bench/modules)
  (`stringio.rb` + `run.sh`). Reproduce:
  `RBGO=./rbgo TRUFFLE=truffleruby bash bench/modules/run.sh 5`.

## Result

!!! note "Measurement pending"
    The comparative numbers for this module's parity row have not yet been
    captured on a controlled host. Rather than print fabricated figures, this
    page documents the methodology; the table will be filled from a real
    best-of-5 run (MRI / YJIT / JRuby / TruffleRuby vs `rbgo` on
    go-ruby-stringio) once recorded. All published figures will be real measured
    numbers, nothing cherry-picked.

| Runtime | time | vs MRI |
| --- | ---: | ---: |
| **rbgo** (go-ruby-stringio) | _pending_ | _pending_ |
| MRI | _pending_ | 1.00× |
| MRI + YJIT | _pending_ | _pending_ |
| JRuby | _pending_ | _pending_ |
| TruffleRuby | _pending_ | _pending_ |

!!! note "Honest framing"
    JRuby and TruffleRuby will be timed **cold, single-shot**, so they carry JVM /
    Graal startup on every run — read them as one-shot `ruby file.rb` costs, the
    same way `rbgo` and MRI are measured, not as steady-state JIT numbers. Rows
    that complete in well under ~200 ms carry the most relative noise; treat their
    ratios as order-of-magnitude.
