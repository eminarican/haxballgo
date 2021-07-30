package room

import (
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-rod/rod"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ysmood/gson"
)

const target = "https://html5.haxball.com/hiF05fAx/__cache_static__/g/headless.html"

type Room struct {
	conf      *Config
	page      *rod.Page
	browser   *rod.Browser
	events    map[string]interface{}
	players   map[int]*Player
	pMutex    sync.RWMutex
	scheduler *Scheduler
	logger    *Logger
	link      string
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
		conf:      &conf,
		page:      page,
		browser:   browser,
		events:    make(map[string]interface{}),
		players:   make(map[int]*Player),
		scheduler: &Scheduler{},
		logger:    &Logger{},
	}

	l := r.Logger()
	if conf.Logging.Pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	zerolog.SetGlobalLevel(func() zerolog.Level {
		if conf.Logging.Debug {
			return zerolog.DebugLevel
		}
		return zerolog.InfoLevel
	}())

	l.Info("Starting room...")

	page.MustEval(conf.String())

	page.MustExpose("emit", func(j gson.JSON) (interface{}, error) {
		return proccessEvent(r, j)
	})

	link := registerEvents(r, page)
	stop := make(chan bool)

	r.Scheduler().Delayed(5*time.Second, func() {
		stop <- true
	})

	select {
	case <-stop:
		l.Error("Token is not valid!")
		os.Exit(1)
	case link := <-link:
		l.Info("Successfully started!")
		l.Infof("Room link: %v", link)
		r.link = link
	}
	return r
}

// Obtains room link.
func (r *Room) Link() string {
	return r.link
}

// Gets logger.
func (r *Room) Logger() *Logger {
	return r.logger
}

// Gets scheduler.
func (r *Room) Scheduler() *Scheduler {
	return r.scheduler
}

// Waits receive signal to shutdown room properly.
func (r *Room) Shutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
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

// Starts the game, if a game is already in progress this method does nothing.
func (r *Room) StartGame() {
	r.page.MustEval(`room.startGame()`)
}

// Stops the game, if no game is in progress this method does nothing.
func (r *Room) StopGame() {
	r.page.MustEval(`room.stopGame()`)
}

// Sets the pause state of the game.
func (r *Room) PauseGame(val bool) {
	r.page.MustEval(`room.pauseGame(` + strconv.FormatBool(val) + `)`)
}

// getScores

// Returns the ball's position in the field or null if no game is in progress.
func (r *Room) GetBallPosition() *mgl32.Vec2 {
	obj := r.page.MustEval(`room.getBallPosition()`).Map()
	if len(obj) == 2 {
		return &mgl32.Vec2{
			float32(obj["x"].Num()),
			float32(obj["y"].Num()),
		}
	}
	return nil
}

// Starts recording of a haxball replay.
//
// Don't forget to call stop recording or it will cause a memory leak.
func (r *Room) StartRecording() {
	r.page.MustEval(`room.startRecording()`)
}

// Stops the recording previously started with startRecording and returns the replay file contents as a []uint8.
//
// Returns null if recording was not started or had already been stopped.
func (r *Room) StopRecording() []uint8 {
	var buf []uint8
	for _, b := range r.page.MustEval(`room.stopRecording()`).Arr() {
		buf = append(buf, uint8(b.Int()))
	}
	if len(buf) == 0 {
		return nil
	}
	return buf
}

// Changes the password of the room, if pass is null the password will be cleared.
func (r *Room) SetPassword(val string) {
	r.page.MustEval(`room.setPassword("` + val + `")`)
}

// Activates or deactivates the recaptcha requirement to join the room.
func (r *Room) SetRequireRecaptcha(val bool) {
	r.page.MustEval(`room.setRequireRecaptcha(` + strconv.FormatBool(val) + `)`)
}

// reorderPlayers

// Sets the room's kick rate limits.
//
// `min` is the minimum number of logic-frames between two kicks. It is impossible to kick faster than this.
//
// `rate` works like `min` but lets players save up extra kicks to use them later depending on the value of `burst`.
//
// `burst` determines how many extra kicks the player is able to save up.
func (r *Room) SetKickRateLimit(min int, rate int, burst int) {
	r.page.MustEval(`room.setKickRateLimit(` + strconv.Itoa(min) + `, ` + strconv.Itoa(rate) + `, ` + strconv.Itoa(burst) + `)`)
}

// setDiscProperties
// getDiscProperties

// Gets the number of discs in the game including the ball and player discs.
func (r *Room) GetDiscCount() int {
	return r.page.MustEval(`room.getDiscCount()`).Int()
}

// CollisionFlags

// Returns the current list of players.
func (r *Room) GetPlayers() []*Player {
	defer r.pMutex.RUnlock()
	r.pMutex.Lock()

	slice := make([]*Player, len(r.players))
	for _, p := range r.players {
		slice = append(slice, p)
	}
	return slice
}

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
