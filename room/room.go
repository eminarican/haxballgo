package room

import (
	"strconv"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/ysmood/gson"
)

const target = "https://html5.haxball.com/hiF05fAx/__cache_static__/g/headless.html"

type Room struct {
	conf    *Config
	page    *rod.Page
	browser *rod.Browser
	events  map[string]interface{}
	players map[int]*Player
	pMutex  sync.RWMutex
}

// Creates a new room
func New() *Room {
	conf, err := readConfig()
	if err != nil {
		panic(err)
	}

	browser := rod.New().MustConnect()
	page := browser.
		MustPage(target).
		MustWaitLoad()
	time.Sleep(time.Second * 2)

	r := &Room{
		conf:    &conf,
		page:    page,
		browser: browser,
		events:  make(map[string]interface{}),
		players: make(map[int]*Player),
	}

	page.MustEval(conf.String())

	page.MustExpose("emit", func(j gson.JSON) (interface{}, error) {
		return proccessEvent(r, j)
	})

	registerEvents(r, page)

	return r
}

// Obtains room link. (if token is invalid just pauses program)
func (r *Room) Link() string {
	// todo: replace this with event based one
	return r.page.MustElement("#roomlink a").MustText()
}

// Shuts down room.
func (r *Room) Shutdown() {
	r.browser.MustClose()
}

// Sends a host announcement with msg as contents. Unlike sendChat,
// announcements will work without a host player and has a larger limit on the number of characters.
func (r *Room) Announce(msg string) {
	// todo: add missing fields
	r.page.MustEval(`room.sendAnnouncement("` + msg + `")`)
}

// Sends a chat message using the host player.
func (r *Room) SendMessage(msg string) {
	r.page.MustEval(`room.sendChat("` + msg + `")`)
}

// Clears the ban for a playerId that belonged 
// to a player that was previously banned.
func (r *Room) ClearBan(id int) {
    r.page.MustEval(`room.clearBan(` + strconv.Itoa(id) + `)`)
}

// Clears the list of banned players.
func (r *Room) ClearBans() {
	r.page.MustEval(`room.clearBans()`)
}

// Sets the time limit of the room. The limit must be specified in number of minutes.
//
// If a game is in progress this method does nothing.
func (r *Room) SetTimeLimit(val int) {
    r.page.MustEval(`room.setTimeLimit(` + strconv.Itoa(val) + `)`)
}

// Parses the value as a .hbs stadium file and sets it as the selected stadium.
//
// There must not be a game in progress, if a game is in progress this method does nothing.
func (r *Room) SetCustomStadium(val string) {
	r.page.MustEval(`room.setCustomStadium(` + val + `)`)
}

// Sets the selected stadium to one of the default stadiums.
// The name must match exactly. (case sensitive)
//
// There must not be a game in progress, if a game is in progress this method does nothing.
func (r *Room) SetDefaultStadium(name string) {
	r.page.MustEval(`room.setDefaultStadium(` + name + `)`)
}

// Sets the teams lock. When teams are locked players are not able to change team unless they are moved by an admin.
func (r *Room) SetTeamsLock(val bool) {
	r.page.MustEval(`room.setTeamsLock(` + strconv.FormatBool(val) + `)`)
}

// setTeamColors

// Gets a player from room. (returns nil if player doesn't exists)
func (r *Room) GetPlayer(id int) *Player {
	defer r.pMutex.RUnlock()
	r.pMutex.RLock()

	return r.players[id]
}

func (r *Room) setPlayer(id int, p *Player) {
	defer r.pMutex.Unlock()
	r.pMutex.Lock()

	r.players[id] = p
}

func (r *Room) removePlayer(id int) {
	defer r.pMutex.Unlock()
	r.pMutex.Lock()

	delete(r.players, id)
}
