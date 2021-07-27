package room

import (
	"fmt"

	"github.com/go-rod/rod"
	"github.com/ysmood/gson"
)

const (
	eventPlayerJoin        = "onPlayerJoin"
	eventPlayerLeave       = "onPlayerLeave"
	eventPlayerChat        = "onPlayerChat"
	eventPlayerBallKick    = "onPlayerBallKick"
	eventGameStart         = "onGameStart"
	eventGameStop          = "onGameStop"
	eventPlayerAdminChange = "onPlayerAdminChange"
	eventGameTick          = "onGameTick"
	eventGamePause         = "onGamePause"
	eventGameUnpause       = "onGameUnpause"
	eventPositionsReset    = "onPositionsReset"
	eventPlayerActivity    = "onPlayerActivity"
	eventStadiumChange     = "onStadiumChange"
	eventRoomLink          = "onRoomLink"
	eventKickRateLimitSet  = "onKickRateLimitSet"
)

func registerEvents(r *Room, p *rod.Page) {
	r.OnPlayerJoin(func(p *Player) {})
	r.OnPlayerLeave(func(p *Player) {})
	r.OnPlayerChat(func(p *Player, msg string) {})
	r.OnPlayerBallKick(func(p *Player) {})
	r.OnGameStart(func(by *Player) {})
	r.OnGameStop(func(by *Player) {})
	r.OnGameTick(func() {})
	r.OnGamePause(func(p *Player) {})
	r.OnGameUnpause(func(p *Player) {})
	r.OnPositionsReset(func() {})
	r.OnPlayerActivity(func(p *Player) {})
	r.OnStadiumChange(func(stadium string, by *Player) {})
	r.OnRoomLink(func(link string) {})
	r.OnKickRateLimitSet(func(min, rate, burst int, by *Player) {})

	p.MustEval(`room.onPlayerJoin = function(player) {
		emit({
			type: "` + eventPlayerJoin + `",
			conn: player.conn,
			name: player.name,
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

	p.MustEval(`room.onPlayerAdminChange = function(player, by) {
		emit({
			type: "` + eventPlayerAdminChange + `",
			id: player.id,
			by: by.id
		})
	}`)
	// onPlayerTeamChange
	// onPlayerKicked

	p.MustEval(`room.onGameTick = function() {
		emit({
			type: "` + eventGameTick + `"
		})
	}`)

	p.MustEval(`room.onGamePause = function(by) {
		emit({
			type: "` + eventGamePause + `",
			id: by.id
		})
	}`)

	p.MustEval(`room.onGameUnpause = function(by) {
		emit({
			type: "` + eventGameUnpause + `",
			id: by.id
		})
	}`)

	p.MustEval(`room.onPositionsReset = function() {
		emit({
			type: "` + eventPositionsReset + `"
		})
	}`)

	p.MustEval(`room.onPlayerActivity = function(player) {
		emit({
			type: "` + eventPlayerActivity + `",
			id: player.id
		})
	}`)

	p.MustEval(`room.onStadiumChange = function(stadium, by) {
		emit({
			type: "` + eventStadiumChange + `",
			stadium: stadium,
			id: by.id
		})
	}`)

	p.MustEval(`room.onRoomLink = function(link) {
		emit({
			type: "` + eventRoomLink + `",
			link: link
		})
	}`)

	p.MustEval(`room.onKickRateLimitSet = function(min, rate, burst, by) {
		emit({
			type: "` + eventKickRateLimitSet + `",
			burst: burst,
			rate: rate,
			min: min,
			id: by.id
		})
	}`)
}

func proccessEvent(r *Room, j gson.JSON) (interface{}, error) {
	obj := j.Map()
	typ := obj["type"].String()

	switch typ {
	case eventPlayerJoin:
		p := newPlayer(r, obj["id"].Int(), obj["conn"].String(), obj["name"].String())
		fun := r.events[typ].(func(*Player))
		r.setPlayer(p.Id(), p)
		fun(p)
	case eventPlayerLeave:
		p := r.GetPlayer(obj["id"].Int())
		fun := r.events[typ].(func(*Player))
		fun(p)
		r.removePlayer(p.Id())
	case eventPlayerChat:
		p := r.GetPlayer(obj["id"].Int())
		msg := obj["message"].String()
		fun := r.events[typ].(func(*Player, string))
		fun(p, msg)
	case eventPlayerBallKick:
		p := r.GetPlayer(obj["id"].Int())
		fun := r.events[typ].(func(*Player))
		fun(p)
	case eventGameStart:
		p := r.GetPlayer(obj["id"].Int())
		fun := r.events[typ].(func(*Player))
		fun(p)
	case eventGameStop:
		p := r.GetPlayer(obj["id"].Int())
		fun := r.events[typ].(func(*Player))
		fun(p)
	case eventPlayerAdminChange:
		p := r.GetPlayer(obj["id"].Int())
		by := r.GetPlayer(obj["by"].Int())
		fun := r.events[typ].(func(*Player, *Player))
		fun(p, by)
	case eventGameTick:
		fun := r.events[typ].(func())
		fun()
	case eventGamePause:
		by := r.GetPlayer(obj["id"].Int())
		fun := r.events[typ].(func(*Player))
		fun(by)
	case eventGameUnpause:
		by := r.GetPlayer(obj["id"].Int())
		fun := r.events[typ].(func(*Player))
		fun(by)
	case eventPositionsReset:
		fun := r.events[typ].(func())
		fun()
	case eventPlayerActivity:
		p := r.GetPlayer(obj["id"].Int())
		fun := r.events[typ].(func(*Player))
		fun(p)
	case eventStadiumChange:
		by := r.GetPlayer(obj["id"].Int())
		stadium := obj["stadium"].String()
		fun := r.events[typ].(func(string, *Player))
		fun(stadium, by)
	case eventRoomLink:
		link := obj["link"].String()
		fun := r.events[typ].(func(string))
		fun(link)
	case eventKickRateLimitSet:
		by := r.GetPlayer(obj["id"].Int())
		burst := obj["burst"].Int()
		rate := obj["rate"].Int()
		min := obj["min"].Int()
		fun := r.events[typ].(func(int, int, int, *Player))
		fun(min, rate, burst, by)
	default:
		return nil, fmt.Errorf("event type %v is invalid", typ)
	}
	return nil, nil
}

func (r *Room) OnPlayerJoin(fun func(*Player)) {
	r.events[eventPlayerJoin] = fun
}

func (r *Room) OnPlayerLeave(fun func(*Player)) {
	r.events[eventPlayerLeave] = fun
}

func (r *Room) OnPlayerChat(fun func(p *Player, msg string)) {
	r.events[eventPlayerChat] = fun
}

func (r *Room) OnPlayerBallKick(fun func(*Player)) {
	r.events[eventPlayerBallKick] = fun
}

func (r *Room) OnGameStart(fun func(by *Player)) {
	r.events[eventGameStart] = fun
}

func (r *Room) OnGameStop(fun func(by *Player)) {
	r.events[eventGameStop] = fun
}

func (r *Room) OnPlayerAdminChange(fun func(p *Player, by *Player)) {
	r.events[eventPlayerAdminChange] = fun
}

func (r *Room) OnGameTick(fun func()) {
	r.events[eventGameTick] = fun
}

func (r *Room) OnGamePause(fun func(*Player)) {
	r.events[eventGamePause] = fun
}

func (r *Room) OnGameUnpause(fun func(*Player)) {
	r.events[eventGameUnpause] = fun
}

func (r *Room) OnPositionsReset(fun func()) {
	r.events[eventPositionsReset] = fun
}

func (r *Room) OnPlayerActivity(fun func(*Player)) {
	r.events[eventPlayerActivity] = fun
}

func (r *Room) OnStadiumChange(fun func(stadium string, by *Player)) {
	r.events[eventStadiumChange] = fun
}

func (r *Room) OnRoomLink(fun func(link string)) {
	r.events[eventRoomLink] = fun
}

func (r *Room) OnKickRateLimitSet(fun func(min int, rate int, burst int, by *Player)) {
	r.events[eventKickRateLimitSet] = fun
}
