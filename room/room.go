package room

import (
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
	pMutex sync.RWMutex
}

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

func (r *Room) Link() string {
	return r.page.MustElement("#roomlink a").MustText()
}

func (r *Room) Shutdown() {
	r.browser.MustClose()
}

func (r *Room) Announce(msg string) {
	r.page.MustEval(`room.sendAnnouncement("` + msg + `")`)
}

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
