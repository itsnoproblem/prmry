package ui

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type MessageFeed struct {
	Input      *widget.Entry
	ListWidget fyne.CanvasObject
	Messages   binding.ExternalStringList
	App        RedditGPTBot
}

func (f *MessageFeed) Len() int {
	list, _ := f.Messages.Get()
	return len(list)
}

func (f *MessageFeed) NewItem() fyne.CanvasObject {
	l := widget.NewLabel("template")
	return l
}

func (f *MessageFeed) UpdateItem(item binding.DataItem, obj fyne.CanvasObject) {
	thing, ok := obj.(*widget.Label)
	if !ok {
		dialog.ShowError(errors.New("failed to cast data item"), f.App.Window)
	}

	thing.Bind(item.(binding.String))

	err := f.Messages.Reload()
	if err != nil {
		dialog.ShowError(err, f.App.Window)
	}

	f.Input.SetText("")
}

func (f *MessageFeed) AddMessages(messages ...string) {
	for _, v := range messages {
		if err := f.Messages.Append(v); err != nil {
			// TODO: log it
		}
	}
	f.ListWidget.Refresh()
}
