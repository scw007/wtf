package wtf

import (
	"time"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/cfg"
)

// Wtfable is the interface that enforces WTF system capabilities on a module
type Wtfable interface {
	Enablable
	Schedulable
	Stoppable

	BorderColor() string
	ConfigText() string
	FocusChar() string
	Focusable() bool
	HelpText() string
	QuitChan() chan bool
	Name() string
	RefreshedAt() time.Time
	SetFocusChar(string)
	TextView() *tview.TextView

	CommonSettings() *cfg.Common
}
