package main

import "sync"

type TailBuffer struct {
	lines []string
	mu    sync.Mutex
	max   int
}

func NewTailBuffer(n int) *TailBuffer {
	return &TailBuffer{
		lines: make([]string, 0, n),
		max:   n,
	}
}

func (t *TailBuffer) Add(line string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.lines) >= t.max {
		t.lines = t.lines[1:]
	}
	t.lines = append(t.lines, line)
}

func (t *TailBuffer) Get() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]string, len(t.lines))
	copy(out, t.lines)
	return out
}
