//nolint:gochecknoglobals,gomnd
package styles

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"log"
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
	AltTextStyle     = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#4272f5"))
	WarningTextStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#ffff00"))
	ErrorTextStyle   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#f54542"))
	FatalTextStyle   = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#e31959")).Bold(true)
)

// text helpers.
func OKTextf(message string, v ...any) string {
	return OKTextStyle.Render(fmt.Sprintf("> "+message, v...))
}
func AltTextf(message string, v ...any) string {
	return AltTextStyle.Render(fmt.Sprintf(message, v...))
}
func WarnTextf(message string, v ...any) string {
	return WarningTextStyle.Render(fmt.Sprintf("⚠ "+message, v...))
}
func ErrorTextf(message string, v ...any) string {
	return ErrorTextStyle.Render(fmt.Sprintf("‼ "+message, v...))
}
func FatalTextf(message string, v ...any) string {
	return FatalTextStyle.Render(fmt.Sprintf("X "+message, v...))
}

type Level int

const (
	OKLevel Level = iota
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Println passes one of the leveled formatters to fmt.Println.
// if level is FatalLevel, it additionally calls log.Fatalln
func Println(level Level, message string, v ...any) {
	var line string

	switch level {
	case WarnLevel:
		line = WarnTextf(message, v...)

	case ErrorLevel:
		line = ErrorTextf(message, v...)

	case FatalLevel:
		line = FatalTextf(message, v...)

	case OKLevel:
		fallthrough
	default:
		line = OKTextf(message, v...)
	}

	fmt.Println(line)
	if level == FatalLevel {
		log.Fatalln("exiting")
	}
}

// component styles.
var (
	// list items.
	ItemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))

	// pagination.
	PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
)
