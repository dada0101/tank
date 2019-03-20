package game

import (
	"TankDemo/db"
	"TankDemo/network"
	"TankDemo/proto"
	"fmt"
	"log"
)

type PlayerData struct {
	name string
	score, win, fail int
}

const (
	NONE = 0
	ROOM = 1
	FIGHT = 2
	PREPARE = 3
)

type TankData struct {
	x, y, z float32
	hp int
}

type ExtraPlayerData struct {
	status int
	room *Room
	isOwner bool
	updateCnt int32
	team int
	tankData TankData
}



type Player struct {
	id int32
	playerData PlayerData
	extraPayerData ExtraPlayerData

	agent *network.Agent
}

func PlayerDump(p *Player) string {
	str := "\nDump player\n"
	str += fmt.Sprintf("id :\t%d\n", p.id)
	str += fmt.Sprintf("playerData[name : %s, score %d, win %d, fail %d]\n", p.playerData.name, p.playerData.score, p.playerData.win, p.playerData.fail)
	str += fmt.Sprintf("extraPlayerData[status : %d, room : %v, isOwner : %v, updateCnt : %v, team : %v]\n",
		p.extraPayerData.status,
		p.extraPayerData.room,
		p.extraPayerData.isOwner,
		p.extraPayerData.updateCnt,
		p.extraPayerData.team)
	return str
}


func NewPlayer(name string, id, score, win, fail int, a *network.Agent) *Player {
	return &Player{
		int32(id),
		PlayerData{name, score, win, fail, },
		ExtraPlayerData{NONE,
										nil,
										false,
										0,
										0,
										TankData{0.0,0.0,0.0,0}},
		a,
	}
}

// 关于重复登陆，不允许下一个登陆
func(p *Player) KickOff() {

}

func(player *Player) Send(p * proto.ProtocolBytes) {
	player.agent.Send(p.GetBuf())
}

func(p *Player) Logout() bool {
	if p.extraPayerData.status == ROOM {
		GetGLobby().LeaveRoom(p)
		room := p.extraPayerData.room
		if room != nil {
			room.Broadcast(room.GetRoomInfo())
		}
	}

	if p.extraPayerData.status == FIGHT {
		room := p.extraPayerData.room
		room.ExitFight(p)
		GetGLobby().LeaveRoom(p)
	}

	if !db.SetUserData(int(p.id), p.playerData.score, p.playerData.win, p.playerData.fail) {
		log.Println("p.id : ", p.id, "save the user data error.")
		return false
	}
	return true
}