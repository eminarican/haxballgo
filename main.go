package main

import (
	"haxball/room"
)

func main() {
	r := room.New()
	defer r.Shutdown()

	r.OnPlayerJoin(func(p room.Player) {
		p.SendMessage("Sunucuya hosgelmise")
	})

	println("Successfully started! Room link:", r.Link())

	select {}
}
