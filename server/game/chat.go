package game

import (
	"TankDemo/network"
	"TankDemo/proto"
	"log"
	"reflect"
)

//func(* Agent, []byte)([]byte, bool)

type chatAgentManager struct {
	agentToPlayer map[*network.Agent]*Player
}
var (
	gChatAgentManager *chatAgentManager
)
func InitChatManager() {
	gChatAgentManager = new(chatAgentManager)
	gChatAgentManager.agentToPlayer = make(map[*network.Agent]*Player)
}
func GetChatAgentManager() *chatAgentManager{
	return gChatAgentManager
}

func(cam *chatAgentManager)FindPlayer(a *network.Agent)(p *Player, ok bool) {
	p, ok = cam.agentToPlayer[a]
	return
}

func ChatProcess(ca *network.Agent, bytes []byte)([]byte, bool) {
	pro := proto.NewProtocolBytes(bytes)
	lenOfMsg := pro.DecodeInt32()
	if lenOfMsg > int32(len(bytes)) - 4 {
		log.Println("something error, because the buffer is too short")
	}
	log.Println(bytes)
	methodName := pro.DecodeString()
	if methodName == "HeatBeat" {
		return nil, true
	}
	md, found := proto.GetGRpcMap().GetMethodDescriptor(methodName)
	if !found {
		log.Println("chat process don't find this method, please try other.",ca , methodName)
		return []byte{1,1,1,1,1}, false
	}
	paramCnt := md.NumIn()
	var params []interface{}
	for i := 0; i < paramCnt; i++ {
		switch md.In(i) {
		case reflect.Int32:
			params = append(params, pro.DecodeInt32())
		case reflect.Float32:
			params = append(params, pro.DecodeFloat32())
		case reflect.String:
			params = append(params, pro.DecodeString())
		default:
			log.Println("error: unknown method param type",md.In(i))
		}
	}
	var m = md.GetMethod().(func(a *network.Agent, params []interface{})(*proto.ProtocolBytes, bool))
	pb, isEmpty := m(ca, params)
	if isEmpty {
		return nil, isEmpty
	}
	return pb.GetBuf(), isEmpty

}

func(cam *chatAgentManager) AddPlayer(a *network.Agent, p *Player) {
	cam.agentToPlayer[a] = p
}
func(cam *chatAgentManager) DelPlayer(a *network.Agent) {
	delete(cam.agentToPlayer, a)
}

func ParseChatName(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	name := params[0].(string)
	player, ok := GetGLobby().FindPlayerByName(name)
	if !ok {
		GetGLobby().LoginChatChan(name, a)
		<- a.ExtraChan
	}
	player, ok = GetGLobby().FindPlayerByName(name)
	if !ok {
		log.Println("not map a player. ", name)
	}
	GetChatAgentManager().AddPlayer(a, player)
	player.AddChatChannel(a)
	return nil, true
}

func ChatAnotherPort(a *network.Agent, params []interface{}) (p *proto.ProtocolBytes, isEmpty bool) {
	p = proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("Chat")
	p.EncodeString(params[0].(string))
	player, ok := GetChatAgentManager().FindPlayer(a)
	if !ok {
		log.Println("player not find.", a)
		return nil, true
	}
	log.Println("team is : ", player.extraPayerData.team)
	p.EncodeInt32(int32(player.extraPayerData.team))
	player.extraPayerData.room.ChatBroadcast(p)
	return nil, true
}