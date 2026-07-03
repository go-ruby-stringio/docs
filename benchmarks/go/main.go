// SPDX-License-Identifier: BSD-3-Clause
//
// Library-level workload for go-ruby-stringio/stringio, mirrored byte-for-byte by
// ruby/stringio.rb. Each sub-benchmark exercises one representative StringIO usage
// pattern; check() emits the canonical output digest so verify.sh can prove the
// pure-Go result equals MRI's before the numbers are trusted.
package main

import (
	"fmt"
	"strings"

	"github.com/go-ruby-stringio/stringio"
)

const nLines = 256    // lines in the read/scan document
const nRows = 128     // rows produced by the puts/printf writer
const readChunk = 64  // bytes per read() in the byte-scan benchmark

// doc builds the deterministic multi-line buffer shared with the Ruby side:
//   "line %05d the quick brown fox jumps %d\n", i, (i*31+7)&0xFF
func doc() string {
	var b strings.Builder
	for i := 0; i < nLines; i++ {
		fmt.Fprintf(&b, "line %05d the quick brown fox jumps %d\n", i, (i*31+7)&0xFF)
	}
	return b.String()
}

// lines returns doc split into its individual "…\n" lines, the chunks the writer
// benchmark streams into a fresh StringIO.
func lines() []string {
	return strings.SplitAfter(strings.TrimSuffix(doc(), "\n")+"\n", "\n")[:nLines]
}

// --- operations (each returns its canonical output for the equality check) ---

// opWrite: open a write buffer, stream every line through Write, read it back with
// String — the "write N chunks then read back" pattern.
func opWrite(ls []string) string {
	io, _ := stringio.New("", "w")
	for _, l := range ls {
		io.Write(l)
	}
	return io.String()
}

// opGets: iterate the whole buffer line by line (Gets/each_line), reassembling the
// lines to prove the split is faithful.
func opGets(d string) string {
	io := stringio.NewString(d)
	var b strings.Builder
	io.Each("\n", func(line string) error {
		b.WriteString(line)
		return nil
	})
	return b.String()
}

// opRead: scan the buffer in fixed-size Read(n) chunks to EOF.
func opRead(d string) string {
	io := stringio.NewString(d)
	var b strings.Builder
	for {
		chunk, ok, _ := io.Read(readChunk)
		if !ok {
			break
		}
		b.WriteString(chunk)
	}
	return b.String()
}

// opGetc: scan the buffer character by character (Getc/each_char).
func opGetc(d string) string {
	io := stringio.NewString(d)
	var b strings.Builder
	io.EachChar(func(ch string) error {
		b.WriteString(ch)
		return nil
	})
	return b.String()
}

// opPuts: build a buffer with Puts + Printf, the line-oriented writer pattern.
func opPuts() string {
	io, _ := stringio.New("", "w")
	for i := 0; i < nRows; i++ {
		io.Puts(fmt.Sprintf("row %d", i))
		io.Printf("=%05d\n", i)
	}
	return io.String()
}

func main() {
	d := doc()
	ls := lines()

	// Verify digests (outside timing): consumed by verify.sh.
	check("write", opWrite(ls))
	check("gets", opGets(d))
	check("read", opRead(d))
	check("getc", opGetc(d))
	check("puts", opPuts())

	// Timed passes.
	bench("write", 300, func() { sink = opWrite(ls) })
	bench("gets", 300, func() { sink = opGets(d) })
	bench("read", 300, func() { sink = opRead(d) })
	bench("getc", 300, func() { sink = opGetc(d) })
	bench("puts", 300, func() { sink = opPuts() })
}
