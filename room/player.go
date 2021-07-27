package room

import "strconv"

type Player struct {
	room *Room
	conn string
	name string
	id   int
}

func newPlayer(r *Room, id int, conn string, name string) *Player {
	return &Player{
		room: r,
		conn: conn,
		name: name,
		id:   id,
	}
}

func (p *Player) Id() int {
    return p.id
}

func (p *Player) Name() string {
    return p.name
}

func (p *Player) Conn() string {
    return p.conn
}

func (p *Player) SendMessage(msg string) {
	p.room.page.MustEval(`room.sendChat("` + msg + `", ` + strconv.Itoa(p.id) + `)`)
}
