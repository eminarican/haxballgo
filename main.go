package main

import (
	"haxball/room"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	r := room.New()
	defer r.Shutdown()

	r.OnPlayerJoin(func(p *room.Player) {
		println("A player joined to room!")
	})

	r.OnPlayerLeave(func(p *room.Player) {
		println("A player leaved from room!")
	})

	r.OnPlayerChat(func(p *room.Player, msg string) (send bool) {
		println(p.Name()+":", msg)
		return true
	})

	println("Successfully started! Room link:", r.Link())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
