package main

import (
	"TankDemo/game"
	"TankDemo/network"
	"TankDemo/rpc"
	"fmt"
	"log"
	"net"
	"reflect"
)


func close(a *network.Agent) {
	game.GetGLobby().DelPlayer(a)
}

func main() {
	network.CloseHandle = close
	rpc.GetGRpcMap().Init()
	rpc.GetGRpcMap().AddActionHandle("HeatBeat", []reflect.Kind{},game.HeartBeat)
	rpc.GetGRpcMap().AddActionHandle("Login", []reflect.Kind{reflect.String, reflect.String,}, game.Login)
	rpc.GetGRpcMap().AddActionHandle("Register", []reflect.Kind{reflect.String, reflect.String,}, game.Register)
	rpc.GetGRpcMap().AddActionHandle("Logout", []reflect.Kind{}, game.LogoutEx)
	rpc.GetGRpcMap().AddActionHandle("StartFight", []reflect.Kind{}, game.StartFight )
	rpc.GetGRpcMap().AddActionHandle("Shooting", []reflect.Kind{
		reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32,
	}, game.Shooting)
	rpc.GetGRpcMap().AddActionHandle( "UpdateUnitInfo", []reflect.Kind{
		reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32, reflect.Float32,
	}, game.UpdateUnitInfo)
	rpc.GetGRpcMap().AddActionHandle( "Hit", []reflect.Kind{reflect.String, reflect.Float32}, game.Hit)
	rpc.GetGRpcMap().AddActionHandle( "GetScore", []reflect.Kind{}, game.GetScore)
	rpc.GetGRpcMap().AddActionHandle( "AddScore", []reflect.Kind{}, game.AddScore)
	rpc.GetGRpcMap().AddActionHandle( "GetList", []reflect.Kind{}, game.GetList)
	rpc.GetGRpcMap().AddActionHandle( "UpdateInfo", []reflect.Kind{}, game.UpdateInfo)
	rpc.GetGRpcMap().AddActionHandle( "GetAchieve", []reflect.Kind{}, game.GetAchieve)
	rpc.GetGRpcMap().AddActionHandle( "GetRoomList", []reflect.Kind{}, game.GetRoomList)
	rpc.GetGRpcMap().AddActionHandle( "CreateRoom", []reflect.Kind{}, game.CreateRoom)
	rpc.GetGRpcMap().AddActionHandle( "EnterRoom", []reflect.Kind{reflect.Int32, }, game.EnterRoom)
	rpc.GetGRpcMap().AddActionHandle( "GetRoomInfo", []reflect.Kind{}, game.GetRoomInfo)
	rpc.GetGRpcMap().AddActionHandle( "LeaveRoom", []reflect.Kind{}, game.LeaveRoom)
//	rpc.GetGRpcMap().AddActionHandle( "")


	_, find := rpc.GetGRpcMap().GetMethodDescriptor("Register")
	fmt.Println(find)

	game.InitGLobby()

	tcpConn, err := net.Listen("tcp", "0.0.0.0:18085")
	if err != nil {
		log.Println(err)
	}
	var uid = 0

	for {
		conn, err := tcpConn.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		a := network.NewAgent(int32(uid), conn)
		uid++
		fmt.Println(network.DumpAgent(a))
		a.Run(game.Process)
	}
}
