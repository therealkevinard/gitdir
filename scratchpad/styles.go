package main

import (
	"errors"
	"fmt"

	"github.com/therealkevinard/gitdir/ui/styles"
)

func main() {
	err := errors.New("ehrmhergherd!")

	for _, s := range []string{
		styles.TitleStyle.Render("title"),
		styles.ItemStyle.Render("item"),
		styles.SelectedItemStyle.Render("selected item"),
		styles.PaginationStyle.Render("pagination"),
		styles.HelpStyle.Render("help"),
		styles.QuitTextStyle.Render("quit"),

		"> some text lines",
		styles.OKTextf("an okay thing happened\nlorem ipsum lipsum is lorem ipsum\nqwe rtyu iop"),
		styles.AltTextf("alttext an okay alttext thing happened\nlorem ipsum alttext lipsum is lorem ipsum\nqwe alttext rtyu iop"),
		styles.WarnTextf("medium-bad thing happened"),
		styles.ErrorTextf("bad thing happened"),
		styles.FatalTextf("awful thing happened: %v", err),
	} {
		fmt.Println(s)
	}
}
