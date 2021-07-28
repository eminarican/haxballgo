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
		r.Logger().Info("A player joined to room!")
	})

	r.OnPlayerLeave(func(p *room.Player) {
		r.Logger().Info("A player leaved from room!")
	})

	r.OnPlayerChat(func(p *room.Player, msg string) (send bool) {
		r.Logger().Infof("%v:%v", p.Name(), msg)
		return true
	})

	r.Logger().Infof("Successfully started! Room link: %v", r.Link())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
