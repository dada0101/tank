package game

import (
	"TankDemo/db"
	"TankDemo/network"
	"TankDemo/proto"
	"log"
	"reflect"
	"sync"
)

type Lobby struct {
	mu sync.Mutex
	agentToPlayer map[*network.Agent]*Player		// 代理到玩家的映射
	roomList []*Room								// 所有房间

	nameToPlayer map[string]*Player					// 玩家姓名到玩家的映射
//	chLeavePlayer chan *Player
	nameToChatAgnetChan map[string]*chan struct{}	// 姓名到chatAgentChan 的映射，用于同步
}

func (l *Lobby)Broadcast() {
	for agent, _ := range l.agentToPlayer {
		if err := agent.Send(l.GetRoomList().GetBuf()); err != nil {
			log.Println(err)
		}
	}
}


func(lobby *Lobby)CreateRoom(p* Player) {
	room := NewRoom()
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
	log.Println("lobby.roomlist: ", lobby.roomList)
	return p
}


var (
	gLobby *Lobby
)

func InitGLobby() {
	gLobby = new(Lobby)
	gLobby.agentToPlayer = make(map[*network.Agent]*Player)
	gLobby.nameToPlayer = make(map[string]*Player)
	gLobby.roomList = make([]*Room, 0)
	gLobby.nameToChatAgnetChan = make(map[string]*chan struct{}, 0)
}

func(l *Lobby) LoginChatChan(name string, a *network.Agent) {
	l.nameToChatAgnetChan[name] = &a.ExtraChan
}
func(l *Lobby) LogoutChatChan(name string) {
	delete(l.nameToChatAgnetChan, name)
}

func GetGLobby() *Lobby {
	return gLobby
}

func(l* Lobby)FindPlayer(a *network.Agent)(p *Player, ok bool) {
	p,ok = l.agentToPlayer[a]
	return
}


func(l *Lobby)FindPlayerByName(name string)(p *Player, ok bool) {
	p, ok = l.nameToPlayer[name]
	return
}

func(l *Lobby)DelPlayer(a *network.Agent) {
	player, ok := l.FindPlayer(a)
	if ok {
		name := player.playerData.name
		delete(l.nameToPlayer, name)
		delete(l.agentToPlayer, a)
		player.Logout()
	}
}

// process the heartbeat
func HeartBeat(a *network.Agent, params []interface{}) (*proto.ProtocolBytes, bool){
	a.NoHBCntZero()
	return nil, true
}
/*
type HeartBeatPackage struct {
}
func(hbp *HeartBeatPackage)Read(stream *proto2.BufferStream) bool {
	return true
}
func(hbp *HeartBeatPackage)Write(stream *proto2.BufferStream) bool {
	return true
}
func(hbp *HeartBeatPackage)Exec(a *network.Agent)  *proto2.BufferStream {
	a.NoHBCntZero()
	return nil
}
*/

func Register(a *network.Agent, params []interface{}) (*proto.ProtocolBytes, bool){
	 name := params[0].(string)
	 pwd := params[1].(string)
	pro := proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	pro.EncodeString("Register")
	uid, can := db.Register(name, pwd)
	if ! can {
		pro.EncodeInt32(-1)
		return pro, false
	}
	 pro.EncodeInt32(0)
	 a.SetUid(int32(uid))
	log.Println("uid", uid)
	db.CreateUserData(uid)
	 return pro, false
}
const (
	FAILED = -1
	SUCCEED = 0
)
/*
type RegisterPackage struct {
	name, pwd string
}

func(rp *RegisterPackage)Read(stream *proto2.BufferStream) bool {
	rp.name = stream.DecodeString()
	rp.pwd = stream.DecodeString()
	return true
}

func(rp *RegisterPackage)Write(stream *proto2.BufferStream) bool {
	stream.EncodeString(rp.name)
	stream.EncodeString(rp.pwd)
	return true
}

func(rp *RegisterPackage)Exec(a *network.Agent) *proto2.BufferStream {
	 bs := proto2.NewBufferStream([]byte{0, 0, 0, 0})
	 bs.EncodeString("Register")
	uid, can := db.Register(rp.name, rp.pwd)
	if ! can {
		bs.EncodeInt32(FAILED)
		return bs
	}
	bs.EncodeInt32(SUCCEED)
	a.SetUid(int32(uid))
	log.Println("uid", uid)
	db.CreateUserData(uid)
	return bs
}
*/
func Login(a *network.Agent, params []interface{})(*proto.ProtocolBytes, bool) {
	name := params[0].(string)
	pwd := params[1].(string)
	pro := proto.NewProtocolBytes([]byte{0, 0, 0, 0})
	pro.EncodeString("Login")
	uid, check, loginCnt := db.CheckPwd(name, pwd)
	if ! check {	// 验证失败
		pro.EncodeInt32(-1)
		return pro,false
	}
	_, ok := GetGLobby().FindPlayerByName(name)
	if ok {		// 如果重复登陆，拒绝
		pro.EncodeInt32(-1)
		return pro,false
	}
	score, win, fail, err :=db.GetUserData(uid)
	if err != nil {
		log.Println("some error about login")
	}
	player := NewPlayer(name, uid, score, win, fail, a, nil, loginCnt)
	GetGLobby().agentToPlayer[a] = player
	GetGLobby().nameToPlayer[name] = player

	pro.EncodeInt32(0)
	if loginCnt == 0 {
		pro.EncodeInt32(int32(1)) // 备注：这里的 1 代表第一次登陆 也就是 新手
	}	else {
		pro.EncodeInt32(0)
	}
	var chanRef *chan struct{}
	chanRef, ok = GetGLobby().nameToChatAgnetChan[name]
	if ok {
		(*chanRef) <- struct{}{}
	}
	return pro, false
}
/*
type LoginPackage struct {
	name, pwd string
}

func(lp *LoginPackage)Read(stream *proto2.BufferStream) bool {
	lp.name = stream.DecodeString()
	lp.pwd = stream.DecodeString()
	return true
}
func(lp *LoginPackage)Write(stream *proto2.BufferStream) bool {
	stream.EncodeString(lp.name)
	stream.EncodeString(lp.pwd)
	return true
}
func(lp *LoginPackage)Exec(a *network.Agent) *proto2.BufferStream {
	bf := proto2.NewBufferStream([]byte{0, 0, 0, 0})
	bf.EncodeString("Login")
	uid, check, _ := db.CheckPwd(lp.name, lp.pwd)
	if ! check {	// 验证失败
		bf.EncodeInt32(FAILED)
		return bf
	}
	_, ok := GetGLobby().FindPlayerByName(lp.name)
	if ok {		// 如果重复登陆，拒绝
		bf.EncodeInt32(FAILED)
		return bf
	}
	score, win, fail, err :=db.GetUserData(uid)
	if err != nil {
		log.Println("some error about login")
	}
	player := NewPlayer(lp.name, uid, score, win, fail, a, nil, 0)
	GetGLobby().agentToPlayer[a] = player
	GetGLobby().nameToPlayer[lp.name] = player

	bf.EncodeInt32(SUCCEED)
	return bf
}
*/

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
		return pb, true
	}
	delete(GetGLobby().nameToPlayer, player.playerData.name)
	delete(GetGLobby().agentToPlayer, a)
	GetChatAgentManager().DelPlayer(player.chatAgent)
	player.Logout()
//	player.chatAgent.Close()
	return pb, false
}
/*
type LogoutPackage struct {
}

func(lp *LogoutPackage)Read(stream *proto2.BufferStream) bool {
	return true
}
func(lp *LogoutPackage)Write(stream *proto2.BufferStream) bool {
	return true
}
func(lp *LogoutPackage)Exec(a *network.Agent) *proto2.BufferStream {
	bs := proto2.NewBufferStream([]byte{0, 0, 0, 0})
	bs.EncodeString("Logout")
	bs.EncodeInt32(SUCCEED)
	player, ok := GetGLobby().agentToPlayer[a]
	if !ok {
		log.Println("this agent dont map a player and it want to logout", a)
		defer func() {
			a.Close()
		}()
		return bs
	}
	delete(GetGLobby().nameToPlayer, player.playerData.name)
	delete(GetGLobby().agentToPlayer, a)
	player.Logout()
	return bs
}

*/
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

//	log.Println(methodName, params)
	var m = md.GetMethod().(func(a *network.Agent, params []interface{})(*proto.ProtocolBytes, bool))
	pb, isEmpty := m(a, params)
	if isEmpty {
		return nil, isEmpty
	}

	return pb.GetBuf(), isEmpty
}

