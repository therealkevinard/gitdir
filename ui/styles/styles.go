//nolint:gochecknoglobals,gomnd
package styles

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// special text styles.
var (
	TitleStyle    = lipgloss.NewStyle().MarginLeft(2)
	HelpStyle     = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	QuitTextStyle = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// basic text styles.
var (
	OKTextStyle      = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#9342f5"))
	WarningTextStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#ffff00"))
	ErrorTextStyle   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#f54542"))
	FatalTextStyle   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#e31959")).Bold(true)
)

// text helpers.
func OKText(message string) string    { return OKTextStyle.Render("> " + message) }
func WarnText(message string) string  { return WarningTextStyle.Render("⚠ " + message) }
func ErrorText(message string) string { return ErrorTextStyle.Render("‼ " + message) }
func FatalText(message string) string { return FatalTextStyle.Render("X️ " + message) }

// component styles.
var (
	// list items.
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))

	// pagination.
	PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
)
