package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	DefaultAppWidth  = 500
	DefaultAppHeight = 500
)

type RedditGPTBot struct {
	App    fyne.App
	Window fyne.Window
}

func NewRGB() RedditGPTBot {
	botApp := app.New()
	window := botApp.NewWindow("Reddit GPT Bot")

	window.Resize(fyne.Size{
		Width:  DefaultAppWidth,
		Height: DefaultAppHeight,
	})

	menu := fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Save log...", func() {
				fmt.Println("Log saved.")
			}),
		),
	)

	window.SetMainMenu(menu)

	return RedditGPTBot{
		App:    botApp,
		Window: window,
	}
}

func NewMessageFeed() MessageFeed {
	msgs := []string{
		"Ready.",
	}
	data := binding.BindStringList(
		&msgs,
	)

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Type something...")
	entry.OnSubmitted = func(s string) {
		if err := data.Append(s); err != nil {
			// TODO: log it
		}
	}

	feed := MessageFeed{
		Input:    entry,
		Messages: data,
	}

	feed.ListWidget = widget.NewListWithData(feed.Messages, feed.NewItem, feed.UpdateItem)

	return feed
}

func (bot *RedditGPTBot) Run() {
	messageFeed := NewMessageFeed()

	logs := container.NewVBox(
		messageFeed.Input,
		container.NewMax(messageFeed.ListWidget),
	)
	logs.Layout = layout.NewVBoxLayout()

	bot.Window.SetContent(logs)

	bot.Window.ShowAndRun()
}
