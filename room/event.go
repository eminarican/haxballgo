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
	eventGameStart      = "onGameStart"
	eventGameStop       = "onGameStop"
	eventGameTick       = "onGameTick"
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

	// onTeamGoal

	p.MustEval(`room.onGameStart = function(player) {
		emit({
			type: "` + eventGameStart + `",
			id: player.id
		})
	}`)

	p.MustEval(`room.onGameStop = function(player) {
		emit({
			type: "` + eventGameStop + `",
			id: player.id
		})
	}`)

	// onPlayerAdminChange
	// onPlayerTeamChange
	// onPlayerKicked

	p.MustEval(`room.onGameTick = function() {
		emit({
			type: "` + eventGameTick + `"
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
	case eventPlayerLeave:
		p := newPlayer(r, obj["id"].Int())
		fun := r.events[eventPlayerLeave].(func(Player))
		fun(p)
	case eventPlayerChat:
		p := newPlayer(r, obj["id"].Int())
		msg := obj["message"].String()
		fun := r.events[eventPlayerChat].(func(Player, string))
		fun(p, msg)
	case eventPlayerBallKick:
		p := newPlayer(r, obj["id"].Int())
		fun := r.events[eventPlayerBallKick].(func(Player))
		fun(p)
	case eventGameStart:
		p := newPlayer(r, obj["id"].Int())
		fun := r.events[eventGameStart].(func(Player))
		fun(p)
	case eventGameStop:
		p := newPlayer(r, obj["id"].Int())
		fun := r.events[eventGameStop].(func(Player))
		fun(p)
	case eventGameTick:
		fun := r.events[eventGameTick].(func())
		fun()
	default:
		return nil, fmt.Errorf("event type %v is invalid", typ)
	}
	return nil, nil
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

func (r *Room) OnPlayerBallKick(fun func(Player)) {
	r.events[eventPlayerBallKick] = fun
}

func (r *Room) OnGameStart(fun func(by Player)) {
	r.events[eventGameStart] = fun
}

func (r *Room) OnGameStop(fun func(by Player)) {
	r.events[eventGameStop] = fun
}

func (r *Room) OnGameTick(fun func()) {
	r.events[eventGameTick] = fun
}
