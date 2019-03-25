package chat

import (
	"TankDemo/game"
	"TankDemo/network"
	"TankDemo/proto"
	"log"
)

//func(* Agent, []byte)([]byte, bool)

type chatAgentManager struct {
	agentToPlayer map[*network.Agent]*game.Player
}
var (
	gChatAgentManager *chatAgentManager
)
func Init() {
	gChatAgentManager = new(chatAgentManager)
	gChatAgentManager.agentToPlayer = make(map[*network.Agent]*game.Player)
}
func GetChatAgentManager() *chatAgentManager{
	return gChatAgentManager
}

func(cam *chatAgentManager)FindPlayer(a *network.Agent)(p *game.Player, ok bool) {
	p, ok = cam.agentToPlayer[a]
	return
}

func Process(ca *network.Agent, bytes []byte)([]byte, bool) {
	p, ok := GetChatAgentManager().FindPlayer(ca)
	if !ok {
		log.Println("the chat agent not map a player, please try other again.", ca)
	}
	pro := proto.NewProtocolBytes(bytes)
	lenOfMsg := pro.DecodeInt32()
	if lenOfMsg > int32(len(bytes)) - 4 {
		log.Println("something error, because the buffer is too short")
	}
	methodName := pro.DecodeString()
	if methodName != "Chat" {
		log.Println("don't find this method, please try other.",ca , methodName)
		return []byte{1,1,1,1,1}, false
	}
	room := p.GetRoom()
	room.ChatBroadcast(pro)
	return nil, true
}


//func Chan(ca *network.Agent, )