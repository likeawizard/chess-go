package board

import (
	"fmt"
	"github.com/rivo/tview"
	"time"
)

type BoardRender struct {
	b           *Board
	elapsed     *time.Duration
	app         *tview.Application
	boardView   *tview.TextView
	currentFen  *tview.TextView
	moves       *tview.TextView
	performance *tview.TextView
	gridLayout  *tview.Grid
}

func New() *BoardRender {
	return &BoardRender{}
}

func (render *BoardRender) InitRender(b *Board, elapsed *time.Duration) *BoardRender {
	render.b = b
	render.elapsed = elapsed
	render.app = tview.NewApplication()
	render.boardView = tview.NewTextView()
	render.currentFen = tview.NewTextView()
	render.moves = tview.NewTextView()
	render.performance = tview.NewTextView()

	render.gridLayout = tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(60, 0).
		SetBorders(true).
		AddItem(render.currentFen, 0, 0, 1, 2, 0, 0, false).
		AddItem(render.performance, 2, 0, 1, 2, 0, 0, false)

	render.gridLayout.AddItem(render.boardView, 1, 0, 1, 1, 0, 100, false).
		AddItem(render.moves, 1, 1, 1, 1, 0, 100, false)

	return render
}

func (render *BoardRender) Update() *BoardRender {
	render.boardView.SetText(render.b.ASCIIRender())
	render.currentFen.SetText(render.b.ExportFEN())
	render.performance.SetText(fmt.Sprintf("Time: %.0f. Evaluations per second %.0f (Cached percent %f)\n", render.elapsed.Seconds(), float64(Evaluations+CachedEvals)/render.elapsed.Seconds(), float64(CachedEvals)/float64(Evaluations+CachedEvals)))
	render.moves.SetText(fmt.Sprintf("%v", render.b.GetMoveList()))

	render.app.Draw()

	return render
}

func (render *BoardRender) Run() {
	render.app.SetRoot(render.gridLayout, true).SetFocus(render.gridLayout).Run()
}

func (b Board) ASCIIRender() string {
	output := ""
	for r := len(b.coords) - 1; r >= 0; r-- {
		output += fmt.Sprintf("%d ", r+1)
		for f := range b.coords[r] {
			if b.coords[f][r] == 0 {
				output += fmt.Sprint("0")
			} else {
				output += fmt.Sprint(Pieces[(b.coords[f][r] - 1)])
			}
		}
		output += fmt.Sprintln("")
	}
	output += fmt.Sprint("  ")
	for _, file := range Files {
		output += fmt.Sprint(file)
	}

	return output
}
