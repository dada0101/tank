package game

import (
	"TankDemo/db"
	"TankDemo/network"
	"TankDemo/proto"
	"log"
	"reflect"
	"sync"
)
/*
type Player struct {

}*/


type Scene struct {
	players []*Player
	rule Rule
}
type Rule struct {

}

type Lobby struct {
	mu sync.Mutex
	agentToPlayer map[*network.Agent]*Player
	uidToPlayer map[int]*Player
	roomList []*Room
//	chLeavePlayer chan *Player
}

func (l *Lobby)Broadcast() {
	for agent, _ := range l.agentToPlayer {
		if err := agent.Send(l.GetRoomList().GetBuf()); err != nil {
			log.Println(err)
		}
	}
}


func(lobby *Lobby)CreateRoom(p* Player) {
	room := NewRoom(network.NewGroup(0))
	room.AddPlayer(p)
	room.SwitchOwner()
	p.extraPayerData.room = room
	lobby.roomList = append(lobby.roomList, room)
	log.Println(RoomDump(room))
}

func(lobby *Lobby)LeaveRoom(p *Player) {
//	if p.extraPayerData.status == NONE {
//		return
//	}
	room := p.extraPayerData.room
	room.DelPlayer(p)
	log.Println("room players count: ", room.playerCnt)
	if room.playerCnt == 0 {
		lobby.mu.Lock()
		var idx = 0
		var r *Room
		for idx, r = range lobby.roomList {
			if r == room {
				break
			}
		}
		copy(lobby.roomList[idx:], lobby.roomList[idx+1:])
		lobby.roomList = lobby.roomList[:len(lobby.roomList) -1]
		lobby.mu.Unlock()
	}
	log.Println(lobby.roomList)
}

func(lobby *Lobby)GetRoomList() *proto.ProtocolBytes{
	p := proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	p.EncodeString("GetRoomList")
	p.EncodeInt32(int32(len(lobby.roomList)))
	for _, room := range lobby.roomList {
		p.EncodeInt32(int32(room.playerCnt))
		p.EncodeInt32(int32(room.status))
	}
	p.SetLength()
	log.Println("lobby.roomlist: ", lobby.roomList)
	return p
}


var (
	gLobby *Lobby
)

func InitGLobby() {
	gLobby = new(Lobby)
	gLobby.agentToPlayer = make(map[*network.Agent]*Player)
	gLobby.uidToPlayer = make(map[int]*Player)
	gLobby.roomList = make([]*Room, 0)

}

func GetGLobby() *Lobby {
	return gLobby
}

func(l* Lobby)FindPlayer(a *network.Agent)(p *Player, ok bool) {
	p,ok = l.agentToPlayer[a]
	return
}
func(l *Lobby)FindPlayerByUid(uid int)(p *Player, ok bool ){
	p, ok = l.uidToPlayer[uid]
	return
}

func(l *Lobby)DelPlayer(a *network.Agent) {
	player, ok := l.FindPlayer(a)
	if ok {
		uid := player.id
		delete(l.uidToPlayer, int(uid))
		delete(l.agentToPlayer, a)
		player.Logout()
	}
}

// process the heartbeat
func HeartBeat(a *network.Agent, params []interface{}) (*proto.ProtocolBytes, bool){
	a.NoHBCntZero()
	return nil, true
}

type HeartBeatPackage struct {
}
func(hbp *HeartBeatPackage)Init([]byte) {

}
func(hbp *HeartBeatPackage)Decode() bool {
	return true
}
func(hbp *HeartBeatPackage)Encode() bool {
	return true
}
func(hbp *HeartBeatPackage)Exec(a *network.Agent) {
	a.NoHBCntZero()
}


func Register(a *network.Agent, params []interface{}) (*proto.ProtocolBytes, bool){
	 name := params[0].(string)
	 pwd := params[1].(string)
	pro := proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	pro.EncodeString("Register")
	uid, can := db.Register(name, pwd)
	if ! can {
		pro.EncodeInt32(-1)
		pro.SetLength()
		return pro, false
	}
	 pro.EncodeInt32(0)
	 pro.SetLength()
	 a.SetUid(int32(uid))
	log.Println("uid", uid)
	db.CreateUserData(uid)
	 return pro, false
}

type RegisterPackage struct {
	name, pwd string

}

func(rp *RegisterPackage)Init([]byte) {

}


func Login(a *network.Agent, params []interface{})(*proto.ProtocolBytes, bool) {
	name := params[0].(string)
	pwd := params[1].(string)
	pro := proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	pro.EncodeString("Login")
	uid, check := db.CheckPwd(name, pwd)
	if ! check {	// 验证失败
		pro.EncodeInt32(-1)
		pro.SetLength()
		return pro,false
	}
	_, ok := GetGLobby().FindPlayerByUid(uid)
	if ok {		// 如果重复登陆，拒绝
		pro.EncodeInt32(-1)
		pro.SetLength()
		return pro,false
	}
	score, win, fail, err :=db.GetUserData(uid)
	if err != nil {
		log.Println("some error about login")
	}
	player := NewPlayer(name, uid, score, win, fail, a)
	GetGLobby().agentToPlayer[a] = player
	GetGLobby().uidToPlayer[uid] = player

	pro.EncodeInt32(0)
	pro.SetLength()
	return pro, false
}



func LogoutEx(a *network.Agent, params []interface{}) (*proto.ProtocolBytes, bool) {
	pb := proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	pb.EncodeString("Logout")
	pb.EncodeInt32(0)

	player, ok := GetGLobby().agentToPlayer[a]
	if !ok {
		log.Println("this agent dont map a player and it want to logout", a)
		defer func() {
			a.Close()
		}()
		pb.SetLength()
		return pb, true
	}
	uid := player.id
	delete(GetGLobby().uidToPlayer, int(uid))
	delete(GetGLobby().agentToPlayer, a)
	player.Logout()
	pb.SetLength()
	return pb, false
}


func Process(a *network.Agent, bytes []byte) ([]byte,bool) {
	pro := proto.NewProtocolBytes(bytes)
	lenOfMsg := pro.DecodeInt32()
	if lenOfMsg > int32(len(bytes)) - 4 {
		log.Println("something error, because the buffer is too short")
	}
	methodName := pro.DecodeString()
	md, found := proto.GetGRpcMap().GetMethodDescriptor(methodName)
	if !found {
		log.Println("don't find this method, please try other.",a , methodName)
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

	log.Println(methodName, params)
	var m = md.GetMethod().(func(a *network.Agent, params []interface{})(*proto.ProtocolBytes, bool))
	pb, isEmpty := m(a, params)
	if isEmpty {
		return nil, isEmpty
	}

	return pb.GetBuf(), isEmpty
}

