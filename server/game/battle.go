package game

import (
	"TankDemo/network"
	"TankDemo/proto"
	"log"
)

// 部分方法实现的地方

func Prepare(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player and he want to fight....", a)
		return nil, true
	}
	player.extraPayerData.status = PREPARE
	p = proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("Prepare")
	p.EncodeInt32(1)
	room := player.extraPayerData.room
	room.Broadcast(room.GetRoomInfo())
	return p, false
}

func Cancel(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player and he want to fight....", a)
		return nil, true
	}
	player.extraPayerData.status = ROOM
	p = proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("Cancel")
	p.EncodeInt32(0)
	room := player.extraPayerData.room
	room.Broadcast(room.GetRoomInfo())
	return p, false
}

func StartFight(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool){
	p = proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	isEmpty = false
	p.EncodeString("StartFight")
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	if player.extraPayerData.isOwner == false {
		p.EncodeInt32(-1)
		return
	}

	room := player.extraPayerData.room
	if room.CanStart() == false {
		p.EncodeInt32(-1)
		return
	}
	p.EncodeInt32(0)
	log.Println("start - fight succeed")
	room.StartFight()
	GetGLobby().Broadcast()
	return
}

func UpdateUnitInfo(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	if player.extraPayerData.status != FIGHT {
		return nil,true
	}
	player.extraPayerData.tankData.x = params[0].(float32)
	player.extraPayerData.tankData.y = params[1].(float32)
	player.extraPayerData.tankData.z = params[2].(float32)
	player.extraPayerData.updateCnt++
	p = proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("UpdateUnitInfo")
	p.EncodeString(player.playerData.name)
	for _, val := range params {
		p.EncodeFloat32(val.(float32))
	}
	player.extraPayerData.room.Broadcast(p)
	return nil, true
}

func Shooting(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	if player.extraPayerData.status != FIGHT {
		log.Println(player.id, "player not fight??")
		return nil,true
	}
	p = proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("Shooting")
	p.EncodeString(player.playerData.name)
	for _, val := range params {
		p.EncodeFloat32(val.(float32))
	}
	player.extraPayerData.room.Broadcast(p)
	log.Println("shooting will quit...")
	return nil, true
}

func Hit(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	if player.extraPayerData.status != FIGHT {
		log.Println(player.id, "player not fight??")
		return nil,true
	}
	room := player.extraPayerData.room
	enemy, ok := room.FindPlayer(params[0].(string))
	if !ok {
		log.Println("this name cant find a player", params[0].(string))
		return nil,true
	}
	if enemy.extraPayerData.tankData.hp <= 0 {
		return nil, true
	}
	log.Println("the hit point is : ", params[1].(float32))
	enemy.extraPayerData.tankData.hp -= int((params[1].(float32)))
	p = proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("Hit")
	p.EncodeString(player.playerData.name)
	p.EncodeString(enemy.playerData.name)
	p.EncodeFloat32(params[1].(float32))
	room.Broadcast(p)
	room.UpdateWin()
	log.Println("hit will quit...")
	return nil, true
}

func GetScore(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	p = proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("GetScore")
	p.EncodeInt32(int32(player.playerData.score))
	return p,false
}

func AddScore(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	player.playerData.score ++
	return nil,true
}

func GetList(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	p = proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("GetList")
	room := player.extraPayerData.room
	p.EncodeInt32(int32(room.playerCnt))
	for _, player := range room.players {
		if player != nil {
			p.EncodeString(player.playerData.name)
			p.EncodeFloat32(float32(player.extraPayerData.tankData.x))
			p.EncodeFloat32(float32(player.extraPayerData.tankData.y))
			p.EncodeFloat32(float32(player.extraPayerData.tankData.y))
			p.EncodeInt32(int32(player.playerData.score))
		}
	}
	return p, false
}

// maybe dont need. so wan can nop :), until the panic occurred
func UpdateInfo(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	return nil,true
}

func GetAchieve(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	p = proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("GetAchieve")
	p.EncodeInt32(int32(player.playerData.win))
	p.EncodeInt32(int32(player.playerData.fail))
	return p,false
}


func SwitchTeam(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player and he want to fight....", a)
		return nil, true
	}
	var res = -1
	p = proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("SwitchTeam")
	var team = int(params[0].(int32))
	if team == TEAM_RED || team == TEAM_BLUE {
		player.extraPayerData.team = team
		res = 0
		room := player.extraPayerData.room
		room.Broadcast(room.GetRoomInfo())
	}
	p.EncodeInt32(int32(res))
	return p, false
}

func GetRoomList(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	return GetGLobby().GetRoomList(), false
}


func CreateRoom(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	p = proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("CreateRoom")
	if player.extraPayerData.status != NONE {
		p.EncodeInt32(-1)
		return p,false
	}
	GetGLobby().CreateRoom(player)
	GetGLobby().Broadcast()

	p.EncodeInt32(0)
	return p,false
}


func EnterRoom(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	defer func() {
		GetGLobby().Broadcast()
	}()
	room := GetGLobby().roomList[params[0].(int32)]
	player, _ := GetGLobby().FindPlayer(a)		// TODO forget the ok
	p = proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("EnterRoom")
	if room.status != ROOM_STATUS_PREPARE {
		p.EncodeInt32(-1)

		return p, false
	}
	if room.AddPlayer(player) {
		p.EncodeInt32(0)
		log.Println("enter room :", RoomDump(room))
		room.Broadcast(room.GetRoomInfo())
		return p, false
	} else {
		p.EncodeInt32(-1)
		return p, false
	}
}

func GetRoomInfo(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	player, _ := GetGLobby().FindPlayer(a)		// forget the ok
	log.Println("lobby: player ", player)
	if player.extraPayerData.status != ROOM {
		return nil, true
	}
	room := player.extraPayerData.room
	p = room.GetRoomInfo()
	return p, false
}

func LeaveRoom(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	p = proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("LeaveRoom")
	player, _ := GetGLobby().FindPlayer(a)		// forget the ok
	if player.extraPayerData.status != ROOM {
		p.EncodeInt32(-1)
		log.Println("the player : ", player.playerData.name, " status is not ROOM.\n", PlayerDump(player))
		return p, false
	}
	p.EncodeInt32(0)
	room := player.extraPayerData.room
	GetGLobby().LeaveRoom(player)

	if room != nil {
		room.Broadcast(room.GetRoomInfo())
	}
	GetGLobby().Broadcast()
	return p, false
}

//maybe delete the function
func Chat(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	p = proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("Chat")
	p.EncodeString(params[0].(string))
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("player not find.", a)
		return nil, true
	}
	log.Println("team is : ", player.extraPayerData.team)
	p.EncodeInt32(int32(player.extraPayerData.team))
//	player.extraPayerData.room.ChatBroadcast(p)

	player.extraPayerData.room.Broadcast(p)
	return nil, true
}