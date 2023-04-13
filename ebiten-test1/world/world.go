package world

import (
	"bufio"
	"bytes"
	"log"
	"strings"
)

type World struct {
	data [][]byte
	cols int
}

func (w *World) Dims() (rows, cols int) { return len(w.data), w.cols }
func (w *World) EmptyAt(row, col int) bool {
	if row < 0 || row >= len(w.data) {
		return false
	}
	if col < 0 || col >= len(w.data[row]) {
		return false
	}
	return w.data[row][col] == '.'
}

func MustGet() (w World) {
	scanner := bufio.NewScanner(strings.NewReader(`
		++++++++++
		+.+....+.+
		+.+....+.+
		+.+++..+.+
		+.+....+.+
		+.+....+.+
		+........+
		+.+++.++++
		+.+......+
		++++++++++
	`))
	for scanner.Scan() {
		b := bytes.TrimSpace(scanner.Bytes())
		if lb := len(b); lb == 0 {
			continue

		} else if lb > w.cols {
			w.cols = lb
		}
		w.data = append(w.data, b)

	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("unable to parse the world: %v", err)
	}
	return w
}
