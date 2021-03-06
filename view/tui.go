package view

import (
	"chatroom/client"
	"fmt"
	"io"

	"github.com/marcusolsson/tui-go"
)

func UI(c client.ChatClient) {
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
		c.SetName(username)
		ui.SetWidget(chatView)
	})

	chatView.OnSubmit(func(msg string) {
		c.SendMessage(msg)
	})

	go func() {
		for {
			select {
			case err := <-c.Error():

				if err == io.EOF {
					ui.Update(func() {
						chatView.AddMessage("Connection closed connection from server.")
					})
				} else {
					panic(err)
				}
			case msg := <-c.Incoming():
				ui.Update(func() {
					chatView.AddMessage(fmt.Sprintf("%v: %v", msg.Name, msg.Message))
				})
			}
		}
	}()

	if err := ui.Run(); err != nil {
		panic(err)
	}
}
