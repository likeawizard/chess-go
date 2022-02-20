package board

import (
	"fmt"
	"github.com/rivo/tview"
	"os"
	"time"
)

const (
	performanceFormat = "Time: %.0f. Evaluations per second %.0f (Cached %d Cached percent %f)\n"
)

type BoardRender interface {
	InitRender(b *Board, elapsed *time.Duration)
	Update()
	Run()
}

type TviewBoardRender struct {
	b           *Board
	elapsed     *time.Duration
	app         *tview.Application
	boardView   *tview.TextView
	currentFen  *tview.TextView
	moves       *tview.TextView
	performance *tview.TextView
	gridLayout  *tview.Grid
}

type SimpleAsciiRender struct {
	b       *Board
	elapsed *time.Duration
}

var i BoardRender = &SimpleAsciiRender{}

func New() BoardRender {
	rendererType := os.Getenv("BOARD_RENDER")
	var renderer BoardRender
	switch rendererType {
	case "tview":
		renderer = &TviewBoardRender{}
	case "simple":
		fallthrough
	default:
		renderer = &SimpleAsciiRender{}
	}

	return renderer
}

func (render *TviewBoardRender) InitRender(b *Board, elapsed *time.Duration) {
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
}

func (render *SimpleAsciiRender) InitRender(b *Board, elapsed *time.Duration) {
	render.b = b
	render.elapsed = elapsed
}

func (render *TviewBoardRender) Update() {
	render.boardView.SetText(render.b.ASCIIRender())
	render.currentFen.SetText(render.b.ExportFEN())
	render.performance.SetText(fmt.Sprintf(performanceFormat, render.elapsed.Seconds(), float64(Evaluations+CachedEvals)/render.elapsed.Seconds(), CachedEvals, float64(CachedEvals)/float64(Evaluations+CachedEvals)))
	render.moves.SetText(fmt.Sprintf("%v", render.b.GetMoveList()))

	render.app.Draw()

}

func (render *SimpleAsciiRender) Update() {
	fmt.Println(render.b.ExportFEN())

	fmt.Println(render.b.GetMoveList()[:1])
	fmt.Printf(performanceFormat+"\n", render.elapsed.Seconds(), float64(Evaluations+CachedEvals)/render.elapsed.Seconds(), CachedEvals, float64(CachedEvals)/float64(Evaluations+CachedEvals))
}
func (render *TviewBoardRender) Run() {
	render.app.SetRoot(render.gridLayout, true).SetFocus(render.gridLayout).Run()
}

func (render *SimpleAsciiRender) Run() {
	fmt.Scanln()
	return
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
