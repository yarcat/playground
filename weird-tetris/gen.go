// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"

	"github.com/fatih/color"
)

type (
	figure       []string
	figureStates []figure
	glass        [9][4]byte
)

var (
	three = figureStates{
		{
			" * ",
			"***",
			"   "},
		{
			" * ",
			" **",
			" * "},
		{
			"***",
			" * ",
			"   "},
		{
			" * ",
			"** ",
			" * "}}
	Z = figureStates{
		{
			"** ",
			" * ",
			" **"},
		{
			"  *",
			"***",
			"*  "},
		{
			"** ",
			" * ",
			" **"},
		{
			"  *",
			"***",
			"*  "}}

	L = figureStates{
		{
			"*  ",
			"*  ",
			"** "},
		{
			"   ",
			"***",
			"*  "},
		{
			"** ",
			" * ",
			" * "},
		{
			"  *",
			"***",
			"   "}}

	cube = figureStates{
		{
			"** ",
			"** ",
			"   "},
		{
			"** ",
			"** ",
			"   "},
		{
			"** ",
			"** ",
			"   "},
		{
			"** ",
			"** ",
			"   "}}

	v = figureStates{
		{
			"*  ",
			"** ",
			"   "},
		{
			"** ",
			"*  ",
			"   "},
		{
			"** ",
			" * ",
			"   "},
		{
			" * ",
			"** ",
			"   "}}

	x = figureStates{
		{
			" * ",
			"***",
			" * "},
		{
			" * ",
			"***",
			" * "},
		{
			" * ",
			"***",
			" * "},
		{
			" * ",
			"***",
			" * "}}

	horns = figureStates{
		{
			"* *",
			"***",
			"   "},
		{
			"** ",
			"*  ",
			"** "},
		{
			"***",
			"* *",
			"   "},
		{
			" **",
			"  *",
			" **"}}

	pipe = figureStates{
		{
			"***",
			"   ",
			"   "},
		{
			"*  ",
			"*  ",
			"*  "},
		{
			"***",
			"   ",
			"   "},
		{
			"*  ",
			"*  ",
			"*  "}}

	figures = []figureStates{
		three,
		Z,
		L,
		cube,
		v,
		v,
		x,
		horns,
		pipe,
	}
)

func generateDrops(idx []int, f func(glass)) {
	var g glass
	genDrops(&g, idx, len(idx)-1, f)
}

func genDrops(g *glass, idx []int, i int, f func(glass)) {
	if i < 0 {
		f(*g)
		return
	}
	for rot := range 4 {
		for col := range colType(4) {
			row := rowType(-1)
			for r := range rowType(9) {
				if tryDrop(g, figures[idx[i]][rot], 8-r, colType(col)) {
					row = 8 - r
					break
				}
			}
			if row < 0 {
				continue
			}
			n := idx[i]
			setDrop(g, figures[n][rot], row, col, byte(n+1))
			genDrops(g, idx, i-1, f)
			setDrop(g, figures[n][rot], row, col, 0)
		}
	}
}

type (
	rowType int
	colType int
)

func tryDrop(g *glass, f figure, row rowType, col colType) bool {
	for r := range f {
		for c := range f[r] {
			if f[r][c] != ' ' {
				rr, cc := int(row)+r, int(col)+c
				if rr >= 9 || cc >= 4 || g[rr][cc] != 0 {
					return false
				}
			}
		}
	}
	return true
}

func setDrop(g *glass, f figure, row rowType, col colType, x byte) {
	for r := range f {
		for c := range f[r] {
			if f[r][c] != ' ' {
				rr, cc := int(row)+r, int(col)+c
				g[rr][cc] = x
			}
		}
	}
}

func generate(f func(glass)) {
	n := 1
	for i := range len(figures) {
		n *= i + 1
	}

	gen(len(figures), func(idx []int) {
		generateDrops(idx, f)
		n--
		fmt.Println("*********", n)
	})

}

func gen(n int, f func([]int)) {
	v := make([]int, n)
	for i := range v {
		v[i] = i
	}
	genImpl(n-1, v, f)
}

func genImpl(n int, v []int, f func([]int)) {
	if n == 0 {
		f(v)
	} else {
		genImpl(n-1, v, f)
		for k := n - 1; k >= 0; k-- {
			v[n], v[k] = v[k], v[n]
			genImpl(n-1, v, f)
			v[n], v[k] = v[k], v[n]
		}
	}
}

func pr(a color.Attribute) func(rune) {
	return func(r rune) {
		color.Set(a)
		fmt.Print(string(r))
	}
}

var colors = [9]func(rune){
	pr(color.FgCyan),
	pr(color.FgRed),
	pr(color.FgGreen),
	pr(color.FgYellow),
	pr(color.FgBlue),
	pr(color.FgMagenta),
	pr(color.FgCyan),
	pr(color.FgGreen),
	pr(color.FgRed),
}

func printGlass(g glass) {
	for _, row := range g {
		for _, c := range row {
			if c > 0 {
				colors[c-1]('\u2588')
				colors[c-1]('\u2588')
			} else {
				fmt.Print(".")
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	gen := make(map[glass]struct{})
	generate(func(g glass) {
		if _, ok := gen[g]; ok {
			return
		}
		gen[g] = struct{}{}
		printGlass(g)
	})
	//var g glass
	//setDrop(&g, figures[0][0], 4, 1, 9)
	//printGlass(g)
}
