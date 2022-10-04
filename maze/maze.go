// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

type Dir byte

const (
	N Dir = 1 << iota
	S
	W
	E
)

/*
	┌───┬───┐
	│   │   │
	├───┼───┤
	│   │   │
	└───┴───┘
*/

func HasDir(e [][]Dir, row, col int, d Dir) bool {
	if len(e) == 0 {
		return false
	}
	if row < 0 || row >= len(e) || col < 0 || col >= len(e[0]) {
		return false
	}
	return e[row][col]&d == d
}

func updateTop(last []rune, e [][]Dir, r int, s map[Coord]tag) {
	hasS := func(x, y int) bool { _, ok := s[Coord{x, y}]; return ok }
	row := e[r]
	hasDir := func(col int, d Dir) bool {
		if col < 0 || col >= len(row) {
			return false
		}
		return row[col]&d == d
	}
	noDir := func(col int, d Dir) bool { return !hasDir(col, d) }
	forCornerVert := func(corner rune, d Dir, hasVert bool) rune {
		switch corner {
		case '└':
			return '├'
		case '┴':
			return '┼'
		case '┘':
			return '┤'
		case '┌':
			if d == E {
				return '┬'
			}
		case '┐':
			if d == W && !hasVert {
				return '┬'
			}
		case '─':
			if d == W && !hasVert {
				return '┌'
			} else if d == E && !hasVert {
				return '┐' // ┌
			}
		case ' ':
			if hasVert {
				return '│'
			}
		}
		return corner
	}

	forCornerHorz := func(corner rune, d Dir, hasHor bool) rune {
		switch corner {
		case ' ':
			if hasHor {
				return '─'
			} else if d == W {
				return '┌'
			} else if d == E {
				return '┐'
			}
		}
		return corner
	}
	for col := range row {
		if noDir(col, W) {
			last[col*4] = forCornerVert(last[col*4], W, hasDir(col, N))
		}
		if noDir(col, E) {
			last[col*4+4] = forCornerVert(last[col*4+4], E, hasDir(col, N))
		}
		if noDir(col, N) {
			last[col*4] = forCornerHorz(last[col*4], W, hasDir(col, W))
			for i := 1; i < 4; i++ {
				last[col*4+i] = '─'
			}
			last[col*4+4] = forCornerHorz(last[col*4+4], E, hasDir(col, E))
		}
		if !hasS(col, r) || noDir(col, N) {
			continue
		}
		if hasS(col, r-1) {
			last[col*4+2] = '┋'
		}
	}
}

func updateMid(last []rune, e [][]Dir, r int, s map[Coord]tag) {
	hasS := func(x, y int) bool { _, ok := s[Coord{x, y}]; return ok }
	row := e[r]
	hasDir := func(col int, d Dir) bool {
		if col < 0 || col >= len(row) {
			return false
		}
		return row[col]&d == d
	}
	noDir := func(col int, d Dir) bool { return !hasDir(col, d) }
	for col := range row {
		if noDir(col, W) {
			last[col*4] = '│'
		}
		if noDir(col, E) {
			last[col*4+4] = '│'
		}
		if !hasS(col, r) {
			continue
		}
		switch {
		case hasS(col, r-1) && hasS(col, r+1) && hasDir(col, N|S):
			last[col*4+2] = '┋'
		case hasS(col-1, r) && hasS(col+1, r) && hasDir(col, E|W):
			last[col*4+0] = '╍'
			last[col*4+1] = '╍'
			last[col*4+2] = '╍'
			last[col*4+3] = '╍'
		case hasS(col-1, r) && hasS(col, r+1) && hasDir(col, W|S):
			last[col*4+0] = '╍'
			last[col*4+1] = '╍'
			last[col*4+2] = '┓'
		case hasS(col+1, r) && hasS(col, r+1) && hasDir(col, E|S):
			last[col*4+2] = '┏'
			last[col*4+3] = '╍'
		case hasS(col, r-1) && hasS(col-1, r) && hasDir(col, N|W):
			last[col*4+0] = '╍'
			last[col*4+1] = '╍'
			last[col*4+2] = '┛'
		case hasS(col, r-1) && hasS(col+1, r) && hasDir(col, N|E):
			last[col*4+2] = '┗'
			last[col*4+3] = '╍'
		}
	}
}

func updateBottom(last []rune, e [][]Dir, r int, s map[Coord]tag) {
	row := e[r]
	hasDir := func(col int, d Dir) bool {
		if col < 0 || col >= len(row) {
			return false
		}
		return row[col]&d == d
	}
	noDir := func(col int, d Dir) bool { return !hasDir(col, d) }
	for col := range row {
		if noDir(col, W) {
			if noDir(col, S) {
				if last[col*4] == '┘' {
					last[col*4] = '┴'
				} else {
					last[col*4] = '└'
				}
				for i := 1; i < 4; i++ {
					last[col*4+i] = '─'
				}
			} else {
				if col > 0 && last[col*4-1] == '─' {
					last[col*4] = '┘'
				} else {
					last[col*4] = '│'
				}
			}
		} else if noDir(col, S) {
			for i := 0; i < 5; i++ {
				last[col*4+i] = '─'
			}
		}
		if noDir(col, E) {
			if noDir(col, S) {
				last[col*4+4] = '┘'
			} else {
				last[col*4+4] = '│'
			}
		}
	}
}

type tag struct{}

func printExits(w io.Writer, e [][]Dir, sol []Coord) {
	var s map[Coord]tag
	if len(sol) > 0 {
		s = make(map[Coord]tag)
		for i := range sol {
			s[sol[i]] = tag{}
		}
	}
	lastRow := make([]rune, 4*len(e[0])+1)
	resetLast := func() {
		for i := range lastRow {
			lastRow[i] = ' '
		}
	}
	outputLast := func() { fmt.Fprintln(w, string(lastRow)) }
	resetLast()
	for row := range e {
		updateTop(lastRow, e, row, s)
		outputLast()
		resetLast()
		updateMid(lastRow, e, row, s)
		outputLast()
		resetLast()
		updateBottom(lastRow, e, row, s)
	}
	outputLast()
}

var (
	allDirs = []Dir{N, W, E, S}
	dirUpd  = map[Dir]func(col int, row int) (ncol, nrow int){
		N: func(c int, r int) (int, int) { return c, r - 1 },
		S: func(c int, r int) (int, int) { return c, r + 1 },
		W: func(c int, r int) (int, int) { return c - 1, r },
		E: func(c int, r int) (int, int) { return c + 1, r },
	}
	dirOpos = map[Dir]Dir{N: S, S: N, W: E, E: W}
)

func shuffledDirs() (dirs [4]Dir) {
	for i, p := range rand.Perm(4) {
		dirs[i] = allDirs[p]
	}
	return dirs
}

func genExits(e [][]Dir, col, row int) {
	for _, d := range shuffledDirs() {
		ncol, nrow := dirUpd[d](col, row)
		if ncol < 0 || ncol >= len(e[0]) || nrow < 0 || nrow >= len(e) {
			continue
		}
		if e[nrow][ncol] != 0 { // && rand.Float64() > 0.01 {
			continue
		}
		e[row][col] |= d
		e[nrow][ncol] |= dirOpos[d]
		genExits(e, ncol, nrow)
	}

}

type Coord struct{ X, Y int }

func (c Coord) Eq(other Coord) bool { return c.X == other.X && c.Y == other.Y }

func solve(e [][]Dir, from, to Coord) (p []Coord) {
	type tag struct{}
	visited := make(map[Coord]tag)

	var v func(Coord) bool
	v = func(c Coord) bool {
		if _, ok := visited[c]; ok {
			return false
		}
		visited[c] = tag{}
		if to.Eq(c) {
			p = append(p, c)
			return true
		}
		for _, d := range shuffledDirs() {
			row, col := c.Y, c.X
			if !HasDir(e, row, col, d) {
				continue
			}
			ncol, nrow := dirUpd[d](col, row)
			nc := Coord{X: ncol, Y: nrow}
			if ncol < 0 || ncol >= len(e[0]) || nrow < 0 || nrow >= len(e) {
				continue
			}
			if v(nc) {
				p = append(p, c)
				return true
			}
		}
		return false
	}
	v(from)
	return p
}

func main() {
	rand.Seed(time.Now().UnixNano())
	printExits(os.Stdout, [][]Dir{
		{S, E, W | E, W | S | E, W},
		{N | E | S, W | E, W | S, N | S, E | S},
		{N | S, 0, N | S, N, N | E | S},
		{N | E, W | E, N | W, E | S, W | N | E | S},
	}, nil)
	const rows, cols = 15, 20
	exits := make([][]Dir, rows)
	for i := range exits {
		exits[i] = make([]Dir, cols)
	}
	genExits(exits, 0, 0)
	printExits(os.Stdout, exits, solve(exits, Coord{X: 0, Y: 0}, Coord{X: cols - 1, Y: 0}))
}
