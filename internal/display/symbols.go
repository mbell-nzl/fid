package display

// Symbol marks a line's kind. fid has no display config, so the symbols are
// always Unicode.
type Symbol string

const (
	SymbolSuccess   Symbol = "✓"
	SymbolError     Symbol = "✗"
	SymbolWarning   Symbol = "⚠"
	SymbolStarted   Symbol = "●"
	SymbolUnstarted Symbol = "○"
	SymbolEpic      Symbol = "◆"
	SymbolNext      Symbol = "→"
)

// Ellipsis ends text cut to fit its column.
const Ellipsis = "…"
