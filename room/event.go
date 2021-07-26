package room

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/ysmood/gson"
)

const (
	eventPlayerJoin     = "onPlayerJoin"
	eventPlayerLeave    = "onPlayerLeave"
	eventPlayerChat     = "onPlayerChat"
	eventPlayerBallKick = "onPlayerBallKick"
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

	// onTeamVictory

	p.MustEval(`room.onPlayerChat = function(player, message) {
		emit({
			type: "` + eventPlayerChat + `",
			message: message,
			id: player.id
		})
	}`)

	p.MustEval(`room.onPlayerBallKick = function(player) {
		emit({
			type: "` + eventPlayerBallKick + `",
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
	case eventPlayerChat:
		p := newPlayer(r, obj["id"].Int())
		msg := obj["message"].String()
		fun := r.events[eventPlayerChat].(func(Player, string))
		fun(p, msg)
		return nil, nil
	case eventPlayerBallKick:
		p := newPlayer(r, obj["id"].Int())
		fun := r.events[eventPlayerBallKick].(func(Player))
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

func (r *Room) OnPlayerChat(fun func(p Player, msg string)) {
	r.events[eventPlayerChat] = fun
}

func (r *Room) OnPlayerBallKick(fun func(p Player)) {
	r.events[eventPlayerBallKick] = fun
}
