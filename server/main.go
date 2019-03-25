package main

import (
	"TankDemo/game"
	"TankDemo/network"
	"TankDemo/proto"
	"log"
	"net"
	"reflect"
	"time"
)


func close(a *network.Agent) {
	game.GetGLobby().DelPlayer(a)
}

func Init() {
	proto.GetGRpcMap().Init()
	proto.GetGRpcMap().AddActionHandle("HeatBeat", []reflect.Kind{},game.HeartBeat)
	proto.GetGRpcMap().AddActionHandle("Login", []reflect.Kind{reflect.String, reflect.String,}, game.Login)
	proto.GetGRpcMap().AddActionHandle("Register", []reflect.Kind{reflect.String, reflect.String,}, game.Register)
	proto.GetGRpcMap().AddActionHandle("Logout", []reflect.Kind{}, game.LogoutEx)
	proto.GetGRpcMap().AddActionHandle("Prepare", []reflect.Kind{}, game.Prepare )
	proto.GetGRpcMap().AddActionHandle("Cancel", []reflect.Kind{}, game.Cancel )
	proto.GetGRpcMap().AddActionHandle("StartFight", []reflect.Kind{}, game.StartFight )
	proto.GetGRpcMap().AddActionHandle("Shooting", []reflect.Kind{
		reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32,
	}, game.Shooting)
	proto.GetGRpcMap().AddActionHandle( "UpdateUnitInfo", []reflect.Kind{
		reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32,
	}, game.UpdateUnitInfo)
	proto.GetGRpcMap().AddActionHandle( "Hit", []reflect.Kind{reflect.String, reflect.Float32}, game.Hit)
	proto.GetGRpcMap().AddActionHandle( "GetScore", []reflect.Kind{}, game.GetScore)
	proto.GetGRpcMap().AddActionHandle( "AddScore", []reflect.Kind{}, game.AddScore)
	proto.GetGRpcMap().AddActionHandle( "GetList", []reflect.Kind{}, game.GetList)
	proto.GetGRpcMap().AddActionHandle( "UpdateInfo", []reflect.Kind{}, game.UpdateInfo)
	proto.GetGRpcMap().AddActionHandle( "GetAchieve", []reflect.Kind{}, game.GetAchieve)
	proto.GetGRpcMap().AddActionHandle( "GetRoomList", []reflect.Kind{}, game.GetRoomList)
	proto.GetGRpcMap().AddActionHandle( "CreateRoom", []reflect.Kind{}, game.CreateRoom)
	proto.GetGRpcMap().AddActionHandle( "EnterRoom", []reflect.Kind{reflect.Int32, }, game.EnterRoom)
	proto.GetGRpcMap().AddActionHandle( "GetRoomInfo", []reflect.Kind{}, game.GetRoomInfo)
	proto.GetGRpcMap().AddActionHandle( "LeaveRoom", []reflect.Kind{}, game.LeaveRoom)
	proto.GetGRpcMap().AddActionHandle( "SwitchTeam", []reflect.Kind{reflect.Int32, }, game.SwitchTeam)

	proto.GetGRpcMap().AddActionHandle("Chat", []reflect.Kind{reflect.String,}, game.ChatAnotherPort)
	proto.GetGRpcMap().AddActionHandle("ChatName", []reflect.Kind{reflect.String, }, game.ParseChatName)
}

func main() {
	Init()
	game.InitGLobby()
	network.CloseHandle = close

//	logFile, err := os.OpenFile("log.txt", os.O_CREATE | os.O_WRONLY | os.O_APPEND, 0666)
//	if err != nil {
//		log.Fatalln("log file can not to open!! ", err)
//	}
//	log.SetOutput(logFile)
//	defer func () {
//		logFile.Close()
//	} ()
	tcpConn, err := net.Listen("tcp", "0.0.0.0:18085")
	if err != nil {
		log.Println(err)
	}
	var uid = 10000
	go chat()
	for {
		conn, err := tcpConn.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		a := network.NewAgent(int32(uid), conn)
		uid++
		a.Run(game.Process, 30 * time.Second)
	}

}

func chat() {
	game.InitChatManager()
	chatConn, err := net.Listen("tcp", "0.0.0.0:18086")
	if err != nil {
		log.Println(err)
	}
	var cid = 0
	for {
		conn, err := chatConn.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		a := network.NewAgent(int32(cid), conn)
		cid++
		a.Run(game.ChatProcess, 100 * time.Hour)
	}
}
//ParseChatName

