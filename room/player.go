package room

import (
	"fmt"
	"strconv"

	"github.com/go-gl/mathgl/mgl32"
)

type Player struct {
	room *Room
	conn string
	name string
	id   int
}

// Creates new player.
func newPlayer(r *Room, id int, conn string, name string) *Player {
	return &Player{
		room: r,
		conn: conn,
		name: name,
		id:   id,
	}
}

// Returns player id.
func (p *Player) Id() int {
	return p.id
}

// Returns player name.
func (p *Player) Name() string {
	return p.name
}

// Returns player conn.
func (p *Player) Conn() string {
	return p.conn
}

// The player's position in the field, if the player is not in the field the value will be null.
func (p *Player) Position() *mgl32.Vec2 {
	obj := p.room.page.MustEval(`room.getPlayer(` + strconv.Itoa(p.id) + `).position`).Map()
	if len(obj) == 2 {
		return &mgl32.Vec2{
			float32(obj["x"].Num()),
			float32(obj["y"].Num()),
		}
	}
	return nil
}

// Sends a formatted chat message to player using the host player.
func (p *Player) Messagef(format string, v ...interface{}) {
	p.Message(fmt.Sprintf(format, v...))
}

// Sends a chat message to player using the host player.
func (p *Player) Message(msg string) {
	p.room.page.MustEval(`room.sendChat("` + msg + `", ` + strconv.Itoa(p.id) + `)`)
}

// Sets player admin privileges.
func (p *Player) SetAdmin(val bool) {
	p.room.page.MustEval(`room.setPlayerAdmin(` + strconv.Itoa(p.id) + `, ` + strconv.FormatBool(val) + `)`)
}

// Whether the player has admin rights.
func (p *Player) IsAdmin() bool {
	return p.room.page.MustEval(`room.getPlayer(` + strconv.Itoa(p.id) + `).admin`).Bool()
}

// Overrides the avatar of the player.
func (p *Player) SetAvatar(val string) {
	p.room.page.MustEval(`room.setPlayerAvatar(` + strconv.Itoa(p.id) + `, "` + val + `")`)
}

// setDiscProperties
// getDiscProperties

// Kicks player from room with aditional ban option.
func (p *Player) Kick(reason string, ban bool) {
	p.room.page.MustEval(`room.kickPlayer(` + strconv.Itoa(p.id) + `, "` + reason + `", ` + strconv.FormatBool(ban) + `)`)
}
