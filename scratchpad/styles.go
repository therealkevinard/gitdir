package main

import (
	"fmt"

	"github.com/therealkevinard/gitdir/ui/styles"
)

func main() {
	for _, s := range []string{
		styles.TitleStyle.Render("title"),
		styles.ItemStyle.Render("item"),
		styles.SelectedItemStyle.Render("selected item"),
		styles.PaginationStyle.Render("pagination"),
		styles.HelpStyle.Render("help"),
		styles.QuitTextStyle.Render("quit"),

		"> some text lines",
		styles.OKText("an okay thing happened"),
		styles.WarnText("medium-bad thing happened"),
		styles.ErrorText("bad thing happened"),
		styles.FatalText("awful thing happened"),
	} {
		fmt.Println(s)
	}
}
