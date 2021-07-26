package room

import "strconv"

type Player struct {
	room *Room
	id   int
}

func newPlayer(r *Room, id int) Player {
	return Player{
		room: r,
		id:   id,
	}
}

func (p *Player) SendMessage(msg string) {
	p.room.page.MustEval(`room.sendChat("` + msg + `", ` + strconv.Itoa(p.id) + `)`)
}
