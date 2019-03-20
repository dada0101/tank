package game

import (
	"TankDemo/network"
	"TankDemo/proto"
	"fmt"
	"log"
	"sync"
)

const (
	ROOM_STATUS_PREPARE = 1
	ROME_STATUS_FIGHT = 2

	ROOM_SIZE = 6
)


type Room struct {
	mu sync.Mutex
	players [ROOM_SIZE]*Player
	status int
	playerCnt int

	group *network.Group
}

func RoomDump(r *Room) string{
	str := "Room Dump\n"
	str += fmt.Sprintf("players :\t %v\n", r.players)
	if r.status == ROOM_STATUS_PREPARE {
		str += fmt.Sprintf("state :\t ROOM_STATUS_PREPARE\n")
	} else {
		str += fmt.Sprintf("state:\t ROME_STATUS_FIGHT\n")
	}
	str += fmt.Sprintf("group %v", r.group)
	str += fmt.Sprintf("playerCnt:\t %v", r.playerCnt)
	return str
}


func NewRoom(g *network.Group) *Room{
	r := new(Room)
	for i := 0; i < len(r.players); i++ {
		r.players[i] = nil
	}
	r.status = ROOM_STATUS_PREPARE
	r.playerCnt = 0
	r.group = g
	return r
}

func(room *Room)FindPlayer(name string) (*Player, bool) {
	for _, player := range room.players {
		if player.playerData.name == name {
			return player, true
		}
	}
	return nil, false
}


func(r *Room)AddPlayer(p *Player) bool {
	if r.playerCnt >= ROOM_SIZE {
		return false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := 0; i < ROOM_SIZE; i++ {
		if r.players[i] == nil {
			r.playerCnt++
			p.extraPayerData.status = ROOM
			p.extraPayerData.team = r.SwitchTeam()
			r.players[i] = p
			r.group.Add(p.agent)
			p.extraPayerData.room = r
			return true
		}
	}
	log.Println("room.players has a error", r.players)
	return false
}

const (
	NOONE = 0
	TEAM_RED = 1
	TEAM_BLUE = 2
)

func(r *Room) SwitchTeam() int{
	cnt1,cnt2 := r.cnts()
	if cnt1 < cnt2 {
		return TEAM_RED
	}
	return TEAM_BLUE
}


func(r *Room)DelPlayer(p *Player) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := 0; i < ROOM_SIZE; i++ {
		if r.players[i] == p {
			p.extraPayerData.status = NONE
			r.players[i] = nil
			r.playerCnt--
			r.group.Del(p.agent)
			if p.extraPayerData.isOwner {
					p.extraPayerData.isOwner = false
					r.SwitchOwner()
			}
			return true
		}
	}
	return false
}

func(r *Room)SwitchOwner() {
	if r.playerCnt <= 0 {
		return
	}
	for i := 0; i < ROOM_SIZE; i++ {
		if r.players[i] != nil {
			r.players[i].extraPayerData.isOwner = true
			return
		}
	}
}

func(r *Room)Broadcast(p *proto.ProtocolBytes) {
	r.group.Broadcast(p.GetBuf())
}

func(r *Room)GetRoomInfo() *proto.ProtocolBytes{
	p := proto.NewProtocolBytes([]byte{0,0,0,0})
	p.EncodeString("GetRoomInfo")
	p.EncodeInt32(int32(r.playerCnt))
	for i := 0; i < ROOM_SIZE; i++ {
		if r.players[i] != nil {
			p.EncodeString(r.players[i].playerData.name)
			p.EncodeInt32(int32(r.players[i].extraPayerData.team))
			p.EncodeInt32(int32(r.players[i].playerData.win))
			p.EncodeInt32(int32(r.players[i].playerData.fail))
			if r.players[i].extraPayerData.isOwner {
				p.EncodeInt32(int32(1))
			} else {
				p.EncodeInt32(int32(0))
			}

			// 追加 player 是否准备好信息
			if r.players[i].extraPayerData.status == PREPARE {
				p.EncodeInt32(1)
			} else {
				p.EncodeInt32(0)
			}

		}
	}
	p.SetLength()
	return p
}

func(r *Room)cnts() (int, int){
	cnt1,cnt2 := 0, 0
	for _, p := range r.players {
		if p != nil {
			if p.extraPayerData.team == TEAM_RED {
				cnt1++
			} else {
				cnt2++
			}
		}
	}
	return cnt1, cnt2
}


func(r *Room)CanStart() bool{
	if r.status == ROME_STATUS_FIGHT {
		return false
	}
	cnt1,cnt2 := r.cnts()
	if cnt1 < 1 || cnt2 < 1 {
		return false
	}

	for _, player := range r.players {
		if player == nil {
			continue
		}
		if player.extraPayerData.isOwner {
			continue
		}
		// 如果有玩家没有准备好，或者他不是房主，那么就返回false
		if player.extraPayerData.status != PREPARE  {
			return false
		}
	}
	return true
}

func(r *Room)StartFight() {
	p := proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("Fight")
	r.status = ROME_STATUS_FIGHT
	//TODO i dont know whats means of it
	teamPos1, teamPos2 := 1, 1
	r.mu.Lock()
	defer r.mu.Unlock()
	p.EncodeInt32(int32(r.playerCnt))
	for _, player := range r.players {
		if player == nil {
			continue
		}
		player.extraPayerData.tankData.hp = 200
		p.EncodeString(player.playerData.name)
		p.EncodeInt32(int32(player.extraPayerData.team))
		if player.extraPayerData.team == TEAM_RED {
			teamPos1++
			p.EncodeInt32(int32(teamPos1))
		} else {
			teamPos2++
			p.EncodeInt32(int32(teamPos2))
		}
		player.extraPayerData.status = FIGHT
	}
	p.SetLength()
	r.Broadcast(p)
}


func(r *Room)IsWin() int {
	if r.status != ROME_STATUS_FIGHT {
		return NOONE
	}
	cnt1, cnt2 := 0, 0
	for _, p := range r.players {
		if p != nil {
			if p.extraPayerData.team == TEAM_RED && p.extraPayerData.tankData.hp >0{
				cnt1++
			}
			if p.extraPayerData.team == TEAM_BLUE && p.extraPayerData.tankData.hp >0{
				cnt2++
			}
		}
	}
	log.Println("cnt1 : ", cnt1, " , cnt2 : ", cnt2)
	if cnt1 <= 0 {
		return TEAM_BLUE
	} else if cnt2 <= 0 {
		return TEAM_RED
	} else {
		return NONE
	}
}

func(r *Room)UpdateWin() {
	whichTeam := r.IsWin()
	if whichTeam == NOONE {
		log.Println("no one team was win")
		return
	}
	r.mu.Lock()
	for _, player := range r.players {
		if player == nil {
			continue
		}
		player.extraPayerData.status = ROOM
		if player.extraPayerData.team == whichTeam {
			player.playerData.win++
		} else {
			player.playerData.fail++
		}
	}
	r.mu.Unlock()
	r.status = ROOM_STATUS_PREPARE
	p := proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("Result")
	p.EncodeInt32(int32(whichTeam))
	p.SetLength()
	r.Broadcast(p)
	GetGLobby().Broadcast()
}

func(r *Room) ExitFight(player *Player) {
	player.extraPayerData.tankData.hp = -1
	p := proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("Hit")
	p.EncodeString(player.playerData.name)
	p.EncodeString(player.playerData.name)
	p.EncodeFloat32(999.0)
	p.SetLength()
	r.Broadcast(p)
	if r.IsWin() == NOONE {
		player.playerData.fail++
	}
	r.UpdateWin()
}



