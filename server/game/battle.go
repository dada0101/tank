package game

import (
	"TankDemo/network"
	"TankDemo/rpc"
	"log"
)

func StartFight(a *network.Agent, params []interface{}) (p *rpc.ProtocolBytes, isEmpty bool){
	p = rpc.NewProtocolBytes([]byte{0, 0, 0, 0})
	isEmpty = false
	p.EncodeString("StartFight")
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}

	if player.extraPayerData.status != ROOM {
		p.EncodeInt32(-1)
		p.SetLength()
		return
	}
	if player.extraPayerData.isOwner == false {
		p.EncodeInt32(-1)
		p.SetLength()
		return
	}
	room := player.extraPayerData.room
	if room.CanStart() == false {
		p.EncodeInt32(-1)
		p.SetLength()
		return
	}
	p.EncodeInt32(0)
	p.SetLength()
	log.Println("start - fight succeed")
	room.StartFight()
	return
}

func UpdateUnitInfo(a *network.Agent, params []interface{}) (p *rpc.ProtocolBytes, isEmpty bool) {
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
	p = rpc.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("UpdateUnitInfo")
	p.EncodeString(player.playerData.name)
	for _, val := range params {
		p.EncodeFloat32(val.(float32))
	}
	p.SetLength()
	player.extraPayerData.room.Broadcast(p)
	return nil, true
}

func Shooting(a *network.Agent, params []interface{}) (p *rpc.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	if player.extraPayerData.status != FIGHT {
		log.Println(player.id, "player not fight??")
		return nil,true
	}
	p = rpc.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("Shooting")
	p.EncodeString(player.playerData.name)
	for _, val := range params {
		p.EncodeFloat32(val.(float32))
	}
	p.SetLength()
	player.extraPayerData.room.Broadcast(p)
	log.Println("shooting will quit...")
	return nil, true
}

func Hit(a *network.Agent, params []interface{}) (p *rpc.ProtocolBytes, isEmpty bool) {
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
	enemy.extraPayerData.tankData.hp -= 50
	p = rpc.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("Hit")
	p.EncodeString(player.playerData.name)
	p.EncodeString(enemy.playerData.name)
	p.EncodeFloat32(50.0)
	p.SetLength()
	room.Broadcast(p)
	room.UpdateWin()
	log.Println("hit will quit...")
	return nil, true
}


func GetScore(a *network.Agent, params []interface{}) (p *rpc.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	p = rpc.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("GetScore")
	p.EncodeInt32(int32(player.playerData.score))
	return p,false
}

func AddScore(a *network.Agent, params []interface{}) (p *rpc.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	player.playerData.score ++
	return nil,true
}

func GetList(a *network.Agent, params []interface{}) (p *rpc.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	p = rpc.NewProtocolBytes([]byte{0, 0, 0, 0})
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
func UpdateInfo(a *network.Agent, params []interface{}) (p *rpc.ProtocolBytes, isEmpty bool) {
	return nil,true
}

func GetAchieve(a *network.Agent, params []interface{}) (p *rpc.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	p = rpc.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("GetAchieve")
	p.EncodeInt32(int32(player.playerData.win))
	p.EncodeInt32(int32(player.playerData.fail))
	p.SetLength()
	return p,false
}


func GetRoomList(a *network.Agent, params []interface{}) (p *rpc.ProtocolBytes, isEmpty bool) {
	return GetGLobby().GetRoomList(), false
}


func CreateRoom(a *network.Agent, params []interface{}) (p *rpc.ProtocolBytes, isEmpty bool) {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	p = rpc.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("CreateRoom")
	if player.extraPayerData.status != NONE {
		p.EncodeInt32(-1)
		p.SetLength()
		return p,false
	}
	GetGLobby().CreateRoom(player)


	p.EncodeInt32(0)
	p.SetLength()
	return p,false
}


func EnterRoom(a *network.Agent, params []interface{}) (p *rpc.ProtocolBytes, isEmpty bool) {
	room := GetGLobby().roomList[params[0].(int32)]
	player, _ := GetGLobby().FindPlayer(a)		// TODO forget the ok
	p = rpc.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("EnterRoom")
	if room.status != ROOM_STATUS_PREPARE {
		p.EncodeInt32(-1)
		p.SetLength()
		return p, false
	}
	if room.AddPlayer(player) {
		p.EncodeInt32(0)
		p.SetLength()
		log.Println("enter room :", RoomDump(room))
		room.Broadcast(room.GetRoomInfo())
		return p, false
	} else {
		p.EncodeInt32(-1)
		p.SetLength()
		return p, false
	}
}

func GetRoomInfo(a *network.Agent, params []interface{}) (p *rpc.ProtocolBytes, isEmpty bool) {
	player, _ := GetGLobby().FindPlayer(a)		// forget the ok
	log.Println("lobby: player ", player)
	if player.extraPayerData.status != ROOM {
		return nil, true
	}
	room := player.extraPayerData.room
	p = room.GetRoomInfo()
	p.SetLength()
	return p, false
}

func LeaveRoom(a *network.Agent, params []interface{}) (p *rpc.ProtocolBytes, isEmpty bool) {
	p = rpc.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("LeaveRoom")
	player, _ := GetGLobby().FindPlayer(a)		// forget the ok
	if player.extraPayerData.status != ROOM {
		p.EncodeInt32(-1)
		p.SetLength()
		return p, false
	}
	p.EncodeInt32(0)
	room := player.extraPayerData.room
	GetGLobby().LeaveRoom(player)

	if room != nil {
		room.Broadcast(room.GetRoomInfo())
	}
	p.SetLength()
	return p, false
}