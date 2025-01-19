package main

import (
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func (cfg *Config) makeUI() (*widget.Entry, *widget.RichText) {
	cfg.Title = "Fed"
	cfg.Edit = widget.NewMultiLineEntry()
	cfg.Preview = widget.NewRichTextFromMarkdown("")
	cfg.Edit.OnChanged = cfg.Preview.ParseMarkdown
	cfg.Filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

	return cfg.Edit, cfg.Preview
}

func (cfg *Config) createMemu(win fyne.Window) {
	open := fyne.NewMenuItem("Open...", cfg.openFunc(win))
	save := fyne.NewMenuItem("Save", cfg.saveFunc(win))
	saveAs := fyne.NewMenuItem("Save as...", cfg.saveAsFunc(win))
	fileMenu := fyne.NewMenu("File", open, save, saveAs)

	cfg.Menu = fyne.NewMainMenu(fileMenu)
	cfg.SaveMenuItem = save
	cfg.SaveMenuItem.Disabled = true

	win.SetMainMenu(cfg.Menu)
}

func (cfg *Config) saveAsFunc(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if writer == nil {
				return
			}
			if !strings.HasSuffix(strings.ToLower(writer.URI().String()), ".md") {
				dialog.ShowInformation("Warning", "File must with suffix \".md\" or \".MD\"", win)
				return
			}

			writer.Write([]uint8(cfg.Edit.Text))
			cfg.CurrentFile = writer.URI()
			win.SetTitle(cfg.Title + " - " + writer.URI().Name())
			cfg.SaveMenuItem.Disabled = false

			defer writer.Close()
		}, win)
		saveDialog.SetFileName("NewFile.md")
		saveDialog.SetFilter(cfg.Filter)
		saveDialog.Show()
	}
}

func (cfg *Config) openFunc(win fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if reader == nil {
				return
			}

			data, err := io.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			cfg.Edit.SetText(string(data))
			cfg.CurrentFile = reader.URI()
			win.SetTitle(cfg.Title + " - " + reader.URI().Name())
			cfg.SaveMenuItem.Disabled = false

			defer reader.Close()
		}, win)
		openDialog.SetFilter(cfg.Filter)
		openDialog.Show()
	}
}

func (cfg *Config) saveFunc(win fyne.Window) func() {
	return func() {
		if cfg.CurrentFile != nil {
			writer, err := storage.Writer(cfg.CurrentFile)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			writer.Write([]byte(cfg.Edit.Text))

			defer writer.Close()
		}
	}
}
