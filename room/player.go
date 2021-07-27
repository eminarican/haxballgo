package room

import "strconv"

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

// Sends a chat message to player using the host player.
func (p *Player) SendMessage(msg string) {
	p.room.page.MustEval(`room.sendChat("` + msg + `", ` + strconv.Itoa(p.id) + `)`)
}

// Sets player admin privileges
func (p *Player) SetAdmin(val bool) {
	p.room.page.MustEval(`room.setPlayerAdmin(` + strconv.Itoa(p.id) + `, ` + strconv.FormatBool(val) + `)`)
}

// Kicks player from room with aditional ban option
func (p *Player) Kick(reason string, ban bool) {
	p.room.page.MustEval(`room.kickPlayer(` + strconv.Itoa(p.id) + `, "` + reason + `", ` + strconv.FormatBool(ban) + `)`)
}
