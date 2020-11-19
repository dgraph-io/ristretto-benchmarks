package zeroout

import "unsafe"

var zeroBuf []byte

func init() {
	zeroBuf = make([]byte, 1024)
	for i := 0; i < len(zeroBuf); i++ {
		zeroBuf[i] = 0x00
	}
}

// ZeroOut zeroes out all the bytes in the range [start, end).
func ZeroOut(dst []byte, start, end int) {
	if start < 0 || start >= len(dst) {
		return // BAD
	}
	if end >= len(dst) {
		end = len(dst)
	}

	if end-start <= 0 {
		return
	}
	buf := dst[start:end]
	n := copy(buf, zeroBuf)
	if n < len(zeroBuf) {
		return
	}
	for i := n; i < len(buf); i *= 2 {
		copy(buf[i:], buf[:i])
	}
}

//go:linkname zeroout runtime.memclrNoHeapPointers
func zeroout(ptr unsafe.Pointer, n uintptr)

func ZeroOutLN(dst []byte, start, end int) {
	if start < 0 || start >= len(dst) {
		return // BAD
	}
	if end >= len(dst) {
		end = len(dst)
	}
	n := end - start
	if n <= 0 {
		return
	}
	zeroout(unsafe.Pointer(&dst[start]), uintptr(n))
}

func ZeroOutNaive(dst []byte, start, end int) {
	if start < 0 || start >= len(dst) {
		return // BAD
	}
	if end >= len(dst) {
		end = len(dst)
	}
	n := end - start
	if n <= 0 {
		return
	}
	b := dst[start:end]
	for i := range b {
		b[i] = 0
	}
}
