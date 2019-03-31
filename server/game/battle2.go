package game
// 文件没有用到
// this fie will to do for change the struct of rpc
// TODO ALL
import (
	"TankDemo/network"
	"TankDemo/proto2"
	"log"
)

type PreparePackage struct {}
type CancelPackage struct {}

func(pp *PreparePackage)Read(s *proto2.BufferStream) bool {return true}
func(pp *PreparePackage)Write(s *proto2.BufferStream) bool {return true}
func(pp *PreparePackage)Exec(a *network.Agent) *proto2.BufferStream {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player and he want to fight....", a)
		return nil
	}
	player.extraPayerData.status = PREPARE
	bs := proto2.NewBufferStream([]byte{0, 0, 0, 0})
	bs.EncodeString("Prepare")
	bs.EncodeInt32(1)		//TODO
	room := player.extraPayerData.room
	room.Broadcast(room.GetRoomInfo())
	return bs
}

func(cp *CancelPackage)Read(s *proto2.BufferStream) bool {return true}
func(cp *CancelPackage)Write(s *proto2.BufferStream) bool {return true}
func(cp *CancelPackage)Exec(a *network.Agent) *proto2.BufferStream {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player and he want to fight....", a)
		return nil
	}
	player.extraPayerData.status = ROOM
	bs := proto2.NewBufferStream([]byte{0, 0, 0, 0})
	bs.EncodeString("Cancel")
	bs.EncodeInt32(1)		//TODO
	room := player.extraPayerData.room
	room.Broadcast(room.GetRoomInfo())
	return bs
}


type StartFightPackage struct {}

func(p *StartFightPackage)Read(s *proto2.BufferStream) bool {return true}
func(p *StartFightPackage)Write(s *proto2.BufferStream) bool {return true}
func(p *StartFightPackage)Exec(a *network.Agent) *proto2.BufferStream {
	bs := proto2.NewBufferStream([]byte{0, 0, 0, 0})
	bs.EncodeString("StartFight")
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	if player.extraPayerData.isOwner == false {
		bs.EncodeInt32(FAILED)
		return bs
	}
	room := player.extraPayerData.room
	if ! room.CanStart() {
		bs.EncodeInt32(FAILED)
		return bs
	}
	bs.EncodeInt32(SUCCEED)
	log.Println("start - fight succeed")
	room.StartFight()
	GetGLobby().Broadcast()
	return bs
}

type UpdateUnitPackage struct {
	posX, posY, posZ float32
	rotX, rotY, rotZ float32
	gunRot, gunRoll float32
}

func(p *UpdateUnitPackage)Read(s *proto2.BufferStream)bool {
	p.posX = s.DecodeFloat32()
	p.posY = s.DecodeFloat32()
	p.posZ = s.DecodeFloat32()
	p.rotX = s.DecodeFloat32()
	p.rotY = s.DecodeFloat32()
	p.rotZ = s.DecodeFloat32()
	p.gunRot = s.DecodeFloat32()
	p.gunRoll = s.DecodeFloat32()
	return true
}
func(p *UpdateUnitPackage)Write(s *proto2.BufferStream)bool {
	s.EncodeFloat32(p.posX)
	s.EncodeFloat32(p.posY)
	s.EncodeFloat32(p.posZ)
	s.EncodeFloat32(p.rotX)
	s.EncodeFloat32(p.rotY)
	s.EncodeFloat32(p.rotZ)
	s.EncodeFloat32(p.gunRot)
	s.EncodeFloat32(p.gunRoll)
	return true
}
func(p *UpdateUnitPackage)Exec(a *network.Agent) *proto2.BufferStream {
	player, ok := GetGLobby().FindPlayer(a)
	if !ok {
		log.Println("this agent not map a player", a)
	}
	if player.extraPayerData.status != FIGHT {
		return nil
	}
	player.extraPayerData.tankData.x = p.posX
	player.extraPayerData.tankData.y = p.posY
	player.extraPayerData.tankData.z = p.posZ
	player.extraPayerData.updateCnt++
	bs := proto2.NewBufferStream([]byte{0, 0, 0, 0})
	bs.EncodeString("UpdateUnitInfo")
	bs.EncodeString(player.playerData.name)

	bs.EncodeInt32(int32(p.posX))
	bs.EncodeInt32(int32(p.posY))
	bs.EncodeInt32(int32(p.posZ))
	bs.EncodeInt32(int32(p.rotX))
	bs.EncodeInt32(int32(p.rotY))
	bs.EncodeInt32(int32(p.rotZ))
	bs.EncodeInt32(int32(p.gunRot))
	bs.EncodeInt32(int32(p.gunRoll))

//TODO	player.extraPayerData.room.Broadcast(bs)
	return nil
}


type ShootPackage struct {

}