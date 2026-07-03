# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
#
# Reference StringIO workload, mirrored byte-for-byte by ../go/main.go. Ruby's
# StringIO is a C extension, so the MRI / YJIT columns measure hand-written C
# doing this byte shuffling; the go column is the pure-Go go-ruby-stringio library.
require "stringio"
require_relative "_harness"

N_LINES   = 256
N_ROWS    = 128
READ_CHUNK = 64

# Deterministic multi-line document shared with the Go driver.
def doc
  (0...N_LINES).map { |i| format("line %05d the quick brown fox jumps %d\n", i, (i * 31 + 7) & 0xFF) }.join
end

DOC   = doc
LINES = DOC.lines           # each "…\n" line, the chunks the writer streams

# --- operations (each returns its canonical output for the equality check) ---

def op_write
  io = StringIO.new(+"", "w")
  LINES.each { |l| io.write(l) }
  io.string
end

def op_gets
  io = StringIO.new(DOC)
  out = +""
  io.each_line("\n") { |line| out << line }
  out
end

def op_read
  io = StringIO.new(DOC)
  out = +""
  while (chunk = io.read(READ_CHUNK))
    out << chunk
  end
  out
end

def op_getc
  io = StringIO.new(DOC)
  out = +""
  while (ch = io.getc)
    out << ch
  end
  out
end

def op_puts
  io = StringIO.new(+"", "w")
  N_ROWS.times do |i|
    io.puts("row #{i}")
    io.printf("=%05d\n", i)
  end
  io.string
end

# Verify digests (outside timing): consumed by verify.sh.
check("write", op_write)
check("gets",  op_gets)
check("read",  op_read)
check("getc",  op_getc)
check("puts",  op_puts)

# Timed passes.
bench("write", 300) { op_write }
bench("gets",  300) { op_gets }
bench("read",  300) { op_read }
bench("getc",  300) { op_getc }
bench("puts",  300) { op_puts }
