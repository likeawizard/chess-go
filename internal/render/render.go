package render

import (
	"fmt"
	"os"
	"time"

	"github.com/likeawizard/chess-go/internal/board"
	"github.com/rivo/tview"
)

const (
	performanceFormat = "Time: %.0f. Evaluations per second %.0f (Cached %d Cached percent %f)\n"
)

type BoardRender interface {
	InitRender(b *board.Board, elapsed *time.Duration)
	Update()
	Run()
}

type TviewBoardRender struct {
	b           *board.Board
	elapsed     *time.Duration
	app         *tview.Application
	boardView   *tview.TextView
	currentFen  *tview.TextView
	moves       *tview.TextView
	performance *tview.TextView
	gridLayout  *tview.Grid
}

type SimpleAsciiRender struct {
	b       *board.Board
	elapsed *time.Duration
}

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

func (render *TviewBoardRender) InitRender(b *board.Board, elapsed *time.Duration) {
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

func (render *SimpleAsciiRender) InitRender(b *board.Board, elapsed *time.Duration) {
	render.b = b
	render.elapsed = elapsed
}

func (render *TviewBoardRender) Update() {
	render.boardView.SetText(ASCIIRender(*render.b))
	render.currentFen.SetText(render.b.ExportFEN())
	moveTime := render.elapsed.Seconds()
	evaluationsPerSecond := float64(board.Evaluations+board.CachedEvals) / render.elapsed.Seconds()
	cachedPercentage := float64(board.CachedEvals) / float64(board.Evaluations+board.CachedEvals)
	render.performance.SetText(fmt.Sprintf(performanceFormat, moveTime, evaluationsPerSecond, board.CachedEvals, cachedPercentage))
	render.moves.SetText(fmt.Sprintf("%v", render.b.GetMoveList()))

	render.app.Draw()

}

func (render *SimpleAsciiRender) Update() {
	fmt.Println(render.b.ExportFEN())

	fmt.Println(render.b.GetLastMove())
	moveTime := render.elapsed.Seconds()
	evaluationsPerSecond := float64(board.Evaluations+board.CachedEvals) / render.elapsed.Seconds()
	cachedPercentage := float64(board.CachedEvals) / float64(board.Evaluations+board.CachedEvals)
	fmt.Printf(performanceFormat+"\n", moveTime, evaluationsPerSecond, board.CachedEvals, cachedPercentage)
}
func (render *TviewBoardRender) Run() {
	render.app.SetRoot(render.gridLayout, true).SetFocus(render.gridLayout).Run()
}

func (render *SimpleAsciiRender) Run() {
	fmt.Scanln()
}

func ASCIIRender(b board.Board) string {
	output := ""
	for r := len(b.Coords) - 1; r >= 0; r-- {
		output += fmt.Sprintf("%d ", r+1)
		for f := range b.Coords[r] {
			if b.Coords[f][r] == 0 {
				output += "0"
			} else {
				output += fmt.Sprint(board.Pieces[(b.Coords[f][r] - 1)])
			}
		}
		output += fmt.Sprintln("")
	}
	output += "  "
	for _, file := range board.Files {
		output += fmt.Sprint(file)
	}

	return output
}
