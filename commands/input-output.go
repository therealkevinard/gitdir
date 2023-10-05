package commands

import (
	"github.com/fatih/color"
)

type lineDecorator string

const (
	sparkleDecorator   lineDecorator = "✨"
	bangBangDecorator  lineDecorator = "‼️"
	infoDecorator      lineDecorator = "💬"
	handraiseDecorator lineDecorator = "🙌"
	cabinetDecorator   lineDecorator = "🗄"
)

type NotifyKind string

const (
	NotifyCreate NotifyKind = "create"
	NotifyClone  NotifyKind = "clone"
	NotifyDone   NotifyKind = "done"

	NotifyInfo  NotifyKind = "info"
	NotifyError NotifyKind = "error"
)

type lineConfig struct {
	decorator lineDecorator
	color     *color.Color
}

//nolint:gochecknoglobals
var decoratorMap = map[NotifyKind]lineConfig{
	NotifyCreate: {
		decorator: sparkleDecorator,
		color:     color.New(color.FgGreen),
	},
	NotifyClone: {
		decorator: cabinetDecorator,
		color:     color.New(color.FgGreen),
	},
	NotifyDone: {
		decorator: handraiseDecorator,
		color:     color.New(color.FgYellow),
	},

	NotifyInfo: {
		decorator: infoDecorator,
		color:     color.New(color.FgWhite),
	},
	NotifyError: {
		decorator: bangBangDecorator,
		color:     color.New(color.FgRed),
	},
}

func Notify(kind NotifyKind, message string) {
	lc, ok := decoratorMap[kind]
	if !ok {
		lc = decoratorMap[NotifyInfo]
	}

	_, _ = lc.color.Println(string(lc.decorator), "\t", message)
}
