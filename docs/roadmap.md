# Roadmap

`go-ruby-stringio/stringio` is grown **test-first**, each capability differential-tested against MRI
rather than built in isolation. Ruby's StringIO — the
deterministic, interpreter-independent in-memory IO — is
**complete**.

| Stage | What | Status |
| --- | --- | --- |
| Modes | `r` / `w` / `a` and the `r+` / `w+` / `a+` read-write variants, with `w`/`w+` truncating the seed and append mode writing at the end regardless of the cursor. | **Done** |
| Reading | `read(n)` / `read`, `gets` with separator / byte limit / paragraph mode, `readline` / `readlines` / `each_line`, `getc` / `each_char` / `readchar`, `getbyte` / `each_byte` / `readbyte`. | **Done** |
| Writing | `write` / `<<`, `puts` / `print` / `printf`, `putc`, and seek-past-end writes that extend and NUL-pad the buffer exactly like MRI. | **Done** |
| Positioning | `pos` / `pos=` / `tell`, `seek` (SEEK_SET / CUR / END), `rewind`; a negative position raises `Errno::EINVAL`. | **Done** |
| Content & state | `string` / `string=`, `truncate`, `size` / `length`, `eof?`, `close` / `closed?`, `flush`, `lineno` / `lineno=`, `ungetc` / `ungetbyte`. | **Done** |
| Exact MRI raises & coverage | `IOError` on closed / wrong-mode streams, `EOFError` on `readline`/`readchar`/`readbyte` at EOF, `ArgumentError` on a negative `read`; a StringIO program corpus run by `ruby` reproduced byte-for-byte. 100% coverage, green across six arches and three OSes. | **Done** |

## Documented out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **No interpreter.** The library implements the deterministic in-memory IO; it never runs arbitrary Ruby. Wiring `$stdout = StringIO.new` or routing `Kernel#puts` through it is the consumer's job — that is why `rbgo` binds this module rather than the reverse.
- **Reference is reference Ruby (MRI).** Byte-for-byte conformance targets MRI's StringIO behaviour, pinned by the differential oracle.
- **Standalone & reusable.** The module has no dependency on the Ruby runtime; the dependency runs the other way.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
deterministic/interpreter split.
