package tui

import (
	"github.com/marcusolsson/tui-go"
)

func StartUi() {
	loginView := NewLoginView()
	chatView := NewChatView()

	ui, err := tui.New(loginView)
	if err != nil {
		panic(err)
	}

	quit := func() { ui.Quit() }

	ui.SetKeybinding("Esc", quit)
	ui.SetKeybinding("Ctrl+c", quit)

	loginView.OnLogin(func(username string) {
		ui.SetWidget(chatView)
	})

	if err := ui.Run(); err != nil {
		panic(err)
	}
}
