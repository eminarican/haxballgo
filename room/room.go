package room

import (
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
