package cd

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/therealkevinard/gitdir/ui/styles"
)

type itemDelegate struct {
	stripPrefix string
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

//nolint:gocritic
func (d itemDelegate) Render(writer io.Writer, model list.Model, index int, listItem list.Item) {
	i, ok := listItem.(stringListItem)
	if !ok {
		return
	}

	cleanstr := strings.TrimLeft(strings.TrimPrefix(string(i), d.stripPrefix), "/")
	str := fmt.Sprintf("%d. %s", index+1, cleanstr)

	fn := styles.ItemStyle.Render
	if index == model.Index() {
		fn = func(s ...string) string {
			return styles.SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(writer, fn(str))
}
