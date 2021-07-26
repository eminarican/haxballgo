package main

import (
	"haxball/room"

	"github.com/go-rod/rod"
)

func main() {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	r := room.New(browser)

	r.OnPlayerJoin(func(p room.Player) {
		p.SendMessage("Sunucuya hosgelmise")
	})

	println("Successfully started! Room link:", r.Link())

	select {}
}
