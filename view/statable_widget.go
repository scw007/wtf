package view

// StatableWidget
//
// A statable widget is one that is able to display information about its stats and settings.
// This view must be used in conjunction with KeyboardWidget (otherwise there's no way for it
// to display onscreen).

import (
	// "fmt"
	// "strings"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/cfg"
	// "github.com/wtfutil/wtf/wtf"
)

// StatableWidget defines the data necessary to make a statable widget
type StatableWidget struct {
	app      *tview.Application
	pages    *tview.Pages
	settings *cfg.Common
	view     *tview.TextView
}

// NewStatableWidget creates and returns an instance of StatableWidget
func NewStatableWidget(app *tview.Application, pages *tview.Pages, settings *cfg.Common) StatableWidget {
	widget := StatableWidget{
		app:      app,
		pages:    pages,
		settings: settings,
	}

	return widget
}

/* -------------------- Exported Functions -------------------- */

// SetView sets the parent view of this widget
func (widget *StatableWidget) SetView(view *tview.TextView) {
	widget.view = view
}

// ShowStatsModal creates, populates, and shows the stats modal
func (widget *StatableWidget) ShowStatsModal() {
	closeFunc := func() {
		widget.pages.RemovePage("statable")
		widget.app.SetFocus(widget.view)
	}

	modal := NewStatableModal(widget, closeFunc)

	widget.pages.AddPage("statable", modal.Frame, false, true)
	widget.app.SetFocus(modal.Frame)

	widget.app.QueueUpdate(func() {
		widget.app.Draw()
	})
}
