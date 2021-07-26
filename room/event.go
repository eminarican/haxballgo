package room

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/ysmood/gson"
)

const (
	eventPlayerJoin  = "onPlayerJoin"
	eventPlayerLeave = "onPlayerLeave"
)

func registerEvents(p *rod.Page) {
	p.MustEval(`room.onPlayerJoin = function(player) {
		emit({
			type: "` + eventPlayerJoin + `",
			id: player.id
		})
	}`)

	p.MustEval(`room.onPlayerLeave = function(player) {
		emit({
			type: "` + eventPlayerLeave + `",
			id: player.id
		})
	}`)
}

func proccessEvent(r *Room, j gson.JSON) (interface{}, error) {
	obj := j.Map()
	typ := obj["type"].String()

	switch typ {
	case eventPlayerJoin:
		p := newPlayer(r, obj["id"].Int())
		fun := r.events[eventPlayerJoin].(func(Player))
		fun(p)
		return nil, nil
	case eventPlayerLeave:
		p := newPlayer(r, obj["id"].Int())
		fun := r.events[eventPlayerLeave].(func(Player))
		fun(p)
		return nil, nil
	}
	return nil, fmt.Errorf("event type %v is invalid", typ)
}

func (r *Room) OnPlayerJoin(fun func(Player)) {
	r.events[eventPlayerJoin] = fun
}

func (r *Room) OnPlayerLeave(fun func(Player)) {
	r.events[eventPlayerLeave] = fun
}
