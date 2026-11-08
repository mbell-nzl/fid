package display

const ansiReset = "\033[0m"

// Role names what a piece of text is, so its colour stays consistent.
type Role string

const (
	RoleNormal  Role = ""
	RoleHeader  Role = "header"
	RoleID      Role = "id"
	RoleError   Role = "error"
	RoleWarning Role = "warning"
	RoleSuccess Role = "success"
	RoleMuted   Role = "muted"
	RoleLabel   Role = "label"
)

var colorCodes = map[Role]string{
	RoleHeader:  "\033[1m",
	RoleID:      "\033[36m",
	RoleError:   "\033[31m",
	RoleWarning: "\033[33m",
	RoleSuccess: "\033[32m",
	RoleMuted:   "\033[2m",
	RoleLabel:   "\033[35m",
}

// Colorize wraps text in its role's colour when colour is on. Pad text
// before colouring it, since padding counts runes.
func (r Renderer) Colorize(role Role, text string) string {
	code := colorCodes[role]
	if !r.Color || code == "" {
		return text
	}
	return code + text + ansiReset
}

// Line is a symbol and its text, both in the role's colour.
func (r Renderer) Line(role Role, symbol Symbol, text string) string {
	return r.Colorize(role, string(symbol)) + " " + r.Colorize(role, text)
}
