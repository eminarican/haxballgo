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
	eventPlayerTeamChange  = "onPlayerTeamChange"
	eventPlayerKicked      = "onPlayerKicked"
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
	r.OnPlayerChat(func(p *Player, msg string) (send bool) { return true })
	r.OnPlayerBallKick(func(p *Player) {})
	r.OnGameStart(func(by *Player) {})
	r.OnGameStop(func(by *Player) {})
	r.OnPlayerAdminChange(func(p, by *Player) {})
	r.OnPlayerTeamChange(func(p, by *Player) {})
	r.OnPlayerKicked(func(p *Player, reason string, ban bool, by *Player) {})
	r.OnGameTick(func() {})
	r.OnGamePause(func(by *Player) {})
	r.OnGameUnpause(func(by *Player) {})
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

	// todo: find out why it's not working
	p.MustEval(`room.onPlayerChat = function(player, message) {
		return emit({
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

	p.MustEval(`room.onPlayerTeamChange = function(player, by) {
		emit({
			type: "` + eventPlayerAdminChange + `",
			id: player.id,
			by: by.id
		})
	}`)

	p.MustEval(`room.onPlayerKicked = function(player, reason, ban, by) {
		emit({
			type: "` + eventPlayerAdminChange + `",
			reason: reason,
			ban: ban,
			id: player.id,
			by: by.id
		})
	}`)

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
		fun := r.events[typ].(func(*Player, string) (send bool))
		if !fun(p, msg) {
			return false, nil
		}
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
	case eventPlayerTeamChange:
		p := r.GetPlayer(obj["id"].Int())
		by := r.GetPlayer(obj["by"].Int())
		fun := r.events[typ].(func(*Player, *Player))
		fun(p, by)
	case eventPlayerKicked:
		p := r.GetPlayer(obj["id"].Int())
		by := r.GetPlayer(obj["by"].Int())
		ban := obj["ban"].Bool()
		reason := obj["reason"].String()
		fun := r.events[typ].(func(*Player, string, bool, *Player))
		fun(p, reason, ban, by)
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

// Event called when a new player joins the room.
func (r *Room) OnPlayerJoin(fun func(*Player)) {
	r.events[eventPlayerJoin] = fun
}

// Event called when a player leaves the room.
func (r *Room) OnPlayerLeave(fun func(*Player)) {
	r.events[eventPlayerLeave] = fun
}

// Event called when a player sends a chat message.
//
// The event function can return `false` in order to filter the chat message.
// This prevents the chat message from reaching other players in the room.
func (r *Room) OnPlayerChat(fun func(p *Player, msg string) (send bool)) {
	r.events[eventPlayerChat] = fun
}

// Event called when a player kicks the ball.
func (r *Room) OnPlayerBallKick(fun func(*Player)) {
	r.events[eventPlayerBallKick] = fun
}

// Event called when a game starts.
//
// `by` is the player which caused the event (can be null if the event wasn't caused by a player).
func (r *Room) OnGameStart(fun func(by *Player)) {
	r.events[eventGameStart] = fun
}

// Event called when a game stops.
//
// `by` is the player which caused the event (can be null if the event wasn't caused by a player).
func (r *Room) OnGameStop(fun func(by *Player)) {
	r.events[eventGameStop] = fun
}

// Event called when a player's admin rights are changed.
//
// `by` is the player which caused the event (can be null if the event wasn't caused by a player).
func (r *Room) OnPlayerAdminChange(fun func(p *Player, by *Player)) {
	r.events[eventPlayerAdminChange] = fun
}

// Event called when a player team is changed.
//
// `by` is the player which caused the event (can be null if the event wasn't caused by a player).
func (r *Room) OnPlayerTeamChange(fun func(p *Player, by *Player)) {
	r.events[eventPlayerTeamChange] = fun
}

// Event called when a player has been kicked from the room. This is always called after the onPlayerLeave event.
//
// `by` is the player which caused the event (can be null if the event wasn't caused by a player).
func (r *Room) OnPlayerKicked(fun func(p *Player, reason string, ban bool, by *Player)) {
	r.events[eventPlayerKicked] = fun
}

// Event called once for every game tick (happens 60 times per second).
// This is useful if you want to monitor the player and ball positions without missing any ticks.
//
// This event is not called if the game is paused or stopped.
func (r *Room) OnGameTick(fun func()) {
	r.events[eventGameTick] = fun
}

// Event called when the game is paused.
func (r *Room) OnGamePause(fun func(by *Player)) {
	r.events[eventGamePause] = fun
}

// Event called when the game is unpaused.
//
// After this event there's a timer before the game is fully unpaused,
// to detect when the game has really resumed you can listen for
// the first onGameTick event after this event is called.
func (r *Room) OnGameUnpause(fun func(by *Player)) {
	r.events[eventGameUnpause] = fun
}

// Event called when the players and ball positions are reset after a goal happens.
func (r *Room) OnPositionsReset(fun func()) {
	r.events[eventPositionsReset] = fun
}

// Event called when a player gives signs of activity, such as pressing a key. This is useful for detecting inactive players.
func (r *Room) OnPlayerActivity(fun func(*Player)) {
	r.events[eventPlayerActivity] = fun
}

// Event called when the stadium is changed.
func (r *Room) OnStadiumChange(fun func(stadium string, by *Player)) {
	r.events[eventStadiumChange] = fun
}

// Event called when the room link is obtained.
func (r *Room) OnRoomLink(fun func(link string)) {
	r.events[eventRoomLink] = fun
}

// Event called when the kick rate is set.
func (r *Room) OnKickRateLimitSet(fun func(min int, rate int, burst int, by *Player)) {
	r.events[eventKickRateLimitSet] = fun
}
