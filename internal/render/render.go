package render

import (
	"fmt"

	"github.com/likeawizard/chess-go/internal/board"
	"github.com/likeawizard/chess-go/internal/config"
	eval "github.com/likeawizard/chess-go/internal/evaluation"
	"github.com/rivo/tview"
)

const (
	performanceFormat = "Time: %.0f. Evaluations per second %.0f\n"
)

type BoardRender interface {
	InitRender(b *board.Board, e *eval.EvalEngine)
	Update(move board.Move)
	Run()
}

type TviewBoardRender struct {
	b           *board.Board
	e           *eval.EvalEngine
	app         *tview.Application
	boardView   *tview.TextView
	currentFen  *tview.TextView
	moves       *tview.TextView
	performance *tview.TextView
	gridLayout  *tview.Grid
}

type SimpleAsciiRender struct {
	b *board.Board
	e *eval.EvalEngine
}

func New(c *config.Config) BoardRender {
	rendererType := c.Render.Mode
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

func (render *TviewBoardRender) InitRender(b *board.Board, e *eval.EvalEngine) {
	render.b = b
	render.e = e
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

func (render *SimpleAsciiRender) InitRender(b *board.Board, e *eval.EvalEngine) {
	render.b = b
	render.e = e
}

func (render *TviewBoardRender) Update(move board.Move) {
	render.boardView.SetText(ASCIIRender(*render.b))
	render.currentFen.SetText(render.b.ExportFEN())
	moveTime := render.e.MoveTime.Seconds()
	evaluationsPerSecond := float64(eval.Evaluations+eval.CachedEvals) / moveTime
	render.performance.SetText(fmt.Sprintf(performanceFormat, moveTime, evaluationsPerSecond))
	// render.moves.SetText(fmt.Sprintf("%v", render.b.GetMoveList()))

	render.app.Draw()

}

func (render *SimpleAsciiRender) Update(move board.Move) {
	fmt.Println(render.b.ExportFEN())

	fmt.Println(move)
	moveTime := render.e.MoveTime.Seconds()
	evaluationsPerSecond := float64(render.e.Evaluations) / moveTime
	fmt.Printf(performanceFormat+"\n", moveTime, evaluationsPerSecond)
}
func (render *TviewBoardRender) Run() {
	render.app.SetRoot(render.gridLayout, true).SetFocus(render.gridLayout).Run()
}

func (render *SimpleAsciiRender) Run() {
	fmt.Scanln()
}

func ASCIIRender(b board.Board) string {
	output := ""
	for r := 7; r >= 0; r-- {
		output += fmt.Sprintf("%d ", r+1)
		for f := 0; f < 8; f++ {
			if b.Coords[f+r*8] == 0 {
				output += "0"
			} else {
				output += fmt.Sprint(board.Pieces[(b.Coords[f+r*8] - 1)])
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
