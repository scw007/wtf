package view

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/wtf"
)

const (
	offscreen   = -1000
	modalWidth  = 80
	modalHeight = 22
)

// StatableModal represents a table-based modal dialog
type StatableModal struct {
	Frame *tview.Frame
	Table *tview.Table
}

// NewStatableModal creates and returns a table-based modal dialog
func NewStatableModal(widget wtf.Wtfable, closeFunc func()) StatableModal {
	keyboardIntercept := buildKeyboardIntercept(closeFunc)

	table := buildTable(widget)
	table.SetInputCapture(keyboardIntercept)

	modal := StatableModal{
		Frame: buildFrame(table),
		Table: table,
	}

	return modal
}

/* -------------------- Helper Functions -------------------- */

func buildFrame(table *tview.Table) *tview.Frame {
	frame := tview.NewFrame(table)
	frame.SetBorder(true)
	frame.SetBorders(1, 1, 0, 0, 1, 1)
	frame.SetRect(offscreen, offscreen, modalWidth, modalHeight)

	drawFunc := func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		w, h := screen.Size()
		frame.SetRect((w/2)-(width/2), (h/2)-(height/2), width, height)
		return x, y, width, height
	}

	frame.SetDrawFunc(drawFunc)

	return frame
}

func buildKeyboardIntercept(closeFunc func()) func(event *tcell.EventKey) *tcell.EventKey {
	keyboardIntercept := func(event *tcell.EventKey) *tcell.EventKey {
		if string(event.Rune()) == "/" {
			closeFunc()
			return nil
		}

		if string(event.Rune()) == "?" {
			closeFunc()
			return nil
		}

		switch event.Key() {
		case tcell.KeyEsc:
			closeFunc()
			return nil
		case tcell.KeyTab:
			return nil
		default:
			return event
		}
	}

	return keyboardIntercept
}

func buildTable(widget wtf.Wtfable) *tview.Table {
	// caption := " [green::b]Settings for " + strings.Title(widget.settings.Module.Type) + "[white]\n\n"

	table := tview.NewTable()
	table.SetBorders(false)

	// Create table here
	table.SetCellSimple(0, 0, "Type:")
	table.SetCellSimple(0, 1, widget.CommonSettings().Module.Type)

	table.SetCellSimple(1, 0, "Refresh:")
	table.SetCellSimple(1, 1, fmt.Sprintf("%d", widget.CommonSettings().RefreshInterval))
	table.SetCellSimple(2, 0, "Refreshed at:")
	table.SetCellSimple(2, 1, widget.RefreshedAt().String())

	table.SetCellSimple(3, 0, "Top:")
	table.SetCellSimple(3, 1, fmt.Sprintf("%d", widget.CommonSettings().PositionSettings.Top))
	table.SetCellSimple(4, 0, "Left:")
	table.SetCellSimple(4, 1, fmt.Sprintf("%d", widget.CommonSettings().PositionSettings.Left))
	table.SetCellSimple(5, 0, "Width:")
	table.SetCellSimple(5, 1, fmt.Sprintf("%d", widget.CommonSettings().PositionSettings.Width))
	table.SetCellSimple(6, 0, "Height:")
	table.SetCellSimple(6, 1, fmt.Sprintf("%d", widget.CommonSettings().PositionSettings.Height))

	for i := 0; i < 6; i++ {
		leftCell := table.GetCell(i, 0)
		leftCell.SetAlign(2)
		leftCell.SetExpansion(1)

		rightCell := table.GetCell(i, 1)
		rightCell.SetExpansion(1)
	}

	return table
}
