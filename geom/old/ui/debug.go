package ui

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func debugPrint(screen *ebiten.Image) {
	const debugTmpl = `
FPS : %.2f
TPS : %.2f
`
	msg := fmt.Sprintf(
		debugTmpl,
		ebiten.CurrentFPS(),
		ebiten.CurrentTPS(),
	)
	ebitenutil.DebugPrint(screen, strings.TrimSpace(msg))
}
