# Usage & API

The public API lives at the module root (`github.com/go-ruby-stringio/stringio`). It is **Ruby-shaped but Go-idiomatic**: every method maps to its Ruby StringIO counterpart, while the surface follows Go conventions — methods that Ruby implements as queries return a value, and the ones that can raise return a typed `error` a host maps onto Ruby's exception classes.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-stringio/stringio`, bound into
    `rbgo` as a native module; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-stringio/stringio
```

## Worked example

```go
// Write into an in-memory buffer.
w := stringio.NewString("")
w.Puts("hello", "world")
fmt.Printf("%q\n", w.String()) // "hello\nworld\n"

// Read it back line by line.
r := stringio.NewString(w.String())
for {
    line, ok, _ := r.Gets("\n")
    if !ok {
        break
    }
    fmt.Printf("%q\n", line) // "hello\n" then "world\n"
}

// Seek-past-end writes NUL-pad, exactly like MRI.
g := stringio.NewString("")
g.Write("abc")
g.Seek(6, stringio.SeekSet)
g.Write("z")
fmt.Printf("%q\n", g.String()) // "abc\x00\x00\x00z"
```

## Shape

A `StringIO` is constructed with `New(s, mode)` or `NewString(s)` (the read-write
default).

```go
func New(s, mode string) (*StringIO, error) // StringIO.new(s, mode)
func NewString(s string) *StringIO          // StringIO.new(s) — read-write

// Reading.
func (s *StringIO) Read(n int) (data string, ok bool, err error) // ok=false ⇒ MRI nil
func (s *StringIO) ReadAll() (string, error)
func (s *StringIO) Gets(sep string) (line string, ok bool, err error)
func (s *StringIO) GetsLimit(sep string, limit int) (string, bool, error)
func (s *StringIO) ReadLine(sep string) (string, error)   // EOFError at end
func (s *StringIO) ReadLines(sep string) ([]string, error)
func (s *StringIO) Each(sep string, fn func(line string) error) error
func (s *StringIO) Getc() (ch string, ok bool, err error)
func (s *StringIO) ReadChar() (string, error)             // EOFError at end
func (s *StringIO) Getbyte() (b byte, ok bool, err error)
func (s *StringIO) ReadByte() (byte, error)               // EOFError at end
func (s *StringIO) EachChar(fn func(ch string) error) error
func (s *StringIO) EachByte(fn func(b byte) error) error
func (s *StringIO) Ungetc(ch string) error
func (s *StringIO) Ungetbyte(b byte) error

// Writing.
func (s *StringIO) Write(str string) (int, error)
func (s *StringIO) Puts(args ...string) error
func (s *StringIO) Print(args ...string) error
func (s *StringIO) Printf(format string, args ...any) error
func (s *StringIO) Putc(c byte) (byte, error)
func (s *StringIO) PutString(str string) (string, error)  // IO#putc with a String

// Positioning, content, state.
func (s *StringIO) Pos() int
func (s *StringIO) Tell() int
func (s *StringIO) SetPos(n int) error
func (s *StringIO) Seek(off, whence int) (int, error)     // SeekSet / SeekCur / SeekEnd
func (s *StringIO) Rewind() int
func (s *StringIO) String() string
func (s *StringIO) SetString(str string)
func (s *StringIO) Truncate(n int) (int, error)
func (s *StringIO) Size() int
func (s *StringIO) Length() int
func (s *StringIO) Eof() (bool, error)
func (s *StringIO) Close()
func (s *StringIO) Closed() bool
func (s *StringIO) Flush() *StringIO
func (s *StringIO) Lineno() int
func (s *StringIO) SetLineno(n int)
```

### Errors

The typed errors below model the exceptions MRI's StringIO raises; a host maps
each onto its Ruby counterpart when binding the type into an interpreter.

| Go error            | Ruby exception                          |
| ------------------- | --------------------------------------- |
| `ErrClosed`         | `IOError`, "closed stream"              |
| `ErrNotReadable`    | `IOError`, "not opened for reading"     |
| `ErrNotWritable`    | `IOError`, "not opened for writing"     |
| `ErrEOF`            | `EOFError`, "end of file reached"       |
| `ErrNegativeLength` | `ArgumentError`, "negative length"      |
| `ErrInvalidSeek`    | `Errno::EINVAL`, "Invalid argument"     |

## MRI conformance

Correctness is defined by reference Ruby. A **differential oracle** runs a corpus
of StringIO programs through both the system `ruby` and this library and asserts
the two agree **byte-for-byte**. The oracle scripts `$stdout.binmode` (and
binmode stdin) so Windows text-mode never rewrites the bytes, gate themselves on
`RUBY_VERSION >= "4.0"`, and skip where `ruby` is absent — so the cross-arch and
Windows lanes still validate the library.

## Relationship to Ruby

`go-ruby-stringio/stringio` is **standalone and reusable**, and is the backend bound
into [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby) by `rbgo` as a
native module — the same way [go-ruby-regexp](https://github.com/go-ruby-regexp)
and [go-ruby-erb](https://github.com/go-ruby-erb) are bound. Binding the type into
a live Ruby object model — wiring `$stdout = StringIO.new` or routing
`Kernel#puts` through it — is the host's job; this library hands back an idiomatic
Go `StringIO`. The dependency runs the other way: this library has no dependency
on the Ruby runtime.
