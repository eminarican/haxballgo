package main

import (
	"haxball/room"
)

func main() {
	r := room.New()
	defer r.Shutdown()

	r.OnPlayerJoin(func(p room.Player) {
		println("A player joined to room!")
	})

	r.OnPlayerLeave(func(p room.Player) {
		println("A player leaved from room!")
	})

	println("Successfully started! Room link:", r.Link())

	select {}
}
