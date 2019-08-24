package view

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/cfg"
	"github.com/wtfutil/wtf/wtf"
)

type helpItem struct {
	Key  string
	Text string
}

// KeyboardWidget manages keyboard control for a widget
type KeyboardWidget struct {
	app      *tview.Application
	pages    *tview.Pages
	view     *tview.TextView
	settings *cfg.Common

	charMap  map[string]func()
	keyMap   map[tcell.Key]func()
	charHelp []helpItem
	keyHelp  []helpItem
	maxKey   int
}

// NewKeyboardWidget creates and returns a new instance of KeyboardWidget
func NewKeyboardWidget(app *tview.Application, pages *tview.Pages, settings *cfg.Common) KeyboardWidget {
	keyWidget := KeyboardWidget{
		app:      app,
		pages:    pages,
		settings: settings,
		charMap:  make(map[string]func()),
		keyMap:   make(map[tcell.Key]func()),
		charHelp: []helpItem{},
		keyHelp:  []helpItem{},
	}

	return keyWidget
}

// SetKeyboardChar sets a character/function combination that responds to key presses
// Example:
//
//    widget.SetKeyboardChar("d", widget.deleteSelectedItem)
//
func (widget *KeyboardWidget) SetKeyboardChar(char string, fn func(), helpText string) {
	if char == "" {
		return
	}

	widget.charMap[char] = fn
	widget.charHelp = append(widget.charHelp, helpItem{char, helpText})
}

// SetKeyboardKey sets a tcell.Key/function combination that responds to key presses
// Example:
//
//    widget.SetKeyboardKey(tcell.KeyCtrlD, widget.deleteSelectedItem)
//
func (widget *KeyboardWidget) SetKeyboardKey(key tcell.Key, fn func(), helpText string) {
	widget.keyMap[key] = fn
	widget.keyHelp = append(widget.keyHelp, helpItem{tcell.KeyNames[key], helpText})

	if len(tcell.KeyNames[key]) > widget.maxKey {
		widget.maxKey = len(tcell.KeyNames[key])
	}
}

// InitializeCommonKeys sets up the keyboard controls that are common to
// all widgets that accept keyboard input
func (widget *KeyboardWidget) InitializeCommonKeys(refreshFunc func()) {
	// Opens a modal dialog that displays the available keyboard controls for the module
	widget.SetKeyboardChar("/", widget.showHelp, "Show/hide this help prompt")

	if refreshFunc != nil {
		// Calls the refresh function on the module to update its data
		widget.SetKeyboardChar("r", refreshFunc, "Refresh widget")
	}

	// Opens a modal dialog that displays the settings and stats for the module
	widget.SetKeyboardChar("?", widget.showSettings, "Show settings and stats for this widget")
}

// InputCapture is the function passed to tview's SetInputCapture() function
// This is done during the main widget's creation process using the following code:
//
//    widget.View.SetInputCapture(widget.InputCapture)
//
func (widget *KeyboardWidget) InputCapture(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return nil
	}

	fn := widget.charMap[string(event.Rune())]
	if fn != nil {
		fn()
		return nil
	}

	fn = widget.keyMap[event.Key()]
	if fn != nil {
		fn()
		return nil
	}

	return event
}

// HelpText returns the help text and keyboard command info for this widget
func (widget *KeyboardWidget) HelpText() string {
	str := " [green::b]Keyboard commands for " + widget.moduleName() + "[white]\n\n"

	for _, item := range widget.charHelp {
		str += fmt.Sprintf("  %s\t%s\n", item.Key, item.Text)
	}
	str += "\n\n"

	for _, item := range widget.keyHelp {
		str += fmt.Sprintf("  %-*s\t%s\n", widget.maxKey, item.Key, item.Text)
	}

	return str
}

func (widget *KeyboardWidget) SetView(view *tview.TextView) {
	widget.view = view
}

/* -------------------- Unexported Functions -------------------- */

func (widget *KeyboardWidget) moduleName() string {
	return strings.Title(widget.settings.Module.Type)
}

func (widget *KeyboardWidget) showHelp() {
	closeFunc := func() {
		widget.pages.RemovePage("help")
		widget.app.SetFocus(widget.view)
	}

	modal := wtf.NewBillboardModal(widget.HelpText(), closeFunc)

	widget.pages.AddPage("help", modal, false, true)
	widget.app.SetFocus(modal)

	widget.app.QueueUpdate(func() {
		widget.app.Draw()
	})
}

func (widget *KeyboardWidget) showSettings() {
	closeFunc := func() {
		widget.pages.RemovePage("settings")
		widget.app.SetFocus(widget.view)
	}

	str := " [green::b]Settings for " + widget.moduleName() + "[white]\n\n"

	str += fmt.Sprintf(" Type: %s\n", widget.settings.Module.Type)
	str += fmt.Sprint("\n")
	str += fmt.Sprintf(" Refresh: %d seconds\n", widget.settings.RefreshInterval)
	str += fmt.Sprint("\n")
	str += fmt.Sprintf(" Top: %d\n", widget.settings.PositionSettings.Top)
	str += fmt.Sprintf(" Left: %d\n", widget.settings.PositionSettings.Left)
	str += fmt.Sprintf(" Width: %d\n", widget.settings.PositionSettings.Width)
	str += fmt.Sprintf(" Height: %d\n", widget.settings.PositionSettings.Height)

	modal := wtf.NewBillboardModal(str, closeFunc)

	widget.pages.AddPage("settings", modal, false, true)
	widget.app.SetFocus(modal)

	widget.app.QueueUpdate(func() {
		widget.app.Draw()
	})
}
