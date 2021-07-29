package main

import (
	"haxball/room"
)

func main() {
	r := room.New()
	l := r.Logger()
	defer r.Shutdown()

	r.OnPlayerJoin(func(p *room.Player) {
		l.Info("A player joined to room!")
	})

	r.OnPlayerLeave(func(p *room.Player) {
		l.Info("A player leaved from room!")
	})

	r.OnPlayerChat(func(p *room.Player, msg string) (send bool) {
		l.Infof("%v:%v", p.Name(), msg)
		return true
	})
}
