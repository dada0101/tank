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

// 玩家状态
const (
	NONE = 0		//  无
	ROOM = 1		// 在房间中
	FIGHT = 2		// 在战斗中
	PREPARE = 3		// 准备状态
)

type TankData struct {
	x, y, z float32
	hp int
}

type ExtraPlayerData struct {
	status int			// 玩家状态
	room *Room			// 所属房间
	isOwner bool		// 是否是房主
	updateCnt int32		// 修改计数
	team int			// 所属队伍
	tankData TankData	// 坦克数据
}



type Player struct {
	id int32						// 玩家id
	playerData PlayerData			// tank数据
	extraPayerData ExtraPlayerData	// 玩家扩展数据

	agent *network.Agent			// 游戏逻辑代理
	chatAgent *network.Agent		// 聊天代理

	agentList map[*network.Agent]bool	// TODO  暂时没有用到 一直为空

	loginCnt int		// 登陆计数
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


func NewPlayer(name string, id, score, win, fail int, gameAgent, chatAgent *network.Agent, loginCnt int) *Player {
	return &Player{
		int32(id),
		PlayerData{name, score, win, fail, },
		ExtraPlayerData{NONE,
										nil,
										false,
										0,
										0,
										TankData{0.0,0.0,0.0,0}},
		gameAgent,
		chatAgent,
		make(map[*network.Agent]bool),
		loginCnt,
	}
}

// TODO 关于重复登陆，不允许下一个登陆
func(p *Player) KickOff() {

}

func(player *Player) Send(p * proto.ProtocolBytes) {
	player.agent.Send(p.GetBuf())
	for a, _ := range player.agentList {
		a.Send(p.GetBuf())
	}
}

func(player *Player) Hello(p * proto.ProtocolBytes) {
	player.chatAgent.Send(p.GetBuf())
}

func(player *Player) AddSpectator(a *network.Agent) {
	player.agentList[a] = true
}
func(player *Player) DelSpectator(a *network.Agent) {
	delete(player.agentList, a)
}

func(p *Player) Logout() bool {
	p.loginCnt++
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
	if !db.SetUserData(int(p.id), p.playerData.score, p.playerData.win, p.playerData.fail, p.loginCnt) {
		log.Println("p.id : ", p.id, "save the user data error.")
		return false
	}
	return true
}

func(p *Player)AddChatChannel(chatAgent *network.Agent) {
	p.chatAgent = chatAgent
}

func(p *Player)GetRoom() *Room {
	return p.extraPayerData.room
}

func(p *Player)GetTeam() int {
	return p.extraPayerData.team
}