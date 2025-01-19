package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type Config struct {
	Title        string
	Edit         *widget.Entry
	Preview      *widget.RichText
	Menu         *fyne.MainMenu
	CurrentFile  fyne.URI
	SaveMenuItem *fyne.MenuItem
	Filter       storage.FileFilter
}

var cfg Config

func main() {
	a := app.NewWithID("fed")
	edit, preview := cfg.makeUI()
	w := a.NewWindow(cfg.Title)
	w.SetContent(container.NewHSplit(edit, preview))
	cfg.createMemu(w)
	w.Resize(fyne.Size{Width: 1000, Height: 600})
	w.CenterOnScreen()

	w.ShowAndRun()
}
