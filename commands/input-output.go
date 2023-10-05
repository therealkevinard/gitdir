package commands

import "fmt"

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

var decoratorMap = map[NotifyKind]lineDecorator{
	NotifyCreate: sparkleDecorator,
	NotifyClone:  cabinetDecorator,
	NotifyDone:   handraiseDecorator,

	NotifyInfo:  infoDecorator,
	NotifyError: bangBangDecorator,
}

func Notify(kind NotifyKind, message string) {
	i, ok := decoratorMap[kind]
	if !ok {
		i = " "
	}

	fmt.Println(string(i), "\t", message)
}
